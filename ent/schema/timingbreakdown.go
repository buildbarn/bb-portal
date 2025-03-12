package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TimingBreakdown holds the schema definition for the TimingBreakdown entity.
type TimingBreakdown struct {
	ent.Schema
}

// Fields of the TimingBreakdown.
func (TimingBreakdown) Fields() []ent.Field {
	return []ent.Field{
		// The name of the activity.
		field.String("name").Optional(),

		// Time spent ding the activity (duration).
		// NOTE: proto has this as an int, but implemented as a string
		field.String("time").Optional(),

		field.Int("execution_info_id").Optional(),
	}
}

// Edges of TimingBreakdown.
func (TimingBreakdown) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the execution info object
		edge.From("execution_info", ExectionInfo.Type).
			Ref("timing_breakdown").
			Unique().
			Field("execution_info_id"),

		// Timing children (this could probably be better reempleted as a node to itself.
		// except the relationship to the executio info object.  maybe we don't care about that?
		// for now, an intermediate 'parent' object is used)
		edge.To("child", TimingChild.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
