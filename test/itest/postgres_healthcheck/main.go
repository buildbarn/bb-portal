package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	// Register the pgx database/sql driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: postgres_healthcheck CONNECTION_STRING")
	}
	database, err := sql.Open("pgx", os.Args[1])
	if err != nil {
		return fmt.Errorf("open PostgreSQL connection: %w", err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL is not ready: %w", err)
	}
	return nil
}
