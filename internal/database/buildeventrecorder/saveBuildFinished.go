package buildeventrecorder

import (
	"context"
	"log/slog"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) saveBuildFinished(ctx context.Context, tx database.Tx, finished *bes.BuildFinished) error {
	if finished == nil {
		return nil
	}

	var endedAt time.Time
	if finished.GetFinishTime() != nil {
		endedAt = finished.GetFinishTime().AsTime()
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		endedAt = time.UnixMilli(finished.GetFinishTimeMillis())
	}

	err := tx.Ent().BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventBuildFinished(false),
		).
		SetProcessedEventBuildFinished(true).
		SetEndedAt(endedAt).
		SetExitCodeCode(finished.GetExitCode().GetCode()).
		SetExitCodeName(finished.GetExitCode().GetName()).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "BuildFinished event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to update bazel invocation with build finished BES message")
	}

	if err := r.flushArtifactGraph(ctx, tx); err != nil {
		return util.StatusWrap(err, "Failed to flush artifact graph at completion")
	}
	return nil
}

// flushArtifactGraph finalizes the in-memory streaming buffer and writes
// a single invocation_artifact_graphs row. No-op when the buffer is nil
// (save level below basic_and_target_and_artifacts) or empty (no
// NamedSetOfFiles or TargetCompleted events were seen).
func (r *buildEventRecorder) flushArtifactGraph(ctx context.Context, tx database.Tx) error {
	if r.artifactGraph == nil {
		return nil
	}
	buf := r.artifactGraph
	r.artifactGraph = nil
	payload, uncompressed, err := buf.Finalize()
	if err != nil {
		return util.StatusWrap(err, "Failed to finalize artifact graph buffer")
	}
	if uncompressed == 0 {
		return nil
	}
	if buf.Capped() {
		slog.Warn("Artifact graph buffer hit uncompressed size cap; partial graph stored",
			"invocation_id", r.InvocationID,
			"uncompressed_bytes", uncompressed)
	}
	return tx.Sqlc().InsertInvocationArtifactGraph(ctx, sqlc.InsertInvocationArtifactGraphParams{
		Payload:           payload,
		UncompressedSize:  uncompressed,
		BazelInvocationID: r.InvocationDbID,
	})
}
