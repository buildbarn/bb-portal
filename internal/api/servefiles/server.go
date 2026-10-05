package servefiles

import (
	"bufio"
	"context"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-remote-execution/pkg/builder"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/blobstore/chunk"
	"github.com/buildbarn/bb-storage/pkg/capabilities"
	"github.com/buildbarn/bb-storage/pkg/cas"
	"github.com/buildbarn/bb-storage/pkg/cas/reader"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var digestFunctionStrings = map[string]remoteexecution.DigestFunction_Value{}

var (
	// For /blobs/<digest>/file/<hash>-<size>/<name>
	rxFile = regexp.MustCompile(`^/api/v1/servefile/(.*?/?)blobs/([^/]+)/file/([^/-]+)-([^/]+)/(.*)$`)

	// For /blobs/<digest>/command/<hash>-<size>/
	rxCommand = regexp.MustCompile(`^/api/v1/servefile/(.*?/?)blobs/([^/]+)/command/([^/-]+)-([^/]+)/?$`)

	// For /blobs/<digest>/directory/<hash>-<size>/
	rxDirectory = regexp.MustCompile(`^/api/v1/servefile/(.*?/?)blobs/([^/]+)/directory/([^/-]+)-([^/]+)/?$`)
)

func init() {
	for _, digestFunction := range digest.SupportedDigestFunctions {
		digestFunctionStrings[strings.ToLower(digestFunction.String())] = digestFunction
	}
}

type digestParams struct {
	instanceName   string
	digestFunction string
	hash           string
	sizeBytes      string
}

type handleFileParams struct {
	digestParams
	name string
}

// Dispatcher dispatches requests to the appropriate handler based on the URL
// path
func Dispatcher(serveFilesService *FileServerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if m := rxFile.FindStringSubmatch(path); m != nil {
			params := handleFileParams{
				digestParams: digestParams{
					instanceName:   m[1],
					digestFunction: m[2],
					hash:           m[3],
					sizeBytes:      m[4],
				},
				name: m[5],
			}
			serveFilesService.HandleFile(w, r, params)
			return
		}

		if m := rxCommand.FindStringSubmatch(path); m != nil {
			params := digestParams{
				instanceName:   m[1],
				digestFunction: m[2],
				hash:           m[3],
				sizeBytes:      m[4],
			}
			serveFilesService.HandleCommand(w, r, params)
			return
		}

		if m := rxDirectory.FindStringSubmatch(path); m != nil {
			params := digestParams{
				instanceName:   m[1],
				digestFunction: m[2],
				hash:           m[3],
				sizeBytes:      m[4],
			}
			serveFilesService.HandleDirectory(w, r, params)
			return
		}

		http.NotFound(w, r)
	}
}

func getDigestFromParams(params digestParams) (digest.Digest, error) {
	instanceNameStr := strings.TrimSuffix(params.instanceName, "/")
	instanceName, err := digest.NewInstanceName(instanceNameStr)
	if err != nil {
		return digest.BadDigest, util.StatusWrapf(err, "Invalid instance name %#v", instanceNameStr)
	}
	digestFunctionEnum, ok := digestFunctionStrings[params.digestFunction]
	if !ok {
		return digest.BadDigest, status.Errorf(codes.InvalidArgument, "Unknown digest function %#v", params.digestFunction)
	}
	digestFunction, err := instanceName.GetDigestFunction(digestFunctionEnum, 0)
	if err != nil {
		return digest.BadDigest, err
	}
	sizeBytes, err := strconv.ParseInt(params.sizeBytes, 10, 64)
	if err != nil {
		return digest.BadDigest, util.StatusWrapf(err, "Invalid blob size %#v", params.sizeBytes)
	}
	return digestFunction.NewDigest(params.hash, sizeBytes)
}

// FileServerService is a service that serves files from the Content
// Addressable Storage (CAS) over HTTP. It also serves shell scripts generated
// from Command messages, and directories as Tarballs.
type FileServerService struct {
	chunkMappingStorage  blobstore.BlobAccess[chunk.Mapping]
	chunkBytesReader     reader.Reader[[]byte]
	commandReader        reader.Reader[*remoteexecution.Command]
	directoryReader      reader.Reader[*remoteexecution.Directory]
	cdcParametersFetcher capabilities.CDCParametersFetcher
	authorizer           auth.Authorizer
}

// NewFileServerService creates a new ServeFilesService that reads blobs
// directly from the raw Chunk Storage (CS) and Chunk Mapping Storage
// (CMS). Every request is authorized with the provided authorizer, so
// the underlying storages do not need to be wrapped in
// AuthorizingBlobAccess.
func NewFileServerService(chunkStorage blobstore.BlobAccess[*chunk.Chunk], chunkMappingStorage blobstore.BlobAccess[chunk.Mapping], commandReader reader.Reader[*remoteexecution.Command], directoryReader reader.Reader[*remoteexecution.Directory], cdcParametersFetcher capabilities.CDCParametersFetcher, authorizer auth.Authorizer) *FileServerService {
	return &FileServerService{
		chunkMappingStorage:  chunkMappingStorage,
		chunkBytesReader:     cas.NewChunkBytesReader(chunkStorage),
		commandReader:        commandReader,
		directoryReader:      directoryReader,
		cdcParametersFetcher: cdcParametersFetcher,
		authorizer:           authorizer,
	}
}

// authorizeInstanceName checks whether the request is allowed to access
// the given instance name.
func (s FileServerService) authorizeInstanceName(ctx context.Context, instanceName digest.InstanceName) error {
	if err := auth.AuthorizeSingleInstanceName(ctx, s.authorizer, instanceName); err != nil {
		return util.StatusWrap(err, "Authorization")
	}
	return nil
}

// getChunkMapping returns the ordered chunk digests that compose the
// blob with the given digest. Blobs that fit in a single chunk have no
// chunk mapping in storage; the blob is the chunk itself.
func (s FileServerService) getChunkMapping(ctx context.Context, cdcParams *remoteexecution.RepMaxCdcParams, d digest.Digest) ([]digest.Digest, error) {
	if cas.IsSingleChunk(cdcParams, d) {
		if d.GetSizeBytes() == 0 {
			// The empty blob is always present and yields no
			// data.
			return nil, nil
		}
		// Blobs that fit in a single chunk have no chunk mappings in
		// storage; the blob is the chunk itself. Its contents cannot
		// mismatch its digest, so a read suffices as an existence
		// check.
		if _, err := s.chunkBytesReader.Read(ctx, d); err != nil {
			return nil, err
		}
		return []digest.Digest{d}, nil
	}
	chunkMapping, err := s.chunkMappingStorage.Get(ctx, d)
	if err != nil {
		return nil, err
	}
	return chunkMapping.GetDigests(), nil
}

// HandleFile serves a file from the Content Addressable Storage (CAS) over HTTP.
func (s FileServerService) HandleFile(w http.ResponseWriter, req *http.Request, params handleFileParams) {
	digest, err := getDigestFromParams(params.digestParams)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}

	ctx := common.ExtractContextFromRequest(req)
	if err := s.authorizeInstanceName(ctx, digest.GetInstanceName()); err != nil {
		http.Error(w, "Digest not found", http.StatusForbidden)
		return
	}
	cdcParams, err := s.cdcParametersFetcher.FetchCDCParameters(ctx, digest.GetInstanceName())
	if err != nil {
		http.Error(w, "Could not determine chunking parameters", http.StatusInternalServerError)
		return
	}

	// Determine the chunk mapping before emitting the headers, so that
	// failures can still be reported as an HTTP error.
	chunkDigests, err := s.getChunkMapping(ctx, cdcParams, digest)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}

	urlPath := req.URL.Path
	contentType := ""
	if extensionStartIndex := strings.LastIndex(urlPath, "."); extensionStartIndex != -1 {
		contentType = mime.TypeByExtension(urlPath[extensionStartIndex:])
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	for _, chunkDigest := range chunkDigests {
		chunkBytes, err := s.chunkBytesReader.Read(ctx, chunkDigest)
		if err != nil {
			log.Print(util.StatusWrapf(err, "Failed to fetch chunk %s of blob %s", chunkDigest, digest))
			panic(http.ErrAbortHandler)
		}
		if _, err := w.Write(chunkBytes); err != nil {
			return
		}
	}
}

// HandleCommand serves a Command message from the Content Addressable Storage
// (CAS) as a shell script over HTTP.
func (s FileServerService) HandleCommand(w http.ResponseWriter, req *http.Request, params digestParams) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.URL.Query().Get("format") != "sh" {
		http.Error(w, "Invalid format. Only supports \"sh\"", http.StatusNotFound)
		return
	}

	digest, err := getDigestFromParams(params)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}
	ctx := common.ExtractContextFromRequest(req)
	if err := s.authorizeInstanceName(ctx, digest.GetInstanceName()); err != nil {
		http.Error(w, "Digest not found", http.StatusForbidden)
		return
	}
	command, err := s.commandReader.Read(ctx, digest)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	bw := bufio.NewWriter(w)
	if err := builder.ConvertCommandToShellScript(command, bw); err != nil {
		log.Print(err)
		panic(http.ErrAbortHandler)
	}
	if err := bw.Flush(); err != nil {
		log.Print(err)
		panic(http.ErrAbortHandler)
	}
}
