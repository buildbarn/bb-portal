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

func (r *BuildEventRecorder) saveTargetConfigured(ctx context.Context, tx *ent.Client, targetConfigured *bes.TargetConfigured, targetConfiguredID *bes.BuildEventId_TargetConfiguredId) error {
	if r.saveTargetDataLevel.GetNone() != nil {
		return nil
	}

	if targetConfigured == nil {
		return nil
	}
	if targetConfiguredID == nil {
		return fmt.Errorf("missing TargetCompletedId for TargetConfigured BES message")
	}

	if targetConfiguredID.Label == "" {
		return fmt.Errorf("missing label in TargetCompletedId for TargetConfigured BES message")
	}

	targetID, err := tx.Target.Create().
		SetInstanceNameID(r.InstanceNameDbID).
		SetLabel(targetConfiguredID.Label).
		SetAspect(targetConfiguredID.Aspect).
		SetTargetKind(targetConfigured.TargetKind).
		OnConflictColumns(
			target.FieldLabel,
			target.FieldAspect,
			target.FieldTargetKind,
			target.InstanceNameColumn,
		).
		Ignore().
		ID(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to create target for TargetConfigured BES message")
	}

	create := tx.TargetKindMapping.Create().
		SetBazelInvocationID(r.InvocationDbID).
		SetTargetID(targetID)

	if r.IsRealTime {
		create.SetStartTimeInMs(time.Now().UnixMilli())
	}

	err = create.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to create TargetKindMapping for target configured BES message")
	}
	return nil
}
