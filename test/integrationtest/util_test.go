package integrationtest

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/api/http/bepuploader"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/mock/gomock"
)

var dbProvider *embedded.DatabaseProvider

func TestMain(m *testing.M) {
	var err error
	tmpDir, err := os.MkdirTemp(os.TempDir(), "embedded_db_test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create temp dir: %v\n", err)
		os.Exit(1)
	}

	dbProvider, err = embedded.NewDatabaseProvider(tmpDir, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		dbProvider.Cleanup()
		os.RemoveAll(tmpDir)
	}()

	code := m.Run()
	os.Exit(code)
}

func setupTestDB(t testing.TB) database.Client {
	conn, err := dbProvider.CreateDatabase()
	require.NoError(t, err)
	db, err := database.New("postgres", conn)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	err = db.Ent().Schema.Create(context.Background())
	require.NoError(t, err)
	return db
}

func createMockUUIDGenerator(t *testing.T, uuidString string, times int) util.UUIDGenerator {
	ctrl := gomock.NewController(t)
	uuidGeneratorRecorder := mock.NewMockUUIDGenerator(ctrl)
	uuidGeneratorRecorder.EXPECT().Call().Return(uuid.MustParse(uuidString), nil).Times(times)
	return uuidGeneratorRecorder.Call
}

func setupTestBepUploader(t *testing.T, db database.Client, testCase testCase, uuidGenerator util.UUIDGenerator) *bepuploader.BepUploader {
	config := &bb_portal.ApplicationConfiguration{
		InstanceNameAuthorizer: &auth.AuthorizerConfiguration{
			Policy: &auth.AuthorizerConfiguration_Allow{},
		},
		BesServiceConfiguration: &bb_portal.BuildEventStreamService{
			SaveDataLevel:                testCase.saveDataLevel,
			AuthMetadataKeyConfiguration: testCase.extractors,
		},
	}
	bepUploader, err := bepuploader.NewBepUploader(db, config, nil, nil, noop.NewTracerProvider(), uuidGenerator)
	require.NoError(t, err)
	return bepUploader
}

func startGraphqlHTTPServer(t *testing.T, client *ent.Client) *httptest.Server {
	srv := graphql.NewGraphqlHandler(client, trace.NewNoopTracerProvider())

	// Bypass DB auth service for integration tests.
	srv.AroundOperations(func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
		return next(dbauthservice.NewContextWithDbAuthServiceBypass(ctx))
	})

	server := httptest.NewServer(srv)
	t.Cleanup(func() { server.Close() })
	return server
}

func checkIfErrorMatches(t *testing.T, wantErr, err error) {
	if wantErr != nil {
		require.Error(t, err)
		if !assert.Contains(t, err.Error(), wantErr.Error()) {
			require.NoError(t, err, "unexpected error received")
		}
	} else {
		require.NoError(t, err)
	}
}
