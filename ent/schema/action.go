package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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

		field.String("label"),
		field.String("type").Optional(),
		field.String("runner").
			Comment("The runner reported by Bazel's compact execution log").
			Optional(),

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

		// Keep remote references as URIs instead of duplicating CAS metadata in
		// the files and digests tables. file:// URIs are intentionally ignored,
		// because they only make sense on the Bazel client that ran the build.
		field.String("primary_output_uri").Optional(),
		field.String("stdout_uri").Optional(),
		field.String("stderr_uri").Optional(),
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
			Immutable(),

		edge.From("action_digest", Digest.Type).
			Field("action_digest_id").
			Ref("actions").
			Unique(),
	}
}

// Indexes for Action.
func (Action) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("label"),
		index.Edges("bazel_invocation"),
		index.Fields("type").Edges("bazel_invocation"),
		index.Edges("configuration"),
		index.Edges("action_digest"),
	}
}

// Annotations for Action.
func (Action) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findActions"),
	}
}

// Mixin of the Action.
func (Action) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
