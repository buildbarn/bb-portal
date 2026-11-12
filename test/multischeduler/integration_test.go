// Package multischeduler tests Portal's production HTTP routing against two
// independent BuildQueueState gRPC servers, without a database or live services.
package multischeduler

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/schedulerservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	bqs "github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/program"
	authcfg "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	grpccfg "github.com/buildbarn/bb-storage/pkg/proto/configuration/grpc"
	jmespathcfg "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeScheduler struct {
	bqs.UnimplementedBuildQueueStateServer
	marker   string
	mu       sync.Mutex
	killed   []string
	canceled chan struct{}
	started  chan struct{}
}

func operation(marker, name string) *bqs.OperationState {
	return &bqs.OperationState{Name: name, TargetId: marker, InvocationName: &bqs.InvocationName{SizeClassQueueName: &bqs.SizeClassQueueName{PlatformQueueName: &bqs.PlatformQueueName{InstanceNamePrefix: "shared-instance"}}}}
}

func (s *fakeScheduler) GetOperation(ctx context.Context, req *bqs.GetOperationRequest) (*bqs.GetOperationResponse, error) {
	if req.OperationName == "slow" {
		select {
		case s.started <- struct{}{}:
		default:
		}
		<-ctx.Done()
		select {
		case s.canceled <- struct{}{}:
		default:
		}
		return nil, ctx.Err()
	}
	return &bqs.GetOperationResponse{Operation: operation(s.marker, req.OperationName)}, nil
}

func (s *fakeScheduler) ListPlatformQueues(context.Context, *emptypb.Empty) (*bqs.ListPlatformQueuesResponse, error) {
	return &bqs.ListPlatformQueuesResponse{PlatformQueues: []*bqs.PlatformQueueState{{Name: &bqs.PlatformQueueName{InstanceNamePrefix: "shared-instance"}, SizeClassQueues: []*bqs.SizeClassQueueState{{SizeClass: uint32(len(s.marker))}}}}}, nil
}

func (s *fakeScheduler) ListWorkers(context.Context, *bqs.ListWorkersRequest) (*bqs.ListWorkersResponse, error) {
	return &bqs.ListWorkersResponse{Workers: []*bqs.WorkerState{{Id: map[string]string{"id": "same-worker"}, CurrentOperation: operation(s.marker, "same-operation")}}, PaginationInfo: &bqs.PaginationInfo{TotalEntries: 1}}, nil
}

func (s *fakeScheduler) ListOperations(_ context.Context, req *bqs.ListOperationsRequest) (*bqs.ListOperationsResponse, error) {
	cursor := req.GetStartAfter().GetOperationName()
	var ops []*bqs.OperationState
	for _, name := range []string{"one", "same-operation", "three"} {
		if name > cursor {
			ops = append(ops, operation(s.marker, name))
			if uint32(len(ops)) == req.PageSize {
				break
			}
		}
	}
	return &bqs.ListOperationsResponse{Operations: ops, PaginationInfo: &bqs.PaginationInfo{TotalEntries: 3}}, nil
}

func (s *fakeScheduler) KillOperations(_ context.Context, req *bqs.KillOperationsRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.killed = append(s.killed, req.GetFilter().GetOperationName())
	return &emptypb.Empty{}, nil
}

type clientFactory struct{ t *testing.T }

func (f clientFactory) NewClientFromConfiguration(c *grpccfg.ClientConfiguration, _ program.Group) (grpc.ClientConnInterface, error) {
	conn, err := grpc.NewClient(c.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		f.t.Cleanup(func() { conn.Close() })
	}
	return conn, err
}

func startScheduler(t *testing.T, marker string) (*fakeScheduler, string) {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := &fakeScheduler{marker: marker, canceled: make(chan struct{}, 1), started: make(chan struct{}, 1)}
	server := grpc.NewServer()
	bqs.RegisterBuildQueueStateServer(server, s)
	go server.Serve(lis)
	t.Cleanup(server.Stop)
	return s, lis.Addr().String()
}

func authorizer(allow bool) *authcfg.AuthorizerConfiguration {
	if allow {
		return &authcfg.AuthorizerConfiguration{Policy: &authcfg.AuthorizerConfiguration_Allow{Allow: &emptypb.Empty{}}}
	}
	return &authcfg.AuthorizerConfiguration{Policy: &authcfg.AuthorizerConfiguration_Deny{Deny: &emptypb.Empty{}}}
}

func config(address string, read, kill bool) *bb_portal.SchedulerService {
	return &bb_portal.SchedulerService{BuildQueueStateClient: &grpccfg.ClientConfiguration{Address: address}, ReadAuthorizer: authorizer(read), KillOperationsAuthorizer: authorizer(kill), ListOperationsPageSize: 2}
}

func portal(t *testing.T, c *bb_portal.ApplicationConfiguration) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	if err := schedulerservice.NewSchedulerServices(c, nil, clientFactory{t}, mux); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(common.ForwardHeadersAsGrpcMetadata(mux))
	t.Cleanup(server.Close)
	return server
}

// grpcWeb exercises real framing and service paths, including trailer failures.
func grpcWeb(ctx context.Context, url, method string, in, out proto.Message, headers ...http.Header) (string, error) {
	payload, err := proto.Marshal(in)
	if err != nil {
		return "", err
	}
	frame := make([]byte, 5+len(payload))
	binary.BigEndian.PutUint32(frame[1:5], uint32(len(payload)))
	copy(frame[5:], payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url+"/buildbarn.buildqueuestate.BuildQueueState/"+method, bytes.NewReader(frame))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/grpc-web+proto")
	req.Header.Set("X-Grpc-Web", "1")
	for _, h := range headers {
		for key, values := range h {
			req.Header[key] = values
		}
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status %d", response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	trailers := ""
	for len(data) >= 5 {
		n := int(binary.BigEndian.Uint32(data[1:5]))
		if n > len(data)-5 {
			return "", fmt.Errorf("truncated grpc-web response")
		}
		body := data[5 : 5+n]
		if data[0]&0x80 != 0 {
			trailers = string(body)
		} else if out != nil {
			if err := proto.Unmarshal(body, out); err != nil {
				return "", err
			}
		}
		data = data[5+n:]
	}
	if len(data) != 0 {
		return "", fmt.Errorf("unexpected trailing bytes")
	}
	if trailers == "" {
		trailers = fmt.Sprintf("grpc-status: %s\r\ngrpc-message: %s\r\ngrpc-status-details-bin: %s\r\n", response.Header.Get("Grpc-Status"), response.Header.Get("Grpc-Message"), response.Header.Get("Grpc-Status-Details-Bin"))
	}
	return trailers, nil
}

func call(t *testing.T, url, method string, in, out proto.Message, want codes.Code) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	trailer, err := grpcWeb(ctx, url, method, in, out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(trailer, fmt.Sprintf("grpc-status: %d", want)) {
		t.Fatalf("%s: expected %s, got %q", method, want, trailer)
	}
}

func permission(t *testing.T, url string, want bool) {
	t.Helper()
	r, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	var body struct {
		Allowed bool `json:"allowed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 200 || body.Allowed != want {
		t.Fatalf("permission %s: HTTP %d allowed %t", url, r.StatusCode, body.Allowed)
	}
}

func TestTwoSchedulersThroughHTTP(t *testing.T) {
	a, aa := startScheduler(t, "alpha")
	b, ba := startScheduler(t, "beta-beta")
	p := portal(t, &bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{{Id: "a", Label: "Alpha", Configuration: config(aa, true, false)}, {Id: "b", Label: "Beta", Configuration: config(ba, true, true)}, {Id: "hidden", Label: "Hidden", Configuration: config(ba, false, true)}}, DefaultSchedulerId: "a"})
	base := func(id string) string { return p.URL + "/api/v1/schedulers/" + id + "/grpcweb" }
	for _, tc := range []struct{ id, marker string }{{"a", "alpha"}, {"b", "beta-beta"}} {
		t.Run(tc.id, func(t *testing.T) {
			var op bqs.GetOperationResponse
			call(t, base(tc.id), "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &op, codes.OK)
			if op.GetOperation().GetTargetId() != tc.marker {
				t.Fatal("operation crossed schedulers")
			}
			var queues bqs.ListPlatformQueuesResponse
			call(t, base(tc.id), "ListPlatformQueues", &emptypb.Empty{}, &queues, codes.OK)
			if queues.GetPlatformQueues()[0].GetSizeClassQueues()[0].GetSizeClass() != uint32(len(tc.marker)) {
				t.Fatal("queues crossed schedulers")
			}
			var workers bqs.ListWorkersResponse
			call(t, base(tc.id), "ListWorkers", &bqs.ListWorkersRequest{Filter: &bqs.ListWorkersRequest_Filter{Type: &bqs.ListWorkersRequest_Filter_All{All: &bqs.SizeClassQueueName{PlatformQueueName: &bqs.PlatformQueueName{InstanceNamePrefix: "shared-instance"}}}}}, &workers, codes.OK)
			if workers.GetWorkers()[0].GetId()["id"] != "same-worker" || workers.GetWorkers()[0].GetCurrentOperation().GetTargetId() != tc.marker {
				t.Fatal("worker crossed schedulers")
			}
			var ops bqs.ListOperationsResponse
			call(t, base(tc.id), "ListOperations", &bqs.ListOperationsRequest{PageSize: 2}, &ops, codes.OK)
			if len(ops.Operations) != 2 || ops.PaginationInfo.TotalEntries != 3 {
				t.Fatalf("pagination: %v", &ops)
			}
			call(t, base(tc.id), "ListOperations", &bqs.ListOperationsRequest{PageSize: 2, StartAfter: &bqs.ListOperationsRequest_StartAfter{OperationName: ops.Operations[1].Name}}, &ops, codes.OK)
			if len(ops.Operations) != 1 || ops.Operations[0].TargetId != tc.marker {
				t.Fatal("page cursor crossed schedulers")
			}
		})
	}
	permission(t, p.URL+"/api/v1/schedulers/a/checkPermissions/killOperation/same-operation", false)
	permission(t, p.URL+"/api/v1/schedulers/b/checkPermissions/killOperation/same-operation", true)
	permission(t, p.URL+"/api/v1/schedulers/hidden/checkPermissions/killOperation/same-operation", false)
	kill := &bqs.KillOperationsRequest{Filter: &bqs.KillOperationsRequest_Filter{Type: &bqs.KillOperationsRequest_Filter_OperationName{OperationName: "same-operation"}}}
	call(t, base("a"), "KillOperations", kill, &emptypb.Empty{}, codes.PermissionDenied)
	call(t, base("b"), "KillOperations", kill, &emptypb.Empty{}, codes.OK)
	call(t, base("hidden"), "KillOperations", kill, &emptypb.Empty{}, codes.PermissionDenied)
	call(t, base("hidden"), "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, codes.NotFound)
	a.mu.Lock()
	ak := len(a.killed)
	a.mu.Unlock()
	b.mu.Lock()
	bk := append([]string(nil), b.killed...)
	b.mu.Unlock()
	if ak != 0 || len(bk) != 1 || bk[0] != "same-operation" {
		t.Fatalf("crossed cancellation: a=%d b=%v", ak, bk)
	}
	var legacy bqs.GetOperationResponse
	call(t, p.URL+"/api/v1/grpcweb", "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &legacy, codes.OK)
	if legacy.Operation.TargetId != "alpha" {
		t.Fatal("legacy default wrong")
	}
	permission(t, p.URL+"/api/v1/checkPermissions/killOperation/same-operation", false)
	for _, path := range []string{"/api/v1/schedulers/missing/grpcweb/buildbarn.buildqueuestate.BuildQueueState/GetOperation", "/api/v1/schedulers/missing/checkPermissions/killOperation/same-operation", "/api/v1/schedulers/INVALID/checkPermissions/killOperation/same-operation"} {
		r, err := http.Get(p.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		if r.StatusCode != http.StatusNotFound {
			t.Fatalf("unknown scheduler path %s: %d", path, r.StatusCode)
		}
	}
	// Cancellation of a browser request must propagate through Portal and gRPC.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	requestDone := make(chan error, 1)
	go func() {
		_, err := grpcWeb(ctx, base("a"), "GetOperation", &bqs.GetOperationRequest{OperationName: "slow"}, &bqs.GetOperationResponse{})
		requestDone <- err
	}()
	select {
	case <-a.started:
	case <-time.After(3 * time.Second):
		t.Fatal("slow upstream did not start")
	}
	// B must remain usable while A's upstream request is still blocked.
	call(t, base("b"), "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, codes.OK)
	cancel()
	if err := <-requestDone; err == nil {
		t.Fatal("slow request did not cancel")
	}

	select {
	case <-a.canceled:
	case <-time.After(3 * time.Second):
		t.Fatal("browser cancellation did not reach scheduler")
	}
	call(t, base("b"), "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, codes.OK)
}

// The production HTTP authentication layer places authentication metadata in
// the context before scheduler dispatch. This fixture uses a test-only trusted
// identity header to verify both transports preserve that context.
func TestAuthenticationContextAcrossScopedRoutes(t *testing.T) {
	_, address := startScheduler(t, "authenticated")
	policy := &authcfg.AuthorizerConfiguration{Policy: &authcfg.AuthorizerConfiguration_JmespathExpression{JmespathExpression: &jmespathcfg.Expression{Expression: "authenticationMetadata.public.principal == 'alice'"}}}
	cfg := config(address, true, true)
	cfg.ReadAuthorizer = policy
	cfg.KillOperationsAuthorizer = policy
	mux := http.NewServeMux()
	if err := schedulerservice.NewSchedulerServices(&bb_portal.ApplicationConfiguration{NamedSchedulers: []*bb_portal.NamedScheduler{{Id: "auth", Configuration: cfg}}}, nil, clientFactory{t}, mux); err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadata, err := auth.NewAuthenticationMetadataFromRaw(map[string]interface{}{"public": map[string]interface{}{"principal": r.Header.Get("X-Test-Principal")}})
		if err != nil {
			t.Error(err)
			http.Error(w, "auth fixture", 500)
			return
		}
		mux.ServeHTTP(w, r.WithContext(auth.NewContextWithAuthenticationMetadata(r.Context(), metadata)))
	})
	p := httptest.NewServer(common.ForwardHeadersAsGrpcMetadata(handler))
	defer p.Close()
	for _, tc := range []struct {
		principal string
		allowed   bool
		code      codes.Code
	}{{"alice", true, codes.OK}, {"bob", false, codes.NotFound}} {
		t.Run(tc.principal, func(t *testing.T) {
			headers := http.Header{"X-Test-Principal": []string{tc.principal}}
			trailer, err := grpcWeb(context.Background(), p.URL+"/api/v1/schedulers/auth/grpcweb", "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, headers)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(trailer, fmt.Sprintf("grpc-status: %d", tc.code)) {
				t.Fatalf("authentication context lost: %s", trailer)
			}
			req, err := http.NewRequest("GET", p.URL+"/api/v1/schedulers/auth/checkPermissions/killOperation/same-operation", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header = headers
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			var permission struct {
				Allowed bool `json:"allowed"`
			}
			if err := json.NewDecoder(response.Body).Decode(&permission); err != nil {
				t.Fatal(err)
			}
			if permission.Allowed != tc.allowed {
				t.Fatal("permission lost authentication context")
			}
		})
	}
}

func TestSingleSchedulerConfigurationsThroughHTTP(t *testing.T) {
	for _, mode := range []string{"legacy", "single-named"} {
		t.Run(mode, func(t *testing.T) {
			_, address := startScheduler(t, "single-scheduler")
			configuration := &bb_portal.ApplicationConfiguration{}
			id := "default"
			if mode == "legacy" {
				configuration.SchedulerServiceConfiguration = config(address, true, true)
			} else {
				id = "only"
				configuration.NamedSchedulers = []*bb_portal.NamedScheduler{{Id: id, Configuration: config(address, true, true)}}
			}
			// Neither configuration supplies an explicit default.
			p := portal(t, configuration)
			r, err := http.Get(p.URL + "/api/v1/schedulers")
			if err != nil {
				t.Fatal(err)
			}
			defer r.Body.Close()
			var catalog struct {
				Schedulers []struct {
					ID string `json:"id"`
				} `json:"schedulers"`
				DefaultSchedulerID string `json:"defaultSchedulerId"`
			}
			if err := json.NewDecoder(r.Body).Decode(&catalog); err != nil {
				t.Fatal(err)
			}
			if r.StatusCode != http.StatusOK || catalog.DefaultSchedulerID != id || len(catalog.Schedulers) != 1 || catalog.Schedulers[0].ID != id {
				t.Fatalf("unexpected single-scheduler catalog: HTTP %d, %+v", r.StatusCode, catalog)
			}
			for _, prefix := range []string{"/api/v1/grpcweb", "/api/v1/schedulers/" + id + "/grpcweb"} {
				var response bqs.GetOperationResponse
				call(t, p.URL+prefix, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &response, codes.OK)
				if response.GetOperation().GetTargetId() != "single-scheduler" {
					t.Fatalf("%s returned the wrong operation: %v", prefix, &response)
				}
			}
			permission(t, p.URL+"/api/v1/checkPermissions/killOperation/same-operation", true)
			permission(t, p.URL+"/api/v1/schedulers/"+id+"/checkPermissions/killOperation/same-operation", true)
		})
	}
}

func TestUnreachableSchedulerDoesNotBlockHealthyScheduler(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	unreachableAddress := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	_, healthyAddress := startScheduler(t, "healthy")
	// The default scheduler is already unreachable when Portal initializes.
	p := portal(t, &bb_portal.ApplicationConfiguration{
		NamedSchedulers: []*bb_portal.NamedScheduler{
			{Id: "down", Configuration: config(unreachableAddress, true, false)},
			{Id: "healthy", Configuration: config(healthyAddress, true, false)},
		},
		DefaultSchedulerId: "down",
	})
	for _, prefix := range []string{"/api/v1/schedulers/down/grpcweb", "/api/v1/grpcweb"} {
		var response bqs.GetOperationResponse
		call(t, p.URL+"/api/v1/schedulers/healthy/grpcweb", "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &response, codes.OK)
		if response.GetOperation().GetTargetId() != "healthy" {
			t.Fatalf("wrong healthy scheduler response: %v", &response)
		}
		call(t, p.URL+prefix, "ListPlatformQueues", &emptypb.Empty{}, &bqs.ListPlatformQueuesResponse{}, codes.Unavailable)
	}
	var response bqs.GetOperationResponse
	call(t, p.URL+"/api/v1/schedulers/healthy/grpcweb", "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &response, codes.OK)
	if response.GetOperation().GetTargetId() != "healthy" {
		t.Fatalf("unreachable scheduler affected healthy response: %v", &response)
	}
}
