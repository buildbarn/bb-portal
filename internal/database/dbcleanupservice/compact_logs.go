package dbcleanupservice

import (
	"bytes"
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (dc *DbCleanupService) normalizeInvocation(ctx context.Context, invocationDbID int64, invocationID uuid.UUID) (err error) {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.normalizeInvocation")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Could not normalize invocation")
		}
		span.End()
	}()

	span.SetAttributes(attribute.String("invocation.id", invocationID.String()))

	normalizedLogs, err := bes.GetNormalizedIncompleteBuildLogs(ctx, dc.db.Ent(), invocationID)
	if err != nil {
		return util.StatusWrap(err, "Could not normalize log")
	}

	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return util.StatusWrap(err, "Could not create zstd encoder")
	}
	defer encoder.Close()

	const chunkSize = 8 * 1024 * 1024 // 8MiB
	line := int64(0)
	buffer := make([]byte, 0, chunkSize/50)

	chunkBuilders := make([]*ent.BuildLogChunkCreate, 0)

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

		chunkBuilders = append(
			chunkBuilders,
			dc.db.Ent().BuildLogChunk.Create().
				SetBazelInvocationID(invocationDbID).
				SetData(data).
				SetChunkIndex(index).
				SetFirstLineIndex(firstLineIndex).
				SetLastLineIndex(lastLineIndex),
		)
		if isLastLineComplete {
			line = lastLineIndex + 1
		} else {
			line = lastLineIndex
		}
	}

	if err := dc.db.Ent().BuildLogChunk.CreateBulk(chunkBuilders...).Exec(ctx); err != nil {
		alreadyNormalized, checkErr := dc.db.Ent().BazelInvocation.Query().
			Where(
				bazelinvocation.IDEQ(invocationDbID),
				bazelinvocation.BepCompleted(true),
				bazelinvocation.HasIncompleteBuildLogs(),
				bazelinvocation.Not(bazelinvocation.HasBuildLogChunks()),
			).Exist(ctx)

		// The invocation was already normalized so we ignore the
		// error.
		if checkErr == nil && alreadyNormalized {
			return nil
		}

		return util.StatusWrap(err, "Could not save log chunk")
	}
	return nil
}

// CompactLogs compacts incomplete build logs by merging log entries for
// the same invocation.
func (dc *DbCleanupService) CompactLogs(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, "bazel_invocations")
	if err != nil {
		return 0, err
	}

	compacted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		invocations, err := dc.db.Sqlc().GetInvocationsForLogCompactionFromPages(ctx, sqlc.GetInvocationsForLogCompactionFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
		if err != nil {
			return 0, util.StatusWrap(err, "Failed to query invocations for log compaction")
		}

		if len(invocations) == 0 {
			return 0, nil
		}

		errs := make([]error, 0, len(invocations))
		for _, inv := range invocations {
			err = dc.normalizeInvocation(ctx, inv.ID, inv.InvocationID)
			if err != nil {
				errs = append(errs, err)
			}
		}

		succeeded := int64(len(invocations) - len(errs))
		if len(errs) != 0 {
			return succeeded, util.StatusFromMultiple(errs)
		}

		return int64(len(invocations)), nil
	})

	return compacted, err
}
