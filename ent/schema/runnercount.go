package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RunnerCount holds the schema definition for the RunnerCount entity.
type RunnerCount struct {
	ent.Schema
}

// Fields of the RunnerCount.
func (RunnerCount) Fields() []ent.Field {
	return []ent.Field{
		// The name of the runner.
		field.String("name").Optional(),

		// The execition kind (local, remote, etc).
		field.String("exec_kind").Optional(),

		// Count of actions of this type executed.
		field.Int64("actions_executed").Optional(),
	}
}

// Edges of the RunnerCount.
func (RunnerCount) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the action summary.
		edge.From("action_summary", ActionSummary.Type).
			Ref("runner_count").
			Unique(),
	}
}

// Indexes of the RunnerCount.
func (RunnerCount) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("action_summary"),
	}
}
