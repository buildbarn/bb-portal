package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) saveBuildToolLogs(ctx context.Context, tx database.Handle, buildToolLogs *bes.BuildToolLogs) error {
	if buildToolLogs == nil {
		return nil
	}

	update := tx.Ent().BazelInvocation.UpdateOneID(r.InvocationDbID)

	for _, log := range buildToolLogs.GetLog() {
		file := files.ParseBepFile(log)
		if file == nil {
			continue
		}
		fileDbID, err := SaveSingleFile(ctx, tx, r.InstanceNameDbID, *file)
		if err != nil {
			return util.StatusWrap(err, "Failed to save build tool log to database")
		}
		update.AddBuildToolLogIDs(fileDbID)
	}
	err := update.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to add build tool logs to invocation")
	}
	return nil
}
