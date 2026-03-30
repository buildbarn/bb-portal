package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TestTarget holds the schema definition for the TestTarget entity.
type TestTarget struct {
	ent.Schema
}

// Fields of the TestTarget.
func (TestTarget) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("target_id").Immutable(),
	}
}

// Edges of the TestTarget.
func (TestTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("target", Target.Type).
			Ref("test_target").
			Field("target_id").
			Unique().
			Required().
			Immutable().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Mixin of the TestTarget.
func (TestTarget) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
