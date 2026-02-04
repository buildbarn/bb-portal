package buildeventrecorder

import (
	"context"
	"database/sql"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// saveTestSummaryBatch efficiently saves a batch of target
// completed for a set of events where the corresponding target
// configured event has already been handled.
func (r *BuildEventRecorder) saveTestSummaryBatch(ctx context.Context, batch []BuildEventWithInfo) error {
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
		"BuildEventRecorder.saveTestSummaryBatch",
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
		if x.Event.GetTestSummary() != nil {
			filteredBatch = append(filteredBatch, x)
		}
	}

	params := sqlc.UpdateTestSummariesBulkParams{
		BazelInvocationID: int64(r.InvocationDbID),
		InstanceNameID:    int64(r.InstanceNameDbID),
		Labels:            make([]string, len(filteredBatch)),
		ConfigIds:         make([]string, len(filteredBatch)),
		OverallStatuses:   make([]string, len(filteredBatch)),
		TotalRunCounts:    make([]int32, len(filteredBatch)),
		RunCounts:         make([]int32, len(filteredBatch)),
		AttemptCounts:     make([]int32, len(filteredBatch)),
		ShardCounts:       make([]int32, len(filteredBatch)),
		TotalNumCacheds:   make([]int32, len(filteredBatch)),
		StartTimes:        make([]time.Time, len(filteredBatch)),
		StopTimes:         make([]time.Time, len(filteredBatch)),
		Durations:         make([]int64, len(filteredBatch)),
	}

	for i, x := range filteredBatch {
		be := x.Event
		testSummaryID := be.GetId().GetTestSummary()
		testSummary := be.GetTestSummary()

		if testSummaryID == nil {
			return status.Error(codes.InvalidArgument, "Test summary event does not have test summary identifier set.")
		}
		if testSummary == nil {
			return status.Error(codes.InvalidArgument, "Test summary event does not have test summary set.")
		}

		params.Labels[i] = testSummaryID.Label
		params.ConfigIds[i] = testSummaryID.Configuration.GetId()
		params.TotalRunCounts[i] = testSummary.TotalRunCount
		params.RunCounts[i] = testSummary.RunCount
		params.AttemptCounts[i] = testSummary.AttemptCount
		params.ShardCounts[i] = testSummary.ShardCount
		params.TotalNumCacheds[i] = testSummary.TotalNumCached

		if status, ok := bes.TestStatus_name[int32(testSummary.OverallStatus)]; ok {
			params.OverallStatuses[i] = status
		} else {
			params.OverallStatuses[i] = "NO_STATUS"
		}

		if testSummary.FirstStartTime != nil && testSummary.FirstStartTime.AsTime().IsZero() == false {
			params.StartTimes[i] = testSummary.FirstStartTime.AsTime()
		}
		if testSummary.LastStopTime != nil && testSummary.LastStopTime.AsTime().IsZero() == false {
			params.StopTimes[i] = testSummary.LastStopTime.AsTime()
		}
		if testSummary.TotalRunDuration != nil && testSummary.TotalRunDuration.AsDuration() != 0 {
			params.Durations[i] = int64(testSummary.TotalRunDuration.AsDuration().Milliseconds())
		}
	}

	affectedRows, err := tx.Sqlc().UpdateTestSummariesBulk(ctx, params)
	if err != nil {
		return util.StatusWrap(err, "Failed to bulk update test summaries")
	}
	if affectedRows != int64(len(filteredBatch)) {
		return status.Errorf(codes.Internal, "Expected to update %d test summaries, but only updated %d", len(batch), affectedRows)
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err = tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit batch of target configured events")
	}

	return nil
}
