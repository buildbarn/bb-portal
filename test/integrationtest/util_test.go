package integrationtest

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/buildbarn/bb-portal/internal/api/http/bepuploader"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/internal/graphql"
	"github.com/buildbarn/bb-portal/pkg/invocation/files"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/blobstore"
	"github.com/buildbarn/bb-storage/pkg/digest"
	jmespath "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var dbProvider *embedded.DatabaseProvider

type fixtureContentAddressableStorage struct {
	blobstore.BlobAccess
	availableDigests map[digest.Digest]struct{}
}

func newFixtureContentAddressableStorage(t *testing.T) *fixtureContentAddressableStorage {
	t.Helper()
	availableURIs := []string{
		"bytestream://cache.example.com/blobs/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/42",
		"bytestream://cache.example.com/blobs/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb/12",
		"bytestream://cache.example.com/blobs/cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc/84",
		"bytestream://cache.example.com/blobs/dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd/24",
		"bytestream://cache.example.com/blobs/eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee/36",
		"bytestream://cache.example.com/blobs/ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff/48",
	}
	availableDigests := make(map[digest.Digest]struct{}, len(availableURIs))
	for _, uri := range availableURIs {
		parsedURI, err := url.Parse(uri)
		require.NoError(t, err)
		require.Equal(t, "cache.example.com", parsedURI.Host)
		parsedDigest := files.GetDigestFromURI(uri)
		require.NotEqual(t, digest.BadDigest, parsedDigest)
		availableDigests[parsedDigest] = struct{}{}
	}
	return &fixtureContentAddressableStorage{availableDigests: availableDigests}
}

func (cas *fixtureContentAddressableStorage) FindMissing(_ context.Context, digests digest.Set) (digest.Set, error) {
	missingDigests := digest.NewSetBuilder(digests.Length())
	for _, blobDigest := range digests.Items() {
		if _, ok := cas.availableDigests[blobDigest]; !ok {
			missingDigests = missingDigests.Add(blobDigest)
		}
	}
	return missingDigests.Build(), nil
}

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = embedded.NewDatabaseProvider(os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start embedded DB: %v\n", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	if err := dbProvider.Cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not clean up embedded DB: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func setupTestBepUploader(t *testing.T, db database.Client, testCase testCase) *bepuploader.BepUploader {
	var authExtractors *bb_portal.AuthMetadataExtractorConfiguration
	var invocationExtractor *jmespath.Expression
	if testCase.dataExtractors != nil {
		authExtractors = testCase.dataExtractors.authMetadataExtractors
		invocationExtractor = testCase.dataExtractors.invocationMetadataExtractor
	}

	authorizer := auth.NewStaticAuthorizer(func(in digest.InstanceName) bool { return true })
	besConfig := &bb_portal.BuildEventStreamService{
		SaveDataLevel:                testCase.saveDataLevel,
		AuthMetadataKeyConfiguration: authExtractors,
		InvocationMetadataExtractor:  invocationExtractor,
		BuildKey:                     testCase.buildKey,
	}

	bepUploader, err := bepuploader.NewBepUploader(db, besConfig, newFixtureContentAddressableStorage(t), authorizer, nil, nil, noop.NewTracerProvider())
	require.NoError(t, err)
	return bepUploader
}

func startGraphqlHTTPServer(t *testing.T, db database.Client) *httptest.Server {
	srv := graphql.NewGraphqlHandler(db, trace.NewNoopTracerProvider())

	// Bypass DB auth service for integration tests.
	srv.AroundOperations(func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
		return next(dbauthservice.NewContextWithDbAuthServiceBypass(ctx))
	})

	server := httptest.NewServer(srv)
	t.Cleanup(func() { server.Close() })
	return server
}

func checkIfErrorMatches(t *testing.T, wantErr, err error) {
	if wantErr != nil {
		require.Error(t, err)
		if !assert.Contains(t, err.Error(), wantErr.Error()) {
			require.NoError(t, err, "unexpected error received")
		}
	} else {
		require.NoError(t, err)
	}
}

func githubActionsExtractor() *jmespath.Expression {
	s := ""

	// This was the easiest way to build this string. We cannot use multiline
	// strings since it contains backticks
	s += "{"
	s += "  \"username\": env.USER"
	s += "  \"hostname\": env.HOSTNAME"
	s += "  \"sourceControls\": ["
	s += "    {"
	s += "      \"repo\": env.GITHUB_REPOSITORY"
	s += "      \"repoUrl\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY) && join('/', [env.GITHUB_SERVER_URL , env.GITHUB_REPOSITORY]) || `null`"
	s += "      \"ref\": env.GITHUB_REF"
	s += "      \"refUrl\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY && env.GITHUB_REF) && join('/', [env.GITHUB_SERVER_URL, env.GITHUB_REPOSITORY, 'tree', env.GITHUB_REF]) || `null`"
	s += "      \"commit\": env.GITHUB_SHA"
	s += "      \"commitUrl\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY && env.GITHUB_SHA) && join('/', [env.GITHUB_SERVER_URL, env.GITHUB_REPOSITORY, 'commit', env.GITHUB_SHA]) || `null`"
	s += "    }"
	s += "  ]"
	s += "  \"invocationTags\": {"
	s += "    \"workflow\": env.GITHUB_WORKFLOW"
	s += "    \"workflow_url\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY && env.GITHUB_RUN_ID) && join('/', [env.GITHUB_SERVER_URL , env.GITHUB_REPOSITORY, 'actions', 'runs', env.GITHUB_RUN_ID]) || `null`"
	s += "    \"job\": env.GITHUB_JOB"
	s += "    \"action\": env.GITHUB_ACTION"
	s += "  }"
	s += "  \"buildTags\": {"
	s += "    \"repo\": env.GITHUB_REPOSITORY"
	s += "    \"repo_url\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY) && join('/', [env.GITHUB_SERVER_URL , env.GITHUB_REPOSITORY]) || `null`"
	s += "    \"workflow\": env.GITHUB_WORKFLOW"
	s += "    \"workflow_url\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY && env.GITHUB_RUN_ID) && join('/', [env.GITHUB_SERVER_URL , env.GITHUB_REPOSITORY, 'actions', 'runs', env.GITHUB_RUN_ID]) || `null`"
	s += "    \"build_id\": (env.GITHUB_SERVER_URL && env.GITHUB_REPOSITORY && env.GITHUB_RUN_ID) && join('/', [env.GITHUB_SERVER_URL , env.GITHUB_REPOSITORY, 'actions', 'runs', env.GITHUB_RUN_ID]) || `null`"
	s += "  }"
	s += "}"

	expr := jmespath.Expression{Expression: s}
	return &expr
}
