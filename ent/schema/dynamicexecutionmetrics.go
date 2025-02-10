package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

// DynamicExecutionMetrics holds the schema definition for the DynamicExecutionMetrics entity.
type DynamicExecutionMetrics struct {
	ent.Schema
}

// Fields of the DynamicExecutionMetrics.
func (DynamicExecutionMetrics) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the DynamicExecutionMetrics.
func (DynamicExecutionMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object.
		edge.From("metrics", Metrics.Type).
			Ref("dynamic_execution_metrics").
			Unique(),

		// Race statistics grouped by mnemonic, local_name, remote_name.
		edge.To("race_statistics", RaceStatistics.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
