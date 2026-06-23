package integrationtest

import (
	"context"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/testkit"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/google/go-cmp/cmp"
	gql "github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestArtifactGraphEndToEnd uploads the artifact fixture at the
// basicAndTargetAndArtifacts save level, fetches the resulting
// artifactGraph via GraphQL, and asserts the decoded structure matches
// the expected NamedSetOfFiles + TargetCompleted shape for the fixture
// build.
//
// Compaction (dbcleanupservice) does not run in this test, so the graph is
// served from the incomplete_artifact_graphs staging table via the
// resolver's partial-state path. The server decodes the events and returns
// structured data, so the test walks the named-set graph directly rather
// than parsing a wire format.
func TestArtifactGraphEndToEnd(t *testing.T) {
	ctx := context.Background()

	tc := testCase{
		name: "TestArtifactGraphEndToEnd",
		saveDataLevel: &bb_portal.BuildEventStreamService_SaveDataLevel{
			Level: &bb_portal.BuildEventStreamService_SaveDataLevel_BasicAndTargetAndArtifacts{
				BasicAndTargetAndArtifacts: &emptypb.Empty{},
			},
		},
	}

	db := testutils.SetupTestDB(t, dbProvider)
	uploader := setupTestBepUploader(t, db, tc)

	file, err := os.Open(bepFolderPath + "/" + artifactsEndToEndBuild.filename)
	require.NoError(t, err)
	t.Cleanup(func() { _ = file.Close() })
	_, _, err = uploader.RecordEventNdjsonFile(ctx, file)
	require.NoError(t, err)

	server := startGraphqlHTTPServer(t, db)
	queryRegistry := testkit.LoadQueryRegistry(t, "", consumerContractFile)

	resp := runArtifactGraphQuery(ctx, t, server, queryRegistry, artifactsEndToEndBuild.invocationID)
	graph := resp.GetBazelInvocation.ArtifactGraph
	require.NotNil(t, graph, "artifactGraph missing from response")

	targets := resolveTargets(graph)

	wantTargets := []targetEntry{
		{
			label: "//:hello",
			outputGroups: []outputGroupEntry{
				{name: "default", fileNames: []string{"hello"}},
			},
		},
	}

	if diff := cmp.Diff(wantTargets, targets, cmp.AllowUnexported(targetEntry{}, outputGroupEntry{})); diff != "" {
		t.Fatalf("decoded targets mismatch (-want +got):\n%s", diff)
	}
	if len(graph.NamedSets) == 0 {
		t.Fatalf("expected at least one NamedSetOfFiles in graph; got 0")
	}
}

type artifactFileResp struct {
	Name        string  `json:"name"`
	URI         *string `json:"uri"`
	Digest      *string `json:"digest"`
	SizeBytes   *int    `json:"sizeBytes"`
	DownloadURL *string `json:"downloadUrl"`
}

type artifactNamedSetResp struct {
	ID          string             `json:"id"`
	ChildSetIds []string           `json:"childSetIds"`
	Files       []artifactFileResp `json:"files"`
}

type artifactOutputGroupResp struct {
	Name       string   `json:"name"`
	Incomplete bool     `json:"incomplete"`
	RootSetIds []string `json:"rootSetIds"`
}

type artifactTargetResp struct {
	Label        string                    `json:"label"`
	Aspect       *string                   `json:"aspect"`
	OutputGroups []artifactOutputGroupResp `json:"outputGroups"`
}

type artifactGraphResp struct {
	NamedSets []artifactNamedSetResp `json:"namedSets"`
	Targets   []artifactTargetResp   `json:"targets"`
}

type artifactGraphResponse struct {
	GetBazelInvocation struct {
		ID            string             `json:"id"`
		ArtifactGraph *artifactGraphResp `json:"artifactGraph"`
	} `json:"getBazelInvocation"`
}

func runArtifactGraphQuery(
	ctx context.Context,
	t *testing.T,
	server *httptest.Server,
	queryRegistry *testkit.QueryRegistry,
	invocationID string,
) artifactGraphResponse {
	t.Helper()
	client := gql.NewClient(server.URL)
	req := queryRegistry.NewRequest("ArtifactGraph")
	req.Var("id", invocationID)
	var got artifactGraphResponse
	require.NoError(t, client.Run(ctx, req, &got))
	return got
}

type outputGroupEntry struct {
	name      string
	fileNames []string
}

type targetEntry struct {
	label        string
	outputGroups []outputGroupEntry
}

// resolveTargets walks each target's output groups through the named-set
// graph and returns the resolved file names, sorted for stable comparison.
func resolveTargets(graph *artifactGraphResp) []targetEntry {
	sets := make(map[string]artifactNamedSetResp, len(graph.NamedSets))
	for _, s := range graph.NamedSets {
		sets[s.ID] = s
	}

	targets := make([]targetEntry, 0, len(graph.Targets))
	for _, t := range graph.Targets {
		entry := targetEntry{label: t.Label}
		for _, og := range t.OutputGroups {
			files := walkFileNames(sets, og.RootSetIds)
			sort.Strings(files)
			entry.outputGroups = append(entry.outputGroups, outputGroupEntry{
				name:      og.Name,
				fileNames: files,
			})
		}
		sort.Slice(entry.outputGroups, func(i, j int) bool {
			return entry.outputGroups[i].name < entry.outputGroups[j].name
		})
		targets = append(targets, entry)
	}
	sort.Slice(targets, func(i, j int) bool { return targets[i].label < targets[j].label })
	return targets
}

func walkFileNames(sets map[string]artifactNamedSetResp, roots []string) []string {
	visited := map[string]struct{}{}
	queue := append([]string(nil), roots...)
	var out []string
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if _, seen := visited[id]; seen {
			continue
		}
		visited[id] = struct{}{}
		set, ok := sets[id]
		if !ok {
			continue
		}
		for _, f := range set.Files {
			out = append(out, f.Name)
		}
		queue = append(queue, set.ChildSetIds...)
	}
	return out
}
