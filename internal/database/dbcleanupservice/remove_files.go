package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// RemoveUnusedFilePaths removes file paths that are no longer referenced by
// any files.
func (dc *DbCleanupService) RemoveUnusedFilePaths(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, bazelinvocation.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteUnusedFilePathsFromPages(ctx, sqlc.DeleteUnusedFilePathsFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old invocations")
	}
	return deleted, nil
}

// RemoveUnusedDigests removes digests that are no longer referenced by
// any files.
func (dc *DbCleanupService) RemoveUnusedDigests(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, bazelinvocation.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteUnusedDigestFromPages(ctx, sqlc.DeleteUnusedDigestFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old invocations")
	}
	return deleted, nil
}

// RemoveUnusedFiles removes files that are no longer used
func (dc *DbCleanupService) RemoveUnusedFiles(ctx context.Context) (int64, error) {
	start, count, err := dc.nextSlice(ctx, bazelinvocation.Table)
	if err != nil {
		return 0, err
	}

	deleted, err := dc.batcher.Batch(ctx, func(ctx context.Context, limit int64) (int64, error) {
		return dc.db.Sqlc().DeleteUnusedFilesFromPages(ctx, sqlc.DeleteUnusedFilesFromPagesParams{
			FromPage:   start,
			Pages:      count,
			BatchLimit: limit,
		})
	})
	if err != nil {
		return 0, util.StatusWrap(err, "Failed to remove old invocations")
	}
	return deleted, nil
}
