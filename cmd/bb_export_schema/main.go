package main

import (
	"context"
	"os"

	entsql "entgo.io/ent/dialect/sql"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		if len(os.Args) != 1 {
			return status.Error(codes.InvalidArgument, "Usage: bb_export_schema")
		}

		dbProvider, err := embedded.NewDatabaseProvider(os.Stderr)
		if err != nil {
			return util.StatusWrap(err, "Failed to create database provider for schema export")
		}
		defer dbProvider.Cleanup()
		db, err := dbProvider.CreateDatabase()
		if err != nil {
			return util.StatusWrap(err, "Failed to create database for schema export")
		}
		defer db.Close()

		driver := entsql.OpenDB("postgres", db)
		client := ent.NewClient(ent.Driver(driver))
		if err = client.Schema.WriteTo(ctx, os.Stdout, migrate.WithDropIndex(true)); err != nil {
			return util.StatusWrap(err, "Failed to output schema")
		}

		return nil
	})
}
