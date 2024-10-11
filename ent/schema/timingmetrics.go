package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TimingMetrics holds the schema definition for the TimingMetrics entity.
type TimingMetrics struct {
	ent.Schema
}

// Fields of the TimingMetrics.
func (TimingMetrics) Fields() []ent.Field {
	return []ent.Field{
		// The CPU time in milliseconds consumed during this build.
		// For Skymeld, it's possible that
		// analysis_phase_time_in_ms + execution_phase_time_in_ms >= wall_time_in_ms
		field.Int64("cpu_time_in_ms").Optional(),

		// The elapsed wall time in milliseconds during this build.
		field.Int64("wall_time_in_ms").Optional(),

		// The elapsed wall time in milliseconds during the analysis phase.
		// When analysis and execution phases are interleaved, this measures the
		// elapsed time from the first analysis work to the last.
		field.Int64("analysis_phase_time_in_ms").Optional(),

		// The elapsed wall time in milliseconds during the execution phase.
		// When analysis and execution phases are interleaved, this measures the
		// elapsed time from the first action execution (excluding workspace status
		// actions) to the last.
		field.Int64("execution_phase_time_in_ms").Optional(),

		// The elapsed wall time in milliseconds until the first action execution.
		// started (excluding workspace status actions).
		field.Int64("actions_execution_start_in_ms").Optional(),
	}
}

// Edges of the TimingMetrics.
func (TimingMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the metrics object.
		edge.From("metrics", Metrics.Type).
			Ref("timing_metrics").
			Unique(),
	}
}
