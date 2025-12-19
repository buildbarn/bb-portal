package buildeventrecorder

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/invocationtarget"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/target"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetkindmapping"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveTargetCompleted(ctx context.Context, tx *ent.Client, targetCompleted *bes.TargetComplete, targetCompletedID *bes.BuildEventId_TargetCompletedId, aborted *bes.Aborted) error {
	if r.saveTargetDataLevel.GetNone() != nil {
		return nil
	}

	if targetCompletedID == nil {
		return fmt.Errorf("missing TargetCompletedId for TargetCompleted BES message")
	}
	if targetCompletedID.Label == "" {
		return fmt.Errorf("missing label in TargetCompletedId for TargetCompleted BES message")
	}

	targetKindMapping, err := tx.TargetKindMapping.Query().
		Where(
			// This creates a more efficient query than using `HasBazelInvocationWith`
			predicate.TargetKindMapping(sql.FieldEQ(targetkindmapping.BazelInvocationColumn, r.InvocationDbID)),
			targetkindmapping.HasTargetWith(
				target.LabelEQ(targetCompletedID.Label),
				target.AspectEQ(targetCompletedID.Aspect),
			),
		).
		WithTarget().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find target kind mapping for label %q: %w", targetCompletedID.Label, err)
	}

	update := tx.InvocationTarget.Create().
		SetBazelInvocationID(r.InvocationDbID).
		SetTargetID(targetKindMapping.Edges.Target.ID)

	if r.IsRealTime {
		now := time.Now().UnixMilli()
		update.SetStartTimeInMs(targetKindMapping.StartTimeInMs)
		update.SetEndTimeInMs(now)
		update.SetDurationInMs(now - targetKindMapping.StartTimeInMs)
	}

	if targetCompleted != nil {
		update.SetSuccess(targetCompleted.Success)
		if targetCompleted.FailureDetail.GetMessage() != "" {
			update.SetFailureMessage(targetCompleted.FailureDetail.GetMessage())
		}
		if len(targetCompleted.Tag) > 0 {
			update.SetTags(targetCompleted.Tag)
		}
	}

	if aborted != nil {
		update.SetAbortReason(invocationtarget.AbortReason(bes.Aborted_AbortReason_name[int32(aborted.Reason)]))
	} else {
		update.SetAbortReason(invocationtarget.AbortReasonNONE)
	}

	err = update.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to update invocation target for target completed BEP message")
	}
	return nil
}
