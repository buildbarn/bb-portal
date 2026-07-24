package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func getProfileFile(buildToolLogs *bes.BuildToolLogs) *files.ParsedBepFile {
	for _, log := range buildToolLogs.GetLog() {
		if log.Name == "command.profile.gz" || log.Name == "command.profile.json" {
			return files.ParseBepFile(log)
		}
	}
	return nil
}

func (r *buildEventRecorder) saveBuildToolLogs(ctx context.Context, tx database.Handle, buildToolLogs *bes.BuildToolLogs) error {
	if buildToolLogs == nil {
		return nil
	}
	profileFile := getProfileFile(buildToolLogs)
	if profileFile == nil {
		return nil
	}
	fileDbID, err := SaveSingleFile(ctx, tx, r.InstanceNameDbID, *profileFile)
	if err != nil {
		return util.StatusWrap(err, "Failed to save profile file to database")
	}
	err = tx.Ent().BazelInvocation.UpdateOneID(r.InvocationDbID).SetProfileID(fileDbID).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to add profile file to invocation")
	}
	return nil
}
