package main

import (
	"context"
	"os"
	"reflect"
	"time"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"

	"github.com/buildbarn/bb-portal/internal/bep"
	"github.com/buildbarn/bb-portal/pkg/frontend"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/blobstoreservice"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/schedulerservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/global"
	http_server "github.com/buildbarn/bb-storage/pkg/http/server"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		lifecycleState, grpcClientFactory, err := global.ApplyConfiguration(configuration.Global, dependenciesGroup)
		if err != nil {
			return util.StatusWrap(err, "Failed to apply global configuration options")
		}

		tracerProvider := otel.GetTracerProvider()
		if tracerProvider == nil || reflect.ValueOf(tracerProvider).IsNil() {
			return status.Error(codes.Internal, "Otel tracer provider is nil")
		}

		if configuration.InstanceNameAuthorizer == nil {
			return status.Error(codes.NotFound, "No InstanceNameAuthorizer configured")
		}
		instanceNameAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(
			configuration.InstanceNameAuthorizer,
			dependenciesGroup,
			grpcClientFactory,
		)
		if err != nil {
			return util.StatusWrap(err, "Failed to create InstanceNameAuthorizer")
		}

		router := mux.NewRouter()
		router.Use(otelmux.Middleware("bb-portal-http", otelmux.WithTracerProvider(tracerProvider)))

		if err = blobstoreservice.NewBlobstoreService(
			&configuration,
			siblingsGroup,
			dependenciesGroup,
			grpcClientFactory,
			instanceNameAuthorizer,
			router,
		); err != nil {
			return util.StatusWrap(err, "Failed to create Blobstore service")
		}
		if configuration.SchedulerServiceConfiguration != nil {
			if err = schedulerservice.NewSchedulerService(
				configuration.SchedulerServiceConfiguration,
				siblingsGroup,
				dependenciesGroup,
				grpcClientFactory,
				instanceNameAuthorizer,
				router,
			); err != nil {
				return util.StatusWrap(err, "Failed to create Scheduler service")
			}
		}
		if configuration.BesServiceConfiguration != nil {
			if err = bep.NewBuildEventProtocolService(
				configuration.BesServiceConfiguration,
				siblingsGroup,
				dependenciesGroup,
				grpcClientFactory,
				instanceNameAuthorizer,
				router,
				tracerProvider,
			); err != nil {
				return util.StatusWrap(err, "Failed to create BES service")
			}
		}

		// This must be the last service created for the router, as it will
		// handle all unmatched requests.
		err = frontend.ServeFrontend(configuration.FrontendServiceConfiguration, router)
		if err != nil {
			return util.StatusWrap(err, "Failed to create frontend proxy service")
		}

		http_server.NewServersFromConfigurationAndServe(
			configuration.HttpServers,
			http_server.NewMetricsHandler(router, "PortalUI"),
			siblingsGroup,
			grpcClientFactory,
		)

		lifecycleState.MarkReadyAndWait(siblingsGroup)
		return nil
	})
}
