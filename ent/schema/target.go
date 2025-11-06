package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Target holds the schema definition for the Target entity.
type Target struct {
	ent.Schema
}

// Fields of the Target.
func (Target) Fields() []ent.Field {
	return []ent.Field{
		// The label of the target ex: //foo:bar.
		field.String("label"),

		// List of all tags associated with this target (for all possible
		// configurations).
		field.Strings("tag").Optional(),

		// The kind of target.
		// (e.g.,  e.g. "cc_library rule", "source file",
		// "generated file") where the completion is reported.
		field.String("target_kind").Optional(),

		// The size of the test, if the target is a test target. Unset otherwise.
		field.Enum("test_size").
			Values("UNKNOWN",
				"SMALL",
				"MEDIUM",
				"LARGE",
				"ENORMOUS").
			Optional(),

		// Overall success of the target (defaults to false).
		field.Bool("success").
			Optional().
			Default(false),

		// The timeout specified for test actions under this configured target.
		field.Int64("test_timeout").Optional(),

		// First time we saw this target.
		field.Int64("start_time_in_ms").Optional(),

		// Time we saw the event complete for this target in unix.
		field.Int64("end_time_in_ms").Optional(),

		// Duration in Milliseconds.
		// Time from target configured message received and processed until target completed message received and processed, calculated on build complete
		field.Int64("duration_in_ms").
			Optional().
			Annotations(entgql.OrderField("DURATION")),

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

// Edges of the Target.
func (Target) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("targets").
			Unique(),

		// TODO: Add these back
		// // Temporarily, also report the important outputs directly.
		// // This is only to allow existing clients help transition to the deduplicated representation;
		// // new clients should not use it.
		// edge.To("important_output", TestFile.Type).
		// 	Annotations(
		// 		entsql.OnDelete(entsql.Cascade),
		// 	),

		// // Report output artifacts (referenced transitively via output_group) which
		// // emit directories instead of singleton files. These directory_output entries
		// // will never include a uri.
		// edge.To("directory_output", TestFile.Type).
		// 	Annotations(
		// 		entsql.OnDelete(entsql.Cascade),
		// 	),

		// // The output files are arranged by their output group. If an output file
		// // is part of multiple output groups, it appears once in each output
		// // group.
		// edge.To("output_group", OutputGroup.Type).Unique().
		// 	Annotations(
		// 		entsql.OnDelete(entsql.Cascade),
		// 	),
	}
}

// Indexes of the Target.
func (Target) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("label"),
		index.Edges("bazel_invocation"),
		index.Fields("label").
			Edges("bazel_invocation").
			Unique(),
	}
}

// Annotations of the Target
func (Target) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findTargets"),
	}
}
