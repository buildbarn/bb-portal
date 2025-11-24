package buildeventrecorder

import (
	"context"
	"fmt"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveStarted(ctx context.Context, tx *ent.Client, event *bes.BuildStarted) error {
	if event == nil {
		return nil
	}

	// For some unknown reason bazel emits build event streams for
	// "query" which do not contain an invocation uuid. It still has one
	// as part of the stream, but not in the started at event.
	if event.GetCommand() != "query" || event.GetUuid() != "" {
		// Atleast give an error if we for some reason have a uuid
		// missmatch for any other reason.
		if event.GetUuid() != r.InvocationID {
			return fmt.Errorf("Invocation ID mismatch: event has %q, but expected %q", event.GetUuid(), r.InvocationID)
		}
	}

	var startedAt time.Time
	if event.GetStartTime() != nil {
		startedAt = event.GetStartTime().AsTime()
		if time.Until(startedAt) > 5*time.Minute {
			return fmt.Errorf("Excessive time travel, start time %v is more than 5 minutes in the future", startedAt)
		}
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		startedAt = time.UnixMilli(event.GetStartTimeMillis())
	}

	// Don't use `UpdateOne()` or `UpdateOneID()` as the internal ent
	// implementation also does a query after the update to return the updated
	// object, even if it is not returned by using `Exec()`.
	err := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventStarted(false),
		).
		SetProcessedEventStarted(true).
		SetBazelVersion(event.GetBuildToolVersion()).
		SetStartedAt(startedAt).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "BuildStarted event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to save started event to database")
	}
	return nil
}
