package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TargetComplete holds the schema definition for the TargetComplete entity.
type TargetComplete struct {
	ent.Schema
}

// Fields of the TargetComplete.
func (TargetComplete) Fields() []ent.Field {
	return []ent.Field{
		// Did the target build successfully.
		field.Bool("success").Optional(),

		// List of tags associated with this configured target.
		field.Strings("tag").Optional(),

		// Target Kind.
		// The kind of target (e.g.,  e.g. "cc_library rule", "source file",
		// "generated file") where the completion is reported.
		// Deprecated: use the target_kind field in TargetConfigured instead.
		field.String("target_kind").Optional(),

		// Time we saw the event complete for this target in unix.
		field.Int64("end_time_in_ms").Optional(),

		// The timeout specified for test actions under this configured target.
		// Deprecated, use `test_timeout` instead.
		field.Int64("test_timeout_seconds").Optional(),

		// The timeout specified for test actions under this configured target.
		field.Int64("test_timeout").Optional(),

		// The size of the test, if the target is a test target. Unset otherwise.
		// Deprecated: use the test_size field in TargetConfigured instead.
		field.Enum("test_size").
			Values("UNKNOWN",
				"SMALL",
				"MEDIUM",
				"LARGE",
				"ENORMOUS").
			Optional(),

		// TODO: implement failure detail.
	}
}

// Edges of the TargetComplete.
func (TargetComplete) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the target pair.
		edge.From("target_pair", TargetPair.Type).
			Ref("completion").
			Unique(),

		// Temporarily, also report the important outputs directly.
		// This is only to allow existing clients help transition to the deduplicated representation;
		// new clients should not use it.
		edge.To("important_output", TestFile.Type),

		// Report output artifacts (referenced transitively via output_group) which
		// emit directories instead of singleton files. These directory_output entries
		// will never include a uri.
		edge.To("directory_output", TestFile.Type),

		// The output files are arranged by their output group. If an output file
		// is part of multiple output groups, it appears once in each output
		// group.
		edge.To("output_group", OutputGroup.Type).Unique(),
	}
}
