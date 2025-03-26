package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MissDetail holds the schema definition for the MissDetail entity.
type MissDetail struct {
	ent.Schema
}

// Fields of the MissDetail.
func (MissDetail) Fields() []ent.Field {
	return []ent.Field{
		// Reasons for not finding an action in the cache.
		field.Enum("reason").
			Values("DIFFERENT_ACTION_KEY",
				"DIFFERENT_DEPS",
				"DIFFERENT_ENVIRONMENT",
				"DIFFERENT_FILES",
				"CORRUPTED_CACHE_ENTRY",
				"NOT_CACHED",
				"UNCONDITIONAL_EXECUTION",
				"UNKNOWN").
			Default("UNKNOWN").Optional(),

		// Counter for this type.
		field.Int32("count").Optional(),

		// foreign key to the action cache statistics object.
		field.Int("action_cache_statistics_id").Optional(),
	}
}

// Edges of the MissDetail.
func (MissDetail) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the action cache statistics object.
		edge.From("action_cache_statistics", ActionCacheStatistics.Type).
			Ref("miss_details").
			Unique().
			Field("action_cache_statistics_id"), // New field added
	}
}

// Indexes of the MissDetail.
func (MissDetail) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("action_cache_statistics_id"),
	}
}
