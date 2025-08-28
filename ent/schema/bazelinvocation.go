package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// BazelInvocation holds the schema definition for the BazelInvocation entity.
type BazelInvocation struct {
	ent.Schema
}

// Fields of the BazelInvocation.
func (BazelInvocation) Fields() []ent.Field {
	return []ent.Field{
		// The bazel client invocation ID.
		field.UUID("invocation_id", uuid.UUID{}).Unique().Immutable(),

		// Time the event started.
		field.Time("started_at").Optional().Annotations(entgql.OrderField("STARTED_AT")),

		// Time the event ended
		field.Time("ended_at").Optional(),

		// Rethink? Keep for now to capture existing processing.
		field.Int("change_number").Optional(),

		// Rethink? Keep for now.
		field.Int("patchset_number").Optional(),

		// Build Event Protocol completed successfuly.
		field.Bool("bep_completed").Default(false),

		// Rethink, keep for now.
		// A step label pulled from the metada
		field.String("step_label").Optional(),

		// Email address of the user who launched the invocation if provided.
		field.String("user_email").Optional(),

		// Ldap (username) of the user who launched the invocation if provided.
		field.String("user_ldap").Optional().Annotations(entgql.OrderField("USER_LDAP")),

		// The full logs from the build..
		field.String("build_logs").Optional().Annotations(entgql.Skip(entgql.SkipType)),

		// The cpu type from the configuration event(s).
		field.String("cpu").Optional(),

		// The platform name from the configuration event(s).
		field.String("platform_name").Optional(),

		// The host name from the system where the invocation was launched
		field.String("hostname").Optional(),

		// If this invocation is part of CI
		field.Bool("is_ci_worker").Optional(),

		// The name from the configuration event(s).
		field.String("configuration_mnemonic").Optional(),

		// The number of successful fetch events seen.
		field.Int64("num_fetches").Optional(),

		// The name of the build profile.
		field.String("profile_name").Optional().Annotations(entgql.Skip(entgql.SkipType)),

		// The instance name for the invocation.
		field.String("instance_name").Optional(),

		field.String("bazel_version").Optional(),

		field.String("exit_code_name").Optional(),

		field.Int32("exit_code_code").Optional(),

		field.String("command_line_command").Optional().Annotations(entgql.Skip()),
		field.String("command_line_executable").Optional().Annotations(entgql.Skip()),
		field.String("command_line_residual").Optional().Annotations(entgql.Skip()),
		field.Strings("command_line").Optional().Annotations(entgql.Skip()),
		field.Strings("explicit_command_line").Optional().Annotations(entgql.Skip()),
		field.Strings("startup_options").Optional().Annotations(entgql.Skip()),
		field.Strings("explicit_startup_options").Optional().Annotations(entgql.Skip()),

		// Track which event types have been processed. Used to block duplicate
		// events.
		field.Bool("processed_event_started").Default(false).Annotations(entgql.Skip()),
		field.Bool("processed_event_build_metadata").Default(false).Annotations(entgql.Skip()),
		field.Bool("processed_event_options_parsed").Default(false).Annotations(entgql.Skip()),
		field.Bool("processed_event_build_finished").Default(false).Annotations(entgql.Skip()),
		field.Bool("processed_event_structured_command_line").Default(false).Annotations(entgql.Skip()),
		field.Bool("processed_event_workspace_status").Default(false).Annotations(entgql.Skip()),
	}
}

// Edges of the BazelInvocation.
func (BazelInvocation) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back from the Build.
		edge.From("build", Build.Type).
			Ref("invocations").
			Unique(),

		// Event metadata for all events processed for this invocation.
		edge.To("event_metadata", EventMetadata.Type).
			Annotations(
				entgql.Skip(entgql.SkipType),
				entsql.OnDelete(entsql.Cascade),
			),

		// Info about the grpc connection that this event was sent over.
		edge.To("connection_metadata", ConnectionMetadata.Type).
			Annotations(
				entgql.Skip(entgql.SkipType),
				entsql.OnDelete(entsql.Cascade),
			),

		// Edge to any probles detected.
		// NOTE: Uses custom resolver / types.
		edge.To("problems", BazelInvocationProblem.Type).
			Annotations(
				entgql.Skip(entgql.SkipType),
				entsql.OnDelete(entsql.Cascade),
			),

		// Build Metrics for the Completed Invocation
		edge.To("metrics", Metrics.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			).
			Unique(),

		// Incomplete Build Log snippets for the Invocation
		edge.To("incomplete_build_logs", IncompleteBuildLog.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Files for the Invocation
		edge.To("invocation_files", InvocationFiles.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Test Data for the completed Invocation
		edge.To("test_collection", TestCollection.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Target Data for the completed Invocation
		edge.To("targets", Target.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Edge to source control information
		edge.To("source_control", SourceControl.Type).
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes for Bazel Invocation.
func (BazelInvocation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("invocation_id"),
		index.Fields("change_number", "patchset_number"),
	}
}

// Annotations for basel invocation.
func (BazelInvocation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findBazelInvocations"),
	}
}
