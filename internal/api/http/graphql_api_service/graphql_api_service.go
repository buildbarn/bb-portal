package graphqlapiservice

import (
	"context"
	"net/http"
	"time"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buildbarn/bb-portal/internal/api/http/loghandler"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/clock"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StartGraphqlAPIService creates a new service that serves the graphql API and
// associated endpoints.
func StartGraphqlAPIService(
	configuration *bb_portal.GraphqlApiService,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *http.ServeMux,
	tracerProvider trace.TracerProvider,
	dbClient database.Client,
) error {
	if dbClient == nil {
		return status.Error(codes.NotFound, "No Database configured")
	}
	if router == nil {
		return status.Error(codes.NotFound, "Failed to create Graphql endpoint. No http server configured")
	}

	if configuration.ReadAuthorizer == nil {
		return status.Error(codes.NotFound, "No ReadAuthorizer configured")
	}
	readAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.ReadAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create ReadAuthorizer")
	}

	dbAuthService := dbauthservice.NewDbAuthService(dbClient.Ent(), clock.SystemClock, readAuthorizer, time.Second*5)

	srv := graphql.NewGraphqlHandler(dbClient, tracerProvider)
	srv.AroundOperations(func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
		return next(dbauthservice.NewContextWithDbAuthService(ctx, dbAuthService))
	})

	// Handle Graphql requests.
	router.Handle("/graphql", srv)
	router.Handle("/graphql/", srv)

	// Graphql Playground
	router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))

	// Handle log requests.
	logHandler := loghandler.NewLogHandler(dbClient.Ent(), dbAuthService, tracerProvider)
	router.Handle(
		"GET /api/v1/invocations/{invocation_id}/log",
		logHandler,
	)
	return nil
}
