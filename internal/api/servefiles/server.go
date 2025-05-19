package servefiles

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-remote-execution/pkg/builder"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
)

var digestFunctionStrings = map[string]remoteexecution.DigestFunction_Value{}

func init() {
	for _, digestFunction := range digest.SupportedDigestFunctions {
		digestFunctionStrings[strings.ToLower(digestFunction.String())] = digestFunction
	}
}

func getDigestFromRequest(req *http.Request) (digest.Digest, error) {
	vars := mux.Vars(req)
	instanceNameStr := strings.TrimSuffix(vars["instanceName"], "/")
	instanceName, err := digest.NewInstanceName(instanceNameStr)
	if err != nil {
		return digest.BadDigest, util.StatusWrapf(err, "Invalid instance name %#v", instanceNameStr)
	}
	digestFunctionStr := vars["digestFunction"]
	digestFunctionEnum, ok := digestFunctionStrings[digestFunctionStr]
	if !ok {
		return digest.BadDigest, status.Errorf(codes.InvalidArgument, "Unknown digest function %#v", digestFunctionStr)
	}
	digestFunction, err := instanceName.GetDigestFunction(digestFunctionEnum, 0)
	if err != nil {
		return digest.BadDigest, err
	}
	sizeBytes, err := strconv.ParseInt(vars["sizeBytes"], 10, 64)
	if err != nil {
		return digest.BadDigest, util.StatusWrapf(err, "Invalid blob size %#v", vars["sizeBytes"])
	}
	return digestFunction.NewDigest(vars["hash"], sizeBytes)
}

// FileServerService is a service that serves files from the Content
// Addressable Storage (CAS) over HTTP. It also serves shell scripts generated
// from Command messages, and directories as Tarballs.
type FileServerService struct {
	contentAddressableStorage blobstore.BlobAccess
	maximumMessageSizeBytes   int
}

// NewFileServerService creates a new ServeFilesService
// with an authorizing CAS if ServeFilesCasConfiguration is configured.
func NewFileServerService(blobAccess blobstore.BlobAccess, instanceNameAuthorizer auth.Authorizer, maximumMessageSizeBytes int) *FileServerService {
	return &FileServerService{
		blobstore.NewAuthorizingBlobAccess(blobAccess, instanceNameAuthorizer, nil, nil),
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
	r := s.contentAddressableStorage.Get(ctx, digest).ToReader()
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

	w.Header().Set("Content-Length", strconv.FormatInt(digest.GetSizeBytes(), 10))
	if utf8.ValidString(string(first[:])) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
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

	commandMessage, err := s.contentAddressableStorage.Get(ctx, digest).ToProto(&remoteexecution.Command{}, s.maximumMessageSizeBytes)
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
