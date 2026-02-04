package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("overall_status").Optional(),

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
		field.Time("first_start_time").Optional().
			Annotations(entgql.OrderField("FIRST_START_TIME")),

		// When the test last finished running.
		field.Time("last_stop_time").Optional(),

		// The total runtime of the test.
		field.Int64("total_run_duration_in_ms").Optional().Nillable().
			Annotations(entgql.OrderField("TOTAL_RUN_DURATION_IN_MS")),
	}
}

// Edges of the TestSummary.
func (TestSummary) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("invocation_target", InvocationTarget.Type).
			Ref("test_summary").
			Required().
			Unique(),

		// A collection of test results associated with this collection
		edge.To("test_results", TestResult.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the TestSummary.
func (TestSummary) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("invocation_target"),
	}
}

// Annotations for TestSummary.
func (TestSummary) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findTestSummaries"),
	}
}

// Mixin of the TestSummary.
func (TestSummary) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
