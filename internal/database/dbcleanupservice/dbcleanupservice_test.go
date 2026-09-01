package dbcleanupservice_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"

	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/durationpb"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = embedded.NewDatabaseProvider(os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	defer dbProvider.Cleanup()
	m.Run()
}

func getNewDbCleanupService(db database.Client, c clock.Clock, traceProvider trace.TracerProvider) (*dbcleanupservice.DbCleanupService, error) {
	cleanupConfiguration := &bb_portal.Database_CleanupConfiguration{
		CleanupInterval:          durationpb.New(1 * time.Minute),
		InvocationMessageTimeout: durationpb.New(30 * time.Second),
		InvocationRetention:      durationpb.New(30 * time.Minute),
	}
	batcher := dbcleanupservice.NewTimedBatcher(clock.SystemClock, 1*time.Second, 128, 1<<20)
	return dbcleanupservice.NewDbCleanupService(db, c, batcher, cleanupConfiguration, traceProvider)
}
