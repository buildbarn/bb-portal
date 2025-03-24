package buildqueuestateproxy

import (
	"context"

	"github.com/buildbarn/bb-portal/internal/api/grpcweb"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// BuildQueueStateServerImpl is a gRPC server that forwards requests to a BuildQueueStateClient.
type BuildQueueStateServerImpl struct {
	client     buildqueuestate.BuildQueueStateClient
	authorizer auth.Authorizer
}

// NewBuildQueueStateServerImpl creates a new BuildQueueStateServerImpl from a
// given client. It also takes an authorizer to filter out the queues that the
// user is not allowed to see.
func NewBuildQueueStateServerImpl(client buildqueuestate.BuildQueueStateClient, authorizer auth.Authorizer) *BuildQueueStateServerImpl {
	return &BuildQueueStateServerImpl{client: client, authorizer: authorizer}
}

// GetOperation proxies GetOperation requests to the client.
func (s *BuildQueueStateServerImpl) GetOperation(ctx context.Context, req *buildqueuestate.GetOperationRequest) (*buildqueuestate.GetOperationResponse, error) {
	response, err := s.client.GetOperation(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Operation was not found")
	}

	platformQueueName := response.GetOperation().GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName()

	if platformQueueName == nil || !grpcweb.IsInstanceNameAllowed(ctx, s.authorizer, platformQueueName.InstanceNamePrefix) {
		return nil, status.Errorf(codes.NotFound, "Operation was not found")
	}

	return response, err
}

// ListOperations proxies ListOperations requests to the client.
func (s *BuildQueueStateServerImpl) ListOperations(ctx context.Context, req *buildqueuestate.ListOperationsRequest) (*buildqueuestate.ListOperationsResponse, error) {
	response, err := s.client.ListOperations(ctx, req)
	if err != nil {
		return nil, err
	}
	response.Operations = filterOperations(ctx, response, s.authorizer)
	// We modify the pagination info to not leak information about operations
	// that the user is not allowed to see.
	response.PaginationInfo.StartIndex = 0
	response.PaginationInfo.TotalEntries = 0
	return response, err
}

// KillOperations proxies KillOperations requests to the client.
func (s *BuildQueueStateServerImpl) KillOperations(ctx context.Context, req *buildqueuestate.KillOperationsRequest) (*emptypb.Empty, error) {
	// Check if the filter is of type OperationName.
	if filter, ok := req.GetFilter().GetType().(*buildqueuestate.KillOperationsRequest_Filter_OperationName); ok {
		// Calls GetOperation to check if the operation exists and the user is allowed to kill it.
		_, err := s.GetOperation(ctx, &buildqueuestate.GetOperationRequest{OperationName: filter.OperationName})
		if err != nil {
			return nil, err
		}

		return s.client.KillOperations(ctx, req)
	}

	return nil, status.Errorf(codes.InvalidArgument, "Can only kill operations by operation name")
}

// ListPlatformQueues proxies ListPlatformQueues requests to the client.
func (s *BuildQueueStateServerImpl) ListPlatformQueues(ctx context.Context, req *emptypb.Empty) (*buildqueuestate.ListPlatformQueuesResponse, error) {
	response, err := s.client.ListPlatformQueues(ctx, req)
	if err != nil {
		return nil, err
	}
	response.PlatformQueues = filterPlatormQueues(ctx, response, s.authorizer)
	return response, err
}

// ListInvocationChildren proxies ListInvocationChildren requests to the client.
func (s *BuildQueueStateServerImpl) ListInvocationChildren(ctx context.Context, req *buildqueuestate.ListInvocationChildrenRequest) (*buildqueuestate.ListInvocationChildrenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

// ListQueuedOperations proxies ListQueuedOperations requests to the client.
func (s *BuildQueueStateServerImpl) ListQueuedOperations(ctx context.Context, req *buildqueuestate.ListQueuedOperationsRequest) (*buildqueuestate.ListQueuedOperationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

// ListWorkers proxies ListWorkers requests to the client.
func (s *BuildQueueStateServerImpl) ListWorkers(ctx context.Context, req *buildqueuestate.ListWorkersRequest) (*buildqueuestate.ListWorkersResponse, error) {
	instanceNamePrefix, err := getInstanceNamePrefixFromListWorkersRequest(req)
	if err != nil {
		return nil, err
	}

	if !grpcweb.IsInstanceNameAllowed(ctx, s.authorizer, instanceNamePrefix) {
		return nil, status.Errorf(codes.PermissionDenied, "Not allowed to list workers for instance name prefix %s", instanceNamePrefix)
	}
	return s.client.ListWorkers(ctx, req)
}

// TerminateWorkers proxies TerminateWorkers requests to the client.
func (s *BuildQueueStateServerImpl) TerminateWorkers(ctx context.Context, req *buildqueuestate.TerminateWorkersRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

// ListDrains proxies ListDrains requests to the client.
func (s *BuildQueueStateServerImpl) ListDrains(ctx context.Context, req *buildqueuestate.ListDrainsRequest) (*buildqueuestate.ListDrainsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

// AddDrain proxies AddDrain requests to the client.
func (s *BuildQueueStateServerImpl) AddDrain(ctx context.Context, req *buildqueuestate.AddOrRemoveDrainRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

// RemoveDrain proxies RemoveDrain requests to the client.
func (s *BuildQueueStateServerImpl) RemoveDrain(ctx context.Context, req *buildqueuestate.AddOrRemoveDrainRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

func filterPlatormQueues(ctx context.Context, response *buildqueuestate.ListPlatformQueuesResponse, authorizer auth.Authorizer) []*buildqueuestate.PlatformQueueState {
	queues := response.GetPlatformQueues()
	// Filter out the queues that the user is not allowed to see.
	allowedQueues := make([]*buildqueuestate.PlatformQueueState, 0, len(queues))
	for _, queue := range queues {

		name := queue.GetName()

		if name != nil && grpcweb.IsInstanceNameAllowed(ctx, authorizer, name.InstanceNamePrefix) {
			allowedQueues = append(allowedQueues, queue)
		}
	}
	return allowedQueues
}

func filterOperations(ctx context.Context, response *buildqueuestate.ListOperationsResponse, authorizer auth.Authorizer) []*buildqueuestate.OperationState {
	operations := response.GetOperations()
	// Filter out the operations that the user is not allowed to see.
	allowedOperations := make([]*buildqueuestate.OperationState, 0, len(operations))
	for _, operation := range operations {

		platformQueueName := operation.GetInvocationName().GetSizeClassQueueName().GetPlatformQueueName()

		if platformQueueName != nil && grpcweb.IsInstanceNameAllowed(ctx, authorizer, platformQueueName.InstanceNamePrefix) {
			allowedOperations = append(allowedOperations, operation)
		}
	}
	return allowedOperations
}

func getInstanceNamePrefixFromListWorkersRequest(req *buildqueuestate.ListWorkersRequest) (string, error) {
	if platformQueueName := req.GetFilter().GetAll().GetPlatformQueueName(); platformQueueName != nil {
		return platformQueueName.InstanceNamePrefix, nil
	}
	if platformQueueName := req.GetFilter().GetExecuting().GetSizeClassQueueName().GetPlatformQueueName(); platformQueueName != nil {
		return platformQueueName.InstanceNamePrefix, nil
	}
	if platformQueueName := req.GetFilter().GetIdleSynchronizing().GetSizeClassQueueName().GetPlatformQueueName(); platformQueueName != nil {
		return platformQueueName.InstanceNamePrefix, nil
	}

	return "", status.Errorf(codes.InvalidArgument, "Request does not contain a valid InstanceNamePrefix")
}
