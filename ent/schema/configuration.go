package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Configuration holds the schema definition for the Configuration entity.
type Configuration struct {
	ent.Schema
}

// Fields of the Configuration.
func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.String("configuration_id").Immutable(),

		field.String("mnemonic").Optional(),

		field.String("platform_name").Optional(),

		field.String("cpu").Optional(),

		field.JSON("make_variables", map[string]string{}).
			Optional().
			Annotations(
				entgql.Type("Map"),
			),

		field.Bool("is_tool").Optional(),

		field.Int64("bazel_invocation_id").
			Comment("The id of the bazel invocation").
			Immutable().
			Unique().
			Annotations(
				entgql.Skip(),
			),
	}
}

// Edges of the Configuration.
func (Configuration) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("configurations").
			Unique().
			Required().
			Immutable(),

		edge.From("invocation_targets", InvocationTarget.Type).
			Ref("configuration"),

		edge.From("actions", Action.Type).
			Ref("configuration"),
	}
}

// Indexes for Configuration.
func (Configuration) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("configuration_id"),
		index.Edges("bazel_invocation"),
		index.Fields("configuration_id").
			Edges("bazel_invocation").
			Unique(),
	}
}

// Mixin of the Configuration.
func (Configuration) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
