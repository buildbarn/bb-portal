package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

// NamedSetOfFiles holds the schema definition for the NamedSetOfFiles entity.
type NamedSetOfFiles struct {
	ent.Schema
}

// Fields of the NamedSetOfFiles.
func (NamedSetOfFiles) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the NamedSetOfFiles.
func (NamedSetOfFiles) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to output group.
		edge.From("output_group", OutputGroup.Type).
			Ref("file_sets").
			Unique(),

		// Files that belong to this named set of files.
		edge.To("files", TestFile.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// Other named sets whose members also belong to this set.
		edge.To("file_sets", NamedSetOfFiles.Type).Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Payload of a message to describe a set of files, usually build artifacts, to
// be referred to later by their name. In this way, files that occur identically
// as outputs of several targets have to be named only once.
