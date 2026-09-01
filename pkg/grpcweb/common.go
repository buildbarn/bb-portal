package grpcweb

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
)

// grpcWebEndpointPrefix is the API prefix for all gRPC Web requests.
const grpcWebEndpointPrefix = "/api/v1/grpcweb"

// AddGrpcWebEndpoint adds a gRPC Web endpoint to a router, with the
// grpcWebEndpointPrefix added to the path.
func AddGrpcWebEndpoint(router *http.ServeMux, grpcWebServer *grpcweb.WrappedGrpcServer, path string) {
	router.Handle(
		grpcWebEndpointPrefix+path,
		http.StripPrefix(
			grpcWebEndpointPrefix,
			grpcWebServer,
		),
	)
}
