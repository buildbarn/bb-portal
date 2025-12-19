package buildeventrecorder

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventmetadata"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
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

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return util.StatusWrap(err, "Failed to create transaction")
	}

	err = r.createEventMetadata(ctx, tx, sequenceNumber, eventHash)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return util.StatusWrap(errors.Join(err, err2), "Failed to create event metadata, and failed to rollback transaction")
		}
		eventMetadata, err2 := r.db.Ent().EventMetadata.Query().
			Where(
				eventmetadata.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
				eventmetadata.SequenceNumberEQ(sequenceNumber),
			).
			Only(ctx)
		if err2 != nil {
			return util.StatusWrap(errors.Join(err, err2), "Failed to create event metadata, and failed to query existing event metadata")
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

	err = r.saveBazelInvocationProblems(ctx, tx.Ent(), buildEvent)
	if err != nil {
		return common.RollbackAndWrapError(tx, util.StatusWrap(err, "Failed to save bazel invocation problems"))
	}

	if buildEvent.LastMessage {
		err = tx.Ent().BazelInvocation.Update().
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
		exist, qerr := r.db.Ent().EventMetadata.Query().
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
	tx database.Tx,
	sequenceNumber int64,
	eventHash string,
) error {
	result, err := tx.Sqlc().RecordEventMetadata(ctx, sqlc.RecordEventMetadataParams{
		SequenceNumber:  sequenceNumber,
		EventHash:       eventHash,
		EventReceivedAt: time.Now(),
		InvocationID:    uuid.MustParse(r.InvocationID),
	})
	if err != nil {
		return util.StatusWrap(err, "Failed recording event metadata")
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
	case *bes.BuildEventId_TargetConfigured:
		return r.saveTargetConfigured(ctx, tx.Ent(), buildEvent.GetConfigured(), buildEvent.GetId().GetTargetConfigured())
	case *bes.BuildEventId_TargetCompleted:
		return r.saveTargetCompleted(ctx, tx.Ent(), buildEvent.GetCompleted(), buildEvent.GetId().GetTargetCompleted(), buildEvent.GetAborted())
	case *bes.BuildEventId_Fetch:
		return r.saveFetch(ctx, tx.Ent(), buildEvent.GetFetch())
	case *bes.BuildEventId_TestResult:
		return r.saveTestResult(ctx, tx.Ent(), buildEvent.GetTestResult(), buildEvent.GetId().GetTestResult().Label)
	case *bes.BuildEventId_TestSummary:
		return r.saveTestSummary(ctx, tx.Ent(), buildEvent.GetTestSummary(), buildEvent.GetId().GetTestSummary().Label)
	case *bes.BuildEventId_BuildToolLogs:
		return r.saveBuildToolLogs(ctx, tx.Ent(), buildEvent.GetBuildToolLogs())
	case *bes.BuildEventId_Progress:
		return r.saveBuildProgress(ctx, tx.Ent(), buildEvent, buildEvent.GetProgress())
	case *bes.BuildEventId_WorkspaceStatus:
		return r.saveWorkspaceStatus(ctx, tx.Ent(), buildEvent.GetWorkspaceStatus())
	default:
		return nil
	}
}
