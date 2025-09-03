package main

import (
	"log"
	"net/http"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/internal/api/servefiles"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	blobstore_configuration "github.com/buildbarn/bb-storage/pkg/blobstore/configuration"
	"github.com/buildbarn/bb-storage/pkg/blobstore/grpcservers"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	auth_proto "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
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
	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	// Authorizer used to deny all write requests.
	denyAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(&auth_proto.AuthorizerConfiguration{
		Policy: &auth_proto.AuthorizerConfiguration_Deny{},
	}, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)

	// Content Addressable Storage (CAS).
	var contentAddressableStorageInfo *blobstore_configuration.BlobAccessInfo
	if browserConfiguration.ContentAddressableStorage != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfiguration.ContentAddressableStorage,
			blobstore_configuration.NewCASBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Content Addressable Storage")
		}
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		remoteexecution.RegisterContentAddressableStorageServer(grpcServer, grpcservers.NewContentAddressableStorageServer(blobAccess, configuration.MaximumMessageSizeBytes))
		bytestream.RegisterByteStreamServer(grpcServer, grpcservers.NewByteStreamServer(blobAccess, 1<<16))
		router.PathPrefix(grpcWebEndpointPrefix + "/google.bytestream.ByteStream/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))

		// Serve files from the Content Addressable Storage (CAS) over HTTP.
		serveFilesService := servefiles.NewFileServerService(
			blobAccess,
			int(configuration.MaximumMessageSizeBytes),
		)
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/file/{hash}-{sizeBytes}/{name}", serveFilesService.HandleFile).Methods("GET")
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/command/{hash}-{sizeBytes}/", serveFilesService.HandleCommand).Methods("GET")
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/directory/{hash}-{sizeBytes}/", serveFilesService.HandleDirectory).Methods("GET")

		contentAddressableStorageInfo = &info
	} else {
		return status.Error(codes.NotFound, "No ContentAddressableStorage configured")
	}

	// Action Cache (AC).
	if browserConfiguration.ActionCache != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfiguration.ActionCache,
			blobstore_configuration.NewACBlobAccessCreator(contentAddressableStorageInfo, grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Action Cache")
		}
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		remoteexecution.RegisterActionCacheServer(grpcServer, grpcservers.NewActionCacheServer(blobAccess, int(configuration.MaximumMessageSizeBytes)))
		router.PathPrefix(grpcWebEndpointPrefix + "/build.bazel.remote.execution.v2.ActionCache/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start ActionCache service because browserServiceConfiguration.actionCache is not configured")
	}

	// Initial Size Class Cache (ISCC).
	if browserConfiguration.InitialSizeClassCache != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfiguration.InitialSizeClassCache,
			blobstore_configuration.NewISCCBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Initial Size Class Cache")
		}
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		iscc.RegisterInitialSizeClassCacheServer(grpcServer, grpcservers.NewInitialSizeClassCacheServer(blobAccess, int(configuration.MaximumMessageSizeBytes)))
		router.PathPrefix(grpcWebEndpointPrefix + "/buildbarn.iscc.InitialSizeClassCache/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start InitialSizeClassCache service because browserServiceConfiguration.initialSizeClassCache is not configured")
	}

	// File System Access Cache (FSAC).
	if browserConfiguration.FileSystemAccessCache != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfiguration.FileSystemAccessCache,
			blobstore_configuration.NewFSACBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create File System Access Cache")
		}
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		fsac.RegisterFileSystemAccessCacheServer(grpcServer, grpcservers.NewFileSystemAccessCacheServer(blobAccess, int(configuration.MaximumMessageSizeBytes)))
		router.PathPrefix(grpcWebEndpointPrefix + "/buildbarn.fsac.FileSystemAccessCache/").Handler(http.StripPrefix(grpcWebEndpointPrefix, grpcWebServer))
	} else {
		log.Printf("Did not start FileSystemAccessCache service because browserServiceConfiguration.fileSystemAccessCache is not configured")
	}

	return nil
}

// NewGrpcWebSchedulerService initializes and configures a gRPC-Web proxy
// server the BuildQueueState API, and starts an endpoint used for checking
// permissions. It registers all routes it handles with the provided router.
func NewGrpcWebSchedulerService(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *mux.Router,
) error {
	schedulerConfiguration := configuration.SchedulerServiceConfiguration
	if schedulerConfiguration == nil {
		log.Printf("Did not start Scheduler service because schedulerConfiguration is not configured")
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

	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}
	killOperationsAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(schedulerConfiguration.KillOperationsAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create KillOperationsAuthorizer")
	}

	if schedulerConfiguration.BuildQueueStateClient == nil {
		return util.StatusWrap(err, "No buildQueueStateClient configured")
	}
	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(schedulerConfiguration.BuildQueueStateClient, dependenciesGroup)
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
