package dbcleanupservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTimedBatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClock := mock.NewMockClock(ctrl)
	ctx := context.Background()

	targetDuration := 1 * time.Second
	minBatchSize := int64(100)
	maxBatchSize := int64(1000)

	t.Run("ImmediateCompletion", func(t *testing.T) {
		batcher := dbcleanupservice.NewTimedBatcher(mockClock, targetDuration, minBatchSize, maxBatchSize)

		mockClock.EXPECT().Now().Return(time.Unix(0, 0))

		operation := func(ctx context.Context, limit int64) (int64, error) {
			require.Equal(t, int64(100), limit)
			return 50, nil // Processed less than limit, should terminate immediately
		}

		total, err := batcher.Batch(ctx, operation)
		require.NoError(t, err)
		require.Equal(t, int64(50), total)
	})

	t.Run("ExponentialGrowthAndClamping", func(t *testing.T) {
		batcher := dbcleanupservice.NewTimedBatcher(mockClock, targetDuration, minBatchSize, maxBatchSize)
		startTime := time.Unix(0, 0)

		// Batch 1: Limit 100
		mockClock.EXPECT().Now().Return(startTime)
		mockClock.EXPECT().Now().Return(startTime.Add(100 * time.Millisecond)) // fast

		// Batch 2: Limit 200
		mockClock.EXPECT().Now().Return(startTime)
		mockClock.EXPECT().Now().Return(startTime.Add(100 * time.Millisecond)) // fast

		// Batch 3: Limit 400
		mockClock.EXPECT().Now().Return(startTime)
		mockClock.EXPECT().Now().Return(startTime.Add(100 * time.Millisecond)) // fast

		// Batch 4: Limit 800
		mockClock.EXPECT().Now().Return(startTime)
		mockClock.EXPECT().Now().Return(startTime.Add(100 * time.Millisecond)) // fast

		// Batch 5: Limit 1000 (Clamped from 1600)
		mockClock.EXPECT().Now().Return(startTime)

		expectedLimits := []int64{100, 200, 400, 800, 1000}
		callCount := 0

		operation := func(ctx context.Context, limit int64) (int64, error) {
			require.Equal(t, expectedLimits[callCount], limit)
			callCount++

			if limit == 1000 {
				return 42, nil // Finish on the last capped batch
			}
			return limit, nil // Return full limit to trigger next loop
		}

		total, err := batcher.Batch(ctx, operation)
		require.NoError(t, err)

		// 100 + 200 + 400 + 800 + 42 = 1542
		require.Equal(t, int64(1542), total)
	})

	t.Run("SlowQueryHaltsGrowth", func(t *testing.T) {
		batcher := dbcleanupservice.NewTimedBatcher(mockClock, targetDuration, minBatchSize, maxBatchSize)
		startTime := time.Unix(0, 0)

		// Batch 1: Limit 100, takes 2 seconds (slower than target)
		mockClock.EXPECT().Now().Return(startTime)
		mockClock.EXPECT().Now().Return(startTime.Add(2 * time.Second))

		// Batch 2: Limit should remain 100
		mockClock.EXPECT().Now().Return(startTime)

		callCount := 0
		operation := func(ctx context.Context, limit int64) (int64, error) {
			require.Equal(t, int64(100), limit)
			callCount++
			if callCount == 2 {
				return 0, nil
			}
			return limit, nil
		}

		total, err := batcher.Batch(ctx, operation)
		require.NoError(t, err)
		require.Equal(t, int64(100), total)
	})

	t.Run("PropagatesOperationError", func(t *testing.T) {
		batcher := dbcleanupservice.NewTimedBatcher(mockClock, targetDuration, minBatchSize, maxBatchSize)
		expectedErr := errors.New("database connection lost")

		mockClock.EXPECT().Now().Return(time.Unix(0, 0))

		operation := func(ctx context.Context, limit int64) (int64, error) {
			return 50, expectedErr
		}

		total, err := batcher.Batch(ctx, operation)
		require.ErrorIs(t, err, expectedErr)
		// Should still return the rows processed before the error
		require.Equal(t, int64(50), total)
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		batcher := dbcleanupservice.NewTimedBatcher(mockClock, targetDuration, minBatchSize, maxBatchSize)

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel before we even start

		operation := func(ctx context.Context, limit int64) (int64, error) {
			t.Fatal("Operation should not be called if context is cancelled")
			return 0, nil
		}

		total, err := batcher.Batch(cancelCtx, operation)
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, int64(0), total)
	})
}
