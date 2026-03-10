package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BuildGraphMetrics holds the schema definition for the BuildGraphMetrics entity.
type BuildGraphMetrics struct {
	ent.Schema
}

// Fields of the BuildGraphMetrics.
func (BuildGraphMetrics) Fields() []ent.Field {
	return []ent.Field{
		// Action Value Lookup Count.
		// How many configured targets/aspects were in this build, including any
		// that were analyzed on a prior build and are still valid. May not be
		// populated if analysis phase was fully cached. Note: for historical
		// reasons this includes input/output files and other configured targets
		// that do not actually have associated actions.
		field.Int32("action_lookup_value_count").Optional(),

		// Action Value Lookup Count Not Including Aspects.
		// How many configured targets alone were in this build: always at most
		// action_lookup_value_count. Useful mainly for historical comparisons to
		// TargetMetrics.targets_configured, which used to not count aspects. This
		// also includes configured targets that do not have associated actions.
		field.Int32("action_lookup_value_count_not_including_aspects").Optional(),

		// Action Count.
		// How many actions belonged to the configured targets/aspects above. It may
		// not be necessary to execute all of these actions to build the requested
		// targets. May not be populated if analysis phase was fully cached.
		field.Int32("action_count").Optional(),

		// Action Count Not Including Aspects.
		// How many configured targets alone were in this build: always at most
		// action_lookup_value_count. Useful mainly for historical comparisons to
		// TargetMetrics.targets_configured, which used to not count aspects. This
		// also includes configured targets that do not have associated actions.
		field.Int32("action_count_not_including_aspects").Optional(),

		// Input File Configured Target Count.
		// How many "input file" configured targets there were: one per source file.
		// Should agree with artifact_metrics.source_artifacts_read.count above,
		field.Int32("input_file_configured_target_count").Optional(),

		// Output File Configured Target Count.
		// How many "output file" configured targets there were: output files that
		// are targets (not implicit outputs).
		field.Int32("output_file_configured_target_count").Optional(),

		// Other Configured Target Count.
		// How many "other" configured targets there were (like alias,
		// package_group, and other non-rule non-file configured targets).
		field.Int32("other_configured_target_count").Optional(),

		// Output Artifact Count.
		// How many artifacts are outputs of the above actions. May not be populated
		// if analysis phase was fully cached.
		field.Int32("output_artifact_count").Optional(),

		// Post Invocation Skyframe Node Count.
		// How many Skyframe nodes there are in memory at the end of the build. This
		// may underestimate the number of nodes when running with memory-saving
		// settings or with Skybuild, and may overestimate if there are nodes from
		// prior evaluations still in the cache.
		field.Int32("post_invocation_skyframe_node_count").Optional(),
	}
}

// Edges of the BuildGraphMetrics.
func (BuildGraphMetrics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("build_graph_metrics").
			Unique(),
	}
}

// Indexes of the BuildGraphMetrics.
func (BuildGraphMetrics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("metrics"),
	}
}

// Mixin of the BuildGraphMetrics.
func (BuildGraphMetrics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
