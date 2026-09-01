package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CompletedAction holds the schema definition for the Completed Action Logger entity.
type CompletedAction struct {
	ent.Schema
}

// Fields of the CompletedAction.
func (CompletedAction) Fields() []ent.Field {
	return []ent.Field{
		// invocation id for each action
		field.String("invocation_id"),

		// target id for each action
		field.String("target_id"),

		// mnemonic for each action
		field.String("mnemonic"),

		// tool_name for each action
		field.String("tool_name"),

		// action_id for each action
		field.String("action_id"),

		// Time when the completed action was first inserted into the database.
		field.Time("created_timestamp").Default(time.Now).Immutable().Annotations(entgql.Skip()),
	}
}

// Edges of the CompletedAction.
func (CompletedAction) Edges() []ent.Edge {
	return []ent.Edge{
		// Action Digest
		edge.To("action_digest", Digest.Type).
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),

		// Historical Execute Response Digest
		edge.To("historical_execute_response_digest", Digest.Type).
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}

// Indexes of the CompletedAction.
func (CompletedAction) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("action_digest"),
		index.Edges("historical_execute_response_digest").Unique(),
	}
}

// Annotations of the CompletedAction.
func (CompletedAction) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findCompletedActions"),
	}
}

// Mixin of the CompletedAction.
func (CompletedAction) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
