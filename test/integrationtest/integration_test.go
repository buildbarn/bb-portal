package integrationtest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	databasecommon "github.com/buildbarn/bb-portal/internal/database/common"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/testkit"
	"github.com/buildbarn/bb-storage/pkg/auth"
	jmespath "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"github.com/buildbarn/bb-storage/pkg/util"
	gql "github.com/machinebox/graphql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Run `bazel run //test/integrationtest:integrationtest_test -- --update-golden` to update golden files.

const (
	bepFolderPath        = "testdata/bepfiles"
	goldenFolderPath     = "testdata/golden"
	consumerContractFile = "../../frontend/src/graphql/__generated__/persisted-documents.json"
)

var (
	updateGoldenFiles = flag.Bool("update-golden", false, "update golden (.golden.json) files")

	errInvocationNotFound        = errors.New("graphql: ent: bazel_invocation not found")
	errBuildNotFound             = errors.New("graphql: ent: build not found")
	errInvocationLocked          = errors.New("already exists and is locked for writing")
	errAuthenticatedUserNotFound = errors.New("ent: authenticated_user not found")

	// Defined BEP files for various test cases.
	successfulBazelBuild = bepFile{
		filename:     "nextjs_build.bep.ndjson",
		invocationID: "fd03240f-697e-4b64-95bc-888e27445bf9",
	}
	failedBazelBuild = bepFile{
		filename:     "nextjs_build_fail.bep.ndjson",
		invocationID: "08ae089d-4c85-405c-83fc-dbe9fc1dc942",
	}
	successfulBazelTest = bepFile{
		filename:     "nextjs_test.bep.ndjson",
		invocationID: "10a37e86-6e2b-4adb-83dd-c2906f42bdd6",
	}
	failedBazelTest = bepFile{
		filename:     "nextjs_test_fail.bep.ndjson",
		invocationID: "571d0839-fd63-4442-bb4d-61f7bfa4ddae",
	}
	failedBazelAnalysis = bepFile{
		filename:     "nextjs_analysis_fail.bep.ndjson",
		invocationID: "df7178e2-a815-4654-a409-d18e845d1e35",
	}
	authenticatedUserUUID = "8bdb3187-e36c-487e-95b8-f8ca28a82068"

	// An authenticated user UUID not present in any BEP file.
	authenticatedUserUUIDNotFound = "A80031E0-1A11-4543-894C-13C48056074A"

	// An invocation ID that is not present in any BEP file.
	invocationIDNotFound = "4FF1C8C5-E51F-4ED1-8197-870FC389DA12" // uuidgen

	// testCases defines all integration test cases to be run.
	testCases = []testCase{
		{
			name: "TestAllBepUploadsAndGraphqlQueries",
			saveTargetDataLevel: &bb_portal.BuildEventStreamService_SaveTargetDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTargetDataLevel_Enriched{
					Enriched: &emptypb.Empty{},
				},
			},
			saveTestDataLevel: &bb_portal.BuildEventStreamService_SaveTestDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTestDataLevel_Enriched{
					Enriched: &emptypb.Empty{},
				},
			},
			bepFileTestCases: []bepFileTestCase{
				{bepFile: successfulBazelBuild},
				{bepFile: failedBazelBuild},
				{bepFile: successfulBazelTest},
				{bepFile: failedBazelTest},
				{bepFile: failedBazelAnalysis},
			},
			graphqlTestCases: graphqlTestTable{
				"LoadFullBazelInvocationDetails": {
					"get successful bazel build": {
						variables: testkit.Variables{
							"invocationID": successfulBazelBuild.invocationID,
						},
					},
					"get successful bazel test": {
						variables: testkit.Variables{
							"invocationID": successfulBazelTest.invocationID,
						},
					},
					"get single failed bazel invocation": {
						variables: testkit.Variables{
							"invocationID": failedBazelBuild.invocationID,
						},
					},
					"get single bazel invocation ignoring target and error progress if action has output": {
						variables: testkit.Variables{
							"invocationID": failedBazelAnalysis.invocationID,
						},
					},
					"get single bazel invocation analysis failed target": {
						variables: testkit.Variables{
							"invocationID": failedBazelTest.invocationID,
						},
					},
					"bazel invocation not found": {
						variables: testkit.Variables{
							"invocationID": invocationIDNotFound,
						},
						wantErr: errInvocationNotFound,
					},
				},
				"FindBuildByUUID": {
					"found": {
						variables: testkit.Variables{
							"uuid": databasecommon.CalculateBuildUUID("https://example.com/build/1234", ""),
						},
					},
					"not found": {
						variables: testkit.Variables{
							"uuid": databasecommon.CalculateBuildUUID("https://example.com/build/4321", ""),
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
				"GetInvocationTargetsForInvocation": {
					"get targets for successfull build": {
						variables: testkit.Variables{
							"invocationID": successfulBazelBuild.invocationID,
						},
					},
					"get targets for failed analysis": {
						variables: testkit.Variables{
							"invocationID": failedBazelAnalysis.invocationID,
						},
					},
				},
				"GetTargetDetails": {
					"get details for a target": {
						variables: testkit.Variables{
							"instanceName": "",
							"label":        "//packages/one:one",
							"aspect":       "",
							"targetKind":   "_npm_package rule",
							"where": map[string]interface{}{
								"hasTargetWith": []interface{}{
									map[string]interface{}{
										"hasInstanceNameWith": map[string]interface{}{
											"name": "",
										},
										"label":      "//packages/one:one",
										"aspect":     "",
										"targetKind": "_npm_package rule",
									},
								},
							},
						},
					},
				},
				"GetTargetsList": {
					"get all targets": {
						variables: testkit.Variables{
							"first": 1000,
						},
					},
				},
				"GetInvocationTargetsForTarget": {
					"get invocation targets for target": {
						variables: testkit.Variables{
							"instanceName": "",
							"label":        "//packages/one:one",
							"aspect":       "",
							"targetKind":   "_npm_package rule",
						},
					},
				},
			},
		},
		{
			name: "TestDuplicateBepUploads",
			saveTargetDataLevel: &bb_portal.BuildEventStreamService_SaveTargetDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTargetDataLevel_Enriched{
					Enriched: &emptypb.Empty{},
				},
			},
			saveTestDataLevel: &bb_portal.BuildEventStreamService_SaveTestDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTestDataLevel_Enriched{
					Enriched: &emptypb.Empty{},
				},
			},
			bepFileTestCases: []bepFileTestCase{
				{bepFile: successfulBazelBuild},
				{bepFile: successfulBazelBuild, wantErr: errInvocationLocked},
			},
			graphqlTestCases: graphqlTestTable{},
		},
		{
			name: "TestGraphqlQueriesWithAuthMetadata",
			saveTargetDataLevel: &bb_portal.BuildEventStreamService_SaveTargetDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTargetDataLevel_Enriched{
					Enriched: &emptypb.Empty{},
				},
			},
			saveTestDataLevel: &bb_portal.BuildEventStreamService_SaveTestDataLevel{
				Level: &bb_portal.BuildEventStreamService_SaveTestDataLevel_Enriched{
					Enriched: &emptypb.Empty{},
				},
			},
			bepFileTestCases: []bepFileTestCase{
				{bepFile: successfulBazelBuild},
			},
			extractors: &bb_portal.AuthMetadataExtractorConfiguration{
				ExternalIdExtractionJmespathExpression:  &jmespath.Expression{Expression: "authenticationMetadata.private.external_id"},
				DisplayNameExtractionJmespathExpression: &jmespath.Expression{Expression: "authenticationMetadata.private.display_name"},
				UserInfoExtractionJmespathExpression:    &jmespath.Expression{Expression: "authenticationMetadata.public"},
			},
			ctx: auth.NewContextWithAuthenticationMetadata(context.Background(), util.Must(auth.NewAuthenticationMetadataFromRaw(map[string]any{
				"private": map[string]any{
					"external_id":  "6b939a5f-95b9-4a7c-ad89-961fb96c5cd1",
					"display_name": "example_username",
				},
				"public": map[string]any{
					"age": 30,
					"contact_information": map[string]any{
						"email":        "user@example.com",
						"phone_number": "800-555-0199",
					},
				},
			}))),
			mockUUID: &authenticatedUserUUID,
			graphqlTestCases: graphqlTestTable{
				"GetAuthenticatedUser": {
					"authenticated user not found": {
						variables: testkit.Variables{
							"userUUID": &authenticatedUserUUIDNotFound,
						},
						wantErr: errAuthenticatedUserNotFound,
					},
					"authenticated user found": {
						variables: testkit.Variables{
							"userUUID": &authenticatedUserUUID,
						},
					},
				},
			},
		},
	}
)

// TestFromBesToGraphql tests the full flow from ingesting BEP files into the
// database to querying the data through the GraphQL API. It does so by loading
// all BEP files in the testdata/bepfiles folder into the database, and then
// running a series of queries against the GraphQL API, comparing the results
// against golden files stored in the testdata/golden folder.
func TestFromBesToGraphql(t *testing.T) {
	queryRegistry := testkit.LoadQueryRegistry(t, "", consumerContractFile)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runTestCase(t, queryRegistry, testCase)
		})
	}

	t.Run("No Unused Operations", func(t *testing.T) {
		t.Skip("WIP: Will add tests for other operations in future commits / PRs")
		require.Empty(t, queryRegistry.UnusedOperations())
	})
}

func runTestCase(t *testing.T, queryRegistry *testkit.QueryRegistry, testCase testCase) {
	ctx := testCase.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	var uuidGenerator util.UUIDGenerator
	if testCase.mockUUID != nil {
		uuidGenerator = createMockUUIDGenerator(t, *testCase.mockUUID, len(testCase.bepFileTestCases))
	} else {
		uuidGenerator = uuid.NewRandom
	}

	db := setupTestDB(t)

	bepUploader := setupTestBepUploader(t, db, testCase, uuidGenerator)

	for _, bepFileTestCase := range testCase.bepFileTestCases {
		t.Run(fmt.Sprintf("SavingFileToDb_%s", bepFileTestCase.bepFile.filename), func(t *testing.T) {
			file, err := os.Open(bepFolderPath + "/" + bepFileTestCase.bepFile.filename)
			require.NoError(t, err)
			_, _, err = bepUploader.RecordEventNdjsonFile(ctx, file)
			checkIfErrorMatches(t, bepFileTestCase.wantErr, err)
			require.NoError(t, file.Close())
		})
	}

	graphqlServer := startGraphqlHTTPServer(t, db)

	runGraphqlTestCases(ctx, t, graphqlServer.URL, queryRegistry, testCase)
}

func runGraphqlTestCases(ctx context.Context, t *testing.T, graphqlServerURL string, queryRegistry *testkit.QueryRegistry, testCase testCase) {
	for gqlRequestName, gqlTestCases := range testCase.graphqlTestCases {
		t.Run(gqlRequestName, func(t *testing.T) {
			for gqlTestCaseName, gqlTestCase := range gqlTestCases {
				t.Run(gqlTestCaseName, func(t *testing.T) {
					if gqlTestCase.skip {
						t.Skip("Test is skipped, needs a input fixture")
					}

					graphQLClient := gql.NewClient(graphqlServerURL)
					req := queryRegistry.NewRequest(gqlRequestName)
					for k, v := range gqlTestCase.variables {
						req.Var(k, v)
					}

					var got map[string]interface{}
					err := graphQLClient.Run(ctx, req, &got)

					checkIfErrorMatches(t, gqlTestCase.wantErr, err)
					if gqlTestCase.wantErr != nil {
						return
					}

					testkit.CheckAgainstGoldenFile(t, got, goldenFolderPath, testCase.name+"/"+gqlRequestName+"/"+gqlTestCaseName, updateGoldenFiles, &testkit.CompareOptions{DateTimeAgnostic: true})
				})
			}
		})
	}
}
