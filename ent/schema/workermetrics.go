package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WorkerMetrics holds information about a worker alive during an invocation.
type WorkerMetrics struct {
	ent.Schema
}

// Fields of the WorkerMetrics struct.
func (WorkerMetrics) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("process_id").Optional(),
		field.String("mnemonic").Optional(),
		field.Bool("is_multiplex").Optional(),
		field.Bool("is_sandbox").Optional(),
		field.Bool("is_measurable").Optional(),
		field.Int64("worker_key_hash").Optional(),
		field.String("worker_status").Optional(),
		field.String("code").Optional(),
		field.Int64("actions_executed").Optional(),
		field.Int64("prior_actions_executed").Optional(),
	}
}

// Edges of the WorkerMetrics struct.
func (WorkerMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("worker_metrics").
			Unique(),
		edge.To("worker_ids", WorkerID.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.To("worker_stats", WorkerStats.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the WorkerMetrics struct.
func (WorkerMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("metrics"),
	}
}

// Mixin of the WorkerMetrics struct.
func (WorkerMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
