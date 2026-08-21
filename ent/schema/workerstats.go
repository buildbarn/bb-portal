package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WorkerStats holds a point-in-time worker measurement.
type WorkerStats struct {
	ent.Schema
}

// Fields of the WorkerStats struct.
func (WorkerStats) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("collect_time_in_ms").Optional(),
		field.Int32("worker_memory_in_kb").Optional(),
		field.Int32("prior_worker_memory_in_kb").Optional(),
		field.Int64("last_action_start_time_in_ms").Optional(),
	}
}

// Edges of the WorkerStats struct.
func (WorkerStats) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("worker_metrics", WorkerMetrics.Type).
			Ref("worker_stats").
			Unique(),
	}
}

// Indexes of the WorkerStats struct.
func (WorkerStats) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("worker_metrics"),
	}
}

// Mixin of the WorkerStats struct.
func (WorkerStats) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
