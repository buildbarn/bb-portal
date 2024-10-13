package archive

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

// errNoArchiver An error helper.
var errNoArchiver = errors.New("no archiver registered")

// BlobArchiver A blob arhiver interace.
type BlobArchiver interface {
	ArchiveBlob(ctx context.Context, blobURI detectors.BlobURI) ent.Blob
}

// BlobMultiArchiver A blob Multi Archiver.
type BlobMultiArchiver struct {
	archivers map[string]BlobArchiver
}

// NewBlobMultiArchiver A blob multi archiver constructor.
func NewBlobMultiArchiver() BlobMultiArchiver {
	return BlobMultiArchiver{
		archivers: map[string]BlobArchiver{},
	}
}

// RegisterArchiver Regsters an archiver.
func (ma *BlobMultiArchiver) RegisterArchiver(schema string, archiver BlobArchiver) {
	ma.archivers[schema] = archiver
}

// ArchiveBlobs Archives blobs.
func (ma *BlobMultiArchiver) ArchiveBlobs(ctx context.Context, blobURIs []detectors.BlobURI) ([]ent.Blob, error) {
	if len(blobURIs) == 0 {
		return nil, nil
	}
	if len(ma.archivers) == 0 {
		return nil, nil
	}
	uri, err := url.Parse(string(blobURIs[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid blob URI: %s: %w", blobURIs[0], err)
	}
	archiver, ok := ma.archivers[uri.Scheme]
	if !ok {
		return nil, fmt.Errorf("scheme %s: %w", uri.Scheme, errNoArchiver)
	}
	blobs := make([]ent.Blob, 0, len(blobURIs))
	for _, blobURI := range blobURIs {
		b := archiver.ArchiveBlob(ctx, blobURI)
		blobs = append(blobs, b)
	}
	return blobs, nil
}
