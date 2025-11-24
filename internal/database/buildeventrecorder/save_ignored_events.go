package buildeventrecorder

import (
	"context"

	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (r *BuildEventRecorder) saveIgnoredEventsBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}

	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveIgnoredEventsBatch",
		trace.WithAttributes(
			attribute.Int("batch_size", len(batch)),
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
		),
	)
	defer span.End()

	if err := r.createEventMetadatasForBatch(ctx, batch, r.db); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	return nil
}
