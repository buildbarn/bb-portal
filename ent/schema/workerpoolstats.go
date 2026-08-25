package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WorkerPoolStats holds lifecycle counts for one worker pool.
type WorkerPoolStats struct {
	ent.Schema
}

// Fields of the WorkerPoolStats struct.
func (WorkerPoolStats) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("hash").Optional(),
		field.String("mnemonic").Optional(),
		field.Int64("created_count").Optional(),
		field.Int64("destroyed_count").Optional(),
		field.Int64("evicted_count").Optional(),
		field.Int64("user_exec_exception_destroyed_count").Optional(),
		field.Int64("io_exception_destroyed_count").Optional(),
		field.Int64("interrupted_exception_destroyed_count").Optional(),
		field.Int64("unknown_destroyed_count").Optional(),
		field.Int64("alive_count").Optional(),
	}
}

// Edges of the WorkerPoolStats struct.
func (WorkerPoolStats) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("worker_pool_metrics", WorkerPoolMetrics.Type).
			Ref("worker_pool_stats").
			Unique(),
	}
}

// Indexes of the WorkerPoolStats struct.
func (WorkerPoolStats) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("worker_pool_metrics"),
	}
}

// Mixin of the WorkerPoolStats struct.
func (WorkerPoolStats) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
