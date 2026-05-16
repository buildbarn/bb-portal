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
func downloadURLFor(uri, digest, digestFunction string, sizeBytes int64, name string) string {
	if uri == "" {
		return ""
	}
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return uri
	}
	if !strings.HasPrefix(uri, "bytestream://") {
		return ""
	}
	if digest == "" || sizeBytes == 0 || digestFunction == "" {
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
	fmt.Fprintf(&b, "blobs/%s/file/%s-%d/%s",
		strings.ToLower(digestFunction), digest, sizeBytes, name)
	return b.String()
}
