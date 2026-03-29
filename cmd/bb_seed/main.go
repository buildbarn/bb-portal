package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_seed"
	"github.com/buildbarn/bb-storage/pkg/global"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Needed to avoid cyclic dependencies in ent.
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	// Register the pgx stdlib driver for ent.
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) (err error) {
		if len(os.Args) != 2 {
			return status.Error(codes.InvalidArgument, "Usage: bb_seed bb_seed.jsonnet")
		}
		var configuration bb_seed.ApplicationConfiguration
		if err := util.UnmarshalConfigurationFromFile(os.Args[1], &configuration); err != nil {
			return util.StatusWrapf(err, "Failed to read configuration from %s", os.Args[1])
		}

		_, _, err = global.ApplyConfiguration(configuration.Global, dependenciesGroup)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		tracerProvider := otel.GetTracerProvider()
		if tracerProvider == nil || reflect.ValueOf(tracerProvider).IsNil() {
			return status.Error(codes.Internal, "Otel tracer provider is nil")
		}

		dialect, connection, err := common.NewSQLConnectionFromConfiguration(configuration.Database, tracerProvider)
		if err != nil {
			return util.StatusWrap(err, "Failed to connect to database for BuildEventStreamService")
		}

		ctx = dbauthservice.NewContextWithDbAuthServiceBypass(ctx)

		db, err := database.New(dialect, connection)
		if err != nil {
			return util.StatusWrap(err, "Failed to create database client from connection")
		}

		if err = db.Ent().Schema.Create(ctx, migrate.WithDropIndex(true)); err != nil {
			return util.StatusWrap(err, "Could not automatically migrate to desired schema")
		}

		tracer := tracerProvider.Tracer("github.com/buildbarn/bb-portal/cmd/bb_seed")
		if err = performSeed(ctx, db, tracer, newSeedParametersFromConfiguration(&configuration)); err != nil {
			slog.Error("There was an error seeding the database")
			return util.StatusWrap(err, "Failed to seed database")
		}
		return nil
	})
}

type seedParameters struct {
	instances              int32
	users                  int32
	invocationsPerInstance int32
	targetsPerInvocation   int32
	from                   time.Time
	timeSpan               time.Duration
}

func newSeedParametersFromConfiguration(configuration *bb_seed.ApplicationConfiguration) seedParameters {
	return seedParameters{
		instances:              configuration.Instances,
		users:                  configuration.Users,
		invocationsPerInstance: configuration.InvocationsPerInstance,
		targetsPerInvocation:   configuration.TargetsPerInvocation,
		from:                   time.Now().Add(-configuration.TimeSpan.AsDuration()),
		timeSpan:               configuration.TimeSpan.AsDuration(),
	}
}

type instanceSeeder struct {
	instanceName    *ent.InstanceName
	users           []*ent.AuthenticatedUser
	db              database.Handle
	invocationCount int32
	targetCount     int32
	from            time.Time
	duration        time.Duration
	random          *rand.Rand
	tracer          trace.Tracer
}

func newInstanceSeeder(db database.Handle, instanceName *ent.InstanceName, users []*ent.AuthenticatedUser, tracer trace.Tracer, params seedParameters) *instanceSeeder {
	return &instanceSeeder{
		db:              db,
		tracer:          tracer,
		instanceName:    instanceName,
		users:           users,
		random:          rand.New(rand.NewSource(instanceName.ID)),
		invocationCount: params.invocationsPerInstance,
		targetCount:     params.targetsPerInvocation,
		from:            params.from,
		duration:        params.timeSpan,
	}
}

const batchSize = 500

type batcher[T, U any] struct {
	chunk     []T
	execute   func(chunk []T) ([]U, error)
	output    []U
	lastPrint time.Time
	processed int
	total     int
}

func newBatcher[T, U any](total int, execute func(chunk []T) ([]U, error)) *batcher[T, U] {
	return &batcher[T, U]{
		chunk:     make([]T, 0, batchSize),
		execute:   execute,
		output:    []U{},
		lastPrint: time.Now(),
		processed: 0,
		total:     total,
	}
}

func (b *batcher[T, U]) add(item T) error {
	b.chunk = append(b.chunk, item)
	if len(b.chunk) >= batchSize {
		return b.flush()
	}
	return nil
}

func (b *batcher[T, U]) flush() error {
	if len(b.chunk) > 0 {
		o, err := b.execute(b.chunk)
		if err != nil {
			return err
		}
		b.output = append(b.output, o...)
		b.processed += len(b.chunk)
		b.chunk = make([]T, 0, batchSize)

		if time.Since(b.lastPrint) > 1*time.Second {
			slog.Info(fmt.Sprintf("Processed %d/%d", b.processed, b.total))
			b.lastPrint = time.Now()
		}
	}
	return nil
}

func (b *batcher[T, U]) result() ([]U, error) {
	err := b.flush()
	if err != nil {
		return nil, err
	}
	return b.output, nil
}

func performSeed(ctx context.Context, client database.Handle, tracer trace.Tracer, params seedParameters) error {
	ctx, span := tracer.Start(ctx, "performSeed")
	defer span.End()

	slog.Info("Starting database reset")
	err := resetTables(ctx, client)
	if err != nil {
		slog.Error("Failed to reset tables")
		return util.StatusWrap(err, "Failed to reset tables")
	}
	slog.Info("Reset complete")
	slog.Info("Starting database seed")
	startTime := time.Now()
	err = seed(
		ctx,
		client,
		tracer,
		params,
	)
	if err != nil {
		slog.Error("Failed to reset tables")
		return util.StatusWrap(err, "Failed to seed database")
	}
	slog.Info(fmt.Sprintf("Seed complete. It took %s", time.Since(startTime)))
	return nil
}

func resetTables(ctx context.Context, db database.Handle) error {
	_, err := db.Ent().ExecContext(ctx, `TRUNCATE TABLE instance_names, authenticated_users RESTART IDENTITY CASCADE`)
	if err != nil {
		util.StatusWrap(err, "Failed to truncate tables: ")
	}
	return nil
}

func seed(ctx context.Context, db database.Handle, tracer trace.Tracer, params seedParameters) error {
	instanceNames, err := seedInstanceNames(ctx, db, params.instances)
	if err != nil {
		return util.StatusWrap(err, "Failed to seed instance names")
	}

	users, err := seedUsers(ctx, db, params.users)
	if err != nil {
		return util.StatusWrap(err, "Failed to seed authenticated users")
	}

	var g errgroup.Group
	for _, instance := range instanceNames {
		g.Go(func() error {
			ctx, span := tracer.Start(ctx, "instance")
			seeder := newInstanceSeeder(db, instance, users, tracer, params)
			defer span.End()
			invocations, err := seeder.seedBazelInvocations(ctx)
			if err != nil {
				return util.StatusWrap(err, "Failed to seed bazel invocations")
			}
			targets, err := seeder.seedTargets(ctx)
			if err != nil {
				return util.StatusWrap(err, "Failed to seed targets")
			}
			return seeder.seedInvocationTargets(ctx, invocations, targets)
		})
	}
	return g.Wait()
}

func seedUsers(ctx context.Context, db database.Handle, n int32) ([]*ent.AuthenticatedUser, error) {
	if n == 0 {
		return nil, nil
	}

	e := db.Ent()
	batcher := newBatcher(int(n), func(chunk []*ent.AuthenticatedUserCreate) ([]*ent.AuthenticatedUser, error) {
		return e.AuthenticatedUser.CreateBulk(chunk...).Save(ctx)
	})
	for i := range n {
		user := fmt.Sprintf("seed_user_%d", i)
		err := batcher.add(e.AuthenticatedUser.Create().
			SetUserUUID(uuid.NewSHA1(uuid.NameSpaceURL, []byte(user))).
			SetExternalID(user).
			SetDisplayName(user))
		if err != nil {
			return nil, err
		}
	}

	users, err := batcher.result()
	if err != nil {
		return nil, util.StatusWrap(err, "Could not batch create authenticated users")
	}
	return users, nil
}

func (s *instanceSeeder) seedTargets(ctx context.Context) ([]*ent.Target, error) {
	n := s.targetCount
	if n == 0 {
		return nil, nil
	}
	ctx, span := s.tracer.Start(ctx, "seedTargets")
	defer span.End()
	e := s.db.Ent()
	batcher := newBatcher(int(n), func(chunk []*ent.TargetCreate) ([]*ent.Target, error) {
		return e.Target.CreateBulk(chunk...).Save(ctx)
	})

	indexToTargetKind := func() string {
		if s.random.Intn(5) == 0 {
			return "cc_test rule"
		}
		return "cc_library rule"
	}
	for i := range int(n) {
		err := batcher.add(e.Target.Create().
			SetLabel(fmt.Sprintf("//seeded/project:target-%d", i)).
			SetTargetKind(indexToTargetKind()).
			SetAspect("").
			SetInstanceName(s.instanceName))
		if err != nil {
			return nil, err
		}
	}

	targets, err := batcher.result()
	if err != nil {
		return nil, util.StatusWrap(err, "Could not batch create targets")
	}

	return targets, nil
}

func (s *instanceSeeder) seedInvocationTargets(ctx context.Context, invocations []*ent.BazelInvocation, targets []*ent.Target) error {
	if len(invocations) == 0 || len(targets) == 0 {
		return nil
	}
	ctx, span := s.tracer.Start(ctx, "seedInvocationTargets")
	defer span.End()
	var buildTargets []*ent.Target
	var testTargets []*ent.Target

	for _, tgt := range targets {
		if tgt.TargetKind == "cc_test rule" {
			testTargets = append(testTargets, tgt)
		} else {
			buildTargets = append(buildTargets, tgt)
		}
	}

	if err := s.seedBuildInvocationTargets(ctx, invocations, buildTargets); err != nil {
		return util.StatusWrap(err, "Failed to seed build targets")
	}

	if err := s.seedTestInvocationTargets(ctx, invocations, testTargets); err != nil {
		return util.StatusWrap(err, "Failed to seed test targets")
	}

	return nil
}

// seedBuildInvocationTargets handles standard build targets
func (s *instanceSeeder) seedBuildInvocationTargets(ctx context.Context, invocations []*ent.BazelInvocation, targets []*ent.Target) error {
	if len(invocations) == 0 || len(targets) == 0 {
		return nil
	}
	ctx, span := s.tracer.Start(ctx, "seedBuildInvocationTargets")
	defer span.End()
	e := s.db.Ent()

	batcher := newBatcher(len(invocations)*len(targets), func(chunk []*ent.InvocationTargetCreate) ([]*ent.InvocationTarget, error) {
		return e.InvocationTarget.CreateBulk(chunk...).Save(ctx)
	})
	for _, invocation := range invocations {
		for _, target := range targets {
			err := batcher.add(
				e.InvocationTarget.Create().
					SetBazelInvocation(invocation).
					SetTarget(target).
					SetSuccess(true).
					SetDurationInMs(int64(1500 + (target.ID * 10))).
					SetAbortReason("NONE"),
			)
			if err != nil {
				return err
			}
		}
	}
	_, err := batcher.result()
	if err != nil {
		return err
	}

	return nil
}

func (s *instanceSeeder) seedTestInvocationTargets(ctx context.Context, invocations []*ent.BazelInvocation, targets []*ent.Target) error {
	if len(invocations) == 0 || len(targets) == 0 {
		return nil
	}
	ctx, span := s.tracer.Start(ctx, "seedTestInvocationTargets")
	defer span.End()
	e := s.db.Ent()

	testTargetBatcher := newBatcher(len(targets), func(chunk []*ent.TestTargetCreate) ([]*ent.TestTarget, error) {
		return e.TestTarget.CreateBulk(chunk...).Save(ctx)
	})

	for _, target := range targets {
		err := testTargetBatcher.add(e.TestTarget.Create().SetTarget(target))
		if err != nil {
			return err
		}
	}
	if _, err := testTargetBatcher.result(); err != nil {
		return util.StatusWrap(err, "Failed to seed TestTarget markers")
	}

	invocationTargetBatcher := newBatcher(len(invocations)*len(targets), func(chunk []*ent.InvocationTargetCreate) ([]*ent.InvocationTarget, error) {
		createdInvTargets, err := e.InvocationTarget.CreateBulk(chunk...).Save(ctx)
		if err != nil {
			return nil, err
		}

		var summaryBuilders []*ent.TestSummaryCreate
		for _, invTarget := range createdInvTargets {
			summaryBuilders = append(summaryBuilders, e.TestSummary.Create().
				SetInvocationTarget(invTarget).
				SetOverallStatus("PASSED").
				SetTotalRunCount(1).
				SetRunCount(1).
				SetAttemptCount(1).
				SetShardCount(1).
				SetTotalRunDurationInMs(1000).
				SetTotalNumCached(0))
		}

		createdSummaries, err := e.TestSummary.CreateBulk(summaryBuilders...).Save(ctx)
		if err != nil {
			return nil, err
		}

		var resultBuilders []*ent.TestResultCreate
		for _, summary := range createdSummaries {
			resultBuilders = append(resultBuilders, e.TestResult.Create().
				SetTestSummary(summary).
				SetRun(1).
				SetShard(1).
				SetAttempt(1).
				SetStatus("PASSED").
				SetTestAttemptDurationInMs(1000).
				SetStrategy("local").
				SetCachedLocally(false).
				SetCachedRemotely(false))
		}

		if len(resultBuilders) > 0 {
			if _, err := e.TestResult.CreateBulk(resultBuilders...).Save(ctx); err != nil {
				return nil, err
			}
		}

		return createdInvTargets, nil
	})

	for _, invocation := range invocations {
		for _, target := range targets {
			err := invocationTargetBatcher.add(e.InvocationTarget.Create().
				SetBazelInvocation(invocation).
				SetTarget(target).
				SetSuccess(true).
				SetDurationInMs(int64(1500 + (target.ID * 10))).
				SetAbortReason("NONE"))
			if err != nil {
				return err
			}
		}
	}

	_, err := invocationTargetBatcher.result()
	if err != nil {
		return err
	}

	return nil
}

func seedInstanceNames(ctx context.Context, db database.Handle, n int32) ([]*ent.InstanceName, error) {
	if n == 0 {
		return nil, nil
	}
	e := db.Ent()
	batcher := newBatcher(int(n), func(chunk []*ent.InstanceNameCreate) ([]*ent.InstanceName, error) {
		return e.InstanceName.CreateBulk(chunk...).Save(ctx)
	})

	for i := range n {
		err := batcher.add(
			e.InstanceName.Create().SetName(fmt.Sprintf("instance-%d", i)),
		)
		if err != nil {
			return nil, util.StatusWrap(err, "Could not add instance to batch")
		}
	}

	instances, err := batcher.result()
	if err != nil {
		return nil, util.StatusWrap(err, "Could not batch create instances")
	}
	return instances, nil
}

func (s *instanceSeeder) seedBazelInvocations(ctx context.Context) ([]*ent.BazelInvocation, error) {
	n := s.invocationCount
	if n == 0 {
		return nil, nil
	}
	ctx, span := s.tracer.Start(ctx, "seedBazelInvocations")
	defer span.End()
	e := s.db.Ent()

	invocationBatcher := newBatcher(int(n), func(chunk []*ent.BazelInvocationCreate) ([]*ent.BazelInvocation, error) {
		return e.BazelInvocation.CreateBulk(chunk...).Save(ctx)
	})
	step := s.duration.Nanoseconds() / int64(n)
	for i := range n {
		invocationID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("inst-%d;index-%d", s.instanceName.ID, i)))
		seedTime := s.from.Add(time.Duration(int64(i) * step))

		builder := e.BazelInvocation.Create().
			SetInvocationID(invocationID).
			SetCreatedTimestamp(seedTime).
			SetStartedAt(seedTime).
			SetInstanceName(s.instanceName).
			SetHostname(fmt.Sprintf("seed-host-%d", i)).
			SetUsername(fmt.Sprintf("seed_user-%d", i))

		if len(s.users) > 0 {
			builder.SetAuthenticatedUser(s.users[s.random.Intn(len(s.users))])
		} else {
			builder.SetUsername("anonymous-user")
		}

		if i%5 != 0 {
			builder.SetEndedAt(seedTime.Add(2 * time.Minute)).
				SetExitCodeCode(0).
				SetExitCodeName("SUCCESS").
				SetBepCompleted(true)
		}
		err := invocationBatcher.add(builder)
		if err != nil {
			return nil, err
		}
	}
	invocations, err := invocationBatcher.result()
	if err != nil {
		return nil, util.StatusWrap(err, "Could not batch create bazel invocations")
	}

	connectionMetadataBatcher := newBatcher(int(n), func(chunk []*ent.ConnectionMetadataCreate) ([]*ent.ConnectionMetadata, error) {
		return e.ConnectionMetadata.CreateBulk(chunk...).Save(ctx)
	})
	for i := range n {
		seedTime := s.from.Add(time.Duration(int64(i) * step))
		err := connectionMetadataBatcher.add(e.ConnectionMetadata.Create().
			SetBazelInvocation(invocations[i]).
			SetConnectionLastOpenAt(seedTime.Add(2 * time.Minute)))
		if err != nil {
			return nil, err
		}
	}
	_, err = connectionMetadataBatcher.result()
	if err != nil {
		return nil, util.StatusWrap(err, "Could not batch create connection metadata")
	}

	eventMetadataBatcher := newBatcher(int(n), func(chunk []*ent.EventMetadataCreate) ([]*ent.EventMetadata, error) {
		return e.EventMetadata.CreateBulk(chunk...).Save(ctx)
	})
	for i := range n {
		seedTime := s.from.Add(time.Duration(int64(i) * step))
		err := eventMetadataBatcher.add(e.EventMetadata.Create().
			SetBazelInvocation(invocations[i]).
			SetEventReceivedAt(seedTime.Add(2 * time.Minute)).
			SetHandled(make([]byte, 0)).
			SetVersion(1))
		if err != nil {
			return nil, err
		}
	}
	_, err = eventMetadataBatcher.result()
	if err != nil {
		return nil, util.StatusWrap(err, "Could not batch create event metadata")
	}

	return invocations, nil
}
