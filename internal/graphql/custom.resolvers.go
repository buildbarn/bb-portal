package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetpair"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testcollection"
	"github.com/buildbarn/bb-portal/internal/graphql/helpers"
	"github.com/buildbarn/bb-portal/internal/graphql/model"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
	"github.com/google/uuid"
)

// Stdout is the resolver for the stdout field.
func (r *actionProblemResolver) Stdout(ctx context.Context, obj *model.ActionProblem) (*model.BlobReference, error) {
	return helpers.BlobReferenceForFile(ctx, r.client, func(ctx context.Context) (*bes.File, error) {
		action, err := helpers.GetAction(ctx, obj.Problem)
		if err != nil {
			return nil, err
		}
		return action.GetStdout(), nil
	})
}

// Stderr is the resolver for the stderr field.
func (r *actionProblemResolver) Stderr(ctx context.Context, obj *model.ActionProblem) (*model.BlobReference, error) {
	return helpers.BlobReferenceForFile(ctx, r.client, func(ctx context.Context) (*bes.File, error) {
		action, err := helpers.GetAction(ctx, obj.Problem)
		if err != nil {
			return nil, err
		}
		return action.GetStderr(), nil
	})
}

// BazelCommand is the resolver for the bazelCommand field.
func (r *bazelInvocationResolver) BazelCommand(ctx context.Context, obj *ent.BazelInvocation) (*model.BazelCommand, error) {
	bcl := obj.Summary.BazelCommandLine
	return &model.BazelCommand{
		// TODO: Scalar ID
		Command:                bcl.Command,
		Executable:             bcl.Executable,
		Residual:               bcl.Residual,
		ExplicitCmdLine:        strings.Join(bcl.ExplicitCmdLine, " "),
		CmdLine:                helpers.StringSliceArrayToPointerArray(bcl.CmdLine),
		StartupOptions:         helpers.StringSliceArrayToPointerArray(bcl.StartUpOptions),
		ExplicitStartupOptions: helpers.StringSliceArrayToPointerArray(bcl.ExplicitStartupOptions),
	}, nil
}

// State is the resolver for the state field.
func (r *bazelInvocationResolver) State(ctx context.Context, obj *ent.BazelInvocation) (*model.BazelInvocationState, error) {
	return &model.BazelInvocationState{
		// TODO: Scalar ID
		BuildEndTime:   obj.EndedAt,
		BuildStartTime: obj.EndedAt,
		ExitCode: &model.ExitCode{
			// TODO: Scalar ID
			Code: obj.Summary.ExitCode.Code,
			Name: obj.Summary.ExitCode.Name,
		},
		BepCompleted: obj.BepCompleted,
	}, nil
}

// User is the resolver for the user field.
func (r *bazelInvocationResolver) User(ctx context.Context, obj *ent.BazelInvocation) (*model.User, error) {
	return &model.User{
		Email: obj.UserEmail,
		Ldap:  obj.UserLdap,
	}, nil
}

// RelatedFiles is the resolver for the relatedFiles field.
func (r *bazelInvocationResolver) RelatedFiles(ctx context.Context, obj *ent.BazelInvocation) ([]*model.NamedFile, error) {
	namedFiles := make([]*model.NamedFile, 0, len(obj.RelatedFiles))
	for relatedFileName, relatedFileURL := range obj.RelatedFiles {
		namedFiles = append(namedFiles, &model.NamedFile{
			Name: relatedFileName,
			URL:  relatedFileURL,
		})
	}
	slices.SortFunc(namedFiles, func(a, b *model.NamedFile) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return namedFiles, nil
}

// Problems is the resolver for the problems field.
func (r *bazelInvocationResolver) Problems(ctx context.Context, obj *ent.BazelInvocation) ([]model.Problem, error) {
	problems, err := obj.QueryProblems().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch problems: %w", err)
	}

	return r.helper.DBProblemsToAPIProblems(ctx, problems)
}

// Profile is the resolver for the profile field.
func (r *bazelInvocationResolver) Profile(ctx context.Context, obj *ent.BazelInvocation) (*model.Profile, error) {
	uriString := strings.ReplaceAll(obj.RelatedFiles[obj.ProfileName], "/"+obj.ProfileName, "")
	uri, err := url.Parse(uriString)
	if err != nil {
		return nil, err
	}

	// Can only provide Profile download for profiles in CAS
	if uri.Scheme != "bytestream" {
		return nil, nil
	}

	components := strings.FieldsFunc(uriString, func(r rune) bool { return r == '/' })

	sizeString := components[len(components)-1]
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		return nil, err
	}
	digest := components[len(components)-2]

	return &model.Profile{
		Name:        obj.ProfileName,
		Digest:      digest,
		SizeInBytes: size,
	}, nil
}

// DownloadURL is the resolver for the downloadURL field.
func (r *blobReferenceResolver) DownloadURL(ctx context.Context, obj *model.BlobReference) (string, error) {
	if obj.Blob == nil {
		panic("Got a name but not blob")
	}
	return fmt.Sprintf("/api/v1/blobs/%d/%s", obj.Blob.ID, url.PathEscape(obj.Name)), nil
}

// SizeInBytes is the resolver for the sizeInBytes field.
func (r *blobReferenceResolver) SizeInBytes(ctx context.Context, obj *model.BlobReference) (*int, error) {
	if obj.Blob == nil {
		return nil, nil
	}
	v := int(obj.Blob.SizeBytes)
	return &v, nil
}

// AvailabilityStatus is the resolver for the availabilityStatus field.
func (r *blobReferenceResolver) AvailabilityStatus(ctx context.Context, obj *model.BlobReference) (model.ActionOutputStatus, error) {
	if obj.Blob == nil {
		return model.ActionOutputStatusUnavailable, nil
	}
	switch obj.Blob.ArchivingStatus {
	case blob.ArchivingStatusQUEUED:
		fallthrough
	case blob.ArchivingStatusARCHIVING:
		return model.ActionOutputStatusProcessing, nil
	case blob.ArchivingStatusSUCCESS:
		return model.ActionOutputStatusAvailable, nil
	case blob.ArchivingStatusFAILED:
		fallthrough
	default:
		return model.ActionOutputStatusUnavailable, nil
	}
}

// EphemeralUrl is the resolver for the downloadURL field.
func (r *blobReferenceResolver) EphemeralURL(ctx context.Context, obj *model.BlobReference) (string, error) {
	if obj.Blob == nil {
		panic("Got a name but not blob")
	}
	if obj.Blob.ArchivingStatus != blob.ArchivingStatusBYTESTREAM {
		return "", nil
	}
	_tmp := strings.Split(obj.Blob.URI, "/blobs/")[1]
	_split := strings.Split(_tmp, "/")
	hash, size := _split[0], _split[1]

	return fmt.Sprintf("/blobs/sha256/file/%s-%s/%s", hash, size, url.PathEscape(obj.Name)), nil
}

// Env is the resolver for the env field.
func (r *buildResolver) Env(ctx context.Context, obj *ent.Build) ([]*model.EnvVar, error) {
	envVars := make([]*model.EnvVar, 0, len(obj.Env))
	for k, v := range obj.Env {
		envVars = append(envVars, &model.EnvVar{
			Key:   k,
			Value: v,
		})
	}
	// Stable sort, mostly for testing purposes.
	slices.SortFunc(envVars, func(a, b *model.EnvVar) int {
		return cmp.Compare(a.Key, b.Key)
	})
	return envVars, nil
}

// DeleteInvocation is the resolver for the deleteInvocation field.
func (r *mutationResolver) DeleteInvocation(ctx context.Context, invocationID uuid.UUID) (bool, error) {
	invocation, err := r.client.BazelInvocation.Query().Where(bazelinvocation.InvocationID(invocationID)).First(ctx)
	if err != nil {
		return false, fmt.Errorf("could not find invocation: %w", err)
	}

	// Delete the invocation
	err = r.client.BazelInvocation.DeleteOne(invocation).Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("could not delete invocation with: %w", err)
	}
	return true, nil
}

// DeleteBuild is the resolver for the deleteBuild field.
func (r *mutationResolver) DeleteBuild(ctx context.Context, buildUUID uuid.UUID) (bool, error) {
	build, err := r.client.Build.Query().Where(build.BuildUUID(buildUUID)).First(ctx)
	if err != nil {
		return false, fmt.Errorf("could not find build: %w", err)
	}
	for _, invocations := range build.Edges.Invocations {
		err = r.client.BazelInvocation.DeleteOne(invocations).Exec(ctx)
		if err != nil {
			return false, fmt.Errorf("could not delete build: %w", err)
		}
	}
	err = r.client.Build.DeleteOne(build).Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("could not delete build with: %w", err)
	}
	return true, nil
}

// DeleteInvocationsBefore is the resolver for the deleteInvocationsBefore field.
func (r *mutationResolver) DeleteInvocationsBefore(ctx context.Context, time time.Time) (*model.DeleteResult, error) {
	// find invocaations before time
	result := &model.DeleteResult{Deleted: 0, Successful: false}
	invocations, err := r.client.BazelInvocation.Query().Where(bazelinvocation.EndedAtLT(time)).All(ctx)
	if err != nil {
		return result, fmt.Errorf("could not find invocations: %w", err)
	}
	result.Found = len(invocations)
	for _, invocation := range invocations {
		err = r.client.BazelInvocation.DeleteOne(invocation).Exec(ctx)
		if err != nil {
			return result, fmt.Errorf("could not delete invocations: %w", err)
		}
		result.Deleted++
	}
	result.Successful = true
	return result, nil
}

// DeleteBuildsBefore is the resolver for the deleteBuildsBefore field.
func (r *mutationResolver) DeleteBuildsBefore(ctx context.Context, time time.Time) (*model.DeleteResult, error) {
	// find builds before time
	result := &model.DeleteResult{Deleted: 0, Successful: false}
	builds, err := r.client.Build.Query().Where(build.TimestampLT(time)).All(ctx)
	if err != nil {
		return result, fmt.Errorf("could not find builds: %w", err)
	}
	result.Found = len(builds)
	for _, build := range builds {
		err = r.client.Build.DeleteOne(build).Exec(ctx)
		if err != nil {
			return result, fmt.Errorf("could not delete build with: %w", err)
		}
		result.Deleted++
	}
	result.Successful = true
	return result, nil
}

// BazelInvocation is the resolver for the bazelInvocation field.
func (r *queryResolver) BazelInvocation(ctx context.Context, invocationID string) (*ent.BazelInvocation, error) {
	invocationUUID, err := uuid.Parse(invocationID)
	if err != nil {
		return nil, fmt.Errorf("invocationID was not a UUID: %w", err)
	}
	return r.client.BazelInvocation.Query().Where(bazelinvocation.InvocationID(invocationUUID)).First(ctx)
}

// GetBuild is the resolver for the getBuild field.
func (r *queryResolver) GetBuild(ctx context.Context, buildURL *string, buildUUID *uuid.UUID) (*ent.Build, error) {
	if buildURL == nil && buildUUID == nil {
		return nil, helpers.ErrOnlyURLOrUUID
	}
	if buildURL != nil && *buildURL != "" && buildUUID != nil {
		return nil, helpers.ErrOnlyURLOrUUID
	}
	if buildURL != nil && *buildURL != "" {
		calculatedBuildUUID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(*buildURL))
		buildUUID = &calculatedBuildUUID
	}
	return r.client.Build.Query().Where(build.BuildUUID(*buildUUID)).First(ctx)
}

// GetUniqueTestLabels is the resolver for the getUniqueTestLabels field.
func (r *queryResolver) GetUniqueTestLabels(ctx context.Context, param *string) ([]*string, error) {
	query := r.client.TestCollection.Query().Limit(100)
	// started := time.Now()
	if param != nil && *param != "" {
		query = query.Where(testcollection.LabelContains(*param))
	}
	res, err := query.Unique(true).Select(testcollection.FieldLabel).Strings(ctx)
	if err != nil {
		return nil, err
	}
	// elapsed := time.Since(started)
	// slog.Info("GetUniqueTestLabels", "elapsed:", elapsed.String())
	return helpers.StringSliceArrayToPointerArray(res), nil
}

// GetUniqueTargetLabels is the resolver for the getUniqueTargetLabels field.
func (r *queryResolver) GetUniqueTargetLabels(ctx context.Context, param *string) ([]*string, error) {
	// started := time.Now()
	query := r.client.TargetPair.Query().Limit(100)
	if param != nil && *param != "" {
		query = query.Where(targetpair.LabelContains(*param))
	}
	res, err := query.Unique(true).Select(targetpair.FieldLabel).Strings(ctx)
	if err != nil {
		return nil, err
	}
	// elapsed := time.Since(started)
	// slog.Info("GetUniqueTargetLabels", "elapsed:", elapsed.String())
	return helpers.StringSliceArrayToPointerArray(res), nil
}

// GetTestDurationAggregation is the resolver for the getTestDurationAggregation field.
func (r *queryResolver) GetTestDurationAggregation(ctx context.Context, label *string) ([]*model.TargetAggregate, error) {
	var result []*model.TargetAggregate
	err := r.client.TestCollection.Query().
		Where(testcollection.LabelContains(*label)).
		GroupBy(testcollection.FieldLabel).
		Aggregate(ent.Count(),
			ent.Sum(testcollection.FieldDurationMs),
			ent.Min(testcollection.FieldDurationMs),
			ent.Max(testcollection.FieldDurationMs)).
		Scan(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTestPassAggregation is the resolver for the getTestPassAggregation field.
func (r *queryResolver) GetTestPassAggregation(ctx context.Context, label *string) ([]*model.TargetAggregate, error) {
	var result []*model.TargetAggregate
	err := r.client.TestCollection.Query().
		Where(testcollection.And(
			testcollection.LabelContains(*label),
			testcollection.OverallStatusEQ(testcollection.OverallStatusPASSED),
		)).
		GroupBy(testcollection.FieldLabel).
		Aggregate(ent.Count()).
		Scan(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTargetDurationAggregation is the resolver for the getTargetDurationAggregation field.
func (r *queryResolver) GetTargetDurationAggregation(ctx context.Context, label *string) ([]*model.TargetAggregate, error) {
	panic(fmt.Errorf("not implemented: GetTargetDurationAggregation - getTargetDurationAggregation"))
}

// GetTargetPassAggregation is the resolver for the getTargetPassAggregation field.
func (r *queryResolver) GetTargetPassAggregation(ctx context.Context, label *string) ([]*model.TargetAggregate, error) {
	panic(fmt.Errorf("not implemented: GetTargetPassAggregation - getTargetPassAggregation"))
}

// GetTestsWithOffset is the resolver for the getTestsWithOffset field.
func (r *queryResolver) GetTestsWithOffset(ctx context.Context, label *string, offset, limit *int, sortBy, direction *string) (*model.TestGridResult, error) {
	maxLimit := 100
	take := 10
	skip := 0
	if limit != nil {
		if *limit > maxLimit {
			return nil, fmt.Errorf("limit cannot exceed %d", maxLimit)
		}
		take = *limit
	}
	if offset != nil {
		skip = *offset
	}

	var result []*model.TestGridRow
	query := r.client.TestCollection.Query()

	if label != nil && *label != "" {
		query = query.Where(testcollection.LabelContains(*label))
	}

	err := query.
		Limit(take).
		Offset(skip).
		GroupBy(testcollection.FieldLabel).
		Aggregate(ent.Count(),
			ent.As(ent.Mean(testcollection.FieldDurationMs), "avg"),
			ent.Sum(testcollection.FieldDurationMs),
			ent.Min(testcollection.FieldDurationMs),
			ent.Max(testcollection.FieldDurationMs)).
		Scan(ctx, &result)
	if err != nil {
		return nil, err
	}
	totalCount := 0
	response := &model.TestGridResult{
		Result: result,
		Total:  &totalCount,
	}
	return response, nil
}

// GetTargetsWithOffset is the resolver for the GetTargetsWithOffset field.
func (r *queryResolver) GetTargetsWithOffset(ctx context.Context, label *string, offset, limit *int, sortBy, direction *string) (*model.TargetGridResult, error) {
	maxLimit := 100
	take := 10
	skip := 0
	if limit != nil {
		if *limit > maxLimit {
			return nil, fmt.Errorf("limit cannot exceed %d", maxLimit)
		}
		take = *limit
	}
	if offset != nil {
		skip = *offset
	}

	var result []*model.TargetGridRow
	query := r.client.TargetPair.Query()

	if label != nil && *label != "" {
		query = query.Where(targetpair.LabelContains(*label))
	}

	err := query.
		Limit(take).
		Offset(skip).
		GroupBy(targetpair.FieldLabel).
		Aggregate(ent.Count(),
			ent.As(ent.Mean(targetpair.FieldDurationInMs), "avg"),
			ent.Sum(targetpair.FieldDurationInMs),
			ent.Min(targetpair.FieldDurationInMs),
			ent.Max(targetpair.FieldDurationInMs)).
		Scan(ctx, &result)
	if err != nil {
		return nil, err
	}
	totalCount := 0
	response := &model.TargetGridResult{
		Result: result,
		Total:  &totalCount,
	}
	return response, nil
}

// GetAveragePassPercentageForLabel is the resolver for the getAveragePassPercentageForLabel field.
func (r *queryResolver) GetAveragePassPercentageForLabel(ctx context.Context, label string) (*float64, error) {
	// TODO: maybe there is a more elegant/faster way to do this with aggregaate
	passCount, err := r.client.TestCollection.Query().
		Where(testcollection.And(
			testcollection.LabelEQ(label),
			testcollection.OverallStatusEQ(testcollection.OverallStatusPASSED),
		)).Count(ctx)
	if err != nil {
		return nil, err
	}
	totalCount, err := r.client.TestCollection.Query().
		Where(testcollection.LabelEQ(label)).
		Count(ctx)
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		result := 0.0
		return helpers.GetFloatPointer(&result), nil
	}
	result := float64(passCount/totalCount) * 100.0
	return helpers.GetFloatPointer(&result), nil
}

// ActionLogOutput is the resolver for the actionLogOutput field.
func (r *testResultResolver) ActionLogOutput(ctx context.Context, obj *model.TestResult) (*model.BlobReference, error) {
	return helpers.GetTestResultActionLogOutput(ctx, r.client, obj)
}

// UndeclaredTestOutputs is the resolver for the undeclaredTestOutputs field.
func (r *testResultResolver) UndeclaredTestOutputs(ctx context.Context, obj *model.TestResult) (*model.BlobReference, error) {
	return helpers.GetTestResultUndeclaredTestOutputs(ctx, r.client, obj)
}

// ActionProblem returns ActionProblemResolver implementation.
func (r *Resolver) ActionProblem() ActionProblemResolver { return &actionProblemResolver{r} }

// BlobReference returns BlobReferenceResolver implementation.
func (r *Resolver) BlobReference() BlobReferenceResolver { return &blobReferenceResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// TestResult returns TestResultResolver implementation.
func (r *Resolver) TestResult() TestResultResolver { return &testResultResolver{r} }

type (
	actionProblemResolver struct{ *Resolver }
	blobReferenceResolver struct{ *Resolver }
	mutationResolver      struct{ *Resolver }
	testResultResolver    struct{ *Resolver }
)
