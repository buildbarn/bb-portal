package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildTag holds the schema definition for the BuildTag entity.
type BuildTag struct {
	ent.Schema
}

// Fields of the BuildTag object.
func (BuildTag) Fields() []ent.Field {
	return []ent.Field{
		// Foreign key to build
		field.Int64("build_id").
			Immutable().
			Annotations(
				entgql.Skip(),
			),

		// A tag consists of a key-value pair
		field.String("key").Immutable().Annotations(entgql.OrderField("KEY")),
		field.String("value").Immutable(),
	}
}

// Edges of BuildTag.
func (BuildTag) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the build
		edge.From("build", Build.Type).
			Field("build_id").
			Ref("tags").
			Unique().
			Required().
			Immutable(),
	}
}

// Indexes for BuildTag.
func (BuildTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("build"),

		// Duplicate keys are allowed for a build, as long as they have different values.
		index.Fields("key", "value").
			Edges("build").
			Unique(),
	}
}

// Annotations of the BuildTag
func (BuildTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}

// Mixin of the BuildTag.
func (BuildTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
