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

// TestCollection holds the schema definition for the TestCollection entity.
type TestCollection struct {
	ent.Schema
}

// Fields of the TestCollection.
func (TestCollection) Fields() []ent.Field {
	return []ent.Field{
		// The label associated with this test.
		field.String("label").
			Optional(),

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
		field.String("strategy").
			Optional(),

		// If the test was cached locally.
		field.Bool("cached_locally").
			Optional(),

		// If the test was cached remotely.
		field.Bool("cached_remotely").
			Optional(),

		field.Time("first_seen").
			Optional().
			Nillable().
			Annotations(entgql.OrderField("FIRST_SEEN")),

		// The test duration in milliseconds.
		field.Int64("duration_ms").
			Optional().
			Annotations(entgql.OrderField("DURATION")),
	}
}

// Edges of the TestCollection.
func (TestCollection) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocaiton.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("test_collection").Unique(),

		// The test summary aossicated with the test.
		edge.To("test_summary", TestSummary.Type).Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// A collection of test results associated with this collection
		edge.To("test_results", TestResultBES.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the TestCollection.
func (TestCollection) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("label"),
		index.Edges("bazel_invocation"),
		// Make each label unique per invocation.
		index.Fields("label").
			Edges("bazel_invocation").
			Unique(),
	}
}

// Annotations of the Test Collection
func (TestCollection) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findTests"),
	}
}
