// Code generated by ent, DO NOT EDIT.

package memorymetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLTE(FieldID, id))
}

// PeakPostGcHeapSize applies equality check predicate on the "peak_post_gc_heap_size" field. It's identical to PeakPostGcHeapSizeEQ.
func PeakPostGcHeapSize(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldPeakPostGcHeapSize, v))
}

// UsedHeapSizePostBuild applies equality check predicate on the "used_heap_size_post_build" field. It's identical to UsedHeapSizePostBuildEQ.
func UsedHeapSizePostBuild(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldUsedHeapSizePostBuild, v))
}

// PeakPostGcTenuredSpaceHeapSize applies equality check predicate on the "peak_post_gc_tenured_space_heap_size" field. It's identical to PeakPostGcTenuredSpaceHeapSizeEQ.
func PeakPostGcTenuredSpaceHeapSize(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// MetricsID applies equality check predicate on the "metrics_id" field. It's identical to MetricsIDEQ.
func MetricsID(v int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldMetricsID, v))
}

// PeakPostGcHeapSizeEQ applies the EQ predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeEQ(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldPeakPostGcHeapSize, v))
}

// PeakPostGcHeapSizeNEQ applies the NEQ predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeNEQ(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNEQ(FieldPeakPostGcHeapSize, v))
}

// PeakPostGcHeapSizeIn applies the In predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeIn(vs ...int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIn(FieldPeakPostGcHeapSize, vs...))
}

// PeakPostGcHeapSizeNotIn applies the NotIn predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeNotIn(vs ...int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotIn(FieldPeakPostGcHeapSize, vs...))
}

// PeakPostGcHeapSizeGT applies the GT predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeGT(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGT(FieldPeakPostGcHeapSize, v))
}

// PeakPostGcHeapSizeGTE applies the GTE predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeGTE(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGTE(FieldPeakPostGcHeapSize, v))
}

// PeakPostGcHeapSizeLT applies the LT predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeLT(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLT(FieldPeakPostGcHeapSize, v))
}

// PeakPostGcHeapSizeLTE applies the LTE predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeLTE(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLTE(FieldPeakPostGcHeapSize, v))
}

// PeakPostGcHeapSizeIsNil applies the IsNil predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeIsNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIsNull(FieldPeakPostGcHeapSize))
}

// PeakPostGcHeapSizeNotNil applies the NotNil predicate on the "peak_post_gc_heap_size" field.
func PeakPostGcHeapSizeNotNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotNull(FieldPeakPostGcHeapSize))
}

// UsedHeapSizePostBuildEQ applies the EQ predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildEQ(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldUsedHeapSizePostBuild, v))
}

// UsedHeapSizePostBuildNEQ applies the NEQ predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildNEQ(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNEQ(FieldUsedHeapSizePostBuild, v))
}

// UsedHeapSizePostBuildIn applies the In predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildIn(vs ...int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIn(FieldUsedHeapSizePostBuild, vs...))
}

// UsedHeapSizePostBuildNotIn applies the NotIn predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildNotIn(vs ...int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotIn(FieldUsedHeapSizePostBuild, vs...))
}

// UsedHeapSizePostBuildGT applies the GT predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildGT(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGT(FieldUsedHeapSizePostBuild, v))
}

// UsedHeapSizePostBuildGTE applies the GTE predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildGTE(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGTE(FieldUsedHeapSizePostBuild, v))
}

// UsedHeapSizePostBuildLT applies the LT predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildLT(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLT(FieldUsedHeapSizePostBuild, v))
}

// UsedHeapSizePostBuildLTE applies the LTE predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildLTE(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLTE(FieldUsedHeapSizePostBuild, v))
}

// UsedHeapSizePostBuildIsNil applies the IsNil predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildIsNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIsNull(FieldUsedHeapSizePostBuild))
}

// UsedHeapSizePostBuildNotNil applies the NotNil predicate on the "used_heap_size_post_build" field.
func UsedHeapSizePostBuildNotNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotNull(FieldUsedHeapSizePostBuild))
}

// PeakPostGcTenuredSpaceHeapSizeEQ applies the EQ predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeEQ(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// PeakPostGcTenuredSpaceHeapSizeNEQ applies the NEQ predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeNEQ(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNEQ(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// PeakPostGcTenuredSpaceHeapSizeIn applies the In predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeIn(vs ...int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIn(FieldPeakPostGcTenuredSpaceHeapSize, vs...))
}

// PeakPostGcTenuredSpaceHeapSizeNotIn applies the NotIn predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeNotIn(vs ...int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotIn(FieldPeakPostGcTenuredSpaceHeapSize, vs...))
}

// PeakPostGcTenuredSpaceHeapSizeGT applies the GT predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeGT(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGT(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// PeakPostGcTenuredSpaceHeapSizeGTE applies the GTE predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeGTE(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldGTE(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// PeakPostGcTenuredSpaceHeapSizeLT applies the LT predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeLT(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLT(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// PeakPostGcTenuredSpaceHeapSizeLTE applies the LTE predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeLTE(v int64) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldLTE(FieldPeakPostGcTenuredSpaceHeapSize, v))
}

// PeakPostGcTenuredSpaceHeapSizeIsNil applies the IsNil predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeIsNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIsNull(FieldPeakPostGcTenuredSpaceHeapSize))
}

// PeakPostGcTenuredSpaceHeapSizeNotNil applies the NotNil predicate on the "peak_post_gc_tenured_space_heap_size" field.
func PeakPostGcTenuredSpaceHeapSizeNotNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotNull(FieldPeakPostGcTenuredSpaceHeapSize))
}

// MetricsIDEQ applies the EQ predicate on the "metrics_id" field.
func MetricsIDEQ(v int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldEQ(FieldMetricsID, v))
}

// MetricsIDNEQ applies the NEQ predicate on the "metrics_id" field.
func MetricsIDNEQ(v int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNEQ(FieldMetricsID, v))
}

// MetricsIDIn applies the In predicate on the "metrics_id" field.
func MetricsIDIn(vs ...int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIn(FieldMetricsID, vs...))
}

// MetricsIDNotIn applies the NotIn predicate on the "metrics_id" field.
func MetricsIDNotIn(vs ...int) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotIn(FieldMetricsID, vs...))
}

// MetricsIDIsNil applies the IsNil predicate on the "metrics_id" field.
func MetricsIDIsNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldIsNull(FieldMetricsID))
}

// MetricsIDNotNil applies the NotNil predicate on the "metrics_id" field.
func MetricsIDNotNil() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.FieldNotNull(FieldMetricsID))
}

// HasMetrics applies the HasEdge predicate on the "metrics" edge.
func HasMetrics() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, MetricsTable, MetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMetricsWith applies the HasEdge predicate on the "metrics" edge with a given conditions (other predicates).
func HasMetricsWith(preds ...predicate.Metrics) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(func(s *sql.Selector) {
		step := newMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasGarbageMetrics applies the HasEdge predicate on the "garbage_metrics" edge.
func HasGarbageMetrics() predicate.MemoryMetrics {
	return predicate.MemoryMetrics(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, GarbageMetricsTable, GarbageMetricsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGarbageMetricsWith applies the HasEdge predicate on the "garbage_metrics" edge with a given conditions (other predicates).
func HasGarbageMetricsWith(preds ...predicate.GarbageMetrics) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(func(s *sql.Selector) {
		step := newGarbageMetricsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.MemoryMetrics) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.MemoryMetrics) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.MemoryMetrics) predicate.MemoryMetrics {
	return predicate.MemoryMetrics(sql.NotPredicates(p))
}
