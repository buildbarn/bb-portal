package common

import (
	"context"
	"log"

	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
)

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
