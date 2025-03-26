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

	"github.com/buildbarn/bb-portal/pkg/summary"
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
		field.Time("started_at").Annotations(entgql.OrderField("STARTED_AT")),

		// Time the event ended
		field.Time("ended_at").Optional(),

		// Rethink? Keep for now to capture existing processing.
		field.Int("change_number").Optional(),

		// Rethink? Keep for now.
		field.Int("patchset_number").Optional(),

		// NOTE: Internal model, not exposed to API.
		// contains invocation information
		field.JSON("summary", summary.InvocationSummary{}).Annotations(entgql.Skip()),

		// Build Event Protocol completed successfuly.
		field.Bool("bep_completed").Optional(),

		// Rethink, keep for now.
		// A step label pulled from the metada
		field.String("step_label"),

		// NOTE: Uses custom resolver.
		// Log snippets of error saved to disk.  Rethink and store in db?
		field.JSON("related_files", map[string]string{}).Annotations(entgql.Skip()),

		// Email address of the user who launched the invocation if provided.
		field.String("user_email").Optional(),

		// Ldap (username) of the user who launched the invocation if provided.
		field.String("user_ldap").Optional().Annotations(entgql.OrderField("USER_LDAP")),

		// The full logs from the build..
		field.String("build_logs").Optional(),

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
		field.String("profile_name").Annotations(entgql.Skip(entgql.SkipType)),

		field.Int("event_file_id").Optional(),
		field.Int("build_id").Optional(),
	}
}

// Edges of the BazelInvocation.
func (BazelInvocation) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back from the Event Files.
		edge.From("event_file", EventFile.Type).
			Ref("bazel_invocation").
			Unique().
			// Required().
			Field("event_file_id"),

		// Edge back from the Build.
		edge.From("build", Build.Type).
			Ref("invocations").
			Unique().
			Field("build_id"),

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

		// Test Data for the completed Invocation
		edge.To("test_collection", TestCollection.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Target Data for the completed Invocation
		edge.To("targets", TargetPair.Type).
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
		index.Fields("change_number", "patchset_number"),
		index.Fields("invocation_id"),
	}
}

// Annotations for basel invocation.
func (BazelInvocation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findBazelInvocations"),
	}
}
