package graphql_test

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	gqlgengraphql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/enttest"
	"github.com/buildbarn/bb-portal/internal/graphql"
)

var (
	update = flag.Bool("update-golden", false, "update golden (.out.json) files")
)

type setupFuncEnt func(ctx context.Context, client *ent.Client)

const (
	queryDir             = "testdata/queries"
	consumerContractFile = "../../frontend/src/graphql/__generated__/persisted-documents.json"
	snapshotDir          = "testdata/snapshots"

	readonlyFixtureEntDatasource = "./testdata/snapshot.db"
)

type mockServer struct {
	Server *httptest.Server
	URL    string
	ctx    context.Context
	client *ent.Client
}

func newMockServer(t *testing.T, entDataSource string) *mockServer {
	ctx := context.Background()
	if entDataSource == "" {
		entDataSource = readonlyFixtureEntDatasource
	}

	require.FileExists(t, entDataSource)
	entDataSource = "file:" + entDataSource + "?mode=ro&_fk=1"
	client := enttest.Open(t, "sqlite3", entDataSource)

	graphQLHandler := handler.NewDefaultServer(graphql.NewSchema(client))

	// Limit concurrency to 1. This prevents GraphQL resolver from creating additional database connections
	// that don't know about the in-memory db schema, thus resulting in an error.
	mu := sync.Mutex{}
	graphQLHandler.AroundFields(func(ctx context.Context, next gqlgengraphql.Resolver) (res interface{}, err error) {
		mu.Lock()
		defer mu.Unlock()
		return next(ctx)
	})

	mux := http.NewServeMux()
	mux.Handle("/graphql", graphQLHandler)
	server := httptest.NewServer(mux)

	srv := mockServer{
		Server: server,
		URL:    server.URL,
		ctx:    ctx,
		//ctrl:   ctrl,
		client: client,
	}
	t.Cleanup(func() {
		require.NoError(t, client.Close())
		srv.close()
	})
	return &srv
}

func (srv *mockServer) close() {
	srv.Server.Close()
}
