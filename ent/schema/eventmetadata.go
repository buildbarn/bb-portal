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
		// The sequence number of the event in the invocation.
		field.Int64("sequence_number").Immutable(),

		// The time when the event was saved received.
		field.Time("event_received_at").Immutable(),

		// The hash of the event proto message.
		field.String("event_hash").Immutable(),

		// Foreign key to bazel_invocation
		field.Int("bazel_invocation_id").Immutable(),
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
		index.Edges("bazel_invocation"),
		index.Fields("sequence_number").
			Edges("bazel_invocation").
			Unique(),
	}
}

// Annotations for basel invocation.
func (EventMetadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}
