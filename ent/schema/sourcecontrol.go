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

		// The source control refs associated with the invocation
		field.String("refs").Optional(),

		// The source control run id associated with the invocation
		field.String("run_id").Optional(),

		// The source control workflow associated with the invocation
		field.String("workflow").Optional(),

		// The source control action associated with the invocation
		field.String("action").Optional(),

		// The source control workspace associated with the invocation
		field.String("workspace").Optional(),

		// The source control event name associated with the invocation
		field.String("event_name").Optional(),

		// The source control job associated with the invocation
		field.String("job").Optional(),

		// The source control job associated with the invocation (Possible duplicate)
		field.String("runner_name").Optional(),

		// The source control runner architecture associated with the invocation (Possible duplicate)
		field.String("runner_arch").Optional(),

		// The source control runner architecture associated with the invocation (Possible duplicate)
		field.String("runner_os").Optional(),

		// The bazel invocation id associated with the invocation
		field.Int("bazel_invocation_id").Optional(),
	}
}

// Edges of SourceControl.
func (SourceControl) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("source_control").
			Unique().
			Field("bazel_invocation_id"),
	}
}
