package blobstoreservice

import (
	"net/http"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/servefiles"
	bb_grpcweb "github.com/buildbarn/bb-portal/pkg/grpcweb"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	blobstore_configuration "github.com/buildbarn/bb-storage/pkg/blobstore/configuration"
	"github.com/buildbarn/bb-storage/pkg/blobstore/grpcservers"
	"github.com/buildbarn/bb-storage/pkg/digest"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/genproto/googleapis/bytestream"
	go_grpc "google.golang.org/grpc"
)

// NewBlobstoreService initializes and configures a gRPC-Web proxy server the
// ActionCache, ContentAddressableStorage, InitialSizeClassCache, and
// FileSystemAccessCache services, as well as serving files from the Content
// Addressable Storage. It registers all routes it handles with the provided
// router.
func NewBlobstoreService(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	instanceNameAuthorizer auth.Authorizer,
	router *mux.Router,
) error {
	browserConfig := configuration.BrowserServiceConfiguration
	if browserConfig == nil {
		return nil
	}

	// Authorizer used to deny all write requests.
	denyAuthorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return false })

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)

	// Content Addressable Storage (CAS).
	var contentAddressableStorageInfo *blobstore_configuration.BlobAccessInfo
	if browserConfig.ContentAddressableStorage != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfig.ContentAddressableStorage,
			blobstore_configuration.NewCASBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Content Addressable Storage")
		}
		// Add the instanceNameAuthorizer to the blobAccess and make it readonly. BB-portal should not have write access.
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		remoteexecution.RegisterContentAddressableStorageServer(grpcServer, grpcservers.NewContentAddressableStorageServer(blobAccess, configuration.MaximumMessageSizeBytes))
		bytestream.RegisterByteStreamServer(grpcServer, grpcservers.NewByteStreamServer(blobAccess, 1<<16))
		router.PathPrefix(bb_grpcweb.GrpcWebEndpointPrefix + "/google.bytestream.ByteStream/").Handler(http.StripPrefix(bb_grpcweb.GrpcWebEndpointPrefix, grpcWebServer))

		// Serve files from the Content Addressable Storage (CAS) over HTTP.
		serveFilesService := servefiles.NewFileServerService(
			blobAccess,
			int(configuration.MaximumMessageSizeBytes),
		)
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/file/{hash}-{sizeBytes}/{name}", serveFilesService.HandleFile).Methods("GET")
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/command/{hash}-{sizeBytes}/", serveFilesService.HandleCommand).Methods("GET")
		router.HandleFunc("/api/v1/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/directory/{hash}-{sizeBytes}/", serveFilesService.HandleDirectory).Methods("GET")

		contentAddressableStorageInfo = &info
	}

	// Action Cache (AC).
	if browserConfig.ActionCache != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfig.ActionCache,
			blobstore_configuration.NewACBlobAccessCreator(contentAddressableStorageInfo, grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Action Cache")
		}
		// Add the instanceNameAuthorizer to the blobAccess and make it readonly. BB-portal should not have write access.
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		remoteexecution.RegisterActionCacheServer(grpcServer, grpcservers.NewActionCacheServer(blobAccess, int(configuration.MaximumMessageSizeBytes)))
		router.PathPrefix(bb_grpcweb.GrpcWebEndpointPrefix + "/build.bazel.remote.execution.v2.ActionCache/").Handler(http.StripPrefix(bb_grpcweb.GrpcWebEndpointPrefix, grpcWebServer))
	}

	// Initial Size Class Cache (ISCC).
	if browserConfig.InitialSizeClassCache != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfig.InitialSizeClassCache,
			blobstore_configuration.NewISCCBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Initial Size Class Cache")
		}
		// Add the instanceNameAuthorizer to the blobAccess and make it readonly. BB-portal should not have write access.
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		iscc.RegisterInitialSizeClassCacheServer(grpcServer, grpcservers.NewInitialSizeClassCacheServer(blobAccess, int(configuration.MaximumMessageSizeBytes)))
		router.PathPrefix(bb_grpcweb.GrpcWebEndpointPrefix + "/buildbarn.iscc.InitialSizeClassCache/").Handler(http.StripPrefix(bb_grpcweb.GrpcWebEndpointPrefix, grpcWebServer))
	}

	// File System Access Cache (FSAC).
	if browserConfig.FileSystemAccessCache != nil {
		info, err := blobstore_configuration.NewBlobAccessFromConfiguration(
			dependenciesGroup,
			browserConfig.FileSystemAccessCache,
			blobstore_configuration.NewFSACBlobAccessCreator(grpcClientFactory, int(configuration.MaximumMessageSizeBytes)),
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create File System Access Cache")
		}
		// Add the instanceNameAuthorizer to the blobAccess and make it readonly. BB-portal should not have write access.
		blobAccess := blobstore.NewAuthorizingBlobAccess(info.BlobAccess, instanceNameAuthorizer, denyAuthorizer, denyAuthorizer)
		fsac.RegisterFileSystemAccessCacheServer(grpcServer, grpcservers.NewFileSystemAccessCacheServer(blobAccess, int(configuration.MaximumMessageSizeBytes)))
		router.PathPrefix(bb_grpcweb.GrpcWebEndpointPrefix + "/buildbarn.fsac.FileSystemAccessCache/").Handler(http.StripPrefix(bb_grpcweb.GrpcWebEndpointPrefix, grpcWebServer))
	}

	return nil
}
