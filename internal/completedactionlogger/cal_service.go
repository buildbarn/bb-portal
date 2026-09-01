package completedactionlogger

import (
	"log"

	"github.com/buildbarn/bb-portal/internal/api/grpc/cal"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/completedactionlogger"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewCompletedActionLoggerServiceAndServe creates a new service that in turn creates
// everything needed in the backend to handle the Completed Action Logger,
// allows receiving completed actions logs
func NewCompletedActionLoggerServiceAndServe(
	configuration *bb_portal.CompletedActionLoggerService,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	cas *blobstore.BlobAccess,
	dbClient database.Client,
) error {
	if dbClient == nil {
		return status.Error(codes.NotFound, "No Database configured")
	}
	if cas == nil {
		return status.Error(codes.NotFound, "No CAS configured")
	}

	if configuration.LogAuthorizer == nil {
		return status.Error(codes.NotFound, "No Instancename Authorizer configured")
	}
	logAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(
		configuration.LogAuthorizer,
		dependenciesGroup,
		grpcClientFactory,
	)
	if err != nil {
		return util.StatusWrap(err, "Failed to create Instancename Authorizer ")
	}

	// Handle the completed action logger gRPC Stream
	actionLoggerServer, err := cal.NewCompletedActionLoggerServer(
		dbClient,
		configuration,
		logAuthorizer,
		*cas,
	)
	if err != nil {
		return util.StatusWrap(err, "Failed to create CompletedActionLoggerServer")
	}

	if err := bb_grpc.NewServersFromConfigurationAndServe(
		configuration.GrpcServers,
		func(s go_grpc.ServiceRegistrar) {
			completedactionlogger.RegisterCompletedActionLoggerServer(s, actionLoggerServer)
		},
		siblingsGroup,
		grpcClientFactory,
	); err != nil {
		log.Fatalf("Failed to start CAL gRPC server: %v", err)
	}

	return nil
}
