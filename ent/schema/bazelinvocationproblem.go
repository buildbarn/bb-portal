package schema

import (
	"encoding/json"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BazelInvocationProblem holds the schema definition for the BazelInvocationProblem entity.
type BazelInvocationProblem struct {
	ent.Schema
}

// Fields of the BazelInvocationProblem.
func (BazelInvocationProblem) Fields() []ent.Field {
	return []ent.Field{
		field.String("problem_type"),
		field.String("label"),
		field.JSON("bep_events", json.RawMessage{}).Annotations(entgql.Skip()), // NOTE: Internal model, not exposed to API.
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
