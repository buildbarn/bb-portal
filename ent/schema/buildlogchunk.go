package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildLogChunk holds a chunk of a bazel build log.
type BuildLogChunk struct {
	ent.Schema
}

// Fields of the BuildLogChunk.
func (BuildLogChunk) Fields() []ent.Field {
	return []ent.Field{
		// The encoded data of the chunk, this is raw bytes that have
		// been compressed with zstd go side after chunking. Once
		// decompressed the data may not be split according to valid
		// utf-8 delimitations and to fully decode the neighbouring
		// chunks may be necessary.
		field.Bytes("data"),
		// The index of the chunk for within this build.
		field.Int("chunk_index"),
		// The inclusive index of the first line in this chunk. Note,
		// this may be a line started in an earlier chunk.
		field.Int64("first_line_index"),
		// The inclusive index of the last line in this chunk. Note,
		// this line may end in a later chunk.
		field.Int64("last_line_index"),
	}
}

// Edges of the BuildLogChunk.
func (BuildLogChunk) Edges() []ent.Edge {
	return []ent.Edge{
		// The chunks in this table belong to one and only one bazel
		// invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("build_log_chunks").
			Unique().
			Required(),
	}
}

// Indexes of the BuildLogChunk.
func (BuildLogChunk) Indexes() []ent.Index {
	return []ent.Index{
		// The chunk index for a specific bazel invocation is unique
		// (and ordered and starts at zero but this is not encoded in
		// the database model).
		index.Edges("bazel_invocation").Fields("chunk_index").Unique(),
	}
}

// Annotations of the BuildLogChunk.
func (BuildLogChunk) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Skip(),
	}
}

// Mixin of the BuildLogChunk.
func (BuildLogChunk) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
