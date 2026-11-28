// Package buildqueuestateproxy contains BuildQueueState gRPC proxies.
package buildqueuestateproxy

import (
	"context"
	"strings"

	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SelectedInstanceNameMetadataKey carries the instance name selected by the
// web UI for requests that do not identify a platform queue themselves.
const SelectedInstanceNameMetadataKey = "x-bb-portal-scheduler-instance-name"

// Route is one BuildQueueState API reachable through Portal.
type Route struct {
	// Requests for instance names with this prefix use this route. When
	// several routes match, the longest prefix wins.
	MatchInstanceNamePrefix digest.InstanceName
	// Replaces MatchInstanceNamePrefix in requests forwarded to Client,
	// and is replaced by it again in responses.
	AddInstanceNamePrefix digest.InstanceName
	Client                buildqueuestate.BuildQueueStateClient
}

type resolvedRoute struct {
	Route
	index       int
	toBackend   digest.InstanceNamePatcher
	fromBackend digest.InstanceNamePatcher
}

type demultiplexingClient struct {
	routes []resolvedRoute
	trie   *digest.InstanceNameTrie
}

// NewDemultiplexingClient creates a BuildQueueStateClient that forwards each
// request to the route responsible for its instance name. Instance names in
// requests are rewritten into the scheduler's namespace, and instance names
// in responses back into Portal's, so callers and their authorizers only see
// Portal's names.
//
// Callers must not pass duplicate MatchInstanceNamePrefix values.
func NewDemultiplexingClient(routes []Route) buildqueuestate.BuildQueueStateClient {
	c := &demultiplexingClient{trie: digest.NewInstanceNameTrie()}
	for i, r := range routes {
		c.trie.Set(r.MatchInstanceNamePrefix, i)
		c.routes = append(c.routes, resolvedRoute{
			Route:       r,
			index:       i,
			toBackend:   digest.NewInstanceNamePatcher(r.MatchInstanceNamePrefix, r.AddInstanceNamePrefix),
			fromBackend: digest.NewInstanceNamePatcher(r.AddInstanceNamePrefix, r.MatchInstanceNamePrefix),
		})
	}
	return c
}

// resolve selects a route from a request's platform queue name, or from the
// selected scheduler metadata when the request has no platform queue name.
func (c *demultiplexingClient) resolve(ctx context.Context, platformQueueName *buildqueuestate.PlatformQueueName) (*resolvedRoute, error) {
	if platformQueueName != nil {
		return c.resolveInstanceName(platformQueueName.InstanceNamePrefix)
	}

	values := metadata.ValueFromIncomingContext(ctx, SelectedInstanceNameMetadataKey)
	if len(values) > 0 {
		for _, value := range values[1:] {
			if value != values[0] {
				return nil, status.Errorf(codes.InvalidArgument, "Conflicting values for metadata key %#v", SelectedInstanceNameMetadataKey)
			}
		}
		return c.resolveInstanceName(values[0])
	}
	if len(c.routes) == 1 {
		return &c.routes[0], nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "Select a scheduler: more than one is configured, so metadata key %#v is required", SelectedInstanceNameMetadataKey)
}

func (c *demultiplexingClient) resolveInstanceName(instanceName string) (*resolvedRoute, error) {
	name, err := digest.NewInstanceName(instanceName)
	if err != nil {
		return nil, util.StatusWrapWithCode(err, codes.InvalidArgument, "Invalid instance name")
	}
	i := c.trie.GetLongestPrefix(name)
	if i < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "No scheduler is configured for instance name %#v", name.String())
	}
	return &c.routes[i], nil
}

func extractListWorkersPlatformQueueName(req *buildqueuestate.ListWorkersRequest) *buildqueuestate.PlatformQueueName {
	filter := req.GetFilter()
	if queue := filter.GetAll(); queue != nil {
		return queue.GetPlatformQueueName()
	}
	if invocation := filter.GetExecuting(); invocation != nil {
		return invocation.GetSizeClassQueueName().GetPlatformQueueName()
	}
	if invocation := filter.GetIdleSynchronizing(); invocation != nil {
		return invocation.GetSizeClassQueueName().GetPlatformQueueName()
	}
	return nil
}

func rewritePlatformQueueNameToBackend(r *resolvedRoute, name *buildqueuestate.PlatformQueueName) error {
	rewritten, err := rewrite(name.InstanceNamePrefix, r.MatchInstanceNamePrefix, r.toBackend)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Request mixes instance names of different schedulers: %#v is not within %#v", name.InstanceNamePrefix, r.MatchInstanceNamePrefix.String())
	}
	name.InstanceNamePrefix = rewritten.String()
	return nil
}

func (c *demultiplexingClient) rewritePlatformQueueNameFromBackend(r *resolvedRoute, name *buildqueuestate.PlatformQueueName) bool {
	rewritten, err := rewrite(name.InstanceNamePrefix, r.AddInstanceNamePrefix, r.fromBackend)
	if err != nil || c.trie.GetLongestPrefix(rewritten) != r.index {
		return false
	}
	name.InstanceNamePrefix = rewritten.String()
	return true
}

func (c *demultiplexingClient) rewriteOperationStateFromBackend(r *resolvedRoute, op *buildqueuestate.OperationState, requirePlatformQueueName bool) bool {
	if op == nil {
		return !requirePlatformQueueName
	}
	name := op.GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName()
	if name == nil {
		if requirePlatformQueueName {
			return false
		}
	} else if !c.rewritePlatformQueueNameFromBackend(r, name) {
		return false
	}
	if currentOperation := op.GetExecuting().GetCurrentOperation(); currentOperation != nil && !c.rewriteOperationStateFromBackend(r, currentOperation, false) {
		return false
	}
	return true
}

func (c *demultiplexingClient) rewriteWorkerStateFromBackend(r *resolvedRoute, worker *buildqueuestate.WorkerState) bool {
	currentOperation := worker.GetCurrentOperation()
	if currentOperation == nil || currentOperation.GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName() == nil {
		return true
	}
	return c.rewriteOperationStateFromBackend(r, currentOperation, false)
}

func rewrite(instanceName string, from digest.InstanceName, patcher digest.InstanceNamePatcher) (digest.InstanceName, error) {
	name, err := digest.NewInstanceName(instanceName)
	if err != nil {
		return digest.InstanceName{}, err
	}
	// InstanceNamePatcher requires the prefix to be present.
	if prefix := from.String(); prefix != "" && instanceName != prefix && !strings.HasPrefix(instanceName, prefix+"/") {
		return digest.InstanceName{}, status.Error(codes.OutOfRange, "Instance name outside prefix")
	}
	return patcher.PatchInstanceName(name), nil
}

func (c *demultiplexingClient) instanceNameOutsideRouteError(r *resolvedRoute) error {
	return status.Errorf(codes.Internal, "Scheduler for instance name prefix %#v returned an instance name outside %#v", r.MatchInstanceNamePrefix.String(), r.AddInstanceNamePrefix.String())
}

func (c *demultiplexingClient) GetOperation(ctx context.Context, in *buildqueuestate.GetOperationRequest, opts ...grpc.CallOption) (*buildqueuestate.GetOperationResponse, error) {
	r, err := c.resolve(ctx, nil)
	if err != nil {
		return nil, err
	}
	response, err := r.Client.GetOperation(ctx, proto.CloneOf(in), opts...)
	if err != nil {
		return nil, err
	}
	if !c.rewriteOperationStateFromBackend(r, response.GetOperation(), true) {
		// The operation exists, but belongs to another route or has no attribution.
		return nil, status.Error(codes.NotFound, "Operation was not found")
	}
	return response, nil
}

// ListOperations drops operations outside the route. Callers page until they
// receive a short page, so pages shortened by dropping are refilled from the
// scheduler before returning.
func (c *demultiplexingClient) ListOperations(ctx context.Context, in *buildqueuestate.ListOperationsRequest, opts ...grpc.CallOption) (*buildqueuestate.ListOperationsResponse, error) {
	r, err := c.resolve(ctx, nil)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	response := &buildqueuestate.ListOperationsResponse{}
	for {
		page, err := r.Client.ListOperations(ctx, backendRequest, opts...)
		if err != nil {
			return nil, err
		}
		response.PaginationInfo = page.PaginationInfo
		for _, operation := range page.Operations {
			if c.rewriteOperationStateFromBackend(r, operation, true) {
				response.Operations = append(response.Operations, operation)
			}
		}
		if in.PageSize == 0 || uint32(len(page.Operations)) < in.PageSize || uint32(len(response.Operations)) >= in.PageSize {
			if in.PageSize != 0 && uint32(len(response.Operations)) > in.PageSize {
				response.Operations = response.Operations[:in.PageSize]
			}
			return response, nil
		}
		lastName := page.Operations[len(page.Operations)-1].Name
		if backendRequest.StartAfter != nil && lastName <= backendRequest.StartAfter.OperationName {
			return nil, status.Error(codes.Internal, "Scheduler pagination cursor did not advance")
		}
		backendRequest.StartAfter = &buildqueuestate.ListOperationsRequest_StartAfter{OperationName: lastName}
	}
}

func (c *demultiplexingClient) KillOperations(ctx context.Context, in *buildqueuestate.KillOperationsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	var platformQueueName *buildqueuestate.PlatformQueueName
	if queue, ok := in.GetFilter().GetType().(*buildqueuestate.KillOperationsRequest_Filter_SizeClassQueueWithoutWorkers); ok {
		platformQueueName = queue.SizeClassQueueWithoutWorkers.GetPlatformQueueName()
	}
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if queue, ok := backendRequest.GetFilter().GetType().(*buildqueuestate.KillOperationsRequest_Filter_SizeClassQueueWithoutWorkers); ok && queue.SizeClassQueueWithoutWorkers.GetPlatformQueueName() != nil {
		if err := rewritePlatformQueueNameToBackend(r, queue.SizeClassQueueWithoutWorkers.PlatformQueueName); err != nil {
			return nil, err
		}
	}
	return r.Client.KillOperations(ctx, backendRequest, opts...)
}

// ListPlatformQueues drops queues outside the route or without attribution.
func (c *demultiplexingClient) ListPlatformQueues(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*buildqueuestate.ListPlatformQueuesResponse, error) {
	r, err := c.resolve(ctx, nil)
	if err != nil {
		return nil, err
	}
	response, err := r.Client.ListPlatformQueues(ctx, proto.CloneOf(in), opts...)
	if err != nil {
		return nil, err
	}
	queues := response.PlatformQueues[:0]
	for _, queue := range response.PlatformQueues {
		name := queue.GetName()
		if name != nil && c.rewritePlatformQueueNameFromBackend(r, name) {
			queues = append(queues, queue)
		}
	}
	response.PlatformQueues = queues
	return response, nil
}

func (c *demultiplexingClient) ListInvocationChildren(ctx context.Context, in *buildqueuestate.ListInvocationChildrenRequest, opts ...grpc.CallOption) (*buildqueuestate.ListInvocationChildrenResponse, error) {
	platformQueueName := in.GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName()
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if platformQueueName != nil {
		name := backendRequest.GetInvocationName().GetSizeClassQueueName().PlatformQueueName
		if err := rewritePlatformQueueNameToBackend(r, name); err != nil {
			return nil, err
		}
	}
	return r.Client.ListInvocationChildren(ctx, backendRequest, opts...)
}

func (c *demultiplexingClient) ListQueuedOperations(ctx context.Context, in *buildqueuestate.ListQueuedOperationsRequest, opts ...grpc.CallOption) (*buildqueuestate.ListQueuedOperationsResponse, error) {
	platformQueueName := in.GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName()
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if platformQueueName != nil {
		name := backendRequest.GetInvocationName().GetSizeClassQueueName().PlatformQueueName
		if err := rewritePlatformQueueNameToBackend(r, name); err != nil {
			return nil, err
		}
	}
	response, err := r.Client.ListQueuedOperations(ctx, backendRequest, opts...)
	if err != nil {
		return nil, err
	}
	for _, operation := range response.QueuedOperations {
		if !c.rewriteOperationStateFromBackend(r, operation, false) {
			return nil, c.instanceNameOutsideRouteError(r)
		}
	}
	return response, nil
}

func (c *demultiplexingClient) ListWorkers(ctx context.Context, in *buildqueuestate.ListWorkersRequest, opts ...grpc.CallOption) (*buildqueuestate.ListWorkersResponse, error) {
	platformQueueName := extractListWorkersPlatformQueueName(in)
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if platformQueueName != nil {
		if err := rewritePlatformQueueNameToBackend(r, extractListWorkersPlatformQueueName(backendRequest)); err != nil {
			return nil, err
		}
	}
	response, err := r.Client.ListWorkers(ctx, backendRequest, opts...)
	if err != nil {
		return nil, err
	}
	for _, worker := range response.Workers {
		if !c.rewriteWorkerStateFromBackend(r, worker) {
			return nil, c.instanceNameOutsideRouteError(r)
		}
	}
	return response, nil
}

func (c *demultiplexingClient) TerminateWorkers(ctx context.Context, in *buildqueuestate.TerminateWorkersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	r, err := c.resolve(ctx, nil)
	if err != nil {
		return nil, err
	}
	return r.Client.TerminateWorkers(ctx, proto.CloneOf(in), opts...)
}

func (c *demultiplexingClient) ListDrains(ctx context.Context, in *buildqueuestate.ListDrainsRequest, opts ...grpc.CallOption) (*buildqueuestate.ListDrainsResponse, error) {
	platformQueueName := in.GetSizeClassQueueName().GetPlatformQueueName()
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if platformQueueName != nil {
		if err := rewritePlatformQueueNameToBackend(r, backendRequest.GetSizeClassQueueName().PlatformQueueName); err != nil {
			return nil, err
		}
	}
	return r.Client.ListDrains(ctx, backendRequest, opts...)
}

func (c *demultiplexingClient) AddDrain(ctx context.Context, in *buildqueuestate.AddOrRemoveDrainRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	platformQueueName := in.GetSizeClassQueueName().GetPlatformQueueName()
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if platformQueueName != nil {
		if err := rewritePlatformQueueNameToBackend(r, backendRequest.GetSizeClassQueueName().PlatformQueueName); err != nil {
			return nil, err
		}
	}
	return r.Client.AddDrain(ctx, backendRequest, opts...)
}

func (c *demultiplexingClient) RemoveDrain(ctx context.Context, in *buildqueuestate.AddOrRemoveDrainRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	platformQueueName := in.GetSizeClassQueueName().GetPlatformQueueName()
	r, err := c.resolve(ctx, platformQueueName)
	if err != nil {
		return nil, err
	}
	backendRequest := proto.CloneOf(in)
	if platformQueueName != nil {
		if err := rewritePlatformQueueNameToBackend(r, backendRequest.GetSizeClassQueueName().PlatformQueueName); err != nil {
			return nil, err
		}
	}
	return r.Client.RemoveDrain(ctx, backendRequest, opts...)
}
