package schema

import (
	"encoding/json"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BazelInvocationProblem holds the schema definition for the BazelInvocationProblem entity.
type BazelInvocationProblem struct {
	ent.Schema
}

// Fields of the BazelInvocationProblem.
func (BazelInvocationProblem) Fields() []ent.Field {
	return []ent.Field{
		// The Problem Type.
		field.String("problem_type"),

		// The Problem Label.
		field.String("label"),

		// The bep_events raw message associated with the field.
		// NOTE: Internal model, not exposed to API.
		field.JSON("bep_events", json.RawMessage{}).Annotations(entgql.Skip()),
	}
}

// Edges of the BazelInvocationProblem.
func (BazelInvocationProblem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bazel_invocation", BazelInvocation.Type).
			Ref("problems").
			Unique(),
	}
}

// Indexes for BazelInvocationProblems.
func (BazelInvocationProblem) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("bazel_invocation"),
	}
}
