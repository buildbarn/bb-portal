package archive

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bazelbuild/remote-apis-sdks/go/pkg/digest"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

// LocalFileArchiver Load file archiver struct.
type LocalFileArchiver struct {
	blobArchiveFolder string
}

// NewLocalFileArchiver Constrctor for file archiver.
func NewLocalFileArchiver(blobArchiveFolder string) LocalFileArchiver {
	return LocalFileArchiver{blobArchiveFolder: blobArchiveFolder}
}

// ArchiveBlob Archive Blob function.
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

// A function to archive a blob.
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
