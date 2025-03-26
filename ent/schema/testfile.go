package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TestFile holds the schema definition for the TestFile entity.
type TestFile struct {
	ent.Schema
}

// Fields of the TestResult.
func (TestFile) Fields() []ent.Field {
	return []ent.Field{
		// Digest of the file.
		// using the build tool's configured digest algorithm,
		// hex-encoded.
		field.String("digest").Optional(),

		// TODO: implement file more fully, right now, just storing the URI.
		field.String("file").Optional(),

		// Length of the file in bytes.
		field.Int64("length").Optional(),

		// Identifier indicating the nature of the file (e.g., "stdout", "stderr").
		field.String("name").Optional(),

		// Prefix.
		// A sequence of prefixes to apply to the file name to construct a full path.
		// In most but not all cases, there will be 3 entries:
		//  1. A root output directory, eg "bazel-out"
		//  2. A configuration mnemonic, eg "k8-fastbuild"
		//  3. An output category, eg "genfiles"
		field.Strings("prefix").Optional(),

		field.Int("test_result_id").Optional(),
	}
}

// Edges of TestFile.
func (TestFile) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the test result.
		edge.From("test_result", TestResultBES.Type).
			Ref("test_action_output").
			Unique().
			Field("test_result_id"),
	}
}
