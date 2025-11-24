package buildeventrecorder

import (
	"context"
	"slices"

	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventmetadata"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *BuildEventRecorder) filterHandledEvents(ctx context.Context, batch []BuildEventWithInfo) ([]BuildEventWithInfo, error) {
	if len(batch) == 0 {
		return nil, nil
	}

	for i := 1; i < len(batch); i++ {
		if batch[i-1].SequenceNumber >= batch[i].SequenceNumber {
			return nil, status.Error(codes.FailedPrecondition, "Sequence numbers out of order")
		}
	}
	minSeq, maxSeq := batch[0].SequenceNumber, batch[len(batch)-1].SequenceNumber

	// Fetches all rows with sequence number in the range but checks the
	// hash go side.
	foundRows, err := r.db.Ent().EventMetadata.Query().
		Where(
			eventmetadata.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
			eventmetadata.SequenceNumberGTE(minSeq),
			eventmetadata.SequenceNumberLTE(maxSeq),
		).
		Select(eventmetadata.FieldSequenceNumber, eventmetadata.FieldEventHash).
		All(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to query existing events")
	}

	existingEvents := make(map[int64][]byte, len(foundRows))
	for _, row := range foundRows {
		existingEvents[row.SequenceNumber] = row.EventHash
	}
	unhandled := make([]BuildEventWithInfo, 0, len(batch)-len(existingEvents))
	for _, info := range batch {
		storedHash, exists := existingEvents[info.SequenceNumber]
		if exists && slices.Equal(storedHash, info.EventHash) {
			continue
		}
		unhandled = append(unhandled, info)
	}

	return unhandled, nil
}
