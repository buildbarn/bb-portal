package testutils

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
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
