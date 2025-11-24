package buildeventrecorder

import (
	"context"
	"time"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
)

func (r *BuildEventRecorder) createEventMetadatasForBatch(ctx context.Context, batch []BuildEventWithInfo, handle database.Handle) error {
	params := sqlc.CreateEventMetadataBulkParams{
		BazelInvocationID: int64(r.InvocationDbID),
		EventHashes:       make([][]byte, len(batch)),
		EventReceivedAts:  make([]time.Time, len(batch)),
		SequenceNumbers:   make([]int64, len(batch)),
	}
	for i, x := range batch {
		params.EventHashes[i] = x.EventHash
		params.EventReceivedAts[i] = x.AddedAt
		params.SequenceNumbers[i] = x.SequenceNumber
	}
	return handle.Sqlc().CreateEventMetadataBulk(ctx, params)
}
