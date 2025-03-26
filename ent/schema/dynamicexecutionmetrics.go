package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DynamicExecutionMetrics holds the schema definition for the DynamicExecutionMetrics entity.
type DynamicExecutionMetrics struct {
	ent.Schema
}

// Fields of the DynamicExecutionMetrics.
func (DynamicExecutionMetrics) Fields() []ent.Field {
	return []ent.Field{
		field.Int("metrics_id").Optional(),
	}
}

// Edges of the DynamicExecutionMetrics.
func (DynamicExecutionMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object.
		edge.From("metrics", Metrics.Type).
			Ref("dynamic_execution_metrics").
			Unique().
			Field("metrics_id"),

		// Race statistics grouped by mnemonic, local_name, remote_name.
		edge.To("race_statistics", RaceStatistics.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
