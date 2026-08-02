package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PackageMetrics holds package loading metrics for an invocation.
type PackageMetrics struct {
	ent.Schema
}

// Fields of the PackageMetrics struct.
func (PackageMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Number of BUILD files successfully loaded during the build.
		field.Int64("packages_loaded").Optional(),
	}
}

// Edges of the PackageMetrics struct.
func (PackageMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("package_metrics").
			Unique(),
	}
}

// Indexes of the PackageMetrics struct.
func (PackageMetrics) Indexes() []ent.Index {
	return []ent.Index{}
}

// Mixin of the PackageMetrics struct.
func (PackageMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
