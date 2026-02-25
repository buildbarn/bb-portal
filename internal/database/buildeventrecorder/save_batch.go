package buildeventrecorder

import (
	"context"
	"math/rand/v2"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (r *buildEventRecorder) saveBatch(ctx context.Context, batch []BuildEventWithInfo) (err error) {
	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveBatch",
		trace.WithAttributes(
			attribute.Int("batch_size", len(batch)),
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
		),
	)
	defer span.End()

	batch, rest := filterProgress(batch)
	if err = r.saveProgressBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch progress events")
	}

	batch, rest = filterConfigurationBatch(rest)
	if err = r.saveConfigurationBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch configuration events")
	}

	batch, rest = filterTargetConfiguredBatch(rest)
	if err = r.saveTargetConfiguredBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch target configured events")
	}

	batch, rest = filterTargetCompletedBatch(rest)
	if err = r.saveTargetCompletedBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch target completed events")
	}

	batch, rest = filterTestResultBatch(rest)
	if err = r.saveTestResultBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch test result events")
	}

	batch, rest = filterTestSummaryBatch(rest)
	if err = r.saveTestSummaryBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch test summary events")
	}

	if err = r.saveRemainingBatch(ctx, rest); err != nil {
		return util.StatusWrap(err, "Failed to save individual events")
	}
	return nil
}

// SaveBatch saves a batch of build events to the database.
func (r *buildEventRecorder) SaveBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.SaveBatch",
		trace.WithAttributes(
			attribute.Int("batch_size", len(batch)),
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
		),
	)
	defer span.End()
	retryCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer func() { cancel() }()

	batch = filterNilEvents(batch)
	backoff := 1 * time.Millisecond
	var errs []error
	for {
		var batchErr error
		if batchErr = r.loadHandledEvents(ctx); batchErr == nil {
			batch = r.filterHandledEvents(batch)
			if batchErr = r.saveBatch(ctx, batch); batchErr == nil {
				return nil
			}
		}
		// Check for forward progress. As long as we have forward
		// progress we reset the retry loop.
		if err := r.loadHandledEvents(ctx); err == nil {
			postBatch := r.filterHandledEvents(batch)
			if len(postBatch) < len(batch) {
				cancel()
				retryCtx, cancel = context.WithTimeout(ctx, 1*time.Second)
				backoff = 1 * time.Millisecond
				batch = postBatch
				errs = nil
				continue
			}
		}
		// No forward progress, perform sleep and retry.
		errs = append(errs, batchErr)
		jitter := time.Duration(float64(backoff) * (0.9 + 0.2*rand.Float64()))
		span.AddEvent("Retrying", trace.WithAttributes(
			attribute.String("last_error", batchErr.Error()),
			attribute.Int64("backoff", int64(jitter)),
		))
		select {
		case <-retryCtx.Done():
			return util.StatusFromMultiple(errs)
		case <-time.After(jitter):
		}
		backoff *= 2
	}
}

func filterNilEvents(batch []BuildEventWithInfo) []BuildEventWithInfo {
	ret := make([]BuildEventWithInfo, 0, len(batch))
	for _, x := range batch {
		if x.Event != nil {
			ret = append(ret, x)
		}
	}
	return ret
}

func filterProgress(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_Progress:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}

func filterConfigurationBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_Configuration:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}

func filterTargetCompletedBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TargetCompleted:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}

func filterTargetConfiguredBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TargetConfigured:
			// Skip target configured events for targets which did not
			// become configured.
			if x.Event.GetConfigured() != nil {
				filtered = append(filtered, x)
			}
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}

func filterTestResultBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TestResult:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}

func filterTestSummaryBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TestSummary:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest
}
