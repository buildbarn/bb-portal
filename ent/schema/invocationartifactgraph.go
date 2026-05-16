package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InvocationArtifactGraph stores a zstd-compressed blob of the build's
// artifact graph: every NamedSetOfFiles event plus every TargetCompleted
// event's output-group references, length-prefixed and concatenated.
// Clients decompress and walk the graph themselves; the server only
// hands the payload back through GraphQL.
//
// One row per invocation, written once at BuildFinished, never updated.
type InvocationArtifactGraph struct {
	ent.Schema
}

func (InvocationArtifactGraph) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("payload"),
		field.Int64("uncompressed_size"),
	}
}

func (InvocationArtifactGraph) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("artifact_graph").
			Required().
			Unique(),
	}
}

func (InvocationArtifactGraph) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation").Unique(),
	}
}

func (InvocationArtifactGraph) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

func (InvocationArtifactGraph) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
