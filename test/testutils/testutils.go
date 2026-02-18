package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// SetupTestDB creates a new database used for testing.
func SetupTestDB(t testing.TB, dbProvider *embedded.DatabaseProvider) database.Client {
	conn, err := dbProvider.CreateDatabase()
	require.NoError(t, err)
	conn.SetMaxOpenConns(1)
	db, err := database.New("postgres", conn)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	err = db.Ent().Schema.Create(context.Background())
	require.NoError(t, err)
	return db
}

// CreateInstanceName creates a instance name in the testing database.
func CreateInstanceName(ctx context.Context, t *testing.T, client *ent.Client, name string) *ent.InstanceName {
	instanceName, err := client.InstanceName.Create().
		SetName(name).
		Save(ctx)
	require.NoError(t, err)
	return instanceName
}

// StartCreateInvocation startes a BazelInvocationCreate query with required
// variables to avoid having to repeat them every time. They can be overridden
// by specifying them agin.
func StartCreateInvocation(client *ent.Client, instanceName *ent.InstanceName) *ent.BazelInvocationCreate {
	return client.BazelInvocation.Create().
		SetInvocationID(uuid.New()).
		SetInstanceName(instanceName).
		SetCreatedTimestamp(time.Time{})
}
