package buildeventrecorder

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	bescore "github.com/bazelbuild/bazel/src/main/protobuf"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func parseEnvVarsFromSectionOptions(section *bescore.CommandLineSection, destMap *map[string]string) {
	if section.GetOptionList() == nil {
		return
	}
	options := section.GetOptionList().GetOption()
	for _, option := range options {
		if option.GetOptionName() != "client_env" {
			// Only looking for env vars from the client env
			continue
		}
		envPair := option.GetOptionValue()
		equalIndex := strings.Index(envPair, "=")
		if equalIndex <= 0 {
			// Skip anything missing an equals sign. The env vars come in the format key=value
			continue
		}
		envName := envPair[:equalIndex]
		envValue := envPair[equalIndex+1:]
		(*destMap)[envName] = envValue
	}
}

func parseProfileNameFromSectionOptions(section *bescore.CommandLineSection) string {
	if section.GetOptionList() != nil {
		options := section.GetOptionList().GetOption()
		for _, option := range options {
			if option.GetOptionName() == "profile" {
				return option.GetOptionValue()
			}
		}
	}

	// Default value if --profile is not set
	return "command.profile.gz"
}

// envToI
func envToI(envVars map[string]string, name string) (int, error) {
	res, err := strconv.Atoi(envVars[name])
	if err != nil {
		return 0, util.StatusWrapf(err, "failed to parse %s (value: %s) as an int", name, envVars[name])
	}
	return res, nil
}

func (r *BuildEventRecorder) recordSourceControl(ctx context.Context, tx *ent.Tx, envVars map[string]string) error {
	sc := tx.SourceControl.Create().
		SetBazelInvocationID(r.InvocationDbID)

	shouldSave := false

	if instanceURL, ok := envVars["GITHUB_SERVER_URL"]; ok {
		// Github
		sc.SetProvider("GITHUB")
		sc.SetInstanceURL(instanceURL)
		shouldSave = true
	} else if instanceURL, ok := envVars["CI_SERVER_URL"]; ok {
		// Gitlab
		sc.SetProvider("GITLAB")
		sc.SetInstanceURL(instanceURL)
		shouldSave = true
	}

	if repo, ok := envVars["GITHUB_REPOSITORY"]; ok {
		// Github
		sc.SetRepo(repo)
		shouldSave = true
	} else if repo, ok := envVars["CI_PROJECT_PATH"]; ok {
		// Gitlab
		sc.SetRepo(repo)
		shouldSave = true
	}

	if commitSha, ok := envVars["GITHUB_SHA"]; ok {
		// Github
		sc.SetCommitSha(commitSha)
		shouldSave = true
	} else if commitSha, ok := envVars["CI_COMMIT_SHA"]; ok {
		// Gitlab
		sc.SetCommitSha(commitSha)
		shouldSave = true
	}

	if refs, ok := envVars["GITHUB_REF"]; ok {
		// Github
		sc.SetRefs(refs)
		shouldSave = true
	} else if refs, ok := envVars["CI_COMMIT_REF_NAME"]; ok {
		// Gitlab
		sc.SetRefs(refs)
		shouldSave = true
	}

	if user, ok := envVars["GITHUB_ACTOR"]; ok {
		// Github
		sc.SetActor(user)
		shouldSave = true
	} else if user, ok := envVars["GITLAB_USER_LOGIN"]; ok {
		// Gitlab
		sc.SetActor(user)
		shouldSave = true
	}

	if eventName, ok := envVars["GITHUB_EVENT_NAME"]; ok {
		// Github
		sc.SetEventName(eventName)
		shouldSave = true
	} else if eventName, ok := envVars["CI_PIPELINE_SOURCE"]; ok {
		// Gitlab
		sc.SetEventName(eventName)
		shouldSave = true
	}

	if workflowName, ok := envVars["GITHUB_WORKFLOW"]; ok {
		// Github
		sc.SetWorkflow(workflowName)
		shouldSave = true
	} else if workflowName, ok := envVars["CI_JOB_NAME"]; ok {
		// Gitlab
		sc.SetWorkflow(workflowName)
		shouldSave = true
	}

	if runID, ok := envVars["GITHUB_RUN_ID"]; ok {
		// Github
		sc.SetRunID(runID)
		shouldSave = true
	} else if runID, ok := envVars["CI_JOB_ID"]; ok {
		// Gitlab
		sc.SetRunID(runID)
		shouldSave = true
	}

	if runNumber, ok := envVars["GITHUB_RUN_NUMBER"]; ok {
		// Github
		sc.SetRunNumber(runNumber)
		shouldSave = true
	}

	if job, ok := envVars["GITHUB_JOB"]; ok {
		// Github
		sc.SetJob(job)
		shouldSave = true
	} else if job, ok := envVars["CI_JOB_STAGE"]; ok {
		// Gitlab
		sc.SetJob(job)
		shouldSave = true
	}

	if action, ok := envVars["GITHUB_ACTION"]; ok {
		// Github
		sc.SetAction(action)
		shouldSave = true
	}

	if runnerName, ok := envVars["RUNNER_NAME"]; ok {
		// Github
		sc.SetRunnerName(runnerName)
		shouldSave = true
	} else if runnerName, ok := envVars["CI_RUNNER_DESCRIPTION"]; ok {
		// Gitlab
		sc.SetRunnerName(runnerName)
		shouldSave = true
	}

	if runnerArch, ok := envVars["RUNNER_ARCH"]; ok {
		// Github
		sc.SetRunnerArch(runnerArch)
		shouldSave = true
	} else if runnerInfo, ok := envVars["CI_RUNNER_EXECUTABLE_ARCH"]; ok {
		// Gitlab
		// This variable comes in the format "os/arch", e.g. "linux/amd64"
		_, arch, _ := strings.Cut(runnerInfo, "/")
		sc.SetRunnerArch(arch)
		shouldSave = true
	}

	if runnerOs, ok := envVars["RUNNER_OS"]; ok {
		// Github
		sc.SetRunnerOs(runnerOs)
		shouldSave = true
	} else if runnerInfo, ok := envVars["CI_RUNNER_EXECUTABLE_ARCH"]; ok {
		// Gitlab
		// This variable comes in the format "os/arch", e.g. "linux/amd64"
		os, _, _ := strings.Cut(runnerInfo, "/")
		sc.SetRunnerOs(os)
		shouldSave = true
	}

	if workspace, ok := envVars["GITHUB_WORKSPACE"]; ok {
		// Github
		sc.SetWorkspace(workspace)
		shouldSave = true
	} else if workspace, ok := envVars["CI_PROJECT_DIR"]; ok {
		// Gitlab
		sc.SetWorkspace(workspace)
		shouldSave = true
	}

	if !shouldSave {
		// No source control data found, skip creating the entry.
		return nil
	}

	err := sc.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to save source control data to database")
	}

	return nil
}

func (r *BuildEventRecorder) getBuildURL(envVars map[string]string) string {
	if buildURL, ok := envVars["BUILD_URL"]; ok {
		return buildURL
	}
	if buildURL, ok := envVars["CI_PIPELINE_URL"]; ok {
		// Gitlab
		return buildURL
	}
	if envVars["GITHUB_SERVER_URL"] != "" && envVars["GITHUB_REPOSITORY"] != "" && envVars["GITHUB_RUN_ID"] != "" {
		// Github
		return fmt.Sprintf("%s/%s/actions/runs/%s", envVars["GITHUB_SERVER_URL"], envVars["GITHUB_REPOSITORY"], envVars["GITHUB_RUN_ID"])
	}
	return ""
}

func (r *BuildEventRecorder) recordBuild(ctx context.Context, tx *ent.Tx, envVars map[string]string) error {
	buildURL := r.getBuildURL(envVars)
	if buildURL == "" {
		return nil
	}

	invocation, err := tx.BazelInvocation.Query().
		Where(bazelinvocation.IDEQ(r.InvocationDbID)).
		Select(bazelinvocation.FieldStartedAt).
		Only(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to query for invocation start time")
	}

	buildUUID := common.CalculateBuildUUID(buildURL, r.InstanceName)

	buildDb, err := tx.Build.Query().
		Where(build.BuildUUIDEQ(buildUUID)).
		Only(ctx)

	switch {
	case ent.IsNotFound(err):
		err = tx.Build.Create().
			SetBuildURL(buildURL).
			SetBuildUUID(buildUUID).
			SetTimestamp(invocation.StartedAt).
			SetInstanceNameID(r.InstanceNameDbID).
			AddInvocationIDs(r.InvocationDbID).
			Exec(ctx)
		if err != nil {
			return util.StatusWrap(err, "Failed to save build information")
		}
		return nil
	case err == nil:
		update := tx.Build.
			Update().
			Where(build.ID(buildDb.ID)).
			AddInvocationIDs(r.InvocationDbID)
		if invocation.StartedAt.Before(buildDb.Timestamp) {
			update.SetTimestamp(invocation.StartedAt)
		}
		if err := update.Exec(ctx); err != nil {
			return util.StatusWrap(err, "Failed to update build information")
		}
		return nil
	default:
		return util.StatusWrap(err, "Failed to query for existing build")
	}
}

func (r *BuildEventRecorder) saveStructuredCommandLine(ctx context.Context, tx *ent.Tx, structuredCommandLine *bescore.CommandLine) error {
	if structuredCommandLine == nil {
		return nil
	}
	if structuredCommandLine.GetCommandLineLabel() != "original" {
		return nil
	}

	update := tx.BazelInvocation.
		Update().
		Where(
			bazelinvocation.ID(r.InvocationDbID),
			bazelinvocation.ProcessedEventStructuredCommandLine(false),
		).
		SetProcessedEventStructuredCommandLine(true)

	envVars := map[string]string{}

	sections := structuredCommandLine.GetSections()
	for _, section := range sections {
		label := section.GetSectionLabel()
		if label == "command options" {
			parseEnvVarsFromSectionOptions(section, &envVars)
			update.SetProfileName(parseProfileNameFromSectionOptions(section))
		} else if section.GetChunkList() != nil {
			sectionChunksStr := strings.Join(section.GetChunkList().GetChunk(), " ")
			switch label {
			case "executable":
				update.SetCommandLineExecutable(sectionChunksStr)
			case "command":
				update.SetCommandLineCommand(sectionChunksStr)
			case "residual":
				update.SetCommandLineResidual(sectionChunksStr)
			}
		}
	}

	// Parse Gerrit change number if available.
	if changeNumberStr, ok := envVars["GERRIT_CHANGE_NUMBER"]; ok && changeNumberStr != "" {
		changeNumber, err := envToI(envVars, "GERRIT_CHANGE_NUMBER")
		if err != nil {
			return util.StatusWrap(err, "failed to parse GERRIT_CHANGE_NUMBER from structured command line")
		}
		update.SetChangeNumber(changeNumber)
	}

	// Parse Gerrit patchset number if available.
	if patchsetNumberStr, ok := envVars["GERRIT_PATCHSET_NUMBER"]; ok && patchsetNumberStr != "" {
		patchsetNumber, err := envToI(envVars, "GERRIT_PATCHSET_NUMBER")
		if err != nil {
			return util.StatusWrap(err, "failed to parse GERRIT_PATCHSET_NUMBER from structured command line")
		}
		update.SetPatchsetNumber(patchsetNumber)
	}

	// Set Hostname
	if hostNameVal, ok := envVars["BB_PORTAL_HOSTNAME"]; ok {
		update.SetHostname(hostNameVal)
	} else if hostNameVal, ok := envVars["HOSTNAME"]; ok {
		update.SetHostname(hostNameVal)
	} else if hostNameVal, ok := envVars["RUNNER_NAME"]; ok {
		update.SetHostname(hostNameVal)
	}

	// Set CI Worker Role from environment variables (can also come from metadata)
	if isCiWorkerVal, ok := envVars["BB_PORTAL_IS_CI_WORKER"]; ok {
		update.SetIsCiWorker(isCiWorkerVal == "true" || isCiWorkerVal == "True" || isCiWorkerVal == "TRUE")
	}

	// github/gitlab actions default env var
	if isCiWorkerVal, ok := envVars["CI"]; ok {
		update.SetIsCiWorker(isCiWorkerVal == "true")
	}

	// Set Step Label from environment variables
	if stepLabelVal, ok := envVars["BB_PORTAL_STEP_LABEL"]; ok {
		update.SetStepLabel(stepLabelVal)
	} else if glWfVal, ok := envVars["CI_JOB_STAGE"]; ok {
		// Gitlab default step label to workfow + job
		if glJobNameVal, ok := envVars["CI_JOB_NAME"]; ok {
			update.SetStepLabel(glWfVal + "+" + glJobNameVal)
		} else {
			update.SetStepLabel(glWfVal)
		}
	} else if ghWfVal, ok := envVars["GITHUB_WORKFLOW"]; ok {
		// Github default step label to workfow + job
		if ghJobNameVal, ok := envVars["GITHUB_JOB"]; ok {
			update.SetStepLabel(ghWfVal + "+" + ghJobNameVal)
		} else {
			update.SetStepLabel(ghWfVal)
		}
	}

	if user, ok := envVars["GITHUB_ACTOR"]; ok {
		// Github
		update.SetUserLdap(user)
	} else if user, ok := envVars["GITLAB_USER_LOGIN"]; ok {
		// Gitlab
		update.SetUserLdap(user)
	} else if user, ok := envVars["USER"]; ok {
		// Local
		update.SetUserLdap(user)
	}

	if email, ok := envVars["GITLAB_USER_EMAIL"]; ok {
		// Gitlab
		update.SetUserEmail(email)
	}

	err := update.Exec(ctx)
	if ent.IsNotFound(err) {
		return util.StatusWrapf(err, "StructuredCommandline event has already been processed for invocation %s", r.InvocationID)
	}
	if err != nil {
		return util.StatusWrap(err, "Failed to save structured command line to database")
	}

	err = r.recordBuild(ctx, tx, envVars)
	if err != nil {
		return util.StatusWrap(err, "Failed to save build data")
	}

	err = r.recordSourceControl(ctx, tx, envVars)
	if err != nil {
		return util.StatusWrap(err, "Failed to save source control data")
	}
	return nil
}
