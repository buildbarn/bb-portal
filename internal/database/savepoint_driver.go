package database

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/google/uuid"
)

// savepointDriver implements an ent driver on top of a transaction that
// uses savepoints to emulate nested transactions.
type savepointDriver struct {
	dialect.Driver
}

// commonSQLDriver is the interface that represents an driver capable of
// running arbitrary raw.
type commonSQLDriver interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// Tx overrides the standard driver's Tx method.
func (d *savepointDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	id := uuid.New().String()

	if err := d.Exec(ctx, fmt.Sprintf("SAVEPOINT %q", id), []any{}, nil); err != nil {
		return nil, err
	}

	return &savepointTx{
		Driver: d.Driver,
		id:     id,
	}, nil
}

// QueryContext passthrough to driver, required because ent performs
// reflection to determine our capabilities.
func (d *savepointDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if drv, ok := d.Driver.(commonSQLDriver); ok {
		return drv.QueryContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("Driver.QueryContext is not supported")
}

// ExecContext passthrough to driver, required because ent performs
// reflection to determine our capabilities.
func (d *savepointDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if drv, ok := d.Driver.(commonSQLDriver); ok {
		return drv.ExecContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("Driver.ExecContext is not supported")
}

type savepointTx struct {
	dialect.Driver
	id string
}

// Commit implements commit by simply releasing the parent savepoint
// (i.e. deferring to the parent transaction).
func (t *savepointTx) Commit() error {
	return t.Exec(context.Background(), fmt.Sprintf("RELEASE SAVEPOINT %q", t.id), []any{}, nil)
}

// Rollback implements rollback by doing rollback to the savepoint,
// making the transaction be in it's equivalent earlier state.
func (t *savepointTx) Rollback() error {
	return t.Exec(context.Background(), fmt.Sprintf("ROLLBACK TO SAVEPOINT %q", t.id), []any{}, nil)
}

// QueryContext passthrough to driver, required because ent performs
// reflection to determine our capabilities.
func (t *savepointTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if drv, ok := t.Driver.(commonSQLDriver); ok {
		return drv.QueryContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("Driver.QueryContext is not supported")
}

// ExecContext passthrough to driver, required because ent performs
// reflection to determine our capabilities.
func (t *savepointTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if drv, ok := t.Driver.(commonSQLDriver); ok {
		return drv.ExecContext(ctx, query, args...)
	}
	return nil, fmt.Errorf("Driver.ExecContext is not supported")
}
