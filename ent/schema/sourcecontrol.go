package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SourceControl holds the schema definition for the SourceControl entity.
type SourceControl struct {
	ent.Schema
}

// Fields of the SourceControl object.
func (SourceControl) Fields() []ent.Field {
	return []ent.Field{
		// The Repository Url associated wth the invocation
		field.String("repo_url").Optional(),

		// The Branch associated with the invocation
		field.String("branch").Optional(),

		// The Commit SHA of the invocation
		field.String("commit_sha").Optional(),

		// The source control actor associated with the invocation
		field.String("actor").Optional(),
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
