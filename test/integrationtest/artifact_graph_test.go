package integrationtest

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"net/http/httptest"
	"os"
	"sort"
	"testing"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/testkit"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/google/go-cmp/cmp"
	gql "github.com/machinebox/graphql"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestArtifactGraphEndToEnd uploads the artifact fixture at the
// basicAndTargetAndArtifacts save level, fetches the resulting
// artifactGraph via GraphQL, decodes the zstd-compressed BEP-graph
// blob, and asserts the decoded structure matches the expected
// NamedSetOfFiles + TargetCompleted shape for the fixture build.
//
// This test deliberately avoids the golden-file framework because the
// zstd-encoded payload is not byte-stable across encoder versions; we
// compare the decoded structure instead.
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

	got := runArtifactGraphQuery(ctx, t, server, queryRegistry, artifactsEndToEndBuild.invocationID)
	payload := decodeArtifactGraphPayload(t, got)

	targets, sets := parseArtifactGraphEvents(t, payload)

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
	if len(sets) == 0 {
		t.Fatalf("expected at least one NamedSetOfFiles in graph; got 0")
	}
}

func runArtifactGraphQuery(
	ctx context.Context,
	t *testing.T,
	server *httptest.Server,
	queryRegistry *testkit.QueryRegistry,
	invocationID string,
) map[string]interface{} {
	t.Helper()
	client := gql.NewClient(server.URL)
	req := queryRegistry.NewRequest("ArtifactGraph")
	req.Var("id", invocationID)
	var got map[string]interface{}
	require.NoError(t, client.Run(ctx, req, &got))
	return got
}

func decodeArtifactGraphPayload(t *testing.T, got map[string]interface{}) []byte {
	t.Helper()
	inv, ok := got["getBazelInvocation"].(map[string]interface{})
	require.True(t, ok, "getBazelInvocation missing from response: %v", got)
	graph, ok := inv["artifactGraph"].(map[string]interface{})
	require.True(t, ok, "artifactGraph missing from response: %v", inv)
	payloadStr, ok := graph["payload"].(string)
	require.True(t, ok, "payload missing from artifactGraph: %v", graph)
	compressed, err := base64.StdEncoding.DecodeString(payloadStr)
	require.NoError(t, err)
	decoder, err := zstd.NewReader(nil)
	require.NoError(t, err)
	defer decoder.Close()
	decompressed, err := decoder.DecodeAll(compressed, nil)
	require.NoError(t, err)
	return decompressed
}

type outputGroupEntry struct {
	name      string
	fileNames []string
}

type targetEntry struct {
	label        string
	outputGroups []outputGroupEntry
}

// parseArtifactGraphEvents walks length-prefixed serialized bes.BuildEvent
// messages and returns the targets seen (with their output groups and
// resolved file names) plus the set of NamedSetOfFiles set IDs encountered.
//
// Two-pass: the recorder may emit TargetCompleted events to the buffer
// before the NamedSetOfFiles events they reference (saveBatch processes
// TargetCompleted as a dedicated batch step, NamedSetOfFiles as part of
// the catch-all step). We collect all sets in pass 1, then walk files in
// pass 2.
func parseArtifactGraphEvents(t *testing.T, stream []byte) ([]targetEntry, map[string]*bes.NamedSetOfFiles) {
	t.Helper()
	sets := map[string]*bes.NamedSetOfFiles{}
	type pendingTarget struct {
		label  string
		groups []struct {
			name     string
			rootIDs  []string
		}
	}
	var pending []pendingTarget

	offset := 0
	for offset < len(stream) {
		size, n := binary.Uvarint(stream[offset:])
		require.Greater(t, n, 0, "invalid varint at offset %d", offset)
		offset += n
		require.LessOrEqual(t, offset+int(size), len(stream), "truncated BuildEvent at offset %d", offset)
		evt := &bes.BuildEvent{}
		require.NoError(t, proto.Unmarshal(stream[offset:offset+int(size)], evt))
		offset += int(size)

		switch id := evt.GetId().GetId().(type) {
		case *bes.BuildEventId_NamedSet:
			sets[id.NamedSet.GetId()] = evt.GetNamedSetOfFiles()
		case *bes.BuildEventId_TargetCompleted:
			completed := evt.GetCompleted()
			if completed == nil {
				continue
			}
			pt := pendingTarget{label: id.TargetCompleted.GetLabel()}
			for _, og := range completed.GetOutputGroup() {
				var rootIDs []string
				for _, ref := range og.GetFileSets() {
					rootIDs = append(rootIDs, ref.GetId())
				}
				pt.groups = append(pt.groups, struct {
					name    string
					rootIDs []string
				}{name: og.GetName(), rootIDs: rootIDs})
			}
			pending = append(pending, pt)
		}
	}

	var targets []targetEntry
	for _, pt := range pending {
		entry := targetEntry{label: pt.label}
		for _, g := range pt.groups {
			files := walkFileNames(sets, g.rootIDs)
			sort.Strings(files)
			entry.outputGroups = append(entry.outputGroups, outputGroupEntry{
				name:      g.name,
				fileNames: files,
			})
		}
		sort.Slice(entry.outputGroups, func(i, j int) bool {
			return entry.outputGroups[i].name < entry.outputGroups[j].name
		})
		targets = append(targets, entry)
	}

	sort.Slice(targets, func(i, j int) bool { return targets[i].label < targets[j].label })
	return targets, sets
}

func walkFileNames(sets map[string]*bes.NamedSetOfFiles, roots []string) []string {
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
		for _, f := range set.GetFiles() {
			out = append(out, f.GetName())
		}
		for _, child := range set.GetFileSets() {
			queue = append(queue, child.GetId())
		}
	}
	return out
}
