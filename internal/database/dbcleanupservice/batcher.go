package dbcleanupservice

import (
	"context"
	"time"

	"github.com/buildbarn/bb-storage/pkg/clock"
)

// Batcher is a utility that batches executions of an operation into
// chunks.
type Batcher interface {
	Batch(context.Context, func(context.Context, int64) (int64, error)) (int64, error)
}

type timedBatcher struct {
	clock          clock.Clock
	targetDuration time.Duration
	minBatchSize   int64
	maxBatchSize   int64
}

// NewTimedBatcher returns a batcher that uses exponential growth to
// reach a target duration.
func NewTimedBatcher(clock clock.Clock, targetDuration time.Duration, minBatchSize, maxBatchSize int64) Batcher {
	return &timedBatcher{
		clock:          clock,
		targetDuration: targetDuration,
		minBatchSize:   minBatchSize,
		maxBatchSize:   maxBatchSize,
	}
}

func (b *timedBatcher) Batch(ctx context.Context, operation func(context.Context, int64) (int64, error)) (int64, error) {
	var totalProcessed int64

	batchLimit := b.minBatchSize

	for {
		if err := ctx.Err(); err != nil {
			return totalProcessed, err
		}

		start := b.clock.Now()
		processed, err := operation(ctx, batchLimit)
		totalProcessed += processed
		if err != nil {
			return totalProcessed, err
		}
		// Fewer rows were processed than the batch limit, we are done.
		if processed < batchLimit {
			return totalProcessed, nil
		}

		duration := b.clock.Now().Sub(start)
		if duration < b.targetDuration {
			// Query was fast, double limit and try again.
			batchLimit = min(2*batchLimit, b.maxBatchSize)
		}
	}
}
