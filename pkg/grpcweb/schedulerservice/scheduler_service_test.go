package schedulerservice_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/grpcweb/schedulerservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/program"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/auth"
	"github.com/buildbarn/bb-storage/pkg/proto/configuration/builder"
	grpc_pb "github.com/buildbarn/bb-storage/pkg/proto/configuration/grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// lazyClientFactory creates connections without dialing, like production.
type lazyClientFactory struct{}

func (lazyClientFactory) NewClientFromConfiguration(c *grpc_pb.ClientConfiguration, _ program.Group) (grpc.ClientConnInterface, error) {
	return grpc.NewClient(c.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func allow() *auth_pb.AuthorizerConfiguration {
	return &auth_pb.AuthorizerConfiguration{Policy: &auth_pb.AuthorizerConfiguration_Allow{Allow: &emptypb.Empty{}}}
}

func service(mutate func(*bb_portal.SchedulerService)) *bb_portal.SchedulerService {
	s := &bb_portal.SchedulerService{
		ListOperationsPageSize:   500,
		ReadAuthorizer:           allow(),
		KillOperationsAuthorizer: allow(),
	}
	mutate(s)
	return s
}

func scheduler(address, addInstanceNamePrefix string) *builder.SchedulerConfiguration {
	return &builder.SchedulerConfiguration{
		Endpoint:              &grpc_pb.ClientConfiguration{Address: address},
		AddInstanceNamePrefix: addInstanceNamePrefix,
	}
}

func catalog(t *testing.T, router *http.ServeMux) string {
	t.Helper()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/schedulers", nil))
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "no-store", recorder.Header().Get("Cache-Control"))
	body, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)
	return string(body)
}

func TestCatalogListsInstanceNamePrefixes(t *testing.T) {
	t.Run("SingleClientServesEveryInstanceName", func(t *testing.T) {
		router := http.NewServeMux()
		require.NoError(t, schedulerservice.NewSchedulerService(service(func(s *bb_portal.SchedulerService) {
			s.BuildQueueStateClient = &grpc_pb.ClientConfiguration{Address: "scheduler:8984"}
		}), nil, nil, lazyClientFactory{}, router))
		require.JSONEq(t, `{"schedulers":[{"instanceNamePrefix":""}]}`, catalog(t, router))
	})

	t.Run("SchedulersSortedByPrefix", func(t *testing.T) {
		router := http.NewServeMux()
		require.NoError(t, schedulerservice.NewSchedulerService(service(func(s *bb_portal.SchedulerService) {
			s.Schedulers = map[string]*builder.SchedulerConfiguration{
				"ci/gcp": scheduler("gcp:8984", "ci/gcp"),
				"ci/aws": scheduler("aws:8984", ""),
				"":       scheduler("default:8984", ""),
			}
		}), nil, nil, lazyClientFactory{}, router))
		require.JSONEq(t, `{"schedulers":[{"instanceNamePrefix":""},{"instanceNamePrefix":"ci/aws"},{"instanceNamePrefix":"ci/gcp"}]}`, catalog(t, router))
	})
}

func TestRejectsInvalidConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name    string
		mutate  func(*bb_portal.SchedulerService)
		message string
	}{
		{
			name:    "NoSchedulers",
			mutate:  func(*bb_portal.SchedulerService) {},
			message: "No build_queue_state_client or schedulers configured",
		},
		{
			name: "ClientCombinedWithSchedulers",
			mutate: func(s *bb_portal.SchedulerService) {
				s.BuildQueueStateClient = &grpc_pb.ClientConfiguration{Address: "scheduler:8984"}
				s.Schedulers = map[string]*builder.SchedulerConfiguration{"ci": scheduler("ci:8984", "ci")}
			},
			message: "build_queue_state_client and schedulers cannot be combined",
		},
		{
			name: "InvalidPrefix",
			mutate: func(s *bb_portal.SchedulerService) {
				s.Schedulers = map[string]*builder.SchedulerConfiguration{"ci//aws": scheduler("aws:8984", "")}
			},
			message: `Invalid instance name "ci//aws"`,
		},
		{
			name: "InvalidAddedPrefix",
			mutate: func(s *bb_portal.SchedulerService) {
				s.Schedulers = map[string]*builder.SchedulerConfiguration{"ci": scheduler("ci:8984", "/ci")}
			},
			message: `Invalid instance name "/ci"`,
		},
		{
			name: "MissingEndpoint",
			mutate: func(s *bb_portal.SchedulerService) {
				s.Schedulers = map[string]*builder.SchedulerConfiguration{"ci": {AddInstanceNamePrefix: "ci"}}
			},
			message: `No endpoint configured for instance name "ci"`,
		},
		{
			name: "MissingReadAuthorizer",
			mutate: func(s *bb_portal.SchedulerService) {
				s.BuildQueueStateClient = &grpc_pb.ClientConfiguration{Address: "scheduler:8984"}
				s.ReadAuthorizer = nil
			},
			message: "No ReadAuthorizer configured",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := schedulerservice.NewSchedulerService(service(tc.mutate), nil, nil, lazyClientFactory{}, http.NewServeMux())
			require.ErrorContains(t, err, tc.message)
		})
	}
}
