package isccproxy

import (
	"context"

	"github.com/buildbarn/bb-portal/internal/api/grpcweb"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/proto/iscc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// IsccServerImpl is a gRPC server that forwards requests to an InitialSizeClassCacheClient.
type IsccServerImpl struct {
	client     iscc.InitialSizeClassCacheClient
	authorizer auth.Authorizer
}

// NewIsccServerImpl creates a new IsccServerImpl from a given client.
func NewIsccServerImpl(client iscc.InitialSizeClassCacheClient, authorizer auth.Authorizer) *IsccServerImpl {
	return &IsccServerImpl{client: client, authorizer: authorizer}
}

// GetPreviousExecutionStats proxies GetPreviousExecutionStats requests to the client.
func (s *IsccServerImpl) GetPreviousExecutionStats(ctx context.Context, req *iscc.GetPreviousExecutionStatsRequest) (*iscc.PreviousExecutionStats, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	if !grpcweb.IsInstanceNameAllowed(ctx, s.authorizer, req.InstanceName) {
		return nil, status.Errorf(codes.PermissionDenied, "Not authorized")
	}
	return s.client.GetPreviousExecutionStats(ctx, req)
}

// UpdatePreviousExecutionStats proxies UpdatePreviousExecutionStats requests to the client.
func (s *IsccServerImpl) UpdatePreviousExecutionStats(ctx context.Context, req *iscc.UpdatePreviousExecutionStatsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}
