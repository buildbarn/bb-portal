package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TargetKindMapping holds the schema definition for the TargetKindMapping entity.
type TargetKindMapping struct {
	ent.Schema
}

// Fields of the TargetKindMapping.
func (TargetKindMapping) Fields() []ent.Field {
	return []ent.Field{
		// First time we saw this InvocationTarget.
		field.Int64("start_time_in_ms").Optional(),
	}
}

// Edges of the TargetKindMapping.
func (TargetKindMapping) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("target_kind_mappings").
			Required().
			Unique(),

		// Edge back to the target
		edge.From("target", Target.Type).
			Ref("target_kind_mappings").
			Required().
			Unique(),
	}
}

// Indexes of the TargetKindMapping.
func (TargetKindMapping) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
		index.Edges("target"),
		index.Edges("bazel_invocation", "target").
			Unique(),
	}
}

// Annotations of the TargetKindMapping.
func (TargetKindMapping) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}
