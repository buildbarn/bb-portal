package integrationtest

import (
	"context"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/testkit"
	jmespath "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
)

type bepFile struct {
	filename     string
	invocationID string
}

type bepFileTestCase struct {
	bepFile bepFile
	wantErr error
}

type graphqlTestCase struct {
	variables testkit.Variables
	wantErr   error
	skip      bool
}

// testCases are grouped by operation name then by test name.
type graphqlTestTable map[string]map[string]graphqlTestCase

type dataExtractors struct {
	authMetadataExtractors      *bb_portal.AuthMetadataExtractorConfiguration
	invocationMetadataExtractor *jmespath.Expression
}

type testCase struct {
	name             string
	ctx              context.Context
	saveDataLevel    *bb_portal.BuildEventStreamService_SaveDataLevel
	dataExtractors   *dataExtractors
	buildKey         string
	mockUUID         *string
	bepFileTestCases []bepFileTestCase
	graphqlTestCases graphqlTestTable
}
