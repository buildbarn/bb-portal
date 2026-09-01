package blobstoreservice

import (
	"net/http"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/servefiles"
	bb_grpcweb "github.com/buildbarn/bb-portal/pkg/grpcweb"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/blobstore/grpcservers"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	bb_zstd "github.com/buildbarn/bb-storage/pkg/zstd"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/genproto/googleapis/bytestream"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BlobAccess contains the BlobAccessInfo for the ActionCache,
// ContentAddressableStorage, InitialSizeClassCache,
// and FileSystemAccessCache
type BlobAccess struct {
	UnauthorizedContentAddressableStorage *blobstore.BlobAccess
	ContentAddressableStorage             *blobstore.BlobAccess
	ActionCache                           *blobstore.BlobAccess
	InitialSizeClassCache                 *blobstore.BlobAccess
	FileSystemAccessCache                 *blobstore.BlobAccess
}

// NewBlobstoreService initializes and configures a gRPC-Web proxy server the
// ActionCache, ContentAddressableStorage, InitialSizeClassCache, and
// FileSystemAccessCache services, as well as serving files from the Content
// Addressable Storage. It registers all routes it handles with the provided
// router.
func NewBlobstoreService(
	configuration *bb_portal.ApplicationConfiguration,
	zstdPool bb_zstd.Pool,
	router *http.ServeMux,
	blobAccess *BlobAccess,
) error {
	if router == nil {
		return status.Error(codes.NotFound, "Failed to create Graphql endpoint. No http server configured")
	}

	if blobAccess.ContentAddressableStorage == nil &&
		blobAccess.ActionCache == nil &&
		blobAccess.InitialSizeClassCache == nil &&
		blobAccess.FileSystemAccessCache == nil {
		return status.Error(codes.InvalidArgument, "No BlobAccess found. Please configure at least one of CAS, AC, ISCC or FSAC")
	}

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)

	// Content Addressable Storage (CAS).
	if blobAccess.ContentAddressableStorage != nil {
		remoteexecution.RegisterContentAddressableStorageServer(grpcServer, grpcservers.NewContentAddressableStorageServer(*blobAccess.ContentAddressableStorage, configuration.MaximumMessageSizeBytes))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/build.bazel.remote.execution.v2.ContentAddressableStorage/")

		bytestream.RegisterByteStreamServer(grpcServer, grpcservers.NewByteStreamServer(*blobAccess.ContentAddressableStorage, 1<<16, zstdPool))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/google.bytestream.ByteStream/")

		// Serve files from the Content Addressable Storage (CAS) over HTTP.
		serveFilesService := servefiles.NewFileServerService(
			*blobAccess.ContentAddressableStorage,
			int(configuration.MaximumMessageSizeBytes),
		)
		router.HandleFunc("GET /api/v1/servefile/", servefiles.Dispatcher(serveFilesService))
	}

	// Action Cache (AC).
	if blobAccess.ActionCache != nil {
		remoteexecution.RegisterActionCacheServer(grpcServer, grpcservers.NewActionCacheServer(*blobAccess.ActionCache, int(configuration.MaximumMessageSizeBytes)))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/build.bazel.remote.execution.v2.ActionCache/")
	}

	// Initial Size Class Cache (ISCC).
	if blobAccess.InitialSizeClassCache != nil {
		iscc.RegisterInitialSizeClassCacheServer(grpcServer, grpcservers.NewInitialSizeClassCacheServer(*blobAccess.InitialSizeClassCache, int(configuration.MaximumMessageSizeBytes)))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/buildbarn.iscc.InitialSizeClassCache/")
	}

	// File System Access Cache (FSAC).
	if blobAccess.FileSystemAccessCache != nil {
		fsac.RegisterFileSystemAccessCacheServer(grpcServer, grpcservers.NewFileSystemAccessCacheServer(*blobAccess.FileSystemAccessCache, int(configuration.MaximumMessageSizeBytes)))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/buildbarn.fsac.FileSystemAccessCache/")
	}
	return nil
}
