package buildeventrecorder

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func eventTypeName(buildEvent *events.BuildEvent) string {
	if eventPayload := buildEvent.GetId().GetId(); eventPayload != nil {
		return reflect.TypeOf(eventPayload).Elem().Name()
	}
	return "<nil>"
}

func (r *BuildEventRecorder) saveEvent(
	ctx context.Context,
	info BuildEventWithInfo,
) (err error) {
	buildEvent := info.Event
	sequenceNumber := info.SequenceNumber
	ctx, span := r.tracer.Start(ctx,
		fmt.Sprintf("BuildEventRecorder.saveEvent_%s", eventTypeName(buildEvent)),
		trace.WithAttributes(
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
			attribute.Int("build_event.sequence_number", int(sequenceNumber)),
			attribute.String("build_event.type", eventTypeName(buildEvent)),
		),
	)
	defer func() {
		if err != nil {
			span.SetStatus(otelcodes.Error, err.Error())
			span.RecordError(err)
		}
		span.End()
	}()

	// ReadCommitted is most often going to be sufficient, but there may
	// be edge cases in some saveBuildEvent implementations.
	//
	// The outer logic revolving around bep_completed and event metadata
	// is guaranteed to be consistent during the transaction due to the
	// use of a shared read lock for the invocation and an optimistic
	// lock for the event metadata.
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
		return status.Error(codes.FailedPrecondition, "Attempted to insert a build event but  but the invocation was already completed.")
	}

	err = r.saveBuildEvent(ctx, tx, buildEvent)
	if err != nil {
		return util.StatusWrapf(err, "Failed to save build event of type %T", buildEvent.GetId().GetId())
	}

	err = r.saveBazelInvocationProblems(ctx, tx.Ent(), buildEvent)
	if err != nil {
		return util.StatusWrap(err, "Failed to save bazel invocation problems")
	}

	r.handledEvents.bitmap.Add(sequenceNumber)
	r.saveHandledEvents(ctx, tx, time.Now())

	err = tx.Commit()
	if err != nil {
		return util.StatusWrap(err, "Failed to commit transaction")
	}
	return nil
}

func (r *BuildEventRecorder) saveBuildEvent(
	ctx context.Context,
	tx database.Tx,
	buildEvent *events.BuildEvent,
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
	case *bes.BuildEventId_Configuration:
		return r.saveBuildConfiguration(ctx, tx.Ent(), buildEvent.GetConfiguration())
	case *bes.BuildEventId_Fetch:
		return r.saveFetch(ctx, tx.Ent(), buildEvent.GetFetch())
	case *bes.BuildEventId_TestResult:
		return r.saveTestResult(ctx, tx.Ent(), buildEvent.GetTestResult(), buildEvent.GetId().GetTestResult().Label)
	case *bes.BuildEventId_TestSummary:
		return r.saveTestSummary(ctx, tx.Ent(), buildEvent.GetTestSummary(), buildEvent.GetId().GetTestSummary().Label)
	case *bes.BuildEventId_BuildToolLogs:
		return r.saveBuildToolLogs(ctx, tx.Ent(), buildEvent.GetBuildToolLogs())
	case *bes.BuildEventId_WorkspaceStatus:
		return r.saveWorkspaceStatus(ctx, tx.Ent(), buildEvent.GetWorkspaceStatus())
	default:
		return nil
	}
}
