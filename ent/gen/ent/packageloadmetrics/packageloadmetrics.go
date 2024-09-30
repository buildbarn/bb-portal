// Code generated by ent, DO NOT EDIT.

package packageloadmetrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the packageloadmetrics type in the database.
	Label = "package_load_metrics"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldLoadDuration holds the string denoting the load_duration field in the database.
	FieldLoadDuration = "load_duration"
	// FieldNumTargets holds the string denoting the num_targets field in the database.
	FieldNumTargets = "num_targets"
	// FieldComputationSteps holds the string denoting the computation_steps field in the database.
	FieldComputationSteps = "computation_steps"
	// FieldNumTransitiveLoads holds the string denoting the num_transitive_loads field in the database.
	FieldNumTransitiveLoads = "num_transitive_loads"
	// FieldPackageOverhead holds the string denoting the package_overhead field in the database.
	FieldPackageOverhead = "package_overhead"
	// EdgePackageMetrics holds the string denoting the package_metrics edge name in mutations.
	EdgePackageMetrics = "package_metrics"
	// Table holds the table name of the packageloadmetrics in the database.
	Table = "package_load_metrics"
	// PackageMetricsTable is the table that holds the package_metrics relation/edge. The primary key declared below.
	PackageMetricsTable = "package_metrics_package_load_metrics"
	// PackageMetricsInverseTable is the table name for the PackageMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "packagemetrics" package.
	PackageMetricsInverseTable = "package_metrics"
)

// Columns holds all SQL columns for packageloadmetrics fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldLoadDuration,
	FieldNumTargets,
	FieldComputationSteps,
	FieldNumTransitiveLoads,
	FieldPackageOverhead,
}

var (
	// PackageMetricsPrimaryKey and PackageMetricsColumn2 are the table columns denoting the
	// primary key for the package_metrics relation (M2M).
	PackageMetricsPrimaryKey = []string{"package_metrics_id", "package_load_metrics_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the PackageLoadMetrics queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByLoadDuration orders the results by the load_duration field.
func ByLoadDuration(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLoadDuration, opts...).ToFunc()
}

// ByNumTargets orders the results by the num_targets field.
func ByNumTargets(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNumTargets, opts...).ToFunc()
}

// ByComputationSteps orders the results by the computation_steps field.
func ByComputationSteps(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldComputationSteps, opts...).ToFunc()
}

// ByNumTransitiveLoads orders the results by the num_transitive_loads field.
func ByNumTransitiveLoads(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNumTransitiveLoads, opts...).ToFunc()
}

// ByPackageOverhead orders the results by the package_overhead field.
func ByPackageOverhead(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPackageOverhead, opts...).ToFunc()
}

// ByPackageMetricsCount orders the results by package_metrics count.
func ByPackageMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPackageMetricsStep(), opts...)
	}
}

// ByPackageMetrics orders the results by package_metrics terms.
func ByPackageMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPackageMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newPackageMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PackageMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, PackageMetricsTable, PackageMetricsPrimaryKey...),
	)
}
