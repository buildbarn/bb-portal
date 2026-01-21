package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func (r *buildEventRecorder) saveBuildMetadata(ctx context.Context, tx *ent.Client, metadata *bes.BuildMetadata) error {
	metadataMap := metadata.GetMetadata()
	if metadataMap == nil {
		return nil
	}

	update := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventBuildMetadata(false),
		).
		SetProcessedEventBuildMetadata(true)

	if stepLabel, ok := metadataMap["BUILD_STEP_LABEL"]; ok {
		update.SetStepLabel(stepLabel)
	}
	if userEmail, ok := metadataMap["user_email"]; ok {
		update.SetUserEmail(userEmail)
	}
	if userLdap, ok := metadataMap["user_ldap"]; ok {
		update.SetUserLdap(userLdap)
	}
	if isCiWorkerVal, ok := metadataMap["is_ci_worker"]; ok {
		update.SetIsCiWorker(isCiWorkerVal == "true" || isCiWorkerVal == "True" || isCiWorkerVal == "TRUE")
	}
	if hostnameVal, ok := metadataMap["hostname"]; ok {
		update.SetHostname(hostnameVal)
	}

	err := update.Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "BuildMetadata event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to update bazel invocation with build metadata BES message")
	}
	return nil
}
