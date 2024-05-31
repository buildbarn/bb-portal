package graphql_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	gql "github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/buildbarn/bb-portal/internal/graphql/helpers"
	"github.com/buildbarn/bb-portal/pkg/testkit"
)

var (
	errInvocationNotFound = errors.New("graphql: ent: bazel_invocation not found")
	errBuildNotFound      = errors.New("graphql: ent: build not found")
)

const (
	successfulBazelBuild                                 = "fd03240f-697e-4b64-95bc-888e27445bf9" // nextjs_build.bep.ndjson
	successfulBazelTest                                  = "10a37e86-6e2b-4adb-83dd-c2906f42bdd6" // nextjs_test.bep.ndjson
	basicFailedInvocation                                = "08ae089d-4c85-405c-83fc-dbe9fc1dc942" // nextjs_build_fail.bep.ndjson
	invocationWithTargetAndErrorProgressWithActionOutput = "df7178e2-a815-4654-a409-d18e845d1e35" // nextjs_error_progress.bep.ndjson ?
	invocationWithAnalysisFailed                         = "571d0839-fd63-4442-bb4d-61f7bfa4ddae" // nextjs_test_fail.bep.ndjson
	invocationWithTargetFailed                           = ""                                     //
	invocationNotFound                                   = "4FF1C8C5-E51F-4ED1-8197-870FC389DA12" // uuidgen

	buildURLFound    = "https://example.com/build/1234"
	buildURLNotFound = "https://example.com/build/4321"

	actionProblemID = "QWN0aW9uUHJvYmxlbTox"
	testResultID    = "VGVzdFJlc3VsdDoz"
)

type apiTestCase struct {
	variables     testkit.Variables
	entDataSource string
	setupEnt      setupFuncEnt
	wantErr       error
	skip          bool
}

// testCases are grouped by operation name then by test name.
type testTable map[string]map[string]apiTestCase

func TestGraphQLAPI_Snapshots(t *testing.T) {
	queryRegistry := testkit.LoadQueryRegistry(t, queryDir, consumerContractFile)
	testCases := testTable{
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
			"found (by URL)": {
				variables: testkit.Variables{
					"url": buildURLFound,
				},
			},
			"found (by UUID)": {
				variables: testkit.Variables{
					"uuid": uuid.NewSHA1(uuid.NameSpaceURL, []byte(buildURLFound)),
				},
			},
			"not found (by URL)": {
				variables: testkit.Variables{
					"url": buildURLNotFound,
				},
				wantErr: errBuildNotFound,
			},
			"not found (by UUID)": {
				variables: testkit.Variables{
					"url": uuid.NewSHA1(uuid.NameSpaceURL, []byte(buildURLNotFound)),
				},
				wantErr: errBuildNotFound,
			},
			"error when both specified": {
				variables: testkit.Variables{
					"url":  buildURLFound,
					"uuid": uuid.NewSHA1(uuid.NameSpaceURL, []byte(buildURLFound)),
				},
				wantErr: helpers.ErrOnlyURLOrUUID,
			},
			"error when neither are specified": {
				wantErr: helpers.ErrOnlyURLOrUUID,
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

	for opName, opTestCases := range testCases {
		t.Run(opName, func(t *testing.T) {
			for name, testCase := range opTestCases {
				opName, name, testCase := opName, name, testCase // Fixes: Using the variable on range scope `x` in function literal (scopelint)
				t.Run(name, func(t *testing.T) {
					if testCase.skip {
						t.Skip("Test is skipped, needs a input fixture")
					}
					mockServer := newMockServer(t, testCase.entDataSource)
					if testCase.setupEnt != nil {
						testCase.setupEnt(mockServer.ctx, mockServer.client)
					}

					graphQLClient := gql.NewClient(mockServer.URL + "/graphql")
					req := queryRegistry.NewRequest(opName)
					for k, v := range testCase.variables {
						req.Var(k, v)
					}

					var got map[string]interface{}
					err := graphQLClient.Run(mockServer.ctx, req, &got)

					// Verify response.
					if testCase.wantErr != nil {
						require.Error(t, err)
						if !assert.Contains(t, err.Error(), testCase.wantErr.Error()) {
							require.NoError(t, err, "unexpected error received")
						}
						return
					}
					require.NoError(t, err)

					testkit.CheckAgainstGoldenFile(t, got, snapshotDir, opName+"/"+name, update, &testkit.CompareOptions{DateTimeAgnostic: true})
				})
			}
		})
	}

	t.Run("No Unused Operations", func(t *testing.T) {
		t.Skip("WIP: Will add tests for other operations in future commits / PRs")
		require.Empty(t, queryRegistry.UnusedOperations())
	})
}
