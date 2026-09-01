package bep

import (
	"net/http"

	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-portal/internal/api/http/bepuploader"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/trace"
	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewBuildEventProtocolService creates a new service that accepts the Build
// Event Stream and manually updated BEP files, and stores the result in the
// database.
func NewBuildEventProtocolService(
	configuration *bb_portal.BuildEventStreamService,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *http.ServeMux,
	tracerProvider trace.TracerProvider,
	dbClient database.Client,
) error {
	if dbClient == nil {
		return status.Error(codes.NotFound, "No Database configured")
	}

	if configuration.PublishAuthorizer == nil {
		return status.Error(codes.NotFound, "No PublishAuthorizer configured")
	}
	publishAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(
		configuration.PublishAuthorizer,
		dependenciesGroup,
		grpcClientFactory,
	)
	if err != nil {
		return util.StatusWrap(err, "Failed to create PublishAuthorizer")
	}

	// Handle BEP file uploads over HTTP.
	if configuration.EnableBepFileUpload {
		if router == nil {
			return status.Error(codes.NotFound, "Failed to create BEP upload endpoint. No http server configured")
		}
		bepUploader, err := bepuploader.NewBepUploader(dbClient, configuration, publishAuthorizer, dependenciesGroup, grpcClientFactory, tracerProvider)
		if err != nil {
			return util.StatusWrap(err, "Failed to create BEP file upload handler")
		}
		router.Handle(
			"POST /api/v1/bep/upload",
			bepUploader,
		)
	}

	// Handle the Build Event gRPC Stream.
	if len(configuration.GrpcServers) != 0 {
		buildEventServer, err := bes.NewBuildEventServer(dbClient, configuration, publishAuthorizer, dependenciesGroup, grpcClientFactory, tracerProvider)
		if err != nil {
			return util.StatusWrap(err, "Failed to create BuildEventServer")
		}
		if err := bb_grpc.NewServersFromConfigurationAndServe(
			configuration.GrpcServers,
			func(s go_grpc.ServiceRegistrar) {
				build.RegisterPublishBuildEventServer(s.(*go_grpc.Server), buildEventServer)
			},
			siblingsGroup,
			grpcClientFactory,
		); err != nil {
			return util.StatusWrap(err, "gRPC server failure")
		}
	}
	return nil
}
