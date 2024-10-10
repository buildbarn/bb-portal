package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ActionCacheStatistics holds the schema definition for the ActionCacheStatistics entity.
type ActionCacheStatistics struct {
	ent.Schema
}

// Fields of the ActionCacheStatistics.
func (ActionCacheStatistics) Fields() []ent.Field {
	return []ent.Field{
		// Size of the action cache in bytes.
		// This is computed by the code that persists the action cache to disk and
		// represents the size of the written files, which has no direct relation to
		// the number of entries in the cache.
		field.Uint64("size_in_bytes").Optional(),

		// Time it took to save the action cache to disk.
		field.Uint64("save_time_in_ms").Optional(),

		// Time it took to load the action cache from disk. Reported as 0 if the
		// action cache has not been loaded in this invocation.
		field.Int64("load_time_in_ms").Optional(),

		// Cache counters.
		field.Int32("hits").Optional(),
		field.Int32("misses").Optional(),
	}
}

// Edges of the ActionCacheStatistics.
func (ActionCacheStatistics) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the associated action summary.
		edge.From("action_summary", ActionSummary.Type).
			Ref("action_cache_statistics").
			Unique(),

		// Breakdown of the cache misses based on the reasons behind them.
		edge.To("miss_details", MissDetail.Type),
	}
}
