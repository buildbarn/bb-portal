package bepuploader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	mb            = 1024 * 1024
	maxUploadSize = 500 * mb
)

// BepUploader handles upload of Build Event Protocol files via HTTP.
type BepUploader struct {
	db                     *ent.Client
	instanceNameAuthorizer auth.Authorizer
	blobArchiver           processing.BlobMultiArchiver
	saveTargetDataLevel    *bb_portal.BuildEventStreamService_SaveTargetDataLevel
	tracerProvider         trace.TracerProvider
}

// NewBepUploader creates a new BepUploader
func NewBepUploader(db *ent.Client, blobArchiver processing.BlobMultiArchiver, configuration *bb_portal.ApplicationConfiguration, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory, tracerProvider trace.TracerProvider) (*BepUploader, error) {
	if configuration.InstanceNameAuthorizer == nil {
		return nil, status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}
	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	besConfiguration := configuration.BesServiceConfiguration
	if besConfiguration == nil {
		return nil, fmt.Errorf("No BesServiceConfiguration configured")
	}

	saveTargetDataLevel := besConfiguration.SaveTargetDataLevel
	if saveTargetDataLevel == nil || saveTargetDataLevel.Level == nil {
		return nil, fmt.Errorf("No saveTargetDataLevel configured")
	}

	return &BepUploader{
		db:                     db,
		instanceNameAuthorizer: instanceNameAuthorizer,
		blobArchiver:           blobArchiver,
		saveTargetDataLevel:    saveTargetDataLevel,
		tracerProvider:         tracerProvider,
	}, nil
}

// RecordEventNdjsonFile records all build events from an ndjson bep file.
func (b *BepUploader) RecordEventNdjsonFile(ctx context.Context, file io.Reader) (string, int, error) {
	scanner := bufio.NewScanner(file)
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	var buildEventRecorder *buildeventrecorder.BuildEventRecorder = nil

	sequenceNumber := int64(0)
	for scanner.Scan() {
		// When reading from the BES stream, the first event has sequence
		// number 1, so this is fine.
		sequenceNumber++
		lineBytes, err := scanner.Bytes(), scanner.Err()
		if err != nil {
			return "", http.StatusInternalServerError, util.StatusWrap(err, "Failed to read a line from build event file")
		}
		bazelEvent := bes.BuildEvent{}
		err = unmarshaler.Unmarshal(lineBytes, &bazelEvent)
		if err != nil {
			return "", http.StatusBadRequest, util.StatusWrap(err, "Failed to unmarshal JSON BES event")
		}

		if buildEventRecorder == nil {
			buildEventRecorder, err = buildeventrecorder.NewBuildEventRecorder(
				ctx,
				b.db,
				b.instanceNameAuthorizer,
				b.blobArchiver,
				b.saveTargetDataLevel,
				b.tracerProvider,
				"",                                // instanceName
				bazelEvent.GetStarted().GetUuid(), // invocationID
				"",                                // correlatedInvocationID
				false,                             // isRealTime
			)
			if err != nil {
				return "", gprcErrorCodeToHTTPStatus(err), util.StatusWrap(err, "Failed to create BuildEventRecorder")
			}
		}

		// TODO (isakstenstrom): Remove this and send the raw BES event instead. This can only be
		// done when we no longer need JSON serialization of events, like we do for
		// BazelInvocationProblems.
		buildEvent := events.NewBuildEvent(&bazelEvent, json.RawMessage(protojson.Format(&bazelEvent)))
		if err = buildEventRecorder.RecordEvent(ctx, &buildEvent, sequenceNumber); err != nil {
			return "", gprcErrorCodeToHTTPStatus(err), util.StatusWrap(err, "Failed to record build event")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", http.StatusInternalServerError, util.StatusWrap(err, "Failed to read build event file")
	}
	return buildEventRecorder.InvocationID, http.StatusOK, nil
}

// ServeHTTP handles upload of build event files via HTTP.
func (b *BepUploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		msg := fmt.Sprintf("The uploaded file is too big. Please choose an file that's less than %dMB in size", maxUploadSize/mb)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	slog.Info("Received file", "name", fileHeader.Filename, "size", fileHeader.Size)

	invocationID, httpStatus, err := b.RecordEventNdjsonFile(r.Context(), file)
	if err != nil {
		http.Error(w, err.Error(), httpStatus)
		return
	}

	location := fmt.Sprintf("/bazel-invocations/%s", invocationID)
	// NOTE: Want to do http.Redirect(w, r, location, http.StatusSeeOther), but can't get it working with antd Upload widget.
	writeLocationResponse(w, location)
}

// A function to write location responses.
func writeLocationResponse(w http.ResponseWriter, location string) {
	w.WriteHeader(http.StatusOK)
	resp := struct {
		Location string
	}{
		Location: location,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respBody)
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}

func gprcErrorCodeToHTTPStatus(err error) int {
	grpcErr := status.Convert(err)
	switch grpcErr.Code() {
	case 0: // OK
		return http.StatusOK
	case 2: // Unknown
		return http.StatusInternalServerError
	case 3: // InvalidArgument
		return http.StatusBadRequest
	case 5: // NotFound
		return http.StatusNotFound
	case 6: // AlreadyExists
		return http.StatusConflict
	case 7: // PermissionDenied
		return http.StatusForbidden
	case 8: // ResourceExhausted
		return http.StatusTooManyRequests
	case 9: // FailedPrecondition
		return http.StatusPreconditionFailed
	case 12: // Unimplemented
		return http.StatusNotImplemented
	case 13: // Internal
		return http.StatusInternalServerError
	case 14: // Unavailable
		return http.StatusServiceUnavailable
	case 15: // DataLoss
		return http.StatusInternalServerError
	case 16: // Unauthenticated
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
