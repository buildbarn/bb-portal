package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FilesMetric holds the schema definition for the FilesMetric entity.
type FilesMetric struct {
	ent.Schema
}

// Fields of the FilesMetric.
func (FilesMetric) Fields() []ent.Field {
	return []ent.Field{
		// The total size in bytes.
		field.Int64("size_in_bytes").Optional(),

		// The total Coount.
		field.Int32("count").Optional(),
	}
}

// Edges of the FilesMetric.
func (FilesMetric) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("artifact_metrics", ArtifactMetrics.Type).
			// Source Artifacts Read.
			// Measures all source files newly read this build. Does not include
			// unchanged sources on incremental builds.
			Ref("source_artifacts_read").

			// Output Artifacts Seen.
			// Measures all output artifacts from executed actions. This includes
			// actions that were cached locally (via the action cache) or remotely (via
			// a remote cache or executor), but does *not* include outputs of actions
			// that were cached internally in Skyframe.
			Ref("output_artifacts_seen").

			// Output Artifacts From Cache.
			// Measures all output artifacts from actions that were cached locally
			// via the action cache. These artifacts were already present on disk at the
			// start of the build. Does not include Skyframe-cached actions' outputs.
			Ref("output_artifacts_from_action_cache").

			// Top Level Artifacts.
			// Measures all artifacts that belong to a top-level output group. Does not
			// deduplicate, so if there are two top-level targets in this build that
			// share an artifact, it will be counted twice.
			Ref("top_level_artifacts").
			Unique(),
	}
}
