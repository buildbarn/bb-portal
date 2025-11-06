package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// IncompleteBuildLog holds the schema definition for the IncompleteBuildLog entity.
type IncompleteBuildLog struct {
	ent.Schema
}

// Fields of the IncompleteBuildLog.
func (IncompleteBuildLog) Fields() []ent.Field {
	return []ent.Field{
		// The id of the snippet, used for ordering.
		field.Int32("snippet_id"),

		// A log snippet
		field.String("log_snippet"),
	}
}

// Edges of the IncompleteBuildLog.
func (IncompleteBuildLog) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("incomplete_build_logs").
			Unique(),
	}
}

// Indexes for IncompleteBuildLog.
func (IncompleteBuildLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
		index.Fields("snippet_id").
			Edges("bazel_invocation").
			Unique(),
	}
}
