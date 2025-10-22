package integrationtest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	databasecommon "github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/internal/mock"
	"github.com/buildbarn/bb-portal/pkg/testkit"
	gql "github.com/machinebox/graphql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Run `bazel run //test/integrationtest:integrationtest_test -- --update-golden` to update golden files.

var updateGoldenFiles = flag.Bool("update-golden", false, "update golden (.golden.json) files")

const (
	bepFolderPath        = "testdata/bepfiles"
	goldenFolderPath     = "testdata/golden"
	consumerContractFile = "../../frontend/src/graphql/__generated__/persisted-documents.json"
)

type testCase struct {
	variables testkit.Variables
	wantErr   error
	skip      bool
}

// testCases are grouped by operation name then by test name.
type testTable map[string]map[string]testCase

func getTestCases() testTable {
	var (
		errInvocationNotFound = errors.New("graphql: ent: bazel_invocation not found")
		errBuildNotFound      = errors.New("graphql: ent: build not found")
	)

	const (
		successfulBazelBuild                                 = "fd03240f-697e-4b64-95bc-888e27445bf9" // nextjs_build.bep.ndjson
		successfulBazelTest                                  = "10a37e86-6e2b-4adb-83dd-c2906f42bdd6" // nextjs_test.bep.ndjson
		basicFailedInvocation                                = "08ae089d-4c85-405c-83fc-dbe9fc1dc942" // nextjs_build_fail.bep.ndjson
		invocationWithTargetAndErrorProgressWithActionOutput = "df7178e2-a815-4654-a409-d18e845d1e35" // nextjs_error_progress.bep.ndjson
		invocationWithAnalysisFailed                         = "571d0839-fd63-4442-bb4d-61f7bfa4ddae" // nextjs_test_fail.bep.ndjson
		invocationWithTargetFailed                           = ""                                     //
		invocationNotFound                                   = "4FF1C8C5-E51F-4ED1-8197-870FC389DA12" // uuidgen

		buildURLFound        = "https://example.com/build/1234"
		buildURLNotFound     = "https://example.com/build/4321"
		buildURLInstanceName = ""

		actionProblemID = "QWN0aW9uUHJvYmxlbTox"
		testResultID    = "VGVzdFJlc3VsdDoz"
	)

	return testTable{
		"LoadFullBazelInvocationDetails": {
			"get successful bazel build": {
				variables: testkit.Variables{
					"invocationID": successfulBazelBuild,
				},
			},
			"get successful bazel test": {
				variables: testkit.Variables{
					"invocationID": successfulBazelTest,
				},
			},
			"get single failed bazel invocation": {
				variables: testkit.Variables{
					"invocationID": basicFailedInvocation,
				},
			},
			"get single bazel invocation ignoring target and error progress if action has output": {
				variables: testkit.Variables{
					"invocationID": invocationWithTargetAndErrorProgressWithActionOutput,
				},
			},
			"get single bazel invocation analysis failed target": {
				variables: testkit.Variables{
					"invocationID": invocationWithAnalysisFailed,
				},
			},
			"get single bazel invocation failed target": {
				variables: testkit.Variables{
					"invocationID": invocationWithTargetFailed,
				},
				skip: true,
			},
			"bazel invocation not found": {
				variables: testkit.Variables{
					"invocationID": invocationNotFound,
				},
				wantErr: errInvocationNotFound,
			},
		},
		"FindBuildByUUID": {
			"found": {
				variables: testkit.Variables{
					"uuid": databasecommon.CalculateBuildUUID(buildURLFound, buildURLInstanceName),
				},
			},
			"not found": {
				variables: testkit.Variables{
					"uuid": databasecommon.CalculateBuildUUID(buildURLNotFound, buildURLInstanceName),
				},
				wantErr: errBuildNotFound,
			},
			"error when neither are specified": {
				wantErr: errors.New("buildUUID must be provided"),
			},
		},
		"FindBuilds": {
			"get all builds": {
				variables: testkit.Variables{
					"first": 1000,
					"where": nil,
				},
			},
		},
		"GetActionProblem": {
			"get action problem by ID output blob archiving queued": {
				variables: testkit.Variables{
					"id": actionProblemID,
				},
			},
		},
	}
}

// TestFromBesToGraphql tests the full flow from ingesting BEP files into the
// database to querying the data through the GraphQL API. It does so by loading
// all BEP files in the testdata/bepfiles folder into the database, and then
// running a series of queries against the GraphQL API, comparing the results
// against golden files stored in the testdata/golden folder.
func TestFromBesToGraphql(t *testing.T) {
	ctrl, ctx := gomock.WithContext(context.Background(), t)

	traceProvider := mock.NewMockTracerProvider(ctrl)
	tracer := mock.NewMockTracer(ctrl)
	traceProvider.BareMockTracerProvider.EXPECT().Tracer("github.com/buildbarn/bb-portal/internal/database/buildeventrecorder").Return(tracer).AnyTimes()
	span := mock.NewMockSpan(ctrl)
	tracer.BareMockTracer.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any()).Return(context.Background(), span).AnyTimes()
	span.BareMockSpan.EXPECT().End().AnyTimes()

	db := setupTestDB(t)
	bepUploader := setupTestBepUploader(t, db, traceProvider)

	// Read BEP file from testdata folder
	dirEntries, err := os.ReadDir(bepFolderPath)
	require.NoError(t, err)

	// Iterate over all files in the folder and load them one by one into the
	// database.
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		t.Run(fmt.Sprintf("SavingFileToDb_%s", entry.Name()), func(t *testing.T) {
			file, err := os.Open(bepFolderPath + "/" + entry.Name())
			require.NoError(t, err)
			_, _, err = bepUploader.RecordEventNdjsonFile(ctx, file)
			require.NoError(t, err)
			err = file.Close()
			require.NoError(t, err)
		})
	}

	graphqlServer := startGraphqlHTTPServer(t, db)

	queryRegistry := testkit.LoadQueryRegistry(t, "", consumerContractFile)

	for requestName, requestTestCases := range getTestCases() {
		t.Run(requestName, func(t *testing.T) {
			for testCaseName, testCase := range requestTestCases {
				t.Run(testCaseName, func(t *testing.T) {
					if testCase.skip {
						t.Skip("Test is skipped, needs a input fixture")
					}

					graphQLClient := gql.NewClient(graphqlServer.URL)
					req := queryRegistry.NewRequest(requestName)
					for k, v := range testCase.variables {
						req.Var(k, v)
					}

					var got map[string]interface{}
					err := graphQLClient.Run(ctx, req, &got)

					// Verify response.
					if testCase.wantErr != nil {
						require.Error(t, err)
						if !assert.Contains(t, err.Error(), testCase.wantErr.Error()) {
							require.NoError(t, err, "unexpected error received")
						}
						return
					}
					require.NoError(t, err)

					testkit.CheckAgainstGoldenFile(t, got, goldenFolderPath, requestName+"/"+testCaseName, updateGoldenFiles, &testkit.CompareOptions{DateTimeAgnostic: true})
				})
			}
		})
	}

	t.Run("No Unused Operations", func(t *testing.T) {
		t.Skip("WIP: Will add tests for other operations in future commits / PRs")
		require.Empty(t, queryRegistry.UnusedOperations())
	})
}
