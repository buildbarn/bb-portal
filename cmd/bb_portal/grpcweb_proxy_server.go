package main

import (
	"log"
	"slices"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/actioncacheproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/casproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/fsacproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/isccproxy"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	bb_http "github.com/buildbarn/bb-storage/pkg/http"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/genproto/googleapis/bytestream"
	go_grpc "google.golang.org/grpc"
)

func registerAndStartServer(
	globalConfig *bb_portal.ApplicationConfiguration,
	proxyConfig *bb_portal.GrpcWebProxyConfiguration,
	siblingsGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	metricsName string,
	registerServer func(*go_grpc.Server, go_grpc.ClientConnInterface),
) {
	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(proxyConfig.Client)
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}

	options := []grpcweb.Option{
		grpcweb.WithOriginFunc(func(origin string) bool {
			return slices.Contains(globalConfig.AllowedOrigins, origin) || slices.Contains(globalConfig.AllowedOrigins, "*")
		}),
	}

	grpcServer := go_grpc.NewServer()
	registerServer(grpcServer, grpcClient)
	grpcWebServer := grpcweb.WrapServer(grpcServer, options...)
	bb_http.NewServersFromConfigurationAndServe(
		proxyConfig.HttpServers,
		bb_http.NewMetricsHandler(grpcWebServer, metricsName),
		siblingsGroup,
		grpcClientFactory,
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

	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, grpcClientFactory)
	if err != nil {
		log.Fatalf("Failed to create InstanceNameAuthorizer: %v", err)
	}

	if configuration.BuildQueueStateProxy != nil {
		registerAndStartServer(
			configuration,
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

	if configuration.ActionCacheProxy != nil {
		registerAndStartServer(
			configuration,
			configuration.ActionCacheProxy,
			siblingsGroup,
			grpcClientFactory,
			"ActionCacheProxy",
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := remoteexecution.NewActionCacheClient(grpcClient)
				remoteexecution.RegisterActionCacheServer(grpcServer, actioncacheproxy.NewAcctionCacheServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not start ActionCache proxy because ActionCacheProxy is not configured")
	}

	if configuration.ContentAddressableStorageProxy != nil {
		registerAndStartServer(
			configuration,
			configuration.ContentAddressableStorageProxy,
			siblingsGroup,
			grpcClientFactory,
			"ContentAddressableStorageProxy",
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := bytestream.NewByteStreamClient(grpcClient)
				bytestream.RegisterByteStreamServer(grpcServer, casproxy.NewCasServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not start ContentAddressableStorage proxy because ContentAddressableStorageProxy is not configured")
	}

	if configuration.InitialSizeClassCacheProxy != nil {
		registerAndStartServer(
			configuration,
			configuration.InitialSizeClassCacheProxy,
			siblingsGroup,
			grpcClientFactory,
			"InitialSizeClassCacheProxy",
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := iscc.NewInitialSizeClassCacheClient(grpcClient)
				iscc.RegisterInitialSizeClassCacheServer(grpcServer, isccproxy.NewIsccServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not start InitialSizeClassCache proxy because InitialSizeClassCacheProxy is not configured")
	}

	if configuration.FileSystemAccessCacheProxy != nil {
		registerAndStartServer(
			configuration,
			configuration.FileSystemAccessCacheProxy,
			siblingsGroup,
			grpcClientFactory,
			"FileSystemAccessCacheProxy",
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := fsac.NewFileSystemAccessCacheClient(grpcClient)
				fsac.RegisterFileSystemAccessCacheServer(grpcServer, fsacproxy.NewFsacServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not start FileSystemAccessCache proxy because FileSystemAccessCacheProxy is not configured")
	}
}
