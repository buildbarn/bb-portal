package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TargetMetrics holds the schema definition for the Blob entity.
type TargetMetrics struct {
	ent.Schema
}

// Fields of the TargetMetrics.
func (TargetMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Targets Loaded.
		// Size of the JVM heap post build in bytes. This is only collected if
		// --memory_profile is set, since it forces a full GC.
		field.Int64("targets_loaded").Optional(),

		// Targets Configured.
		// Size of the peak JVM heap size in bytes post GC. Note that this reports 0
		// if there was no major GC during the build.
		field.Int64("targets_configured").Optional(),

		// Target Configured Not Including Aspects.
		// Size of the peak tenured space JVM heap size event in bytes post GC. Note
		// that this reports 0 if there was no major GC during the build.
		field.Int64("targets_configured_not_including_aspects").Optional(),
	}
}

// Edges of TargetMetrics.
func (TargetMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object.
		edge.From("metrics", Metrics.Type).
			Ref("target_metrics").
			Unique(),
	}
}

// Indexes of the TargetMetrics.
func (TargetMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("metrics"),
	}
}

// Mixin of the TargetMetrics.
func (TargetMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
