package processing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/bazelbuild/remote-apis-sdks/go/pkg/digest"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

var errNoArchiver = errors.New("no archiver registered")

type BlobArchiver interface {
	ArchiveBlob(ctx context.Context, blobURI detectors.BlobURI) ent.Blob
}

type BlobMultiArchiver struct {
	archivers map[string]BlobArchiver
}

func NewBlobMultiArchiver() BlobMultiArchiver {
	return BlobMultiArchiver{
		archivers: map[string]BlobArchiver{},
	}
}

func (ma *BlobMultiArchiver) RegisterArchiver(schema string, archiver BlobArchiver) {
	ma.archivers[schema] = archiver
}

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

type LocalFileArchiver struct {
	blobArchiveFolder string
}

func NewLocalFileArchiver(blobArchiveFolder string) LocalFileArchiver {
	return LocalFileArchiver{blobArchiveFolder: blobArchiveFolder}
}

func (lfa LocalFileArchiver) ArchiveBlob(_ context.Context, blobURI detectors.BlobURI) ent.Blob {
	b, err := lfa.archiveBlob(blobURI)
	if err != nil {
		return ent.Blob{
			URI:             string(blobURI),
			ArchivingStatus: blob.ArchivingStatusFAILED,
			Reason:          err.Error(),
		}
	}
	return ent.Blob{
		URI:             string(blobURI),
		ArchivingStatus: blob.ArchivingStatusSUCCESS,
		ArchiveURL:      b.ArchiveURL,
	}
}

func (lfa LocalFileArchiver) archiveBlob(blobURI detectors.BlobURI) (*ent.Blob, error) {
	sourcePath := strings.TrimPrefix(string(blobURI), "file://")
	d, err := digest.NewFromFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create digest for path %s: %w", sourcePath, err)
	}
	// Avoid using path.Join() as it removes the "./" prefix for a relative path.
	destPath := lfa.blobArchiveFolder + "/" + d.Hash + "-" + strconv.FormatInt(d.Size, 10)
	destPath = strings.ReplaceAll(destPath, "//", "/")

	var source *os.File
	if source, err = os.Open(sourcePath); err != nil {
		return nil, fmt.Errorf("failed to open source: %w", err)
	}
	defer source.Close()

	var dest *os.File
	if dest, err = os.Create(destPath); err != nil {
		return nil, fmt.Errorf("failed to create destination %s: %w", destPath, err)
	}
	defer dest.Close()

	if _, err = io.Copy(dest, source); err != nil {
		return nil, fmt.Errorf("failed to copy from source to destination %s: %w", destPath, err)
	}
	return &ent.Blob{
		URI:        destPath,
		SizeBytes:  d.Size,
		ArchiveURL: destPath,
	}, nil
}
