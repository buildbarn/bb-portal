package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	build "google.golang.org/genproto/googleapis/devtools/build/v1"
	go_grpc "google.golang.org/grpc"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/api"
	"github.com/buildbarn/bb-portal/internal/api/grpc/bes"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/processing"
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
	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		flag.Parse()

		if *configFile == "" {
			flag.Usage()
		}

		var configuration bb_portal.ApplicationConfiguration
		if err := util.UnmarshalConfigurationFromFile(*configFile, &configuration); err != nil {
			return util.StatusWrapf(err, "Failed to read configuration from %s", os.Args[1])
		}

		lifecycleState, _, err := global.ApplyConfiguration(configuration.Global)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		dbClient, err := ent.Open(
			*dsDriver,
			*dsURL,
		)
		if err != nil {
			return util.StatusWrapf(err, "Failed to open ent client")
		}
		if err = dbClient.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
			return util.StatusWrapf(err, "Failed to run schema migration")
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

		router := mux.NewRouter()
		newPortalService(blobArchiver, dbClient, router)
		bb_http.NewServersFromConfigurationAndServe(
			configuration.HttpServers,
			bb_http.NewMetricsHandler(router, "PortalUI"),
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

func newPortalService(archiver processing.BlobMultiArchiver, dbClient *ent.Client, router *mux.Router) {
	srv := handler.NewDefaultServer(graphql.NewSchema(dbClient))
	srv.Use(entgql.Transactioner{TxOpener: dbClient})

	router.PathPrefix("/graphql").Handler(srv)
	router.Handle("/graphiql", playground.Handler("GraphQL Playground", "/graphql"))
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
