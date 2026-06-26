package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildToolLog holds the schema definition for the BuildToolLog entity.
type BuildToolLog struct {
	ent.Schema
}

// Fields of the BuildToolLog object.
func (BuildToolLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("bazel_invocation_id").Immutable(),
		field.Int64("file_id").Immutable(),
	}
}

// Edges of BuildToolLog.
func (BuildToolLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("bazel_invocation", BazelInvocation.Type).
			Unique().
			Required().
			Immutable().
			Field("bazel_invocation_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.To("file", File.Type).
			Unique().
			Required().
			Immutable().
			Field("file_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes for BuildToolLog.
func (BuildToolLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("bazel_invocation_id", "file_id").Unique(),
		index.Fields("file_id"),
	}
}

// Annotations of the BuildToolLog
func (BuildToolLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

// Mixin of the BuildToolLog.
func (BuildToolLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
