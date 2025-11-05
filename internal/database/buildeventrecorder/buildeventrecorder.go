package buildeventrecorder

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/connectionmetadata"
	apicommon "github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BuildEventRecorder contains information about a Bazel invocation that is
// required when processing individual build events.
type BuildEventRecorder struct {
	db                  *ent.Client
	problemDetector     detectors.ProblemDetector
	blobArchiver        processing.BlobMultiArchiver
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel
	tracer              trace.Tracer

	InstanceName           string
	InvocationID           string
	InvocationDbID         int
	CorrelatedInvocationID string
	IsRealTime             bool
}

// NewBuildEventRecorder creates a new BuildEventRecorder
func NewBuildEventRecorder(
	ctx context.Context,
	db *ent.Client,
	instanceNameAuthorizer auth.Authorizer,
	blobArchiver processing.BlobMultiArchiver,
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel,
	tracerProvider trace.TracerProvider,
	instanceName string,
	invocationID string,
	correlatedInvocationID string,
	isRealTime bool,
) (*BuildEventRecorder, error) {
	if invocationID == "" {
		return nil, status.Error(codes.InvalidArgument, "Invocation ID is required")
	}
	if !apicommon.IsInstanceNameAllowed(ctx, instanceNameAuthorizer, instanceName) {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Instance name %q is not allowed", instanceName))
	}

	invocationDb, err := findOrCreateInvocation(ctx, db, invocationID, instanceName)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find or create bazel invocation")
	}

	return &BuildEventRecorder{
		db:                  db,
		problemDetector:     detectors.NewProblemDetector(),
		blobArchiver:        blobArchiver,
		saveTargetDataLevel: saveTargetDataLevel,
		tracer:              tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"),

		InstanceName:           instanceName,
		InvocationID:           invocationID,
		InvocationDbID:         invocationDb.ID,
		CorrelatedInvocationID: correlatedInvocationID,
		IsRealTime:             isRealTime,
	}, nil
}

func findOrCreateInvocation(
	ctx context.Context,
	db *ent.Client,
	invocationID string,
	instanceName string,
) (*ent.BazelInvocation, error) {
	invocationIDUUID, err := uuid.Parse(invocationID)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to parse invocation ID")
	}

	invocationDbID, err := db.BazelInvocation.Create().
		SetInvocationID(invocationIDUUID).
		SetInstanceName(instanceName).
		OnConflictColumns(bazelinvocation.FieldInvocationID).
		Ignore().
		ID(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create or find bazel invocation")
	}

	invocationDb, err := db.BazelInvocation.Query().
		Where(bazelinvocation.ID(invocationDbID)).
		Select(
			bazelinvocation.FieldID,
			bazelinvocation.FieldBepCompleted,
			bazelinvocation.FieldInstanceName,
		).
		Only(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to query bazel invocation by ID")
	}

	if invocationDb.InstanceName != instanceName {
		return nil, status.Errorf(codes.FailedPrecondition, "Invocation with ID %q already exists with different instance name. Possible UUID collision.", invocationID)
	}
	if invocationDb.BepCompleted {
		return nil, status.Errorf(codes.FailedPrecondition, "Invocation with ID %q already exists and is locked for writing", invocationID)
	}
	return invocationDb, nil
}

// StartLoggingConnectionMetadata starts a goroutine that periodically
// logs connection metadata for the invocation associated with this
// BuildEventRecorder. It continues to do so until the provided context
// is canceled.
func (r *BuildEventRecorder) StartLoggingConnectionMetadata(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := r.db.ConnectionMetadata.Create().
					SetBazelInvocationID(r.InvocationDbID).
					SetConnectionLastOpenAt(time.Now()).
					OnConflictColumns(connectionmetadata.BazelInvocationColumn).
					UpdateNewValues().
					Exec(ctx)
				if err != nil {
					slog.Warn("Failed to log connection metadata", "err", err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
