package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EventFile holds the schema definition for the EventFile entity.
type EventFile struct {
	ent.Schema
}

// Fields of the EventFile.
func (EventFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").Unique().Immutable(),
		field.Time("mod_time"),
		field.String("protocol"), // *.bep, *.log, etc
		field.String("mime_type"),
		field.String("status").Default("DETECTED"),
		field.String("reason").Optional(),
	}
}

// Edges of the EventFile.
func (EventFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("bazel_invocation", BazelInvocation.Type).
			Unique(),
	}
}

// Indexes of the EventFile.
func (EventFile) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
