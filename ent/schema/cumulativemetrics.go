package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CumulativeMetrics holds the schema definition for the CumulativeMetrics entity.
type CumulativeMetrics struct {
	ent.Schema
}

// Fields of the CumulativeMetrics.
func (CumulativeMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Number Of Analyses.
		// One-indexed number of "analyses" the server has run, including the
		// current one. Will be incremented for every build/test/cquery/etc. command
		// that reaches the analysis phase.
		field.Int32("num_analyses").Optional(),

		// Number of Builds.
		// One-indexed number of "builds" the server has run, including the current
		// one. Will be incremented for every build/test/run/etc. command that
		// reaches the execution phase.
		field.Int32("num_builds").Optional(),
	}
}

// Edges of the TimingMetrics.
func (CumulativeMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metircs object.
		edge.From("metrics", Metrics.Type).Ref("cumulative_metrics"),
	}
}
