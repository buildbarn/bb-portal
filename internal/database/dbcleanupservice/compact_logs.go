package dbcleanupservice

import (
	"bytes"
	"context"
	"fmt"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/klauspost/compress/zstd"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (dc *DbCleanupService) normalizeInvocation(ctx context.Context, invocation *ent.BazelInvocation) (err error) {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.normalizeInvocation")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Could not normalize invocation")
		}
		span.End()
	}()

	span.SetAttributes(attribute.String("invocation.id", invocation.InvocationID.String()))

	normalizedLogs, err := bes.GetNormalizedIncompleteBuildLogs(ctx, dc.db.Ent(), invocation.InvocationID)
	if err != nil {
		return util.StatusWrap(err, "Could not normalize log")
	}

	tx, err := dc.db.Ent().Tx(ctx)
	if err != nil {
		return util.StatusWrap(err, "Could not start transaction")
	}

	defer tx.Rollback()

	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return util.StatusWrap(err, "Could not create zstd encoder")
	}
	defer encoder.Close()

	const chunkSize = 8 * 1024 * 1024 // 8MiB
	line := int64(0)
	buffer := make([]byte, 0, chunkSize/50)
	for index := 0; index*chunkSize < len(normalizedLogs); index++ {
		start := index * chunkSize
		end := min((index+1)*chunkSize, len(normalizedLogs))
		data := normalizedLogs[start:end]
		lineCount := bytes.Count(data, []byte{'\n'})
		isLastLineComplete := data[len(data)-1] == '\n'
		firstLineIndex := line
		var lastLineIndex int64
		if isLastLineComplete {
			lastLineIndex = firstLineIndex + int64(lineCount) - 1
		} else {
			lastLineIndex = firstLineIndex + int64(lineCount)
		}
		data = encoder.EncodeAll(data, buffer[:0])
		_, err = tx.BuildLogChunk.Create().
			SetBazelInvocation(invocation).
			SetData(data).
			SetChunkIndex(index).
			SetFirstLineIndex(firstLineIndex).
			SetLastLineIndex(lastLineIndex).
			Save(ctx)
		if err != nil {
			return util.StatusWrap(err, "Could not save log chunk")
		}

		if isLastLineComplete {
			line = lastLineIndex + 1
		} else {
			line = lastLineIndex
		}
	}

	if err := tx.Commit(); err != nil {
		return util.StatusWrap(err, "Could not commit transaction")
	}

	return nil
}

// CompactLogs compacts incomplete build logs by merging log entries for
// the same invocation.
func (dc *DbCleanupService) CompactLogs(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.CompactLogs")
	defer span.End()

	invocations, err := dc.db.Ent().BazelInvocation.Query().
		Where(
			bazelinvocation.BepCompleted(true),
			bazelinvocation.HasIncompleteBuildLogs(),
			bazelinvocation.Not(
				bazelinvocation.HasBuildLogChunks(),
			),
		).
		All(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to query invocations with incomplete build logs")
		return util.StatusWrap(err, "Failed to query invocations with incomplete build logs")
	}

	errs := make([]error, 0, len(invocations))
	for _, invocation := range invocations {
		err = dc.normalizeInvocation(ctx, invocation)
		if err != nil {
			errs = append(errs, err)
		}
	}

	span.SetAttributes(
		attribute.Int("invocations.total", len(invocations)),
		attribute.Int("invocations.succeeded", len(invocations)-len(errs)),
	)

	if len(errs) != 0 {
		err = util.StatusFromMultiple(errs)
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("%d/%d invocation logs could not be compacted", len(errs), len(invocations)))
		return err
	}

	return nil
}
