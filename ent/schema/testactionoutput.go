package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TestActionOutput holds the schema definition for the TestActionOutput entity.
type TestActionOutput struct {
	ent.Schema
}

// Fields of the TestActionOutput object.
func (TestActionOutput) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("test_result_id").Immutable(),
		field.Int64("file_id").Immutable(),
	}
}

// Edges of TestActionOutput.
func (TestActionOutput) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("test_result", TestResult.Type).
			Unique().
			Required().
			Immutable().
			Field("test_result_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.To("file", File.Type).
			Unique().
			Required().
			Immutable().
			Field("file_id").
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes for TestActionOutput.
func (TestActionOutput) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("test_result_id", "file_id").Unique(),
		index.Fields("file_id"),
	}
}

// Annotations of the TestActionOutput
func (TestActionOutput) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

// Mixin of the TestActionOutput.
func (TestActionOutput) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
