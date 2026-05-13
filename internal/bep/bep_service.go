package bep

import (
	"context"
	"time"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-portal/internal/api/http/bepuploader"
	"github.com/buildbarn/bb-portal/internal/api/http/loghandler"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/buildbarn/bb-portal/internal/graphql"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/clock"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewBuildEventProtocolService creates a new service that in turn creates
// everything needed in the backend to handle the Build Event Stream, accept
// manual upload of BEP files and serve the Graphql API
func NewBuildEventProtocolService(
	configuration *bb_portal.BuildEventStreamService,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	instanceNameAuthorizer auth.Authorizer,
	router *mux.Router,
	tracerProvider trace.TracerProvider,
) error {
	dialect, connection, err := common.NewSQLConnectionFromConfiguration(configuration.Database, tracerProvider)
	if err != nil {
		return util.StatusWrap(err, "Failed to connect to database for BuildEventStreamService")
	}

	dbClient, err := database.New(dialect, connection)
	if err != nil {
		return util.StatusWrap(err, "Failed to create database client from connection")
	}

	// Attempt to migrate towards ents model.
	if err = dbClient.Ent().Schema.Create(context.Background(), migrate.WithDropIndex(true)); err != nil {
		return util.StatusWrap(err, "Could not automatically migrate to desired schema")
	}

	prometheusmetrics.SyncMetrics(dbClient.Ent())

	// Configure the database cleanup service.
	if configuration.DatabaseCleanupConfiguration == nil {
		return status.Error(codes.InvalidArgument, "No databaseCleanupConfiguration configured for BuildEventStreamService")
	}
	databaseCleanupService, err := dbcleanupservice.NewDbCleanupService(
		dbClient,
		clock.SystemClock,
		configuration.DatabaseCleanupConfiguration,
		tracerProvider,
	)
	if err != nil {
		return util.StatusWrap(err, "Failed to create DatabaseCleanupService")
	}
	databaseCleanupService.StartDbCleanupService(context.Background(), dependenciesGroup)

	dbAuthService := dbauthservice.NewDbAuthService(dbClient.Ent(), clock.SystemClock, instanceNameAuthorizer, time.Second*5)

	// Handle Graphql requests.
	srv := graphql.NewGraphqlHandler(dbClient, tracerProvider)
	srv.AroundOperations(func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
		return next(dbauthservice.NewContextWithDbAuthService(ctx, dbAuthService))
	})
	router.PathPrefix("/graphql").Handler(srv)
	if configuration.EnableGraphqlPlayground {
		router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))
	}

	// Handle log requests.
	logHandler, err := loghandler.NewLogHandler(dbClient.Ent(), dbAuthService, tracerProvider)
	router.Path("/api/v1/invocations/{invocation_id}/log").Methods("GET").Handler(logHandler)

	// Handle BEP file uploads over HTTP.
	if configuration.EnableBepFileUpload {
		bepUploader, err := bepuploader.NewBepUploader(dbClient, configuration, instanceNameAuthorizer, dependenciesGroup, grpcClientFactory, tracerProvider)
		if err != nil {
			return util.StatusWrap(err, "Failed to create BEP file upload handler")
		}
		router.Path("/api/v1/bep/upload").Methods("POST").Handler(bepUploader)
	}

	// Handle the Build Event gRPC Stream.
	buildEventServer, err := bes.NewBuildEventServer(dbClient, configuration, instanceNameAuthorizer, dependenciesGroup, grpcClientFactory, tracerProvider)
	if err != nil {
		return util.StatusWrap(err, "Failed to create BuildEventServer")
	}
	if err := bb_grpc.NewServersFromConfigurationAndServe(
		configuration.GrpcServers,
		func(s go_grpc.ServiceRegistrar) {
			build.RegisterPublishBuildEventServer(s.(*go_grpc.Server), buildEventServer)
		},
		siblingsGroup,
		grpcClientFactory,
	); err != nil {
		return util.StatusWrap(err, "gRPC server failure")
	}
	return nil
}
