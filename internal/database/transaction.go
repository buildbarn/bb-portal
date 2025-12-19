package database

import (
	"context"
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/sqlc"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// Tx implements a database Handle as well as commit/rollback functions.
// While the Tx interface itself does not support nested invocations the
// resulting *ent.Client may attempt to make one. For this purpose ent
// is supplied with a transaction wrapper that implements nested
// transactions using sql savepoints.
type Tx interface {
	Handle
	Commit() error
	Rollback() error
}

type tx struct {
	tx      *sql.Tx
	dialect string
	ent     *ent.Client
	sqlc    sqlc.Querier
}

func newTx(ctx context.Context, dialect string, conn *sql.DB, opts *sql.TxOptions) (Tx, error) {
	t, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, util.StatusWrap(err, "Unable to start transaction")
	}
	ret := tx{dialect: dialect, tx: t}
	drv := entsql.NewDriver(ret.dialect, entsql.Conn{ExecQuerier: ret.tx})
	txDrv := &savepointDriver{Driver: drv}
	ret.ent = ent.NewClient(ent.Driver(txDrv))
	ret.sqlc = sqlc.New(ret.tx)
	return &ret, nil
}

// Ent returns the ent client.
func (t *tx) Ent() *ent.Client {
	return t.ent
}

// Sqlc returns the sqlc client.
func (t *tx) Sqlc() sqlc.Querier {
	return t.sqlc
}

// Commit commits the underlying transaction.
func (t *tx) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the underlying transaction.
func (t *tx) Rollback() error {
	return t.tx.Rollback()
}
