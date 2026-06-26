package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// InstanceName holds the schema definition for the InstanceName entity.
type InstanceName struct {
	ent.Schema
}

// Fields of the InstanceName.
func (InstanceName) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().Immutable(),
	}
}

// Edges of InstanceName.
func (InstanceName) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("bazel_invocations", BazelInvocation.Type),
		edge.To("builds", Build.Type),
		edge.To("targets", Target.Type),
		edge.To("file_paths", FilePath.Type),
	}
}

// Mixin of the InstanceName.
func (InstanceName) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
