package blobstoreservice

import (
	"net/http"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/servefiles"
	bb_grpcweb "github.com/buildbarn/bb-portal/pkg/grpcweb"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/blobstore/chunk"
	"github.com/buildbarn/bb-storage/pkg/blobstore/grpcservers"
	"github.com/buildbarn/bb-storage/pkg/cas"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	bb_zstd "github.com/buildbarn/bb-storage/pkg/zstd"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/genproto/googleapis/bytestream"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ContentAddressableStorage contains the constituent parts of the
// Content Addressable Storage (CAS), which in this code base is built
// on top of a Chunk Storage (CS) and Chunk Mapping Storage (CMS). The
// gRPC servers use the authorizing variants of the storages, while
// specialized consumers such as the HTTP file server read through
// ReadInfo and apply the authorizer themselves.
type ContentAddressableStorage struct {
	ChunkStorage              blobstore.BlobAccess[*chunk.Chunk]
	ChunkMappingStorage       blobstore.BlobAccess[chunk.Mapping]
	ReadInfo                  cas.ReadInfo
	ReadAuthorizer            auth.Authorizer
	PutChunkMappingAuthorizer auth.Authorizer
}

// BlobAccess contains the BlobAccessInfo for the ActionCache,
// ContentAddressableStorage, InitialSizeClassCache,
// and FileSystemAccessCache
type BlobAccess struct {
	ContentAddressableStorage *ContentAddressableStorage
	ActionCache               *blobstore.BlobAccess[*remoteexecution.ActionResult]
	InitialSizeClassCache     *blobstore.BlobAccess[*iscc.PreviousExecutionStats]
	FileSystemAccessCache     *blobstore.BlobAccess[*fsac.FileSystemAccessProfile]
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
	if contentAddressableStorage := blobAccess.ContentAddressableStorage; contentAddressableStorage != nil {
		readerPutter := cas.NewReaderPutter(
			contentAddressableStorage.ChunkStorage,
			contentAddressableStorage.ChunkMappingStorage,
		)
		remoteexecution.RegisterContentAddressableStorageServer(
			grpcServer,
			grpcservers.NewContentAddressableStorageServer(
				contentAddressableStorage.ChunkStorage,
				contentAddressableStorage.ChunkMappingStorage,
				contentAddressableStorage.ReadInfo.CDCParametersFetcher,
				readerPutter,
				contentAddressableStorage.ReadAuthorizer,
				contentAddressableStorage.PutChunkMappingAuthorizer,
				configuration.MaximumMessageSizeBytes,
				/* maximumChunkCount = */ 0,
			),
		)
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/build.bazel.remote.execution.v2.ContentAddressableStorage/")

		bytestream.RegisterByteStreamServer(
			grpcServer,
			grpcservers.NewByteStreamServer(
				contentAddressableStorage.ChunkStorage,
				contentAddressableStorage.ChunkMappingStorage,
				contentAddressableStorage.ReadInfo.CDCParametersFetcher,
				zstdPool,
				readerPutter,
			),
		)
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/google.bytestream.ByteStream/")

		// Serve files from the Content Addressable Storage (CAS) over HTTP.
		serveFilesService := servefiles.NewFileServerService(
			contentAddressableStorage.ReadInfo,
			cas.NewMessageReader[remoteexecution.Command](
				contentAddressableStorage.ReadInfo,
				int(configuration.MaximumMessageSizeBytes),
			),
			cas.NewMessageReader[remoteexecution.Directory](
				contentAddressableStorage.ReadInfo,
				int(configuration.MaximumMessageSizeBytes),
			),
			contentAddressableStorage.ReadAuthorizer,
		)
		router.HandleFunc("GET /api/v1/servefile/", servefiles.Dispatcher(serveFilesService))
	}

	// Action Cache (AC).
	if blobAccess.ActionCache != nil {
		remoteexecution.RegisterActionCacheServer(grpcServer, grpcservers.NewActionCacheServer(*blobAccess.ActionCache))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/build.bazel.remote.execution.v2.ActionCache/")
	}

	// Initial Size Class Cache (ISCC).
	if blobAccess.InitialSizeClassCache != nil {
		iscc.RegisterInitialSizeClassCacheServer(grpcServer, grpcservers.NewInitialSizeClassCacheServer(*blobAccess.InitialSizeClassCache))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/buildbarn.iscc.InitialSizeClassCache/")
	}

	// File System Access Cache (FSAC).
	if blobAccess.FileSystemAccessCache != nil {
		fsac.RegisterFileSystemAccessCacheServer(grpcServer, grpcservers.NewFileSystemAccessCacheServer(*blobAccess.FileSystemAccessCache))
		bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/buildbarn.fsac.FileSystemAccessCache/")
	}
	return nil
}
