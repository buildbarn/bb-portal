package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildGraphEvaluationStat holds a Skyframe operation count for a Skyfunction.
type BuildGraphEvaluationStat struct {
	ent.Schema
}

// Fields of the BuildGraphEvaluationStat struct.
func (BuildGraphEvaluationStat) Fields() []ent.Field {
	return []ent.Field{
		field.String("operation").Optional(),
		field.String("skyfunction_name").Optional(),
		field.Int64("count").Optional(),
	}
}

// Edges of the BuildGraphEvaluationStat struct.
func (BuildGraphEvaluationStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("build_graph_metrics", BuildGraphMetrics.Type).
			Ref("evaluation_stats").
			Unique(),
	}
}

// Indexes of the BuildGraphEvaluationStat struct.
func (BuildGraphEvaluationStat) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("build_graph_metrics"),
	}
}

// Mixin of the BuildGraphEvaluationStat struct.
func (BuildGraphEvaluationStat) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
