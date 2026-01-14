package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EventMetadata holds the schema definition for the EventMetadata entity.
type EventMetadata struct {
	ent.Schema
}

// Fields of the EventMetadata.
func (EventMetadata) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("handled").
			Comment("Binary representation of the events that have been handled"),
		field.Time("event_received_at").
			Comment("Last time an event was received"),
		field.Int64("version").
			Comment("Optimistic lock version number"),
		field.Int64("bazel_invocation_id").
			Comment("The id of the bazel invocation").
			Immutable().
			Unique(),
	}
}

// Edges of the EventMetadata.
func (EventMetadata) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("event_metadata").
			Unique().
			Required().
			Immutable(),
	}
}

// Indexes of the EventMetadata.
func (EventMetadata) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_received_at"),
	}
}

// Annotations for basel invocation.
func (EventMetadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

// Mixin of the EventMetadata.
func (EventMetadata) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
