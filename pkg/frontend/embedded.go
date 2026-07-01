package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
)

//go:embed all:embedded_frontend
var embeddedFiles embed.FS

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

	router.PathPrefix("/").Handler(handler)
	return nil
}
