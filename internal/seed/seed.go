package seed

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/buildbarn/bb-portal/internal/seed/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// The maximum number of identifiers that can be allocated per instance for
// each of the seeded tables. Identifiers are assigned deterministically in
// non-overlapping per-instance ranges, so that parallel seeders never
// collide. The values below only need to be generous upper bounds.
const (
	maxInvocationsPerInstance       = 1_000_000
	maxTargetsPerInstance           = 1_000_000
	maxInvocationTargetsPerInstance = 100_000_000
	maxSourceControlsPerInvocation  = 100
	maxInvocationTagsPerInvocation  = 100
	maxBuildTagsPerInstance         = 100
)

// idAllocator assigns non-overlapping ranges of identifiers to each seeded
// instance: instance N starts allocating at N*capacity + 1.
type idAllocator struct {
	capacity int64
}

func (a idAllocator) forInstance(instanceID int64) *instanceIDAllocator {
	return &instanceIDAllocator{
		next: instanceID*a.capacity + 1,
	}
}

type instanceIDAllocator struct {
	next int64
}

func (a *instanceIDAllocator) nextID() int64 {
	id := a.next
	a.next++
	return id
}

type idAllocators struct {
	invocations        idAllocator
	targets            idAllocator
	builds             idAllocator
	buildTags          idAllocator
	sourceControls     idAllocator
	invocationTags     idAllocator
	connectionMetadata idAllocator
	eventMetadata      idAllocator
	invocationTargets  idAllocator
	testTargets        idAllocator
	testSummaries      idAllocator
	testResults        idAllocator
}

func newIDAllocators() *idAllocators {
	return &idAllocators{
		invocations:        idAllocator{capacity: maxInvocationsPerInstance},
		targets:            idAllocator{capacity: maxTargetsPerInstance},
		builds:             idAllocator{capacity: maxInvocationsPerInstance},
		buildTags:          idAllocator{capacity: maxBuildTagsPerInstance},
		sourceControls:     idAllocator{capacity: maxInvocationsPerInstance * maxSourceControlsPerInvocation},
		invocationTags:     idAllocator{capacity: maxInvocationsPerInstance * maxInvocationTagsPerInvocation},
		connectionMetadata: idAllocator{capacity: maxInvocationsPerInstance},
		eventMetadata:      idAllocator{capacity: maxInvocationsPerInstance},
		invocationTargets:  idAllocator{capacity: maxInvocationTargetsPerInstance},
		testTargets:        idAllocator{capacity: maxTargetsPerInstance},
		testSummaries:      idAllocator{capacity: maxInvocationTargetsPerInstance},
		testResults:        idAllocator{capacity: maxInvocationTargetsPerInstance},
	}
}

// advanceSequences advances the identity sequence of all seeded tables, so
// that future inserts by the application do not collide with the explicitly
// assigned seed identifiers.
func advanceSequences(ctx context.Context, db *sqlc.Queries) error {
	for _, fn := range []func(ctx context.Context) error{
		db.SeedAdvanceInstanceNamesSequence,
		db.SeedAdvanceAuthenticatedUsersSequence,
		db.SeedAdvanceBazelInvocationsSequence,
		db.SeedAdvanceSourceControlsSequence,
		db.SeedAdvanceInvocationTagsSequence,
		db.SeedAdvanceConnectionMetadataSequence,
		db.SeedAdvanceEventMetadataSequence,
		db.SeedAdvanceBuildsSequence,
		db.SeedAdvanceBuildTagsSequence,
		db.SeedAdvanceTargetsSequence,
		db.SeedAdvanceTestTargetsSequence,
		db.SeedAdvanceInvocationTargetsSequence,
		db.SeedAdvanceTestSummariesSequence,
		db.SeedAdvanceTestResultsSequence,
	} {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func seed(ctx context.Context, db *sqlc.Queries, params Parameters) error {
	instanceNames, err := seedInstanceNames(ctx, db, params.instances)
	if err != nil {
		return util.StatusWrap(err, "Failed to seed instance names")
	}

	users, err := seedUsers(ctx, db, params.users)
	if err != nil {
		return util.StatusWrap(err, "Failed to seed authenticated users")
	}

	buildTagKeys, err := getTagKeys("buildTags", params.invocationMetadataExtractor)
	if err != nil {
		return util.StatusWrap(err, "Failed to extract build tag keys")
	}
	if params.buildKey != "" && !slices.Contains(buildTagKeys, params.buildKey) {
		slog.Warn("Build tags do not include the specified build key", "buildKey", params.buildKey)
	}

	invocationTagKeys, err := getTagKeys("invocationTags", params.invocationMetadataExtractor)
	if err != nil {
		return util.StatusWrap(err, "Failed to extract invocation tag keys")
	}

	sourceControlCount, err := getSourceControlCount(params.invocationMetadataExtractor)
	if err != nil {
		return util.StatusWrap(err, "Failed to extract source control keys")
	}

	idAllocators := newIDAllocators()
	var g errgroup.Group
	for _, instance := range instanceNames {
		g.Go(func() error {
			seeder := newInstanceSeeder(db, instance, users, buildTagKeys, invocationTagKeys, sourceControlCount, params, idAllocators)
			invocations, err := seeder.seedBazelInvocations(ctx)
			if err != nil {
				return util.StatusWrap(err, "Failed to seed bazel invocations")
			}

			if err := seeder.seedBuild(ctx, invocations, int(instance.ID)-1); err != nil {
				return util.StatusWrap(err, "Failed to seed build")
			}

			targets, err := seeder.seedTargets(ctx)
			if err != nil {
				return util.StatusWrap(err, "Failed to seed targets")
			}
			return seeder.seedInvocationTargets(ctx, invocations, targets)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func seedInstanceNames(ctx context.Context, db *sqlc.Queries, n int32) ([]sqlc.InstanceName, error) {
	if n == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, n)
	names := make([]string, 0, n)
	for i := range n {
		ids = append(ids, int64(i+1))
		names = append(names, fmt.Sprintf("instance-%d", i))
	}

	if _, err := db.SeedCreateInstanceNames(ctx, sqlc.SeedCreateInstanceNamesParams{
		Ids:   ids,
		Names: names,
	}); err != nil {
		return nil, err
	}

	instanceNames := make([]sqlc.InstanceName, 0, n)
	for i := range n {
		instanceNames = append(instanceNames, sqlc.InstanceName{
			ID:   ids[i],
			Name: names[i],
		})
	}
	return instanceNames, nil
}

func seedUsers(ctx context.Context, db *sqlc.Queries, n int32) ([]int64, error) {
	if n == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, n)
	userUUIDs := make([]uuid.UUID, 0, n)
	externalIDs := make([]string, 0, n)
	displayNames := make([]string, 0, n)
	for i := range n {
		user := fmt.Sprintf("seed_user_%d", i)
		ids = append(ids, int64(i+1))
		userUUIDs = append(userUUIDs, uuid.NewSHA1(uuid.NameSpaceURL, []byte(user)))
		externalIDs = append(externalIDs, user)
		displayNames = append(displayNames, user)
	}

	if _, err := db.SeedCreateAuthenticatedUsers(ctx, sqlc.SeedCreateAuthenticatedUsersParams{
		Ids:          ids,
		UserUuids:    userUUIDs,
		ExternalIds:  externalIDs,
		DisplayNames: displayNames,
	}); err != nil {
		return nil, err
	}
	return ids, nil
}

// PerformSeed resets all seeded tables and fills the database with seeded
// data according to the given parameters.
func PerformSeed(ctx context.Context, db *sqlc.Queries, params Parameters) error {
	slog.Info("Starting database reset")
	if err := db.SeedTruncateEverythingCascade(ctx); err != nil {
		slog.Error("Failed to reset tables")
		return util.StatusWrap(err, "Failed to reset tables")
	}
	slog.Info("Reset complete")
	slog.Info("Starting database seed")
	startTime := time.Now()
	if err := seed(ctx, db, params); err != nil {
		return util.StatusWrap(err, "Failed to seed database")
	}
	if err := advanceSequences(ctx, db); err != nil {
		return util.StatusWrap(err, "Failed to advance identity sequences")
	}
	slog.Info(fmt.Sprintf("Seed complete. It took %s", time.Since(startTime)))
	return nil
}
