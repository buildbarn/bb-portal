package schedulerservice

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	grpc_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/grpc"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	go_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testService() *bb_portal.SchedulerService {
	return &bb_portal.SchedulerService{
		BuildQueueStateClient:    &grpc_pb.ClientConfiguration{Address: "scheduler:4888"},
		ListOperationsPageSize:   500,
		ReadAuthorizer:           &auth_pb.AuthorizerConfiguration{},
		KillOperationsAuthorizer: &auth_pb.AuthorizerConfiguration{},
	}
}

func named(id string) *bb_portal.NamedScheduler {
	return &bb_portal.NamedScheduler{Id: id, Label: id, Configuration: testService()}
}

func TestNormalizeConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		config        *bb_portal.ApplicationConfiguration
		defaultID     string
		count         int
		errorContains string
	}{
		{name: "disabled", config: &bb_portal.ApplicationConfiguration{}},
		{name: "legacy", config: &bb_portal.ApplicationConfiguration{SchedulerServiceConfiguration: testService()}, defaultID: "default", count: 1},
		{name: "named single", config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{named("ci")}}, defaultID: "ci", count: 1},
		{name: "named multiple", config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{named("ci"), named("dev")}, DefaultSchedulerId: "dev"}, defaultID: "dev", count: 2},
		{name: "mixed", config: &bb_portal.ApplicationConfiguration{SchedulerServiceConfiguration: testService(), NamedSchedulers: []*bb_portal.NamedScheduler{named("ci")}}, errorContains: "cannot be combined"},
		{name: "duplicate", config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{named("ci"), named("ci")}}, errorContains: "duplicate"},
		{name: "unknown default", config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{named("ci")}, DefaultSchedulerId: "dev"}, errorContains: "does not identify"},
		{name: "disabled with default", config: &bb_portal.ApplicationConfiguration{DefaultSchedulerId: "ci"}, errorContains: "does not identify"},
		{name: "ambiguous default", config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{named("ci"), named("dev")}}, errorContains: "required"},
		{name: "nil scheduler", config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{nil}}, errorContains: "invalid scheduler ID"},
	}
	for _, id := range []string{"", "../ci", "CI", "a/b", "a%2fb", "-ci", strings.Repeat("a", 64)} {
		tests = append(tests, struct {
			name          string
			config        *bb_portal.ApplicationConfiguration
			defaultID     string
			count         int
			errorContains string
		}{name: "invalid " + id, config: &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{named(id)}}, errorContains: "invalid scheduler ID"})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schedulers, defaultID, err := normalizeConfiguration(test.config)
			if test.errorContains != "" {
				if err == nil || !strings.Contains(err.Error(), test.errorContains) {
					t.Fatalf("expected %q error, got %v", test.errorContains, err)
				}
				if status.Code(err) != codes.InvalidArgument {
					t.Fatalf("expected InvalidArgument, got %v", err)
				}
				return
			}
			if err != nil || defaultID != test.defaultID || len(schedulers) != test.count {
				t.Fatalf("got %d schedulers, default %q, error %v", len(schedulers), defaultID, err)
			}
		})
	}
	for _, field := range []string{"configuration", "client", "address", "page size", "read authorizer", "kill authorizer"} {
		t.Run("missing "+field, func(t *testing.T) {
			scheduler := named("ci")
			switch field {
			case "configuration":
				scheduler.Configuration = nil
			case "client":
				scheduler.Configuration.BuildQueueStateClient = nil
			case "address":
				scheduler.Configuration.BuildQueueStateClient.Address = " "
			case "page size":
				scheduler.Configuration.ListOperationsPageSize = 0
			case "read authorizer":
				scheduler.Configuration.ReadAuthorizer = nil
			case "kill authorizer":
				scheduler.Configuration.KillOperationsAuthorizer = nil
			}
			if _, _, err := normalizeConfiguration(&bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{scheduler}}); status.Code(err) != codes.InvalidArgument {
				t.Fatalf("expected InvalidArgument, got %v", err)
			}
		})
	}
}

func TestRoutesKeepSchedulerAndOperationTogether(t *testing.T) {
	mux := http.NewServeMux()
	backends := map[string]*schedulerBackend{}
	var called string
	for _, id := range []string{"ci", "dev"} {
		backends[id] = &schedulerBackend{
			grpc: grpcweb.WrapServer(go_grpc.NewServer(go_grpc.UnknownServiceHandler(func(_ any, stream go_grpc.ServerStream) error {
				method, _ := go_grpc.MethodFromServerStream(stream)
				called = id + ":" + method
				return status.Error(codes.Unimplemented, "test endpoint")
			}))),
			permissions: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(id + ":" + r.PathValue("operationName"))) }),
		}
	}
	if err := registerSchedulerRoutes(mux, backends, schedulerCatalog{Schedulers: []schedulerMetadata{{ID: "ci", Label: "CI"}, {ID: "dev", Label: "Dev"}}, DefaultSchedulerID: "ci"}); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		url, want string
		status    int
	}{
		{"/api/v1/schedulers/ci/grpcweb" + servicePath + "GetOperation", "ci:" + servicePath + "GetOperation", 200},
		{"/api/v1/schedulers/dev/grpcweb" + servicePath + "GetOperation", "dev:" + servicePath + "GetOperation", 200},
		{"/api/v1/grpcweb" + servicePath + "GetOperation", "ci:" + servicePath + "GetOperation", 200},
		{"/api/v1/schedulers/ci/checkPermissions/killOperation/collision", "ci:collision", 200},
		{"/api/v1/schedulers/dev/checkPermissions/killOperation/collision", "dev:collision", 200},
		{"/api/v1/checkPermissions/killOperation/collision", "ci:collision", 200},
		{"/api/v1/schedulers/unknown/grpcweb" + servicePath + "GetOperation", "Unknown scheduler ID\n", 404},
		{"/api/v1/schedulers/CI/checkPermissions/killOperation/collision", "Unknown scheduler ID\n", 404},
		{"/api/v1/schedulers/dev/grpcweb/other.Service/Method", "404 page not found\n", 404},
	}
	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			called = ""
			r := httptest.NewRequest(http.MethodGet, test.url, nil)
			if strings.Contains(test.url, "/grpcweb/") {
				r = httptest.NewRequest(http.MethodPost, test.url, bytes.NewReader(make([]byte, 5)))
				r.Header.Set("Content-Type", "application/grpc-web+proto")
				r.Header.Set("X-Grpc-Web", "1")
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			got := w.Body.String()
			if strings.Contains(test.url, "/grpcweb/") && test.status == http.StatusOK {
				got = called
			}
			if w.Code != test.status || got != test.want {
				t.Fatalf("got %d %q, want %d %q", w.Code, w.Body.String(), test.status, test.want)
			}
		})
	}
	for _, url := range []string{"/api/v1/schedulers", "/api/v1/schedulers/"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, url, nil))
		var catalog schedulerCatalog
		if err := json.Unmarshal(w.Body.Bytes(), &catalog); err != nil || catalog.DefaultSchedulerID != "ci" || len(catalog.Schedulers) != 2 {
			t.Fatalf("invalid public catalog: %s (%v)", w.Body.String(), err)
		}
	}
}

func TestDisabledCatalog(t *testing.T) {
	mux := http.NewServeMux()
	if err := NewSchedulerServices(&bb_portal.ApplicationConfiguration{}, nil, nil, mux); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/schedulers", nil))
	if w.Code != 200 || w.Body.String() != "{\"schedulers\":[],\"defaultSchedulerId\":\"\"}\n" {
		t.Fatalf("unexpected disabled catalog %d %s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/schedulers/ci/grpcweb"+servicePath+"GetOperation", nil))
	if w.Code != 404 {
		t.Fatalf("disabled scheduler request returned %d", w.Code)
	}
}

func TestInitializationWithoutHTTPServer(t *testing.T) {
	for _, config := range []*bb_portal.ApplicationConfiguration{nil, {}} {
		if err := NewSchedulerServices(config, nil, nil, nil); err != nil {
			t.Fatalf("disabled scheduler service requires no HTTP server: %v", err)
		}
	}
	for _, config := range []*bb_portal.ApplicationConfiguration{
		{NamedSchedulers: []*bb_portal.NamedScheduler{named("ci")}},
		{NamedSchedulers: []*bb_portal.NamedScheduler{named("../ci")}},
		{DefaultSchedulerId: "missing"},
	} {
		if err := NewSchedulerServices(config, nil, nil, nil); status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument for configured/invalid scheduler without HTTP server, got %v", err)
		}
	}
}
