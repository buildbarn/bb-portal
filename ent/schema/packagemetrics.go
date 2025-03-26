package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PackageMetrics holds the schema definition for the PackageMetrics entity.
type PackageMetrics struct {
	ent.Schema
}

// Fields of the PackageMetrics.
func (PackageMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Packages Loaded Count.
		// Number of BUILD files (aka packages) successfully loaded during this
		// build.
		//
		// [For Bazel binaries built at source states] Before Dec 2021, this value
		// was the number of packages attempted to be loaded, for a particular
		// definition of "attempted".
		//
		// After Dec 2021, this value would sometimes overcount because the same
		// package could sometimes be attempted to be loaded multiple times due to
		// memory pressure.
		//
		// After Feb 2022, this value is the number of packages successfully
		// loaded.
		field.Int64("packages_loaded").Optional(),
		field.Int("metrics_id").Optional(),
	}
}

// Edges of PackageMetrics.
func (PackageMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object.
		edge.From("metrics", Metrics.Type).
			Ref("package_metrics").
			Unique().
			Field("metrics_id"),

		// Loading time metrics per package.
		edge.To("package_load_metrics", PackageLoadMetrics.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
