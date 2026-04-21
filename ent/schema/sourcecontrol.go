package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SourceControl holds the schema definition for the SourceControl entity.
type SourceControl struct {
	ent.Schema
}

// Fields of the SourceControl object.
func (SourceControl) Fields() []ent.Field {
	return []ent.Field{
		// The repo used for a invocation
		field.String("repo").Optional(),
		field.String("repo_url").Optional(),

		// Git ref used for the invocation, such as branch or pull request
		field.String("ref").Optional(),
		field.String("ref_url").Optional(),

		// Commit used for the invocation
		field.String("commit").Optional(),
		field.String("commit_url").Optional(),
	}
}

// Edges of SourceControl.
func (SourceControl) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("source_control").
			Unique(),
	}
}

// Indexes for SourceControl.
func (SourceControl) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
	}
}

// Mixin of the SourceControl.
func (SourceControl) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
