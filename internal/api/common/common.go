package common

import (
	"context"
	"log"
	"net/http"

	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"google.golang.org/grpc/metadata"
)

// ExtractContextFromRequest extracts a Context from an incoming HTTP request,
// forwarding any request headers as gRPC metadata.
func ExtractContextFromRequest(req *http.Request) context.Context {
	var pairs []string
	for key, values := range req.Header {
		for _, value := range values {
			pairs = append(pairs, key, value)
		}
	}
	return metadata.NewIncomingContext(req.Context(), metadata.Pairs(pairs...))
}

// IsInstanceNameAllowed checks whether the given instance name is allowed by
// the authorizer.
func IsInstanceNameAllowed(ctx context.Context, authorizer auth.Authorizer, instanceNameString string) bool {
	instanceName, err := digest.NewInstanceName(instanceNameString)
	if err != nil {
		log.Println("Error parsing instance name from operation: ", err)
		return false
	}
	return auth.AuthorizeSingleInstanceName(ctx, authorizer, instanceName) == nil
}
