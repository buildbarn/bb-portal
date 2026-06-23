package buildeventrecorder

import (
	"context"
	"database/sql"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// insertIncompleteArtifactGraphEvents serializes the artifact-graph-relevant
// BuildEvents in the batch (NamedSetOfFiles and TargetCompleted variants)
// and bulk-inserts them into the incomplete_artifact_graphs staging table
// using the supplied transaction. No-op unless the save level requests
// artifact data. Persisting incrementally keeps the recorder stateless so
// it survives failover, mirroring how progress events accumulate in
// incomplete_build_logs.
func (r *buildEventRecorder) insertIncompleteArtifactGraphEvents(ctx context.Context, tx database.Tx, batch []BuildEventWithInfo) error {
	if !r.artifactsEnabled() || len(batch) == 0 {
		return nil
	}
	params := sqlc.CreateIncompleteArtifactGraphsParams{
		BazelInvocationID: r.InvocationDbID,
		SeqIds:            make([]int32, 0, len(batch)),
		Events:            make([][]byte, 0, len(batch)),
	}
	for _, x := range batch {
		if x.Event == nil {
			continue
		}
		// Only the NamedSetOfFiles and TargetCompleted variants describe
		// the artifact graph; skip anything else routed here.
		if x.Event.GetNamedSetOfFiles() == nil && x.Event.GetCompleted() == nil {
			continue
		}
		payload, err := proto.Marshal(x.Event)
		if err != nil {
			return util.StatusWrap(err, "Failed to marshal artifact graph BuildEvent")
		}
		params.SeqIds = append(params.SeqIds, int32(x.SequenceNumber))
		params.Events = append(params.Events, payload)
	}
	if len(params.Events) == 0 {
		return nil
	}
	if err := tx.Sqlc().CreateIncompleteArtifactGraphs(ctx, params); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert incomplete artifact graph events")
	}
	return nil
}

// saveNamedSetOfFilesBatch persists a batch of NamedSetOfFiles events to
// the artifact-graph staging table. Only invoked when the save level
// requests artifact data; see saveBatch.
func (r *buildEventRecorder) saveNamedSetOfFilesBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}

	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveNamedSetOfFilesBatch",
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
		return status.Error(codes.InvalidArgument, "Attempted to insert named set of files for an invocation but the invocation was already completed.")
	}

	if err := r.insertIncompleteArtifactGraphEvents(ctx, tx, batch); err != nil {
		return err
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err = tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit transaction")
	}
	return nil
}

func filterNamedSetOfFilesBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_NamedSet:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}
