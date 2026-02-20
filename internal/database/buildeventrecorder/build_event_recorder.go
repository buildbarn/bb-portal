package buildeventrecorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RoaringBitmap/roaring"
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent/connectionmetadata"
	apicommon "github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sqlc-dev/pqtype"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BuildEventWithInfo couples a build event with additional metadata
// required for processing.
type BuildEventWithInfo struct {
	Event *bes.BuildEvent
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

// BuildEventRecorder performs the actual recording of build events for
// an invocation to the storage layer.
type BuildEventRecorder interface {
	// StartLoggingConnectionMetadata starts a logging routine that will
	// log that there is a connection alive and well for the specific
	// invocation.
	StartLoggingConnectionMetadata(ctx context.Context)
	// SaveBatch saves a batch of build events to the storage layer.
	SaveBatch(ctx context.Context, batch []BuildEventWithInfo) error
}

type buildEventRecorder struct {
	db            database.Client
	handledEvents handledEvents
	saveDataLevel *bb_portal.BuildEventStreamService_SaveDataLevel
	tracer        trace.Tracer

	InstanceName     string
	InstanceNameDbID int64
	InvocationID     string
	InvocationDbID   int64
	IsRealTime       bool
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
	saveDataLevel *bb_portal.BuildEventStreamService_SaveDataLevel,
	tracerProvider trace.TracerProvider,
	instanceName string,
	invocationID string,
	isRealTime bool,
	extractors *authmetadataextraction.AuthMetadataExtractors,
) (BuildEventRecorder, error) {
	if invocationID == "" {
		return nil, status.Error(codes.InvalidArgument, "Invocation ID is required")
	}
	if !apicommon.IsInstanceNameAllowed(ctx, instanceNameAuthorizer, instanceName) {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Instance name %q is not allowed", instanceName))
	}

	tracer := tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/database/buildeventrecorder")

	userDb, err := FindOrCreateAuthenticatedUser(ctx, db, extractors, prometheusmetrics.AuthenticatedUsersCount)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find or create authenticated user")
	}

	instanceNameDbID, invocationDbID, err := FindOrCreateInvocation(ctx, db, invocationID, instanceName, tracer, userDb, prometheusmetrics.Invocations)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find or create bazel invocation")
	}

	return &buildEventRecorder{
		db:            db,
		saveDataLevel: saveDataLevel,
		tracer:        tracer,

		InstanceName:     instanceName,
		InstanceNameDbID: instanceNameDbID,
		InvocationID:     invocationID,
		InvocationDbID:   invocationDbID,
		IsRealTime:       isRealTime,
	}, nil
}

// FindOrCreateAuthenticatedUser creates a user or
// returns its ID if it already exists.
func FindOrCreateAuthenticatedUser(
	ctx context.Context,
	db database.Client,
	extractors *authmetadataextraction.AuthMetadataExtractors,
	authenticatedUsersGauge prometheus.Gauge,
) (*int64, error) {
	userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, extractors)
	if userSummary == nil {
		return nil, nil
	}

	var displayName sql.NullString
	if userSummary.DisplayName != nil {
		displayName = sql.NullString{
			String: *userSummary.DisplayName,
			Valid:  true,
		}
	}

	var userInfo pqtype.NullRawMessage
	if userSummary.UserInfo != nil {
		userInfoRaw, err := json.Marshal(userSummary.UserInfo)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to marshal UserInfo")
		}
		userInfo = pqtype.NullRawMessage{
			RawMessage: userInfoRaw,
			Valid:      true,
		}
	}

	authenticatedUserDb, err := db.Sqlc().CreateAuthenticatedUser(ctx,
		sqlc.CreateAuthenticatedUserParams{
			UserUuid:    uuid.NewSHA1(uuid.NameSpaceURL, []byte(userSummary.ExternalID)),
			ExternalID:  userSummary.ExternalID,
			DisplayName: displayName,
			UserInfo:    userInfo,
		},
	)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create authenticated user")
	}
	if authenticatedUserDb.Created {
		authenticatedUsersGauge.Inc()
	}

	return &authenticatedUserDb.ID, nil
}

// FindOrCreateInvocation creates an invocation or
// returns its ID and associated instance name ID
// if it already exists.
func FindOrCreateInvocation(
	ctx context.Context,
	db database.Client,
	invocationID string,
	instanceName string,
	tracer trace.Tracer,
	authenticatedUserDbID *int64,
	invocationsGaugeVec *prometheus.GaugeVec,
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
	invocationDb, err := db.Sqlc().CreateBazelInvocation(ctx, sqlc.CreateBazelInvocationParams{
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
	if invocationDb.Created {
		if authenticatedUserDbID != nil {
			invocationsGaugeVec.WithLabelValues(prometheusmetrics.AuthenticatedUsersLabel).Inc()
		} else {
			invocationsGaugeVec.WithLabelValues(prometheusmetrics.UnauthenticatedUsersLabel).Inc()
		}
	}

	return instanceNameDbID, invocationDb.ID, nil
}

// StartLoggingConnectionMetadata starts a goroutine that periodically
// logs connection metadata for the invocation associated with this
// BuildEventRecorder. It continues to do so until the provided context
// is canceled.
func (r *buildEventRecorder) StartLoggingConnectionMetadata(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// It's both fine and expected that we fail to log
				// connection metadata. That simply means we do not have
				// a timestamp for the connection.
				_ = r.db.Ent().ConnectionMetadata.Create().
					SetBazelInvocationID(r.InvocationDbID).
					SetConnectionLastOpenAt(time.Now().UTC()).
					OnConflictColumns(connectionmetadata.BazelInvocationColumn).
					UpdateNewValues().
					Exec(ctx)
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
