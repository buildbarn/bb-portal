package buildeventrecorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent/invocationtarget"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"

	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// saveTargetCompletedBatch efficiently saves a batch of target
// completed for a set of events where the corresponding target
// configured event has already been handled.
func (r *BuildEventRecorder) saveTargetCompletedBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if r.saveTargetDataLevel.GetNone() != nil || len(batch) == 0 {
		return nil
	}

	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveTargetCompletedBatch",
		trace.WithAttributes(
			attribute.Int("batch_size", len(batch)),
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
		),
	)
	defer span.End()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return util.StatusWrap(err, "Failed to start transaction")
	}
	defer func() { tx.Rollback() }()

	row, err := tx.Sqlc().LockBazelInvocationCompletion(ctx, r.InvocationDbID)
	if err != nil {
		return util.StatusWrap(err, "Failed to lock bep completed for invocation")
	}
	if row.BepCompleted {
		return status.Error(codes.InvalidArgument, "Attempted to complete targets for an invocation but the invocation was already completed.")
	}

	// Lookup target info.
	targetInfoMap, err := r.resolveTargetInfo(ctx, tx, batch)
	if err != nil {
		return util.StatusWrap(err, "Failed to get target info mapping")
	}

	if err = createInvocationTargetsBulk(ctx, r.IsRealTime, r.InvocationDbID, tx, batch, targetInfoMap); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert invocation targets")
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err = tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit batch of target completed events")
	}

	return nil
}

func createInvocationTargetsBulk(ctx context.Context, isRealTime bool, invocationDbID int64, tx database.Handle, batch []BuildEventWithInfo, targetInfoMap map[invocationTargetKey]completedTargetInfo) error {
	params := sqlc.CreateInvocationTargetsBulkParams{
		BazelInvocationID: int64(invocationDbID),
		TargetIds:         make([]int64, len(batch)),
		Successes:         make([]bool, len(batch)),
		TagsList:          make([]string, len(batch)),
		StartTimes:        make([]int64, len(batch)),
		EndTimes:          make([]int64, len(batch)),
		Durations:         make([]int64, len(batch)),
		FailureMessages:   make([]string, len(batch)),
		AbortReasons:      make([]string, len(batch)),
	}
	for i, x := range batch {
		be := x.Event
		targetCompletedID := be.GetId().GetTargetCompleted()
		targetCompleted := be.GetCompleted()
		aborted := be.GetAborted()
		key := invocationTargetKey{
			label:  targetCompletedID.Label,
			aspect: stripParams(targetCompletedID.Aspect),
		}
		targetInfo := targetInfoMap[key]

		params.TargetIds[i] = int64(targetInfo.targetID)

		// TODO: This logic is really bothersome, maybe we should just
		// remove it completely and/or fetch it from the profile or
		// something.
		if isRealTime && !targetInfo.startTime.IsZero() {
			params.StartTimes[i] = targetInfo.startTime.UnixMilli()
			params.EndTimes[i] = x.AddedAt.UnixMilli()
			params.Durations[i] = x.AddedAt.Sub(targetInfo.startTime).Milliseconds()
		}

		params.Successes[i] = false
		if targetCompleted != nil {
			params.Successes[i] = targetCompleted.Success
			if msg := targetCompleted.FailureDetail.GetMessage(); msg != "" {
				params.FailureMessages[i] = msg
			}
			if len(targetCompleted.Tag) > 0 {
				if b, err := json.Marshal(targetCompleted.Tag); err == nil {
					params.TagsList[i] = string(b)
				}
			}
		}

		if aborted != nil {
			reasonStr, exists := bes.Aborted_AbortReason_name[int32(aborted.Reason)]
			if !exists {
				reasonStr = bes.Aborted_UNKNOWN.String()
			}
			params.AbortReasons[i] = reasonStr
		} else {
			// TODO: shouldn't this really be NULL sql side?
			params.AbortReasons[i] = string(invocationtarget.AbortReasonNONE)
		}
	}
	return tx.Sqlc().CreateInvocationTargetsBulk(ctx, params)
}

type invocationTargetKey struct {
	aspect string
	label  string
}

type completedTargetInfo struct {
	targetID  int
	startTime time.Time
}

// resolveTargetInfo fetches TargetID and StartTimeInMs from the database
// by joining TargetKindMapping with Target.
func (r *BuildEventRecorder) resolveTargetInfo(ctx context.Context, tx database.Handle, batch []BuildEventWithInfo) (map[invocationTargetKey]completedTargetInfo, error) {
	// Deduplicate keys, this logic should not be required as each
	// completed event should refer back to a single configured event.
	keys := make([]invocationTargetKey, 0, len(batch))
	added := make(map[invocationTargetKey]struct{})

	for _, info := range batch {
		id := info.Event.GetId().GetTargetCompleted()
		if id == nil || id.Label == "" {
			return nil, status.Error(codes.InvalidArgument, "Received invalid target completed events to batch target completed method.")
		}
		k := invocationTargetKey{
			aspect: stripParams(id.Aspect),
			label:  id.Label,
		}
		if _, ok := added[k]; !ok {
			keys = append(keys, k)
			added[k] = struct{}{}
		}
	}

	params := sqlc.FindMappedTargetsParams{
		BazelInvocationID: int64(r.InvocationDbID),
		Aspects:           make([]string, len(batch)),
		Labels:            make([]string, len(batch)),
	}
	for i, k := range keys {
		params.Aspects[i] = k.aspect
		params.Labels[i] = k.label
	}
	mappedRows, err := tx.Sqlc().FindMappedTargets(ctx, params)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to query target info")
	}
	result := make(map[invocationTargetKey]completedTargetInfo, len(keys))

	for _, row := range mappedRows {
		var startTime int64
		if row.StartTimeInMs.Valid {
			startTime = row.StartTimeInMs.Int64
		}
		key := invocationTargetKey{label: row.Label, aspect: row.Aspect}
		value := completedTargetInfo{
			targetID:  int(row.TargetID),
			startTime: time.Unix(startTime/1000, (startTime%1000)*1_000_000),
		}
		result[key] = value
	}

	for _, k := range keys {
		if _, found := result[k]; !found {
			return nil, status.Errorf(codes.FailedPrecondition, "Attempted to complete a target for which no configuration was found. %v", k)
		}
	}

	return result, nil
}

func stripParams(aspect string) string {
	if idx := strings.Index(aspect, "["); idx != -1 {
		return aspect[:idx]
	}
	return aspect
}
