package database

import (
	"context"
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

type client struct {
	conn    *sql.DB
	dialect string
	ent     *ent.Client
	sqlc    sqlc.Querier
}

// Client is a handle capable of starting new transactions.
type Client interface {
	Handle
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

// Handle has two different entry points, Ent and Sqlc style. They can
// be mixed and matched as desired.
type Handle interface {
	Ent() *ent.Client
	Sqlc() sqlc.Querier
}

// New creates a new Client for the specified dialect and connection.
func New(dialect string, conn *sql.DB) (Client, error) {
	ret := client{conn: conn, dialect: dialect, sqlc: sqlc.New(conn)}
	driver := entsql.OpenDB(dialect, conn)
	ret.ent = ent.NewClient(ent.Driver(driver))
	return &ret, nil
}

// BeginTx starts a new transaction.
func (client *client) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := newTx(ctx, client.dialect, client.conn, opts)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to start transaction")
	}
	return tx, nil
}

// Ent returns the ent client.
func (client *client) Ent() *ent.Client {
	return client.ent
}

// Sqlc returns the sqlc client.
func (client *client) Sqlc() sqlc.Querier {
	return client.sqlc
}
