package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	// Register the pgx database/sql driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const requestTimeout = 2 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 5 {
		return fmt.Errorf("usage: stack_healthcheck POSTGRESQL_CONNECTION_STRING CAS_ADDRESS SCHEDULER_ADDRESS JAEGER_URL")
	}
	if err := checkPostgres(os.Args[1]); err != nil {
		return err
	}
	if err := checkCAS(os.Args[2]); err != nil {
		return err
	}
	if err := checkWorker(os.Args[3]); err != nil {
		return err
	}
	return checkJaeger(os.Args[4])
}

func checkPostgres(connectionString string) error {
	database, err := sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("open PostgreSQL connection: %w", err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL is not ready: %w", err)
	}
	return nil
}

func checkCAS(address string) error {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create CAS client: %w", err)
	}
	defer connection.Close()

	// A batch of absent digests makes the sharding frontend contact both
	// storage backends with overwhelming probability, without writing data.
	digests := make([]*remoteexecution.Digest, 32)
	for i := range digests {
		hash := sha256.Sum256([]byte(fmt.Sprintf("bb-portal-itest-health-check-%d", i)))
		digests[i] = &remoteexecution.Digest{
			Hash:      hex.EncodeToString(hash[:]),
			SizeBytes: 1,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err = remoteexecution.NewContentAddressableStorageClient(connection).FindMissingBlobs(ctx, &remoteexecution.FindMissingBlobsRequest{
		InstanceName: "hardlinking",
		BlobDigests:  digests,
	})
	if err != nil {
		return fmt.Errorf("Buildbarn CAS is not ready: %w", err)
	}
	return nil
}

func checkWorker(address string) error {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create scheduler client: %w", err)
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	response, err := buildqueuestate.NewBuildQueueStateClient(connection).ListPlatformQueues(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("Buildbarn scheduler is not ready: %w", err)
	}
	for _, platformQueue := range response.GetPlatformQueues() {
		if platformQueue.GetName().GetInstanceNamePrefix() == "hardlinking" {
			for _, sizeClassQueue := range platformQueue.GetSizeClassQueues() {
				if sizeClassQueue.GetWorkersCount() > 0 {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("Buildbarn worker is not registered yet")
}

func checkJaeger(address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return fmt.Errorf("create Jaeger health check request: %w", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("Jaeger is not ready: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Jaeger is not ready: %s", response.Status)
	}
	return nil
}
