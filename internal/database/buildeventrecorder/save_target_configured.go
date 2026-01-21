package buildeventrecorder

import (
	"context"
	"database/sql"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type targetKey struct {
	label, aspect, kind string
}

// saveTargetConfiguredBatch efficiently saves a batch of target
// configured events.
func (r *buildEventRecorder) saveTargetConfiguredBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}

	switch r.saveDataLevel.GetLevel().(type) {
	case *bb_portal.BuildEventStreamService_SaveDataLevel_Basic:
		return nil
	case *bb_portal.BuildEventStreamService_SaveDataLevel_BasicAndTarget:
		// Continue processing.
	default:
		return status.Error(codes.Internal, "Attempted to save target completed events when `saveDataLevel` is not recognized. This is probably a bug.")
	}

	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveTargetConfiguredBatch",
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

	row, err := tx.Sqlc().LockBazelInvocationCompletion(ctx, int64(r.InvocationDbID))
	if err != nil {
		return util.StatusWrap(err, "Failed to lock bep completed for invocation")
	}
	if row.BepCompleted {
		return status.Error(codes.InvalidArgument, "Attempted to configure targets for an invocation but the invocation was already completed.")
	}

	// Don't handle aborted events.
	filteredBatch := make([]BuildEventWithInfo, 0, len(batch))
	for _, x := range batch {
		if x.Event.GetConfigured() != nil {
			filteredBatch = append(filteredBatch, x)
		}
	}

	targetIds, err := getOrCreateTargets(ctx, r.InstanceNameDbID, tx, filteredBatch)
	if err != nil {
		return util.StatusWrap(err, "Failed to get or create targets")
	}

	if err := createTargetKindMappingsBulk(ctx, r.IsRealTime, r.InvocationDbID, tx, filteredBatch, targetIds); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err = tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit batch of target configured events")
	}

	return nil
}

func createTargetKindMappingsBulk(ctx context.Context, isRealTime bool, invocationDbID int64, tx database.Handle, batch []BuildEventWithInfo, targetIds map[targetKey]int) error {
	params := sqlc.CreateTargetKindMappingsBulkParams{
		BazelInvocationID: int64(invocationDbID),
		TargetIds:         make([]int64, len(batch)),
		StartTimes:        make([]int64, len(batch)),
	}
	for i, info := range batch {
		be := info.Event
		configured := be.GetConfigured()
		targetConfigured := be.GetId().GetTargetConfigured()
		params.TargetIds[i] = int64(targetIds[targetKey{targetConfigured.Label, targetConfigured.Aspect, configured.TargetKind}])
		if isRealTime {
			params.StartTimes[i] = info.AddedAt.UnixMilli()
		}
	}
	return tx.Sqlc().CreateTargetKindMappingsBulk(ctx, params)
}

// getOrCreateTargets performs an select or insert in an efficient
// manner. Most of the time the targets will already exist in which case
// this will only be a single select. Otherwise it inserts the missing
// targets.
func getOrCreateTargets(ctx context.Context, instanceNameID int64, tx database.Handle, batch []BuildEventWithInfo) (map[targetKey]int, error) {
	var labels, aspects, kinds []string
	uniqueKeys := make(map[targetKey]struct{}, len(batch))

	for _, info := range batch {
		be := info.Event
		configured := be.GetConfigured()
		targetConfigured := be.GetId().GetTargetConfigured()
		if configured == nil || targetConfigured == nil {
			return nil, status.Error(codes.InvalidArgument, "Received non target configured event")
		}
		key := targetKey{label: targetConfigured.Label, aspect: targetConfigured.Aspect, kind: configured.TargetKind}
		if _, exists := uniqueKeys[key]; !exists {
			uniqueKeys[key] = struct{}{}
			labels = append(labels, key.label)
			aspects = append(aspects, key.aspect)
			kinds = append(kinds, key.kind)
		}
	}

	// First and typically only select.
	foundRows, err := tx.Sqlc().FindTargets(ctx, sqlc.FindTargetsParams{
		InstanceNameID: int64(instanceNameID),
		Labels:         labels,
		Aspects:        aspects,
		TargetKinds:    kinds,
	})
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find existing targets")
	}
	result := make(map[targetKey]int, len(uniqueKeys))
	for _, row := range foundRows {
		result[targetKey{label: row.Label, aspect: row.Aspect, kind: row.TargetKind}] = int(row.ID)
	}
	if len(result) == len(uniqueKeys) {
		return result, nil
	}

	// Some targets were missing, insert them.
	missLabels, missAspects, missKinds := findMissingRows(labels, aspects, kinds, result)
	newRows, err := tx.Sqlc().CreateTargets(ctx, sqlc.CreateTargetsParams{
		InstanceNameID: int64(instanceNameID),
		Labels:         missLabels,
		Aspects:        missAspects,
		TargetKinds:    missKinds,
	})
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to insert missing targets")
	}
	for _, row := range newRows {
		result[targetKey{label: row.Label, aspect: row.Aspect, kind: row.TargetKind}] = int(row.ID)
	}
	if len(result) != len(uniqueKeys) {
		return nil, status.Error(codes.Unavailable, "Not all missing targets were created.")
	}
	return result, nil
}

func findMissingRows(labels, aspects, kinds []string, result map[targetKey]int) (missLabels, missAspects, missKinds []string) {
	for i, label := range labels {
		key := targetKey{label: label, aspect: aspects[i], kind: kinds[i]}
		if _, found := result[key]; !found {
			missLabels = append(missLabels, label)
			missAspects = append(missAspects, aspects[i])
			missKinds = append(missKinds, kinds[i])
		}
	}
	return missLabels, missAspects, missKinds
}
