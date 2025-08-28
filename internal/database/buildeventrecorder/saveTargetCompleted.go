package buildeventrecorder

import (
	"context"
	"fmt"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/target"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveTargetCompleted(ctx context.Context, tx *ent.Tx, targetCompleted *bes.TargetComplete, label string, aborted *bes.Aborted) error {
	if r.saveTargetDataLevel.GetNone() != nil {
		return nil
	}

	if label == "" {
		return fmt.Errorf("missing label for TargetCompleted BES message")
	}

	targetDb, err := tx.Target.Query().Where(
		target.LabelEQ(label),
		target.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
	).Only(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to find target for target completed BEP message")
	}

	update := tx.Target.Update().
		Where(
			target.ID(targetDb.ID),
		)

	if r.IsRealTime {
		now := time.Now().UnixMilli()
		update.SetEndTimeInMs(now)
		update.SetDurationInMs(now - targetDb.StartTimeInMs)
	}
	if targetCompleted != nil {
		update.SetSuccess(targetCompleted.GetSuccess())
		if r.saveTargetDataLevel.GetEnriched() != nil && targetCompleted.TestTimeout != nil {
			update.SetTestTimeout(targetCompleted.TestTimeout.Seconds)
		}
	} else {
		update.SetSuccess(false)
	}
	if aborted != nil {
		update.SetAbortReason(target.AbortReason(bes.Aborted_AbortReason_name[int32(aborted.Reason)]))
	}

	err = update.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to update target for target completed BEP message")
	}
	return nil
}
