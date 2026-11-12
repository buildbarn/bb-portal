package schedulerservice

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	bb_grpcweb "github.com/buildbarn/bb-portal/pkg/grpcweb"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	auth_configuration "github.com/buildbarn/bb-storage/pkg/auth/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// newSchedulerBackend creates one reusable proxy with its own authorization policy.
func newSchedulerBackend(
	configuration *bb_portal.SchedulerService,
	dependenciesGroup program.Group,
	grpcClientFactory bb_grpc.ClientFactory,
) (*schedulerBackend, error) {
	readAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.ReadAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create ReadAuthorizer")
	}

	killOperationsAuthorizer, err := auth_configuration.DefaultAuthorizerFactory.NewAuthorizerFromConfiguration(configuration.KillOperationsAuthorizer, dependenciesGroup, grpcClientFactory)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create KillOperationsAuthorizer")
	}

	grpcClient, err := grpcClientFactory.NewClientFromConfiguration(configuration.BuildQueueStateClient, dependenciesGroup)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create gRPC client for BuildQueueState")
	}
	buildQueueStateClient := buildqueuestate.NewBuildQueueStateClient(grpcClient)
	buildQueueStateServer := buildqueuestateproxy.NewBuildQueueStateServerImpl(buildQueueStateClient, readAuthorizer, killOperationsAuthorizer, configuration.ListOperationsPageSize)

	grpcServer := go_grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	buildqueuestate.RegisterBuildQueueStateServer(grpcServer, buildQueueStateServer)

	return &schedulerBackend{grpc: grpcWebServer, permissions: http.HandlerFunc(buildQueueStateServer.CheckKillOperationAuthorization)}, nil
}

const servicePath = "/buildbarn.buildqueuestate.BuildQueueState/"

var schedulerIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,62}$`)

type schedulerBackend struct {
	grpc        *grpcweb.WrappedGrpcServer
	permissions http.Handler
}

type schedulerMetadata struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type schedulerCatalog struct {
	Schedulers         []schedulerMetadata `json:"schedulers"`
	DefaultSchedulerID string              `json:"defaultSchedulerId"`
}

// normalizeConfiguration validates identities and configuration without contacting
// schedulers. A temporary backend outage must not prevent another scheduler's use.
func normalizeConfiguration(configuration *bb_portal.ApplicationConfiguration) ([]*bb_portal.NamedScheduler, string, error) {
	schedulers := configuration.GetNamedSchedulers()
	defaultID := configuration.GetDefaultSchedulerId()
	if legacy := configuration.GetSchedulerServiceConfiguration(); legacy != nil {
		if len(schedulers) != 0 {
			return nil, "", status.Errorf(codes.InvalidArgument, "scheduler_service_configuration and named_schedulers cannot be combined")
		}
		schedulers = []*bb_portal.NamedScheduler{{Id: "default", Label: "Default", Configuration: legacy}}
	}
	seen := map[string]bool{}
	for _, scheduler := range schedulers {
		if scheduler == nil || !schedulerIDPattern.MatchString(scheduler.GetId()) {
			return nil, "", status.Errorf(codes.InvalidArgument, "invalid scheduler ID %q", scheduler.GetId())
		}
		if seen[scheduler.Id] {
			return nil, "", status.Errorf(codes.InvalidArgument, "duplicate scheduler ID %q", scheduler.Id)
		}
		seen[scheduler.Id] = true
		service := scheduler.Configuration
		if service == nil || service.BuildQueueStateClient == nil || strings.TrimSpace(service.BuildQueueStateClient.Address) == "" {
			return nil, "", status.Errorf(codes.InvalidArgument, "scheduler %q requires a BuildQueueState client address", scheduler.Id)
		}
		if service.ListOperationsPageSize == 0 || service.ReadAuthorizer == nil || service.KillOperationsAuthorizer == nil {
			return nil, "", status.Errorf(codes.InvalidArgument, "scheduler %q requires page size, read and cancellation authorizers", scheduler.Id)
		}
	}
	if defaultID == "" && len(schedulers) == 1 {
		defaultID = schedulers[0].Id
	}
	if defaultID != "" && !seen[defaultID] {
		return nil, "", status.Errorf(codes.InvalidArgument, "default_scheduler_id %q does not identify a configured scheduler", defaultID)
	}
	if len(schedulers) > 1 && defaultID == "" {
		return nil, "", status.Errorf(codes.InvalidArgument, "default_scheduler_id is required with multiple schedulers")
	}
	return schedulers, defaultID, nil
}

// NewSchedulerServices configures named scheduler proxies and legacy default aliases.
// The catalog contains only public identity metadata; backend configuration stays private.
func NewSchedulerServices(configuration *bb_portal.ApplicationConfiguration, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory, router *http.ServeMux) error {
	schedulers, defaultID, err := normalizeConfiguration(configuration)
	if err != nil {
		return err
	}
	if router == nil {
		if len(schedulers) == 0 {
			return nil
		}
		return status.Error(codes.InvalidArgument, "No HTTP server configured for schedulers")
	}
	backends := map[string]*schedulerBackend{}
	catalog := schedulerCatalog{Schedulers: []schedulerMetadata{}, DefaultSchedulerID: defaultID}
	for _, scheduler := range schedulers {
		backend, err := newSchedulerBackend(scheduler.Configuration, dependenciesGroup, grpcClientFactory)
		if err != nil {
			return util.StatusWrapf(err, "scheduler %q", scheduler.Id)
		}
		backends[scheduler.Id] = backend
		label := scheduler.Label
		if strings.TrimSpace(label) == "" {
			label = scheduler.Id
		}
		catalog.Schedulers = append(catalog.Schedulers, schedulerMetadata{ID: scheduler.Id, Label: label})
	}
	return registerSchedulerRoutes(router, backends, catalog)
}

func registerSchedulerRoutes(router *http.ServeMux, backends map[string]*schedulerBackend, catalog schedulerCatalog) error {
	// Marshal before registering handlers, so serialization failures cannot produce
	// partial successful HTTP responses. The catalog is immutable after startup.
	catalogJSON, err := json.Marshal(catalog)
	if err != nil {
		return util.StatusWrap(err, "Failed to marshal scheduler catalog")
	}
	catalogJSON = append(catalogJSON, '\n')
	serveCatalog := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if _, err := w.Write(catalogJSON); err != nil {
			log.Printf("Failed to write scheduler catalog: %v", err)
		}
	}
	router.HandleFunc("GET /api/v1/schedulers/{$}", serveCatalog)
	// Avoid a redirect for the endpoint used by the frontend.
	router.HandleFunc("GET /api/v1/schedulers", serveCatalog)
	scoped := map[string]http.Handler{}
	for id, backend := range backends {
		mux := http.NewServeMux()
		mux.Handle("/grpcweb"+servicePath, http.StripPrefix("/grpcweb", backend.grpc))
		mux.Handle("GET /checkPermissions/killOperation/{operationName}", backend.permissions)
		scoped[id] = http.StripPrefix("/api/v1/schedulers/"+id, mux)
	}
	router.HandleFunc("/api/v1/schedulers/{schedulerID}/", func(w http.ResponseWriter, r *http.Request) {
		handler, ok := scoped[r.PathValue("schedulerID")]
		if !ok {
			http.Error(w, "Unknown scheduler ID", http.StatusNotFound)
			return
		}
		handler.ServeHTTP(w, r)
	})
	if backend := backends[catalog.DefaultSchedulerID]; backend != nil {
		bb_grpcweb.AddGrpcWebEndpoint(router, backend.grpc, servicePath)
		router.Handle("GET /api/v1/checkPermissions/killOperation/{operationName}", backend.permissions)
	}
	return nil
}

// NewSchedulerService preserves the original initialization API for embedders.
func NewSchedulerService(configuration *bb_portal.SchedulerService, _, dependenciesGroup program.Group, grpcClientFactory bb_grpc.ClientFactory, router *http.ServeMux) error {
	return NewSchedulerServices(&bb_portal.ApplicationConfiguration{SchedulerServiceConfiguration: configuration}, dependenciesGroup, grpcClientFactory, router)
}
