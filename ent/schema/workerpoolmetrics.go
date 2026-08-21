package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

// WorkerPoolMetrics groups worker pool lifecycle statistics.
type WorkerPoolMetrics struct {
	ent.Schema
}

// Fields of the WorkerPoolMetrics struct.
func (WorkerPoolMetrics) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the WorkerPoolMetrics struct.
func (WorkerPoolMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("worker_pool_metrics").
			Unique(),
		edge.To("worker_pool_stats", WorkerPoolStats.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the WorkerPoolMetrics struct.
func (WorkerPoolMetrics) Indexes() []ent.Index {
	return []ent.Index{}
}

// Mixin of the WorkerPoolMetrics struct.
func (WorkerPoolMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
