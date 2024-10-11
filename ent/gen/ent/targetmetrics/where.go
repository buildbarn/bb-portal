// Code generated by ent, DO NOT EDIT.

package targetmetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLTE(FieldID, id))
}

// TargetsLoaded applies equality check predicate on the "targets_loaded" field. It's identical to TargetsLoadedEQ.
func TargetsLoaded(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldTargetsLoaded, v))
}

// TargetsConfigured applies equality check predicate on the "targets_configured" field. It's identical to TargetsConfiguredEQ.
func TargetsConfigured(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldTargetsConfigured, v))
}

// TargetsConfiguredNotIncludingAspects applies equality check predicate on the "targets_configured_not_including_aspects" field. It's identical to TargetsConfiguredNotIncludingAspectsEQ.
func TargetsConfiguredNotIncludingAspects(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsLoadedEQ applies the EQ predicate on the "targets_loaded" field.
func TargetsLoadedEQ(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldTargetsLoaded, v))
}

// TargetsLoadedNEQ applies the NEQ predicate on the "targets_loaded" field.
func TargetsLoadedNEQ(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNEQ(FieldTargetsLoaded, v))
}

// TargetsLoadedIn applies the In predicate on the "targets_loaded" field.
func TargetsLoadedIn(vs ...int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIn(FieldTargetsLoaded, vs...))
}

// TargetsLoadedNotIn applies the NotIn predicate on the "targets_loaded" field.
func TargetsLoadedNotIn(vs ...int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotIn(FieldTargetsLoaded, vs...))
}

// TargetsLoadedGT applies the GT predicate on the "targets_loaded" field.
func TargetsLoadedGT(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGT(FieldTargetsLoaded, v))
}

// TargetsLoadedGTE applies the GTE predicate on the "targets_loaded" field.
func TargetsLoadedGTE(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGTE(FieldTargetsLoaded, v))
}

// TargetsLoadedLT applies the LT predicate on the "targets_loaded" field.
func TargetsLoadedLT(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLT(FieldTargetsLoaded, v))
}

// TargetsLoadedLTE applies the LTE predicate on the "targets_loaded" field.
func TargetsLoadedLTE(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLTE(FieldTargetsLoaded, v))
}

// TargetsLoadedIsNil applies the IsNil predicate on the "targets_loaded" field.
func TargetsLoadedIsNil() predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIsNull(FieldTargetsLoaded))
}

// TargetsLoadedNotNil applies the NotNil predicate on the "targets_loaded" field.
func TargetsLoadedNotNil() predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotNull(FieldTargetsLoaded))
}

// TargetsConfiguredEQ applies the EQ predicate on the "targets_configured" field.
func TargetsConfiguredEQ(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldTargetsConfigured, v))
}

// TargetsConfiguredNEQ applies the NEQ predicate on the "targets_configured" field.
func TargetsConfiguredNEQ(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNEQ(FieldTargetsConfigured, v))
}

// TargetsConfiguredIn applies the In predicate on the "targets_configured" field.
func TargetsConfiguredIn(vs ...int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIn(FieldTargetsConfigured, vs...))
}

// TargetsConfiguredNotIn applies the NotIn predicate on the "targets_configured" field.
func TargetsConfiguredNotIn(vs ...int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotIn(FieldTargetsConfigured, vs...))
}

// TargetsConfiguredGT applies the GT predicate on the "targets_configured" field.
func TargetsConfiguredGT(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGT(FieldTargetsConfigured, v))
}

// TargetsConfiguredGTE applies the GTE predicate on the "targets_configured" field.
func TargetsConfiguredGTE(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGTE(FieldTargetsConfigured, v))
}

// TargetsConfiguredLT applies the LT predicate on the "targets_configured" field.
func TargetsConfiguredLT(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLT(FieldTargetsConfigured, v))
}

// TargetsConfiguredLTE applies the LTE predicate on the "targets_configured" field.
func TargetsConfiguredLTE(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLTE(FieldTargetsConfigured, v))
}

// TargetsConfiguredIsNil applies the IsNil predicate on the "targets_configured" field.
func TargetsConfiguredIsNil() predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIsNull(FieldTargetsConfigured))
}

// TargetsConfiguredNotNil applies the NotNil predicate on the "targets_configured" field.
func TargetsConfiguredNotNil() predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotNull(FieldTargetsConfigured))
}

// TargetsConfiguredNotIncludingAspectsEQ applies the EQ predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsEQ(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldEQ(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsConfiguredNotIncludingAspectsNEQ applies the NEQ predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsNEQ(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNEQ(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsConfiguredNotIncludingAspectsIn applies the In predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsIn(vs ...int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIn(FieldTargetsConfiguredNotIncludingAspects, vs...))
}

// TargetsConfiguredNotIncludingAspectsNotIn applies the NotIn predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsNotIn(vs ...int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotIn(FieldTargetsConfiguredNotIncludingAspects, vs...))
}

// TargetsConfiguredNotIncludingAspectsGT applies the GT predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsGT(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGT(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsConfiguredNotIncludingAspectsGTE applies the GTE predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsGTE(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldGTE(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsConfiguredNotIncludingAspectsLT applies the LT predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsLT(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLT(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsConfiguredNotIncludingAspectsLTE applies the LTE predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsLTE(v int64) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldLTE(FieldTargetsConfiguredNotIncludingAspects, v))
}

// TargetsConfiguredNotIncludingAspectsIsNil applies the IsNil predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsIsNil() predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldIsNull(FieldTargetsConfiguredNotIncludingAspects))
}

// TargetsConfiguredNotIncludingAspectsNotNil applies the NotNil predicate on the "targets_configured_not_including_aspects" field.
func TargetsConfiguredNotIncludingAspectsNotNil() predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.FieldNotNull(FieldTargetsConfiguredNotIncludingAspects))
}

// HasMetrics applies the HasEdge predicate on the "metrics" edge.
func HasMetrics() predicate.TargetMetrics {
	return predicate.TargetMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, MetricsTable, MetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMetricsWith applies the HasEdge predicate on the "metrics" edge with a given conditions (other predicates).
func HasMetricsWith(preds ...predicate.Metrics) predicate.TargetMetrics {
	return predicate.TargetMetrics(func(s *sql.Selector) {
		step := newMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TargetMetrics) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TargetMetrics) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TargetMetrics) predicate.TargetMetrics {
	return predicate.TargetMetrics(sql.NotPredicates(p))
}
