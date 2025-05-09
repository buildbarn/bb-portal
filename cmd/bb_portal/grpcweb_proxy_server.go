package main

import (
	"errors"
	"log"
	"net/http"

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
	"github.com/buildbarn/bb-storage/pkg/program"
	grpc_configuration "github.com/buildbarn/bb-storage/pkg/proto/configuration/grpc"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/genproto/googleapis/bytestream"
	go_grpc "google.golang.org/grpc"
)

func initializeGrpcWebServer(
	grpcClientConfig *grpc_configuration.ClientConfiguration,
	grpcClientFactory bb_grpc.ClientFactory,
	grpcServer *go_grpc.Server,
	registerServer func(*go_grpc.Server, go_grpc.ClientConnInterface),
) {
	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(grpcClientConfig)
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	registerServer(grpcServer, grpcClient)
}

// StartGrpcWebProxyServer initializes and starts a gRPC web proxy server based on the provided configuration.
// It registers and starts a server for BuildQueueState, ActionCache, ContentAddressableStorage, and InitialSizeClassCache.
//
// Parameters:
//   - configuration: A pointer to the ApplicationConfiguration which contains the settings for each proxy server.
//   - instanceNameAuthorizer: A auth.Authorizer that checks that only requests with an approved instance name are forwarded.
//   - siblingsGroup: A program.Group that manages the lifecycle of the servers.
//   - grpcClientFactory: A factory for creating gRPC clients.
//
// Each service is registered and started with its respective configuration and implementation.
func StartGrpcWebProxyServer(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
) (http.Handler, error) {
	if configuration.InstanceNameAuthorizer == nil {
		log.Printf("Did not start gRPC-web proxy because InstanceNameAuthorizer is not configured")
		return nil, errors.New("InstanceNameAuthorizer is not configured")
	}

	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, grpcClientFactory)
	if err != nil {
		log.Fatalf("Failed to create InstanceNameAuthorizer: %v", err)
	}

	grpcServer := go_grpc.NewServer()

	if configuration.BuildQueueStateClient != nil {
		initializeGrpcWebServer(
			configuration.BuildQueueStateClient,
			grpcClientFactory,
			grpcServer,
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := buildqueuestate.NewBuildQueueStateClient(grpcClient)
				buildqueuestate.RegisterBuildQueueStateServer(grpcServer, buildqueuestateproxy.NewBuildQueueStateServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not initialize BuildQueueState proxy because BuildQueueStateClient is not configured")
	}

	if configuration.ActionCacheClient != nil {
		initializeGrpcWebServer(
			configuration.ActionCacheClient,
			grpcClientFactory,
			grpcServer,
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := remoteexecution.NewActionCacheClient(grpcClient)
				remoteexecution.RegisterActionCacheServer(grpcServer, actioncacheproxy.NewAcctionCacheServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not initialize ActionCache proxy because ActionCacheClient is not configured")
	}

	if configuration.ContentAddressableStorageClient != nil {
		initializeGrpcWebServer(
			configuration.ContentAddressableStorageClient,
			grpcClientFactory,
			grpcServer,
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := bytestream.NewByteStreamClient(grpcClient)
				bytestream.RegisterByteStreamServer(grpcServer, casproxy.NewCasServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not initialize ContentAddressableStorage proxy because ContentAddressableStorageClient is not configured")
	}

	if configuration.InitialSizeClassCacheClient != nil {
		initializeGrpcWebServer(
			configuration.InitialSizeClassCacheClient,
			grpcClientFactory,
			grpcServer,
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := iscc.NewInitialSizeClassCacheClient(grpcClient)
				iscc.RegisterInitialSizeClassCacheServer(grpcServer, isccproxy.NewIsccServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not initialize InitialSizeClassCache proxy because InitialSizeClassCacheClient is not configured")
	}

	if configuration.FileSystemAccessCacheClient != nil {
		initializeGrpcWebServer(
			configuration.FileSystemAccessCacheClient,
			grpcClientFactory,
			grpcServer,
			func(grpcServer *go_grpc.Server, grpcClient go_grpc.ClientConnInterface) {
				c := fsac.NewFileSystemAccessCacheClient(grpcClient)
				fsac.RegisterFileSystemAccessCacheServer(grpcServer, fsacproxy.NewFsacServerImpl(c, instanceNameAuthorizer))
			},
		)
	} else {
		log.Printf("Did not initialize FileSystemAccessCache proxy because FileSystemAccessCacheClient is not configured")
	}

	return grpcweb.WrapServer(grpcServer), nil
}
