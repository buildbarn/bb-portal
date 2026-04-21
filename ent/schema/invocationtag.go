package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InvocationTag holds the schema definition for the InvocationTag entity.
type InvocationTag struct {
	ent.Schema
}

// Fields of the InvocationTag object.
func (InvocationTag) Fields() []ent.Field {
	return []ent.Field{
		// Foreign key to bazel invocation
		field.Int64("bazel_invocation_id").
			Immutable().
			Annotations(
				entgql.Skip(),
			),

		// A tag consists of a key-value pair
		field.String("key").Immutable().Annotations(entgql.OrderField("KEY")),
		field.String("value").Immutable(),
	}
}

// Edges of InvocationTag.
func (InvocationTag) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("tags").
			Unique().
			Required().
			Immutable(),
	}
}

// Indexes for InvocationTag.
func (InvocationTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
		index.Fields("key").
			Edges("bazel_invocation").
			Unique(),
	}
}

// Annotations of the InvocationTag
func (InvocationTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}

// Mixin of the InvocationTag.
func (InvocationTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
