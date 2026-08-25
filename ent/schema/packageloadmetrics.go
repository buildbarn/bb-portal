package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PackageLoadMetrics holds loading metrics for one Bazel package.
type PackageLoadMetrics struct {
	ent.Schema
}

// Fields of the PackageLoadMetrics struct.
func (PackageLoadMetrics) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional(),
		field.Int64("load_duration_in_ns").Optional(),
		field.Uint64("num_targets").Optional(),
		field.Uint64("computation_steps").Optional(),
		field.Uint64("num_transitive_loads").Optional(),
		field.Uint64("package_overhead").Optional(),
		field.Uint64("glob_filesystem_operation_cost").Optional(),
	}
}

// Edges of the PackageLoadMetrics struct.
func (PackageLoadMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("package_metrics", PackageMetrics.Type).
			Ref("package_load_metrics").
			Unique(),
	}
}

// Indexes of the PackageLoadMetrics struct.
func (PackageLoadMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("package_metrics"),
	}
}

// Mixin of the PackageLoadMetrics struct.
func (PackageLoadMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
