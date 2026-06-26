package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Digest holds the schema definition for the Digest entity.
type Digest struct {
	ent.Schema
}

// Fields of the Digest.
func (Digest) Fields() []ent.Field {
	return []ent.Field{
		// REv2 instance name, not BEP instance name, since the file can have a
		// different instance name than its invocation.
		field.String("rev2_instance_name").
			Immutable(),

		field.Int16("digest_function").
			Immutable().
			Annotations(entgql.Type("String")),

		field.Bytes("hash").
			Immutable().
			Annotations(entgql.Type("String")),

		field.Int64("size_bytes").
			Immutable(),
	}
}

// Edges of the Digest.
func (Digest) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge to a invocation file with this digest.
		edge.To("files", File.Type),
	}
}

// Indexes of the FilePath.
func (Digest) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rev2_instance_name", "digest_function", "hash", "size_bytes").
			Unique(),
	}
}

// Mixin of the Digest.
func (Digest) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
