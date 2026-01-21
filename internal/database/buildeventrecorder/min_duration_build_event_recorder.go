package buildeventrecorder

import (
	"context"
	"math/rand/v2"
	"time"
)

type minDurationBuildEventRecorder struct {
	BuildEventRecorder
	minDuration time.Duration
}

// NewMinDurationBuildEventRecorder decorates a BuildEventRecorder so
// that it processes a batch for atleast as long as configured. This
// prevents small batches from consuming an undue amount of resources.
func NewMinDurationBuildEventRecorder(buildEventRecorder BuildEventRecorder, minDuration time.Duration) BuildEventRecorder {
	return &minDurationBuildEventRecorder{
		BuildEventRecorder: buildEventRecorder,
		minDuration:        minDuration,
	}
}

// SaveBatch implementation for minDurationBuildEventRecorder. The empty
// batch is simply ignored but otherwise it sleeps if the processing is
// done before the configured duration has passed. On errors it returns
// immediately.
func (r *minDurationBuildEventRecorder) SaveBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}
	start := time.Now()
	ret := r.BuildEventRecorder.SaveBatch(ctx, batch)
	if ret != nil {
		return ret
	}
	if elapsed := time.Since(start); elapsed < r.minDuration {
		// Sleep with jitter
		factor := 0.95 + 0.1*rand.Float64()
		sleepTime := r.minDuration - elapsed
		time.Sleep(time.Duration(float64(sleepTime) * factor))
	}
	return nil
}
