package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildGraphAspectCount holds configured target and action counts for an aspect.
type BuildGraphAspectCount struct {
	ent.Schema
}

// Fields of the BuildGraphAspectCount struct.
func (BuildGraphAspectCount) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").Optional(),
		field.String("aspect_name").Optional(),
		field.Uint64("count").Optional(),
		field.Uint64("action_count").Optional(),
	}
}

// Edges of the BuildGraphAspectCount struct.
func (BuildGraphAspectCount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("build_graph_metrics", BuildGraphMetrics.Type).
			Ref("aspect_counts").
			Unique(),
	}
}

// Indexes of the BuildGraphAspectCount struct.
func (BuildGraphAspectCount) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("build_graph_metrics"),
	}
}

// Mixin of the BuildGraphAspectCount struct.
func (BuildGraphAspectCount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
