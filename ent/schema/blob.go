package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Blob holds the schema definition for the Blob entity.
type Blob struct {
	ent.Schema
}

// Fields of the Blob.
func (Blob) Fields() []ent.Field {
	return []ent.Field{
		field.String("uri").Unique().Immutable(),
		field.Int64("size_bytes").Optional(),
		field.Enum("archiving_status").
			Values("QUEUED", "ARCHIVING", "SUCCESS", "FAILED").
			Default("QUEUED"),
		field.String("reason").Optional(),
		field.String("archive_url").Optional(),
	}
}

// Edges of the Blob.
func (Blob) Edges() []ent.Edge {
	return nil
}
