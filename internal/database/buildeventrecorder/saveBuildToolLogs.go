package buildeventrecorder

import (
	"context"
	"log/slog"
	"strings"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
)

func (r *buildEventRecorder) saveBuildToolLogs(ctx context.Context, tx *ent.Client, buildToolLogs *bes.BuildToolLogs) error {
	if buildToolLogs == nil {
		return nil
	}

	for _, log := range buildToolLogs.GetLog() {
		name := log.GetName()
		if name == "" {
			slog.Warn("Skipping build tool log with empty name")
			continue
		}

		pathPrefix := log.GetPathPrefix()
		pathPrefix = append(pathPrefix, name)
		fullName := strings.Join(pathPrefix, "/")

		file := tx.InvocationFiles.Create().
			SetName(fullName).
			SetBazelInvocationID(r.InvocationDbID)

		if content := string(log.GetContents()); content != "" {
			file.SetContent(content)
		}
		if digest := getFileDigestFromBesFile(log); digest != nil {
			file.SetDigest(*digest)
		}
		if length := getFileSizeBytesFromBesFile(log); length != nil {
			file.SetSizeBytes(*length)
		}
		if digestFunction := getFileDigestFunctionFromBesFile(log); digestFunction != nil {
			file.SetDigestFunction(*digestFunction)
		}

		err := file.Exec(ctx)
		if err != nil {
			slog.Error("Failed to save build tool log", "name", fullName, "err", err)
		}
	}
	return nil
}
