package servefiles

import (
	"bufio"
	"context"
	"io"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-remote-execution/pkg/builder"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
)

var digestFunctionStrings = map[string]remoteexecution.DigestFunction_Value{}

type contextKey string

const varKey contextKey = "routerVars"

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

// Dispatcher dispatches requests to the appropriate
// handler based on the URL path
func Dispatcher(serveFilesService *FileServerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if matchRegex := rxFile.FindStringSubmatch(path); matchRegex != nil {
			r = r.WithContext(context.WithValue(r.Context(), varKey, map[string]string{
				"instanceName":   matchRegex[1],
				"digestFunction": matchRegex[2],
				"hash":           matchRegex[3],
				"sizeBytes":      matchRegex[4],
				"name":           matchRegex[5],
			}))
			serveFilesService.HandleFile(w, r)
			return
		}

		if m := rxCommand.FindStringSubmatch(path); m != nil {
			r = r.WithContext(context.WithValue(r.Context(), varKey, map[string]string{
				"instanceName":   m[1],
				"digestFunction": m[2],
				"hash":           m[3],
				"sizeBytes":      m[4],
			}))
			serveFilesService.HandleCommand(w, r)
			return
		}

		if m := rxDirectory.FindStringSubmatch(path); m != nil {
			r = r.WithContext(context.WithValue(r.Context(), varKey, map[string]string{
				"instanceName":   m[1],
				"digestFunction": m[2],
				"hash":           m[3],
				"sizeBytes":      m[4],
			}))
			serveFilesService.HandleDirectory(w, r)
			return
		}

		http.NotFound(w, r)
	}
}

func getDigestFromRequest(req *http.Request) (digest.Digest, error) {
	instanceNameStr := strings.TrimSuffix(req.PathValue("instanceName"), "/")
	instanceName, err := digest.NewInstanceName(instanceNameStr)
	if err != nil {
		return digest.BadDigest, util.StatusWrapf(err, "Invalid instance name %#v", instanceNameStr)
	}
	digestFunctionStr := req.PathValue("digestFunction")
	digestFunctionEnum, ok := digestFunctionStrings[digestFunctionStr]
	if !ok {
		return digest.BadDigest, status.Errorf(codes.InvalidArgument, "Unknown digest function %#v", digestFunctionStr)
	}
	digestFunction, err := instanceName.GetDigestFunction(digestFunctionEnum, 0)
	if err != nil {
		return digest.BadDigest, err
	}
	sizeBytes, err := strconv.ParseInt(req.PathValue("sizeBytes"), 10, 64)
	if err != nil {
		return digest.BadDigest, util.StatusWrapf(err, "Invalid blob size %#v", req.PathValue("sizeBytes"))
	}
	return digestFunction.NewDigest(req.PathValue("hash"), sizeBytes)
}

// FileServerService is a service that serves files from the Content
// Addressable Storage (CAS) over HTTP. It also serves shell scripts generated
// from Command messages, and directories as Tarballs.
type FileServerService struct {
	blobAccess              blobstore.BlobAccess
	maximumMessageSizeBytes int
}

// NewFileServerService creates a new ServeFilesService
// with an authorizing CAS if ServeFilesCasConfiguration is configured.
func NewFileServerService(blobAccess blobstore.BlobAccess, maximumMessageSizeBytes int) *FileServerService {
	return &FileServerService{
		blobAccess,
		int(maximumMessageSizeBytes),
	}
}

// HandleFile serves a file from the Content Addressable Storage (CAS) over HTTP.
func (s FileServerService) HandleFile(w http.ResponseWriter, req *http.Request) {
	digest, err := getDigestFromRequest(req)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}

	ctx := common.ExtractContextFromRequest(req)
	r := s.blobAccess.Get(ctx, digest).ToReader()
	defer r.Close()

	// Attempt to read the first chunk of data to see whether we can
	// trigger an error. Only when no error occurs, we start setting
	// response headers.
	var first [4096]byte
	n, err := r.Read(first[:])
	if err != nil && err != io.EOF {
		http.Error(w, "Could not send file", http.StatusInternalServerError)
		return
	}

	urlPath := req.URL.Path
	contentType := ""
	if extensionStartIndex := strings.LastIndex(urlPath, "."); extensionStartIndex != -1 {
		contentType = mime.TypeByExtension(urlPath[extensionStartIndex:])
	}
	if contentType == "" {
		if utf8.ValidString(string(first[:])) {
			contentType = "text/plain; charset=utf-8"
		} else {
			contentType = "application/octet-stream"
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(first[:n])
	io.Copy(w, r)
}

// HandleCommand serves a Command message from the Content Addressable Storage
// (CAS) as a shell script over HTTP.
func (s FileServerService) HandleCommand(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if req.URL.Query().Get("format") != "sh" {
		http.Error(w, "Invalid format. Only supports \"sh\"", http.StatusNotFound)
		return
	}

	digest, err := getDigestFromRequest(req)
	if err != nil {
		http.Error(w, "Digest not found", http.StatusNotFound)
		return
	}
	ctx := common.ExtractContextFromRequest(req)

	commandMessage, err := s.blobAccess.Get(ctx, digest).ToProto(&remoteexecution.Command{}, s.maximumMessageSizeBytes)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	command := commandMessage.(*remoteexecution.Command)

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
