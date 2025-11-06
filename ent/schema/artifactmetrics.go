package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ArtifactMetrics holds the schema definition for the ArtifactMetrics entity.
type ArtifactMetrics struct {
	ent.Schema
}

// Fields of the ArtifactMetrics.
func (ArtifactMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Measures all source files newly read this build. Does not include
		// unchanged sources on incremental builds.
		field.Int64("source_artifacts_read_size_in_bytes").Optional(),
		field.Int32("source_artifacts_read_count").Optional(),

		// Measures all output artifacts from executed actions. This includes
		// actions that were cached locally (via the action cache) or remotely (via
		// a remote cache or executor), but does *not* include outputs of actions
		// that were cached internally in Skyframe.
		field.Int64("output_artifacts_seen_size_in_bytes").Optional(),
		field.Int32("output_artifacts_seen_count").Optional(),

		// Measures all output artifacts from actions that were cached locally
		// via the action cache. These artifacts were already present on disk at the
		// start of the build. Does not include Skyframe-cached actions' outputs.
		field.Int64("output_artifacts_from_action_cache_size_in_bytes").Optional(),
		field.Int32("output_artifacts_from_action_cache_count").Optional(),

		// Measures all artifacts that belong to a top-level output group. Does not
		// deduplicate, so if there are two top-level targets in this build that
		// share an artifact, it will be counted twice.
		field.Int64("top_level_artifacts_size_in_bytes").Optional(),
		field.Int32("top_level_artifacts_count").Optional(),
	}
}

// Edges of the ArtifactMetrics.
func (ArtifactMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object
		edge.From("metrics", Metrics.Type).
			Ref("artifact_metrics").
			Unique(),
	}
}

// Indexes of the ArtifactMetrics.
func (ArtifactMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("metrics"),
	}
}
