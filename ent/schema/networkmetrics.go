package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NetworkMetrics holds the schema definition for the NetworkMetrics entity.
type NetworkMetrics struct {
	ent.Schema
}

// Fields of the NetworkMetrics.
func (NetworkMetrics) Fields() []ent.Field {
	return []ent.Field{
		field.Int("metrics_id").Optional(),
	}
}

// Edges of the NetworkMetrics.
func (NetworkMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object.
		edge.From("metrics", Metrics.Type).
			Ref("network_metrics").
			Unique().
			Field("metrics_id"),

		// Information about host network.
		edge.To("system_network_stats", SystemNetworkStats.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
