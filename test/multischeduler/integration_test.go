// Package multischeduler tests Portal's production HTTP routing against two
// BuildQueueState gRPC servers keyed by instance name prefix, without a
// database or live services.
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
	"github.com/buildbarn/bb-portal/pkg/grpcweb/buildqueuestateproxy"
	"github.com/buildbarn/bb-portal/pkg/grpcweb/schedulerservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	bqs "github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/program"
	authcfg "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	"github.com/buildbarn/bb-storage/pkg/proto/configuration/builder"
	grpccfg "github.com/buildbarn/bb-storage/pkg/proto/configuration/grpc"
	jmespathcfg "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// fakeScheduler reports every queue and operation under instanceNamePrefix,
// i.e. in its own namespace. Operation and worker IDs deliberately collide
// across schedulers.
type fakeScheduler struct {
	bqs.UnimplementedBuildQueueStateServer
	marker             string
	instanceNamePrefix string
	mu                 sync.Mutex
	killed             []string
	listWorkersFor     []string
	canceled           chan struct{}
	started            chan struct{}
}

func (s *fakeScheduler) operation(name string) *bqs.OperationState {
	operation := &bqs.OperationState{Name: name, TargetId: s.marker, InvocationName: &bqs.InvocationName{SizeClassQueueName: &bqs.SizeClassQueueName{PlatformQueueName: &bqs.PlatformQueueName{InstanceNamePrefix: s.instanceNamePrefix}}}}
	if name == "restricted-operation" {
		operation.InstanceNameSuffix = "restricted"
	}
	return operation
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
	return &bqs.GetOperationResponse{Operation: s.operation(req.OperationName)}, nil
}

func (s *fakeScheduler) ListPlatformQueues(context.Context, *emptypb.Empty) (*bqs.ListPlatformQueuesResponse, error) {
	return &bqs.ListPlatformQueuesResponse{PlatformQueues: []*bqs.PlatformQueueState{{Name: &bqs.PlatformQueueName{InstanceNamePrefix: s.instanceNamePrefix}, SizeClassQueues: []*bqs.SizeClassQueueState{{SizeClass: uint32(len(s.marker))}}}}}, nil
}

func (s *fakeScheduler) ListWorkers(_ context.Context, req *bqs.ListWorkersRequest) (*bqs.ListWorkersResponse, error) {
	s.mu.Lock()
	s.listWorkersFor = append(s.listWorkersFor, req.GetFilter().GetAll().GetPlatformQueueName().GetInstanceNamePrefix())
	s.mu.Unlock()
	return &bqs.ListWorkersResponse{Workers: []*bqs.WorkerState{{Id: map[string]string{"id": "same-worker"}, CurrentOperation: s.operation("same-operation")}}, PaginationInfo: &bqs.PaginationInfo{TotalEntries: 1}}, nil
}

func (s *fakeScheduler) ListOperations(_ context.Context, req *bqs.ListOperationsRequest) (*bqs.ListOperationsResponse, error) {
	cursor := req.GetStartAfter().GetOperationName()
	var ops []*bqs.OperationState
	for _, name := range []string{"one", "same-operation", "three"} {
		if name > cursor {
			ops = append(ops, s.operation(name))
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

func startScheduler(t *testing.T, marker, instanceNamePrefix string) (*fakeScheduler, string) {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := &fakeScheduler{marker: marker, instanceNamePrefix: instanceNamePrefix, canceled: make(chan struct{}, 1), started: make(chan struct{}, 1)}
	server := grpc.NewServer()
	bqs.RegisterBuildQueueStateServer(server, s)
	go server.Serve(lis)
	t.Cleanup(server.Stop)
	return s, lis.Addr().String()
}

func allow() *authcfg.AuthorizerConfiguration {
	return &authcfg.AuthorizerConfiguration{Policy: &authcfg.AuthorizerConfiguration_Allow{Allow: &emptypb.Empty{}}}
}

func jmespath(expression string) *authcfg.AuthorizerConfiguration {
	return &authcfg.AuthorizerConfiguration{Policy: &authcfg.AuthorizerConfiguration_JmespathExpression{JmespathExpression: &jmespathcfg.Expression{Expression: expression}}}
}

func schedulerAt(address, addInstanceNamePrefix string) *builder.SchedulerConfiguration {
	return &builder.SchedulerConfiguration{Endpoint: &grpccfg.ClientConfiguration{Address: address}, AddInstanceNamePrefix: addInstanceNamePrefix}
}

func portal(t *testing.T, c *bb_portal.SchedulerService) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	if err := schedulerservice.NewSchedulerService(c, nil, nil, clientFactory{t}, mux); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(common.ForwardHeadersAsGrpcMetadata(mux))
	t.Cleanup(server.Close)
	return server
}

// selecting returns the header the web UI sends to scope requests that carry
// no instance name.
func selecting(instanceName string) http.Header {
	return http.Header{http.CanonicalHeaderKey(buildqueuestateproxy.SelectedInstanceNameMetadataKey): []string{instanceName}}
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
	req, err := http.NewRequestWithContext(ctx, "POST", url+"/api/v1/grpcweb/buildbarn.buildqueuestate.BuildQueueState/"+method, bytes.NewReader(frame))
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

func call(t *testing.T, url, method string, in, out proto.Message, want codes.Code, headers ...http.Header) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	trailer, err := grpcWeb(ctx, url, method, in, out, headers...)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(trailer, fmt.Sprintf("grpc-status: %d", want)) {
		t.Fatalf("%s: expected %s, got %q", method, want, trailer)
	}
}

func permission(t *testing.T, url string, want bool, headers http.Header) {
	permissionFor(t, url, "same-operation", want, headers)
}

func permissionFor(t *testing.T, url, operationName string, want bool, headers http.Header) {
	t.Helper()
	req, err := http.NewRequest("GET", url+"/api/v1/checkPermissions/killOperation/"+operationName, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header = headers
	r, err := http.DefaultClient.Do(req)
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
		t.Fatalf("permission with %v: HTTP %d allowed %t", headers, r.StatusCode, body.Allowed)
	}
}

func listWorkers(instanceNamePrefix string) *bqs.ListWorkersRequest {
	return &bqs.ListWorkersRequest{Filter: &bqs.ListWorkersRequest_Filter{Type: &bqs.ListWorkersRequest_Filter_All{All: &bqs.SizeClassQueueName{PlatformQueueName: &bqs.PlatformQueueName{InstanceNamePrefix: instanceNamePrefix}}}}}
}

// aws sees stripped names (bb_storage's default); gcp uses perfect
// forwarding. The kill authorizer is written against Portal's names only,
// which proves one global authorizer works for both.
func TestTwoSchedulersThroughHTTP(t *testing.T) {
	aws, awsAddress := startScheduler(t, "alpha", "mars")
	gcp, gcpAddress := startScheduler(t, "beta-beta", "ci/gcp/europe")
	p := portal(t, &bb_portal.SchedulerService{
		Schedulers: map[string]*builder.SchedulerConfiguration{
			"ci/aws": schedulerAt(awsAddress, ""),
			"ci/gcp": schedulerAt(gcpAddress, "ci/gcp"),
		},
		ReadAuthorizer:           allow(),
		KillOperationsAuthorizer: jmespath("starts_with(instanceName, 'ci/gcp/')"),
		ListOperationsPageSize:   2,
	})

	for _, tc := range []struct{ selected, marker, instanceName string }{
		{"ci/aws", "alpha", "ci/aws/mars"},
		{"ci/gcp", "beta-beta", "ci/gcp/europe"},
	} {
		t.Run(tc.selected, func(t *testing.T) {
			selection := selecting(tc.selected)
			var op bqs.GetOperationResponse
			call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &op, codes.OK, selection)
			if op.GetOperation().GetTargetId() != tc.marker || op.GetOperation().GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName().GetInstanceNamePrefix() != tc.instanceName {
				t.Fatalf("operation crossed schedulers or kept the scheduler's name: %v", &op)
			}
			var queues bqs.ListPlatformQueuesResponse
			call(t, p.URL, "ListPlatformQueues", &emptypb.Empty{}, &queues, codes.OK, selection)
			if len(queues.PlatformQueues) != 1 || queues.PlatformQueues[0].GetName().GetInstanceNamePrefix() != tc.instanceName || queues.PlatformQueues[0].GetSizeClassQueues()[0].GetSizeClass() != uint32(len(tc.marker)) {
				t.Fatalf("queues crossed schedulers: %v", &queues)
			}
			var ops bqs.ListOperationsResponse
			call(t, p.URL, "ListOperations", &bqs.ListOperationsRequest{PageSize: 2}, &ops, codes.OK, selection)
			if len(ops.Operations) != 2 || ops.PaginationInfo.TotalEntries != 3 {
				t.Fatalf("pagination: %v", &ops)
			}
			call(t, p.URL, "ListOperations", &bqs.ListOperationsRequest{PageSize: 2, StartAfter: &bqs.ListOperationsRequest_StartAfter{OperationName: ops.Operations[1].Name}}, &ops, codes.OK, selection)
			if len(ops.Operations) != 1 || ops.Operations[0].TargetId != tc.marker {
				t.Fatal("page cursor crossed schedulers")
			}
		})
	}

	t.Run("ListWorkersRoutesByInstanceName", func(t *testing.T) {
		// The selection points elsewhere on purpose: the queue name wins.
		var workers bqs.ListWorkersResponse
		call(t, p.URL, "ListWorkers", listWorkers("ci/aws/mars"), &workers, codes.OK, selecting("ci/gcp"))
		if workers.GetWorkers()[0].GetCurrentOperation().GetTargetId() != "alpha" {
			t.Fatal("worker request did not follow its instance name")
		}
		call(t, p.URL, "ListWorkers", listWorkers("ci/gcp/europe"), &workers, codes.OK)
		if workers.GetWorkers()[0].GetCurrentOperation().GetTargetId() != "beta-beta" {
			t.Fatal("worker request did not follow its instance name")
		}
		aws.mu.Lock()
		defer aws.mu.Unlock()
		gcp.mu.Lock()
		defer gcp.mu.Unlock()
		if len(aws.listWorkersFor) != 1 || aws.listWorkersFor[0] != "mars" || len(gcp.listWorkersFor) != 1 || gcp.listWorkersFor[0] != "ci/gcp/europe" {
			t.Fatalf("schedulers saw the wrong names: aws=%v gcp=%v", aws.listWorkersFor, gcp.listWorkersFor)
		}
	})

	t.Run("UnscopedRequestsNeedSelection", func(t *testing.T) {
		call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, codes.NotFound)
		call(t, p.URL, "ListPlatformQueues", &emptypb.Empty{}, &bqs.ListPlatformQueuesResponse{}, codes.InvalidArgument)
		call(t, p.URL, "ListPlatformQueues", &emptypb.Empty{}, &bqs.ListPlatformQueuesResponse{}, codes.InvalidArgument, selecting("ci/azure"))
	})

	t.Run("OneKillAuthorizerSeesPortalNames", func(t *testing.T) {
		permission(t, p.URL, false, selecting("ci/aws"))
		permission(t, p.URL, true, selecting("ci/gcp"))
		kill := &bqs.KillOperationsRequest{Filter: &bqs.KillOperationsRequest_Filter{Type: &bqs.KillOperationsRequest_Filter_OperationName{OperationName: "same-operation"}}}
		call(t, p.URL, "KillOperations", kill, &emptypb.Empty{}, codes.PermissionDenied, selecting("ci/aws"))
		call(t, p.URL, "KillOperations", kill, &emptypb.Empty{}, codes.OK, selecting("ci/gcp"))
		aws.mu.Lock()
		awsKilled := len(aws.killed)
		aws.mu.Unlock()
		gcp.mu.Lock()
		gcpKilled := append([]string(nil), gcp.killed...)
		gcp.mu.Unlock()
		if awsKilled != 0 || len(gcpKilled) != 1 || gcpKilled[0] != "same-operation" {
			t.Fatalf("crossed cancellation: aws=%d gcp=%v", awsKilled, gcpKilled)
		}
	})

	t.Run("BrowserCancellationReachesScheduler", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		requestDone := make(chan error, 1)
		go func() {
			_, err := grpcWeb(ctx, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "slow"}, &bqs.GetOperationResponse{}, selecting("ci/aws"))
			requestDone <- err
		}()
		select {
		case <-aws.started:
		case <-time.After(3 * time.Second):
			t.Fatal("slow upstream did not start")
		}
		// gcp must remain usable while aws's upstream request is blocked.
		call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, codes.OK, selecting("ci/gcp"))
		cancel()
		if err := <-requestDone; err == nil {
			t.Fatal("slow request did not cancel")
		}
		select {
		case <-aws.canceled:
		case <-time.After(3 * time.Second):
			t.Fatal("browser cancellation did not reach scheduler")
		}
	})
}

// The production HTTP authentication layer places authentication metadata in
// the context before scheduler dispatch. This fixture uses a test-only trusted
// identity header to verify that the router preserves that context.
func TestAuthenticationContextReachesAuthorizers(t *testing.T) {
	_, address := startScheduler(t, "authenticated", "ci")
	policy := jmespath("authenticationMetadata.public.principal == 'alice'")
	mux := http.NewServeMux()
	if err := schedulerservice.NewSchedulerService(&bb_portal.SchedulerService{
		Schedulers:               map[string]*builder.SchedulerConfiguration{"ci": schedulerAt(address, "ci")},
		ReadAuthorizer:           policy,
		KillOperationsAuthorizer: policy,
		ListOperationsPageSize:   2,
	}, nil, nil, clientFactory{t}, mux); err != nil {
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
			call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &bqs.GetOperationResponse{}, tc.code, headers)
			permission(t, p.URL, tc.allowed, headers)
		})
	}
}

func TestKillAuthorizationUsesFullOperationInstanceName(t *testing.T) {
	aws, address := startScheduler(t, "restricted", "mars")
	p := portal(t, &bb_portal.SchedulerService{
		Schedulers: map[string]*builder.SchedulerConfiguration{
			"ci/aws": schedulerAt(address, ""),
		},
		ReadAuthorizer:           allow(),
		KillOperationsAuthorizer: jmespath("instanceName == 'ci/aws/mars'"),
		ListOperationsPageSize:   2,
	})

	var response bqs.GetOperationResponse
	call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "restricted-operation"}, &response, codes.OK)
	if response.GetOperation().GetInstanceNameSuffix() != "restricted" {
		t.Fatalf("operation suffix was not returned: %v", &response)
	}
	permissionFor(t, p.URL, "restricted-operation", false, nil)
	kill := &bqs.KillOperationsRequest{Filter: &bqs.KillOperationsRequest_Filter{Type: &bqs.KillOperationsRequest_Filter_OperationName{OperationName: "restricted-operation"}}}
	call(t, p.URL, "KillOperations", kill, &emptypb.Empty{}, codes.PermissionDenied)
	aws.mu.Lock()
	defer aws.mu.Unlock()
	if len(aws.killed) != 0 {
		t.Fatalf("denied operation was killed: %v", aws.killed)
	}
}

// Most deployments have one scheduler. Neither form needs a selection.
func TestSingleSchedulerNeedsNoSelection(t *testing.T) {
	for _, tc := range []struct {
		name   string
		config func(address string) *bb_portal.SchedulerService
		prefix string
	}{
		{"BuildQueueStateClient", func(address string) *bb_portal.SchedulerService {
			return &bb_portal.SchedulerService{BuildQueueStateClient: &grpccfg.ClientConfiguration{Address: address}}
		}, ""},
		{"OneSchedulerEntry", func(address string) *bb_portal.SchedulerService {
			return &bb_portal.SchedulerService{Schedulers: map[string]*builder.SchedulerConfiguration{"ci": schedulerAt(address, "ci")}}
		}, "ci"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, address := startScheduler(t, "single-scheduler", "ci")
			c := tc.config(address)
			c.ReadAuthorizer, c.KillOperationsAuthorizer, c.ListOperationsPageSize = allow(), allow(), 2
			p := portal(t, c)

			r, err := http.Get(p.URL + "/api/v1/schedulers")
			if err != nil {
				t.Fatal(err)
			}
			defer r.Body.Close()
			var catalog struct {
				Schedulers []struct {
					InstanceNamePrefix string `json:"instanceNamePrefix"`
				} `json:"schedulers"`
			}
			if err := json.NewDecoder(r.Body).Decode(&catalog); err != nil {
				t.Fatal(err)
			}
			if len(catalog.Schedulers) != 1 || catalog.Schedulers[0].InstanceNamePrefix != tc.prefix {
				t.Fatalf("unexpected catalog: %+v", catalog)
			}
			var response bqs.GetOperationResponse
			call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &response, codes.OK)
			if response.GetOperation().GetTargetId() != "single-scheduler" {
				t.Fatalf("wrong operation: %v", &response)
			}
			permission(t, p.URL, true, nil)
		})
	}
}

// The empty prefix is a valid key. The browser then sends an empty header
// value, which must survive both HTTP-to-gRPC metadata conversions.
func TestEmptyPrefixSelectionThroughHTTP(t *testing.T) {
	_, rootAddress := startScheduler(t, "root", "")
	_, ciAddress := startScheduler(t, "ci", "ci")
	p := portal(t, &bb_portal.SchedulerService{
		Schedulers: map[string]*builder.SchedulerConfiguration{
			"":   schedulerAt(rootAddress, ""),
			"ci": schedulerAt(ciAddress, "ci"),
		},
		ReadAuthorizer:           allow(),
		KillOperationsAuthorizer: allow(),
		ListOperationsPageSize:   2,
	})
	for _, tc := range []struct{ selected, marker string }{{"", "root"}, {"ci", "ci"}} {
		var response bqs.GetOperationResponse
		call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &response, codes.OK, selecting(tc.selected))
		if response.GetOperation().GetTargetId() != tc.marker {
			t.Fatalf("selection %q reached %q", tc.selected, response.GetOperation().GetTargetId())
		}
	}
	permission(t, p.URL, true, selecting(""))
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
	_, healthyAddress := startScheduler(t, "healthy", "ci/healthy")
	// One scheduler is already unreachable when Portal initializes.
	p := portal(t, &bb_portal.SchedulerService{
		Schedulers: map[string]*builder.SchedulerConfiguration{
			"ci/down":    schedulerAt(unreachableAddress, "ci/down"),
			"ci/healthy": schedulerAt(healthyAddress, "ci/healthy"),
		},
		ReadAuthorizer:           allow(),
		KillOperationsAuthorizer: allow(),
		ListOperationsPageSize:   2,
	})
	call(t, p.URL, "ListPlatformQueues", &emptypb.Empty{}, &bqs.ListPlatformQueuesResponse{}, codes.Unavailable, selecting("ci/down"))
	var response bqs.GetOperationResponse
	call(t, p.URL, "GetOperation", &bqs.GetOperationRequest{OperationName: "same-operation"}, &response, codes.OK, selecting("ci/healthy"))
	if response.GetOperation().GetTargetId() != "healthy" {
		t.Fatalf("unreachable scheduler affected healthy response: %v", &response)
	}
}
