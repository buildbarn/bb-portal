package main

import (
	"context"
	"log/slog"
	"os"
	"reflect"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"

	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/seed"
	"github.com/buildbarn/bb-portal/internal/seed/sqlc"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_seed"
	"github.com/buildbarn/bb-storage/pkg/global"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) (err error) {
		if len(os.Args) != 2 {
			return status.Error(codes.InvalidArgument, "Usage: bb_seed bb_seed.jsonnet")
		}
		var configuration bb_seed.ApplicationConfiguration
		if err := util.UnmarshalConfigurationFromFile(os.Args[1], &configuration); err != nil {
			return util.StatusWrapf(err, "Failed to read configuration from %s", os.Args[1])
		}

		_, _, err = global.ApplyConfiguration(configuration.Global, dependenciesGroup)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		tracerProvider := otel.GetTracerProvider()
		if tracerProvider == nil || reflect.ValueOf(tracerProvider).IsNil() {
			return status.Error(codes.Internal, "Otel tracer provider is nil")
		}

		dialect, connection, err := common.NewSQLConnectionFromConfiguration(configuration.Database, tracerProvider)
		if err != nil {
			return util.StatusWrap(err, "Failed to connect to database for BuildEventStreamService")
		}

		db, err := database.New(dialect, connection)
		if err != nil {
			return util.StatusWrap(err, "Failed to create database client from connection")
		}

		if err = db.Ent().Schema.Create(ctx, migrate.WithDropIndex(true)); err != nil {
			return util.StatusWrap(err, "Could not automatically migrate to desired schema")
		}

		if err = seed.PerformSeed(ctx, sqlc.New(connection), seed.NewSeedParametersFromConfiguration(&configuration)); err != nil {
			slog.Error("There was an error seeding the database")
			return util.StatusWrap(err, "Failed to seed database")
		}
		return nil
	})
}
