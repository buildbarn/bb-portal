package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

// Metrics holds the schema definition for the Metrics entity.
type Metrics struct {
	ent.Schema
}

// Fields of the Metrics struct.
func (Metrics) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the Metrics.
func (Metrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("metrics").
			Unique(),

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

		// Timing metrics.
		edge.To("timing_metrics", TimingMetrics.Type).
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

		// Build graph metrics.
		edge.To("build_graph_metrics", BuildGraphMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Package loading metrics.
		edge.To("package_metrics", PackageMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Counts accumulated over the lifetime of the Bazel server.
		edge.To("cumulative_metrics", CumulativeMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Workers that were alive during the invocation.
		edge.To("worker_metrics", WorkerMetrics.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Worker pool lifecycle metrics.
		edge.To("worker_pool_metrics", WorkerPoolMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Results of local and remote dynamic execution races.
		edge.To("dynamic_execution_metrics", DynamicExecutionMetrics.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the Metrics.
func (Metrics) Indexes() []ent.Index {
	return []ent.Index{}
}

// Mixin of the Metrics.
func (Metrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
