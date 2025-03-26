// Code generated by ent, DO NOT EDIT.

package packageloadmetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldName, v))
}

// LoadDuration applies equality check predicate on the "load_duration" field. It's identical to LoadDurationEQ.
func LoadDuration(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldLoadDuration, v))
}

// NumTargets applies equality check predicate on the "num_targets" field. It's identical to NumTargetsEQ.
func NumTargets(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldNumTargets, v))
}

// ComputationSteps applies equality check predicate on the "computation_steps" field. It's identical to ComputationStepsEQ.
func ComputationSteps(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldComputationSteps, v))
}

// NumTransitiveLoads applies equality check predicate on the "num_transitive_loads" field. It's identical to NumTransitiveLoadsEQ.
func NumTransitiveLoads(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldNumTransitiveLoads, v))
}

// PackageOverhead applies equality check predicate on the "package_overhead" field. It's identical to PackageOverheadEQ.
func PackageOverhead(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldPackageOverhead, v))
}

// PackageMetricsID applies equality check predicate on the "package_metrics_id" field. It's identical to PackageMetricsIDEQ.
func PackageMetricsID(v int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldPackageMetricsID, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldHasSuffix(FieldName, v))
}

// NameIsNil applies the IsNil predicate on the "name" field.
func NameIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldName))
}

// NameNotNil applies the NotNil predicate on the "name" field.
func NameNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldName))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldContainsFold(FieldName, v))
}

// LoadDurationEQ applies the EQ predicate on the "load_duration" field.
func LoadDurationEQ(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldLoadDuration, v))
}

// LoadDurationNEQ applies the NEQ predicate on the "load_duration" field.
func LoadDurationNEQ(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldLoadDuration, v))
}

// LoadDurationIn applies the In predicate on the "load_duration" field.
func LoadDurationIn(vs ...int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldLoadDuration, vs...))
}

// LoadDurationNotIn applies the NotIn predicate on the "load_duration" field.
func LoadDurationNotIn(vs ...int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldLoadDuration, vs...))
}

// LoadDurationGT applies the GT predicate on the "load_duration" field.
func LoadDurationGT(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldLoadDuration, v))
}

// LoadDurationGTE applies the GTE predicate on the "load_duration" field.
func LoadDurationGTE(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldLoadDuration, v))
}

// LoadDurationLT applies the LT predicate on the "load_duration" field.
func LoadDurationLT(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldLoadDuration, v))
}

// LoadDurationLTE applies the LTE predicate on the "load_duration" field.
func LoadDurationLTE(v int64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldLoadDuration, v))
}

// LoadDurationIsNil applies the IsNil predicate on the "load_duration" field.
func LoadDurationIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldLoadDuration))
}

// LoadDurationNotNil applies the NotNil predicate on the "load_duration" field.
func LoadDurationNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldLoadDuration))
}

// NumTargetsEQ applies the EQ predicate on the "num_targets" field.
func NumTargetsEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldNumTargets, v))
}

// NumTargetsNEQ applies the NEQ predicate on the "num_targets" field.
func NumTargetsNEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldNumTargets, v))
}

// NumTargetsIn applies the In predicate on the "num_targets" field.
func NumTargetsIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldNumTargets, vs...))
}

// NumTargetsNotIn applies the NotIn predicate on the "num_targets" field.
func NumTargetsNotIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldNumTargets, vs...))
}

// NumTargetsGT applies the GT predicate on the "num_targets" field.
func NumTargetsGT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldNumTargets, v))
}

// NumTargetsGTE applies the GTE predicate on the "num_targets" field.
func NumTargetsGTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldNumTargets, v))
}

// NumTargetsLT applies the LT predicate on the "num_targets" field.
func NumTargetsLT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldNumTargets, v))
}

// NumTargetsLTE applies the LTE predicate on the "num_targets" field.
func NumTargetsLTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldNumTargets, v))
}

// NumTargetsIsNil applies the IsNil predicate on the "num_targets" field.
func NumTargetsIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldNumTargets))
}

// NumTargetsNotNil applies the NotNil predicate on the "num_targets" field.
func NumTargetsNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldNumTargets))
}

// ComputationStepsEQ applies the EQ predicate on the "computation_steps" field.
func ComputationStepsEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldComputationSteps, v))
}

// ComputationStepsNEQ applies the NEQ predicate on the "computation_steps" field.
func ComputationStepsNEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldComputationSteps, v))
}

// ComputationStepsIn applies the In predicate on the "computation_steps" field.
func ComputationStepsIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldComputationSteps, vs...))
}

// ComputationStepsNotIn applies the NotIn predicate on the "computation_steps" field.
func ComputationStepsNotIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldComputationSteps, vs...))
}

// ComputationStepsGT applies the GT predicate on the "computation_steps" field.
func ComputationStepsGT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldComputationSteps, v))
}

// ComputationStepsGTE applies the GTE predicate on the "computation_steps" field.
func ComputationStepsGTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldComputationSteps, v))
}

// ComputationStepsLT applies the LT predicate on the "computation_steps" field.
func ComputationStepsLT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldComputationSteps, v))
}

// ComputationStepsLTE applies the LTE predicate on the "computation_steps" field.
func ComputationStepsLTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldComputationSteps, v))
}

// ComputationStepsIsNil applies the IsNil predicate on the "computation_steps" field.
func ComputationStepsIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldComputationSteps))
}

// ComputationStepsNotNil applies the NotNil predicate on the "computation_steps" field.
func ComputationStepsNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldComputationSteps))
}

// NumTransitiveLoadsEQ applies the EQ predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldNumTransitiveLoads, v))
}

// NumTransitiveLoadsNEQ applies the NEQ predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsNEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldNumTransitiveLoads, v))
}

// NumTransitiveLoadsIn applies the In predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldNumTransitiveLoads, vs...))
}

// NumTransitiveLoadsNotIn applies the NotIn predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsNotIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldNumTransitiveLoads, vs...))
}

// NumTransitiveLoadsGT applies the GT predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsGT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldNumTransitiveLoads, v))
}

// NumTransitiveLoadsGTE applies the GTE predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsGTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldNumTransitiveLoads, v))
}

// NumTransitiveLoadsLT applies the LT predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsLT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldNumTransitiveLoads, v))
}

// NumTransitiveLoadsLTE applies the LTE predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsLTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldNumTransitiveLoads, v))
}

// NumTransitiveLoadsIsNil applies the IsNil predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldNumTransitiveLoads))
}

// NumTransitiveLoadsNotNil applies the NotNil predicate on the "num_transitive_loads" field.
func NumTransitiveLoadsNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldNumTransitiveLoads))
}

// PackageOverheadEQ applies the EQ predicate on the "package_overhead" field.
func PackageOverheadEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldPackageOverhead, v))
}

// PackageOverheadNEQ applies the NEQ predicate on the "package_overhead" field.
func PackageOverheadNEQ(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldPackageOverhead, v))
}

// PackageOverheadIn applies the In predicate on the "package_overhead" field.
func PackageOverheadIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldPackageOverhead, vs...))
}

// PackageOverheadNotIn applies the NotIn predicate on the "package_overhead" field.
func PackageOverheadNotIn(vs ...uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldPackageOverhead, vs...))
}

// PackageOverheadGT applies the GT predicate on the "package_overhead" field.
func PackageOverheadGT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGT(FieldPackageOverhead, v))
}

// PackageOverheadGTE applies the GTE predicate on the "package_overhead" field.
func PackageOverheadGTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldGTE(FieldPackageOverhead, v))
}

// PackageOverheadLT applies the LT predicate on the "package_overhead" field.
func PackageOverheadLT(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLT(FieldPackageOverhead, v))
}

// PackageOverheadLTE applies the LTE predicate on the "package_overhead" field.
func PackageOverheadLTE(v uint64) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldLTE(FieldPackageOverhead, v))
}

// PackageOverheadIsNil applies the IsNil predicate on the "package_overhead" field.
func PackageOverheadIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldPackageOverhead))
}

// PackageOverheadNotNil applies the NotNil predicate on the "package_overhead" field.
func PackageOverheadNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldPackageOverhead))
}

// PackageMetricsIDEQ applies the EQ predicate on the "package_metrics_id" field.
func PackageMetricsIDEQ(v int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldEQ(FieldPackageMetricsID, v))
}

// PackageMetricsIDNEQ applies the NEQ predicate on the "package_metrics_id" field.
func PackageMetricsIDNEQ(v int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNEQ(FieldPackageMetricsID, v))
}

// PackageMetricsIDIn applies the In predicate on the "package_metrics_id" field.
func PackageMetricsIDIn(vs ...int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIn(FieldPackageMetricsID, vs...))
}

// PackageMetricsIDNotIn applies the NotIn predicate on the "package_metrics_id" field.
func PackageMetricsIDNotIn(vs ...int) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotIn(FieldPackageMetricsID, vs...))
}

// PackageMetricsIDIsNil applies the IsNil predicate on the "package_metrics_id" field.
func PackageMetricsIDIsNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldIsNull(FieldPackageMetricsID))
}

// PackageMetricsIDNotNil applies the NotNil predicate on the "package_metrics_id" field.
func PackageMetricsIDNotNil() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.FieldNotNull(FieldPackageMetricsID))
}

// HasPackageMetrics applies the HasEdge predicate on the "package_metrics" edge.
func HasPackageMetrics() predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, PackageMetricsTable, PackageMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPackageMetricsWith applies the HasEdge predicate on the "package_metrics" edge with a given conditions (other predicates).
func HasPackageMetricsWith(preds ...predicate.PackageMetrics) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(func(s *sql.Selector) {
		step := newPackageMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PackageLoadMetrics) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PackageLoadMetrics) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PackageLoadMetrics) predicate.PackageLoadMetrics {
	return predicate.PackageLoadMetrics(sql.NotPredicates(p))
}
