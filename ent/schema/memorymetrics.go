package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemoryMetrics holds the schema definition for the Blob entity.
type MemoryMetrics struct {
	ent.Schema
}

// Fields of the MemoryMetrics.
func (MemoryMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Size of the JVM heap post build in bytes. This is only collected if
		// --memory_profile is set, since it forces a full GC.
		field.Int64("peak_post_gc_heap_size").Optional(),

		// Size of the peak JVM heap size in bytes post GC. Note that this reports 0
		// if there was no major GC during the build.
		field.Int64("used_heap_size_post_build").Optional(),

		// Size of the peak tenured space JVM heap size event in bytes post GC. Note
		// that this reports 0 if there was no major GC during the build.
		field.Int64("peak_post_gc_tenured_space_heap_size").Optional(),
	}
}

// Edges of MemoryMetrics.
func (MemoryMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the memory metrics object
		edge.From("metrics", Metrics.Type).
			Ref("memory_metrics").
			Unique(),

		// Metrics about garbage collection
		edge.To("garbage_metrics", GarbageMetrics.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the MemoryMetrics.
func (MemoryMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("metrics"),
	}
}

// Mixin of the MemoryMetrics.
func (MemoryMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
