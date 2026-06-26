package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// FilePath holds the schema definition for the FilePath entity.
type FilePath struct {
	ent.Schema
}

// Fields of the FilePath.
func (FilePath) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("bep_instance_name_id").
			Immutable().
			Annotations(entgql.Skip()),

		field.String("path").
			Immutable(),
	}
}

// Edges of the FilePath.
func (FilePath) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bep_instance_name", InstanceName.Type).
			Field("bep_instance_name_id").
			Ref("file_paths").
			Unique().
			Required().
			Immutable(),

		edge.To("files", File.Type),
	}
}

// Indexes of the FilePath.
func (FilePath) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("bep_instance_name_id", "path").
			Unique(),
	}
}

// Mixin of the FilePath.
func (FilePath) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
