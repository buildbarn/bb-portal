package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OutputGroup holds the schema definition for the OutputGroup entity.
type OutputGroup struct {
	ent.Schema
}

// Fields of the OutputGroup.
func (OutputGroup) Fields() []ent.Field {
	return []ent.Field{
		// Name of the output group.
		field.String("name").Optional(),

		// Incomplete.
		// Indicates that one or more of the output group's files were not built
		// successfully (the generating action failed).
		field.Bool("incomplete").Optional(),

		field.Int("target_complete_id").Optional(),
	}
}

// Edges of the OutputGroup.
func (OutputGroup) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the target completion object.
		edge.From("target_complete", TargetComplete.Type).
			Ref("output_group").
			Unique().
			Field("target_complete_id"),

		// Inline Files.
		// Inlined files that belong to this output group, requested via
		// --build_event_inline_output_groups.
		edge.To("inline_files", TestFile.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// List of file sets that belong to this output group as well.
		edge.To("file_sets", NamedSetOfFiles.Type).Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
