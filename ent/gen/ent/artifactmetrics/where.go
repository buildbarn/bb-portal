// Code generated by ent, DO NOT EDIT.

package artifactmetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.FieldLTE(FieldID, id))
}

// HasMetrics applies the HasEdge predicate on the "metrics" edge.
func HasMetrics() predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, MetricsTable, MetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMetricsWith applies the HasEdge predicate on the "metrics" edge with a given conditions (other predicates).
func HasMetricsWith(preds ...predicate.Metrics) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := newMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSourceArtifactsRead applies the HasEdge predicate on the "source_artifacts_read" edge.
func HasSourceArtifactsRead() predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SourceArtifactsReadTable, SourceArtifactsReadColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSourceArtifactsReadWith applies the HasEdge predicate on the "source_artifacts_read" edge with a given conditions (other predicates).
func HasSourceArtifactsReadWith(preds ...predicate.FilesMetric) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := newSourceArtifactsReadStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOutputArtifactsSeen applies the HasEdge predicate on the "output_artifacts_seen" edge.
func HasOutputArtifactsSeen() predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OutputArtifactsSeenTable, OutputArtifactsSeenColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOutputArtifactsSeenWith applies the HasEdge predicate on the "output_artifacts_seen" edge with a given conditions (other predicates).
func HasOutputArtifactsSeenWith(preds ...predicate.FilesMetric) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := newOutputArtifactsSeenStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOutputArtifactsFromActionCache applies the HasEdge predicate on the "output_artifacts_from_action_cache" edge.
func HasOutputArtifactsFromActionCache() predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OutputArtifactsFromActionCacheTable, OutputArtifactsFromActionCacheColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOutputArtifactsFromActionCacheWith applies the HasEdge predicate on the "output_artifacts_from_action_cache" edge with a given conditions (other predicates).
func HasOutputArtifactsFromActionCacheWith(preds ...predicate.FilesMetric) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := newOutputArtifactsFromActionCacheStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTopLevelArtifacts applies the HasEdge predicate on the "top_level_artifacts" edge.
func HasTopLevelArtifacts() predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, TopLevelArtifactsTable, TopLevelArtifactsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTopLevelArtifactsWith applies the HasEdge predicate on the "top_level_artifacts" edge with a given conditions (other predicates).
func HasTopLevelArtifactsWith(preds ...predicate.FilesMetric) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(func(s *sql.Selector) {
		step := newTopLevelArtifactsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ArtifactMetrics) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ArtifactMetrics) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ArtifactMetrics) predicate.ArtifactMetrics {
	return predicate.ArtifactMetrics(sql.NotPredicates(p))
}
