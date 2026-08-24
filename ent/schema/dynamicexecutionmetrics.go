package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

// DynamicExecutionMetrics holds metrics for dynamic execution races.
type DynamicExecutionMetrics struct {
	ent.Schema
}

// Fields of the DynamicExecutionMetrics struct.
func (DynamicExecutionMetrics) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the DynamicExecutionMetrics struct.
func (DynamicExecutionMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("dynamic_execution_metrics").
			Unique(),

		edge.To("race_statistics", DynamicExecutionRaceStatistic.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the DynamicExecutionMetrics struct.
func (DynamicExecutionMetrics) Indexes() []ent.Index {
	return []ent.Index{}
}

// Mixin of the DynamicExecutionMetrics struct.
func (DynamicExecutionMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
