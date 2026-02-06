package buildeventrecorder

import (
	"context"
	"database/sql"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *buildEventRecorder) saveRemainingBatch(
	ctx context.Context,
	batch []BuildEventWithInfo,
) (err error) {
	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveIndividualEvents",
		trace.WithAttributes(
			attribute.Int("batch_size", len(batch)),
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
		),
	)
	defer span.End()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return util.StatusWrap(err, "Failed to create transaction")
	}
	defer tx.Rollback()
	lock, err := tx.Sqlc().LockBazelInvocationCompletion(ctx, int64(r.InvocationDbID))
	if err != nil {
		return util.StatusWrap(err, "Failed to lock bep completed for invocation")
	}
	if lock.BepCompleted {
		return status.Error(codes.FailedPrecondition, "Attempted to insert build events but the invocation was already completed.")
	}

	for _, info := range batch {
		buildEvent := info.Event
		err = r.saveBuildEvent(ctx, tx, buildEvent)
		if err != nil {
			return util.StatusWrapf(err, "Failed to save build event of type %T", buildEvent.GetId().GetId())
		}
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	err = tx.Commit()
	if err != nil {
		return util.StatusWrap(err, "Failed to commit transaction")
	}
	return nil
}

func (r *buildEventRecorder) saveBuildEvent(
	ctx context.Context,
	tx database.Tx,
	buildEvent *bes.BuildEvent,
) error {
	switch buildEvent.GetId().GetId().(type) {
	case *bes.BuildEventId_Started:
		return r.saveStarted(ctx, tx.Ent(), buildEvent.GetStarted())
	case *bes.BuildEventId_BuildMetadata:
		return r.saveBuildMetadata(ctx, tx.Ent(), buildEvent.GetBuildMetadata())
	case *bes.BuildEventId_OptionsParsed:
		return r.saveOptionsParsed(ctx, tx.Ent(), buildEvent.GetOptionsParsed())
	case *bes.BuildEventId_BuildFinished:
		return r.saveBuildFinished(ctx, tx.Ent(), buildEvent.GetFinished())
	case *bes.BuildEventId_BuildMetrics:
		return r.saveBuildMetrics(ctx, tx.Ent(), buildEvent.GetBuildMetrics())
	case *bes.BuildEventId_StructuredCommandLine:
		return r.saveStructuredCommandLine(ctx, tx.Ent(), buildEvent.GetStructuredCommandLine())
	case *bes.BuildEventId_ActionCompleted:
		return r.saveActionExecuted(ctx, tx.Ent(), buildEvent.GetAction(), buildEvent.GetId().GetActionCompleted())
	case *bes.BuildEventId_Fetch:
		return r.saveFetch(ctx, tx.Ent(), buildEvent.GetFetch())
	case *bes.BuildEventId_BuildToolLogs:
		return r.saveBuildToolLogs(ctx, tx.Ent(), buildEvent.GetBuildToolLogs())
	case *bes.BuildEventId_WorkspaceStatus:
		return r.saveWorkspaceStatus(ctx, tx.Ent(), buildEvent.GetWorkspaceStatus())
	default:
		return nil
	}
}
