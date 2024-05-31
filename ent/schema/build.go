package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Build struct {
	ent.Schema
}

// Fields of the BazelInvocation.
func (Build) Fields() []ent.Field {
	return []ent.Field{
		field.String("build_url").Unique().Immutable(),
		field.UUID("build_uuid", uuid.UUID{}).Unique().Immutable(),
		field.JSON("env", map[string]string{}).Annotations(entgql.Skip()), // NOTE: Uses custom resolver.
	}
}

// Edges of the BazelInvocation.
func (Build) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("invocations", BazelInvocation.Type),
	}
}

func (Build) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("env"),
	}
}

func (Build) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findBuilds"),
	}
}
