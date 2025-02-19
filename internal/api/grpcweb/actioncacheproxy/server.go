package actioncacheproxy

import (
	"context"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/api/grpcweb"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ActionCacheServerImpl is a gRPC server that forwards requests to an ActionCacheClient.
type ActionCacheServerImpl struct {
	client     remoteexecution.ActionCacheClient
	authorizer auth.Authorizer
}

// NewAcctionCacheServerImpl creates a new ActionCacheServerImpl from a given client.
func NewAcctionCacheServerImpl(client remoteexecution.ActionCacheClient, authorizer auth.Authorizer) *ActionCacheServerImpl {
	return &ActionCacheServerImpl{client: client, authorizer: authorizer}
}

// GetActionResult proxies GetActionResult requests to the client.
func (s *ActionCacheServerImpl) GetActionResult(ctx context.Context, req *remoteexecution.GetActionResultRequest) (*remoteexecution.ActionResult, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	if !grpcweb.IsInstanceNamePrefixAllowed(ctx, s.authorizer, req.InstanceName) {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}

	response, err := s.client.GetActionResult(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}

	return response, err
}

// UpdateActionResult proxies UpdateActionResult requests to the client.
func (s *ActionCacheServerImpl) UpdateActionResult(ctx context.Context, req *remoteexecution.UpdateActionResultRequest) (*remoteexecution.ActionResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}
