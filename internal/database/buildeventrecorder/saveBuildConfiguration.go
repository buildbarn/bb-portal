package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *BuildEventRecorder) saveBuildConfiguration(ctx context.Context, tx *ent.Client, configuration *bes.Configuration) error {
	if configuration == nil {
		return nil
	}

	err := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
		).
		SetCPU(configuration.GetCpu()).
		SetPlatformName(configuration.GetPlatformName()).
		SetConfigurationMnemonic(configuration.GetMnemonic()).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to update bazel invocation with build configuration BES message")
	}
	return nil
}
