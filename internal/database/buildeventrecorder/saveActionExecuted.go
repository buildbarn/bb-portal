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

// getRemoteFileURI returns a URI that can be resolved away from the Bazel
// client. Local file:// references are deliberately not persisted.
func getRemoteFileURI(file *bes.File) string {
	if file == nil {
		return ""
	}
	uri := strings.TrimSpace(file.GetUri())
	if uri == "" {
		return ""
	}
	parsedURI, err := url.Parse(uri)
	if err != nil || parsedURI.Scheme == "" || strings.EqualFold(parsedURI.Scheme, "file") {
		return ""
	}
	return uri
}

type actionFileAvailability struct {
	verifyBytestreamURIs bool
	verified             bool
	missingDigests       map[storagedigest.Digest]struct{}
}

func getBytestreamDigest(uri string) (storagedigest.Digest, bool) {
	parsedURI, err := url.Parse(uri)
	if err != nil || !strings.EqualFold(parsedURI.Scheme, "bytestream") {
		return storagedigest.BadDigest, false
	}
	return files.GetDigestFromURI(uri), true
}

func (a actionFileAvailability) getAvailableRemoteFileURI(file *bes.File) string {
	uri := getRemoteFileURI(file)
	if uri == "" {
		return ""
	}
	digest, isBytestreamURI := getBytestreamDigest(uri)
	if !isBytestreamURI {
		return uri
	}
	if digest == storagedigest.BadDigest {
		return ""
	}
	if !a.verifyBytestreamURIs {
		return uri
	}
	if !a.verified {
		return ""
	}
	if _, missing := a.missingDigests[digest]; missing {
		return ""
	}
	return uri
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
			uri := getRemoteFileURI(file)
			if digest, isBytestreamURI := getBytestreamDigest(uri); isBytestreamURI && digest != storagedigest.BadDigest {
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

func getActionLabel(actionExecuted *bes.ActionExecuted, actionCompletedID *bes.BuildEventId_ActionCompletedId) string {
	if label := actionCompletedID.GetLabel(); label != "" {
		return label
	}
	return actionExecuted.GetLabel()
}

// saveActionExecutedBatch persists all ActionExecuted events in a batch. The
// configuration lookup is also batched, because --build_event_publish_all_actions
// can cause thousands of these events to be reported for one invocation.
func (r *buildEventRecorder) saveActionExecutedBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}
	actionFileAvailability := r.verifyActionFileAvailability(ctx, batch)

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

	creates := make([]*ent.ActionCreate, 0, len(batch))
	for _, info := range batch {
		actionExecuted := info.Event.GetAction()
		actionCompletedID := info.Event.GetId().GetActionCompleted()
		if actionExecuted == nil || actionCompletedID == nil {
			continue
		}

		create := tx.Ent().Action.Create().
			SetBazelInvocationID(r.InvocationDbID).
			SetLabel(getActionLabel(actionExecuted, actionCompletedID)).
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
		if primaryOutputURI := actionFileAvailability.getAvailableRemoteFileURI(actionExecuted.GetPrimaryOutput()); primaryOutputURI != "" {
			create.SetPrimaryOutputURI(primaryOutputURI)
		}
		if !actionExecuted.Success {
			if stdoutURI := actionFileAvailability.getAvailableRemoteFileURI(actionExecuted.GetStdout()); stdoutURI != "" {
				create.SetStdoutURI(stdoutURI)
			}
			if stderrURI := actionFileAvailability.getAvailableRemoteFileURI(actionExecuted.GetStderr()); stderrURI != "" {
				create.SetStderrURI(stderrURI)
			}
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
		if err := tx.Ent().Action.CreateBulk(creates...).Exec(ctx); err != nil {
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
