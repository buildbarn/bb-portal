package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		// Foreign key to a Digest
		//
		// Since entgql seems to require that all entities have atleast one
		// field, we cannot skip this one.
		field.Int64("digest_id").
			Immutable(),

		// Foreign key to a FilePath
		field.Int64("file_path_id").
			Immutable().
			Annotations(entgql.Skip()),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("digest", Digest.Type).
			Field("digest_id").
			Ref("files").
			Unique().
			Required().
			Immutable(),

		edge.From("file_path", FilePath.Type).
			Field("file_path_id").
			Ref("files").
			Unique().
			Required().
			Immutable(),

		edge.From("action_stdout", Action.Type).
			Ref("stdout"),

		edge.From("action_stderr", Action.Type).
			Ref("stderr"),

		edge.From("build_tool_logs", BazelInvocation.Type).
			Ref("build_tool_logs").
			Through("tool_logs", BuildToolLog.Type),

		edge.From("test_action_output", TestResult.Type).
			Ref("test_action_output").
			Through("test_action_output_table", TestActionOutput.Type),
	}
}

// Indexes of the File.
func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("digest"),
		index.Edges(
			"file_path",
			"digest",
		).Unique(),
	}
}

// Mixin of the File.
func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
