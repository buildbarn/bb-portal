package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
		field.Int32("snippet_id").Immutable(),

		// A log snippet
		field.Bytes("log_snippet").Immutable(),

		// Foreign key to bazel invocation
		field.Int64("bazel_invocation_id").Immutable(),
	}
}

// Edges of the IncompleteBuildLog.
func (IncompleteBuildLog) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("incomplete_build_logs").
			Unique().
			Required().
			Immutable(),
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

// Annotations for IncompleteBuildLog
func (IncompleteBuildLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

// Mixin of the IncompleteBuildLog.
func (IncompleteBuildLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
