package schedulerservice

import (
	"net/http"

	bb_grpcweb "github.com/buildbarn/bb-portal/pkg/grpcweb"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewSchedulerService initializes and configures a gRPC-Web proxy server the
// BuildQueueState API, and starts an endpoint used for checking permissions.
// It registers all routes it handles with the provided router.
func NewSchedulerService(
	configuration *bb_portal.SchedulerService,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	instanceNameAuthorizer auth.Authorizer,
	router *mux.Router,
) error {
	if configuration.ListOperationsPageSize <= 0 {
		return status.Error(codes.NotFound, "No ListOperationsPageSize configured (or it is set to 0)")
	}

	if configuration.KillOperationsAuthorizer == nil {
		return status.Error(codes.NotFound, "No KillOperationsAuthorizer configured")
	}
	killOperationsAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.KillOperationsAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create KillOperationsAuthorizer")
	}

	if configuration.BuildQueueStateClient == nil {
		return util.StatusWrap(err, "No BuildQueueStateClient configured")
	}
	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(configuration.BuildQueueStateClient, dependenciesGroup)
	if err != nil {
		return util.StatusWrap(err, "Failed to create gRPC client for BuildQueueState")
	}
	buildQueueStateClient := buildqueuestate.NewBuildQueueStateClient(grpcClient)
	buildQueueStateServer := buildqueuestateproxy.NewBuildQueueStateServerImpl(buildQueueStateClient, instanceNameAuthorizer, killOperationsAuthorizer, configuration.ListOperationsPageSize)

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	buildqueuestate.RegisterBuildQueueStateServer(grpcServer, buildQueueStateServer)

	router.PathPrefix(bb_grpcweb.GrpcWebEndpointPrefix + "/buildbarn.buildqueuestate.BuildQueueState/").Handler(http.StripPrefix(bb_grpcweb.GrpcWebEndpointPrefix, grpcWebServer))
	router.HandleFunc("/api/v1/checkPermissions/killOperation/{operationName}", buildQueueStateServer.CheckKillOperationAuthorization).Methods("GET")

	return nil
}
