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

// InvocationTarget holds the schema definition for the InvocationTarget entity.
type InvocationTarget struct {
	ent.Schema
}

// Fields of the InvocationTarget.
func (InvocationTarget) Fields() []ent.Field {
	return []ent.Field{
		// Overall success of the target.
		field.Bool("success").Default(false),

		// List of all tags associated with this target (for all possible
		// configurations).
		field.Strings("tags").Optional(),

		// First time we saw this target.
		field.Int64("start_time_in_ms").
			Optional().
			Annotations(entgql.OrderField("STARTED_AT")),

		// Time we saw the event complete for this target in unix.
		field.Int64("end_time_in_ms").Optional(),

		// Duration in Milliseconds.
		// Time from target configured message received and processed until target completed message received and processed, calculated on build complete
		field.Int64("duration_in_ms").
			Optional().
			Annotations(entgql.OrderField("DURATION")),

		field.String("failure_message").Optional(),

		// reason the target was aborted if any
		field.Enum("abort_reason").
			Values(
				"ANALYSIS_FAILURE",
				"INCOMPLETE",
				"INTERNAL",
				"LOADING_FAILURE",
				"NO_ANALYZE",
				"NO_BUILD",
				"NONE", // Added NONE to represent no abort reason
				"OUT_OF_MEMORY",
				"REMOTE_ENVIRONMENT_FAILURE",
				"SKIPPED",
				"TIME_OUT",
				"UNKNOWN",
				"USER_INTERRUPTED",
			),
	}
}

// Edges of the InvocationTarget.
func (InvocationTarget) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("invocation_targets").
			Required().
			Unique(),

		// Edge back to the target
		edge.From("target", Target.Type).
			Ref("invocation_targets").
			Required().
			Unique(),

		// Edge to the configuration used for this target
		edge.To("configuration", Configuration.Type).
			Unique(),

		edge.To("test_summary", TestSummary.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the InvocationTarget.
func (InvocationTarget) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
		index.Edges("target"),
		index.Edges("configuration"),
		index.Edges("bazel_invocation", "target", "configuration"),
	}
}

// Annotations of the InvocationTarget
func (InvocationTarget) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}

// Mixin of the InvocationTarget.
func (InvocationTarget) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
