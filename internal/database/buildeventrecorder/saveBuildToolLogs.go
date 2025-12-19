package buildeventrecorder

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
)

func getFileDigest(file *bes.File) string {
	if file.GetDigest() != "" {
		return file.GetDigest()
	}
	if uri := file.GetUri(); uri != "" {
		stringArr := strings.Split(uri, "/")
		return stringArr[len(stringArr)-2]
	}
	return ""
}

func getFileSizeBytes(file *bes.File) int64 {
	if file.GetLength() != 0 {
		return file.GetLength()
	}
	if uri := file.GetUri(); uri != "" {
		stringArr := strings.Split(uri, "/")
		sizeBytes, err := strconv.ParseInt(stringArr[len(stringArr)-1], 10, 64)
		if err == nil {
			return sizeBytes
		}
	}
	return 0
}

func getFileDigestFunction(file *bes.File) string {
	if uri := file.GetUri(); uri != "" {
		stringArr := strings.Split(uri, "/")
		digestFunction := stringArr[len(stringArr)-3]
		for df := range remoteexecution.DigestFunction_Value_value {
			if strings.ToLower(df) == digestFunction {
				return digestFunction
			}
		}
	}
	return strings.ToLower(remoteexecution.DigestFunction_SHA256.String())
}

func (r *BuildEventRecorder) saveBuildToolLogs(ctx context.Context, tx *ent.Client, buildToolLogs *bes.BuildToolLogs) error {
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
		if digest := getFileDigest(log); digest != "" {
			file.SetDigest(digest)
		}
		if length := getFileSizeBytes(log); length != 0 {
			file.SetSizeBytes(length)
		}
		if digestFunction := getFileDigestFunction(log); digestFunction != "" {
			file.SetDigestFunction(digestFunction)
		}

		err := file.Exec(ctx)
		if err != nil {
			slog.Error("Failed to save build tool log", "name", fullName, "err", err)
		}
	}
	return nil
}
