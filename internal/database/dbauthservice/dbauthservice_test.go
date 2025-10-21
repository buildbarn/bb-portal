package dbauthservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t testing.TB) *ent.Client {
	db, err := ent.Open("sqlite3", "file:testDb?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	err = db.Schema.Create(dbauthservice.NewContextWithDbAuthServiceBypass(context.Background()))
	require.NoError(t, err)
	return db
}

func TestGetInstanceNames(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)

	authorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return true })

	t.Run("NoInstanceNames", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		clock.EXPECT().Now().Return(time.Unix(1000, 0))
		got := dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 0)
	})

	t.Run("SkipInvalidNames", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		_, err := db.InstanceName.Create().SetName("validName").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("/invalidName/").Save(ctx)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(time.Unix(1000, 0))
		got := dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 1)
		require.Equal(t, "validName", got[0].String())
	})

	t.Run("MultipleInstanceNames", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		_, err := db.InstanceName.Create().SetName("validName1").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("validName2").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("validName3").Save(ctx)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(time.Unix(1000, 0))
		got := dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 3)
		require.Equal(t, "validName1", got[0].String())
		require.Equal(t, "validName2", got[1].String())
		require.Equal(t, "validName3", got[2].String())
	})

	t.Run("TestCache", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, time.Second*10)

		_, err := db.InstanceName.Create().SetName("validName1").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("validName2").Save(ctx)
		require.NoError(t, err)

		clock.EXPECT().Now().Return(time.Unix(1000, 0))
		got := dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 2)
		require.Equal(t, "validName1", got[0].String())
		require.Equal(t, "validName2", got[1].String())

		clock.EXPECT().Now().Return(time.Unix(1005, 0))
		_, err = db.InstanceName.Create().SetName("validName3").Save(ctx)
		require.NoError(t, err)

		got = dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 2)
		require.Equal(t, "validName1", got[0].String())
		require.Equal(t, "validName2", got[1].String())

		clock.EXPECT().Now().Return(time.Unix(1011, 0))
		got = dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 3)
		require.Equal(t, "validName1", got[0].String())
		require.Equal(t, "validName2", got[1].String())
		require.Equal(t, "validName3", got[2].String())
	})
}

func TestGetAuthorizedInstanceNames(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	clock.EXPECT().Now().AnyTimes().Return(time.Unix(1000, 0))

	authorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("instanceName == 'validName1' || instanceName == 'validName2'"),
	)

	t.Run("NoInstanceNames", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		got := dbAuthService.GetAuthorizedInstanceNames(ctx)
		require.Len(t, got, 0)
	})

	t.Run("MultipleInstanceNames", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		_, err := db.InstanceName.Create().SetName("validName1").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("validName2").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("validName3").Save(ctx)
		require.NoError(t, err)
		_, err = db.InstanceName.Create().SetName("validName4").Save(ctx)
		require.NoError(t, err)

		got := dbAuthService.GetAuthorizedInstanceNames(ctx)
		require.Len(t, got, 2)
		require.Equal(t, "validName1", got[0])
		require.Equal(t, "validName2", got[1])
	})
}

func TestQueryFiltering(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)
	clock.EXPECT().Now().AnyTimes().Return(time.Unix(1000, 0))

	adminCtx := dbauthservice.NewContextWithDbAuthServiceBypass(ctx)
	authorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("instanceName == 'validName1' || instanceName == 'validName2'"),
	)

	t.Run("NoInvocations", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)
		authCtx := dbauthservice.NewContextWithDbAuthService(ctx, dbAuthService)

		got, err := db.BazelInvocation.Query().All(authCtx)
		require.NoError(t, err)
		require.Len(t, got, 0)
		num, err := db.BazelInvocation.Query().Count(authCtx)
		require.NoError(t, err)
		require.Equal(t, 0, num)
	})

	t.Run("MultipleInvocations", func(t *testing.T) {
		db := setupTestDB(t)
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)
		authCtx := dbauthservice.NewContextWithDbAuthService(ctx, dbAuthService)

		instanceName1, err := db.InstanceName.Create().SetName("validName1").Save(adminCtx)
		require.NoError(t, err)
		instanceName2, err := db.InstanceName.Create().SetName("validName2").Save(adminCtx)
		require.NoError(t, err)
		instanceName3, err := db.InstanceName.Create().SetName("validName3").Save(adminCtx)
		require.NoError(t, err)

		bi1, err := db.BazelInvocation.Create().
			SetInstanceNameID(instanceName1.ID).
			SetInvocationID(uuid.New()).
			Save(adminCtx)
		require.NoError(t, err)
		bi2, err := db.BazelInvocation.Create().
			SetInstanceNameID(instanceName2.ID).
			SetInvocationID(uuid.New()).
			Save(adminCtx)
		require.NoError(t, err)
		_, err = db.BazelInvocation.Create().
			SetInstanceNameID(instanceName3.ID).
			SetInvocationID(uuid.New()).
			Save(adminCtx)
		require.NoError(t, err)

		got, err := db.BazelInvocation.Query().All(authCtx)
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.Equal(t, bi1.ID, got[0].ID)
		require.Equal(t, bi2.ID, got[1].ID)

		num, err := db.BazelInvocation.Query().Count(authCtx)
		require.NoError(t, err)
		require.Equal(t, 2, num)
	})
}
