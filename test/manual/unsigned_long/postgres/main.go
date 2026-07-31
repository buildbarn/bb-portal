package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buildbarn/bb-portal/internal/database/embedded"
)

const postgresPassword = "postgres"

func run() (returnErr error) {
	postgresPort := flag.Int("postgres-port", 0, "port for the embedded Postgres server")
	healthPort := flag.Int("health-port", 0, "port for the HTTP health endpoint")
	flag.Parse()

	if *postgresPort == 0 {
		return errors.New("--postgres-port is required")
	}
	if *healthPort == 0 {
		return errors.New("--health-port is required")
	}

	databaseProvider, err := embedded.NewDatabaseProviderAtPort(os.Stderr, *postgresPort, postgresPassword)
	if err != nil {
		return fmt.Errorf("start embedded Postgres: %w", err)
	}
	defer func() {
		returnErr = errors.Join(returnErr, databaseProvider.Cleanup())
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.WriteHeader(http.StatusOK)
	})
	healthServer := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", *healthPort),
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	healthServerErrors := make(chan error, 1)
	go func() {
		healthServerErrors <- healthServer.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
	case err := <-healthServerErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve embedded Postgres health endpoint: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := healthServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down embedded Postgres health endpoint: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
