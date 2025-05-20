package fsacproxy

import (
	"context"

	"github.com/buildbarn/bb-portal/internal/api/common"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/proto/fsac"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// FsacServerImpl is a gRPC server that forwards requests to an FileSystemAccessCacheClient.
type FsacServerImpl struct {
	client     fsac.FileSystemAccessCacheClient
	authorizer auth.Authorizer
}

// NewFsacServerImpl creates a new FsacServerImpl from a given client.
func NewFsacServerImpl(client fsac.FileSystemAccessCacheClient, authorizer auth.Authorizer) *FsacServerImpl {
	return &FsacServerImpl{client: client, authorizer: authorizer}
}

// GetFileSystemAccessProfile proxies GetFileSystemAccessProfile requests to the client.
func (s *FsacServerImpl) GetFileSystemAccessProfile(ctx context.Context, req *fsac.GetFileSystemAccessProfileRequest) (*fsac.FileSystemAccessProfile, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	if !common.IsInstanceNameAllowed(ctx, s.authorizer, req.InstanceName) {
		return nil, status.Errorf(codes.PermissionDenied, "Not authorized")
	}
	return s.client.GetFileSystemAccessProfile(ctx, req)
}

// UpdateFileSystemAccessProfile proxies UpdateFileSystemAccessProfile requests to the client.
func (s *FsacServerImpl) UpdateFileSystemAccessProfile(ctx context.Context, req *fsac.UpdateFileSystemAccessProfileRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}
