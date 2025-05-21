package buildqueuestateproxy

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/stretchr/testify/require"

	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/auth"
	"github.com/jmespath/go-jmespath"
	"go.uber.org/mock/gomock"
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
		operations := []*buildqueuestate.OperationState{}
		allowedOperations := filterOperations(ctx, operations, a)
		if len(allowedOperations) != 0 {
			t.Errorf("Expected no operations, got %d", len(allowedOperations))
		}
	})

	t.Run("FilterOperations", func(t *testing.T) {
		operations := []*buildqueuestate.OperationState{
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
		}
		allowedOperations := filterOperations(ctx, operations, a)
		if len(allowedOperations) != 1 {
			t.Errorf("Expected one operation, got %d", len(allowedOperations))
		}
		expected := operations[0]
		if allowedOperations[0] != expected {
			t.Errorf("Expected operation %+v, got %+v", expected, allowedOperations[0])
		}
	})

	t.Run("AllowEmptyInstanceNames", func(t *testing.T) {
		operations := []*buildqueuestate.OperationState{
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
		}
		allowedOperations := filterOperations(ctx, operations, a)
		if len(allowedOperations) != 1 {
			t.Errorf("Expected one operation, got %d", len(allowedOperations))
		}
		expected := operations[0]
		if allowedOperations[0] != expected {
			t.Errorf("Expected operation %+v, got %+v", expected, allowedOperations[0])
		}
	})

	t.Run("InvalidOperation", func(t *testing.T) {
		log.SetOutput(io.Discard)
		operations := []*buildqueuestate.OperationState{
			{
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "asdff//////DF////",
						},
					},
				},
			},
		}
		allowedOperations := filterOperations(ctx, operations, a)
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

func TestGetPaginationInfo(t *testing.T) {
	t.Run("NoEntries", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(0, 10, func(i int) bool { return true })
		if paginationInfo.StartIndex != 0 {
			t.Errorf("Expected start index 0, got %d", paginationInfo.StartIndex)
		}
		if paginationInfo.TotalEntries != 0 {
			t.Errorf("Expected total entries 0, got %d", paginationInfo.TotalEntries)
		}
		if endIndex != 0 {
			t.Errorf("Expected end index 0, got %d", endIndex)
		}
	})

	t.Run("LessThanPageSize", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(5, 10, func(i int) bool { return true })
		if paginationInfo.StartIndex != 0 {
			t.Errorf("Expected start index 0, got %d", paginationInfo.StartIndex)
		}
		if paginationInfo.TotalEntries != 5 {
			t.Errorf("Expected total entries 5, got %d", paginationInfo.TotalEntries)
		}
		if endIndex != 5 {
			t.Errorf("Expected end index 5, got %d", endIndex)
		}
	})

	t.Run("ExactPageSize", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(10, 10, func(i int) bool { return true })
		if paginationInfo.StartIndex != 0 {
			t.Errorf("Expected start index 0, got %d", paginationInfo.StartIndex)
		}
		if paginationInfo.TotalEntries != 10 {
			t.Errorf("Expected total entries 10, got %d", paginationInfo.TotalEntries)
		}
		if endIndex != 10 {
			t.Errorf("Expected end index 10, got %d", endIndex)
		}
	})

	t.Run("MoreThanPageSize", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(15, 10, func(i int) bool { return true })
		if paginationInfo.StartIndex != 0 {
			t.Errorf("Expected start index 0, got %d", paginationInfo.StartIndex)
		}
		if paginationInfo.TotalEntries != 15 {
			t.Errorf("Expected total entries 15, got %d", paginationInfo.TotalEntries)
		}
		if endIndex != 10 {
			t.Errorf("Expected end index 10, got %d", endIndex)
		}
	})

	t.Run("StartAfterMiddle", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(20, 5, func(i int) bool { return i >= 10 })
		if paginationInfo.StartIndex != 10 {
			t.Errorf("Expected start index 10, got %d", paginationInfo.StartIndex)
		}
		if paginationInfo.TotalEntries != 20 {
			t.Errorf("Expected total entries 20, got %d", paginationInfo.TotalEntries)
		}
		if endIndex != 15 {
			t.Errorf("Expected end index 15, got %d", endIndex)
		}
	})

	t.Run("StartCloseToEnd", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(20, 5, func(i int) bool { return i >= 18 })
		if paginationInfo.StartIndex != 18 {
			t.Errorf("Expected start index 18, got %d", paginationInfo.StartIndex)
		}
		if paginationInfo.TotalEntries != 20 {
			t.Errorf("Expected total entries 20, got %d", paginationInfo.TotalEntries)
		}
		if endIndex != 20 {
			t.Errorf("Expected end index 20, got %d", endIndex)
		}
	})
}

func TestCreatePaginatedListOperationsResponse(t *testing.T) {
	t.Run("NoOperations", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{}
		response := createPaginatedListOperationsResponse(allOperations, 10, nil)
		if len(response.Operations) != 0 {
			t.Errorf("Expected no operations, got %d", len(response.Operations))
		}
		if response.PaginationInfo.TotalEntries != 0 {
			t.Errorf("Expected total entries 0, got %d", response.PaginationInfo.TotalEntries)
		}
	})

	t.Run("LessThanPageSize", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
		}
		response := createPaginatedListOperationsResponse(allOperations, 10, nil)
		if len(response.Operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(response.Operations))
		}
		if response.PaginationInfo.TotalEntries != 2 {
			t.Errorf("Expected total entries 2, got %d", response.PaginationInfo.TotalEntries)
		}
	})

	t.Run("ExactPageSize", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
			{Name: "op3"},
		}
		response := createPaginatedListOperationsResponse(allOperations, 3, nil)
		if len(response.Operations) != 3 {
			t.Errorf("Expected 3 operations, got %d", len(response.Operations))
		}
		if response.PaginationInfo.TotalEntries != 3 {
			t.Errorf("Expected total entries 3, got %d", response.PaginationInfo.TotalEntries)
		}
	})

	t.Run("MoreThanPageSize", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
			{Name: "op3"},
			{Name: "op4"},
		}
		response := createPaginatedListOperationsResponse(allOperations, 2, nil)
		if len(response.Operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(response.Operations))
		}
		if response.PaginationInfo.TotalEntries != 4 {
			t.Errorf("Expected total entries 4, got %d", response.PaginationInfo.TotalEntries)
		}
	})

	t.Run("StartAfterMiddle", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
			{Name: "op3"},
			{Name: "op4"},
		}
		startAfter := &buildqueuestate.ListOperationsRequest_StartAfter{OperationName: "op2"}
		response := createPaginatedListOperationsResponse(allOperations, 2, startAfter)
		if len(response.Operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(response.Operations))
		}
		if response.Operations[0].Name != "op3" {
			t.Errorf("Expected operation 'op3', got %s", response.Operations[0].Name)
		}
		if response.Operations[1].Name != "op4" {
			t.Errorf("Expected operation 'op4', got %s", response.Operations[1].Name)
		}
		if response.PaginationInfo.StartIndex != 2 {
			t.Errorf("Expected start index 2, got %d", response.PaginationInfo.StartIndex)
		}
		if response.PaginationInfo.TotalEntries != 4 {
			t.Errorf("Expected total entries 4, got %d", response.PaginationInfo.TotalEntries)
		}
	})

	t.Run("StartCloseToEnd", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
			{Name: "op3"},
			{Name: "op4"},
		}
		startAfter := &buildqueuestate.ListOperationsRequest_StartAfter{OperationName: "op3"}
		response := createPaginatedListOperationsResponse(allOperations, 2, startAfter)
		if len(response.Operations) != 1 {
			t.Errorf("Expected 1 operation, got %d", len(response.Operations))
		}
		if response.Operations[0].Name != "op4" {
			t.Errorf("Expected operation 'op4', got %s", response.Operations[0].Name)
		}
		if response.PaginationInfo.StartIndex != 3 {
			t.Errorf("Expected start index 3, got %d", response.PaginationInfo.StartIndex)
		}
		if response.PaginationInfo.TotalEntries != 4 {
			t.Errorf("Expected total entries 4, got %d", response.PaginationInfo.TotalEntries)
		}
	})
}

func TestListOperations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(auth.NewContextWithAuthenticationMetadata(context.Background(), auth.MustNewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
					},
				}),
			},
		}),
	})), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	instanceNameAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	bqsServer := NewBuildQueueStateServerImpl(bqsClient, instanceNameAuthorizer, 2)

	operations := []*buildqueuestate.OperationState{
		{
			Name: "op1",
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
		},
		{
			Name: "op2",
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
		},
		{
			Name: "op3",
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
		},
		{
			Name: "op4",
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
		},
		{
			Name: "op5",
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
		},
	}

	t.Run("NoOperations", func(t *testing.T) {
		clientResponse := &buildqueuestate.ListOperationsResponse{
			Operations: []*buildqueuestate.OperationState{},
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 0,
			},
		}

		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(clientResponse, nil)

		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 5,
		})
		require.NoError(t, err)
		require.Equal(t, clientResponse, resp)
	})

	t.Run("FilterOperations", func(t *testing.T) {
		clientResponse := &buildqueuestate.ListOperationsResponse{
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
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 2,
			},
		}
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(clientResponse, nil)

		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 5,
		})
		require.NoError(t, err)
		require.Equal(t, clientResponse.Operations[0:1], resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(1), resp.PaginationInfo.TotalEntries)
	})

	t.Run("AllowEmptyInstanceNames", func(t *testing.T) {
		clientResponse := &buildqueuestate.ListOperationsResponse{
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
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 2,
			},
		}
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(clientResponse, nil)
		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 5,
		})
		require.NoError(t, err)
		require.Equal(t, clientResponse.Operations[0:1], resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(1), resp.PaginationInfo.TotalEntries)
	})

	t.Run("InvalidOperation", func(t *testing.T) {
		clientResponse := &buildqueuestate.ListOperationsResponse{
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
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 1,
			},
		}
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(clientResponse, nil)
		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 5,
		})
		require.NoError(t, err)
		require.Equal(t, []*buildqueuestate.OperationState{}, resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(0), resp.PaginationInfo.TotalEntries)
	})

	t.Run("ClientPaginationWithNumOperationsMultipleOfClientPageSize", func(t *testing.T) {
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[0:2],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 4,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[2:4],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   2,
				TotalEntries: 4,
			},
		}, nil)
		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 10,
		})
		require.NoError(t, err)
		require.Equal(t, operations[0:4], resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(4), resp.PaginationInfo.TotalEntries)
	})

	t.Run("ClientPaginationWithNumOperationsNotMultipleOfClientPageSize", func(t *testing.T) {
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[0:2],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 5,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[2:4],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   2,
				TotalEntries: 5,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[4:5],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   4,
				TotalEntries: 5,
			},
		}, nil)
		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 10,
		})
		require.NoError(t, err)
		require.Equal(t, operations[0:5], resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(5), resp.PaginationInfo.TotalEntries)
	})

	t.Run("ServerPaginationPageSize", func(t *testing.T) {
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[0:2],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 5,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[2:4],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   2,
				TotalEntries: 5,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[4:5],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   4,
				TotalEntries: 5,
			},
		}, nil)
		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 3,
		})
		require.NoError(t, err)
		require.Equal(t, operations[0:3], resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(5), resp.PaginationInfo.TotalEntries)
	})

	t.Run("ServerPaginationStartIndex", func(t *testing.T) {
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[0:2],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 5,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[2:4],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   2,
				TotalEntries: 5,
			},
		}, nil)
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: operations[4:5],
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   4,
				TotalEntries: 5,
			},
		}, nil)
		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 3,
			StartAfter: &buildqueuestate.ListOperationsRequest_StartAfter{
				OperationName: "op1",
			},
		})
		require.NoError(t, err)
		require.Equal(t, operations[1:4], resp.Operations)
		require.Equal(t, uint32(1), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(5), resp.PaginationInfo.TotalEntries)
	})
}
