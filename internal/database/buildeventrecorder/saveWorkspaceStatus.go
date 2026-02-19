package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) saveWorkspaceStatus(ctx context.Context, tx *ent.Client, workspaceStatus *bes.WorkspaceStatus) error {
	if workspaceStatus == nil {
		return nil
	}

	update := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventWorkspaceStatus(false),
		).
		SetProcessedEventWorkspaceStatus(true)

	for _, item := range workspaceStatus.GetItem() {
		switch item.GetKey() {
		case "BUILD_HOST":
			update.SetHostname(item.GetValue())
		case "BUILD_USER":
			update.SetUserLdap(item.GetValue())
		}
	}

	err := update.Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "WorkspaceStatus event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to update bazel invocation with workspace status BES message")
	}
	return nil
}
