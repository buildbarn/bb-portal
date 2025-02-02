// Code generated by ent, DO NOT EDIT.

package metrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Metrics {
	return predicate.Metrics(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Metrics {
	return predicate.Metrics(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Metrics {
	return predicate.Metrics(sql.FieldLTE(FieldID, id))
}

// HasBazelInvocation applies the HasEdge predicate on the "bazel_invocation" edge.
func HasBazelInvocation() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, BazelInvocationTable, BazelInvocationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBazelInvocationWith applies the HasEdge predicate on the "bazel_invocation" edge with a given conditions (other predicates).
func HasBazelInvocationWith(preds ...predicate.BazelInvocation) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newBazelInvocationStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasActionSummary applies the HasEdge predicate on the "action_summary" edge.
func HasActionSummary() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, ActionSummaryTable, ActionSummaryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasActionSummaryWith applies the HasEdge predicate on the "action_summary" edge with a given conditions (other predicates).
func HasActionSummaryWith(preds ...predicate.ActionSummary) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newActionSummaryStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMemoryMetrics applies the HasEdge predicate on the "memory_metrics" edge.
func HasMemoryMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, MemoryMetricsTable, MemoryMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMemoryMetricsWith applies the HasEdge predicate on the "memory_metrics" edge with a given conditions (other predicates).
func HasMemoryMetricsWith(preds ...predicate.MemoryMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newMemoryMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTargetMetrics applies the HasEdge predicate on the "target_metrics" edge.
func HasTargetMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, TargetMetricsTable, TargetMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTargetMetricsWith applies the HasEdge predicate on the "target_metrics" edge with a given conditions (other predicates).
func HasTargetMetricsWith(preds ...predicate.TargetMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newTargetMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPackageMetrics applies the HasEdge predicate on the "package_metrics" edge.
func HasPackageMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, PackageMetricsTable, PackageMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPackageMetricsWith applies the HasEdge predicate on the "package_metrics" edge with a given conditions (other predicates).
func HasPackageMetricsWith(preds ...predicate.PackageMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newPackageMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTimingMetrics applies the HasEdge predicate on the "timing_metrics" edge.
func HasTimingMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, TimingMetricsTable, TimingMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTimingMetricsWith applies the HasEdge predicate on the "timing_metrics" edge with a given conditions (other predicates).
func HasTimingMetricsWith(preds ...predicate.TimingMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newTimingMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCumulativeMetrics applies the HasEdge predicate on the "cumulative_metrics" edge.
func HasCumulativeMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, CumulativeMetricsTable, CumulativeMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCumulativeMetricsWith applies the HasEdge predicate on the "cumulative_metrics" edge with a given conditions (other predicates).
func HasCumulativeMetricsWith(preds ...predicate.CumulativeMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newCumulativeMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasArtifactMetrics applies the HasEdge predicate on the "artifact_metrics" edge.
func HasArtifactMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, ArtifactMetricsTable, ArtifactMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasArtifactMetricsWith applies the HasEdge predicate on the "artifact_metrics" edge with a given conditions (other predicates).
func HasArtifactMetricsWith(preds ...predicate.ArtifactMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newArtifactMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasNetworkMetrics applies the HasEdge predicate on the "network_metrics" edge.
func HasNetworkMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, NetworkMetricsTable, NetworkMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasNetworkMetricsWith applies the HasEdge predicate on the "network_metrics" edge with a given conditions (other predicates).
func HasNetworkMetricsWith(preds ...predicate.NetworkMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newNetworkMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDynamicExecutionMetrics applies the HasEdge predicate on the "dynamic_execution_metrics" edge.
func HasDynamicExecutionMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, DynamicExecutionMetricsTable, DynamicExecutionMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDynamicExecutionMetricsWith applies the HasEdge predicate on the "dynamic_execution_metrics" edge with a given conditions (other predicates).
func HasDynamicExecutionMetricsWith(preds ...predicate.DynamicExecutionMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newDynamicExecutionMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBuildGraphMetrics applies the HasEdge predicate on the "build_graph_metrics" edge.
func HasBuildGraphMetrics() predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, BuildGraphMetricsTable, BuildGraphMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBuildGraphMetricsWith applies the HasEdge predicate on the "build_graph_metrics" edge with a given conditions (other predicates).
func HasBuildGraphMetricsWith(preds ...predicate.BuildGraphMetrics) predicate.Metrics {
	return predicate.Metrics(func(s *sql.Selector) {
		step := newBuildGraphMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Metrics) predicate.Metrics {
	return predicate.Metrics(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Metrics) predicate.Metrics {
	return predicate.Metrics(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Metrics) predicate.Metrics {
	return predicate.Metrics(sql.NotPredicates(p))
}
