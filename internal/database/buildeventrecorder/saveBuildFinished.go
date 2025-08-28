package buildeventrecorder

import (
	"context"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveBuildFinished(ctx context.Context, tx *ent.Tx, finished *bes.BuildFinished) error {
	if finished == nil {
		return nil
	}

	var endedAt time.Time
	if finished.GetFinishTime() != nil {
		endedAt = finished.GetFinishTime().AsTime()
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		endedAt = time.UnixMilli(finished.GetFinishTimeMillis())
	}

	err := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventBuildFinished(false),
		).
		SetProcessedEventBuildFinished(true).
		SetEndedAt(endedAt).
		SetExitCodeCode(finished.GetExitCode().GetCode()).
		SetExitCodeName(finished.GetExitCode().GetName()).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "BuildFinished event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to update bazel invocation with build finished BES message")
	}
	return nil
}
