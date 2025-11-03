package buildeventrecorder

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventmetadata"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func eventTypeName(buildEvent *events.BuildEvent) string {
	if eventPayload := buildEvent.GetId().GetId(); eventPayload != nil {
		return reflect.TypeOf(eventPayload).Elem().Name()
	}
	return "<nil>"
}

// RecordEvent records a build event in the database.
func (r *BuildEventRecorder) RecordEvent(
	ctx context.Context,
	buildEvent *events.BuildEvent,
	sequenceNumber int64,
) error {
	ctx, span := r.tracer.Start(ctx,
		fmt.Sprintf("BuildEventRecorder.recordEvent_%s", eventTypeName(buildEvent)),
		trace.WithAttributes(
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
			attribute.String("build_event.type", eventTypeName(buildEvent)),
		),
	)
	defer span.End()

	// We create the event hash before starting the transaction, as
	// this operation does not need to be part of it.
	eventHash, err := r.getEventHash(buildEvent)
	if err != nil {
		return util.StatusWrap(err, "Failed to get event hash")
	}

	tx, err := r.db.BeginTx(ctx, &entsql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return util.StatusWrap(err, "Failed to create transaction")
	}

	err = r.createEventMetadata(ctx, tx, sequenceNumber, eventHash)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return util.StatusWrap(err, "Failed to create event metadata, and failed to rollback transaction")
		}
		eventMetadata, err := r.db.EventMetadata.Query().
			Where(
				eventmetadata.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
				eventmetadata.SequenceNumberEQ(sequenceNumber),
			).
			Only(ctx)
		if err != nil {
			return util.StatusWrap(err, "Failed to create event metadata, and failed to query existing event metadata")
		}
		if eventMetadata.EventHash == eventHash {
			// This exact event has already been processed. Ignore it and send
			// an ACK back.
			return nil
		}
		return status.Errorf(codes.AlreadyExists, "Event with invocation ID %s and sequence number %d already processed with different content", r.InvocationID, sequenceNumber)
	}

	err = r.saveBuildEvent(ctx, tx, buildEvent)
	if err != nil {
		return common.RollbackAndWrapError(tx, util.StatusWrapf(err, "Failed to save build event of type %T", buildEvent.GetId().GetId()))
	}

	err = r.saveBazelInvocationProblems(ctx, tx, buildEvent)
	if err != nil {
		return common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to save bazel invocation problems"))
	}

	if buildEvent.LastMessage {
		err = tx.BazelInvocation.Update().
			Where(
				bazelinvocation.ID(r.InvocationDbID),
			).
			SetBepCompleted(true).
			Exec(ctx)
		if err != nil {
			return common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to mark BEP as completed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		// If the commit fails, we check if the event has already been handled.
		// This can happen if two identical events are sent concurrently. In
		// this case we should not return an error, and justs send an ACK back.
		exist, qerr := r.db.EventMetadata.Query().
			Where(
				eventmetadata.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
				eventmetadata.SequenceNumber(sequenceNumber),
				eventmetadata.EventHash(eventHash)).
			Exist(ctx)
		if qerr != nil {
			return util.StatusWrap(qerr, "Failed to check if event was already processed")
		}
		if exist {
			// This exact event has already been processed. Ignore it.
			return nil
		}
		return util.StatusWrap(err, "Failed to commit transaction")
	}
	return nil
}

func (r *BuildEventRecorder) getEventHash(buildEvent *events.BuildEvent) (string, error) {
	marshalOptions := proto.MarshalOptions{Deterministic: true}
	data, err := marshalOptions.Marshal(buildEvent)
	if err != nil {
		return "", util.StatusWrap(err, "Failed to marshal build event")
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:]), nil
}

func (r *BuildEventRecorder) createEventMetadata(
	ctx context.Context,
	tx *ent.Tx,
	sequenceNumber int64,
	eventHash string,
) error {
	// TODO: Rewrite error messages

	// No SQL injections are possible here, as all the inputs to the query are
	// table and column names from our own code, or parameters ($1, $2, ...).
	// The parameters are passed separately to the query, so they cannot be
	// interpreted as SQL code.
	query := fmt.Sprintf(`
		INSERT INTO %s (%s, %s, %s, %s)
		SELECT $1, $2, $3, b.id
		FROM %s AS b
		WHERE b.%s = $4 AND b.%s = false`,
		eventmetadata.Table,
		eventmetadata.FieldSequenceNumber,
		eventmetadata.FieldEventReceivedAt,
		eventmetadata.FieldEventHash,
		eventmetadata.BazelInvocationColumn,
		bazelinvocation.Table,
		bazelinvocation.FieldInvocationID,
		bazelinvocation.FieldBepCompleted,
	)
	args := []any{sequenceNumber, time.Now(), eventHash, r.InvocationID}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return util.StatusWrap(err, "failed executing conditional insert")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return util.StatusWrap(err, "failed to get rows affected by insert")
	}
	if rowsAffected != 1 {
		return status.Errorf(codes.AlreadyExists, "Event with invocation ID %s and sequence number %d already processed.", r.InvocationID, sequenceNumber)
	}
	return nil
}

func (r *BuildEventRecorder) saveBuildEvent(
	ctx context.Context,
	tx *ent.Tx,
	buildEvent *events.BuildEvent,
) error {
	switch buildEvent.GetId().GetId().(type) {
	case *bes.BuildEventId_Started:
		return r.saveStarted(ctx, tx, buildEvent.GetStarted())
	case *bes.BuildEventId_BuildMetadata:
		return r.saveBuildMetadata(ctx, tx, buildEvent.GetBuildMetadata())
	case *bes.BuildEventId_OptionsParsed:
		return r.saveOptionsParsed(ctx, tx, buildEvent.GetOptionsParsed())
	case *bes.BuildEventId_BuildFinished:
		return r.saveBuildFinished(ctx, tx, buildEvent.GetFinished())
	case *bes.BuildEventId_BuildMetrics:
		return r.saveBuildMetrics(ctx, tx, buildEvent.GetBuildMetrics())
	case *bes.BuildEventId_StructuredCommandLine:
		return r.saveStructuredCommandLine(ctx, tx, buildEvent.GetStructuredCommandLine())
	case *bes.BuildEventId_Configuration:
		return r.saveBuildConfiguration(ctx, tx, buildEvent.GetConfiguration())
	case *bes.BuildEventId_TargetConfigured:
		return r.saveTargetConfigured(ctx, tx, buildEvent.GetConfigured(), buildEvent.GetTargetConfiguredLabel())
	case *bes.BuildEventId_TargetCompleted:
		return r.saveTargetCompleted(ctx, tx, buildEvent.GetCompleted(), buildEvent.GetTargetCompletedLabel(), buildEvent.GetAborted())
	case *bes.BuildEventId_Fetch:
		return r.saveFetch(ctx, tx, buildEvent.GetFetch())
	case *bes.BuildEventId_TestResult:
		return r.saveTestResult(ctx, tx, buildEvent.GetTestResult(), buildEvent.GetId().GetTestResult().Label)
	case *bes.BuildEventId_TestSummary:
		return r.saveTestSummary(ctx, tx, buildEvent.GetTestSummary(), buildEvent.GetId().GetTestSummary().Label)
	case *bes.BuildEventId_BuildToolLogs:
		return r.saveBuildToolLogs(ctx, tx, buildEvent.GetBuildToolLogs())
	case *bes.BuildEventId_Progress:
		return r.saveBuildProgress(ctx, tx, buildEvent, buildEvent.GetProgress())
	case *bes.BuildEventId_WorkspaceStatus:
		return r.saveWorkspaceStatus(ctx, tx, buildEvent.GetWorkspaceStatus())
	default:
		return nil
	}
}
