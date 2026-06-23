package dbcleanupservice

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/klauspost/compress/zstd"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// compactArtifactGraph folds an invocation's staged artifact-graph events
// into a single zstd-compressed InvocationArtifactGraph blob. The blob is
// a stream of length-prefixed serialized BuildEvents — the format
// internal/graphql/artifact_graph.go decodes.
func (dc *DbCleanupService) compactArtifactGraph(ctx context.Context, invocation *ent.BazelInvocation) (err error) {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.compactArtifactGraph")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Could not compact artifact graph")
		}
		span.End()
	}()
	span.SetAttributes(attribute.String("invocation.id", invocation.InvocationID.String()))

	events, err := dc.db.Sqlc().GetIncompleteArtifactGraphEvents(ctx, invocation.ID)
	if err != nil {
		return util.StatusWrap(err, "Could not read staged artifact graph events")
	}
	if len(events) == 0 {
		return nil
	}

	var raw bytes.Buffer
	var lenBuf [binary.MaxVarintLen64]byte
	for _, event := range events {
		n := binary.PutUvarint(lenBuf[:], uint64(len(event)))
		raw.Write(lenBuf[:n])
		raw.Write(event)
	}

	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return util.StatusWrap(err, "Could not create zstd encoder")
	}
	defer encoder.Close()
	payload := encoder.EncodeAll(raw.Bytes(), nil)

	if err := dc.db.Sqlc().InsertInvocationArtifactGraph(ctx, sqlc.InsertInvocationArtifactGraphParams{
		Payload:           payload,
		BazelInvocationID: invocation.ID,
	}); err != nil {
		return util.StatusWrap(err, "Could not write compacted artifact graph")
	}
	return nil
}

// CompactArtifactGraphs folds staged artifact-graph events into a single
// compressed blob per completed invocation, the same way CompactLogs
// turns incomplete build logs into BuildLogChunks.
func (dc *DbCleanupService) CompactArtifactGraphs(ctx context.Context) (int64, error) {
	invocations, err := dc.db.Ent().BazelInvocation.Query().
		Where(
			bazelinvocation.BepCompleted(true),
			bazelinvocation.HasIncompleteArtifactGraphs(),
			bazelinvocation.Not(
				bazelinvocation.HasArtifactGraph(),
			),
		).
		All(ctx)
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to query invocations with staged artifact graphs")
	}

	errs := make([]error, 0, len(invocations))
	for _, invocation := range invocations {
		if err := dc.compactArtifactGraph(ctx, invocation); err != nil {
			errs = append(errs, err)
		}
	}

	succeeded := int64(len(invocations) - len(errs))
	if len(errs) != 0 {
		return succeeded, util.StatusFromMultiple(errs)
	}
	return succeeded, nil
}
