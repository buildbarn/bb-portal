package dbauthservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	"go.uber.org/mock/gomock"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	"github.com/stretchr/testify/require"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = embedded.NewDatabaseProvider(os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	defer dbProvider.Cleanup()
	os.Exit(m.Run())
}

func TestGetInstanceNames(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)
	clock := mock.NewMockClock(ctrl)

	authorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return true })

	t.Run("NoInstanceNames", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider).Ent()
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		clock.EXPECT().Now().Return(time.Unix(1000, 0))
		got := dbAuthService.GetInstanceNames(ctx)
		require.Len(t, got, 0)
	})

	t.Run("SkipInvalidNames", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider).Ent()
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
		db := testutils.SetupTestDB(t, dbProvider).Ent()
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
		db := testutils.SetupTestDB(t, dbProvider).Ent()
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
		db := testutils.SetupTestDB(t, dbProvider).Ent()
		dbAuthService := dbauthservice.NewDbAuthService(db, clock, authorizer, 0)

		got := dbAuthService.GetAuthorizedInstanceNames(ctx)
		require.Len(t, got, 0)
	})

	t.Run("MultipleInstanceNames", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider).Ent()
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
