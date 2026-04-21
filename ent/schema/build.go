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

// Build holds the schema definition for the Build.
type Build struct {
	ent.Schema
}

// Fields of the Build.
func (Build) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("build_uuid", uuid.UUID{}).Unique().Immutable(),
		field.Time("timestamp").Annotations(entgql.OrderField("TIMESTAMP")),
	}
}

// Edges of the Build.
func (Build) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("instance_name", InstanceName.Type).
			Ref("builds").
			Unique().
			Required(),

		edge.To("invocations", BazelInvocation.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),

		// Metadata for a Build.
		edge.To("tags", BuildTag.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),
	}
}

// Indexes of the Build.
func (Build) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("build_uuid"),
		index.Fields("timestamp"),
		index.Edges("instance_name"),
	}
}

// Annotations of the Build.
func (Build) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findBuilds"),
	}
}

// Mixin of the Build.
func (Build) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
