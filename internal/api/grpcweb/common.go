package grpcweb

import (
	"context"
	"log"

	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
)

// IsInstanceNamePrefixAllowed checks whether the given instance name prefix is
// allowed by the authorizer.
func IsInstanceNamePrefixAllowed(ctx context.Context, authorizer auth.Authorizer, instanceNamePrefix string) bool {
	instanceName, err := digest.NewInstanceName(instanceNamePrefix)
	if err != nil {
		log.Println("Error parsing instance name from operation: ", err)
		return false
	}
	return auth.AuthorizeSingleInstanceName(ctx, authorizer, instanceName) == nil
}
