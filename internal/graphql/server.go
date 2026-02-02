package graphql

import (
	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aereal/otelgqlgen"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"go.opentelemetry.io/otel/trace"
)

// NewGraphqlHandler creates a new GraphQL handler
func NewGraphqlHandler(dbClient *ent.Client, tracerProvider trace.TracerProvider) *handler.Server {
	srv := handler.New(NewSchema(dbClient))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	srv.Use(entgql.Transactioner{TxOpener: dbClient})
	srv.Use(otelgqlgen.New(otelgqlgen.WithTracerProvider(tracerProvider)))
	// A fixed complexity limit for incoming GraphQL queries.
	// See https://gqlgen.com/master/reference/complexity/ for more details.
	srv.Use(extension.FixedComplexityLimit(1000))

	return srv
}
