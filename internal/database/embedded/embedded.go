package embedded

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	// Register the pgx driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/buildbarn/bb-storage/pkg/util"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DatabaseProvider holds the reference to a postgres instance.
type DatabaseProvider struct {
	postgres       *embeddedpostgres.EmbeddedPostgres
	user, password string
	port           int
	db             *sql.DB
	runtimePath    string
}

// NewDatabaseProvider creates and starts an ephemeral Postgres
// instance that can provide databases on demand.
func NewDatabaseProvider(logger io.Writer) (*DatabaseProvider, error) {
	// Ask the kernel for a free port.
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to find free port")
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	const user, database = "postgres", "postgres"
	password := uuid.New().String()

	runtimePath, err := os.MkdirTemp(os.TempDir(), "embedded_db_runtime")
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create temp dir for embedded postgres runtime")
	}

	config := embeddedpostgres.DefaultConfig().
		Port(uint32(port)).
		Database(database).
		Logger(logger).
		Username(user).
		Password(password).
		RuntimePath(runtimePath)

	postgres := embeddedpostgres.NewDatabase(config)

	if err := postgres.Start(); err != nil {
		return nil, util.StatusWrap(err, "Failed to start embedded postgres")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", user, password, port, database)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		postgres.Stop()
		return nil, util.StatusWrap(err, "Failed to open sql connection")
	}

	return &DatabaseProvider{
		postgres:    postgres,
		user:        user,
		password:    password,
		port:        port,
		db:          db,
		runtimePath: runtimePath,
	}, nil
}

// CreateDatabase creates a new empty database.
func (d *DatabaseProvider) CreateDatabase() (*sql.DB, error) {
	if d.db == nil {
		return nil, status.Error(codes.FailedPrecondition, "Creating database on closed connection")
	}
	database := uuid.New().String()
	if _, err := d.db.Exec(fmt.Sprintf("CREATE DATABASE %q", database)); err != nil {
		return nil, util.StatusWrap(err, "Could not create database")
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", d.user, d.password, d.port, database)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, util.StatusWrap(err, "Could not open connection to database")
	}
	return db, nil
}

// Cleanup closes the primary connection and stops the postgres server.
func (d *DatabaseProvider) Cleanup() error {
	var err1, err2, err3 error
	if d.db != nil {
		err1 = d.db.Close()
	}
	if d.postgres != nil {
		err2 = d.postgres.Stop()
	}
	err3 = os.RemoveAll(d.runtimePath)

	return errors.Join(err1, err2, err3)
}
