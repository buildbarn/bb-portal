package dbcleanupservice

import (
	"context"
	"time"

	"github.com/buildbarn/bb-storage/pkg/util"
)

// CalculateNextSlice is exported for testing.
func CalculateNextSlice(counter, totalPages, runsPerHour int64) (from, count int64) {
	count = max(totalPages/runsPerHour, 1)
	from = (counter * count) % totalPages
	prevFrom := ((counter - 1) * count) % totalPages

	if from < prevFrom {
		return 0, count + from
	}
	return from, count
}

// nextSlice calculates the physical page range for the current cleanup
// tick. It automatically adjusts the start page and number of pages
// with the goal of iterating through the entire table in one hour.
// Should the cleanupInteral be greater than one hour it will clean up
// the entire table in one tick. It will always do a minimum of one page
// worth of cleanup.
func (dc *DbCleanupService) nextSlice(ctx context.Context, table string) (from, count int64, err error) {
	pages32, err := dc.db.Sqlc().SelectPages(ctx, table)
	if err != nil {
		return 0, 0, util.StatusWrapf(err, "Failed to get pages for %s", table)
	}
	totalPages := max(int64(pages32), 1)
	runsPerHour := max(int64(time.Hour/dc.cleanupInterval), 1)
	from, count = CalculateNextSlice(dc.counter, totalPages, runsPerHour)
	return from, count, nil
}
