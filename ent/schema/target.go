package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Target holds the schema definition for the Target entity.
type Target struct {
	ent.Schema
}

// Fields of the Target.
func (Target) Fields() []ent.Field {
	return []ent.Field{
		// The label of the target ex: //foo:bar.
		field.String("label"),

		// The aspect of the target completion if any.
		field.String("aspect"),

		// The kind of target.
		// (e.g.,  e.g. "cc_library rule", "source file",
		// "generated file") where the completion is reported.
		field.String("target_kind"),
	}
}

// Edges of the Target.
func (Target) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("instance_name", InstanceName.Type).
			Ref("targets").
			Unique().
			Required(),

		// Target Data for the completed Invocation
		edge.To("invocation_targets", InvocationTarget.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entgql.RelayConnection(),
			),

		edge.To("target_kind_mappings", TargetKindMapping.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the Target.
func (Target) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("instance_name"),
		index.Fields("label", "aspect"),
		index.Fields("label", "aspect", "target_kind").
			Edges("instance_name").
			Unique(),
	}
}

// Annotations of the Target
func (Target) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("findTargets"),
	}
}
