package buildeventrecorder

import (
	"context"
	"database/sql"
	"log/slog"
	"net/url"
	"reflect"
	"strings"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/bazelbuild/bazel/src/main/protobuf"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/configuration"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	storagedigest "github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorCodeFromFailureDetail(failureDetail *protobuf.FailureDetail) string {
	if failureDetail == nil || failureDetail.Category == nil {
		return ""
	}
	detailValue := reflect.ValueOf(failureDetail.Category)
	if detailValue.Kind() == reflect.Ptr {
		detailValue = detailValue.Elem()
	}
	if detailValue.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < detailValue.NumField(); i++ {
		fieldValue := detailValue.Field(i)
		if fieldValue.Kind() != reflect.Ptr {
			continue
		}

		method := fieldValue.MethodByName("GetCode")
		if !method.IsValid() {
			continue
		}

		result := method.Call(nil)
		if len(result) == 0 {
			continue
		}

		stringer := result[0].MethodByName("String")
		if !stringer.IsValid() {
			continue
		}
		return stringer.Call(nil)[0].String()
	}
	return ""
}

type actionFileAvailability struct {
	verifyBytestreamURIs bool
	verified             bool
	missingDigests       map[storagedigest.Digest]struct{}
}

func getBytestreamDigest(uri string) (storagedigest.Digest, bool) {
	uri = strings.TrimSpace(uri)
	parsedURI, err := url.Parse(uri)
	if err != nil || !strings.EqualFold(parsedURI.Scheme, "bytestream") {
		return storagedigest.BadDigest, false
	}
	return files.GetDigestFromURI(uri), true
}

func (a actionFileAvailability) getAvailableFile(file *bes.File, fallbackPath string) *files.ParsedBepFile {
	if file == nil {
		return nil
	}
	digest, isBytestreamURI := getBytestreamDigest(file.GetUri())
	if !isBytestreamURI {
		return nil
	}
	if digest == storagedigest.BadDigest {
		return nil
	}
	if !a.verifyBytestreamURIs {
		return files.ParseBepFileWithFallbackPath(file, fallbackPath)
	}
	if !a.verified {
		return nil
	}
	if _, missing := a.missingDigests[digest]; missing {
		return nil
	}
	return files.ParseBepFileWithFallbackPath(file, fallbackPath)
}

func actionExecutedFiles(actionExecuted *bes.ActionExecuted) []*bes.File {
	if actionExecuted == nil {
		return nil
	}
	return []*bes.File{
		actionExecuted.GetPrimaryOutput(),
		actionExecuted.GetStdout(),
		actionExecuted.GetStderr(),
	}
}

func (r *buildEventRecorder) verifyActionFileAvailability(ctx context.Context, batch []BuildEventWithInfo) actionFileAvailability {
	availability := actionFileAvailability{
		verifyBytestreamURIs: r.contentAddressableStorage != nil,
	}
	if !availability.verifyBytestreamURIs {
		return availability
	}

	digestSetBuilder := storagedigest.NewSetBuilder(len(batch) * 3)
	for _, info := range batch {
		for _, file := range actionExecutedFiles(info.Event.GetAction()) {
			if digest, isBytestreamURI := getBytestreamDigest(file.GetUri()); isBytestreamURI && digest != storagedigest.BadDigest {
				digestSetBuilder = digestSetBuilder.Add(digest)
			}
		}
	}
	digests := digestSetBuilder.Build()
	if digests.Empty() {
		availability.verified = true
		return availability
	}

	missingDigests, err := r.contentAddressableStorage.FindMissing(ctx, digests)
	if err != nil {
		slog.WarnContext(ctx, "Could not verify action file blobs in CAS; file links will be unavailable", "invocation_id", r.InvocationID, "err", err)
		return availability
	}

	availability.verified = true
	availability.missingDigests = make(map[storagedigest.Digest]struct{}, missingDigests.Length())
	for _, digest := range missingDigests.Items() {
		availability.missingDigests[digest] = struct{}{}
	}
	return availability
}

func getActionConfigurationID(actionExecuted *bes.ActionExecuted, actionCompletedID *bes.BuildEventId_ActionCompletedId) string {
	if configurationID := actionCompletedID.GetConfiguration().GetId(); configurationID != "" {
		return configurationID
	}
	return actionExecuted.GetConfiguration().GetId()
}

type actionExecutionFiles struct {
	primaryOutput *files.ParsedBepFile
	stdout        *files.ParsedBepFile
	stderr        *files.ParsedBepFile
}

func getAvailableActionFiles(availability actionFileAvailability, batch []BuildEventWithInfo) ([]actionExecutionFiles, []*files.ParsedBepFile) {
	filesByEvent := make([]actionExecutionFiles, len(batch))
	allFiles := make([]*files.ParsedBepFile, 0, len(batch)*3)
	for i, info := range batch {
		actionExecuted := info.Event.GetAction()
		actionCompletedID := info.Event.GetId().GetActionCompleted()
		if actionExecuted == nil || actionCompletedID == nil {
			continue
		}

		actionFiles := actionExecutionFiles{
			primaryOutput: availability.getAvailableFile(actionExecuted.GetPrimaryOutput(), actionCompletedID.GetPrimaryOutput()),
			stdout:        availability.getAvailableFile(actionExecuted.GetStdout(), "stdout"),
			stderr:        availability.getAvailableFile(actionExecuted.GetStderr(), "stderr"),
		}
		filesByEvent[i] = actionFiles
		for _, file := range []*files.ParsedBepFile{actionFiles.primaryOutput, actionFiles.stdout, actionFiles.stderr} {
			if file != nil {
				allFiles = append(allFiles, file)
			}
		}
	}
	return filesByEvent, allFiles
}

// saveActionExecutedBatch persists all ActionExecuted events in a batch. The
// configuration lookup is also batched, because --build_event_publish_all_actions
// can cause thousands of these events to be reported for one invocation.
func (r *buildEventRecorder) saveActionExecutedBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}
	actionFileAvailability := r.verifyActionFileAvailability(ctx, batch)
	actionFilesByEvent, actionFiles := getAvailableActionFiles(actionFileAvailability, batch)

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return util.StatusWrap(err, "Failed to start transaction")
	}
	defer tx.Rollback()

	row, err := tx.Sqlc().LockBazelInvocationCompletion(ctx, r.InvocationDbID)
	if err != nil {
		return util.StatusWrap(err, "Failed to lock bep completed for invocation")
	}
	if row.BepCompleted {
		return status.Error(codes.FailedPrecondition, "Attempted to insert action events but the invocation was already completed.")
	}

	configurationIDSet := make(map[string]struct{})
	for _, info := range batch {
		actionExecuted := info.Event.GetAction()
		actionCompletedID := info.Event.GetId().GetActionCompleted()
		if actionExecuted == nil || actionCompletedID == nil {
			continue
		}
		if configurationID := getActionConfigurationID(actionExecuted, actionCompletedID); configurationID != "" {
			configurationIDSet[configurationID] = struct{}{}
		}
	}

	configurationIDs := make([]string, 0, len(configurationIDSet))
	for configurationID := range configurationIDSet {
		configurationIDs = append(configurationIDs, configurationID)
	}

	configurationDBIDs := make(map[string]int64, len(configurationIDs))
	if len(configurationIDs) > 0 {
		configurations, err := tx.Ent().Configuration.Query().
			Where(
				configuration.BazelInvocationID(r.InvocationDbID),
				configuration.ConfigurationIDIn(configurationIDs...),
			).
			All(ctx)
		if err != nil {
			return util.StatusWrap(err, "Failed to query configurations for ActionExecuted events")
		}
		for _, config := range configurations {
			configurationDBIDs[config.ConfigurationID] = config.ID
		}
	}

	actionFileIDs, err := saveFilesBatch(ctx, tx, r.InstanceNameDbID, actionFiles)
	if err != nil {
		return util.StatusWrap(err, "Failed to save ActionExecuted output files")
	}
	actionFileIDByParsedFile := make(map[*files.ParsedBepFile]int64, len(actionFiles))
	for i, actionFile := range actionFiles {
		actionFileIDByParsedFile[actionFile] = actionFileIDs[i]
	}

	creates := make([]*ent.ActionExecutionCreate, 0, len(batch))
	for i, info := range batch {
		actionExecuted := info.Event.GetAction()
		actionCompletedID := info.Event.GetId().GetActionCompleted()
		if actionExecuted == nil || actionCompletedID == nil {
			continue
		}

		create := tx.Ent().ActionExecution.Create().
			SetBazelInvocationID(r.InvocationDbID).
			SetLabel(actionCompletedID.GetLabel()).
			SetSuccess(actionExecuted.Success).
			SetExitCode(actionExecuted.ExitCode)

		if configurationID := getActionConfigurationID(actionExecuted, actionCompletedID); configurationID != "" {
			if configurationDBID, ok := configurationDBIDs[configurationID]; ok {
				create.SetConfigurationID(configurationDBID)
			}
		}
		if actionExecuted.Type != "" {
			create.SetType(actionExecuted.Type)
		}
		if len(actionExecuted.CommandLine) > 0 {
			create.SetCommandLine(actionExecuted.CommandLine)
		}
		if primaryOutput := actionCompletedID.GetPrimaryOutput(); primaryOutput != "" {
			create.SetPrimaryOutput(primaryOutput)
		}
		if primaryOutputFile := actionFilesByEvent[i].primaryOutput; primaryOutputFile != nil {
			create.SetPrimaryOutputFileID(actionFileIDByParsedFile[primaryOutputFile])
		}
		if stdoutFile := actionFilesByEvent[i].stdout; stdoutFile != nil {
			create.SetStdoutID(actionFileIDByParsedFile[stdoutFile])
		}
		if stderrFile := actionFilesByEvent[i].stderr; stderrFile != nil {
			create.SetStderrID(actionFileIDByParsedFile[stderrFile])
		}
		if failureMessage := actionExecuted.GetFailureDetail().GetMessage(); failureMessage != "" {
			create.SetFailureMessage(failureMessage)
		}
		if failureCode := getErrorCodeFromFailureDetail(actionExecuted.GetFailureDetail()); failureCode != "" {
			create.SetFailureCode(failureCode)
		}
		if actionExecuted.StartTime != nil {
			create.SetStartTime(actionExecuted.StartTime.AsTime())
		}
		if actionExecuted.EndTime != nil {
			create.SetEndTime(actionExecuted.EndTime.AsTime())
		}
		creates = append(creates, create)
	}

	if len(creates) > 0 {
		if err := tx.Ent().ActionExecution.CreateBulk(creates...).Exec(ctx); err != nil {
			return util.StatusWrap(err, "Failed to save ActionExecuted events")
		}
	}
	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}
	if err := tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit batch of action events")
	}
	return nil
}
