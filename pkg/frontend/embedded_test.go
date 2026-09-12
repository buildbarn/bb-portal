package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestPreloadLinkMiddleware(t *testing.T) {
	frontendFS := fstest.MapFS{
		"route-manifest.json": &fstest.MapFile{
			Data: []byte(`[
			{
				"path": "/bazel-invocations/{invocationID}",
				"js": [
					"bazel-invocations.js"
				],
				"css": [
					"bazel-invocations.css"
				]
			}
		]`),
		},
	}

	t.Run("Add preload links", func(t *testing.T) {
		nextCalled := false

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
		})
		handler := preloadLinkMiddleware(next, frontendFS)

		req := httptest.NewRequest(
			http.MethodGet,
			"/bazel-invocations/12345",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		links := rec.Header().Values("Link")

		require.Len(t, links, 2)
		require.True(t, nextCalled)
		require.Contains(t, links, "</bazel-invocations.css>; rel=preload; as=style; crossorigin")
		require.Contains(t, links, "</bazel-invocations.js>; rel=modulepreload; crossorigin")
	})

	t.Run("Static Assets fallback", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		})
		handler := preloadLinkMiddleware(next, frontendFS)

		req := httptest.NewRequest(
			http.MethodGet,
			"/assets/logo.png",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.True(t, nextCalled)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Empty(t, rec.Header().Values("Link"))
	})
}
