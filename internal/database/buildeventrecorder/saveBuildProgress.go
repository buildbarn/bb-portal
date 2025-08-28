package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveBuildProgress(ctx context.Context, tx *ent.Tx, buildEvent *events.BuildEvent, progress *bes.Progress) error {
	if buildEvent == nil || progress == nil {
		return nil
	}

	opaqueCount := buildEvent.GetId().GetProgress().GetOpaqueCount()
	log := progress.GetStderr()
	if log != progress.GetStdout() {
		log += progress.GetStdout()
	}

	err := tx.IncompleteBuildLog.Create().
		SetBazelInvocationID(r.InvocationDbID).
		SetSnippetID(opaqueCount).
		SetLogSnippet(log).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save incomplete build log snippet")
	}
	return nil
}
