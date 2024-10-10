package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
			Unique(),

		// Details about memory usage and garbage collections.
		edge.To("memory_metrics", MemoryMetrics.Type).
			Unique(),

		// Target metrics.
		edge.To("target_metrics", TargetMetrics.Type).
			Unique(),

		// Package metrics.
		edge.To("package_metrics", PackageMetrics.Type).
			Unique(),

		// Timing metrics.
		edge.To("timing_metrics", TimingMetrics.Type).
			Unique(),

		// Cumulative metrics.
		edge.To("cumulative_metrics", CumulativeMetrics.Type).
			Unique(),

		// Artifact metrics.
		edge.To("artifact_metrics", ArtifactMetrics.Type).
			Unique(),

		// Network metrics if available.
		edge.To("network_metrics", NetworkMetrics.Type).
			Unique(),

		// Dynamic execution metrics if available.
		edge.To("dynamic_execution_metrics", DynamicExecutionMetrics.Type).
			Unique(),

		// Build graph metrics.
		edge.To("build_graph_metrics", BuildGraphMetrics.Type).
			Unique(),
	}
}

// Annotations of the Metrics.
func (Metrics) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findMetrics"),
	}
}
