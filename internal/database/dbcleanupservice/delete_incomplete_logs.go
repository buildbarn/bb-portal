package dbcleanupservice

import (
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

// DeleteIncompleteLogs deletes logs which have had their incomplete
// build logs normalized.
func (dc *DbCleanupService) DeleteIncompleteLogs(ctx context.Context) error {
	ctx, span := dc.tracer.Start(ctx, "DbCleanupService.DeleteIncompleteLogs")
	defer span.End()

	start, count, err := dc.nextSlice(ctx, incompletebuildlog.Table)
	if err != nil {
		return err
	}

	deleted, err := dc.db.Sqlc().DeleteIncompleteLogsFromPages(ctx, sqlc.DeleteIncompleteLogsFromPagesParams{
		FromPage: start,
		Pages:    count,
	})
	if err != nil {
		return util.StatusWrap(err, "Could not delete incompleted build logs")
	}

	span.SetAttributes(attribute.Int64("deleted_invocations", deleted))

	return nil
}
