package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TestSummary holds the schema definition for the TestSummary entity.
type TestSummary struct {
	ent.Schema
}

// Fields of the TestSummary.
func (TestSummary) Fields() []ent.Field {
	return []ent.Field{
		// Wrapper around BlazeTestStatus to support importing that enum to proto3.
		// Overall status of test, accumulated over all runs, shards, and attempts.
		field.Enum("overall_status").Optional().
			Values("NO_STATUS",
				"PASSED",
				"FLAKY",
				"TIMEOUT",
				"FAILED",
				"INCOMPLETE",
				"REMOTE_FAILURE",
				"FAILED_TO_BUILD",
				"TOOL_HALTED_BEFORE_TESTING").
			Default("NO_STATUS"),

		// Total number of shard attempts.
		// E.g., if a target has 4 runs, 3 shards, each with 2 attempts,
		// then total_run_count will be 4*3*2 = 24.
		field.Int32("total_run_count").Optional(),

		// Value of runs_per_test for the test.
		field.Int32("run_count").Optional(),

		// Number of attempts.
		// If there are a different number of attempts per shard, the highest attempt
		// count across all shards for each run is used.
		field.Int32("attempt_count").Optional(),

		// Number of shards.
		field.Int32("shard_count").Optional(),

		// Total number of cached test actions
		field.Int32("total_num_cached").Optional(),

		// When the test first started running.
		field.Int64("first_start_time").Optional(),

		// When the test last finished running.
		field.Int64("last_stop_time").Optional(),

		// The total runtime of the test.
		field.Int64("total_run_duration").Optional(),

		// Test target label, possibly redundant and could be removed.
		field.String("label").Optional(),
	}
}

// Edges of the TestSummary.
func (TestSummary) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back tot he test collection.
		edge.From("test_collection", TestCollection.Type).
			Ref("test_summary"),

		// Path to logs of passed runs.
		edge.To("passed", TestFile.Type),

		// Path to logs of failed runs;
		edge.To("failed", TestFile.Type),
	}
}
