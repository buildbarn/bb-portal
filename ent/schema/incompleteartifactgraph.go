package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// IncompleteArtifactGraph holds a single serialized BEP BuildEvent
// (NamedSetOfFiles or TargetCompleted variant) that contributes to an
// invocation's artifact graph. Rows accumulate in the database as events
// stream in — the recorder stays stateless, so this survives failover and
// the graph is queryable in its partial state mid-build, mirroring how
// IncompleteBuildLog accumulates progress events.
//
// After the build finishes, dbcleanupservice.CompactArtifactGraphs folds
// these rows into a single compressed InvocationArtifactGraph blob and the
// rows are deleted, the same way incomplete build logs are compacted into
// BuildLogChunks.
type IncompleteArtifactGraph struct {
	ent.Schema
}

// Fields of the IncompleteArtifactGraph.
func (IncompleteArtifactGraph) Fields() []ent.Field {
	return []ent.Field{
		// The BEP sequence number of the event, used for ordering and as
		// the per-invocation idempotency key.
		field.Int32("seq_id").Immutable(),

		// A single serialized bes.BuildEvent.
		field.Bytes("event").Immutable(),

		// Foreign key to bazel invocation
		field.Int64("bazel_invocation_id").Immutable(),
	}
}

// Edges of the IncompleteArtifactGraph.
func (IncompleteArtifactGraph) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("incomplete_artifact_graphs").
			Unique().
			Required().
			Immutable(),
	}
}

// Indexes for IncompleteArtifactGraph.
func (IncompleteArtifactGraph) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
		index.Fields("seq_id").
			Edges("bazel_invocation").
			Unique(),
	}
}

// Annotations for IncompleteArtifactGraph.
func (IncompleteArtifactGraph) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

// Mixin of the IncompleteArtifactGraph.
func (IncompleteArtifactGraph) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
