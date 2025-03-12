package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Metrics holds the schema definition for the Metrics entity.
type Metrics struct {
	ent.Schema
}

// Fields of the Metrics struct.
func (Metrics) Fields() []ent.Field {
	return []ent.Field{
		field.Int("bazel_invocation_id").Optional(),
	}
}

// Edges of the Metrics.
func (Metrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("metrics").
			Unique().
			Field("bazel_invocation_id"), // New field added

		// The action summmary with details about actions executed.
		edge.To("action_summary", ActionSummary.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Details about memory usage and garbage collections.
		edge.To("memory_metrics", MemoryMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Target metrics.
		edge.To("target_metrics", TargetMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Package metrics.
		edge.To("package_metrics", PackageMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Timing metrics.
		edge.To("timing_metrics", TimingMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Cumulative metrics.
		edge.To("cumulative_metrics", CumulativeMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Artifact metrics.
		edge.To("artifact_metrics", ArtifactMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Network metrics if available.
		edge.To("network_metrics", NetworkMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Dynamic execution metrics if available.
		edge.To("dynamic_execution_metrics", DynamicExecutionMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Build graph metrics.
		edge.To("build_graph_metrics", BuildGraphMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Annotations of the Metrics.
func (Metrics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findMetrics"),
	}
}
