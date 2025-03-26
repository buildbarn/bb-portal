package casproxy

import (
	"context"
	"io"
	"strings"

	"github.com/buildbarn/bb-portal/internal/api/grpcweb"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"google.golang.org/genproto/googleapis/bytestream"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CasServerImpl is an gRPC server that forwards requests to a ByteStreamClient.
type CasServerImpl struct {
	client     bytestream.ByteStreamClient
	authorizer auth.Authorizer
}

// NewCasServerImpl creates a new CasServerImpl from a given client.
func NewCasServerImpl(client bytestream.ByteStreamClient, authorizer auth.Authorizer) *CasServerImpl {
	return &CasServerImpl{client: client, authorizer: authorizer}
}

// Read proxies Read requests to the client.
func (s *CasServerImpl) Read(req *bytestream.ReadRequest, stream bytestream.ByteStream_ReadServer) error {
	if req == nil {
		return status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	instanceName := getInstanceName(req.ResourceName)
	if !grpcweb.IsInstanceNameAllowed(stream.Context(), s.authorizer, instanceName) {
		return status.Errorf(codes.PermissionDenied, "Not authorized")
	}

	clientStream, err := s.client.Read(context.Background(), req)
	if err != nil {
		return err
	}
	for {
		resp, err := clientStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

// Write proxies Write requests to the client.
func (s *CasServerImpl) Write(stream bytestream.ByteStream_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "Action is not supported")
}

// QueryWriteStatus proxies QueryWriteStatus requests to the client.
func (s *CasServerImpl) QueryWriteStatus(ctx context.Context, req *bytestream.QueryWriteStatusRequest) (*bytestream.QueryWriteStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Action is not supported")
}

func getInstanceName(resourceName string) string {
	splitString := strings.Split(resourceName, "/blobs")[0]
	return strings.TrimPrefix(splitString, "/")
}
