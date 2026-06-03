package zstdmiddleware

import (
	"io"
	"net/http"
	"strings"

	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-storage/pkg/zstd"
)

// zstdResponseWriter wraps the original ResponseWriter and redirects Writes through the zstd encoder.
type zstdResponseWriter struct {
	http.ResponseWriter
	writer      io.Writer
	wroteHeader bool
}

func (zw *zstdResponseWriter) WriteHeader(statusCode int) {
	zw.wroteHeader = true

	zw.ResponseWriter.Header().Set("Content-Encoding", "zstd")
	zw.ResponseWriter.Header().Add("Vary", "Accept-Encoding")
	zw.ResponseWriter.Header().Del("Content-Length")

	zw.ResponseWriter.WriteHeader(statusCode)
}

func (zw *zstdResponseWriter) Write(b []byte) (int, error) {
	// ZSTD requires modifying the headers of the response which must be
	// done before the body is written.
	if !zw.wroteHeader {
		zw.WriteHeader(http.StatusOK)
	}
	return zw.writer.Write(b)
}

// NewZstdMiddleware returns a standard gorilla/mux compatible middleware function
func NewZstdMiddleware(zstdPool zstd.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client supports zstd
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "zstd") {
				next.ServeHTTP(w, r)
				return
			}

			// Let Vite development websockets through unaffected
			if r.Header.Get("Upgrade") == "websocket" {
				next.ServeHTTP(w, r)
				return
			}

			encoder, err := zstdPool.NewEncoder(common.ExtractContextFromRequest(r), w)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			defer encoder.Close()

			next.ServeHTTP(
				&zstdResponseWriter{
					ResponseWriter: w,
					writer:         encoder,
				},
				r,
			)
		})
	}
}
