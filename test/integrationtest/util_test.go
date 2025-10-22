package integrationtest

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/api/http/bepuploader"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func setupTestDB(t *testing.T) *ent.Client {
	db, err := ent.Open("sqlite3", "file:testDb?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	err = db.Schema.Create(context.Background())
	require.NoError(t, err)
	return db
}

func setupTestBepUploader(t *testing.T, client *ent.Client, traceProvider trace.TracerProvider) *bepuploader.BepUploader {
	config := &bb_portal.ApplicationConfiguration{
		InstanceNameAuthorizer: &auth.AuthorizerConfiguration{
			Policy: &auth.AuthorizerConfiguration_Allow{},
		},
		BesServiceConfiguration: &bb_portal.BuildEventStreamService{
			SaveTargetDataLevel: &bb_portal.BuildEventStreamService_SaveTargetDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTargetDataLevel_Enriched{},
			},
		},
	}
	bepUploader, err := bepuploader.NewBepUploader(client, processing.NewBlobMultiArchiver(), config, nil, nil, traceProvider)
	require.NoError(t, err)
	return bepUploader
}

func startGraphqlHTTPServer(t *testing.T, client *ent.Client) *httptest.Server {
	graphQLHandler := handler.NewDefaultServer(graphql.NewSchema(client))

	server := httptest.NewServer(graphQLHandler)
	t.Cleanup(func() { server.Close() })
	return server
}
