package integrationtest

import (
	"context"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/testkit"
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

type testCase struct {
	name                string
	ctx                 context.Context
	saveTargetDataLevel *bb_portal.BuildEventStreamService_SaveTargetDataLevel
	saveTestDataLevel   *bb_portal.BuildEventStreamService_SaveTestDataLevel
	extractors          *bb_portal.AuthMetadataExtractorConfiguration
	mockUUID            *string
	bepFileTestCases    []bepFileTestCase
	graphqlTestCases    graphqlTestTable
}
