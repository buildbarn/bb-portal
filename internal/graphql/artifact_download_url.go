package graphql

import (
	"fmt"
	"net/url"
	"strings"
)

// downloadURLFor returns the bb-portal-served download URL for a file
// with the given URI/digest fields. Returns "" if the file is not
// downloadable through bb-portal.
//
// Patterns:
//
//	bytestream://<host>/<instance>/blobs/<hash>/<size>          (REv2 bytestream)
//	http(s)://...                                                → returned as-is
//	anything else                                                → ""
func downloadURLFor(uri, digest string, sizeBytes int64, name string) string {
	if uri == "" {
		return ""
	}
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return uri
	}
	if !strings.HasPrefix(uri, "bytestream://") {
		return ""
	}
	if digest == "" || sizeBytes == 0 {
		return ""
	}
	// The BEP File proto carries only the hex digest, not which function
	// produced it. The function is implied by the hex digest's length;
	// derive it rather than assuming sha256.
	digestFunction := digestFunctionForHexLength(len(digest))
	if digestFunction == "" {
		return ""
	}
	// Parse the bytestream URI: bytestream://host/<instance...>/blobs/<hash>/<size>
	// Extract <instance> (everything before "/blobs/").
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	p := strings.TrimPrefix(u.Path, "/")
	idx := strings.Index(p, "/blobs/")
	instance := ""
	if idx > 0 {
		instance = p[:idx]
	}
	var b strings.Builder
	b.WriteString("/api/v1/servefile/")
	if instance != "" {
		b.WriteString(instance)
		b.WriteByte('/')
	}
	fmt.Fprintf(&b, "blobs/%s/file/%s-%d/%s", digestFunction, digest, sizeBytes, name)
	return b.String()
}

// digestFunctionForHexLength maps the length of a hex-encoded digest to
// the REv2 digest function that produces it. Returns "" for unrecognized
// lengths. The 64-character case is sha256 by convention (it shares its
// length with other 256-bit functions such as blake3).
func digestFunctionForHexLength(n int) string {
	switch n {
	case 32:
		return "md5"
	case 40:
		return "sha1"
	case 64:
		return "sha256"
	case 96:
		return "sha384"
	case 128:
		return "sha512"
	default:
		return ""
	}
}
