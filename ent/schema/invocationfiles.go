package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InvocationFiles holds the schema definition for the InvocationFiles entity.
type InvocationFiles struct {
	ent.Schema
}

// Fields of the InvocationFiles.
func (InvocationFiles) Fields() []ent.Field {
	return []ent.Field{
		// Name of the file, including path relative to the invocation root.
		field.String("name"),

		// Content of the file, if available.
		field.String("content").Optional(),

		// Digest of the file, if available.
		field.String("digest").Optional(),

		// SizeBytes is the size of the file in bytes, if available.
		field.Int64("size_bytes").Optional(),

		// DigestFunction is the function used to compute the digest, in lower case. Defaults to "sha256".
		field.String("digest_function").Optional(),
	}
}

// Edges of the InvocationFiles.
func (InvocationFiles) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the bazel invocation.
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("invocation_files").
			Unique(),
	}
}

// Indexes of the InvocationFiles.
func (InvocationFiles) Indexes() []ent.Index {
	return []ent.Index{
		// Index making the combination of a name and invocation unique.
		index.Fields("name").
			Edges("bazel_invocation").
			Unique(),
	}
}
