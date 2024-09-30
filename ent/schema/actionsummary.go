package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ActionSummary holds the schema definition for the ActionSummary entity.
type ActionSummary struct {
	ent.Schema
}

// Fields of the ActionSummary.
func (ActionSummary) Fields() []ent.Field {
	return []ent.Field{
		// The total number of actions created and registered during the build,
		// including both aspects and configured targets. This metric includes
		// unused actions that were constructed but not executed during this build.
		// It does not include actions that were created on prior builds that are
		// still valid, even if those actions had to be re-executed on this build.
		// For the total number of actions that would be created if this invocation
		// were "clean", see BuildGraphMetrics below.
		field.Int64("actions_created").Optional(),

		// The total number of actions created this build just by configured
		// targets. Used mainly to allow consumers of actions_created, which used to
		// not include aspects' actions, to normalize across the Blaze release that
		// switched actions_created to include all created actions.
		field.Int64("actions_created_not_including_aspects").Optional(),

		// The total number of actions executed during the build. This includes any
		// remote cache hits, but excludes local action cache hits.
		field.Int64("actions_executed").Optional(),

		// Deprecated. The total number of remote cache hits.
		field.Int64("remote_cache_hits").Optional(),
	}
}

// Edges of the ActionSummary.
func (ActionSummary) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metrics", Metrics.Type).
			Ref("action_summary").
			Unique(),

		// Contains the top N actions by number of actions executed.
		edge.To("action_data", ActionData.Type),

		// Count of which Runner types were executed which actions.
		edge.To("runner_count", RunnerCount.Type),

		// Information about the action cache behavior during a single invocation.
		edge.To("action_cache_statistics", ActionCacheStatistics.Type),
	}
}
