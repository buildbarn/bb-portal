package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.Uint64("size_in_bytes").
			Annotations(entgql.Type("Uint64")).
			Optional(),

		// Time it took to save the action cache to disk.
		field.Uint64("save_time_in_ms").
			Annotations(entgql.Type("Uint64")).
			Optional(),

		// Time it took to load the action cache from disk. Reported as 0 if the
		// action cache has not been loaded in this invocation.
		field.Uint64("load_time_in_ms").
			Annotations(entgql.Type("Uint64")).
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Time spent waiting on the cache check semaphore.
		field.Uint64("cache_check_semaphore_wait_time_in_ms").
			Annotations(entgql.Type("Uint64")).
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

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
		edge.To("miss_details", MissDetail.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

// Indexes of the ActionCacheStatistics.
func (ActionCacheStatistics) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("action_summary"),
	}
}

// Mixin of the ActionCacheStatistics.
func (ActionCacheStatistics) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
