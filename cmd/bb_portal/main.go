package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"slices"
	"time"

	_ "net/http/pprof"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/api"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/processing"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/global"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	bb_http "github.com/buildbarn/bb-storage/pkg/http"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
)

const (
	readHeaderTimeout = 3 * time.Second
	folderPermission  = 0o750
)

var (
	configFile        = flag.String("config-file", "", "bb_portal config file")
	dsDriver          = flag.String("datasource-driver", "sqlite3", "Data source driver to use")
	dsURL             = flag.String("datasource-url", "file:buildportal.db?_journal=WAL&_fk=1", "Data source URL for the DB")
	blobArchiveFolder = flag.String("blob-archive-folder", "./blob-archive/",
		"Folder where blobs (log outputs, stdout, stderr, undeclared test outputs) referenced from failures are archived")
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

		lifecycleState, grpcClientFactory, err := global.ApplyConfiguration(configuration.Global)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		router := mux.NewRouter()

		err = NewGrpcWebSchedulerService(&configuration, siblingsGroup, grpcClientFactory, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC-Web Scheduler service")
		}
		err = NewGrpcWebBrowserService(&configuration, siblingsGroup, dependenciesGroup, grpcClientFactory, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create gRPC-Web Browser service")
		}
		err = newBuildEventStreamService(&configuration, siblingsGroup, grpcClientFactory, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create BES service")
		}
		// This must be the last service created for the router, as it will
		// handle all unmatched requests.
		err = newFrontendProxyService(&configuration, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create frontend proxy service")
		}

		bb_http.NewServersFromConfigurationAndServe(
			configuration.HttpServers,
			bb_http.NewMetricsHandler(allowCorsWrapper(configuration.AllowedOrigins, router), "PortalUI"),
			siblingsGroup,
			grpcClientFactory,
		)

		lifecycleState.MarkReadyAndWait(siblingsGroup)
		return nil
	})
}

func configureBlobArchiving(blobArchiver processing.BlobMultiArchiver, archiveFolder string) {
	err := os.MkdirAll(archiveFolder, folderPermission)
	if err != nil {
		fatal("failed to create blob archive folder", "folder", archiveFolder, "err", err)
	}

	localBlobArchiver := processing.NewLocalFileArchiver(archiveFolder)
	blobArchiver.RegisterArchiver("file", localBlobArchiver)

	noopArchiver := processing.NewNoopArchiver()
	blobArchiver.RegisterArchiver("bytestream", noopArchiver)
}

func fatal(msg string, args ...any) {
	// Workaround: No slog.Fatal.
	slog.Error(msg, args...)
	os.Exit(1)
}

func newBuildEventStreamService(configuration *bb_portal.ApplicationConfiguration, siblingsGroup program.Group, grpcClientFactory bb_grpc.ClientFactory, router *mux.Router) error {
	besConfiguration := configuration.BesServiceConfiguration
	if besConfiguration == nil {
		log.Printf("Did not start BuildEventStream service because buildEventStreamConfiguration is not configured")
		return nil
	}

	var dbClient *ent.Client
	var err error

	switch dbConfig := besConfiguration.Database.Source.(type) {
	case *bb_portal.Database_Sqlite:
		if dbConfig.Sqlite.ConnectionString == "" {
			return status.Error(codes.InvalidArgument, "Empty connection string for sqlite database")
		}
		dbClient, err = ent.Open(
			"sqlite3",
			dbConfig.Sqlite.ConnectionString,
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to open sqlite database")
		}

		if err = dbClient.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
			return util.StatusWrap(err, "Failed to run schema migration")
		}
	case *bb_portal.Database_Postgres:
		if dbConfig.Postgres.ConnectionString == "" {
			return status.Error(codes.InvalidArgument, "Empty connection string for postgres database")
		}
		db, err := sql.Open("pgx", dbConfig.Postgres.ConnectionString)
		if err != nil {
			return util.StatusWrap(err, "Failed to open postgres database")
		}
		drv := entsql.OpenDB(dialect.Postgres, db)
		dbClient = ent.NewClient(ent.Driver(drv))
		if err = dbClient.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true), migrate.WithDropIndex(true)); err != nil {
			return util.StatusWrap(err, "Failed to run schema migration")
		}
	}

	blobArchiveFolder := besConfiguration.BlobArchiveFolder
	if blobArchiveFolder == "" {
		return status.Error(codes.NotFound, "No blobArchiveFolder configured for besServiceConfiguration")
	}
	blobArchiver := processing.NewBlobMultiArchiver()
	configureBlobArchiving(blobArchiver, blobArchiveFolder)

	srv := handler.NewDefaultServer(graphql.NewSchema(dbClient))
	srv.Use(entgql.Transactioner{TxOpener: dbClient})

	router.PathPrefix("/graphql").Handler(srv)
	if besConfiguration.EnableGraphqlPlayground {
		router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))
	}
	if besConfiguration.EnableBepFileUpload {
		router.Handle("/api/v1/bep/upload", api.NewBEPUploadHandler(dbClient, blobArchiver)).Methods("POST")
	}

	if err := bb_grpc.NewServersFromConfigurationAndServe(
		besConfiguration.GrpcServers,
		func(s go_grpc.ServiceRegistrar) {
			build.RegisterPublishBuildEventServer(s.(*go_grpc.Server), bes.NewBuildEventServer(dbClient, blobArchiver))
		},
		siblingsGroup,
		grpcClientFactory,
	); err != nil {
		return util.StatusWrap(err, "gRPC server failure")
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
			AllowedHeaders: []string{"Authorization", "Content-Type", "X-Grpc-Web"},
		},
	).Handler(httpHandler)
}
