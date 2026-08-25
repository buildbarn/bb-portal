package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildGraphRuleClassCount holds configured target and action counts for a rule class.
type BuildGraphRuleClassCount struct {
	ent.Schema
}

// Fields of the BuildGraphRuleClassCount struct.
func (BuildGraphRuleClassCount) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").Optional(),
		field.String("rule_class").Optional(),
		field.Uint64("count").Optional(),
		field.Uint64("action_count").Optional(),
	}
}

// Edges of the BuildGraphRuleClassCount struct.
func (BuildGraphRuleClassCount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("build_graph_metrics", BuildGraphMetrics.Type).
			Ref("rule_class_counts").
			Unique(),
	}
}

// Indexes of the BuildGraphRuleClassCount struct.
func (BuildGraphRuleClassCount) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("build_graph_metrics"),
	}
}

// Mixin of the BuildGraphRuleClassCount struct.
func (BuildGraphRuleClassCount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
