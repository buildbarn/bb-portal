package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TargetConfigured holds the schema definition for the TargetConfigured entity.
type TargetConfigured struct {
	ent.Schema
}

// Fields of the TargetConfigured.
func (TargetConfigured) Fields() []ent.Field {
	return []ent.Field{
		// Tags.
		// List of all tags associated with this target (for all possible
		// configurations).
		field.Strings("tag").Optional(),

		// The kind of target.
		// (e.g.,  e.g. "cc_library rule", "source file",
		// "generated file") where the completion is reported.
		field.String("target_kind").Optional(),

		// First time we saw this target.
		field.Int64("start_time_in_ms").Optional(),

		// The size of the test, if the target is a test target. Unset otherwise.
		// TODO is there somewway to reference this instead of repeating it?
		field.Enum("test_size").
			Values("UNKNOWN",
				"SMALL",
				"MEDIUM",
				"LARGE",
				"ENORMOUS").
			Optional(),

		field.Int("target_pair_id").Optional(),
	}
}

// Edges of the TargetConfigured.
func (TargetConfigured) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the target pair.
		edge.From("target_pair", TargetPair.Type).
			Ref("configuration").
			Unique().
			Field("target_pair_id"),
	}
}
