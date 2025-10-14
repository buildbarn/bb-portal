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

func (r *BuildEventRecorder) saveBatch(ctx context.Context, batch []BuildEventWithInfo) error {
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

	batch, rest, err := filterProgress(batch)
	if err != nil {
		return util.StatusWrap(err, "Failed to filter progress events from batch")
	}
	if err = r.saveProgressBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch progress events")
	}

	batch, rest, err = filterConfigurationBatch(rest)
	if err != nil {
		return util.StatusWrap(err, "Failed to filter configuration events from batch")
	}
	if err = r.saveConfigurationBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch configuration events")
	}

	batch, rest, err = filterTargetConfiguredBatch(rest)
	if err != nil {
		return util.StatusWrap(err, "Failed to filter target configured events from batch")
	}
	if err = r.saveTargetConfiguredBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch target configured events")
	}

	batch, rest, err = filterTargetCompletedBatch(rest)
	if err != nil {
		return util.StatusWrap(err, "Failed to filter target completed events from batch")
	}
	if err = r.saveTargetCompletedBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch target completed events")
	}

	batch, rest, err = filterTestResultBatch(rest)
	if err != nil {
		return util.StatusWrap(err, "Failed to filter test result events from batch")
	}
	if err = r.saveTestResultBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch test result events")
	}

	batch, rest, err = filterTestSummaryBatch(rest)
	if err != nil {
		return util.StatusWrap(err, "Failed to filter test summary events from batch")
	}
	if err = r.saveTestSummaryBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save batch test summary events")
	}

	// These are events which are not handled by the individual event
	// handler, they are filtered out to efficiently save their event
	// metadata in a batch.
	batch, rest, err = filterIgnoredIndividualEvents(rest)
	if err != nil {
		return util.StatusWrap(err, "Failed to ignored events from batch")
	}
	if err = r.saveIgnoredEventsBatch(ctx, batch); err != nil {
		return util.StatusWrap(err, "Failed to save ignored events")
	}

	// The remaining events tend to not arrive in large batches and may
	// as well be handled one by one.
	for _, info := range rest {
		if err := r.saveEvent(ctx, info); err != nil {
			return util.StatusWrap(err, "Failed to save individual event")
		}
	}

	return nil
}

// SaveBatch saves a batch of build events to the database.
func (r *BuildEventRecorder) SaveBatch(ctx context.Context, batch []BuildEventWithInfo) error {
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

func filterProgress(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_Progress:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest, nil
}

func filterConfigurationBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_Configuration:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest, nil
}

func filterNamedSet(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_NamedSet:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest, nil
}

func filterTargetCompletedBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TargetCompleted:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest, nil
}

func filterTargetConfiguredBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
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
	return filtered, rest, nil
}

func filterTestResultBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TestResult:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest, nil
}

func filterTestSummaryBatch(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_TestSummary:
			filtered = append(filtered, x)
		default:
			rest = append(rest, x)
		}
	}
	return filtered, rest, nil
}

func filterIgnoredIndividualEvents(batch []BuildEventWithInfo) (filtered, rest []BuildEventWithInfo, err error) {
	for _, x := range batch {
		switch x.Event.GetId().GetId().(type) {
		case *bes.BuildEventId_Started,
			*bes.BuildEventId_ActionCompleted,
			*bes.BuildEventId_BuildMetadata,
			*bes.BuildEventId_OptionsParsed,
			*bes.BuildEventId_BuildFinished,
			*bes.BuildEventId_BuildMetrics,
			*bes.BuildEventId_StructuredCommandLine,
			*bes.BuildEventId_Fetch,
			*bes.BuildEventId_TestResult,
			*bes.BuildEventId_TestSummary,
			*bes.BuildEventId_BuildToolLogs,
			*bes.BuildEventId_WorkspaceStatus:
			// The above event types are handled.
			rest = append(rest, x)
		default:
			// All other event types are ignored.
			filtered = append(filtered, x)
		}
	}
	return filtered, rest, nil
}
