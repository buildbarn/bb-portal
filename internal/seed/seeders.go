package seed

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/seed/sqlc"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_seed"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	jmespath_proto "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
)

// Parameters describes the shape of the data to seed.
type Parameters struct {
	instances                   int32
	users                       int32
	invocationsPerInstance      int32
	targetsPerInvocation        int32
	from                        time.Time
	timeSpan                    time.Duration
	invocationMetadataExtractor *jmespath_proto.Expression
	buildKey                    string
}

// NewSeedParametersFromConfiguration constructs Parameters from the
// bb_seed configuration.
func NewSeedParametersFromConfiguration(configuration *bb_seed.ApplicationConfiguration) Parameters {
	return Parameters{
		instances:                   configuration.Instances,
		users:                       configuration.Users,
		invocationsPerInstance:      configuration.InvocationsPerInstance,
		targetsPerInvocation:        configuration.TargetsPerInvocation,
		from:                        time.Now().Add(-configuration.TimeSpan.AsDuration()),
		timeSpan:                    configuration.TimeSpan.AsDuration(),
		invocationMetadataExtractor: configuration.InvocationMetadataExtractor,
		buildKey:                    configuration.BuildKey,
	}
}

func getTagKeys(key string, metadataExtractor *jmespath_proto.Expression) ([]string, error) {
	resSlice := make([]string, 0)

	if metadataExtractor == nil {
		return resSlice, nil
	}

	compiledMetadataExtractor := jmespath.MustCompile(metadataExtractor.Expression)
	// Evauluate the expression against nil since we only need the keys
	keyMap, err := compiledMetadataExtractor.Search(nil)
	if err != nil {
		return nil, err
	}

	resMap, ok := keyMap.(map[string]any)
	if !ok {
		return resSlice, nil
	}

	if tags, ok := resMap[key].(map[string]any); ok {
		for key := range tags {
			resSlice = append(resSlice, key)
		}
	}
	return resSlice, nil
}

func getSourceControlCount(metadataExtractor *jmespath_proto.Expression) (int, error) {
	if metadataExtractor == nil {
		return 0, nil
	}

	compiledMetadataExtractor := jmespath.MustCompile(metadataExtractor.Expression)
	// Evauluate the expression against nil since we only need the keys
	keyMap, err := compiledMetadataExtractor.Search(nil)
	if err != nil {
		return 0, err
	}

	resMap, ok := keyMap.(map[string]any)
	if !ok {
		return 0, nil
	}

	if sourceControls, ok := resMap["sourceControls"].([]any); ok {
		return len(sourceControls), nil
	}

	return 0, nil
}

// seededTarget is a target whose row has been seeded, carrying only the
// information needed when seeding invocation targets.
type seededTarget struct {
	id         int64
	targetKind string
}

type instanceSeeder struct {
	instanceName       sqlc.InstanceName
	userIDs            []int64
	buildTagKeys       []string
	invocationTagKeys  []string
	sourceControlCount int
	db                 *sqlc.Queries
	buildKey           string
	invocationCount    int32
	targetCount        int32
	from               time.Time
	duration           time.Duration
	random             *rand.Rand

	invocationIDAllocator         *instanceIDAllocator
	targetIDAllocator             *instanceIDAllocator
	buildIDAllocator              *instanceIDAllocator
	buildTagIDAllocator           *instanceIDAllocator
	sourceControlIDAllocator      *instanceIDAllocator
	invocationTagIDAllocator      *instanceIDAllocator
	connectionMetadataIDAllocator *instanceIDAllocator
	eventMetadataIDAllocator      *instanceIDAllocator
	invocationTargetIDAllocator   *instanceIDAllocator
	testTargetIDAllocator         *instanceIDAllocator
	testSummaryIDAllocator        *instanceIDAllocator
	testResultIDAllocator         *instanceIDAllocator
}

func newInstanceSeeder(db *sqlc.Queries, instanceName sqlc.InstanceName, userIDs []int64, buildTagKeys, invocationTagKeys []string, sourceControlCount int, params Parameters, idAllocators *idAllocators) *instanceSeeder {
	return &instanceSeeder{
		db:                            db,
		instanceName:                  instanceName,
		userIDs:                       userIDs,
		buildTagKeys:                  buildTagKeys,
		invocationTagKeys:             invocationTagKeys,
		sourceControlCount:            sourceControlCount,
		random:                        rand.New(rand.NewSource(instanceName.ID)),
		buildKey:                      params.buildKey,
		invocationCount:               params.invocationsPerInstance,
		targetCount:                   params.targetsPerInvocation,
		from:                          params.from,
		duration:                      params.timeSpan,
		invocationIDAllocator:         idAllocators.invocations.forInstance(instanceName.ID),
		targetIDAllocator:             idAllocators.targets.forInstance(instanceName.ID),
		buildIDAllocator:              idAllocators.builds.forInstance(instanceName.ID),
		buildTagIDAllocator:           idAllocators.buildTags.forInstance(instanceName.ID),
		sourceControlIDAllocator:      idAllocators.sourceControls.forInstance(instanceName.ID),
		invocationTagIDAllocator:      idAllocators.invocationTags.forInstance(instanceName.ID),
		connectionMetadataIDAllocator: idAllocators.connectionMetadata.forInstance(instanceName.ID),
		eventMetadataIDAllocator:      idAllocators.eventMetadata.forInstance(instanceName.ID),
		invocationTargetIDAllocator:   idAllocators.invocationTargets.forInstance(instanceName.ID),
		testTargetIDAllocator:         idAllocators.testTargets.forInstance(instanceName.ID),
		testSummaryIDAllocator:        idAllocators.testSummaries.forInstance(instanceName.ID),
		testResultIDAllocator:         idAllocators.testResults.forInstance(instanceName.ID),
	}
}

// seedBazelInvocations seeds the invocations of a single instance, together
// with all of the tables that reference them. It returns the identifiers of
// the seeded invocations, in seed order.
func (s *instanceSeeder) seedBazelInvocations(ctx context.Context) ([]int64, error) {
	n := s.invocationCount
	if n == 0 {
		return nil, nil
	}

	type invocationRow struct {
		id                  int64
		invocationID        uuid.UUID
		seedTime            time.Time
		ended               bool
		username            string
		authenticatedUserID int64
	}

	invocationIDs := make([]int64, 0, n)
	invocationBatcher := newBatcher("bazel_invocations", int(n), func(chunk []invocationRow) error {
		params := sqlc.SeedCreateBazelInvocationsParams{
			Ids:                  make([]int64, 0, len(chunk)),
			InvocationIds:        make([]uuid.UUID, 0, len(chunk)),
			CreatedTimestamps:    make([]time.Time, 0, len(chunk)),
			StartedAts:           make([]time.Time, 0, len(chunk)),
			EndedAts:             make([]time.Time, 0, len(chunk)),
			BepCompleteds:        make([]bool, 0, len(chunk)),
			Usernames:            make([]string, 0, len(chunk)),
			Hostnames:            make([]string, 0, len(chunk)),
			ExitCodeNames:        make([]string, 0, len(chunk)),
			ExitCodeCodes:        make([]int32, 0, len(chunk)),
			AuthenticatedUserIds: make([]int64, 0, len(chunk)),
			InstanceNameIds:      make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			endedAt := time.Time{}
			exitCodeCode := int32(-1)
			exitCodeName := ""
			if row.ended {
				endedAt = row.seedTime.Add(2 * time.Minute)
				exitCodeCode = 0
				exitCodeName = "SUCCESS"
			}
			params.Ids = append(params.Ids, row.id)
			params.InvocationIds = append(params.InvocationIds, row.invocationID)
			params.CreatedTimestamps = append(params.CreatedTimestamps, row.seedTime)
			params.StartedAts = append(params.StartedAts, row.seedTime)
			params.EndedAts = append(params.EndedAts, endedAt)
			params.BepCompleteds = append(params.BepCompleteds, row.ended)
			params.Usernames = append(params.Usernames, row.username)
			params.Hostnames = append(params.Hostnames, fmt.Sprintf("seed-host-%d", row.id%1_000_000))
			params.ExitCodeNames = append(params.ExitCodeNames, exitCodeName)
			params.ExitCodeCodes = append(params.ExitCodeCodes, exitCodeCode)
			params.AuthenticatedUserIds = append(params.AuthenticatedUserIds, row.authenticatedUserID)
			params.InstanceNameIds = append(params.InstanceNameIds, s.instanceName.ID)
		}
		_, err := s.db.SeedCreateBazelInvocations(ctx, params)
		return err
	})

	step := s.duration.Nanoseconds() / int64(n)
	for i := range n {
		seedTime := s.from.Add(time.Duration(int64(i) * step))

		row := invocationRow{
			id:           s.invocationIDAllocator.nextID(),
			invocationID: uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("inst-%d;index-%d", s.instanceName.ID, i))),
			seedTime:     seedTime,
			username:     fmt.Sprintf("seed_user-%d", i),
		}

		if len(s.userIDs) > 0 {
			row.username = "anonymous-user"
			row.authenticatedUserID = s.userIDs[s.random.Intn(len(s.userIDs))]
		}

		if i%5 != 0 {
			row.ended = true
		}

		invocationIDs = append(invocationIDs, row.id)
		if err := invocationBatcher.add(row); err != nil {
			return nil, err
		}
	}
	if err := invocationBatcher.result(); err != nil {
		return nil, util.StatusWrap(err, "Could not batch create bazel invocations")
	}

	if s.sourceControlCount > 0 {
		total := s.sourceControlCount * len(invocationIDs)
		sourceControlBatcher := newBatcher("source_controls", total, func(chunk []sourceControlRow) error {
			params := sqlc.SeedCreateSourceControlsParams{
				Ids:                make([]int64, 0, len(chunk)),
				Repos:              make([]string, 0, len(chunk)),
				RepoUrls:           make([]string, 0, len(chunk)),
				Refs:               make([]string, 0, len(chunk)),
				RefUrls:            make([]string, 0, len(chunk)),
				Commits:            make([]string, 0, len(chunk)),
				CommitUrls:         make([]string, 0, len(chunk)),
				BazelInvocationIds: make([]int64, 0, len(chunk)),
			}
			for _, row := range chunk {
				params.Ids = append(params.Ids, row.id)
				params.Repos = append(params.Repos, row.repo)
				params.RepoUrls = append(params.RepoUrls, row.repoURL)
				params.Refs = append(params.Refs, row.ref)
				params.RefUrls = append(params.RefUrls, row.refURL)
				params.Commits = append(params.Commits, row.commit)
				params.CommitUrls = append(params.CommitUrls, row.commitURL)
				params.BazelInvocationIds = append(params.BazelInvocationIds, row.invocationID)
			}
			_, err := s.db.SeedCreateSourceControls(ctx, params)
			return err
		})

		for index, invocationID := range invocationIDs {
			for range s.sourceControlCount {
				if err := sourceControlBatcher.add(sourceControlRow{
					id:           s.sourceControlIDAllocator.nextID(),
					repo:         fmt.Sprintf("repo-%d", index),
					repoURL:      fmt.Sprintf("repoUrl-%d", index),
					ref:          fmt.Sprintf("ref-%d", index),
					refURL:       fmt.Sprintf("refUrl-%d", index),
					commit:       fmt.Sprintf("commit-%d", index),
					commitURL:    fmt.Sprintf("commitUrl-%d", index),
					invocationID: invocationID,
				}); err != nil {
					return nil, err
				}
			}
		}
		if err := sourceControlBatcher.result(); err != nil {
			return nil, util.StatusWrap(err, "Could not batch create source controls")
		}
	}

	if len(s.invocationTagKeys) > 0 {
		total := len(s.invocationTagKeys) * len(invocationIDs)
		invocationTagBatcher := newBatcher("invocation_tags", total, func(chunk []invocationTagRow) error {
			params := sqlc.SeedCreateInvocationTagsParams{
				Ids:                make([]int64, 0, len(chunk)),
				Keys:               make([]string, 0, len(chunk)),
				Values:             make([]string, 0, len(chunk)),
				BazelInvocationIds: make([]int64, 0, len(chunk)),
			}
			for _, row := range chunk {
				params.Ids = append(params.Ids, row.id)
				params.Keys = append(params.Keys, row.key)
				params.Values = append(params.Values, row.value)
				params.BazelInvocationIds = append(params.BazelInvocationIds, row.invocationID)
			}
			_, err := s.db.SeedCreateInvocationTags(ctx, params)
			return err
		})

		for index, invocationID := range invocationIDs {
			for _, key := range s.invocationTagKeys {
				if err := invocationTagBatcher.add(invocationTagRow{
					id:           s.invocationTagIDAllocator.nextID(),
					key:          key,
					value:        fmt.Sprintf("%s-%d", key, index),
					invocationID: invocationID,
				}); err != nil {
					return nil, err
				}
			}
		}
		if err := invocationTagBatcher.result(); err != nil {
			return nil, util.StatusWrap(err, "Could not batch create invocation tags")
		}
	}

	connectionMetadataBatcher := newBatcher("connection_metadata", int(n), func(chunk []connectionMetadataRow) error {
		params := sqlc.SeedCreateConnectionMetadataParams{
			Ids:                   make([]int64, 0, len(chunk)),
			ConnectionLastOpenAts: make([]time.Time, 0, len(chunk)),
			BazelInvocationIds:    make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.ConnectionLastOpenAts = append(params.ConnectionLastOpenAts, row.lastOpenAt)
			params.BazelInvocationIds = append(params.BazelInvocationIds, row.invocationID)
		}
		_, err := s.db.SeedCreateConnectionMetadata(ctx, params)
		return err
	})
	for i, invocationID := range invocationIDs {
		seedTime := s.from.Add(time.Duration(int64(i) * step))
		if err := connectionMetadataBatcher.add(connectionMetadataRow{
			id:           s.connectionMetadataIDAllocator.nextID(),
			lastOpenAt:   seedTime.Add(2 * time.Minute),
			invocationID: invocationID,
		}); err != nil {
			return nil, err
		}
	}
	if err := connectionMetadataBatcher.result(); err != nil {
		return nil, util.StatusWrap(err, "Could not batch create connection metadata")
	}

	eventMetadataBatcher := newBatcher("event_metadata", int(n), func(chunk []eventMetadataRow) error {
		params := sqlc.SeedCreateEventMetadataParams{
			Ids:                make([]int64, 0, len(chunk)),
			Handleds:           make([][]byte, 0, len(chunk)),
			EventReceivedAts:   make([]time.Time, 0, len(chunk)),
			Versions:           make([]int64, 0, len(chunk)),
			BazelInvocationIds: make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.Handleds = append(params.Handleds, row.handled)
			params.EventReceivedAts = append(params.EventReceivedAts, row.eventReceivedAt)
			params.Versions = append(params.Versions, row.version)
			params.BazelInvocationIds = append(params.BazelInvocationIds, row.invocationID)
		}
		_, err := s.db.SeedCreateEventMetadata(ctx, params)
		return err
	})
	for i, invocationID := range invocationIDs {
		seedTime := s.from.Add(time.Duration(int64(i) * step))
		if err := eventMetadataBatcher.add(eventMetadataRow{
			id:              s.eventMetadataIDAllocator.nextID(),
			handled:         make([]byte, 0),
			eventReceivedAt: seedTime.Add(2 * time.Minute),
			version:         1,
			invocationID:    invocationID,
		}); err != nil {
			return nil, err
		}
	}
	if err := eventMetadataBatcher.result(); err != nil {
		return nil, util.StatusWrap(err, "Could not batch create event metadata")
	}

	return invocationIDs, nil
}

// seedBuild seeds a single build, referencing all of the invocations that
// were seeded for this instance.
func (s *instanceSeeder) seedBuild(ctx context.Context, invocationIDs []int64, index int) error {
	if s.buildKey == "" || len(s.buildTagKeys) == 0 || !slices.Contains(s.buildTagKeys, s.buildKey) {
		return nil
	}

	buildUUID := common.CalculateBuildUUID(
		fmt.Sprintf("%s-%d", s.buildKey, index),
		s.instanceName.Name,
	)
	buildID := s.buildIDAllocator.nextID()
	if _, err := s.db.SeedCreateBuilds(ctx, sqlc.SeedCreateBuildsParams{
		Ids:             []int64{buildID},
		BuildUuids:      []uuid.UUID{buildUUID},
		Timestamps:      []time.Time{s.from},
		InstanceNameIds: []int64{s.instanceName.ID},
	}); err != nil {
		return util.StatusWrap(err, "Failed to create build")
	}

	// All invocations of this instance are associated with the build by
	// setting their build_invocations foreign key.
	if _, err := s.db.SeedAssignInvocationsToBuild(ctx, sqlc.SeedAssignInvocationsToBuildParams{
		BuildID:       buildID,
		InvocationIds: invocationIDs,
	}); err != nil {
		return util.StatusWrap(err, "Failed to assign invocations to build")
	}

	buildTagBatcher := newBatcher("build_tags", len(s.buildTagKeys), func(chunk []buildTagRow) error {
		params := sqlc.SeedCreateBuildTagsParams{
			Ids:      make([]int64, 0, len(chunk)),
			Keys:     make([]string, 0, len(chunk)),
			Values:   make([]string, 0, len(chunk)),
			BuildIds: make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.Keys = append(params.Keys, row.key)
			params.Values = append(params.Values, row.value)
			params.BuildIds = append(params.BuildIds, row.buildID)
		}
		_, err := s.db.SeedCreateBuildTags(ctx, params)
		return err
	})

	for _, key := range s.buildTagKeys {
		if err := buildTagBatcher.add(buildTagRow{
			id:      s.buildTagIDAllocator.nextID(),
			key:     key,
			value:   fmt.Sprintf("%s-%d", key, index),
			buildID: buildID,
		}); err != nil {
			return err
		}
	}
	if err := buildTagBatcher.result(); err != nil {
		return util.StatusWrap(err, "Could not batch create build tags")
	}

	return nil
}

// seedTargets seeds the targets of a single instance. It returns the seeded
// targets, in seed order.
func (s *instanceSeeder) seedTargets(ctx context.Context) ([]seededTarget, error) {
	n := s.targetCount
	if n == 0 {
		return nil, nil
	}

	targets := make([]seededTarget, 0, n)
	targetBatcher := newBatcher("targets", int(n), func(chunk []targetRow) error {
		params := sqlc.SeedCreateTargetsParams{
			Ids:             make([]int64, 0, len(chunk)),
			Labels:          make([]string, 0, len(chunk)),
			Aspects:         make([]string, 0, len(chunk)),
			TargetKinds:     make([]string, 0, len(chunk)),
			InstanceNameIds: make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.Labels = append(params.Labels, row.label)
			params.Aspects = append(params.Aspects, "")
			params.TargetKinds = append(params.TargetKinds, row.targetKind)
			params.InstanceNameIds = append(params.InstanceNameIds, s.instanceName.ID)
		}
		_, err := s.db.SeedCreateTargets(ctx, params)
		return err
	})

	for i := range int(n) {
		targetKind := "cc_library rule"
		if s.random.Intn(5) == 0 {
			targetKind = "cc_test rule"
		}
		target := seededTarget{
			id:         s.targetIDAllocator.nextID(),
			targetKind: targetKind,
		}
		targets = append(targets, target)
		if err := targetBatcher.add(targetRow{
			id:         target.id,
			label:      fmt.Sprintf("//seeded/project:target-%d", i),
			targetKind: targetKind,
		}); err != nil {
			return nil, err
		}
	}
	if err := targetBatcher.result(); err != nil {
		return nil, util.StatusWrap(err, "Could not batch create targets")
	}

	return targets, nil
}

// seedInvocationTargets seeds the invocation targets for the given
// invocations and targets, splitting off the test targets, which require
// additional test summaries and test results.
func (s *instanceSeeder) seedInvocationTargets(ctx context.Context, invocationIDs []int64, targets []seededTarget) error {
	var buildTargets []seededTarget
	var testTargets []seededTarget

	for _, target := range targets {
		if target.targetKind == "cc_test rule" {
			testTargets = append(testTargets, target)
		} else {
			buildTargets = append(buildTargets, target)
		}
	}

	if err := s.seedBuildInvocationTargets(ctx, invocationIDs, buildTargets); err != nil {
		return util.StatusWrap(err, "Failed to seed build targets")
	}

	return s.seedTestInvocationTargets(ctx, invocationIDs, testTargets)
}

// seedBuildInvocationTargets handles standard build targets.
func (s *instanceSeeder) seedBuildInvocationTargets(ctx context.Context, invocationIDs []int64, targets []seededTarget) error {
	if len(invocationIDs) == 0 || len(targets) == 0 {
		return nil
	}

	total := len(invocationIDs) * len(targets)
	batcher := newBatcher("invocation_targets", total, func(chunk []invocationTargetRow) error {
		params := sqlc.SeedCreateInvocationTargetsParams{
			Ids:                make([]int64, 0, len(chunk)),
			Successes:          make([]bool, 0, len(chunk)),
			AbortReasons:       make([]string, 0, len(chunk)),
			BazelInvocationIds: make([]int64, 0, len(chunk)),
			TargetIds:          make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.Successes = append(params.Successes, true)
			params.AbortReasons = append(params.AbortReasons, "NONE")
			params.BazelInvocationIds = append(params.BazelInvocationIds, row.invocationID)
			params.TargetIds = append(params.TargetIds, row.targetID)
		}
		_, err := s.db.SeedCreateInvocationTargets(ctx, params)
		return err
	})

	for _, invocationID := range invocationIDs {
		for _, target := range targets {
			if err := batcher.add(invocationTargetRow{
				id:           s.invocationTargetIDAllocator.nextID(),
				invocationID: invocationID,
				targetID:     target.id,
			}); err != nil {
				return err
			}
		}
	}
	return batcher.result()
}

// seedTestInvocationTargets handles test targets. Every test target gets a
// TestTarget marker row, and every invocation x test target combination gets
// an invocation target, test summary and test result.
func (s *instanceSeeder) seedTestInvocationTargets(ctx context.Context, invocationIDs []int64, targets []seededTarget) error {
	if len(invocationIDs) == 0 || len(targets) == 0 {
		return nil
	}

	testTargetBatcher := newBatcher("test_targets", len(targets), func(chunk []testTargetRow) error {
		params := sqlc.SeedCreateTestTargetsParams{
			Ids:       make([]int64, 0, len(chunk)),
			TargetIds: make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.TargetIds = append(params.TargetIds, row.targetID)
		}
		_, err := s.db.SeedCreateTestTargets(ctx, params)
		return err
	})

	for _, target := range targets {
		if err := testTargetBatcher.add(testTargetRow{
			id:       s.testTargetIDAllocator.nextID(),
			targetID: target.id,
		}); err != nil {
			return err
		}
	}
	if err := testTargetBatcher.result(); err != nil {
		return util.StatusWrap(err, "Failed to seed test target markers")
	}

	// Test summaries and results are seeded in the same order as the
	// invocation targets, meaning that their identifiers can be derived
	// from the invocation target identifiers by simple offsetting.
	total := len(invocationIDs) * len(targets)
	invocationTargetBatcher := newBatcher("test_invocation_targets", total, func(chunk []invocationTargetRow) error {
		params := sqlc.SeedCreateInvocationTargetsParams{
			Ids:                make([]int64, 0, len(chunk)),
			Successes:          make([]bool, 0, len(chunk)),
			AbortReasons:       make([]string, 0, len(chunk)),
			BazelInvocationIds: make([]int64, 0, len(chunk)),
			TargetIds:          make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.Successes = append(params.Successes, true)
			params.AbortReasons = append(params.AbortReasons, "NONE")
			params.BazelInvocationIds = append(params.BazelInvocationIds, row.invocationID)
			params.TargetIds = append(params.TargetIds, row.targetID)
		}
		_, err := s.db.SeedCreateInvocationTargets(ctx, params)
		return err
	})
	summaryBatcher := newBatcher("test_summaries", total, func(chunk []testSummaryRow) error {
		params := sqlc.SeedCreateTestSummariesParams{
			Ids:                   make([]int64, 0, len(chunk)),
			OverallStatuses:       make([]string, 0, len(chunk)),
			TotalRunCounts:        make([]int32, 0, len(chunk)),
			RunCounts:             make([]int32, 0, len(chunk)),
			AttemptCounts:         make([]int32, 0, len(chunk)),
			ShardCounts:           make([]int32, 0, len(chunk)),
			TotalNumCacheds:       make([]int32, 0, len(chunk)),
			TotalRunDurationInMss: make([]int64, 0, len(chunk)),
			InvocationTargetIds:   make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.OverallStatuses = append(params.OverallStatuses, "PASSED")
			params.TotalRunCounts = append(params.TotalRunCounts, 1)
			params.RunCounts = append(params.RunCounts, 1)
			params.AttemptCounts = append(params.AttemptCounts, 1)
			params.ShardCounts = append(params.ShardCounts, 1)
			params.TotalNumCacheds = append(params.TotalNumCacheds, 0)
			params.TotalRunDurationInMss = append(params.TotalRunDurationInMss, 1000)
			params.InvocationTargetIds = append(params.InvocationTargetIds, row.invocationTargetID)
		}
		_, err := s.db.SeedCreateTestSummaries(ctx, params)
		return err
	})
	resultBatcher := newBatcher("test_results", total, func(chunk []testResultRow) error {
		params := sqlc.SeedCreateTestResultsParams{
			Ids:                      make([]int64, 0, len(chunk)),
			Runs:                     make([]int32, 0, len(chunk)),
			Shards:                   make([]int32, 0, len(chunk)),
			Attempts:                 make([]int32, 0, len(chunk)),
			Statuses:                 make([]string, 0, len(chunk)),
			CachedLocallys:           make([]bool, 0, len(chunk)),
			TestAttemptDurationInMss: make([]int64, 0, len(chunk)),
			Strategies:               make([]string, 0, len(chunk)),
			CachedRemotelys:          make([]bool, 0, len(chunk)),
			TestSummaryIds:           make([]int64, 0, len(chunk)),
		}
		for _, row := range chunk {
			params.Ids = append(params.Ids, row.id)
			params.Runs = append(params.Runs, 1)
			params.Shards = append(params.Shards, 1)
			params.Attempts = append(params.Attempts, 1)
			params.Statuses = append(params.Statuses, "PASSED")
			params.CachedLocallys = append(params.CachedLocallys, false)
			params.TestAttemptDurationInMss = append(params.TestAttemptDurationInMss, 1000)
			params.Strategies = append(params.Strategies, "local")
			params.CachedRemotelys = append(params.CachedRemotelys, false)
			params.TestSummaryIds = append(params.TestSummaryIds, row.summaryID)
		}
		_, err := s.db.SeedCreateTestResults(ctx, params)
		return err
	})

	for _, invocationID := range invocationIDs {
		for _, target := range targets {
			invocationTargetID := s.invocationTargetIDAllocator.nextID()
			summaryID := s.testSummaryIDAllocator.nextID()
			if err := invocationTargetBatcher.add(invocationTargetRow{
				id:           invocationTargetID,
				invocationID: invocationID,
				targetID:     target.id,
			}); err != nil {
				return err
			}
			if err := summaryBatcher.add(testSummaryRow{
				id:                 summaryID,
				invocationTargetID: invocationTargetID,
			}); err != nil {
				return err
			}
			if err := resultBatcher.add(testResultRow{
				id:        s.testResultIDAllocator.nextID(),
				summaryID: summaryID,
			}); err != nil {
				return err
			}
		}
	}
	if err := invocationTargetBatcher.result(); err != nil {
		return util.StatusWrap(err, "Failed to seed test invocation targets")
	}
	if err := summaryBatcher.result(); err != nil {
		return util.StatusWrap(err, "Failed to seed test summaries")
	}
	return resultBatcher.result()
}

type sourceControlRow struct {
	id           int64
	repo         string
	repoURL      string
	ref          string
	refURL       string
	commit       string
	commitURL    string
	invocationID int64
}

type invocationTagRow struct {
	id           int64
	key          string
	value        string
	invocationID int64
}

type connectionMetadataRow struct {
	id           int64
	lastOpenAt   time.Time
	invocationID int64
}

type eventMetadataRow struct {
	id              int64
	handled         []byte
	eventReceivedAt time.Time
	version         int64
	invocationID    int64
}

type buildTagRow struct {
	id      int64
	key     string
	value   string
	buildID int64
}

type targetRow struct {
	id         int64
	label      string
	targetKind string
}

type invocationTargetRow struct {
	id           int64
	invocationID int64
	targetID     int64
}

type testTargetRow struct {
	id       int64
	targetID int64
}

type testSummaryRow struct {
	id                 int64
	invocationTargetID int64
}

type testResultRow struct {
	id        int64
	summaryID int64
}
