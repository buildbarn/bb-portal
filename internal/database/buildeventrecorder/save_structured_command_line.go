package buildeventrecorder

import (
	"context"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildtag"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/pkg/invocation"
	"github.com/buildbarn/bb-portal/pkg/invocationmetadataextraction"
	"github.com/buildbarn/bb-storage/pkg/util"

	bes "github.com/bazelbuild/bazel/src/main/protobuf"
)

func parseEnvVarsFromSectionOptions(section *bes.CommandLineSection) map[string]string {
	if section.GetOptionList() == nil {
		return nil
	}
	ret := make(map[string]string)
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
		ret[envName] = envValue
	}
	return ret
}

func parseProfileNameFromSectionOptions(section *bes.CommandLineSection) string {
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

func (r *buildEventRecorder) recordSourceControl(ctx context.Context, tx *ent.Client, sourceControls []invocationmetadataextraction.SourceControl) error {
	scBuilders := make([]*ent.SourceControlCreate, 0, len(sourceControls))
	for _, sc := range sourceControls {
		create := tx.SourceControl.Create().SetBazelInvocationID(r.InvocationDbID)
		shouldSave := false
		if repo := sc.Repo; repo != nil && *repo != "" {
			create.SetRepo(*repo)
			shouldSave = true
		}
		if repoURL := sc.RepoURL; repoURL != nil && *repoURL != "" {
			create.SetRepoURL(*repoURL)
			shouldSave = true
		}
		if ref := sc.Ref; ref != nil && *ref != "" {
			create.SetRef(*ref)
			shouldSave = true
		}
		if refURL := sc.RefURL; refURL != nil && *refURL != "" {
			create.SetRefURL(*refURL)
			shouldSave = true
		}
		if commit := sc.Commit; commit != nil && *commit != "" {
			create.SetCommit(*commit)
			shouldSave = true
		}
		if commitURL := sc.CommitURL; commitURL != nil && *commitURL != "" {
			create.SetCommitURL(*commitURL)
			shouldSave = true
		}
		if shouldSave {
			scBuilders = append(scBuilders, create)
		}
	}

	if err := tx.SourceControl.CreateBulk(scBuilders...).Exec(ctx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert source controls to the database")
	}
	return nil
}

func (r *buildEventRecorder) recordBuildTags(ctx context.Context, tx *ent.Client, buildDbID int64, tags map[string]string) error {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}

	// The keys need to be sorted. Otherwise they are inserted in a random
	// order and the tests get sad.
	sort.Strings(keys)

	tagBuilders := make([]*ent.BuildTagCreate, 0, len(tags))
	for _, key := range keys {
		value := tags[key]
		if key != "" && value != "" {
			create := tx.BuildTag.Create().
				SetBuildID(buildDbID).
				SetKey(key).
				SetValue(value)
			tagBuilders = append(tagBuilders, create)
		}
	}

	err := tx.BuildTag.CreateBulk(tagBuilders...).
		OnConflictColumns(buildtag.FieldBuildID, buildtag.FieldKey, buildtag.FieldValue).
		DoNothing().
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to bulk insert invocation tags to the database")
	}
	return nil
}

func (r *buildEventRecorder) recordBuild(ctx context.Context, tx database.Tx, invocationMetadata *invocationmetadataextraction.InvocationMetadata) error {
	if r.buildKey == "" {
		return nil
	}
	buildID, ok := invocationMetadata.BuildTags[r.buildKey]
	if !ok {
		return nil
	}
	if buildID == "" {
		return nil
	}

	buildUUID := common.CalculateBuildUUID(buildID, r.InstanceName)

	buildDbID, err := tx.Ent().Build.Create().
		SetInstanceNameID(r.InstanceNameDbID).
		SetBuildUUID(buildUUID).
		SetTimestamp(time.Now()).
		AddInvocationIDs(r.InvocationDbID).
		OnConflictColumns(build.FieldBuildUUID).
		Ignore().
		ID(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to upsert build")
	}

	err = tx.Sqlc().UpdateBuildTimestampFromInvocation(ctx, r.InvocationDbID)
	if err != nil {
		return util.StatusWrap(err, "Failed to update build timestamp")
	}

	if err := r.recordBuildTags(ctx, tx.Ent(), buildDbID, invocationMetadata.BuildTags); err != nil {
		return util.StatusWrap(err, "Failed to record build tags")
	}
	return nil
}

func (r *buildEventRecorder) recordInvocationTags(ctx context.Context, tx *ent.Client, tags map[string]string) error {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}

	// The keys need to be sorted. Otherwise they are inserted in a random
	// order and the tests get sad.
	sort.Strings(keys)

	tagBuilders := make([]*ent.InvocationTagCreate, 0, len(tags))
	for _, key := range keys {
		value := tags[key]
		if key != "" && value != "" {
			create := tx.InvocationTag.Create().
				SetBazelInvocationID(r.InvocationDbID).
				SetKey(key).
				SetValue(value)
			tagBuilders = append(tagBuilders, create)
		}
	}

	err := tx.InvocationTag.CreateBulk(tagBuilders...).Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to bulk insert invocation tags to the database")
	}
	return nil
}

func (r *buildEventRecorder) saveStructuredCommandLine(ctx context.Context, tx database.Tx, buildEvent *bes.CommandLine) error {
	envVars, data, profileName := parseSections(buildEvent.GetSections())

	switch buildEvent.CommandLineLabel {
	case "canonical":
		update := tx.Ent().BazelInvocation.
			UpdateOneID(r.InvocationDbID).
			SetProfileName(profileName).
			SetCanonicalCommandLine(&data)

		invocationMetadata := invocationmetadataextraction.ExtractInvocationMetadata(r.dataExtractors.InvocationMetadataExtractor, envVars)
		if invocationMetadata != nil {
			if username := invocationMetadata.Username; username != nil && *username != "" {
				update.SetUsername(*username)
			}
			if hostname := invocationMetadata.Hostname; hostname != nil && *hostname != "" {
				update.SetHostname(*hostname)
			}

			if err := r.recordSourceControl(ctx, tx.Ent(), invocationMetadata.SourceControls); err != nil {
				return util.StatusWrap(err, "Failed to save source control information")
			}
			if err := r.recordBuild(ctx, tx, invocationMetadata); err != nil {
				return util.StatusWrap(err, "Failed to record build")
			}
			if err := r.recordInvocationTags(ctx, tx.Ent(), invocationMetadata.InvocationTags); err != nil {
				return util.StatusWrap(err, "Failed to record invocation tags")
			}
		}
		if err := update.Exec(ctx); err != nil {
			return util.StatusWrapf(err, "Failed to update invocation with data from StructuredCommandLine event")
		}
	case "original":
		_, err := tx.Ent().BazelInvocation.UpdateOneID(r.InvocationDbID).SetOriginalCommandLine(&data).Save(ctx)
		if err != nil {
			return util.StatusWrapf(err, "Failed to save command line data")
		}
	}

	return nil
}

func parseSections(buildEventSections []*bes.CommandLineSection) (envVars map[string]string, data invocation.CommandLineData, profileName string) {
	// Explicitly set empty slices rather than rely on nil equivalence as the
	// json encoder will explicitly map nil slices to nil rather than empty
	// list.
	data = invocation.CommandLineData{
		Options:        []invocation.CommandLineOption{},
		StartupOptions: []invocation.CommandLineOption{},
		Residual:       []string{},
	}

	for _, section := range buildEventSections {
		switch section.GetSectionLabel() {
		case "executable":
			if list := section.GetChunkList(); list != nil && len(list.Chunk) > 0 {
				data.Executable = list.Chunk[0]
			}
		case "command":
			if list := section.GetChunkList(); list != nil && len(list.Chunk) > 0 {
				data.Command = list.Chunk[0]
			}
		case "startup options":
			data.StartupOptions = extractOptions(section.GetOptionList().GetOption())
		case "command options":
			data.Options = extractOptions(section.GetOptionList().GetOption())
			envVars = parseEnvVarsFromSectionOptions(section)
			profileName = parseProfileNameFromSectionOptions(section)
		case "residual":
			if list := section.GetChunkList(); list != nil {
				data.Residual = append(data.Residual, list.Chunk...)
			}
		}
	}
	return envVars, data, profileName
}

func extractOptions(protoOptions []*bes.Option) []invocation.CommandLineOption {
	result := make([]invocation.CommandLineOption, 0, len(protoOptions))
	for _, opt := range protoOptions {
		// Tags marked as hidden may contain sensitive information or
		// information not relevant to a user.
		if slices.Contains(opt.MetadataTags, bes.OptionMetadataTag_HIDDEN) {
			continue
		}
		result = append(result, invocation.CommandLineOption{
			Option: opt.GetOptionName(),
			Value:  opt.GetOptionValue(),
		})
	}
	return result
}
