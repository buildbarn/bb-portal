package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RaceStatistics holds the schema definition for the RaceStatistics entity.
type RaceStatistics struct {
	ent.Schema
}

// Fields of the RaceStatistics.
func (RaceStatistics) Fields() []ent.Field {
	return []ent.Field{
		// NOTE: Not currently included on the proto, but included here now for completeness.
		// Mnemonic of the action.
		field.String("mnemonic").Optional(),

		// Name of runner of local branch.
		field.String("local_runner").Optional(),

		// Name of runner of remote branch.
		field.String("remote_runner").Optional(),

		// Number of wins of local branch in race.
		field.Int64("local_wins").Optional(),

		// Number of wins of remote branch in race.
		field.Int64("renote_wins").Optional(),
	}
}

// Edges of the RaceStatistics.
func (RaceStatistics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the dynamic execution metrics object.
		edge.From("dynamic_execution_metrics", DynamicExecutionMetrics.Type).
			Ref("race_statistics").
			Unique(),
	}
}
