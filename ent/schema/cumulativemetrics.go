package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CumulativeMetrics holds counts accumulated over a Bazel server's lifetime.
type CumulativeMetrics struct {
	ent.Schema
}

// Fields of the CumulativeMetrics struct.
func (CumulativeMetrics) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("num_analyses").Optional(),
		field.Int32("num_builds").Optional(),
	}
}

// Edges of the CumulativeMetrics struct.
func (CumulativeMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("cumulative_metrics").
			Unique(),
	}
}

// Indexes of the CumulativeMetrics struct.
func (CumulativeMetrics) Indexes() []ent.Index {
	return []ent.Index{}
}

// Mixin of the CumulativeMetrics struct.
func (CumulativeMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
