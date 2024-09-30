package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TestCollection holds the schema definition for the TestCollection entity.
type TestCollection struct {
	ent.Schema
}

// Fields of the TestCollection.
func (TestCollection) Fields() []ent.Field {
	return []ent.Field{
		// The label associated with this test.
		field.String("label").Optional(),

		// The overall status of the test.
		field.Enum("overall_status").
			Optional().
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

		// The strategy of the test.
		field.String("strategy").Optional(),

		// If the test was cached locally.
		field.Bool("cached_locally").Optional(),

		// If the test was cached remotely.
		field.Bool("cached_remotely").Optional(),

		// The test duration in milliseconds.
		field.Int64("duration_ms").Optional(),
	}
}

// Edges of the TestCollection.
func (TestCollection) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocaiton.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("test_collection"),

		// The test summary aossicated with the test.
		edge.To("test_summary", TestSummary.Type).Unique(),

		// A collection of test results associated.
		edge.To("test_results", TestResultBES.Type),
	}
}
