package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ConnectionMetadata holds the schema definition for the ConnectionMetadata entity.
type ConnectionMetadata struct {
	ent.Schema
}

// Fields of the ConnectionMetadata.
func (ConnectionMetadata) Fields() []ent.Field {
	return []ent.Field{
		// The time when the event was saved received.
		field.Time("connection_last_open_at"),
	}
}

// Edges of the ConnectionMetadata.
func (ConnectionMetadata) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("connection_metadata").
			Unique().
			Required(),
	}
}

// Indexes of the ConnectionMetadata.
func (ConnectionMetadata) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation").
			Unique(),
	}
}

// Annotations for basel invocation.
func (ConnectionMetadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}
