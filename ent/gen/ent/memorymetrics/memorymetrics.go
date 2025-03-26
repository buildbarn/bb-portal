// Code generated by ent, DO NOT EDIT.

package memorymetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the memorymetrics type in the database.
	Label = "memory_metrics"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldPeakPostGcHeapSize holds the string denoting the peak_post_gc_heap_size field in the database.
	FieldPeakPostGcHeapSize = "peak_post_gc_heap_size"
	// FieldUsedHeapSizePostBuild holds the string denoting the used_heap_size_post_build field in the database.
	FieldUsedHeapSizePostBuild = "used_heap_size_post_build"
	// FieldPeakPostGcTenuredSpaceHeapSize holds the string denoting the peak_post_gc_tenured_space_heap_size field in the database.
	FieldPeakPostGcTenuredSpaceHeapSize = "peak_post_gc_tenured_space_heap_size"
	// FieldMetricsID holds the string denoting the metrics_id field in the database.
	FieldMetricsID = "metrics_id"
	// EdgeMetrics holds the string denoting the metrics edge name in mutations.
	EdgeMetrics = "metrics"
	// EdgeGarbageMetrics holds the string denoting the garbage_metrics edge name in mutations.
	EdgeGarbageMetrics = "garbage_metrics"
	// Table holds the table name of the memorymetrics in the database.
	Table = "memory_metrics"
	// MetricsTable is the table that holds the metrics relation/edge.
	MetricsTable = "memory_metrics"
	// MetricsInverseTable is the table name for the Metrics entity.
	// It exists in this package in order to avoid circular dependency with the "metrics" package.
	MetricsInverseTable = "metrics"
	// MetricsColumn is the table column denoting the metrics relation/edge.
	MetricsColumn = "metrics_id"
	// GarbageMetricsTable is the table that holds the garbage_metrics relation/edge.
	GarbageMetricsTable = "garbage_metrics"
	// GarbageMetricsInverseTable is the table name for the GarbageMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "garbagemetrics" package.
	GarbageMetricsInverseTable = "garbage_metrics"
	// GarbageMetricsColumn is the table column denoting the garbage_metrics relation/edge.
	GarbageMetricsColumn = "memory_metrics_id"
)

// Columns holds all SQL columns for memorymetrics fields.
var Columns = []string{
	FieldID,
	FieldPeakPostGcHeapSize,
	FieldUsedHeapSizePostBuild,
	FieldPeakPostGcTenuredSpaceHeapSize,
	FieldMetricsID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the MemoryMetrics queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByPeakPostGcHeapSize orders the results by the peak_post_gc_heap_size field.
func ByPeakPostGcHeapSize(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPeakPostGcHeapSize, opts...).ToFunc()
}

// ByUsedHeapSizePostBuild orders the results by the used_heap_size_post_build field.
func ByUsedHeapSizePostBuild(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUsedHeapSizePostBuild, opts...).ToFunc()
}

// ByPeakPostGcTenuredSpaceHeapSize orders the results by the peak_post_gc_tenured_space_heap_size field.
func ByPeakPostGcTenuredSpaceHeapSize(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPeakPostGcTenuredSpaceHeapSize, opts...).ToFunc()
}

// ByMetricsID orders the results by the metrics_id field.
func ByMetricsID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMetricsID, opts...).ToFunc()
}

// ByMetricsField orders the results by metrics field.
func ByMetricsField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMetricsStep(), sql.OrderByField(field, opts...))
	}
}

// ByGarbageMetricsCount orders the results by garbage_metrics count.
func ByGarbageMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newGarbageMetricsStep(), opts...)
	}
}

// ByGarbageMetrics orders the results by garbage_metrics terms.
func ByGarbageMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newGarbageMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, MetricsTable, MetricsColumn),
	)
}
func newGarbageMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(GarbageMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, GarbageMetricsTable, GarbageMetricsColumn),
	)
}
