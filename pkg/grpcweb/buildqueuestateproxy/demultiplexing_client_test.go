package buildqueuestateproxy

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-remote-execution/pkg/proto/buildqueuestate"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/buildbarn/bb-storage/pkg/testutil"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func route(match, add string, client buildqueuestate.BuildQueueStateClient) Route {
	return Route{
		MatchInstanceNamePrefix: util.Must(digest.NewInstanceName(match)),
		AddInstanceNamePrefix:   util.Must(digest.NewInstanceName(add)),
		Client:                  client,
	}
}

func queue(instanceNamePrefix string) *buildqueuestate.SizeClassQueueName {
	return &buildqueuestate.SizeClassQueueName{
		PlatformQueueName: &buildqueuestate.PlatformQueueName{InstanceNamePrefix: instanceNamePrefix},
	}
}

func operation(name, instanceNamePrefix string) *buildqueuestate.OperationState {
	return &buildqueuestate.OperationState{
		Name:           name,
		InvocationName: &buildqueuestate.InvocationName{SizeClassQueueName: queue(instanceNamePrefix)},
	}
}

func listWorkers(instanceNamePrefix string) *buildqueuestate.ListWorkersRequest {
	return &buildqueuestate.ListWorkersRequest{
		Filter: &buildqueuestate.ListWorkersRequest_Filter{
			Type: &buildqueuestate.ListWorkersRequest_Filter_All{All: queue(instanceNamePrefix)},
		},
		PageSize: 10,
	}
}

func platformQueues(instanceNamePrefixes ...string) *buildqueuestate.ListPlatformQueuesResponse {
	response := &buildqueuestate.ListPlatformQueuesResponse{}
	for _, p := range instanceNamePrefixes {
		response.PlatformQueues = append(response.PlatformQueues, &buildqueuestate.PlatformQueueState{
			Name: &buildqueuestate.PlatformQueueName{InstanceNamePrefix: p},
		})
	}
	return response
}

func selecting(instanceName string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(SelectedInstanceNameMetadataKey, instanceName))
}

// Requests that carry an instance name are routed by it, never by the UI's
// selection, mirroring bb_storage's demultiplexing of Execute requests.
func TestRoutesByInstanceNameInRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	stripping := mock.NewMockBuildQueueStateClient(ctrl)
	forwarding := mock.NewMockBuildQueueStateClient(ctrl)
	catchAll := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{
		route("ci", "", catchAll),
		route("ci/aws", "", stripping),
		route("ci/gcp", "ci/gcp", forwarding),
	})

	t.Run("LongestPrefixStripsMatchedPrefix", func(t *testing.T) {
		stripping.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, listWorkers("mars"))).
			Return(&buildqueuestate.ListWorkersResponse{
				Workers: []*buildqueuestate.WorkerState{{CurrentOperation: operation("op", "mars")}},
			}, nil)

		response, err := router.ListWorkers(selecting("ci/gcp"), listWorkers("ci/aws/mars"))
		require.NoError(t, err)
		testutil.RequireEqualProto(t, &buildqueuestate.ListWorkersResponse{
			Workers: []*buildqueuestate.WorkerState{{CurrentOperation: operation("op", "ci/aws/mars")}},
		}, response)
	})

	t.Run("ExactMatchWithPerfectForwarding", func(t *testing.T) {
		forwarding.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, listWorkers("ci/gcp"))).
			Return(&buildqueuestate.ListWorkersResponse{}, nil)

		_, err := router.ListWorkers(context.Background(), listWorkers("ci/gcp"))
		require.NoError(t, err)
	})

	t.Run("ShorterPrefixForOtherNames", func(t *testing.T) {
		catchAll.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, listWorkers("azure"))).
			Return(&buildqueuestate.ListWorkersResponse{}, nil)

		_, err := router.ListWorkers(context.Background(), listWorkers("ci/azure"))
		require.NoError(t, err)
	})

	t.Run("PrefixMatchesWholeComponents", func(t *testing.T) {
		catchAll.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, listWorkers("aws-east"))).
			Return(&buildqueuestate.ListWorkersResponse{}, nil)

		_, err := router.ListWorkers(context.Background(), listWorkers("ci/aws-east"))
		require.NoError(t, err)
	})

	t.Run("UnknownInstanceName", func(t *testing.T) {
		_, err := router.ListWorkers(context.Background(), listWorkers("other"))
		testutil.RequireEqualStatus(t, status.Error(codes.InvalidArgument, "No scheduler is configured for instance name \"other\""), err)
	})

	t.Run("RequestIsNotModified", func(t *testing.T) {
		stripping.EXPECT().ListWorkers(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListWorkersResponse{}, nil)
		request := listWorkers("ci/aws/mars")

		_, err := router.ListWorkers(context.Background(), request)
		require.NoError(t, err)
		testutil.RequireEqualProto(t, listWorkers("ci/aws/mars"), request)
	})
}

func TestRoutesAndRewritesTypedPlatformQueueRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	client := NewDemultiplexingClient([]Route{route("ci", "scheduler", backend)})
	ctx := selecting("not-configured")

	allWorkers := &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_All{All: queue("ci/workers")},
	}}
	executingWorkers := &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_Executing{Executing: &buildqueuestate.InvocationName{SizeClassQueueName: queue("ci/executing")}},
	}}
	idleSynchronizingWorkers := &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_IdleSynchronizing{IdleSynchronizing: &buildqueuestate.InvocationName{SizeClassQueueName: queue("ci/idle")}},
	}}
	backend.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_All{All: queue("scheduler/workers")},
	}})).Return(&buildqueuestate.ListWorkersResponse{Workers: []*buildqueuestate.WorkerState{{CurrentOperation: operation("worker", "scheduler/workers")}}}, nil)
	backend.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_Executing{Executing: &buildqueuestate.InvocationName{SizeClassQueueName: queue("scheduler/executing")}},
	}})).Return(&buildqueuestate.ListWorkersResponse{}, nil)
	backend.EXPECT().ListWorkers(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_IdleSynchronizing{IdleSynchronizing: &buildqueuestate.InvocationName{SizeClassQueueName: queue("scheduler/idle")}},
	}})).Return(&buildqueuestate.ListWorkersResponse{}, nil)

	response, err := client.ListWorkers(ctx, allWorkers)
	require.NoError(t, err)
	testutil.RequireEqualProto(t, &buildqueuestate.ListWorkersResponse{Workers: []*buildqueuestate.WorkerState{{CurrentOperation: operation("worker", "ci/workers")}}}, response)
	_, err = client.ListWorkers(ctx, executingWorkers)
	require.NoError(t, err)
	_, err = client.ListWorkers(ctx, idleSynchronizingWorkers)
	require.NoError(t, err)
	testutil.RequireEqualProto(t, &buildqueuestate.ListWorkersRequest{Filter: &buildqueuestate.ListWorkersRequest_Filter{
		Type: &buildqueuestate.ListWorkersRequest_Filter_All{All: queue("ci/workers")},
	}}, allWorkers)

	invocationName := func(prefix string) *buildqueuestate.InvocationName {
		return &buildqueuestate.InvocationName{SizeClassQueueName: queue(prefix)}
	}
	backend.EXPECT().ListInvocationChildren(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListInvocationChildrenRequest{InvocationName: invocationName("scheduler/children")})).
		Return(&buildqueuestate.ListInvocationChildrenResponse{}, nil)
	_, err = client.ListInvocationChildren(ctx, &buildqueuestate.ListInvocationChildrenRequest{InvocationName: invocationName("ci/children")})
	require.NoError(t, err)

	backend.EXPECT().ListQueuedOperations(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListQueuedOperationsRequest{InvocationName: invocationName("scheduler/queued")})).
		Return(&buildqueuestate.ListQueuedOperationsResponse{QueuedOperations: []*buildqueuestate.OperationState{operation("queued", "scheduler/queued")}}, nil)
	queuedOperations, err := client.ListQueuedOperations(ctx, &buildqueuestate.ListQueuedOperationsRequest{InvocationName: invocationName("ci/queued")})
	require.NoError(t, err)
	testutil.RequireEqualProto(t, &buildqueuestate.ListQueuedOperationsResponse{QueuedOperations: []*buildqueuestate.OperationState{operation("queued", "ci/queued")}}, queuedOperations)

	backend.EXPECT().ListDrains(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListDrainsRequest{SizeClassQueueName: queue("scheduler/drains")})).
		Return(&buildqueuestate.ListDrainsResponse{}, nil)
	_, err = client.ListDrains(ctx, &buildqueuestate.ListDrainsRequest{SizeClassQueueName: queue("ci/drains")})
	require.NoError(t, err)

	for _, tc := range []struct {
		name string
		call func() error
	}{
		{
			name: "AddDrain",
			call: func() error {
				backend.EXPECT().AddDrain(gomock.Any(), testutil.EqProto(t, &buildqueuestate.AddOrRemoveDrainRequest{SizeClassQueueName: queue("scheduler/add")})).Return(&emptypb.Empty{}, nil)
				_, err := client.AddDrain(ctx, &buildqueuestate.AddOrRemoveDrainRequest{SizeClassQueueName: queue("ci/add")})
				return err
			},
		},
		{
			name: "RemoveDrain",
			call: func() error {
				backend.EXPECT().RemoveDrain(gomock.Any(), testutil.EqProto(t, &buildqueuestate.AddOrRemoveDrainRequest{SizeClassQueueName: queue("scheduler/remove")})).Return(&emptypb.Empty{}, nil)
				_, err := client.RemoveDrain(ctx, &buildqueuestate.AddOrRemoveDrainRequest{SizeClassQueueName: queue("ci/remove")})
				return err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) { require.NoError(t, tc.call()) })
	}

	backend.EXPECT().KillOperations(gomock.Any(), testutil.EqProto(t, &buildqueuestate.KillOperationsRequest{Filter: &buildqueuestate.KillOperationsRequest_Filter{
		Type: &buildqueuestate.KillOperationsRequest_Filter_SizeClassQueueWithoutWorkers{SizeClassQueueWithoutWorkers: queue("scheduler/kill")},
	}})).Return(&emptypb.Empty{}, nil)
	_, err = client.KillOperations(ctx, &buildqueuestate.KillOperationsRequest{Filter: &buildqueuestate.KillOperationsRequest_Filter{
		Type: &buildqueuestate.KillOperationsRequest_Filter_SizeClassQueueWithoutWorkers{SizeClassQueueWithoutWorkers: queue("ci/kill")},
	}})
	require.NoError(t, err)
}

func TestHeaderlessTerminateWorkersUsesSelectedScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	aws := mock.NewMockBuildQueueStateClient(ctrl)
	gcp := mock.NewMockBuildQueueStateClient(ctrl)
	client := NewDemultiplexingClient([]Route{route("ci/aws", "", aws), route("ci/gcp", "", gcp)})
	gcp.EXPECT().TerminateWorkers(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)
	_, err := client.TerminateWorkers(selecting("ci/gcp"), &buildqueuestate.TerminateWorkersRequest{})
	require.NoError(t, err)
}

func TestPopulatedWorkerAndQueuedOperationAttributionMustBelongToRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	client := NewDemultiplexingClient([]Route{route("ci", "scheduler", backend)})

	backend.EXPECT().ListWorkers(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListWorkersResponse{
		Workers: []*buildqueuestate.WorkerState{{CurrentOperation: operation("outside", "other")}},
	}, nil)
	_, err := client.ListWorkers(context.Background(), listWorkers("ci/workers"))
	testutil.RequireEqualStatus(t, status.Error(codes.Internal, "Scheduler for instance name prefix \"ci\" returned an instance name outside \"scheduler\""), err)

	backend.EXPECT().ListQueuedOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListQueuedOperationsResponse{
		QueuedOperations: []*buildqueuestate.OperationState{operation("outside", "other")},
	}, nil)
	_, err = client.ListQueuedOperations(context.Background(), &buildqueuestate.ListQueuedOperationsRequest{
		InvocationName: &buildqueuestate.InvocationName{SizeClassQueueName: queue("ci/queued")},
	})
	testutil.RequireEqualStatus(t, status.Error(codes.Internal, "Scheduler for instance name prefix \"ci\" returned an instance name outside \"scheduler\""), err)
}

func TestExecutingOperationAttributionMustBelongToRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	client := NewDemultiplexingClient([]Route{route("ci", "scheduler", backend)})
	operationWithExecutingWorker := operation("outer", "scheduler/outer")
	operationWithExecutingWorker.Stage = &buildqueuestate.OperationState_Executing{
		Executing: &buildqueuestate.WorkerState{CurrentOperation: operation("inner", "other")},
	}
	backend.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
		Operation: operationWithExecutingWorker,
	}, nil)
	_, err := client.GetOperation(selecting("ci"), &buildqueuestate.GetOperationRequest{OperationName: "outer"})
	testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Operation was not found"), err)
}

func TestOptionalWorkerAndQueuedOperationAttribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	client := NewDemultiplexingClient([]Route{route("ci", "scheduler", backend)})
	workersRequest := listWorkers("ci/workers")
	backend.EXPECT().ListWorkers(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListWorkersResponse{
		Workers: []*buildqueuestate.WorkerState{{CurrentOperation: &buildqueuestate.OperationState{Name: "unattributed"}}},
	}, nil)
	response, err := client.ListWorkers(context.Background(), workersRequest)
	require.NoError(t, err)
	testutil.RequireEqualProto(t, &buildqueuestate.ListWorkersResponse{
		Workers: []*buildqueuestate.WorkerState{{CurrentOperation: &buildqueuestate.OperationState{Name: "unattributed"}}},
	}, response)

	backend.EXPECT().ListQueuedOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListQueuedOperationsResponse{
		QueuedOperations: []*buildqueuestate.OperationState{{Name: "unattributed"}},
	}, nil)
	queued, err := client.ListQueuedOperations(context.Background(), &buildqueuestate.ListQueuedOperationsRequest{InvocationName: &buildqueuestate.InvocationName{SizeClassQueueName: queue("ci/queued")}})
	require.NoError(t, err)
	testutil.RequireEqualProto(t, &buildqueuestate.ListQueuedOperationsResponse{
		QueuedOperations: []*buildqueuestate.OperationState{{Name: "unattributed"}},
	}, queued)
}

// Requests without an instance name go to the scheduler the UI selected.
// With a single scheduler no selection is needed.
func TestRoutesBySelection(t *testing.T) {
	t.Run("SingleSchedulerNeedsNoSelection", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		only := mock.NewMockBuildQueueStateClient(ctrl)
		router := NewDemultiplexingClient([]Route{route("", "", only)})
		only.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(platformQueues("ci", ""), nil)

		response, err := router.ListPlatformQueues(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, platformQueues("ci", ""), response)
	})

	ctrl := gomock.NewController(t)
	aws := mock.NewMockBuildQueueStateClient(ctrl)
	gcp := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{
		route("ci/aws", "", aws),
		route("ci/gcp", "ci/gcp", gcp),
	})

	t.Run("MultipleSchedulersRequireSelection", func(t *testing.T) {
		_, err := router.ListPlatformQueues(context.Background(), &emptypb.Empty{})
		testutil.RequireEqualStatus(t, status.Error(codes.InvalidArgument, "Select a scheduler: more than one is configured, so metadata key \"x-bb-portal-scheduler-instance-name\" is required"), err)
	})

	t.Run("SelectionRewritesResponse", func(t *testing.T) {
		aws.EXPECT().GetOperation(gomock.Any(), testutil.EqProto(t, &buildqueuestate.GetOperationRequest{OperationName: "op"})).
			Return(&buildqueuestate.GetOperationResponse{Operation: operation("op", "mars")}, nil)

		response, err := router.GetOperation(selecting("ci/aws"), &buildqueuestate.GetOperationRequest{OperationName: "op"})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, &buildqueuestate.GetOperationResponse{Operation: operation("op", "ci/aws/mars")}, response)
	})

	t.Run("SelectionUsesLongestPrefix", func(t *testing.T) {
		gcp.EXPECT().KillOperations(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

		_, err := router.KillOperations(selecting("ci/gcp/europe"), &buildqueuestate.KillOperationsRequest{
			Filter: &buildqueuestate.KillOperationsRequest_Filter{
				Type: &buildqueuestate.KillOperationsRequest_Filter_OperationName{OperationName: "op"},
			},
		})
		require.NoError(t, err)
	})

	t.Run("UnknownSelection", func(t *testing.T) {
		_, err := router.ListPlatformQueues(selecting("ci/azure"), &emptypb.Empty{})
		testutil.RequireEqualStatus(t, status.Error(codes.InvalidArgument, "No scheduler is configured for instance name \"ci/azure\""), err)
	})

	t.Run("InvalidSelection", func(t *testing.T) {
		_, err := router.ListPlatformQueues(selecting("ci//aws"), &emptypb.Empty{})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}

// A scheduler may report instance names outside the namespace its route
// covers, e.g. when one scheduler serves several routes. Those must not be
// attributed to this route, so authorizers never see a mistranslated name.
func TestDropsInstanceNamesOutsideRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{
		route("ci/aws", "aws", backend),
	})
	ctx := selecting("ci/aws")

	t.Run("ListPlatformQueues", func(t *testing.T) {
		backend.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(platformQueues("aws/mars", "gcp", "aws"), nil)

		response, err := router.ListPlatformQueues(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, platformQueues("ci/aws/mars", "ci/aws"), response)
	})

	t.Run("GetOperation", func(t *testing.T) {
		backend.EXPECT().GetOperation(gomock.Any(), gomock.Any()).
			Return(&buildqueuestate.GetOperationResponse{Operation: operation("op", "gcp")}, nil)

		_, err := router.GetOperation(ctx, &buildqueuestate.GetOperationRequest{OperationName: "op"})
		testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Operation was not found"), err)
	})

	t.Run("ListOperationsKeepsFullPages", func(t *testing.T) {
		// The proxy pages until it receives a short page, so dropping
		// operations must not shorten a page that has more to follow.
		gomock.InOrder(
			backend.EXPECT().ListOperations(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListOperationsRequest{PageSize: 2})).
				Return(&buildqueuestate.ListOperationsResponse{
					Operations:     []*buildqueuestate.OperationState{operation("a", "aws"), operation("b", "gcp")},
					PaginationInfo: &buildqueuestate.PaginationInfo{TotalEntries: 5},
				}, nil),
			backend.EXPECT().ListOperations(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListOperationsRequest{
				PageSize:   2,
				StartAfter: &buildqueuestate.ListOperationsRequest_StartAfter{OperationName: "b"},
			})).
				Return(&buildqueuestate.ListOperationsResponse{
					Operations:     []*buildqueuestate.OperationState{operation("c", "aws/x"), operation("d", "gcp")},
					PaginationInfo: &buildqueuestate.PaginationInfo{TotalEntries: 5},
				}, nil),
		)

		response, err := router.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{PageSize: 2})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, &buildqueuestate.ListOperationsResponse{
			Operations:     []*buildqueuestate.OperationState{operation("a", "ci/aws"), operation("c", "ci/aws/x")},
			PaginationInfo: &buildqueuestate.PaginationInfo{TotalEntries: 5},
		}, response)
	})

	t.Run("ListOperationsStopsAtShortPage", func(t *testing.T) {
		backend.EXPECT().ListOperations(gomock.Any(), gomock.Any()).
			Return(&buildqueuestate.ListOperationsResponse{
				Operations:     []*buildqueuestate.OperationState{operation("e", "gcp")},
				PaginationInfo: &buildqueuestate.PaginationInfo{TotalEntries: 5},
			}, nil)

		response, err := router.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{
			PageSize:   2,
			StartAfter: &buildqueuestate.ListOperationsRequest_StartAfter{OperationName: "d"},
		})
		require.NoError(t, err)
		require.Empty(t, response.Operations)
	})
}

func TestPropagatesBackendErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{route("", "", backend)})

	t.Run("ListPlatformQueues", func(t *testing.T) {
		backend.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unavailable, "down"))

		_, err := router.ListPlatformQueues(context.Background(), &emptypb.Empty{})
		testutil.RequireEqualStatus(t, status.Error(codes.Unavailable, "down"), err)
	})

	// Only operations outside the route turn into NotFound; real failures
	// stay visible.
	t.Run("GetOperation", func(t *testing.T) {
		backend.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Internal, "boom"))

		_, err := router.GetOperation(context.Background(), &buildqueuestate.GetOperationRequest{OperationName: "op"})
		testutil.RequireEqualStatus(t, status.Error(codes.Internal, "boom"), err)
	})
}

func TestNestedRoutesDoNotLeakShadowedResponses(t *testing.T) {
	for _, tc := range []struct {
		name           string
		selected       string
		shadowed       string
		owned          string
		rewrittenOwned string
		routes         func(buildqueuestate.BuildQueueStateClient, buildqueuestate.BuildQueueStateClient) []Route
	}{
		{
			name:           "EmptyAndCi",
			selected:       "",
			shadowed:       "ci/project",
			owned:          "",
			rewrittenOwned: "",
			routes: func(root, nested buildqueuestate.BuildQueueStateClient) []Route {
				return []Route{route("", "", root), route("ci", "", nested)}
			},
		},
		{
			name:           "CiAndCiAws",
			selected:       "ci",
			shadowed:       "aws/mars",
			owned:          "mars",
			rewrittenOwned: "ci/mars",
			routes: func(parent, nested buildqueuestate.BuildQueueStateClient) []Route {
				return []Route{route("ci", "", parent), route("ci/aws", "", nested)}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			backend := mock.NewMockBuildQueueStateClient(ctrl)
			nested := mock.NewMockBuildQueueStateClient(ctrl)
			router := NewDemultiplexingClient(tc.routes(backend, nested))
			ctx := selecting(tc.selected)

			backend.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(platformQueues(tc.owned, tc.shadowed), nil)
			queues, err := router.ListPlatformQueues(ctx, &emptypb.Empty{})
			require.NoError(t, err)
			testutil.RequireEqualProto(t, platformQueues(tc.rewrittenOwned), queues)

			backend.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
				Operations: []*buildqueuestate.OperationState{
					operation("owned", tc.owned),
					operation("shadowed", tc.shadowed),
				},
			}, nil)
			operations, err := router.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{})
			require.NoError(t, err)
			testutil.RequireEqualProto(t, &buildqueuestate.ListOperationsResponse{
				Operations: []*buildqueuestate.OperationState{operation("owned", tc.rewrittenOwned)},
			}, operations)

			backend.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
				Operation: operation("shadowed", tc.shadowed),
			}, nil)
			_, err = router.GetOperation(ctx, &buildqueuestate.GetOperationRequest{OperationName: "shadowed"})
			testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Operation was not found"), err)
		})
	}
}

func TestResponsesWithoutPlatformQueueNameFailClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{route("ci", "", backend)})
	ctx := selecting("ci")

	t.Run("GetOperation", func(t *testing.T) {
		backend.EXPECT().GetOperation(gomock.Any(), gomock.Any()).Return(&buildqueuestate.GetOperationResponse{
			Operation: &buildqueuestate.OperationState{Name: "missing-attribution"},
		}, nil)
		_, err := router.GetOperation(ctx, &buildqueuestate.GetOperationRequest{OperationName: "missing-attribution"})
		testutil.RequireEqualStatus(t, status.Error(codes.NotFound, "Operation was not found"), err)
	})

	t.Run("ListOperations", func(t *testing.T) {
		backend.EXPECT().ListOperations(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListOperationsResponse{
			Operations: []*buildqueuestate.OperationState{{Name: "missing-attribution"}},
		}, nil)
		response, err := router.ListOperations(ctx, &buildqueuestate.ListOperationsRequest{})
		require.NoError(t, err)
		require.Empty(t, response.Operations)
	})

	t.Run("ListPlatformQueues", func(t *testing.T) {
		backend.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(&buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: []*buildqueuestate.PlatformQueueState{
				{},
				{Name: &buildqueuestate.PlatformQueueName{InstanceNamePrefix: "mars"}},
			},
		}, nil)
		response, err := router.ListPlatformQueues(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, platformQueues("ci/mars"), &buildqueuestate.ListPlatformQueuesResponse{
			PlatformQueues: response.PlatformQueues,
		})
	})
}

func TestListOperationsRejectsNonAdvancingCursor(t *testing.T) {
	ctrl := gomock.NewController(t)
	backend := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{route("ci/aws", "aws", backend)})
	gomock.InOrder(
		backend.EXPECT().ListOperations(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListOperationsRequest{PageSize: 2})).
			Return(&buildqueuestate.ListOperationsResponse{Operations: []*buildqueuestate.OperationState{
				operation("a", "aws"), operation("b", "gcp"),
			}}, nil),
		backend.EXPECT().ListOperations(gomock.Any(), testutil.EqProto(t, &buildqueuestate.ListOperationsRequest{
			PageSize: 2, StartAfter: &buildqueuestate.ListOperationsRequest_StartAfter{OperationName: "b"},
		})).Return(&buildqueuestate.ListOperationsResponse{Operations: []*buildqueuestate.OperationState{
			operation("a", "gcp"), operation("b", "gcp"),
		}}, nil),
	)

	_, err := router.ListOperations(selecting("ci/aws"), &buildqueuestate.ListOperationsRequest{PageSize: 2})
	testutil.RequireEqualStatus(t, status.Error(codes.Internal, "Scheduler pagination cursor did not advance"), err)
}

func TestConflictingSelectedSchedulerMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	aws := mock.NewMockBuildQueueStateClient(ctrl)
	gcp := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{
		route("ci/aws", "", aws),
		route("ci/gcp", "", gcp),
	})

	t.Run("ConflictingValuesAreRejected", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			SelectedInstanceNameMetadataKey, "ci/aws",
			SelectedInstanceNameMetadataKey, "ci/gcp",
		))
		_, err := router.ListPlatformQueues(ctx, &emptypb.Empty{})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("IdenticalValuesAreAllowed", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			SelectedInstanceNameMetadataKey, "ci/aws",
			SelectedInstanceNameMetadataKey, "ci/aws",
		))
		aws.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(platformQueues("mars"), nil)
		response, err := router.ListPlatformQueues(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		testutil.RequireEqualProto(t, platformQueues("ci/aws/mars"), response)
	})
}

// Guard against a later edit forgetting that the selection may carry the
// empty instance name, which is a valid route key.
func TestEmptySelectionSelectsEmptyPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	root := mock.NewMockBuildQueueStateClient(ctrl)
	other := mock.NewMockBuildQueueStateClient(ctrl)
	router := NewDemultiplexingClient([]Route{
		route("", "", root),
		route("ci", "ci", other),
	})
	root.EXPECT().ListPlatformQueues(gomock.Any(), gomock.Any()).Return(platformQueues(""), nil)

	response, err := router.ListPlatformQueues(selecting(""), &emptypb.Empty{})
	require.NoError(t, err)
	require.True(t, proto.Equal(platformQueues(""), response))
}
