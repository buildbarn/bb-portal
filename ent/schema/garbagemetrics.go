package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GarbageMetrics holds the schema definition for the GarbageMetrics entity.
type GarbageMetrics struct {
	ent.Schema
}

// Fields of the GarbageMetrics.
func (GarbageMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Type of garbage collected, e.g. G1 Old Gen.
		field.String("type").Optional(),

		// Number of bytes of garbage of the given type collected during this invocation
		field.Int64("garbage_collected").Optional(),
	}
}

// Edges of GarbageMetrics.
func (GarbageMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the memory metrics object
		edge.From("memory_metrics", MemoryMetrics.Type).
			Ref("garbage_metrics").
			Unique(),
	}
}
