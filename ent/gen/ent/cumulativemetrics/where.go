// Code generated by ent, DO NOT EDIT.

package cumulativemetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldLTE(FieldID, id))
}

// NumAnalyses applies equality check predicate on the "num_analyses" field. It's identical to NumAnalysesEQ.
func NumAnalyses(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldEQ(FieldNumAnalyses, v))
}

// NumBuilds applies equality check predicate on the "num_builds" field. It's identical to NumBuildsEQ.
func NumBuilds(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldEQ(FieldNumBuilds, v))
}

// NumAnalysesEQ applies the EQ predicate on the "num_analyses" field.
func NumAnalysesEQ(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldEQ(FieldNumAnalyses, v))
}

// NumAnalysesNEQ applies the NEQ predicate on the "num_analyses" field.
func NumAnalysesNEQ(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNEQ(FieldNumAnalyses, v))
}

// NumAnalysesIn applies the In predicate on the "num_analyses" field.
func NumAnalysesIn(vs ...int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldIn(FieldNumAnalyses, vs...))
}

// NumAnalysesNotIn applies the NotIn predicate on the "num_analyses" field.
func NumAnalysesNotIn(vs ...int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNotIn(FieldNumAnalyses, vs...))
}

// NumAnalysesGT applies the GT predicate on the "num_analyses" field.
func NumAnalysesGT(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldGT(FieldNumAnalyses, v))
}

// NumAnalysesGTE applies the GTE predicate on the "num_analyses" field.
func NumAnalysesGTE(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldGTE(FieldNumAnalyses, v))
}

// NumAnalysesLT applies the LT predicate on the "num_analyses" field.
func NumAnalysesLT(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldLT(FieldNumAnalyses, v))
}

// NumAnalysesLTE applies the LTE predicate on the "num_analyses" field.
func NumAnalysesLTE(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldLTE(FieldNumAnalyses, v))
}

// NumAnalysesIsNil applies the IsNil predicate on the "num_analyses" field.
func NumAnalysesIsNil() predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldIsNull(FieldNumAnalyses))
}

// NumAnalysesNotNil applies the NotNil predicate on the "num_analyses" field.
func NumAnalysesNotNil() predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNotNull(FieldNumAnalyses))
}

// NumBuildsEQ applies the EQ predicate on the "num_builds" field.
func NumBuildsEQ(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldEQ(FieldNumBuilds, v))
}

// NumBuildsNEQ applies the NEQ predicate on the "num_builds" field.
func NumBuildsNEQ(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNEQ(FieldNumBuilds, v))
}

// NumBuildsIn applies the In predicate on the "num_builds" field.
func NumBuildsIn(vs ...int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldIn(FieldNumBuilds, vs...))
}

// NumBuildsNotIn applies the NotIn predicate on the "num_builds" field.
func NumBuildsNotIn(vs ...int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNotIn(FieldNumBuilds, vs...))
}

// NumBuildsGT applies the GT predicate on the "num_builds" field.
func NumBuildsGT(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldGT(FieldNumBuilds, v))
}

// NumBuildsGTE applies the GTE predicate on the "num_builds" field.
func NumBuildsGTE(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldGTE(FieldNumBuilds, v))
}

// NumBuildsLT applies the LT predicate on the "num_builds" field.
func NumBuildsLT(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldLT(FieldNumBuilds, v))
}

// NumBuildsLTE applies the LTE predicate on the "num_builds" field.
func NumBuildsLTE(v int32) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldLTE(FieldNumBuilds, v))
}

// NumBuildsIsNil applies the IsNil predicate on the "num_builds" field.
func NumBuildsIsNil() predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldIsNull(FieldNumBuilds))
}

// NumBuildsNotNil applies the NotNil predicate on the "num_builds" field.
func NumBuildsNotNil() predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.FieldNotNull(FieldNumBuilds))
}

// HasMetrics applies the HasEdge predicate on the "metrics" edge.
func HasMetrics() predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, MetricsTable, MetricsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMetricsWith applies the HasEdge predicate on the "metrics" edge with a given conditions (other predicates).
func HasMetricsWith(preds ...predicate.Metrics) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(func(s *sql.Selector) {
		step := newMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.CumulativeMetrics) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.CumulativeMetrics) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.CumulativeMetrics) predicate.CumulativeMetrics {
	return predicate.CumulativeMetrics(sql.NotPredicates(p))
}
