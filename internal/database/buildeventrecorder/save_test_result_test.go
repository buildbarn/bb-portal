package buildeventrecorder_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/digest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestSaveTestResultBatchLinksOutputsToOwnTestResult checks that when a
// single batch of TestResult events spans several targets, each test
// result is linked to its own test_action_output files. The test
// results are inserted in bulk and the outputs are attached by position
// in the returned ID list, so this fails if the insert does not return
// IDs in batch order (e.g. when a join emits rows in targets order).
func TestSaveTestResultBatchLinksOutputsToOwnTestResult(t *testing.T) {
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())
	db, conn := testutils.SetupTestDBWithConn(t, dbProvider)

	recorder, err := buildeventrecorder.NewBuildEventRecorder(
		ctx,
		db,
		auth.NewStaticAuthorizer(func(digest.InstanceName) bool { return true }),
		&bb_portal.BuildEventStreamService_SaveDataLevel{
			Level: &bb_portal.BuildEventStreamService_SaveDataLevel_BasicAndTarget{
				BasicAndTarget: &emptypb.Empty{},
			},
		},
		noop.NewTracerProvider(),
		"",
		uuid.New().String(),
		&buildeventrecorder.DataExtractors{},
		"",
		10*time.Second,
	)
	require.NoError(t, err)

	const configID = "cfg"
	configuration := &bes.BuildEventId_ConfigurationId{Id: configID}
	// Targets are created in this order, which differs from the order
	// of the test results below.
	labels := []string{"//pkg:a_test", "//pkg:b_test", "//pkg:c_test"}

	var seq uint32
	event := func(be *bes.BuildEvent) buildeventrecorder.BuildEventWithInfo {
		seq++
		return buildeventrecorder.BuildEventWithInfo{Event: be, SequenceNumber: seq}
	}

	setup := []buildeventrecorder.BuildEventWithInfo{
		event(&bes.BuildEvent{
			Id: &bes.BuildEventId{Id: &bes.BuildEventId_Configuration{Configuration: configuration}},
			Payload: &bes.BuildEvent_Configuration{Configuration: &bes.Configuration{
				Mnemonic: "k8-fastbuild",
			}},
		}),
	}
	for _, label := range labels {
		setup = append(setup, event(&bes.BuildEvent{
			Id: &bes.BuildEventId{Id: &bes.BuildEventId_TargetConfigured{
				TargetConfigured: &bes.BuildEventId_TargetConfiguredId{Label: label},
			}},
			Payload: &bes.BuildEvent_Configured{Configured: &bes.TargetConfigured{
				TargetKind: "go_test rule",
			}},
		}))
	}
	for _, label := range labels {
		setup = append(setup, event(&bes.BuildEvent{
			Id: &bes.BuildEventId{Id: &bes.BuildEventId_TargetCompleted{
				TargetCompleted: &bes.BuildEventId_TargetCompletedId{Label: label, Configuration: configuration},
			}},
			Children: []*bes.BuildEventId{{Id: &bes.BuildEventId_TestSummary{
				TestSummary: &bes.BuildEventId_TestSummaryId{Label: label, Configuration: configuration},
			}}},
			Payload: &bes.BuildEvent_Completed{Completed: &bes.TargetComplete{Success: true}},
		}))
	}
	require.NoError(t, recorder.SaveBatch(ctx, setup))

	// Production databases have planner statistics, under which a small
	// targets table makes Postgres hash join the inputs against it.
	_, err = conn.ExecContext(ctx, "ANALYZE")
	require.NoError(t, err)

	// One batch of test results, deliberately not in targets order,
	// including two shards of the same target.
	type resultKey struct {
		label string
		shard int32
	}
	results := []resultKey{
		{"//pkg:c_test", 0},
		{"//pkg:a_test", 0},
		{"//pkg:b_test", 0},
		{"//pkg:a_test", 1},
	}
	expectedHash := make(map[resultKey]string, len(results))
	var testResults []buildeventrecorder.BuildEventWithInfo
	for _, r := range results {
		sum := sha256.Sum256([]byte(fmt.Sprintf("test.log of %s shard %d", r.label, r.shard)))
		hash := hex.EncodeToString(sum[:])
		expectedHash[r] = hash
		testResults = append(testResults, event(&bes.BuildEvent{
			Id: &bes.BuildEventId{Id: &bes.BuildEventId_TestResult{
				TestResult: &bes.BuildEventId_TestResultId{
					Label:         r.label,
					Configuration: configuration,
					Run:           1,
					Shard:         r.shard,
					Attempt:       1,
				},
			}},
			Payload: &bes.BuildEvent_TestResult{TestResult: &bes.TestResult{
				Status:        bes.TestStatus_PASSED,
				ExecutionInfo: &bes.TestResult_ExecutionInfo{Strategy: "remote"},
				TestActionOutput: []*bes.File{{
					Name: "test.log",
					File: &bes.File_Uri{Uri: fmt.Sprintf("bytestream://localhost:8980/blobs/%s/%d", hash, 100+r.shard)},
				}},
			}},
		}))
	}
	require.NoError(t, recorder.SaveBatch(ctx, testResults))

	rows, err := db.Ent().TestResult.Query().
		WithTestSummary(func(q *ent.TestSummaryQuery) {
			q.WithInvocationTarget(func(q *ent.InvocationTargetQuery) { q.WithTarget() })
		}).
		WithTestActionOutput(func(q *ent.FileQuery) { q.WithDigest() }).
		All(ctx)
	require.NoError(t, err)
	require.Len(t, rows, len(results))

	got := make(map[resultKey]string, len(rows))
	for _, row := range rows {
		key := resultKey{
			label: row.Edges.TestSummary.Edges.InvocationTarget.Edges.Target.Label,
			shard: row.Shard,
		}
		require.Len(t, row.Edges.TestActionOutput, 1, "test result %v", key)
		got[key] = hex.EncodeToString(row.Edges.TestActionOutput[0].Edges.Digest.Hash)
	}
	require.Equal(t, expectedHash, got)
}
