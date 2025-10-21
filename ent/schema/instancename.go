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
		edge.To("blobs", Blob.Type),
		edge.To("targets", Target.Type),
	}
}
