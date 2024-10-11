package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ExectionInfo holds the schema definition for the ExectionInfo entity.
type ExectionInfo struct {
	ent.Schema
}

// Fields of the ExectionInfo.
func (ExectionInfo) Fields() []ent.Field {
	return []ent.Field{
		// Deprecated, use TargetComplete.test_timeout instead.
		field.Int32("timeout_seconds").Optional(),

		// Name of the strategy to execute this test action (e.g., "local", "remote").
		field.String("strategy").Optional(),

		// True, if the reported attempt was a cache hit in a remote cache.
		field.Bool("cached_remotely").Optional(),

		// The exit code of the test action.
		field.Int32("exit_code").Optional(),

		// Hostname.
		// The hostname of the machine where the test action was executed (in case
		// of remote execution), if known.
		field.String("hostname").Optional(),
	}
}

// Edges of ExectionInfo.
func (ExectionInfo) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the test result
		edge.From("test_result", TestResultBES.Type).
			Ref("execution_info").
			Unique(),

		// Represents a hierarchical timing breakdown of an activity.
		// The top level time should be the total time of the activity.
		// Invariant: `time` >= sum of `time`s of all direct children.
		edge.To("timing_breakdown", TimingBreakdown.Type).Unique(),

		// resource usage info
		edge.To("resource_usage", ResourceUsage.Type),
	}
}
