package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"slices"
	"time"

	_ "net/http/pprof"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
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
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/global"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	http_server "github.com/buildbarn/bb-storage/pkg/http/server"
	"github.com/buildbarn/bb-storage/pkg/program"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
)

const (
	readHeaderTimeout = 3 * time.Second
	folderPermission  = 0o750
)

func main() {
	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		if len(os.Args) != 2 {
			return status.Error(codes.InvalidArgument, "Usage: bb_portal bb_portal.jsonnet")
		}
		var configuration bb_portal.ApplicationConfiguration
		if err := util.UnmarshalConfigurationFromFile(os.Args[1], &configuration); err != nil {
			return util.StatusWrapf(err, "Failed to read configuration from %s", os.Args[1])
		}

		prometheusmetrics.RegisterMetrics()

		lifecycleState, grpcClientFactory, err := global.ApplyConfiguration(configuration.Global, dependenciesGroup)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		tracerProvider := otel.GetTracerProvider()
		if tracerProvider == nil || reflect.ValueOf(tracerProvider).IsNil() {
			return status.Error(codes.Internal, "Otel tracer provider is nil")
		}
		router := mux.NewRouter()
		router.Use(otelmux.Middleware("bb-portal-http", otelmux.WithTracerProvider(tracerProvider)))

		err = NewGrpcWebSchedulerService(&configuration, siblingsGroup, dependenciesGroup, grpcClientFactory, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC-Web Scheduler service")
		}
		err = NewGrpcWebBrowserService(&configuration, siblingsGroup, dependenciesGroup, grpcClientFactory, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC-Web Browser service")
		}
		err = newBuildEventStreamService(&configuration, siblingsGroup, dependenciesGroup, grpcClientFactory, router, tracerProvider)
		if err != nil {
			return util.StatusWrap(err, "Failed to create BES service")
		}
		// This must be the last service created for the router, as it will
		// handle all unmatched requests.
		err = newFrontendProxyService(&configuration, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create frontend proxy service")
		}

		http_server.NewServersFromConfigurationAndServe(
			configuration.HttpServers,
			http_server.NewMetricsHandler(allowCorsWrapper(configuration.AllowedOrigins, router), "PortalUI"),
			siblingsGroup,
			grpcClientFactory,
		)

		lifecycleState.MarkReadyAndWait(siblingsGroup)
		return nil
	})
}

func newBuildEventStreamService(
	configuration *bb_portal.ApplicationConfiguration,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *mux.Router,
	tracerProvider trace.TracerProvider,
) error {
	besConfiguration := configuration.BesServiceConfiguration
	if besConfiguration == nil {
		log.Printf("Did not start BuildEventStream service because buildEventStreamConfiguration is not configured")
		return nil
	}

	dialect, connection, err := common.NewSQLConnectionFromConfiguration(besConfiguration.Database, tracerProvider)
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

	// Configure the database cleanup service.
	cleanupConfiguration := besConfiguration.DatabaseCleanupConfiguration
	if cleanupConfiguration == nil {
		return status.Error(codes.InvalidArgument, "No databaseCleanupConfiguration configured for BuildEventStreamService")
	}

	databaseCleanerService, err := dbcleanupservice.NewDbCleanupService(
		dbClient,
		clock.SystemClock,
		cleanupConfiguration,
		tracerProvider,
	)
	if err != nil {
		return util.StatusWrap(err, "Failed to create DatabaseCleanupService")
	}

	databaseCleanerService.StartDbCleanupService(context.Background(), dependenciesGroup)

	instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.InstanceNameAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
	}
	dbAuthService := dbauthservice.NewDbAuthService(dbClient.Ent(), clock.SystemClock, instanceNameAuthorizer, time.Second*5)

	err = addGraphqlHandler(configuration, besConfiguration, dbAuthService, dependenciesGroup, grpcClientFactory, router, dbClient.Ent(), tracerProvider)
	if err != nil {
		return util.StatusWrap(err, "Failed to add GraphQL handler for BuildEventStreamService")
	}

	// Handle log requests.
	logHandler, err := loghandler.NewLogHandler(dbClient.Ent(), dbAuthService, tracerProvider)
	router.Path("/api/v1/invocations/{invocation_id}/log").Methods("GET").Handler(logHandler)

	// Handle BEP file uploads over HTTP.
	if besConfiguration.EnableBepFileUpload {
		bepUploader, err := bepuploader.NewBepUploader(dbClient, configuration, dependenciesGroup, grpcClientFactory, tracerProvider, uuid.NewRandom)
		if err != nil {
			return util.StatusWrap(err, "Failed to create BEP file upload handler")
		}
		router.Path("/api/v1/bep/upload").Methods("POST").Handler(bepUploader)
	}

	// Handle the build event stream gRPC strem.
	buildEventServer, err := bes.NewBuildEventServer(dbClient, configuration, dependenciesGroup, grpcClientFactory, tracerProvider, uuid.NewRandom)
	if err != nil {
		return util.StatusWrap(err, "Failed to create BuildEventServer")
	}
	if err := bb_grpc.NewServersFromConfigurationAndServe(
		besConfiguration.GrpcServers,
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

func addGraphqlHandler(
	configuration *bb_portal.ApplicationConfiguration,
	besConfiguration *bb_portal.BuildEventStreamService,
	dbAuthService *dbauthservice.DbAuthService,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *mux.Router,
	dbClient *ent.Client,
	tracerProvider trace.TracerProvider,
) error {
	srv := graphql.NewGraphqlHandler(dbClient, tracerProvider)

	if configuration.InstanceNameAuthorizer == nil {
		return status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
	}

	switch configuration.InstanceNameAuthorizer.Policy.(type) {
	case *auth_pb.AuthorizerConfiguration_Allow:
		// If the policy is "Allow", we add a auth bypass. Using the
		// DbAuthService would just slow us down.
		srv.AroundOperations(func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
			return next(dbauthservice.NewContextWithDbAuthServiceBypass(ctx))
		})
	default:
		srv.AroundOperations(func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
			return next(dbauthservice.NewContextWithDbAuthService(ctx, dbAuthService))
		})
	}

	router.PathPrefix("/graphql").Handler(srv)
	if besConfiguration.EnableGraphqlPlayground {
		router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))
	}
	return nil
}

func newFrontendProxyService(configuration *bb_portal.ApplicationConfiguration, router *mux.Router) error {
	if configuration.FrontendProxyUrl == "" {
		log.Println("No frontend proxy URL specified, skipping proxying")
		return nil
	}
	remote, err := url.Parse(configuration.FrontendProxyUrl)
	if err != nil {
		return util.StatusWrapf(err, "Failed to parse frontend proxy URL")
	}

	// Return 404 for all API requests not already handled.
	router.PathPrefix("/api/").Handler(router.NotFoundHandler)

	log.Println("Proxying frontend to", remote)
	router.PathPrefix("/").Handler(httputil.NewSingleHostReverseProxy(remote))
	return nil
}

func allowCorsWrapper(allowedOrigins []string, httpHandler http.Handler) http.Handler {
	if allowedOrigins == nil {
		log.Println("No allowed origins specified, CORS disabled")
		return httpHandler
	}
	return cors.New(
		cors.Options{
			AllowOriginFunc: func(origin string) bool {
				return slices.Contains(allowedOrigins, origin) || slices.Contains(allowedOrigins, "*")
			},
			AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-Grpc-Web", "X-Requested-With"},
		},
	).Handler(httpHandler)
}
