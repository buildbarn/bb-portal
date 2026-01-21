package buildeventrecorder

import (
	"context"
	"database/sql"

	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// saveProgressBatch is an efficient implementation of save progress for
// a batch of progress events.
func (r *buildEventRecorder) saveProgressBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}

	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveProgressBatch",
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
		return status.Error(codes.InvalidArgument, "Attempted to insert progress for an invocation but the invocation was already completed.")
	}

	params := sqlc.CreateIncompleteBuildLogsParams{
		BazelInvocationID: int64(r.InvocationDbID),
		SnippetIds:        make([]int32, 0, len(batch)),
		LogSnippets:       make([][]byte, 0, len(batch)),
	}
	for _, x := range batch {
		be := x.Event
		progress := be.GetProgress()
		if progress == nil {
			return status.Error(codes.InvalidArgument, "Received non progress event to batch progress method")
		}
		opaqueCount := be.GetId().GetProgress().GetOpaqueCount()
		logText := progress.GetStderr()
		if logText != progress.GetStdout() {
			logText += progress.GetStdout()
		}
		if logText != "" {
			params.SnippetIds = append(params.SnippetIds, opaqueCount)
			params.LogSnippets = append(params.LogSnippets, []byte(logText))
		}
	}

	if err = tx.Sqlc().CreateIncompleteBuildLogs(ctx, params); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert incomplete build logs")
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err = tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit batch of progress events")
	}

	return nil
}
