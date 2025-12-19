package common

import (
	"database/sql"
	"fmt"

	// Register the pgx driver.
	_ "github.com/jackc/pgx/v5/stdlib"

	"entgo.io/ent/dialect"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/util"

	"github.com/google/uuid"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RollbackAndWrapError attempts to roll back the provided transaction.
// If that fails, the rollback error is combined with the original error
// into a single error value.
func RollbackAndWrapError(tx database.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		if err == nil {
			return util.StatusWrap(rerr, "Failed to rollback transaction")
		}
		return util.StatusFromMultiple([]error{
			util.StatusWrap(rerr, "Failed to rollback transaction"),
			err,
		})
	}
	return err
}

// CalculateBuildUUID calculates a UUID for a build, based on the build URL
// and instance name.
func CalculateBuildUUID(buildURL, instanceName string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("instanceName: %s, buildUrl: %s", instanceName, buildURL)))
}

// NewSQLConnectionFromConfiguration creates an otel decorated sql
// connection.
func NewSQLConnectionFromConfiguration(dbConfig *bb_portal.Database, tracerProvider trace.TracerProvider) (string, *sql.DB, error) {
	switch dbConfig := dbConfig.Source.(type) {
	case *bb_portal.Database_Postgres:
		if dbConfig.Postgres.ConnectionString == "" {
			return "", nil, status.Error(codes.InvalidArgument, "Empty connection string for postgres database")
		}
		db, err := otelsql.Open(
			"pgx",
			dbConfig.Postgres.ConnectionString,
			otelsql.WithTracerProvider(tracerProvider),
			otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL),
		)
		if err != nil {
			return "", nil, util.StatusWrap(err, "Failed to open postgres database")
		}
		return dialect.Postgres, db, nil
	default:
		return "", nil, status.Error(codes.InvalidArgument, "Missing database configuration")
	}
}
