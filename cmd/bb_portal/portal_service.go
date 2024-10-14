package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/api"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/cas"
	"github.com/buildbarn/bb-portal/pkg/processing"
	"github.com/gorilla/mux"
)

func newPortalService(archiver processing.BlobMultiArchiver, casManager *cas.ConnectionManager, dbClient *ent.Client, router *mux.Router) {
	srv := handler.NewDefaultServer(graphql.NewSchema(dbClient))
	srv.Use(entgql.Transactioner{TxOpener: dbClient})

	router.PathPrefix("/graphql").Handler(srv)
	router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))
	router.Handle("/api/v1/blobs/{blobID}/{name}", api.NewBlobHandler(dbClient, casManager))
	router.Handle("/api/v1/bep/upload", api.NewBEPUploadHandler(dbClient, archiver)).Methods("POST")
	router.PathPrefix("/").Handler(frontendServer())
}

func frontendServer() http.Handler {
	targetURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:3000",
	}
	return httputil.NewSingleHostReverseProxy(targetURL)
}
