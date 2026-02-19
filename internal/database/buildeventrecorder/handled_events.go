package buildeventrecorder

import (
	"context"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) loadHandledEvents(ctx context.Context) error {
	r.handledEvents.bitmap = roaring.NewBitmap()
	eventMetadata, err := r.db.Sqlc().GetOrCreateEventMetadata(ctx, int64(r.InvocationDbID))
	if err != nil {
		return util.StatusWrap(err, "Could not load handled events")
	}
	r.handledEvents.version = eventMetadata.Version
	r.handledEvents.id = eventMetadata.ID
	if len(eventMetadata.Handled) == 0 {
		return nil
	}
	if _, err = r.handledEvents.bitmap.FromBuffer(eventMetadata.Handled); err != nil {
		return util.StatusWrap(err, "Could not deserialize handled events")
	}
	return nil
}

func (r *buildEventRecorder) saveHandledEvents(ctx context.Context, db database.Handle, timestamp time.Time) error {
	r.handledEvents.bitmap.RunOptimize()
	bytes, err := r.handledEvents.bitmap.ToBytes()
	if err != nil {
		return util.StatusWrap(err, "Could not serialize handled events")
	}
	data := sqlc.UpdateEventMetadataParams{
		Handled:         bytes,
		Version:         r.handledEvents.version,
		ID:              r.handledEvents.id,
		EventReceivedAt: timestamp,
	}
	version, err := db.Sqlc().UpdateEventMetadata(ctx, data)
	if err != nil {
		return util.StatusWrap(err, "Could not update event metadata")
	}
	r.handledEvents.version = version
	return nil
}

func (r *buildEventRecorder) filterHandledEvents(batch []BuildEventWithInfo) []BuildEventWithInfo {
	unhandled := make([]BuildEventWithInfo, 0, len(batch))
	for _, event := range batch {
		if !r.handledEvents.bitmap.Contains(event.SequenceNumber) {
			unhandled = append(unhandled, event)
		}
	}
	return unhandled
}

func (r *buildEventRecorder) saveHandledEventsForBatch(ctx context.Context, batch []BuildEventWithInfo, handle database.Handle) error {
	sequences := make([]uint32, len(batch))
	for i, x := range batch {
		sequences[i] = x.SequenceNumber
	}
	r.handledEvents.bitmap.AddMany(sequences)
	return r.saveHandledEvents(ctx, handle, time.Now())
}
