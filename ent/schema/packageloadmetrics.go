package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PackageLoadMetrics holds the schema definition for the Blob entity.
type PackageLoadMetrics struct {
	ent.Schema
}

// Fields of the PackageLoadMetrics.
func (PackageLoadMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Package Name
		field.String("name").Optional(),

		// How long it took to load this package.
		field.Int64("load_duration").
			Optional(),

		// Bmber of targets using the package.
		field.Uint64("num_targets").Optional(),

		// Computation steps for the package.
		field.Uint64("computation_steps").Optional(),

		// Transitive loads.
		field.Uint64("num_transitive_loads").Optional(),

		// Package overhead.
		field.Uint64("package_overhead").Optional(),
	}
}

// Edges of PackageLoadMetrics.
func (PackageLoadMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the package metrics
		edge.From("package_metrics", PackageMetrics.Type).
			Ref("package_load_metrics").
			Unique(),
	}
}

// Indexes of the PackageLoadMetrics.
func (PackageLoadMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("package_metrics"),
	}
}
