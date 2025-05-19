package main

import (
	"log"
	"net/http"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/actioncacheproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/casproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/fsacproxy"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/isccproxy"
	"github.com/buildbarn/bb-portal/internal/api/servefiles"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	blobstore_configuration "github.com/buildbarn/bb-storage/pkg/blobstore/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/genproto/googleapis/bytestream"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const grpcWebEndpointPrefix = "/api/v1/grpcweb"

// NewGrpcWebBrowserService initializes and configures a gRPC-Web proxy server
// the ActionCache, ContentAddressableStorage, InitialSizeClassCache, and
// FileSystemAccessCache services, as well as serving files from the Content
// Addressable Storage. It registers all routes it handles with the provided
// router.
func NewGrpcWebBrowserService(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *mux.Router,
) error {
	browserConfiguration := configuration.BrowserServiceConfiguration
	if browserConfiguration == nil {
		log.Printf("Did not start gRPC-web Browser proxy because browserConfiguration is not configured")
		return nil
	}

	if configuration.InstanceNameAuthorizer == nil {
		return status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}
	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)

	if browserConfiguration.ActionCacheClient != nil {
		grpcClient, err := grpcClientFactory.NewClientFromConfiguration(browserConfiguration.ActionCacheClient)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC client for ActionCache")
		}
		actionCacheClient := remoteexecution.NewActionCacheClient(grpcClient)
		actionCacheServer := actioncacheproxy.NewAcctionCacheServerImpl(actionCacheClient, instanceNameAuthorizer)
		remoteexecution.RegisterActionCacheServer(grpcServer, actionCacheServer)
		router.PathPrefix(grpcWebEndpointPrefix + "/build.bazel.remote.execution.v2.ActionCache/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start ActionCache proxy because actionCacheClient is not configured")
	}

	if browserConfiguration.ContentAddressableStorageClient != nil {
		grpcClient, err := grpcClientFactory.NewClientFromConfiguration(browserConfiguration.ContentAddressableStorageClient)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC client for ContentAddressableStorage")
		}
		casClient := bytestream.NewByteStreamClient(grpcClient)
		casServer := casproxy.NewCasServerImpl(casClient, instanceNameAuthorizer)
		bytestream.RegisterByteStreamServer(grpcServer, casServer)
		router.PathPrefix(grpcWebEndpointPrefix + "/google.bytestream.ByteStream/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start ContentAddressableStorage proxy because contentAddressableStorageClient is not configured")
	}

	if browserConfiguration.InitialSizeClassCacheClient != nil {
		grpcClient, err := grpcClientFactory.NewClientFromConfiguration(browserConfiguration.ContentAddressableStorageClient)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC client for InitialSizeClassCache")
		}
		isccClient := iscc.NewInitialSizeClassCacheClient(grpcClient)
		isccServer := isccproxy.NewIsccServerImpl(isccClient, instanceNameAuthorizer)
		iscc.RegisterInitialSizeClassCacheServer(grpcServer, isccServer)
		router.PathPrefix(grpcWebEndpointPrefix + "/buildbarn.iscc.InitialSizeClassCache/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start InitialSizeClassCache proxy because initialSizeClassCacheClient is not configured")
	}

	if browserConfiguration.FileSystemAccessCacheClient != nil {
		grpcClient, err := grpcClientFactory.NewClientFromConfiguration(browserConfiguration.ContentAddressableStorageClient)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC client for FileSystemAccessCache")
		}
		fsacClient := fsac.NewFileSystemAccessCacheClient(grpcClient)
		fsacServer := fsacproxy.NewFsacServerImpl(fsacClient, instanceNameAuthorizer)
		fsac.RegisterFileSystemAccessCacheServer(grpcServer, fsacServer)
		router.PathPrefix(grpcWebEndpointPrefix + "/buildbarn.fsac.FileSystemAccessCache/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start FileSystemAccessCache proxy because fileSystemAccessCacheClient is not configured")
	}

	if browserConfiguration.ServeFilesCasConfiguration != nil {
		contentAddressableStorage, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfiguration.ServeFilesCasConfiguration,
			blobstore_configuration.NewCASBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create CAS blob access for serving files")
		}
		serveFilesService := servefiles.NewFileServerService(
			contentAddressableStorage.BlobAccess,
			instanceNameAuthorizer,
			int(configuration.MaximumMessageSizeBytes),
		)
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/file/{hash}-{sizeBytes}/{name}", serveFilesService.HandleFile).Methods("GET")
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/command/{hash}-{sizeBytes}/", serveFilesService.HandleCommand).Methods("GET")
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/directory/{hash}-{sizeBytes}/", serveFilesService.HandleDirectory).Methods("GET")
	} else {
		log.Printf("Did not start serving files from Content Addressable Storage because serveFilesCasConfiguration is not configured")
	}

	return nil
}

// NewGrpcWebSchedulerService initializes and configures a gRPC-Web proxy
// server the BuildQueueState API, and starts an endpoint used for checking
// permissions. It registers all routes it handles with the provided router.
func NewGrpcWebSchedulerService(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *mux.Router,
) error {
	schedulerConfiguration := configuration.SchedulerServiceConfiguration
	if schedulerConfiguration == nil {
		log.Printf("Did not start gRPC-web Scheduler proxy because schedulerConfiguration is not configured")
		return nil
	}
	if configuration.InstanceNameAuthorizer == nil {
		return status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}
	if schedulerConfiguration.KillOperationsAuthorizer == nil {
		return status.Error(codes.NotFound, "No KillOperationsAuthorizer configured")
	}
	if schedulerConfiguration.ListOperationsPageSize <= 0 {
		return status.Error(codes.NotFound, "No ListOperationsPageSize configured (or it is set to 0)")
	}

	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}
	killOperationsAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(schedulerConfiguration.KillOperationsAuthorizer, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create KillOperationsAuthorizer")
	}

	if schedulerConfiguration.BuildQueueStateClient == nil {
		return util.StatusWrap(err, "No buildQueueStateClient configured")
	}
	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(schedulerConfiguration.BuildQueueStateClient)
	if err != nil {
		return util.StatusWrap(err, "Failed to create gRPC client for BuildQueueState")
	}
	buildQueueStateClient := buildqueuestate.NewBuildQueueStateClient(grpcClient)
	buildQueueStateServer := buildqueuestateproxy.NewBuildQueueStateServerImpl(buildQueueStateClient, instanceNameAuthorizer, killOperationsAuthorizer, schedulerConfiguration.ListOperationsPageSize)

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	buildqueuestate.RegisterBuildQueueStateServer(grpcServer, buildQueueStateServer)

	router.PathPrefix(grpcWebEndpointPrefix + "/buildbarn.buildqueuestate.BuildQueueState/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	router.HandleFunc("/api/v1/checkPermissions/killOperation/{operationName}", buildQueueStateServer.CheckKillOperationAuthorization).Methods("GET")

	return nil
}
