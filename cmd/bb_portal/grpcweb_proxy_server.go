package main

import (
	"log"
	"slices"

	"github.com/buildbarn/bb-portal/internal/api/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	bb_http "github.com/buildbarn/bb-storage/pkg/http"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	go_grpc "google.golang.org/grpc"
)

func registerAndStartServer(
	configuration *bb_portal.GrpcWebProxyConfiguration,
	siblingsGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	metricsName string,
	registerServer func(*go_grpc.Server, go_grpc.ClientConnInterface),
) {
	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(configuration.Client)
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}

	options := []grpcweb.Option{
		grpcweb.WithOriginFunc(func(origin string) bool { return slices.Contains(configuration.AllowedOrigins, origin) }),
	}

	grpcServer := go_grpc.NewServer()
	registerServer(grpcServer, grpcClient)
	grpcWebServer := grpcweb.WrapServer(grpcServer, options...)
	bb_http.NewServersFromConfigurationAndServe(
		configuration.HttpServers,
		bb_http.NewMetricsHandler(grpcWebServer, metricsName),
		siblingsGroup,
	)
}

// StartGrpcWebProxyServer initializes and starts multiple gRPC web proxy servers based on the provided configuration.
// It registers and starts servers for BuildQueueStateProxy, ActionCacheProxy, ContentAddressableStorageProxy, and InitialSizeClassCacheProxy.
//
// Parameters:
//   - configuration: A pointer to the ApplicationConfiguration which contains the settings for each proxy server.
//   - instanceNameAuthorizer: A auth.Authorizer that checks that only requests with an approved instance name are forwarded.
//   - siblingsGroup: A program.Group that manages the lifecycle of the servers.
//   - grpcClientFactory: A factory for creating gRPC clients.
//
// Each proxy server is registered and started with its respective configuration and implementation.
func StartGrpcWebProxyServer(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
) {
	if configuration.InstanceNameAuthorizer == nil {
		log.Printf("Did not start gRPC-web proxy because InstanceNameAuthorizer is not configured")
		return
	}

	instanceNameAuthorizer, err := auth.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer)
	if err != nil {
		log.Fatalf("Failed to create InstanceNameAuthorizer: %v", err)
	}

	if configuration.BuildQueueStateProxy != nil {
		registerAndStartServer(
			configuration.BuildQueueStateProxy,
			siblingsGroup,
			grpcClientFactory,
			"BuildQueueStateProxy",
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := buildqueuestate.NewBuildQueueStateClient(grpcClient)
				buildqueuestate.RegisterBuildQueueStateServer(grpcServer, buildqueuestateproxy.NewBuildQueueStateServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not start BuildQueueState proxy because BuildQueueStateProxy is not configured")
	}
}
