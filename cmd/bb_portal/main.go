package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"

	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"
	apicommon "github.com/buildbarn/bb-portal/internal/api/common"
	graphqlapiservice "github.com/buildbarn/bb-portal/internal/api/http/graphql_api_service"
	"github.com/buildbarn/bb-portal/internal/api/http/zstdmiddleware"
	"github.com/buildbarn/bb-portal/internal/bep"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/buildbarn/bb-portal/pkg/frontend"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/blobstoreservice"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/schedulerservice"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	blobstore_configuration "github.com/buildbarn/bb-storage/pkg/blobstore/configuration"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/global"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	http_server "github.com/buildbarn/bb-storage/pkg/http/server"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	bb_zstd "github.com/buildbarn/bb-storage/pkg/zstd"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupDatabase(
	ctx context.Context,
	configuration *bb_portal.Database,
	tracerProvider trace.TracerProvider,
	dependenciesGroup program.Group,
	zstdPool bb_zstd.Pool,
) (database.Client, *dbcleanupservice.DbCleanupService, error) {
	dialect, connection, err := common.NewSQLConnectionFromConfiguration(configuration, tracerProvider)
	if err != nil {
		return nil, nil, util.StatusWrap(err, "Failed to connect to database")
	}
	return initializeDatabase(ctx, configuration, dialect, connection, tracerProvider, zstdPool)
}

// initializeDatabase performs all database setup that happens after a
// connection has been opened. It is separate from setupDatabase so that
// it can be tested against an embedded database.
func initializeDatabase(
	ctx context.Context,
	configuration *bb_portal.Database,
	dialect string,
	connection *sql.DB,
	tracerProvider trace.TracerProvider,
	zstdPool bb_zstd.Pool,
) (database.Client, *dbcleanupservice.DbCleanupService, error) {
	dbClient, err := database.New(dialect, connection)
	if err != nil {
		return nil, nil, util.StatusWrap(err, "Failed to create database client from connection")
	}

	switch configuration.GetSchemaManagement().GetMode().(type) {
	case nil, *bb_portal.Database_SchemaManagement_Apply, *bb_portal.Database_SchemaManagement_ApplyAndExit:
		// Attempt to migrate towards ents model.
		if err = dbClient.Ent().Schema.Create(ctx, migrate.WithDropIndex(true)); err != nil {
			return nil, nil, util.StatusWrap(err, "Could not automatically migrate to desired schema")
		}
	case *bb_portal.Database_SchemaManagement_Verify:
		if err = verifySchema(ctx, dbClient); err != nil {
			return nil, nil, util.StatusWrap(err, "Failed to verify the installed schema")
		}
	case *bb_portal.Database_SchemaManagement_Skip:
		log.Print("Schema management is set to skip, assuming the desired schema is already installed")
	default:
		return nil, nil, status.Error(codes.Internal, "Schema management mode is not recognized. This is probably a bug.")
	}

	prometheusmetrics.SyncMetrics(dbClient.Ent())

	// Configure the database cleanup service.
	var dbCleanupService *dbcleanupservice.DbCleanupService
	if configuration.CleanupConfiguration != nil {
		dbCleanupService, err = dbcleanupservice.NewDbCleanupService(
			dbClient,
			clock.SystemClock,
			dbcleanupservice.NewTimedBatcher(clock.SystemClock, 1*time.Second, 1<<7, 1<<20),
			configuration.CleanupConfiguration,
			tracerProvider,
			zstdPool,
		)
		if err != nil {
			return nil, nil, util.StatusWrap(err, "Failed to create DatabaseCleanupService")
		}
	}
	return dbClient, dbCleanupService, nil
}

// verifySchema checks that the installed schema has everything the
// desired schema has, for when the schema is managed externally.
func verifySchema(ctx context.Context, dbClient database.Client) error {
	// This plans the changes Schema.Create would make, without executing
	// them. The drop options are deliberately omitted, so that objects
	// the installed schema has in addition are not reported.
	var missingChanges bytes.Buffer
	if err := dbClient.Ent().Schema.WriteTo(ctx, &missingChanges); err != nil {
		return util.StatusWrap(err, "Failed to compare the installed schema with the desired schema")
	}
	if missingChanges.Len() > 0 {
		return status.Errorf(codes.FailedPrecondition, "The installed schema differs from the desired schema. Missing changes:\n%s", missingChanges.String())
	}
	return nil
}

func main() {
	program.RunMain(func(ctx context.Context, siblingsGroup, dependenciesGroup program.Group) error {
		if len(os.Args) != 2 {
			return status.Error(codes.InvalidArgument, "Usage: bb_portal bb_portal.jsonnet")
		}
		var configuration bb_portal.ApplicationConfiguration
		if err := util.UnmarshalConfigurationFromFile(os.Args[1], &configuration); err != nil {
			return util.StatusWrapf(err, "Failed to read configuration from %s", os.Args[1])
		}

		lifecycleState, grpcClientFactory, err := global.ApplyConfiguration(configuration.Global, dependenciesGroup)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		tracerProvider := otel.GetTracerProvider()
		if tracerProvider == nil || reflect.ValueOf(tracerProvider).IsNil() {
			return status.Error(codes.Internal, "Otel tracer provider is nil")
		}

		// Create a process-wide ZSTD compression pool.
		zstdPool := bb_zstd.NewPoolFromConfiguration(configuration.ZstdPool)

		// Setup database, if configured
		var dbClient database.Client
		var dbCleanupService *dbcleanupservice.DbCleanupService
		if configuration.Database != nil {
			if dbClient, dbCleanupService, err = setupDatabase(
				ctx,
				configuration.Database,
				tracerProvider,
				dependenciesGroup,
				zstdPool,
			); err != nil {
				return util.StatusWrap(err, "Failed to setup the database")
			}
			if _, ok := configuration.Database.GetSchemaManagement().GetMode().(*bb_portal.Database_SchemaManagement_ApplyAndExit); ok {
				log.Print("Applied the desired schema, exiting as schema management is set to apply_and_exit")
				return nil
			}
		}

		// Setup blobaccess to CAS, AC, ISCC and FSAC.
		blobAccess, err := createBlobAccess(
			&configuration,
			grpcClientFactory,
			dependenciesGroup,
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create Blob Access")
		}

		// Setup http router, if configured
		var router *http.ServeMux
		if len(configuration.HttpServers) != 0 {
			router = http.NewServeMux()
		}

		// Configure and start the blobstore proxy service
		if configuration.BlobstoreServiceConfiguration != nil {
			if err = blobstoreservice.NewBlobstoreService(
				&configuration,
				zstdPool,
				router,
				blobAccess,
			); err != nil {
				return util.StatusWrap(err, "Failed to create Blobstore service")
			}
		}

		// Configure and start the scheduler service
		if configuration.SchedulerServiceConfiguration != nil {
			if err = schedulerservice.NewSchedulerService(
				configuration.SchedulerServiceConfiguration,
				siblingsGroup,
				dependenciesGroup,
				grpcClientFactory,
				router,
			); err != nil {
				return util.StatusWrap(err, "Failed to create Scheduler service")
			}
		}

		// Configure and start the BES upload service
		if configuration.BesServiceConfiguration != nil {
			if err = bep.NewBuildEventProtocolService(
				configuration.BesServiceConfiguration,
				siblingsGroup,
				dependenciesGroup,
				grpcClientFactory,
				router,
				tracerProvider,
				dbClient,
			); err != nil {
				return util.StatusWrap(err, "Failed to create BES service")
			}
		}

		// Configure and start the graphql API
		if configuration.GraphqlApiServiceConfiguration != nil {
			if err := graphqlapiservice.StartGraphqlAPIService(
				configuration.GraphqlApiServiceConfiguration,
				siblingsGroup,
				dependenciesGroup,
				grpcClientFactory,
				router,
				tracerProvider,
				dbClient,
			); err != nil {
				return util.StatusWrap(err, "Failed to start Graphql API service")
			}
		}

		// Serves info about the current user.
		if router != nil {
			router.HandleFunc(
				"GET /api/v1/whoami",
				whoamiHandler,
			)
		}

		// Serve the frontend. This must be the last service created for the
		// router, as it will handle all unmatched http requests.
		if configuration.FrontendServiceConfiguration != nil {
			if err := frontend.ServeFrontend(
				configuration.FrontendServiceConfiguration,
				router,
			); err != nil {
				return util.StatusWrap(err, "Failed to create frontend proxy service")
			}
		}

		// Setup and start the http server
		if len(configuration.HttpServers) != 0 {
			zstdMiddleware := zstdmiddleware.NewZstdMiddleware(zstdPool)

			var handler http.Handler
			handler = apicommon.ForwardHeadersAsGrpcMetadata(router)
			handler = zstdMiddleware(handler)
			handler = otelhttp.NewHandler(
				handler,
				"bb-portal-http",
				otelhttp.WithTracerProvider(tracerProvider),
			)
			handler = http_server.NewMetricsHandler(handler, "PortalUI")

			http_server.NewServersFromConfigurationAndServe(
				configuration.HttpServers,
				handler,
				siblingsGroup,
				grpcClientFactory,
			)
		}

		// We wait with starting the cleanup until everyting else is up and
		// running. Otherwise a faulty config would shut down the application
		// in the middle of the first cleanup.
		if dbCleanupService != nil {
			dbCleanupService.StartDbCleanupService(ctx, dependenciesGroup)
		}

		lifecycleState.MarkReadyAndWait(siblingsGroup)
		return nil
	})
}

func newStorageBlobAccess[T any](
	dependenciesGroup program.Group,
	configuration *bb_portal.StorageBlobAccessConfiguration,
	creator blobstore_configuration.BlobAccessCreator[T],
	grpcClientFactory bb_grpc.ClientFactory,
) (blobstore_configuration.BlobAccessInfo[T], blobstore.BlobAccess[T], error) {
	info, err := blobstore_configuration.NewBlobAccessFromConfiguration(dependenciesGroup, configuration.Backend, creator)
	if err != nil {
		return blobstore_configuration.BlobAccessInfo[T]{},
			nil,
			util.StatusWrap(err, "Failed to create new blob access from configuration")
	}

	readAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(
		configuration.ReadAuthorizer,
		dependenciesGroup,
		grpcClientFactory,
	)
	if err != nil {
		return blobstore_configuration.BlobAccessInfo[T]{},
			nil,
			util.StatusWrap(err, "Failed to create Get() authorizer")
	}
	writeAuthorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return false })

	return info,
		blobstore.NewAuthorizingBlobAccess(info.BlobAccess, readAuthorizer, writeAuthorizer, readAuthorizer),
		nil
}

func createBlobAccess(
	configuration *bb_portal.ApplicationConfiguration,
	grpcClientFactory bb_grpc.ClientFactory,
	dependenciesGroup program.Group,
) (*blobstoreservice.BlobAccess, error) {
	blobAccess := &blobstoreservice.BlobAccess{}

	// Content Addressable Storage (CAS).
	var contentAddressableStorageInfo *blobstore_configuration.ContentAddressableStorageInfo
	if casConfiguration := configuration.ContentAddressableStorageServer; casConfiguration != nil {
		casInfo, err := blobstore_configuration.NewCASFromConfiguration(
			dependenciesGroup,
			casConfiguration.ContentAddressableStorage,
			grpcClientFactory,
			int(configuration.MaximumMessageSizeBytes),
		)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create Content Addressable Storage")
		}
		contentAddressableStorageInfo = &casInfo

		// The raw storages are handed to specialized blob accesses,
		// such as the Action Cache, which perform authorization at
		// their own API boundaries. The user-facing gRPC servers
		// obtain authorized versions here.
		readAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(
			casConfiguration.ReadAuthorizer,
			dependenciesGroup,
			grpcClientFactory,
		)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create Get() authorizer for Content Addressable Storage")
		}
		putAuthorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return false })
		authorizedChunkStorage := blobstore.NewAuthorizingBlobAccess(casInfo.ChunkStorage, readAuthorizer, putAuthorizer, readAuthorizer)
		authorizedChunkMappingStorage := blobstore.NewAuthorizingBlobAccess(casInfo.ChunkMappingStorage, readAuthorizer, putAuthorizer, readAuthorizer)

		blobAccess.ContentAddressableStorage = &blobstoreservice.ContentAddressableStorage{
			ChunkStorage:              authorizedChunkStorage,
			ChunkMappingStorage:       authorizedChunkMappingStorage,
			ReadInfo:                  casInfo.ReadInfo,
			ReadAuthorizer:            readAuthorizer,
			PutChunkMappingAuthorizer: putAuthorizer,
		}
	}

	// Action Cache (AC).
	if configuration.ActionCache != nil {
		_, authorizedBackend, err := newStorageBlobAccess(
			dependenciesGroup,
			configuration.ActionCache,
			blobstore_configuration.NewACBlobAccessCreator(
				contentAddressableStorageInfo,
				grpcClientFactory,
				int(configuration.MaximumMessageSizeBytes),
			),
			grpcClientFactory,
		)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create Action Cache")
		}
		blobAccess.ActionCache = &authorizedBackend
	}

	// Initial Size Class Cache (ISCC).
	if configuration.InitialSizeClassCache != nil {
		_, authorizedBackend, err := newStorageBlobAccess(
			dependenciesGroup,
			configuration.InitialSizeClassCache,
			blobstore_configuration.NewISCCBlobAccessCreator(
				grpcClientFactory,
				int(configuration.MaximumMessageSizeBytes),
			),
			grpcClientFactory,
		)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create Initial Size Class Cache")
		}
		blobAccess.InitialSizeClassCache = &authorizedBackend
	}

	// File System Access Cache (FSAC).
	if configuration.FileSystemAccessCache != nil {
		_, authorizedBackend, err := newStorageBlobAccess(
			dependenciesGroup,
			configuration.FileSystemAccessCache,
			blobstore_configuration.NewFSACBlobAccessCreator(
				grpcClientFactory,
				int(configuration.MaximumMessageSizeBytes),
			),
			grpcClientFactory,
		)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create File System Access Cache")
		}
		blobAccess.FileSystemAccessCache = &authorizedBackend
	}
	return blobAccess, nil
}

// whoamiHandler is a http handler that returns info about the current user.
// This is useful when debugging authentication and authorization.
func whoamiHandler(resp http.ResponseWriter, req *http.Request) {
	authMetadata := auth.AuthenticationMetadataFromContext(req.Context()).GetRaw()

	payload := struct {
		AuthenticationMetadata any `json:"authenticationMetadata"`
	}{
		AuthenticationMetadata: authMetadata,
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(payload); err != nil {
		http.Error(resp, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
