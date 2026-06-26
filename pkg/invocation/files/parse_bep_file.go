package files

import (
	"encoding/hex"
	"net/url"
	"path"
	"strings"

	"github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-storage/pkg/digest"
)

// ParsedBepFile is a BEP file where the URL has been parsed to populate
// missing fields.
type ParsedBepFile struct {
	Path              string `json:"path"`
	SymlinkTargetPath string `json:"symlinkTargetPath"`
	InstanceName      string `json:"instanceName"`
	DigestFunction    int16  `json:"digestFunction"`
	Hash              []byte `json:"hash"`
	SizeBytes         int64  `json:"sizeBytes"`
}

func getDigestFromURI(uri string) digest.Digest {
	// Remove nested protocols if they exist
	if idx := strings.LastIndex(uri, "://"); idx != -1 {
		uri = uri[idx+3:]
	}

	// Re-add a single protocol to make url.Parse happy
	parsed, err := url.Parse("https://" + uri)
	if err != nil {
		return digest.BadDigest
	}

	uriPath := strings.Trim(parsed.Path, "/")
	d, _, err := digest.NewDigestFromByteStreamReadPath(uriPath)
	if err != nil {
		return digest.BadDigest
	}
	return d
}

// ParseBepFile parses a BEP file
func ParseBepFile(file *proto.File) *ParsedBepFile {
	if file == nil {
		return nil
	}

	filePath := path.Join(file.PathPrefix...)
	filePath = path.Join(filePath, file.Name)
	if filePath == "" {
		return nil
	}

	var (
		instanceNameStr   string
		rawDigestFunction = remoteexecution.DigestFunction_UNKNOWN
		hash              = file.Digest
		sizeBytes         = file.Length
	)

	if uri := file.GetUri(); uri != "" {
		if d := getDigestFromURI(uri); d != digest.BadDigest {
			instanceNameStr = d.GetInstanceName().String()
			rawDigestFunction = d.GetDigestFunction().GetEnumValue()

			if hash == "" {
				hash = d.GetHashString()
			}
			if sizeBytes == 0 {
				sizeBytes = d.GetSizeBytes()
			}
		}
	}

	if hash == "" {
		return nil
	}

	// Verify that the instance name is valid
	instanceName, err := digest.NewInstanceName(instanceNameStr)
	if err != nil {
		return nil
	}

	// Derive digest function from the length of the hash
	digestFunction, err := instanceName.GetDigestFunction(rawDigestFunction, len(hash))
	if err != nil {
		return nil
	}

	var hashBytes []byte
	if decodedHashBytes, err := hex.DecodeString(hash); err == nil {
		hashBytes = decodedHashBytes
	}

	return &ParsedBepFile{
		Path:              filePath,
		SymlinkTargetPath: file.GetSymlinkTargetPath(),
		InstanceName:      instanceName.String(),
		DigestFunction:    int16(digestFunction.GetEnumValue()),
		Hash:              hashBytes,
		SizeBytes:         sizeBytes,
	}
}
