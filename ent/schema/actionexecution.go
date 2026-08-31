package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ActionExecution holds the schema definition for an action observed during a
// Bazel invocation.
type ActionExecution struct {
	ent.Schema
}

// Fields of the ActionExecution.
func (ActionExecution) Fields() []ent.Field {
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
			Optional().
			Annotations(
				entgql.Skip(),
			),

		field.Int64("action_digest_id").
			Comment("The REv2 Action digest obtained from Bazel's execution log").
			Optional().
			Annotations(
				entgql.Skip(),
			),

		field.Int64("primary_output_file_id").
			Comment("The normalized CAS-backed primary output file").
			Optional().
			Immutable().
			Annotations(entgql.Skip()),

		field.Int64("stdout_file_id").
			Comment("The normalized CAS-backed standard output file").
			Optional().
			Immutable().
			Annotations(entgql.Skip()),

		field.Int64("stderr_file_id").
			Comment("The normalized CAS-backed standard error file").
			Optional().
			Immutable().
			Annotations(entgql.Skip()),

		field.String("label"),
		field.String("type").Optional(),
		field.String("runner").
			Comment("The runner reported by Bazel's compact execution log").
			Optional(),
		field.Bool("cache_hit").
			Comment("Whether Bazel's compact execution log reported a disk or remote cache hit").
			Optional().
			Nillable(),

		field.Bool("success").Optional(),
		field.Int32("exit_code").Optional(),

		field.Strings("command_line").Optional(),

		field.Time("start_time").Optional(),
		field.Time("end_time").Optional(),

		field.String("failure_code").Optional(),
		field.String("failure_message").Optional(),

		// The path from BuildEventId.ActionCompletedId. This is available for
		// successful actions even when the primary output's File message only
		// contains a URI.
		field.String("primary_output").Optional(),
	}
}

// Edges of the ActionExecution.
func (ActionExecution) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the Bazel invocation in which this execution was observed.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Field("bazel_invocation_id").
			Ref("action_executions").
			Unique().
			Required().
			Immutable(),

		// Edge to the configuration.
		edge.To("configuration", Configuration.Type).
			Field("configuration_id").
			Unique().
			Immutable(),

		edge.From("action_digest", Digest.Type).
			Field("action_digest_id").
			Ref("action_executions").
			Unique(),

		edge.To("primary_output_file", File.Type).
			Field("primary_output_file_id").
			Unique().
			Immutable(),

		edge.To("stdout", File.Type).
			Field("stdout_file_id").
			Unique().
			Immutable(),

		edge.To("stderr", File.Type).
			Field("stderr_file_id").
			Unique().
			Immutable(),
	}
}

// Indexes for ActionExecution.
func (ActionExecution) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("label"),
		index.Edges("bazel_invocation"),
		index.Fields("type").Edges("bazel_invocation"),
		index.Edges("configuration"),
		index.Edges("action_digest"),
		index.Edges("primary_output_file"),
		index.Edges("stdout"),
		index.Edges("stderr"),
	}
}

// Annotations for ActionExecution.
func (ActionExecution) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}

// Mixin of the ActionExecution.
func (ActionExecution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
