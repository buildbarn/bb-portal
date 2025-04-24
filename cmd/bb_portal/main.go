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
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	go_grpc "google.golang.org/grpc"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/api"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-portal/internal/api/servefiles"
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
	bepFolder         = flag.String("bep-folder", "./bep-files/", "Folder to watch for new BEP files")
	blobArchiveFolder = flag.String("blob-archive-folder", "./blob-archive/",
		"Folder where blobs (log outputs, stdout, stderr, undeclared test outputs) referenced from failures are archived")
)

func main() {
	go func() {
		pprof := http.ListenAndServe("localhost:8083", nil)
		slog.Info(pprof.Error())
	}()

	go func() {
		// initialize prometheus metrics
		prometheusmetrics.RegisterMetrics()
		slog.Info("Starting metrics server on :8112")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8112", nil)
	}()

	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		flag.Parse()

		if *configFile == "" {
			flag.Usage()
		}

		var configuration bb_portal.ApplicationConfiguration
		if err := util.UnmarshalConfigurationFromFile(*configFile, &configuration); err != nil {
			return util.StatusWrapf(err, "Failed to read configuration from %s", os.Args[1])
		}

		lifecycleState, grpcClientFactory, err := global.ApplyConfiguration(configuration.Global)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		var dbClient *ent.Client

		if *dsDriver == "pgx" {
			db, err := sql.Open("pgx", *dsURL)
			if err != nil {
				fatal("Failed to open pgx database", "err", err)
			}
			drv := entsql.OpenDB(dialect.Postgres, db)
			dbClient = ent.NewClient(ent.Driver(drv))
		} else {
			dbClient, err = ent.Open(
				*dsDriver,
				*dsURL,
			)
		}

		if err != nil {
			return util.StatusWrapf(err, "Failed to open ent client")
		}

		if *dsDriver == "pgx" {
			if err = dbClient.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true), migrate.WithDropIndex(true)); err != nil {
				return util.StatusWrapf(err, "Failed to run schema migration")
			}
		} else {
			if err = dbClient.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
				return util.StatusWrapf(err, "Failed to run schema migration")
			}
		}

		blobArchiver := processing.NewBlobMultiArchiver()
		configureBlobArchiving(blobArchiver, *blobArchiveFolder)

		// Create new watcher.
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return util.StatusWrapf(err, "Failed to create fsnotify.Watcher")
		}
		defer watcher.Close()
		runWatcher(watcher, dbClient, *bepFolder, blobArchiver)

		serveFileService := servefiles.NewFileServerServiceFromConfiguration(dependenciesGroup, &configuration, grpcClientFactory)

		router := mux.NewRouter()
		c := cors.New(
			cors.Options{
				AllowOriginFunc: func(origin string) bool {
					return slices.Contains(configuration.AllowedOrigins, origin) || slices.Contains(configuration.AllowedOrigins, "*")
				},
				AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
				AllowedHeaders: []string{"Authorization", "Content-Type"},
			},
		)
		newPortalService(&configuration, blobArchiver, dbClient, serveFileService, router)
		bb_http.NewServersFromConfigurationAndServe(
			configuration.HttpServers,
			bb_http.NewMetricsHandler(c.Handler(router), "PortalUI"),
			siblingsGroup,
		)

		if err := bb_grpc.NewServersFromConfigurationAndServe(
			configuration.GrpcServers,
			func(s go_grpc.ServiceRegistrar) {
				build.RegisterPublishBuildEventServer(s.(*go_grpc.Server), bes.NewBuildEventServer(dbClient, blobArchiver))
			},
			siblingsGroup,
		); err != nil {
			return util.StatusWrap(err, "gRPC server failure")
		}

		StartGrpcWebProxyServer(&configuration, siblingsGroup, grpcClientFactory)

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

func runWatcher(watcher *fsnotify.Watcher, client *ent.Client, bepFolder string, blobArchiver processing.BlobMultiArchiver) {
	ctx := context.Background()
	worker := processing.New(client, blobArchiver)
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				slog.Info("Received an event", "event", event)
				if event.Has(fsnotify.Write) {
					slog.Info("modified file", "name", event.Name)
					if _, err := worker.ProcessFile(ctx, event.Name); err != nil {
						slog.Error("Failed to process file", "file", event.Name, "err", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("Received an error from fsnotify", "err", err)
			}
		}
	}()

	// Add a path.
	err := os.MkdirAll(bepFolder, folderPermission)
	if err != nil {
		fatal("failed to create BEP folder", "folder", bepFolder, "err", err)
	}
	err = watcher.Add(bepFolder)
	if err != nil {
		fatal("watched register BEP folder with fsnotify.Watcher", "folder", bepFolder, "err", err)
	}
}

func fatal(msg string, args ...any) {
	// Workaround: No slog.Fatal.
	slog.Error(msg, args...)
	os.Exit(1)
}

func newPortalService(configuration *bb_portal.ApplicationConfiguration, archiver processing.BlobMultiArchiver, dbClient *ent.Client, serveFilesService *servefiles.FileServerService, router *mux.Router) {
	srv := handler.NewDefaultServer(graphql.NewSchema(dbClient))
	srv.Use(entgql.Transactioner{TxOpener: dbClient})

	router.PathPrefix("/graphql").Handler(srv)
	router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))
	router.Handle("/api/v1/bep/upload", api.NewBEPUploadHandler(dbClient, archiver)).Methods("POST")
	if serveFilesService != nil {
		router.HandleFunc("/api/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/file/{hash}-{sizeBytes}/{name}", serveFilesService.HandleFile).Methods("GET")
		router.HandleFunc("/api/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/command/{hash}-{sizeBytes}/", serveFilesService.HandleCommand).Methods("GET")
		router.HandleFunc("/api/servefile/{instanceName:(?:.*?/)?}blobs/{digestFunction}/directory/{hash}-{sizeBytes}/", serveFilesService.HandleDirectory).Methods("GET")
	}
	if configuration.FrontendProxyUrl != "" {
		router.PathPrefix("/").Handler(frontendServer(configuration.FrontendProxyUrl))
	}
}

func frontendServer(proxyURL string) http.Handler {
	remote, err := url.Parse(proxyURL)
	if err != nil {
		log.Fatalf("Could not parse proxy URL: %v", err)
	}

	log.Println("Proxying frontend to", remote)

	return httputil.NewSingleHostReverseProxy(remote)
}
