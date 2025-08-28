package buildeventrecorder

import (
	"context"
	"fmt"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/target"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveTargetConfigured(ctx context.Context, tx *ent.Tx, targetConfigured *bes.TargetConfigured, label string) error {
	if r.saveTargetDataLevel.GetNone() != nil {
		return nil
	}

	if label == "" {
		return fmt.Errorf("missing label for TargetConfigured BES message")
	}
	if targetConfigured == nil {
		return nil
	}

	create := tx.Target.Create().
		SetLabel(label).
		SetTargetKind(targetConfigured.TargetKind).
		SetTestSize(target.TestSize(bes.TestSize_name[int32(targetConfigured.TestSize)])).
		SetSuccess(false).
		SetBazelInvocationID(r.InvocationDbID)

	if r.saveTargetDataLevel.GetEnriched() != nil {
		create.SetTag(targetConfigured.Tag)
	}

	if r.IsRealTime {
		create.SetStartTimeInMs(time.Now().UnixMilli())
	}

	err := create.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to create target for target configured BES message")
	}
	return nil
}
