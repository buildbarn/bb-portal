package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// SqlcDBTX adapts database/sql handles to the pgx-flavored DBTX
// interface that sqlc generates. Arguments are passed through to
// database/sql, meaning that PostgreSQL array parameters are encoded
// by the pgx stdlib driver rather than pgx itself.
type SqlcDBTX struct {
	execQueryRow execQueryRow
}

// execQueryRow contains the subset of database/sql operations needed
// by SqlcDBTX. It is implemented by both *sql.DB and *sql.Tx.
type execQueryRow interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// NewSqlcDBTX wraps a database/sql handle so that it can be passed to
// sqlc.New().
func NewSqlcDBTX(db execQueryRow) SqlcDBTX {
	return SqlcDBTX{execQueryRow: db}
}

// Exec executes the query, translating the resulting row count into a
// pgconn.CommandTag.
func (d SqlcDBTX) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	result, err := d.execQueryRow.ExecContext(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	return pgconn.NewCommandTag(strconv.FormatInt(rowsAffected, 10)), nil
}

// Query executes the query and wraps the returned rows.
func (d SqlcDBTX) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	rows, err := d.execQueryRow.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlcRows{rows: rows}, nil
}

// QueryRow executes the query expecting at most one row back.
func (d SqlcDBTX) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	row := d.execQueryRow.QueryRowContext(ctx, query, args...)
	return &sqlcRow{row: row}
}

type sqlcRows struct {
	rows *sql.Rows
}

func (r *sqlcRows) Close() {
	_ = r.rows.Close()
}

func (r *sqlcRows) Err() error {
	return r.rows.Err()
}

func (r *sqlcRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *sqlcRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *sqlcRows) Next() bool {
	return r.rows.Next()
}

func (r *sqlcRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *sqlcRows) Values() ([]any, error) {
	return nil, fmt.Errorf("not supported")
}

func (r *sqlcRows) RawValues() [][]byte {
	return nil
}

func (r *sqlcRows) Conn() *pgx.Conn {
	return nil
}

type sqlcRow struct {
	row *sql.Row
}

func (r *sqlcRow) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}
