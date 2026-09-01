package prometheusservice

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewPrometheusService registers a reverse-proxy handler at
// /api/v1/prometheus/* that forwards requests to the configured Prometheus URL.
// This lets the frontend query Prometheus without direct cluster access or CORS issues.
func NewPrometheusService(prometheusURL string, router *mux.Router) error {
	target, err := url.Parse(prometheusURL)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid prometheus_url %q: %v", prometheusURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	// Remove Accept-Encoding so Prometheus returns plain (uncompressed) JSON.
	// Without this, the browser's Accept-Encoding header is forwarded verbatim,
	// Prometheus compresses the response, and the Go proxy passes through the
	// compressed bytes without decompressing — breaking JSON parsing in the browser.
	proxy.ModifyResponse = nil
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Del("Accept-Encoding")
	}
	// Strip the /api/v1/prometheus prefix before forwarding so the upstream
	// Prometheus receives its native /api/v1/... paths.
	router.PathPrefix("/api/v1/prometheus/").Handler(
		http.StripPrefix("/api/v1/prometheus", proxy),
	)
	return nil
}
