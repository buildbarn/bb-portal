package buildeventrecorder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/buildbarn/bb-portal/ent/gen/ent/authenticateduser"
	"github.com/buildbarn/bb-portal/ent/gen/ent/connectionmetadata"
	apicommon "github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
	"github.com/buildbarn/bb-portal/pkg/events"
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

// BuildEventWithInfo couples a build event with additional metadata
// required for processing.
type BuildEventWithInfo struct {
	Event *events.BuildEvent
	// Sequence number is int64 by spec, however it is also starting at
	// 1 and incrementing by 1 for each event by spec. We therefore
	// assume we can fit the SequenceNumber in an uint32 and in case we
	// can not we will refuse to process the event.
	SequenceNumber uint32
	// TODO: This field is only used for duration calculation of target
	// completed events which is dubious to begin with. Might be worth
	// removing.
	AddedAt time.Time
}

// BuildEventRecorder contains information about a Bazel invocation that is
// required when processing individual build events.
type BuildEventRecorder struct {
	db                  database.Client
	problemDetector     detectors.ProblemDetector
	handledEvents       handledEvents
	blobArchiver        processing.BlobMultiArchiver
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel
	saveTestDataLevel   *bb_portal.BuildEventStreamService_SaveTestDataLevel
	tracer              trace.Tracer

	InstanceName           string
	InstanceNameDbID       int64
	InvocationID           string
	InvocationDbID         int64
	CorrelatedInvocationID string
	IsRealTime             bool
}

type handledEvents struct {
	bitmap  *roaring.Bitmap
	version int64
	id      int64
}

// NewBuildEventRecorder creates a new BuildEventRecorder
func NewBuildEventRecorder(
	ctx context.Context,
	db database.Client,
	instanceNameAuthorizer auth.Authorizer,
	blobArchiver processing.BlobMultiArchiver,
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel,
	saveTestDataLevel *bb_portal.BuildEventStreamService_SaveTestDataLevel,
	tracerProvider trace.TracerProvider,
	instanceName string,
	invocationID string,
	correlatedInvocationID string,
	isRealTime bool,
	extractors *authmetadataextraction.AuthMetadataExtractors,
	uuidGenerator util.UUIDGenerator,
) (*BuildEventRecorder, error) {
	if invocationID == "" {
		return nil, status.Error(codes.InvalidArgument, "Invocation ID is required")
	}
	if !apicommon.IsInstanceNameAllowed(ctx, instanceNameAuthorizer, instanceName) {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Instance name %q is not allowed", instanceName))
	}

	tracer := tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/buildeventrecorder")

	userDb, err := findOrCreateAuthenticatedUser(ctx, db, extractors, uuidGenerator)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find or create authenticated user")
	}

	instanceNameDbID, invocationDbID, err := findOrCreateInvocation(ctx, db, invocationID, instanceName, tracer, userDb)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find or create bazel invocation")
	}

	return &BuildEventRecorder{
		db:                  db,
		problemDetector:     detectors.NewProblemDetector(),
		blobArchiver:        blobArchiver,
		saveTargetDataLevel: saveTargetDataLevel,
		saveTestDataLevel:   saveTestDataLevel,
		tracer:              tracer,

		InstanceName:           instanceName,
		InstanceNameDbID:       instanceNameDbID,
		InvocationID:           invocationID,
		InvocationDbID:         invocationDbID,
		CorrelatedInvocationID: correlatedInvocationID,
		IsRealTime:             isRealTime,
	}, nil
}

func findOrCreateAuthenticatedUser(
	ctx context.Context,
	db database.Client,
	extractors *authmetadataextraction.AuthMetadataExtractors,
	uuidGenerator util.UUIDGenerator,
) (*int64, error) {
	userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, extractors)
	if userSummary == nil {
		return nil, nil
	}

	// UserUUID is immutable and will not be regenerated
	// if the object already exists.
	userRecord := db.Ent().AuthenticatedUser.Create().
		SetUserUUID(util.Must(uuidGenerator())).
		SetExternalID(userSummary.ExternalID).
		SetNillableDisplayName(userSummary.DisplayName)
	if userSummary.UserInfo != nil {
		userRecord.SetUserInfo(userSummary.UserInfo)
	}
	id, err := userRecord.OnConflictColumns(authenticateduser.FieldExternalID).UpdateNewValues().ID(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to upsert authenticated user")
	}
	return &id, nil
}

func findOrCreateInvocation(
	ctx context.Context,
	db database.Client,
	invocationID string,
	instanceName string,
	tracer trace.Tracer,
	authenticatedUserDbID *int64,
) (int64, int64, error) {
	invocationUUID, err := uuid.Parse(invocationID)
	if err != nil {
		return 0, 0, util.StatusWrap(err, "Failed to parse invocation ID")
	}

	var userID sql.NullInt64
	if authenticatedUserDbID != nil {
		userID = sql.NullInt64{Int64: int64(*authenticatedUserDbID), Valid: true}
	}

	instanceNameDbID, err := db.Sqlc().CreateInstanceName(ctx, instanceName)
	if err != nil {
		return 0, 0, util.StatusWrap(err, "Failed to create instance name")
	}
	// There is a race condition here where the just selected instance
	// name becomes deleted between the select and the insert. This is
	// fine as the client will retry.
	invocationDbID, err := db.Sqlc().CreateBazelInvocation(ctx, sqlc.CreateBazelInvocationParams{
		InvocationID:        invocationUUID,
		InstanceNameID:      instanceNameDbID,
		AuthenticatedUserID: userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, status.Errorf(codes.FailedPrecondition,
				"Failed to create invocation with id %s. The id may refer to a invocation that is locked for writing",
				invocationID,
			)
		}
		return 0, 0, util.StatusWrap(err, "Failed to create bazel invocation")
	}

	return instanceNameDbID, invocationDbID, nil
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
				err := r.db.Ent().ConnectionMetadata.Create().
					SetBazelInvocationID(r.InvocationDbID).
					SetConnectionLastOpenAt(time.Now().UTC()).
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
