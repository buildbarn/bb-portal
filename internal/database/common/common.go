package common

import (
	"fmt"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
)

// RollbackAndWrapError attempts to roll back the provided transaction.
// If that fails, the rollback error is combined with the original error
// into a single error value.
func RollbackAndWrapError(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		if err == nil {
			return util.StatusWrap(rerr, "Failed to rollback transaction")
		}
		return util.StatusFromMultiple([]error{
			util.StatusWrap(rerr, "Failed to rollback transaction"),
			err,
		})
	}
	return err
}

// CalculateBuildUUID calculates a UUID for a build, based on the build URL
// and instance name.
func CalculateBuildUUID(buildURL, instanceName string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("instanceName: %s, buildUrl: %s", instanceName, buildURL)))
}
