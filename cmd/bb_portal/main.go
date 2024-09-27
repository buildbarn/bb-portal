package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/debug"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-sqlite3"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	"github.com/buildbarn/bb-portal/internal/api"
	"github.com/buildbarn/bb-portal/internal/api/grpc"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/cas"
	"github.com/buildbarn/bb-portal/pkg/processing"
)

const (
	readHeaderTimeout = 3 * time.Second
	folderPermission  = 0750
)

var (
	httpBindAddr             = flag.String("bind-http", ":8081", "Bind address for the HTTP server.")
	grpcBindAddr             = flag.String("bind-grpc", ":8082", "Bind address for the gRPC server.")
	enableDebug              = flag.Bool("debug", false, "Enable debugging mode.")
	dsDriver                 = flag.String("datasource-driver", "sqlite3", "Data source driver to use")
	dsURL                    = flag.String("datasource-url", "file:buildportal.db?_journal=WAL&_fk=1", "Data source URL for the DB")
	bepFolder                = flag.String("bep-folder", "./bep-files/", "Folder to watch for new BEP files")
	caFile                   = flag.String("ca-file", "", "Custom CA certificate file")
	credentialsHelperCommand = flag.String("credential_helper", "", "Path to a credential helper. Compatible with Bazel's --credential_helper")
	blobArchiveFolder        = flag.String("blob-archive-folder", "./blob-archive/",
		"Folder where blobs (log outputs, stdout, stderr, undeclared test outputs) referenced from failures are archived")
)

func main() {
	flag.Parse()

	client, err := ent.Open(
		*dsDriver,
		*dsURL,
	)
	if err != nil {
		fatal("opening ent client", "err", err)
	}
	if err = client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
		fatal("running schema migration", "err", err)
	}

	blobArchiver := processing.NewBlobMultiArchiver()
	configureBlobArchiving(blobArchiver, *blobArchiveFolder)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal("failed to create fsnotify.Watcher", "err", err)
	}
	defer watcher.Close()
	runWatcher(watcher, client, *bepFolder, blobArchiver)

	srv := handler.NewDefaultServer(graphql.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})
	if *enableDebug {
		srv.Use(&debug.Tracer{})
	}

	fs := frontendServer()
	http.Handle("/graphql", srv)
	http.Handle("/graphiql",
		playground.Handler("GraphQL Playground", "/graphql"),
	)
	casManager := cas.NewConnectionManager(cas.ManagerParams{
		TLSCACertFile:            *caFile,
		CredentialsHelperCommand: *credentialsHelperCommand,
	})

	http.Handle("/api/v1/blobs/{blobID}/{name}", api.NewBlobHandler(client, casManager))
	http.Handle("POST /api/v1/bep/upload", api.NewBEPUploadHandler(client, blobArchiver))
	http.Handle("/", fs)
	slog.Info("HTTP listening on", "address", *httpBindAddr)

	grpcServer := runGRPCServer(client, *grpcBindAddr, blobArchiver)
	defer grpcServer.GracefulStop()
	slog.Info("gRPC listening on", "address", *grpcBindAddr)

	server := &http.Server{
		Addr:              *httpBindAddr,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	if err = server.ListenAndServe(); err != nil {
		slog.Error("http server terminated", "err", err)
	}
}

func configureBlobArchiving(blobArchiver processing.BlobMultiArchiver, archiveFolder string) {
	err := os.MkdirAll(archiveFolder, folderPermission)
	if err != nil {
		fatal("failed to create blob archive folder", "folder", archiveFolder, "err", err)
	}
	localBlobArchiver := processing.NewLocalFileArchiver(archiveFolder)
	blobArchiver.RegisterArchiver("file", localBlobArchiver)
}

func runGRPCServer(db *ent.Client, bindAddr string, blobArchiver processing.BlobMultiArchiver) *grpc.Server {
	lis, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(db, blobArchiver)
	go func() {
		if err := srv.Serve(lis); err != nil {
			slog.Error("error from gRPC server", "err", err)
		}
	}()
	return srv
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

func frontendServer() http.Handler {
	targetURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:3000",
	}
	return httputil.NewSingleHostReverseProxy(targetURL)
}

func fatal(msg string, args ...any) {
	// Workaround: No slog.Fatal.
	slog.Error(msg, args...)
	os.Exit(1)
}
