package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DynamicExecutionRaceStatistic holds the outcome of a dynamic execution race.
type DynamicExecutionRaceStatistic struct {
	ent.Schema
}

// Fields of the DynamicExecutionRaceStatistic struct.
func (DynamicExecutionRaceStatistic) Fields() []ent.Field {
	return []ent.Field{
		field.String("mnemonic").Optional(),
		field.String("local_runner").Optional(),
		field.String("remote_runner").Optional(),
		field.Int32("local_wins").Optional(),
		field.Int32("remote_wins").Optional(),
	}
}

// Edges of the DynamicExecutionRaceStatistic struct.
func (DynamicExecutionRaceStatistic) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dynamic_execution_metrics", DynamicExecutionMetrics.Type).
			Ref("race_statistics").
			Unique(),
	}
}

// Indexes of the DynamicExecutionRaceStatistic struct.
func (DynamicExecutionRaceStatistic) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("dynamic_execution_metrics"),
	}
}

// Mixin of the DynamicExecutionRaceStatistic struct.
func (DynamicExecutionRaceStatistic) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
