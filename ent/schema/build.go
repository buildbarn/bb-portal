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
		field.String("build_url").Immutable(),
		field.UUID("build_uuid", uuid.UUID{}).Unique().Immutable(),
		field.String("instance_name").Immutable(),
		field.Time("timestamp"),
	}
}

// Edges of the Build.
func (Build) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("invocations", BazelInvocation.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the Build.
func (Build) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("build_uuid"),
		index.Fields("build_url"),
		index.Fields("build_url", "instance_name").Unique(),
	}
}

// Annotations of the Build.
func (Build) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findBuilds"),
	}
}
