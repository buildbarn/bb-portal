package buildqueuestateproxy

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/auth"
	"github.com/jmespath/go-jmespath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestFilterPlatformQueues(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), auth.MustNewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
					},
				}),
			},
		}),
	}))

	t.Run("NoPlatformQueues", func(t *testing.T) {
		platformQueues := buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: []*buildqueuestate.PlatformQueueState{},
		}
		allowedQueues := filterPlatormQueues(ctx, &platformQueues, a)
		if len(allowedQueues) != 0 {
			t.Errorf("Expected no platform queues, got %d", len(allowedQueues))
		}
	})

	t.Run("FilterQueues", func(t *testing.T) {
		platformQueues := buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: []*buildqueuestate.PlatformQueueState{
				{
					Name: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
				{
					Name: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "forbidden",
					},
				},
			},
		}
		allowedQueues := filterPlatormQueues(ctx, &platformQueues, a)
		if len(allowedQueues) != 1 {
			t.Errorf("Expected one platform queue, got %d", len(allowedQueues))
		}
		expected := platformQueues.PlatformQueues[0]
		if allowedQueues[0] != expected {
			t.Errorf("Expected platform queue %+v, got %+v", expected, allowedQueues[0])
		}
	})

	t.Run("AllowEmptyInstanceNames", func(t *testing.T) {
		platformQueues := buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: []*buildqueuestate.PlatformQueueState{
				{
					Name: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "",
					},
				},
				{
					Name: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "forbidden",
					},
				},
			},
		}
		allowedQueues := filterPlatormQueues(ctx, &platformQueues, a)
		if len(allowedQueues) != 1 {
			t.Errorf("Expected one platform queue, got %d", len(allowedQueues))
		}
		expected := platformQueues.PlatformQueues[0]
		if allowedQueues[0] != expected {
			t.Errorf("Expected platform queue %+v, got %+v", expected, allowedQueues[0])
		}
	})

	t.Run("InvalidPlatformQueue", func(t *testing.T) {
		log.SetOutput(io.Discard)
		platformQueues := buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: []*buildqueuestate.PlatformQueueState{
				{
					Name: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "asdff//////DF////",
					},
				},
			},
		}
		allowedQueues := filterPlatormQueues(ctx, &platformQueues, a)
		if len(allowedQueues) != 0 {
			t.Errorf("Expected no platform queues, got %d", len(allowedQueues))
		}
	})
}

func TestFilterOperations(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), auth.MustNewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
					},
				}),
			},
		}),
	}))

	t.Run("NoOperations", func(t *testing.T) {
		operations := buildqueuestate.ListOperationsResponse{
			Operations: []*buildqueuestate.OperationState{},
		}
		allowedOperations := filterOperations(ctx, &operations, a)
		if len(allowedOperations) != 0 {
			t.Errorf("Expected no operations, got %d", len(allowedOperations))
		}
	})

	t.Run("FilterOperations", func(t *testing.T) {
		operations := buildqueuestate.ListOperationsResponse{
			Operations: []*buildqueuestate.OperationState{
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "allowed",
							},
						},
					},
				},
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "forbidden",
							},
						},
					},
				},
			},
		}
		allowedOperations := filterOperations(ctx, &operations, a)
		if len(allowedOperations) != 1 {
			t.Errorf("Expected one operation, got %d", len(allowedOperations))
		}
		expected := operations.Operations[0]
		if allowedOperations[0] != expected {
			t.Errorf("Expected operation %+v, got %+v", expected, allowedOperations[0])
		}
	})

	t.Run("AllowEmptyInstanceNames", func(t *testing.T) {
		operations := buildqueuestate.ListOperationsResponse{
			Operations: []*buildqueuestate.OperationState{
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "",
							},
						},
					},
				},
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "forbidden",
							},
						},
					},
				},
			},
		}
		allowedOperations := filterOperations(ctx, &operations, a)
		if len(allowedOperations) != 1 {
			t.Errorf("Expected one operation, got %d", len(allowedOperations))
		}
		expected := operations.Operations[0]
		if allowedOperations[0] != expected {
			t.Errorf("Expected operation %+v, got %+v", expected, allowedOperations[0])
		}
	})

	t.Run("InvalidOperation", func(t *testing.T) {
		log.SetOutput(io.Discard)
		operations := buildqueuestate.ListOperationsResponse{
			Operations: []*buildqueuestate.OperationState{
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "asdff//////DF////",
							},
						},
					},
				},
			},
		}
		allowedOperations := filterOperations(ctx, &operations, a)
		if len(allowedOperations) != 0 {
			t.Errorf("Expected no operations, got %d", len(allowedOperations))
		}
	})
}

func TestGetInstanceNamePrefixFromListWorkersRequest(t *testing.T) {
	t.Run("NoFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{}
		_, err := getInstanceNamePrefixFromListWorkersRequest(req)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument error, got %v", err)
		}
	})

	t.Run("AllFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{
			Filter: &buildqueuestate.ListWorkersRequest_Filter{
				Type: &buildqueuestate.ListWorkersRequest_Filter_All{
					All: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "all",
						},
					},
				},
			},
		}
		instanceNamePrefix, err := getInstanceNamePrefixFromListWorkersRequest(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if instanceNamePrefix != "all" {
			t.Errorf("Expected instance name prefix 'all', got %s", instanceNamePrefix)
		}
	})

	t.Run("ExecutingFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{
			Filter: &buildqueuestate.ListWorkersRequest_Filter{
				Type: &buildqueuestate.ListWorkersRequest_Filter_Executing{
					Executing: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "executing",
							},
						},
					},
				},
			},
		}
		instanceNamePrefix, err := getInstanceNamePrefixFromListWorkersRequest(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if instanceNamePrefix != "executing" {
			t.Errorf("Expected instance name prefix 'executing', got %s", instanceNamePrefix)
		}
	})

	t.Run("IdleSynchronizingFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{
			Filter: &buildqueuestate.ListWorkersRequest_Filter{
				Type: &buildqueuestate.ListWorkersRequest_Filter_IdleSynchronizing{
					IdleSynchronizing: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "idle",
							},
						},
					},
				},
			},
		}
		instanceNamePrefix, err := getInstanceNamePrefixFromListWorkersRequest(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if instanceNamePrefix != "idle" {
			t.Errorf("Expected instance name prefix 'idle', got %s", instanceNamePrefix)
		}
	})

	t.Run("InvalidFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{
			Filter: &buildqueuestate.ListWorkersRequest_Filter{},
		}
		_, err := getInstanceNamePrefixFromListWorkersRequest(req)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument error, got %v", err)
		}
	})
}
