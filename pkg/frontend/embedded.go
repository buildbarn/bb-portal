package frontend

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"github.com/klauspost/compress/zstd"
)

//go:embed all:embedded_frontend
var embeddedFiles embed.FS

var zstdPool = sync.Pool{
	New: func() any {
		w, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
		return w
	},
}

type compressResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (cw *compressResponseWriter) Write(b []byte) (int, error) {
	return cw.Writer.Write(b)
}

func cacheControlMiddleware(sourceFS fs.FS, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		_, err := fs.Stat(sourceFS, path)
		if err == nil && path != "index.html" {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}
		next.ServeHTTP(w, r)
	})
}

func zstdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "zstd") {
			next.ServeHTTP(w, r)
			return
		}

		zw := zstdPool.Get().(*zstd.Encoder)
		zw.Reset(w)
		defer func() {
			zw.Close()
			zstdPool.Put(zw)
		}()

		w.Header().Set("Content-Encoding", "zstd")
		w.Header().Add("Vary", "Accept-Encoding")
		w.Header().Del("Content-Length")

		crw := &compressResponseWriter{ResponseWriter: w, Writer: zw}
		next.ServeHTTP(crw, r)
	})
}

func setupEmbeddedHandler(router *mux.Router, frontendConfig *frontend.PortalFrontendConfiguration) error {
	embeddedFrontendFS, err := fs.Sub(embeddedFiles, "embedded_frontend")
	if err != nil {
		return util.StatusWrap(err, "Failed to read embedded files")
	}

	spaFS, err := newSpaFS(embeddedFrontendFS, frontendConfig)
	if err != nil {
		return util.StatusWrap(err, "Failed to create SPA file system")
	}

	var handler http.Handler = http.FileServerFS(spaFS)
	handler = cacheControlMiddleware(embeddedFrontendFS, handler)
	handler = zstdMiddleware(handler)

	router.PathPrefix("/").Handler(handler)
	return nil
}
