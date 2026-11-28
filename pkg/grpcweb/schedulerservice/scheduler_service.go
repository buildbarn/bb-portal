package schedulerservice

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	bb_grpcweb "github.com/buildbarn/bb-portal/pkg/grpcweb"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	"github.com/buildbarn/bb-storage/pkg/digest"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/proto/configuration/builder"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewSchedulerService initializes and configures a gRPC-Web proxy server the
// BuildQueueState API, and starts an endpoint used for checking permissions.
// It registers all routes it handles with the provided router.
func NewSchedulerService(
	configuration *bb_portal.SchedulerService,
	siblingsGroup program.Group,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
	router *http.ServeMux,
) error {
	if router == nil {
		return status.Error(codes.NotFound, "Failed to create Graphql endpoint. No http server configured")
	}
	if configuration.ListOperationsPageSize <= 0 {
		return status.Error(codes.NotFound, "No ListOperationsPageSize configured (or it is set to 0)")
	}

	if configuration.ReadAuthorizer == nil {
		return status.Error(codes.NotFound, "No ReadAuthorizer configured")
	}
	readAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.ReadAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create ReadAuthorizer")
	}

	if configuration.KillOperationsAuthorizer == nil {
		return status.Error(codes.NotFound, "No KillOperationsAuthorizer configured")
	}
	killOperationsAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.KillOperationsAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return util.StatusWrap(err, "Failed to create KillOperationsAuthorizer")
	}

	routes, err := newRoutes(configuration, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return err
	}
	// One proxy in front of all schedulers: the router translates instance
	// names back into Portal's namespace, so a single pair of authorizers
	// covers every scheduler.
	buildQueueStateClient := buildqueuestateproxy.NewDemultiplexingClient(routes)
	buildQueueStateServer := buildqueuestateproxy.NewBuildQueueStateServerImpl(buildQueueStateClient, readAuthorizer, killOperationsAuthorizer, configuration.ListOperationsPageSize)

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	buildqueuestate.RegisterBuildQueueStateServer(grpcServer, buildQueueStateServer)

	bb_grpcweb.AddGrpcWebEndpoint(router, grpcWebServer, "/buildbarn.buildqueuestate.BuildQueueState/")
	router.HandleFunc("GET /api/v1/checkPermissions/killOperation/{operationName}", buildQueueStateServer.CheckKillOperationAuthorization)
	return registerCatalog(router, routes)
}

// newRoutes mirrors bb_storage's NewDemultiplexingBuildQueueFromConfiguration,
// so that the same scheduler map means the same thing in both binaries.
func newRoutes(configuration *bb_portal.SchedulerService, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory) ([]buildqueuestateproxy.Route, error) {
	schedulers := configuration.Schedulers
	if configuration.BuildQueueStateClient != nil {
		if len(schedulers) > 0 {
			return nil, status.Error(codes.InvalidArgument, "build_queue_state_client and schedulers cannot be combined")
		}
		schedulers = map[string]*builder.SchedulerConfiguration{"": {Endpoint: configuration.BuildQueueStateClient}}
	}
	if len(schedulers) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No build_queue_state_client or schedulers configured")
	}

	// Sorted, so that startup errors and the catalog are deterministic.
	keys := make([]string, 0, len(schedulers))
	for k := range schedulers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	routes := make([]buildqueuestateproxy.Route, 0, len(keys))
	for _, k := range keys {
		scheduler := schedulers[k]
		matchInstanceNamePrefix, err := digest.NewInstanceName(k)
		if err != nil {
			return nil, util.StatusWrapf(err, "Invalid instance name %#v", k)
		}
		addInstanceNamePrefix, err := digest.NewInstanceName(scheduler.GetAddInstanceNamePrefix())
		if err != nil {
			return nil, util.StatusWrapf(err, "Invalid instance name %#v", scheduler.GetAddInstanceNamePrefix())
		}
		if scheduler.GetEndpoint() == nil {
			return nil, status.Errorf(codes.InvalidArgument, "No endpoint configured for instance name %#v", k)
		}
		// Clients connect lazily, so an unreachable scheduler only fails
		// its own requests.
		endpoint, err := grpcClientFactory.NewClientFromConfiguration(scheduler.Endpoint, dependenciesGroup)
		if err != nil {
			return nil, util.StatusWrapf(err, "Failed to create BuildQueueState client for instance name %#v", k)
		}
		routes = append(routes, buildqueuestateproxy.Route{
			MatchInstanceNamePrefix: matchInstanceNamePrefix,
			AddInstanceNamePrefix:   addInstanceNamePrefix,
			Client:                  buildqueuestate.NewBuildQueueStateClient(endpoint),
		})
	}
	return routes, nil
}

type catalogEntry struct {
	InstanceNamePrefix string `json:"instanceNamePrefix"`
}

type catalog struct {
	Schedulers []catalogEntry `json:"schedulers"`
}

// registerCatalog serves the instance name prefixes the web UI may select.
// Endpoints and authorizers are never exposed.
func registerCatalog(router *http.ServeMux, routes []buildqueuestateproxy.Route) error {
	c := catalog{Schedulers: make([]catalogEntry, 0, len(routes))}
	for _, r := range routes {
		c.Schedulers = append(c.Schedulers, catalogEntry{InstanceNamePrefix: r.MatchInstanceNamePrefix.String()})
	}
	body, err := json.Marshal(c)
	if err != nil {
		return util.StatusWrap(err, "Failed to marshal scheduler catalog")
	}
	router.HandleFunc("GET /api/v1/schedulers", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if _, err := w.Write(body); err != nil {
			log.Printf("Failed to write scheduler catalog: %v", err)
		}
	})
	return nil
}
