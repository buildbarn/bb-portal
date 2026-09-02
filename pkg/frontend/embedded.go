package frontend

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
)

// RouteManifest contains the JavaScript and CSS assets for a route
type RouteManifest struct {
	Path string   `json:"path"`
	JS   []string `json:"js"`
	CSS  []string `json:"css"`
}

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

func preloadLinkMiddleware(next http.Handler, embeddedFrontendFS fs.FS) http.Handler {
	var routeManifest []RouteManifest
	routeMapBytes, err := fs.ReadFile(embeddedFrontendFS, "route-manifest.json")
	if err != nil {
		log.Fatalf("Failed to read route-manifest.json: %v", err)
	}

	if err := json.Unmarshal(routeMapBytes, &routeManifest); err != nil {
		log.Fatalf("Failed to parse route-manifest.json: %v", err)
	}

	shadowRouter := mux.NewRouter()
	for _, route := range routeManifest {
		shadowRouter.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			for _, css := range route.CSS {
				w.Header().Add("Link", fmt.Sprintf("</%s>; rel=preload; as=style; crossorigin", css))
			}
			for _, js := range route.JS {
				w.Header().Add("Link", fmt.Sprintf("</%s>; rel=modulepreload; crossorigin", js))
			}
			next.ServeHTTP(w, r)
		}).Methods("GET")
	}

	// Fallback for static assets (images, fonts, etc.)
	shadowRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})

	return shadowRouter
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
	handler = preloadLinkMiddleware(handler, embeddedFrontendFS)
	handler = cacheControlMiddleware(embeddedFrontendFS, handler)
	router.PathPrefix("/").Handler(handler)
	return nil
}
