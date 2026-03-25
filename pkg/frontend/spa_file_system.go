package frontend

import (
	"bytes"
	"io"
	"io/fs"
	"time"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"github.com/buildbarn/bb-storage/pkg/util"
)

type modifiedFile struct {
	*bytes.Reader
	info *modifiedFileInfo
}

// Close is a no-op
func (m *modifiedFile) Close() error { return nil }

// Stat returns a custom FileInfo
func (m *modifiedFile) Stat() (fs.FileInfo, error) {
	return m.info, nil
}

type modifiedFileInfo struct {
	fs.FileInfo
	newSize    int64
	newModTime time.Time
}

func (m *modifiedFileInfo) Size() int64        { return m.newSize }
func (m *modifiedFileInfo) ModTime() time.Time { return m.newModTime }

type spaFS struct {
	sourceFS  fs.FS
	indexData []byte
	indexInfo *modifiedFileInfo
}

func newSpaFS(sourceFS fs.FS, frontendConfig *frontend.PortalFrontendConfiguration) (*spaFS, error) {
	indexFile, err := sourceFS.Open("index.html")
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to get index.html file from embedded frontend")
	}
	defer indexFile.Close()

	indexFileInfo, err := indexFile.Stat()
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to get index.html file info from embedded frontend")
	}

	indexContent, err := io.ReadAll(indexFile)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to read index.html file from embedded frontend")
	}

	if err := validateFrontendConfig(frontendConfig); err != nil {
		return nil, util.StatusWrap(err, "Error validating frontend config")
	}

	newIndexContent, err := injectFrontendConfigScript(indexContent, frontendConfig)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to insert frontend config into index.html")
	}

	return &spaFS{
		sourceFS:  sourceFS,
		indexData: newIndexContent,
		indexInfo: &modifiedFileInfo{
			FileInfo:   indexFileInfo,
			newSize:    int64(len(newIndexContent)),
			newModTime: time.Now(),
		},
	}, nil
}

func (f *spaFS) serveIndexHTML() fs.File {
	return &modifiedFile{
		Reader: bytes.NewReader(f.indexData),
		info:   f.indexInfo,
	}
}

// Open opens a file in the single page application FS
func (f *spaFS) Open(name string) (fs.File, error) {
	if name == "index.html" {
		return f.serveIndexHTML(), nil
	}

	file, err := f.sourceFS.Open(name)
	if err != nil {
		return f.serveIndexHTML(), nil
	}
	return file, nil
}
