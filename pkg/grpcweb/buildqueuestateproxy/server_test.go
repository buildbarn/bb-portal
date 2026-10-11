package buildqueuestateproxy

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/stretchr/testify/require"

	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/auth"
	"github.com/buildbarn/bb-storage/pkg/testutil"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFilterPlatformQueues(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
					},
				}),
			},
		}),
	})))

	t.Run("NoPlatformQueues", func(t *testing.T) {
		platformQueues := buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: []*buildqueuestate.PlatformQueueState{},
		}
		allowedQueues := filterPlatormQueues(ctx, &platformQueues, a)
		require.Len(t, allowedQueues, 0)
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
		require.Len(t, allowedQueues, 1)
		require.Equal(t, platformQueues.PlatformQueues[0], allowedQueues[0])
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
		require.Len(t, allowedQueues, 1)
		require.Equal(t, platformQueues.PlatformQueues[0], allowedQueues[0])
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
		require.Len(t, allowedQueues, 0)
	})
}

func TestFilterOperations(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
						structpb.NewStringValue("allowed/foo"),
					},
				}),
			},
		}),
	})))

	t.Run("NoOperations", func(t *testing.T) {
		operations := []*buildqueuestate.OperationState{}
		allowedOperations := filterOperations(ctx, operations, a)
		require.Len(t, allowedOperations, 0)
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
							InstanceNamePrefix: "allowed",
						},
					},
				},
				InstanceNameSuffix: "foo",
			},
			{
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "allowed",
						},
					},
				},
				InstanceNameSuffix: "bar",
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
			{
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "forbidden",
						},
					},
				},
				InstanceNameSuffix: "foo",
			},
		}
		allowedOperations := filterOperations(ctx, operations, a)
		require.Len(t, allowedOperations, 2)
		require.Equal(t, operations[0], allowedOperations[0])
		require.Equal(t, operations[1], allowedOperations[1])
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
		require.Len(t, allowedOperations, 1)
		require.Equal(t, operations[0], allowedOperations[0])
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
		require.Len(t, allowedOperations, 0)
	})
}

func TestGetInstanceNamePrefixFromListWorkersRequest(t *testing.T) {
	t.Run("NoFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{}
		_, err := getInstanceNamePrefixFromListWorkersRequest(req)
		require.ErrorContains(t, err, "Request does not contain a valid InstanceNamePrefix")
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
		require.NoError(t, err)
		require.Equal(t, "all", instanceNamePrefix)
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
		require.NoError(t, err)
		require.Equal(t, "executing", instanceNamePrefix)
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
		require.NoError(t, err)
		require.Equal(t, "idle", instanceNamePrefix)
	})

	t.Run("InvalidFilter", func(t *testing.T) {
		req := &buildqueuestate.ListWorkersRequest{
			Filter: &buildqueuestate.ListWorkersRequest_Filter{},
		}
		_, err := getInstanceNamePrefixFromListWorkersRequest(req)
		require.ErrorContains(t, err, "Request does not contain a valid InstanceNamePrefix")
	})
}

func TestIsOperationAllowed(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
						structpb.NewStringValue("allowed/foo"),
					},
				}),
			},
		}),
	})))

	t.Run("NoOperations", func(t *testing.T) {
		allowed := isOperationAllowed(ctx, a, nil)
		require.False(t, allowed)
	})

	t.Run("AllowedPrefix", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
		}
		allowed := isOperationAllowed(ctx, a, operation)
		require.True(t, allowed)
	})

	t.Run("AllowedSuffix", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{
			InstanceNameSuffix: "allowed",
		}
		allowed := isOperationAllowed(ctx, a, operation)
		require.True(t, allowed)
	})

	t.Run("AllowedFull", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
			InstanceNameSuffix: "foo",
		}
		allowed := isOperationAllowed(ctx, a, operation)
		require.True(t, allowed)
	})

	t.Run("DisallowedPrefix", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "disallowed",
					},
				},
			},
		}
		allowed := isOperationAllowed(ctx, a, operation)
		require.False(t, allowed)
	})

	t.Run("DisallowedSuffix", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{
			InstanceNameSuffix: "disallowed",
		}
		allowed := isOperationAllowed(ctx, a, operation)
		require.False(t, allowed)
	})

	t.Run("DisallowedFull", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
					PlatformQueueName: &buildqueuestate.PlatformQueueName{
						InstanceNamePrefix: "allowed",
					},
				},
			},
			InstanceNameSuffix: "bar",
		}
		allowed := isOperationAllowed(ctx, a, operation)
		require.False(t, allowed)
	})
}

func TestIsPrefixSuffixAllowed(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
						structpb.NewStringValue("allowed/foo"),
					},
				}),
			},
		}),
	})))

	t.Run("NoNames", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "", "")
		require.True(t, allowed)
	})

	t.Run("AllowedPrefix", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "allowed", "")
		require.True(t, allowed)
	})

	t.Run("AllowedSuffix", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "", "allowed")
		require.True(t, allowed)
	})

	t.Run("AllowedFull", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "allowed", "foo")
		require.True(t, allowed)
	})

	t.Run("DisallowedPrefix", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "disallowed", "")
		require.False(t, allowed)
	})

	t.Run("DisallowedSuffix", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "", "disallowed")
		require.False(t, allowed)
	})

	t.Run("DisallowedFull", func(t *testing.T) {
		allowed := isPrefixSuffixAllowed(ctx, a, "allowed", "bar")
		require.False(t, allowed)
	})
}

func TestGetPaginationInfo(t *testing.T) {
	t.Run("NoEntries", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(0, 10, func(i int) bool { return true })
		require.Equal(t, uint32(0), paginationInfo.StartIndex)
		require.Equal(t, uint32(0), paginationInfo.TotalEntries)
		require.Equal(t, 0, endIndex)
	})

	t.Run("LessThanPageSize", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(5, 10, func(i int) bool { return true })
		require.Equal(t, uint32(0), paginationInfo.StartIndex)
		require.Equal(t, uint32(5), paginationInfo.TotalEntries)
		require.Equal(t, 5, endIndex)
	})

	t.Run("ExactPageSize", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(10, 10, func(i int) bool { return true })
		require.Equal(t, uint32(0), paginationInfo.StartIndex)
		require.Equal(t, uint32(10), paginationInfo.TotalEntries)
		require.Equal(t, 10, endIndex)
	})

	t.Run("MoreThanPageSize", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(15, 10, func(i int) bool { return true })
		require.Equal(t, uint32(0), paginationInfo.StartIndex)
		require.Equal(t, uint32(15), paginationInfo.TotalEntries)
		require.Equal(t, 10, endIndex)
	})

	t.Run("StartAfterMiddle", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(20, 5, func(i int) bool { return i >= 10 })
		require.Equal(t, uint32(10), paginationInfo.StartIndex)
		require.Equal(t, uint32(20), paginationInfo.TotalEntries)
		require.Equal(t, 15, endIndex)
	})

	t.Run("StartCloseToEnd", func(t *testing.T) {
		paginationInfo, endIndex := getPaginationInfo(20, 5, func(i int) bool { return i >= 18 })
		require.Equal(t, uint32(18), paginationInfo.StartIndex)
		require.Equal(t, uint32(20), paginationInfo.TotalEntries)
		require.Equal(t, 20, endIndex)
	})
}

func TestCreatePaginatedListOperationsResponse(t *testing.T) {
	t.Run("NoOperations", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{}
		response := createPaginatedListOperationsResponse(allOperations, 10, nil)
		require.Len(t, response.Operations, 0)
		require.Equal(t, uint32(0), response.PaginationInfo.TotalEntries)
	})

	t.Run("LessThanPageSize", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
		}
		response := createPaginatedListOperationsResponse(allOperations, 10, nil)
		require.Len(t, response.Operations, 2)
		require.Equal(t, uint32(2), response.PaginationInfo.TotalEntries)
	})

	t.Run("ExactPageSize", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
			{Name: "op3"},
		}
		response := createPaginatedListOperationsResponse(allOperations, 3, nil)
		require.Len(t, response.Operations, 3)
		require.Equal(t, uint32(3), response.PaginationInfo.TotalEntries)
	})

	t.Run("MoreThanPageSize", func(t *testing.T) {
		allOperations := []*buildqueuestate.OperationState{
			{Name: "op1"},
			{Name: "op2"},
			{Name: "op3"},
			{Name: "op4"},
		}
		response := createPaginatedListOperationsResponse(allOperations, 2, nil)
		require.Len(t, response.Operations, 2)
		require.Equal(t, uint32(4), response.PaginationInfo.TotalEntries)
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
		require.Len(t, response.Operations, 2)
		require.Equal(t, uint32(2), response.PaginationInfo.StartIndex)
		require.Equal(t, uint32(4), response.PaginationInfo.TotalEntries)
		require.Equal(t, "op3", response.Operations[0].Name)
		require.Equal(t, "op4", response.Operations[1].Name)
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

		require.Len(t, response.Operations, 1)
		require.Equal(t, uint32(3), response.PaginationInfo.StartIndex)
		require.Equal(t, uint32(4), response.PaginationInfo.TotalEntries)
		require.Equal(t, "op4", response.Operations[0].Name)
	})
}

func TestListOperations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
						structpb.NewStringValue("allowed/foo"),
					},
				}),
			},
		}),
	}))), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)
	killOperationsAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, killOperationsAuthorizer, nil, 2)

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
								InstanceNamePrefix: "allowed",
							},
						},
					},
					InstanceNameSuffix: "foo",
				},
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "allowed",
							},
						},
					},
					InstanceNameSuffix: "bar",
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
				{
					InvocationName: &buildqueuestate.InvocationName{
						SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
							PlatformQueueName: &buildqueuestate.PlatformQueueName{
								InstanceNamePrefix: "forbidden",
							},
						},
					},
					InstanceNameSuffix: "foo",
				},
			},
			PaginationInfo: &buildqueuestate.PaginationInfo{
				StartIndex:   0,
				TotalEntries: 5,
			},
		}
		bqsClient.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(clientResponse, nil)

		resp, err := bqsServer.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize: 5,
		})
		require.NoError(t, err)
		require.Equal(t, clientResponse.Operations[0:2], resp.Operations)
		require.Equal(t, uint32(0), resp.PaginationInfo.StartIndex)
		require.Equal(t, uint32(2), resp.PaginationInfo.TotalEntries)
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

func TestIsAllowedToKillOperation(t *testing.T) {
	ctrl, ctx := gomock.WithContext(auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
					},
				}),
			},
		}),
	}))), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return true })
	killOperationsAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName)"),
	)

	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, killOperationsAuthorizer, nil, 2)

	t.Run("EmptyOperationName", func(t *testing.T) {
		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.NotFound, "Operation was not found"))
		result := bqsServer.IsAllowedToKillOperation(ctx, "")
		require.False(t, result)
	})

	t.Run("OperationNotFound", func(t *testing.T) {
		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.NotFound, "Operation was not found"))
		result := bqsServer.IsAllowedToKillOperation(ctx, "not found")
		require.False(t, result)
	})

	t.Run("OperationFoundButNotAllowedByAuthorizer", func(t *testing.T) {
		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
			Operation: &buildqueuestate.OperationState{
				Name: "op1",
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "forbidden",
						},
					},
				},
			},
		}, nil)
		result := bqsServer.IsAllowedToKillOperation(ctx, "op1")
		require.False(t, result)
	})

	t.Run("OperationFoundButNotAllowedByAuthorizer", func(t *testing.T) {
		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
			Operation: &buildqueuestate.OperationState{
				Name: "op1",
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "allowed",
						},
					},
				},
			},
		}, nil)
		result := bqsServer.IsAllowedToKillOperation(ctx, "op1")
		require.True(t, result)
	})
}

func TestCheckKillOperationAuthorization(t *testing.T) {
	ctrl, ctx := gomock.WithContext(auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
					},
				}),
			},
		}),
	}))), t)

	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return true })
	killOperationsAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName)"),
	)

	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, killOperationsAuthorizer, nil, 2)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/checkPermissions/killOperation/op1", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		bqsServer.CheckKillOperationAuthorization(rr, req)

		require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		require.Equal(t, "Method not allowed\n", rr.Body.String())
	})

	t.Run("AllowedToKillOperation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/checkPermissions/killOperation/op1", nil)
		require.NoError(t, err)

		req = req.WithContext(ctx)
		req.SetPathValue("operationName", "op1")

		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
			Operation: &buildqueuestate.OperationState{
				Name: "op1",
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "allowed",
						},
					},
				},
			},
		}, nil)

		rr := httptest.NewRecorder()
		bqsServer.CheckKillOperationAuthorization(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": true}`, rr.Body.String())
	})

	t.Run("NotAllowedToKillOperation", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/checkPermissions/killOperation/op2", nil)
		require.NoError(t, err)

		req = req.WithContext(ctx)
		req.SetPathValue("operationName", "op2")

		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
			Operation: &buildqueuestate.OperationState{
				Name: "op2",
				InvocationName: &buildqueuestate.InvocationName{
					SizeClassQueueName: &buildqueuestate.SizeClassQueueName{
						PlatformQueueName: &buildqueuestate.PlatformQueueName{
							InstanceNamePrefix: "forbidden",
						},
					},
				},
			},
		}, nil)

		rr := httptest.NewRecorder()
		bqsServer.CheckKillOperationAuthorization(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": false}`, rr.Body.String())
	})

	t.Run("OperationNotFound", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/checkPermissions/killOperation/op3", nil)
		require.NoError(t, err)

		req = req.WithContext(ctx)
		req.SetPathValue("operationName", "op3")

		bqsClient.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(nil, status.Errorf(codes.NotFound, "Operation was not found"))

		rr := httptest.NewRecorder()
		bqsServer.CheckKillOperationAuthorization(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": false}`, rr.Body.String())
	})
}

func TestCensorWorkerState(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)
	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed/worker"),
					},
				}),
			},
		}),
	})))

	t.Run("NoCurrentOperation", func(t *testing.T) {
		worker := &buildqueuestate.WorkerState{}
		censorWorkerState(ctx, a, "allowed", worker)
		require.Nil(t, worker.CurrentOperation)
	})

	t.Run("AllowedOperation", func(t *testing.T) {
		operation := &buildqueuestate.OperationState{InstanceNameSuffix: "worker"}
		worker := &buildqueuestate.WorkerState{CurrentOperation: operation}
		censorWorkerState(ctx, a, "allowed", worker)
		require.Same(t, operation, worker.CurrentOperation)
	})

	t.Run("DeniedOperation", func(t *testing.T) {
		worker := &buildqueuestate.WorkerState{
			CurrentOperation: &buildqueuestate.OperationState{InstanceNameSuffix: "forbidden"},
		}
		censorWorkerState(ctx, a, "allowed", worker)
		require.Nil(t, worker.CurrentOperation)
	})
}

func newContextWithPermittedInstanceNames(names ...string) context.Context {
	values := make([]*structpb.Value, 0, len(names))
	for _, name := range names {
		values = append(values, structpb.NewStringValue(name))
	}
	return auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{Values: values}),
			},
		}),
	})))
}

func newPermittedInstanceNamesAuthorizer() auth.Authorizer {
	return auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)
}

func newSizeClassQueueName(instanceNamePrefix string) *buildqueuestate.SizeClassQueueName {
	return &buildqueuestate.SizeClassQueueName{
		PlatformQueueName: &buildqueuestate.PlatformQueueName{
			InstanceNamePrefix: instanceNamePrefix,
		},
		SizeClass: 1,
	}
}

func TestGetInstanceNamePrefixFromSizeClassQueueName(t *testing.T) {
	t.Run("NilSizeClassQueueName", func(t *testing.T) {
		_, err := getInstanceNamePrefixFromSizeClassQueueName(nil)
		testutil.RequireEqualStatus(t, status.Error(codes.InvalidArgument, "Request does not contain a valid InstanceNamePrefix"), err)
	})

	t.Run("NilPlatformQueueName", func(t *testing.T) {
		_, err := getInstanceNamePrefixFromSizeClassQueueName(&buildqueuestate.SizeClassQueueName{})
		testutil.RequireEqualStatus(t, status.Error(codes.InvalidArgument, "Request does not contain a valid InstanceNamePrefix"), err)
	})

	t.Run("EmptyPrefix", func(t *testing.T) {
		instanceNamePrefix, err := getInstanceNamePrefixFromSizeClassQueueName(newSizeClassQueueName(""))
		require.NoError(t, err)
		require.Equal(t, "", instanceNamePrefix)
	})

	t.Run("Prefix", func(t *testing.T) {
		instanceNamePrefix, err := getInstanceNamePrefixFromSizeClassQueueName(newSizeClassQueueName("a/b"))
		require.NoError(t, err)
		require.Equal(t, "a/b", instanceNamePrefix)
	})
}

func TestCensorQueuedOperationState(t *testing.T) {
	a := newPermittedInstanceNamesAuthorizer()
	ctx := newContextWithPermittedInstanceNames("allowed", "allowed/foo")

	newOperation := func(instanceNameSuffix string) *buildqueuestate.OperationState {
		return &buildqueuestate.OperationState{
			Name:               "op1",
			Priority:           3,
			ExpectedDuration:   &durationpb.Duration{Seconds: 5},
			QueuedTimestamp:    &timestamppb.Timestamp{Seconds: 1000},
			Timeout:            &timestamppb.Timestamp{Seconds: 2000},
			ActionDigest:       &remoteexecution.Digest{Hash: "abc", SizeBytes: 12},
			TargetId:           "//foo:bar",
			InstanceNameSuffix: instanceNameSuffix,
			DigestFunction:     remoteexecution.DigestFunction_SHA256,
			Stage:              &buildqueuestate.OperationState_Queued{Queued: &emptypb.Empty{}},
		}
	}

	t.Run("AllowedSuffix", func(t *testing.T) {
		operation := newOperation("foo")
		censorQueuedOperationState(ctx, a, "allowed", operation)
		testutil.RequireEqualProto(t, newOperation("foo"), operation)
	})

	t.Run("AllowedEmptySuffix", func(t *testing.T) {
		operation := newOperation("")
		censorQueuedOperationState(ctx, a, "allowed", operation)
		testutil.RequireEqualProto(t, newOperation(""), operation)
	})

	t.Run("DisallowedSuffix", func(t *testing.T) {
		operation := newOperation("bar")
		censorQueuedOperationState(ctx, a, "allowed", operation)
		testutil.RequireEqualProto(t, &buildqueuestate.OperationState{
			Priority:         3,
			ExpectedDuration: &durationpb.Duration{Seconds: 5},
			QueuedTimestamp:  &timestamppb.Timestamp{Seconds: 1000},
			Stage:            &buildqueuestate.OperationState_Queued{Queued: &emptypb.Empty{}},
		}, operation)
	})

	t.Run("JoinMatchesIsPrefixSuffixAllowed", func(t *testing.T) {
		for _, tc := range []struct{ prefix, suffix string }{
			{"", ""},
			{"", "allowed"},
			{"", "foo"},
			{"allowed", ""},
			{"allowed", "foo"},
			{"allowed", "bar"},
			{"forbidden", "foo"},
		} {
			operation := newOperation(tc.suffix)
			censorQueuedOperationState(ctx, a, tc.prefix, operation)
			require.Equal(t, isPrefixSuffixAllowed(ctx, a, tc.prefix, tc.suffix), operation.ActionDigest != nil, "prefix %#v suffix %#v", tc.prefix, tc.suffix)
		}
	})
}

func TestListInvocationChildren(t *testing.T) {
	ctrl, ctx := gomock.WithContext(newContextWithPermittedInstanceNames("allowed"), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := newPermittedInstanceNamesAuthorizer()
	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, readAuthorizer, nil, 2)

	t.Run("Allowed", func(t *testing.T) {
		req := &buildqueuestate.ListInvocationChildrenRequest{
			InvocationName: &buildqueuestate.InvocationName{
				SizeClassQueueName: newSizeClassQueueName("allowed"),
				Ids:                []*anypb.Any{{TypeUrl: "type.googleapis.com/foo", Value: []byte{1, 2, 3}}},
			},
			Filter: buildqueuestate.ListInvocationChildrenRequest_QUEUED,
		}
		clientResponse := &buildqueuestate.ListInvocationChildrenResponse{
			Children: []*buildqueuestate.InvocationChildState{
				{Id: &anypb.Any{TypeUrl: "type.googleapis.com/foo", Value: []byte{9}}, State: &buildqueuestate.InvocationState{QueuedChildrenCount: 2}},
				{Id: &anypb.Any{TypeUrl: "type.googleapis.com/foo", Value: []byte{4}}, State: &buildqueuestate.InvocationState{QueuedChildrenCount: 1}},
			},
		}
		bqsClient.EXPECT().ListInvocationChildren(gomock.Any(), req).Return(clientResponse, nil)

		resp, err := bqsServer.ListInvocationChildren(ctx, req)
		require.NoError(t, err)
		require.Same(t, clientResponse, resp)
	})

	t.Run("DeniedPrefix", func(t *testing.T) {
		_, err := bqsServer.ListInvocationChildren(ctx, &buildqueuestate.ListInvocationChildrenRequest{
			InvocationName: &buildqueuestate.InvocationName{SizeClassQueueName: newSizeClassQueueName("forbidden")},
		})
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, "Not allowed to list invocation children for instance name prefix forbidden"), err)
	})

	t.Run("MissingInvocationName", func(t *testing.T) {
		_, err := bqsServer.ListInvocationChildren(ctx, &buildqueuestate.ListInvocationChildrenRequest{})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("InvocationNotFound", func(t *testing.T) {
		bqsClient.EXPECT().ListInvocationChildren(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "Invocation not found"))

		_, err := bqsServer.ListInvocationChildren(ctx, &buildqueuestate.ListInvocationChildrenRequest{
			InvocationName: &buildqueuestate.InvocationName{SizeClassQueueName: newSizeClassQueueName("allowed")},
		})
		testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Failed to list invocation children: Invocation not found"), err)
	})
}

func TestListQueuedOperations(t *testing.T) {
	ctrl, ctx := gomock.WithContext(newContextWithPermittedInstanceNames("allowed", "allowed/foo"), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := newPermittedInstanceNamesAuthorizer()
	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, readAuthorizer, nil, 2)

	invocationName := &buildqueuestate.InvocationName{SizeClassQueueName: newSizeClassQueueName("allowed")}
	startAfter := &buildqueuestate.ListQueuedOperationsRequest_StartAfter{
		Priority:         1,
		ExpectedDuration: &durationpb.Duration{Seconds: 2},
		QueuedTimestamp:  &timestamppb.Timestamp{Seconds: 3},
	}

	t.Run("ClampsPageSize", func(t *testing.T) {
		bqsClient.EXPECT().ListQueuedOperations(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, req *buildqueuestate.ListQueuedOperationsRequest, opts ...any) (*buildqueuestate.ListQueuedOperationsResponse, error) {
				testutil.RequireEqualProto(t, &buildqueuestate.ListQueuedOperationsRequest{
					InvocationName: invocationName,
					PageSize:       2,
					StartAfter:     startAfter,
				}, req)
				return &buildqueuestate.ListQueuedOperationsResponse{PaginationInfo: &buildqueuestate.PaginationInfo{}}, nil
			},
		)

		_, err := bqsServer.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{
			InvocationName: invocationName,
			PageSize:       100,
			StartAfter:     startAfter,
		})
		require.NoError(t, err)
	})

	t.Run("KeepsSmallerPageSize", func(t *testing.T) {
		bqsClient.EXPECT().ListQueuedOperations(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, req *buildqueuestate.ListQueuedOperationsRequest, opts ...any) (*buildqueuestate.ListQueuedOperationsResponse, error) {
				require.Equal(t, uint32(1), req.PageSize)
				require.Nil(t, req.StartAfter)
				return &buildqueuestate.ListQueuedOperationsResponse{PaginationInfo: &buildqueuestate.PaginationInfo{}}, nil
			},
		)

		_, err := bqsServer.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{
			InvocationName: invocationName,
			PageSize:       1,
		})
		require.NoError(t, err)
	})

	t.Run("CensorsOperations", func(t *testing.T) {
		bqsClient.EXPECT().ListQueuedOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListQueuedOperationsResponse{
			QueuedOperations: []*buildqueuestate.OperationState{
				{
					Name:               "op1",
					Priority:           1,
					ActionDigest:       &remoteexecution.Digest{Hash: "abc", SizeBytes: 1},
					TargetId:           "//allowed",
					InstanceNameSuffix: "foo",
					DigestFunction:     remoteexecution.DigestFunction_SHA256,
				},
				{
					Name:               "op2",
					Priority:           2,
					ActionDigest:       &remoteexecution.Digest{Hash: "def", SizeBytes: 2},
					TargetId:           "//forbidden",
					InstanceNameSuffix: "bar",
					DigestFunction:     remoteexecution.DigestFunction_SHA256,
				},
			},
			PaginationInfo: &buildqueuestate.PaginationInfo{StartIndex: 4, TotalEntries: 10},
		}, nil)

		resp, err := bqsServer.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{
			InvocationName: invocationName,
			PageSize:       2,
		})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, &buildqueuestate.ListQueuedOperationsResponse{
			QueuedOperations: []*buildqueuestate.OperationState{
				{
					Name:               "op1",
					Priority:           1,
					ActionDigest:       &remoteexecution.Digest{Hash: "abc", SizeBytes: 1},
					TargetId:           "//allowed",
					InstanceNameSuffix: "foo",
					DigestFunction:     remoteexecution.DigestFunction_SHA256,
				},
				{
					Priority: 2,
				},
			},
			PaginationInfo: &buildqueuestate.PaginationInfo{StartIndex: 4, TotalEntries: 10},
		}, resp)
	})

	t.Run("DeniedPrefix", func(t *testing.T) {
		_, err := bqsServer.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{
			InvocationName: &buildqueuestate.InvocationName{SizeClassQueueName: newSizeClassQueueName("forbidden")},
			PageSize:       2,
		})
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, "Not allowed to list queued operations for instance name prefix forbidden"), err)
	})

	t.Run("MissingInvocationName", func(t *testing.T) {
		_, err := bqsServer.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{PageSize: 2})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("ClientError", func(t *testing.T) {
		bqsClient.EXPECT().ListQueuedOperations(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "Invocation not found"))

		_, err := bqsServer.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{
			InvocationName: invocationName,
			PageSize:       2,
		})
		testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Failed to list queued operations: Invocation not found"), err)
	})
}

func TestListDrains(t *testing.T) {
	ctrl, ctx := gomock.WithContext(newContextWithPermittedInstanceNames("allowed"), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := newPermittedInstanceNamesAuthorizer()
	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, readAuthorizer, nil, 2)

	t.Run("Allowed", func(t *testing.T) {
		req := &buildqueuestate.ListDrainsRequest{SizeClassQueueName: newSizeClassQueueName("allowed")}
		clientResponse := &buildqueuestate.ListDrainsResponse{
			Drains: []*buildqueuestate.DrainState{
				{WorkerIdPattern: map[string]string{"hostname": "a"}, CreatedTimestamp: &timestamppb.Timestamp{Seconds: 1}},
				{WorkerIdPattern: map[string]string{"hostname": "b"}, CreatedTimestamp: &timestamppb.Timestamp{Seconds: 2}},
			},
		}
		bqsClient.EXPECT().ListDrains(gomock.Any(), req).Return(clientResponse, nil)

		resp, err := bqsServer.ListDrains(ctx, req)
		require.NoError(t, err)
		require.Same(t, clientResponse, resp)
	})

	t.Run("DeniedPrefix", func(t *testing.T) {
		_, err := bqsServer.ListDrains(ctx, &buildqueuestate.ListDrainsRequest{SizeClassQueueName: newSizeClassQueueName("forbidden")})
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, "Not allowed to list drains for instance name prefix forbidden"), err)
	})

	t.Run("MissingSizeClassQueueName", func(t *testing.T) {
		_, err := bqsServer.ListDrains(ctx, &buildqueuestate.ListDrainsRequest{})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("ClientError", func(t *testing.T) {
		bqsClient.EXPECT().ListDrains(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "No workers exist"))

		_, err := bqsServer.ListDrains(ctx, &buildqueuestate.ListDrainsRequest{SizeClassQueueName: newSizeClassQueueName("allowed")})
		testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Failed to list drains: No workers exist"), err)
	})
}

func testAddOrRemoveDrain(
	t *testing.T,
	expect func(bqsClient *mock.MockBuildQueueStateClient) *gomock.Call,
	call func(ctx context.Context, bqsServer *BuildQueueStateServerImpl, req *buildqueuestate.AddOrRemoveDrainRequest) (*emptypb.Empty, error),
	errorMessage string,
) {
	ctrl, ctx := gomock.WithContext(newContextWithPermittedInstanceNames("readable", "drainable", "both"), t)
	bqsClient := mock.NewMockBuildQueueStateClient(ctrl)
	readAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) && instanceName != 'drainable'"),
	)
	drainAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) && instanceName != 'readable'"),
	)
	bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, readAuthorizer, drainAuthorizer, 2)
	newRequest := func(instanceNamePrefix string) *buildqueuestate.AddOrRemoveDrainRequest {
		return &buildqueuestate.AddOrRemoveDrainRequest{
			SizeClassQueueName: newSizeClassQueueName(instanceNamePrefix),
			WorkerIdPattern:    map[string]string{"hostname": "worker-1", "thread": "3"},
		}
	}

	t.Run("NoDrainAuthorizer", func(t *testing.T) {
		bqsServer := NewBuildQueueStateServerImpl(bqsClient, readAuthorizer, readAuthorizer, nil, 2)
		_, err := call(ctx, bqsServer, newRequest("both"))
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, "Not allowed to modify drains for instance name prefix both"), err)
	})

	t.Run("ReadAllowedDrainDenied", func(t *testing.T) {
		_, err := call(ctx, bqsServer, newRequest("readable"))
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, "Not allowed to modify drains for instance name prefix readable"), err)
	})

	t.Run("ReadDeniedDrainAllowed", func(t *testing.T) {
		_, err := call(ctx, bqsServer, newRequest("drainable"))
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, "Not allowed to modify drains for instance name prefix drainable"), err)
	})

	t.Run("Allowed", func(t *testing.T) {
		expect(bqsClient).DoAndReturn(
			func(ctx context.Context, req *buildqueuestate.AddOrRemoveDrainRequest, opts ...any) (*emptypb.Empty, error) {
				testutil.RequireEqualProto(t, newRequest("both"), req)
				return &emptypb.Empty{}, nil
			},
		)
		resp, err := call(ctx, bqsServer, newRequest("both"))
		require.NoError(t, err)
		testutil.RequireEqualProto(t, &emptypb.Empty{}, resp)
	})

	t.Run("ClientError", func(t *testing.T) {
		expect(bqsClient).Return(nil, status.Error(codes.PermissionDenied, "Authorization failed"))
		_, err := call(ctx, bqsServer, newRequest("both"))
		testutil.RequireEqualStatus(t, status.Error(codes.PermissionDenied, errorMessage+": Authorization failed"), err)
	})

	t.Run("MissingSizeClassQueueName", func(t *testing.T) {
		_, err := call(ctx, bqsServer, &buildqueuestate.AddOrRemoveDrainRequest{})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

func TestAddDrain(t *testing.T) {
	testAddOrRemoveDrain(
		t,
		func(bqsClient *mock.MockBuildQueueStateClient) *gomock.Call {
			return bqsClient.EXPECT().AddDrain(gomock.Any(), gomock.Any())
		},
		func(ctx context.Context, bqsServer *BuildQueueStateServerImpl, req *buildqueuestate.AddOrRemoveDrainRequest) (*emptypb.Empty, error) {
			return bqsServer.AddDrain(ctx, req)
		},
		"Failed to add drain",
	)
}

func TestRemoveDrain(t *testing.T) {
	testAddOrRemoveDrain(
		t,
		func(bqsClient *mock.MockBuildQueueStateClient) *gomock.Call {
			return bqsClient.EXPECT().RemoveDrain(gomock.Any(), gomock.Any())
		},
		func(ctx context.Context, bqsServer *BuildQueueStateServerImpl, req *buildqueuestate.AddOrRemoveDrainRequest) (*emptypb.Empty, error) {
			return bqsServer.RemoveDrain(ctx, req)
		},
		"Failed to remove drain",
	)
}

func TestIsAllowedToModifyDrains(t *testing.T) {
	ctx := newContextWithPermittedInstanceNames("readable", "drainable", "both")
	readAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) && instanceName != 'drainable'"),
	)
	drainAuthorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) && instanceName != 'readable'"),
	)

	t.Run("NoDrainAuthorizer", func(t *testing.T) {
		bqsServer := NewBuildQueueStateServerImpl(nil, readAuthorizer, readAuthorizer, nil, 2)
		require.False(t, bqsServer.IsAllowedToModifyDrains(ctx, "both"))
	})

	bqsServer := NewBuildQueueStateServerImpl(nil, readAuthorizer, readAuthorizer, drainAuthorizer, 2)

	t.Run("ReadDenied", func(t *testing.T) {
		require.False(t, bqsServer.IsAllowedToModifyDrains(ctx, "drainable"))
	})

	t.Run("DrainDenied", func(t *testing.T) {
		require.False(t, bqsServer.IsAllowedToModifyDrains(ctx, "readable"))
	})

	t.Run("BothDenied", func(t *testing.T) {
		require.False(t, bqsServer.IsAllowedToModifyDrains(ctx, "forbidden"))
	})

	t.Run("Allowed", func(t *testing.T) {
		require.True(t, bqsServer.IsAllowedToModifyDrains(ctx, "both"))
	})
}

func TestCheckModifyDrainsAuthorization(t *testing.T) {
	ctx := newContextWithPermittedInstanceNames("allowed", "a/b")
	readAuthorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return true })
	drainAuthorizer := newPermittedInstanceNamesAuthorizer()
	bqsServer := NewBuildQueueStateServerImpl(nil, readAuthorizer, readAuthorizer, drainAuthorizer, 2)

	check := func(t *testing.T, bqsServer *BuildQueueStateServerImpl, instanceNamePrefix string) *httptest.ResponseRecorder {
		req, err := http.NewRequest("GET", "/api/v1/checkPermissions/modifyDrains/"+instanceNamePrefix, nil)
		require.NoError(t, err)

		req = req.WithContext(ctx)
		req.SetPathValue("instanceNamePrefix", instanceNamePrefix)

		rr := httptest.NewRecorder()
		bqsServer.CheckModifyDrainsAuthorization(rr, req)
		return rr
	}

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/v1/checkPermissions/modifyDrains/allowed", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		bqsServer.CheckModifyDrainsAuthorization(rr, req)

		require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		require.Equal(t, "Method not allowed\n", rr.Body.String())
	})

	t.Run("Allowed", func(t *testing.T) {
		rr := check(t, bqsServer, "allowed")
		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": true}`, rr.Body.String())
	})

	t.Run("Denied", func(t *testing.T) {
		rr := check(t, bqsServer, "forbidden")
		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": false}`, rr.Body.String())
	})

	t.Run("NoDrainAuthorizer", func(t *testing.T) {
		rr := check(t, NewBuildQueueStateServerImpl(nil, readAuthorizer, readAuthorizer, nil, 2), "allowed")
		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": false}`, rr.Body.String())
	})

	t.Run("EmptyPrefix", func(t *testing.T) {
		rr := check(t, bqsServer, "")
		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": true}`, rr.Body.String())
	})

	t.Run("MultiSegmentPrefix", func(t *testing.T) {
		rr := check(t, bqsServer, "a/b")
		require.Equal(t, http.StatusOK, rr.Code)
		require.JSONEq(t, `{"allowed": true}`, rr.Body.String())
	})

	t.Run("Route", func(t *testing.T) {
		router := http.NewServeMux()
		router.HandleFunc("GET /api/v1/checkPermissions/modifyDrains/{instanceNamePrefix...}", bqsServer.CheckModifyDrainsAuthorization)
		for path, expected := range map[string]string{
			"/api/v1/checkPermissions/modifyDrains/":          `{"allowed": true}`,
			"/api/v1/checkPermissions/modifyDrains/allowed":   `{"allowed": true}`,
			"/api/v1/checkPermissions/modifyDrains/a/b":       `{"allowed": true}`,
			"/api/v1/checkPermissions/modifyDrains/a":         `{"allowed": false}`,
			"/api/v1/checkPermissions/modifyDrains/forbidden": `{"allowed": false}`,
		} {
			req := httptest.NewRequest("GET", path, nil).WithContext(ctx)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			require.Equal(t, http.StatusOK, rr.Code, path)
			require.JSONEq(t, expected, rr.Body.String(), path)
		}
	})
}
