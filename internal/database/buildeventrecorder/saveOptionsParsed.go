package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) saveOptionsParsed(ctx context.Context, tx *ent.Client, optionsParsed *bes.OptionsParsed) error {
	if optionsParsed == nil {
		return nil
	}
	err := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventOptionsParsed(false),
		).
		SetProcessedEventOptionsParsed(true).
		SetCommandLine(optionsParsed.GetCmdLine()).
		SetExplicitCommandLine(optionsParsed.GetExplicitCmdLine()).
		SetStartupOptions(optionsParsed.GetStartupOptions()).
		SetExplicitStartupOptions(optionsParsed.GetExplicitStartupOptions()).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "OptionsParsed event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to save options parsed to database")
	}
	return nil
}
