package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveFetch(ctx context.Context, tx *ent.Tx, fetch *bes.Fetch) error {
	if fetch == nil {
		return nil
	}

	err := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
		).
		AddNumFetches(1).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save fetch BES message to database")
	}
	return nil
}
