package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// AuthenticatedUser holds the schema definition for the AuthenticatedUser entity.
type AuthenticatedUser struct {
	ent.Schema
}

// Fields of the AuthenticatedUser.
func (AuthenticatedUser) Fields() []ent.Field {
	return []ent.Field{
		// A UUID generated when the object is first created, used to fetch the
		// user in GraphQL. Will not be change in the upsert in
		// `findOrCreateAuthenticatedUser`.
		field.UUID("user_uuid", uuid.UUID{}).Unique().Immutable(),

		// A string which the external identity provider uses to identify the user.
		// This field is used as the index for upserting AuthenticatedUser objects
		// in `findOrCreateAuthenticatedUser`.
		// Defined by the bb-portal configuration key `external_id_extraction_jmespath_expression`.
		field.String("external_id").Unique().Immutable(),

		// A human-readable name for displaying the user.
		// Defined by the bb-portal configuration key `display_name_extraction_jmespath_expression`.
		// If the extractor returns `nil` (i.e. the configured key is no longer present in
		// the authentication metadata, or if the extraction otherwise fails), this field will
		// not change.
		field.String("display_name").Optional(),

		// Arbitrary structured data about the user.
		// Defined by the bb-portal configuration key `user_info_extraction_jmespath_expression`
		// If the extractor returns `nil` (i.e. the configured key is no longer present in
		// the authentication metadata, or if the extraction otherwise fails), this field will
		// not change.
		field.JSON("user_info", map[string]any{}).Optional(),
	}
}

// Indexes for AuthenticatedUser.
func (AuthenticatedUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_uuid"),
	}
}

// Edges of the AuthenticatedUser.
func (AuthenticatedUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("bazel_invocations", BazelInvocation.Type).
			Annotations(entgql.RelayConnection()),
	}
}
