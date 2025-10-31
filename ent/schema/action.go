package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Action holds the schema definition for the Action entity.
type Action struct {
	ent.Schema
}

// Fields of the Action.
func (Action) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("bazel_invocation_id").
			Comment("The id of the bazel invocation").
			Immutable().
			Unique().
			Annotations(
				entgql.Skip(),
			),

		field.Int64("configuration_id").
			Comment("The id of the configuration").
			Immutable().
			Unique().
			Annotations(
				entgql.Skip(),
			),

		field.String("label"),
		field.String("type").Optional(),

		field.Bool("success").Optional(),
		field.Int32("exit_code").Optional(),

		field.Strings("command_line").Optional(),

		field.Time("start_time").Optional(),
		field.Time("end_time").Optional(),

		field.String("failure_code").Optional(),
		field.String("failure_message").Optional(),

		field.String("stdout_hash").Optional(),
		field.Int64("stdout_size_bytes").Optional(),
		field.String("stdout_hash_function").Optional(),
		field.String("stderr_hash").Optional(),
		field.Int64("stderr_size_bytes").Optional(),
		field.String("stderr_hash_function").Optional(),
	}
}

// Edges of the Action.
func (Action) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("actions").
			Unique().
			Required().
			Immutable(),

		// Edge to the configuration.
		edge.To("configuration", Configuration.Type).
			Field("configuration_id").
			Unique().
			Required().
			Immutable(),
	}
}

// Indexes for Action.
func (Action) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("label"),
		index.Edges("bazel_invocation"),
	}
}

// Mixin of the Action.
func (Action) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
