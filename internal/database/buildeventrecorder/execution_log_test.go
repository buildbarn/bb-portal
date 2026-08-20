package buildeventrecorder

import (
	"bytes"
	"testing"

	bazelprotobuf "github.com/bazelbuild/bazel/src/main/protobuf"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	storagedigest "github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protodelim"
)

func compactExecutionLog(t *testing.T, entries ...*bazelprotobuf.ExecLogEntry) []byte {
	t.Helper()
	var data bytes.Buffer
	writer, err := zstd.NewWriter(&data)
	require.NoError(t, err)
	for _, entry := range entries {
		_, err := protodelim.MarshalTo(writer, entry)
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	return data.Bytes()
}

func TestParseCompactExecutionLog(t *testing.T) {
	const actionHash = "4048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07"
	data := compactExecutionLog(
		t,
		&bazelprotobuf.ExecLogEntry{
			Type: &bazelprotobuf.ExecLogEntry_Invocation_{Invocation: &bazelprotobuf.ExecLogEntry_Invocation{HashFunctionName: "SHA-256"}},
		},
		&bazelprotobuf.ExecLogEntry{
			Id:   7,
			Type: &bazelprotobuf.ExecLogEntry_File_{File: &bazelprotobuf.ExecLogEntry_File{Path: "bazel-out/bin/example"}},
		},
		&bazelprotobuf.ExecLogEntry{
			Type: &bazelprotobuf.ExecLogEntry_Spawn_{Spawn: &bazelprotobuf.ExecLogEntry_Spawn{
				TargetLabel: "//example:example",
				Mnemonic:    "GoLink",
				Runner:      "remote",
				CacheHit:    true,
				Outputs: []*bazelprotobuf.ExecLogEntry_Output{
					{Type: &bazelprotobuf.ExecLogEntry_Output_OutputId{OutputId: 7}},
				},
				Digest: &bazelprotobuf.Digest{Hash: actionHash, SizeBytes: 145},
			}},
		},
	)

	actions, err := parseCompactExecutionLog(bytes.NewReader(data), storagedigest.EmptyInstanceName)
	require.NoError(t, err)
	require.Len(t, actions, 1)
	require.Equal(t, "//example:example", actions[0].targetLabel)
	require.Equal(t, "GoLink", actions[0].mnemonic)
	require.Equal(t, []string{"bazel-out/bin/example"}, actions[0].outputPaths)
	require.Equal(t, "remote", actions[0].runner)
	require.NotNil(t, actions[0].cacheHit)
	require.True(t, *actions[0].cacheHit)
	require.Equal(t, actionHash, actions[0].actionDigest.GetHashString())
	require.EqualValues(t, 145, actions[0].actionDigest.GetSizeBytes())
	require.Equal(t, "", actions[0].actionDigest.GetInstanceName().String())
}

func TestParseCompactExecutionLogIgnoresSpawnsWithoutOutputs(t *testing.T) {
	data := compactExecutionLog(
		t,
		&bazelprotobuf.ExecLogEntry{Type: &bazelprotobuf.ExecLogEntry_Invocation_{Invocation: &bazelprotobuf.ExecLogEntry_Invocation{HashFunctionName: "SHA-256"}}},
		&bazelprotobuf.ExecLogEntry{Type: &bazelprotobuf.ExecLogEntry_Spawn_{Spawn: &bazelprotobuf.ExecLogEntry_Spawn{Mnemonic: "Internal"}}},
	)

	actions, err := parseCompactExecutionLog(bytes.NewReader(data), storagedigest.EmptyInstanceName)
	require.NoError(t, err)
	require.Empty(t, actions)
}

func TestParseCompactExecutionLogClassifiesExecutionLocation(t *testing.T) {
	const actionHash = "4048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07"
	data := compactExecutionLog(
		t,
		&bazelprotobuf.ExecLogEntry{Type: &bazelprotobuf.ExecLogEntry_Invocation_{Invocation: &bazelprotobuf.ExecLogEntry_Invocation{HashFunctionName: "SHA-256"}}},
		&bazelprotobuf.ExecLogEntry{Id: 1, Type: &bazelprotobuf.ExecLogEntry_File_{File: &bazelprotobuf.ExecLogEntry_File{Path: "bazel-out/bin/local"}}},
		&bazelprotobuf.ExecLogEntry{Id: 2, Type: &bazelprotobuf.ExecLogEntry_File_{File: &bazelprotobuf.ExecLogEntry_File{Path: "bazel-out/bin/cached"}}},
		&bazelprotobuf.ExecLogEntry{Type: &bazelprotobuf.ExecLogEntry_Spawn_{Spawn: &bazelprotobuf.ExecLogEntry_Spawn{
			TargetLabel: "//example:local",
			Mnemonic:    "GoCompilePkg",
			Runner:      "darwin-sandbox",
			Outputs: []*bazelprotobuf.ExecLogEntry_Output{
				{Type: &bazelprotobuf.ExecLogEntry_Output_OutputId{OutputId: 1}},
			},
			Digest: &bazelprotobuf.Digest{Hash: actionHash, SizeBytes: 145},
		}}},
		&bazelprotobuf.ExecLogEntry{Type: &bazelprotobuf.ExecLogEntry_Spawn_{Spawn: &bazelprotobuf.ExecLogEntry_Spawn{
			TargetLabel: "//example:cached",
			Mnemonic:    "GoLink",
			Runner:      "remote cache hit",
			CacheHit:    true,
			Outputs: []*bazelprotobuf.ExecLogEntry_Output{
				{Type: &bazelprotobuf.ExecLogEntry_Output_OutputId{OutputId: 2}},
			},
			Digest: &bazelprotobuf.Digest{Hash: actionHash, SizeBytes: 145},
		}}},
		&bazelprotobuf.ExecLogEntry{Type: &bazelprotobuf.ExecLogEntry_SymlinkAction_{SymlinkAction: &bazelprotobuf.ExecLogEntry_SymlinkAction{
			TargetLabel: "//example:symlink",
			Mnemonic:    "Symlink",
			OutputPath:  "bazel-out/bin/symlink",
		}}},
	)

	actions, err := parseCompactExecutionLog(bytes.NewReader(data), storagedigest.EmptyInstanceName)
	require.NoError(t, err)
	require.Len(t, actions, 3)

	require.Equal(t, "darwin-sandbox", actions[0].runner)
	require.NotNil(t, actions[0].cacheHit)
	require.False(t, *actions[0].cacheHit)
	require.Equal(t, actionHash, actions[0].actionDigest.GetHashString())
	require.Equal(t, "remote cache hit", actions[1].runner)
	require.NotNil(t, actions[1].cacheHit)
	require.True(t, *actions[1].cacheHit)
	require.Equal(t, actionHash, actions[1].actionDigest.GetHashString())
	require.Equal(t, internalActionRunner, actions[2].runner)
	require.Nil(t, actions[2].cacheHit)
	require.Equal(t, []string{"bazel-out/bin/symlink"}, actions[2].outputPaths)
	require.Equal(t, storagedigest.BadDigest, actions[2].actionDigest)
}

func TestClearMissingActionDigests(t *testing.T) {
	const (
		availableHash = "4048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07"
		missingHash   = "5048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07"
	)
	availableDigest := testActionDigest(t, availableHash)
	missingDigest := testActionDigest(t, missingHash)
	actions := []executionLogAction{
		{runner: "darwin-sandbox", actionDigest: availableDigest},
		{runner: "remote", actionDigest: missingDigest},
		{runner: internalActionRunner},
	}

	clearMissingActionDigests(
		actions,
		storagedigest.NewSetBuilder(1).Add(missingDigest).Build(),
	)

	require.Equal(t, availableDigest, actions[0].actionDigest)
	require.Equal(t, storagedigest.BadDigest, actions[1].actionDigest)
	require.Equal(t, "remote", actions[1].runner)
	require.Equal(t, storagedigest.BadDigest, actions[2].actionDigest)
}

func testActionDigest(t *testing.T, hash string) storagedigest.Digest {
	t.Helper()
	digestFunction, err := storagedigest.EmptyInstanceName.GetDigestFunction(0, len(hash))
	require.NoError(t, err)
	digest, err := digestFunction.NewDigest(hash, 145)
	require.NoError(t, err)
	return digest
}

func boolPointer(value bool) *bool {
	return &value
}

func TestMatchExecutionLogActions(t *testing.T) {
	const (
		firstHash  = "4048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07"
		secondHash = "5048aad102bbf0ee98cdfc2dc9797d9b01afab79ee77738134828ae637f5be07"
	)

	t.Run("matches exact and canonical-label fallback", func(t *testing.T) {
		databaseActions := []*ent.Action{
			{ID: 1, Label: "//example:first", Type: "GoLink", PrimaryOutput: "bazel-out/bin/first"},
			{ID: 2, Label: "//example:second", Type: "GoLink", PrimaryOutput: "bazel-out/bin/second"},
		}
		firstDigest := testActionDigest(t, firstHash)
		secondDigest := testActionDigest(t, secondHash)
		matches := matchExecutionLogActions(databaseActions, []executionLogAction{
			{targetLabel: "//example:first", mnemonic: "GoLink", outputPaths: []string{"bazel-out/bin/first"}, runner: "remote", cacheHit: boolPointer(false), actionDigest: firstDigest},
			{targetLabel: "@@//example:second", mnemonic: "GoLink", outputPaths: []string{"bazel-out/bin/second"}, runner: "remote cache hit", cacheHit: boolPointer(true), actionDigest: secondDigest},
		})

		require.Equal(t, "remote", matches[1].runner)
		require.NotNil(t, matches[1].cacheHit)
		require.False(t, *matches[1].cacheHit)
		require.Equal(t, firstDigest, matches[1].actionDigest)
		require.Equal(t, "remote cache hit", matches[2].runner)
		require.NotNil(t, matches[2].cacheHit)
		require.True(t, *matches[2].cacheHit)
		require.Equal(t, secondDigest, matches[2].actionDigest)
	})

	t.Run("retains local runner and available digest", func(t *testing.T) {
		databaseActions := []*ent.Action{{ID: 1, Label: "//example:local", Type: "GoLink", PrimaryOutput: "bazel-out/bin/local"}}
		digest := testActionDigest(t, firstHash)
		matches := matchExecutionLogActions(databaseActions, []executionLogAction{{
			targetLabel:  "//example:local",
			mnemonic:     "GoLink",
			outputPaths:  []string{"bazel-out/bin/local"},
			runner:       "darwin-sandbox",
			cacheHit:     boolPointer(false),
			actionDigest: digest,
		}})

		require.Equal(t, "darwin-sandbox", matches[1].runner)
		require.NotNil(t, matches[1].cacheHit)
		require.False(t, *matches[1].cacheHit)
		require.Equal(t, digest, matches[1].actionDigest)
	})

	t.Run("does not guess when output matching is ambiguous", func(t *testing.T) {
		databaseActions := []*ent.Action{
			{ID: 1, Label: "//example:first", Type: "GoLink", PrimaryOutput: "bazel-out/bin/shared"},
			{ID: 2, Label: "//example:second", Type: "GoLink", PrimaryOutput: "bazel-out/bin/shared"},
		}
		matches := matchExecutionLogActions(databaseActions, []executionLogAction{{
			targetLabel:  "@@//example:canonical",
			mnemonic:     "GoLink",
			outputPaths:  []string{"bazel-out/bin/shared"},
			runner:       "remote",
			actionDigest: testActionDigest(t, firstHash),
		}})

		require.Empty(t, matches)
	})

	t.Run("drops conflicting digests for one action", func(t *testing.T) {
		databaseActions := []*ent.Action{{ID: 1, Label: "//example:first", Type: "GoLink", PrimaryOutput: "bazel-out/bin/first"}}
		matches := matchExecutionLogActions(databaseActions, []executionLogAction{
			{targetLabel: "//example:first", mnemonic: "GoLink", outputPaths: []string{"bazel-out/bin/first"}, runner: "remote", actionDigest: testActionDigest(t, firstHash)},
			{targetLabel: "//example:first", mnemonic: "GoLink", outputPaths: []string{"bazel-out/bin/first"}, runner: "remote", actionDigest: testActionDigest(t, secondHash)},
		})

		require.Empty(t, matches)
	})

	t.Run("drops conflicting cache results for one action", func(t *testing.T) {
		databaseActions := []*ent.Action{{ID: 1, Label: "//example:first", Type: "GoLink", PrimaryOutput: "bazel-out/bin/first"}}
		actionDigest := testActionDigest(t, firstHash)
		matches := matchExecutionLogActions(databaseActions, []executionLogAction{
			{targetLabel: "//example:first", mnemonic: "GoLink", outputPaths: []string{"bazel-out/bin/first"}, runner: "remote", cacheHit: boolPointer(false), actionDigest: actionDigest},
			{targetLabel: "//example:first", mnemonic: "GoLink", outputPaths: []string{"bazel-out/bin/first"}, runner: "remote", cacheHit: boolPointer(true), actionDigest: actionDigest},
		})

		require.Empty(t, matches)
	})
}
