package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TargetPair holds the schema definition for the TargetPair entity.
type TargetPair struct {
	ent.Schema
}

// Fields of the TargetPair.
func (TargetPair) Fields() []ent.Field {
	return []ent.Field{
		// The label of the target ex: //foo:bar.
		field.String("label").Optional(),

		// Duration in Milliseconds.
		// Time from target configured message received and processed until target completed message received and processed, calculated on build complete
		field.Int64("duration_in_ms").Optional(),

		// Overall success of the target (defaults to false).
		field.Bool("success").Optional().Default(false),

		// The target kind if available.
		field.String("target_kind").Optional(),

		// The size of the test, if the target is a test target. Unset otherwise.
		field.Enum("test_size").
			Values("UNKNOWN",
				"SMALL",
				"MEDIUM",
				"LARGE",
				"ENORMOUS").
			Default("UNKNOWN").
			Optional(),

		// reason the target was aborted if any
		field.Enum("abort_reason").
			Values("UNKNOWN",
				"USER_INTERRUPTED",
				"NO_ANALYZE",
				"NO_BUILD",
				"TIME_OUT",
				"REMOTE_ENVIRONMENT_FAILURE",
				"INTERNAL",
				"LOADING_FAILURE",
				"ANALYSIS_FAILURE",
				"SKIPPED",
				"INCOMPLETE",
				"OUT_OF_MEMORY").
			Optional(),
	}
}

// Edges of the TargetPair.
func (TargetPair) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("targets"),

		// Edge to the target configuration object.
		edge.To("configuration", TargetConfigured.Type).Unique(),

		// Edge to the target completed object.
		edge.To("completion", TargetComplete.Type).Unique(),
	}
}
