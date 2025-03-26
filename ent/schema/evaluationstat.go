package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EvaluationStat holds the schema definition for the EvaluationStat entity.
type EvaluationStat struct {
	ent.Schema
}

// Fields of the EvaluationStat.
func (EvaluationStat) Fields() []ent.Field {
	return []ent.Field{
		// Name of the Skyfunction.
		field.String("skyfunction_name").Optional(),

		// How many times a given operation was carried out on a Skyfunction.
		field.Int64("count").Optional(),

		field.Int("build_graph_metrics_id").Optional(),
	}
}

// Edges of the EvaluationStat.
func (EvaluationStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("build_graph_metrics", BuildGraphMetrics.Type).

			// NOTE: Not populated on the proto currently, but included here for completeness.

			// Dirtied Values.
			// Number of SkyValues that were dirtied during the build. Dirtied nodes are
			// those that transitively depend on a node that changed by itself (e.g. one
			// representing a file in the file system)
			Ref("dirtied_values").

			// Changed Values.
			// Number of SkyValues that changed by themselves. For example, when a file
			// on the file system changes, the SkyValue representing it will change.
			Ref("changed_values").

			// Built Values.
			// Number of SkyValues that were built. This means that they were evaluated
			// and were found to have changed from their previous version.
			Ref("built_values").

			// Cleaned Values.
			// Number of SkyValues that were evaluated and found clean, i.e. equal to
			// their previous version.
			Ref("cleaned_values").

			// Evaluated Values.
			// Number of evaluations to build SkyValues. This includes restarted
			// evaluations, which means there can be multiple evaluations per built
			// SkyValue. Subtract built_values from this number to get the number of
			// restarted evaluations.
			Ref("evaluated_values").
			Field("build_graph_metrics_id"). // New field added
			Unique(),
	}
}
