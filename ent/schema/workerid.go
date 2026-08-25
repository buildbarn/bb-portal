package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WorkerID stores one ID associated with a worker process.
type WorkerID struct {
	ent.Schema
}

// Fields of the WorkerID struct.
func (WorkerID) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("worker_id"),
	}
}

// Edges of the WorkerID struct.
func (WorkerID) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("worker_metrics", WorkerMetrics.Type).
			Ref("worker_ids").
			Unique(),
	}
}

// Indexes of the WorkerID struct.
func (WorkerID) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("worker_metrics"),
	}
}

// Mixin of the WorkerID struct.
func (WorkerID) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
