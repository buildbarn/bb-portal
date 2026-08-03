package buildeventrecorder

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	bazelprotobuf "github.com/bazelbuild/bazel/src/main/protobuf"
	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/action"
	entdigest "github.com/buildbarn/bb-portal/ent/gen/ent/digest"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	storagedigest "github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/encoding/protodelim"
)

const (
	compactExecutionLogName  = "execution_log.binpb.zst"
	maxExecutionLogEntrySize = 64 << 20
	internalActionRunner     = "internal"
)

type executionLogAction struct {
	targetLabel  string
	mnemonic     string
	outputPaths  []string
	runner       string
	actionDigest storagedigest.Digest
}

func getCompactExecutionLog(buildToolLogs *bes.BuildToolLogs) *bes.File {
	if buildToolLogs == nil {
		return nil
	}
	for _, log := range buildToolLogs.GetLog() {
		if log.GetName() == compactExecutionLogName {
			return log
		}
	}
	return nil
}

func executionLogDigestFunction(instanceName storagedigest.InstanceName, hashFunctionName string, hashLength int) (storagedigest.Function, error) {
	normalizedName := strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToUpper(hashFunctionName))
	rawDigestFunction, ok := remoteexecution.DigestFunction_Value_value[normalizedName]
	if !ok {
		rawDigestFunction = int32(remoteexecution.DigestFunction_UNKNOWN)
	}
	return instanceName.GetDigestFunction(remoteexecution.DigestFunction_Value(rawDigestFunction), hashLength)
}

func parseCompactExecutionLog(reader io.Reader, instanceName storagedigest.InstanceName) ([]executionLogAction, error) {
	decoder, err := zstd.NewReader(reader)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create execution log decoder")
	}
	defer decoder.Close()

	bufferedReader := bufio.NewReader(decoder)
	unmarshalOptions := protodelim.UnmarshalOptions{MaxSize: maxExecutionLogEntrySize}
	pathsByID := map[uint32]string{}
	hashFunctionName := ""
	actions := []executionLogAction{}

	for {
		entry := &bazelprotobuf.ExecLogEntry{}
		if err := unmarshalOptions.UnmarshalFrom(bufferedReader, entry); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, util.StatusWrap(err, "Failed to parse compact execution log entry")
		}

		if invocation := entry.GetInvocation(); invocation != nil {
			hashFunctionName = invocation.GetHashFunctionName()
			continue
		}
		if file := entry.GetFile(); file != nil {
			pathsByID[entry.GetId()] = file.GetPath()
			continue
		}
		if directory := entry.GetDirectory(); directory != nil {
			pathsByID[entry.GetId()] = directory.GetPath()
			continue
		}
		if symlink := entry.GetUnresolvedSymlink(); symlink != nil {
			pathsByID[entry.GetId()] = symlink.GetPath()
			continue
		}

		if symlinkAction := entry.GetSymlinkAction(); symlinkAction != nil {
			if symlinkAction.GetOutputPath() != "" {
				actions = append(actions, executionLogAction{
					targetLabel: symlinkAction.GetTargetLabel(),
					mnemonic:    symlinkAction.GetMnemonic(),
					outputPaths: []string{symlinkAction.GetOutputPath()},
					runner:      internalActionRunner,
				})
			}
			continue
		}

		spawn := entry.GetSpawn()
		if spawn == nil {
			continue
		}

		outputPaths := make([]string, 0, len(spawn.GetOutputs()))
		seenOutputPaths := map[string]struct{}{}
		for _, output := range spawn.GetOutputs() {
			var outputPath string
			switch outputType := output.GetType().(type) {
			case *bazelprotobuf.ExecLogEntry_Output_OutputId:
				outputPath = pathsByID[outputType.OutputId]
			case *bazelprotobuf.ExecLogEntry_Output_InvalidOutputPath:
				outputPath = outputType.InvalidOutputPath
			}
			if outputPath == "" {
				continue
			}
			if _, ok := seenOutputPaths[outputPath]; ok {
				continue
			}
			seenOutputPaths[outputPath] = struct{}{}
			outputPaths = append(outputPaths, outputPath)
		}
		if len(outputPaths) == 0 {
			continue
		}

		executionLogAction := executionLogAction{
			targetLabel: spawn.GetTargetLabel(),
			mnemonic:    spawn.GetMnemonic(),
			outputPaths: outputPaths,
			runner:      spawn.GetRunner(),
		}
		if spawn.GetDigest() != nil && spawn.GetDigest().GetHash() != "" {
			spawnHashFunctionName := spawn.GetDigest().GetHashFunctionName()
			if spawnHashFunctionName == "" {
				spawnHashFunctionName = hashFunctionName
			}
			digestFunction, err := executionLogDigestFunction(instanceName, spawnHashFunctionName, len(spawn.GetDigest().GetHash()))
			if err != nil {
				return nil, util.StatusWrap(err, "Failed to determine execution log digest function")
			}
			actionDigest, err := digestFunction.NewDigest(spawn.GetDigest().GetHash(), spawn.GetDigest().GetSizeBytes())
			if err != nil {
				return nil, util.StatusWrap(err, "Invalid Action digest in execution log")
			}
			executionLogAction.actionDigest = actionDigest
		}
		actions = append(actions, executionLogAction)
	}

	return actions, nil
}

func clearMissingActionDigests(actions []executionLogAction, missingDigests storagedigest.Set) {
	missing := make(map[storagedigest.Digest]struct{}, missingDigests.Length())
	for _, digest := range missingDigests.Items() {
		missing[digest] = struct{}{}
	}
	for i := range actions {
		if _, ok := missing[actions[i].actionDigest]; ok {
			actions[i].actionDigest = storagedigest.BadDigest
		}
	}
}

func clearAllActionDigests(actions []executionLogAction) {
	for i := range actions {
		actions[i].actionDigest = storagedigest.BadDigest
	}
}

func (r *buildEventRecorder) readExecutionLogActions(ctx context.Context, buildToolLogs *bes.BuildToolLogs) ([]executionLogAction, error) {
	if r.contentAddressableStorage == nil {
		return nil, nil
	}
	executionLog := getCompactExecutionLog(buildToolLogs)
	if executionLog == nil || executionLog.GetUri() == "" {
		return nil, nil
	}
	executionLogDigest := files.GetDigestFromURI(executionLog.GetUri())
	if executionLogDigest == storagedigest.BadDigest {
		return nil, fmt.Errorf("execution log has an invalid ByteStream URI")
	}

	reader := r.contentAddressableStorage.Get(ctx, executionLogDigest).ToReader()
	defer reader.Close()
	actions, err := parseCompactExecutionLog(reader, executionLogDigest.GetInstanceName())
	if err != nil {
		return nil, err
	}

	digestSetBuilder := storagedigest.NewSetBuilder(len(actions))
	for _, action := range actions {
		if action.actionDigest != storagedigest.BadDigest {
			digestSetBuilder = digestSetBuilder.Add(action.actionDigest)
		}
	}
	actionDigests := digestSetBuilder.Build()
	if actionDigests.Empty() {
		return actions, nil
	}
	missingDigests, err := r.contentAddressableStorage.FindMissing(ctx, actionDigests)
	if err != nil {
		clearAllActionDigests(actions)
		slog.WarnContext(ctx, "Could not verify Action blobs in CAS; bb-browser links will be unavailable", "invocation_id", r.InvocationID, "err", err)
		return actions, nil
	}
	clearMissingActionDigests(actions, missingDigests)
	return actions, nil
}

func (r *buildEventRecorder) readExecutionLogActionsFromBatch(ctx context.Context, batch []BuildEventWithInfo) []executionLogAction {
	var actions []executionLogAction
	for _, info := range batch {
		if _, ok := info.Event.GetId().GetId().(*bes.BuildEventId_BuildToolLogs); !ok {
			continue
		}
		parsedActions, err := r.readExecutionLogActions(ctx, info.Event.GetBuildToolLogs())
		if err != nil {
			slog.WarnContext(ctx, "Could not read Bazel execution log; Action Cache links will be unavailable", "invocation_id", r.InvocationID, "err", err)
			continue
		}
		actions = append(actions, parsedActions...)
	}
	return actions
}

type actionMatchKey struct {
	targetLabel string
	mnemonic    string
	outputPath  string
}

type actionOutputKey struct {
	mnemonic   string
	outputPath string
}

func addUniqueActionMatch[K comparable](matches map[K]int64, ambiguous map[K]struct{}, key K, actionID int64) {
	if _, ok := ambiguous[key]; ok {
		return
	}
	if existingID, ok := matches[key]; ok && existingID != actionID {
		delete(matches, key)
		ambiguous[key] = struct{}{}
		return
	}
	matches[key] = actionID
}

type executionLogActionMetadata struct {
	runner       string
	actionDigest storagedigest.Digest
}

func executionLogActionMetadataEqual(left, right executionLogActionMetadata) bool {
	if left.runner != right.runner {
		return false
	}
	if left.actionDigest == storagedigest.BadDigest || right.actionDigest == storagedigest.BadDigest {
		return left.actionDigest == right.actionDigest
	}
	return left.actionDigest.String() == right.actionDigest.String()
}

func matchExecutionLogActions(databaseActions []*ent.Action, executionLogActions []executionLogAction) map[int64]executionLogActionMetadata {
	exactMatches := map[actionMatchKey]int64{}
	ambiguousExactMatches := map[actionMatchKey]struct{}{}
	outputMatches := map[actionOutputKey]int64{}
	ambiguousOutputMatches := map[actionOutputKey]struct{}{}
	for _, databaseAction := range databaseActions {
		exactKey := actionMatchKey{
			targetLabel: databaseAction.Label,
			mnemonic:    databaseAction.Type,
			outputPath:  databaseAction.PrimaryOutput,
		}
		addUniqueActionMatch(exactMatches, ambiguousExactMatches, exactKey, databaseAction.ID)
		outputKey := actionOutputKey{mnemonic: databaseAction.Type, outputPath: databaseAction.PrimaryOutput}
		addUniqueActionMatch(outputMatches, ambiguousOutputMatches, outputKey, databaseAction.ID)
	}

	matchedActions := map[int64]executionLogActionMetadata{}
	ambiguousActionIDs := map[int64]struct{}{}
	for _, executionLogAction := range executionLogActions {
		for _, outputPath := range executionLogAction.outputPaths {
			exactKey := actionMatchKey{
				targetLabel: executionLogAction.targetLabel,
				mnemonic:    executionLogAction.mnemonic,
				outputPath:  outputPath,
			}
			actionID, ok := exactMatches[exactKey]
			if !ok {
				actionID, ok = outputMatches[actionOutputKey{mnemonic: executionLogAction.mnemonic, outputPath: outputPath}]
			}
			if !ok {
				continue
			}
			if _, ambiguous := ambiguousActionIDs[actionID]; ambiguous {
				continue
			}
			metadata := executionLogActionMetadata{
				runner:       executionLogAction.runner,
				actionDigest: executionLogAction.actionDigest,
			}
			if existingMetadata, exists := matchedActions[actionID]; exists && !executionLogActionMetadataEqual(existingMetadata, metadata) {
				delete(matchedActions, actionID)
				ambiguousActionIDs[actionID] = struct{}{}
				continue
			}
			matchedActions[actionID] = metadata
		}
	}
	return matchedActions
}

func (r *buildEventRecorder) saveExecutionLogActionMetadata(ctx context.Context, tx database.Handle, executionLogActions []executionLogAction) error {
	if len(executionLogActions) == 0 {
		return nil
	}

	databaseActions, err := tx.Ent().Action.Query().
		Where(
			action.BazelInvocationID(r.InvocationDbID),
			action.PrimaryOutputNotNil(),
		).
		All(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to query actions for execution log correlation")
	}

	for actionID, metadata := range matchExecutionLogActions(databaseActions, executionLogActions) {
		update := tx.Ent().Action.UpdateOneID(actionID)
		if metadata.runner != "" {
			update.SetRunner(metadata.runner)
		}
		if metadata.actionDigest != storagedigest.BadDigest {
			actionDigest := metadata.actionDigest
			digestID, err := tx.Ent().Digest.Create().
				SetRev2InstanceName(actionDigest.GetInstanceName().String()).
				SetDigestFunction(int16(actionDigest.GetDigestFunction().GetEnumValue())).
				SetHash(actionDigest.GetHashBytes()).
				SetSizeBytes(actionDigest.GetSizeBytes()).
				OnConflictColumns(entdigest.FieldRev2InstanceName, entdigest.FieldDigestFunction, entdigest.FieldHash, entdigest.FieldSizeBytes).
				Ignore().
				ID(ctx)
			if err != nil {
				return util.StatusWrap(err, "Failed to save Action digest")
			}
			update.SetActionDigestID(digestID)
		}
		if err := update.Exec(ctx); err != nil {
			return util.StatusWrap(err, "Failed to associate Action execution metadata")
		}
	}
	return nil
}
