// Code generated by ent, DO NOT EDIT.

package metrics

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the metrics type in the database.
	Label = "metrics"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeBazelInvocation holds the string denoting the bazel_invocation edge name in mutations.
	EdgeBazelInvocation = "bazel_invocation"
	// EdgeActionSummary holds the string denoting the action_summary edge name in mutations.
	EdgeActionSummary = "action_summary"
	// EdgeMemoryMetrics holds the string denoting the memory_metrics edge name in mutations.
	EdgeMemoryMetrics = "memory_metrics"
	// EdgeTargetMetrics holds the string denoting the target_metrics edge name in mutations.
	EdgeTargetMetrics = "target_metrics"
	// EdgePackageMetrics holds the string denoting the package_metrics edge name in mutations.
	EdgePackageMetrics = "package_metrics"
	// EdgeTimingMetrics holds the string denoting the timing_metrics edge name in mutations.
	EdgeTimingMetrics = "timing_metrics"
	// EdgeCumulativeMetrics holds the string denoting the cumulative_metrics edge name in mutations.
	EdgeCumulativeMetrics = "cumulative_metrics"
	// EdgeArtifactMetrics holds the string denoting the artifact_metrics edge name in mutations.
	EdgeArtifactMetrics = "artifact_metrics"
	// EdgeNetworkMetrics holds the string denoting the network_metrics edge name in mutations.
	EdgeNetworkMetrics = "network_metrics"
	// EdgeDynamicExecutionMetrics holds the string denoting the dynamic_execution_metrics edge name in mutations.
	EdgeDynamicExecutionMetrics = "dynamic_execution_metrics"
	// EdgeBuildGraphMetrics holds the string denoting the build_graph_metrics edge name in mutations.
	EdgeBuildGraphMetrics = "build_graph_metrics"
	// Table holds the table name of the metrics in the database.
	Table = "metrics"
	// BazelInvocationTable is the table that holds the bazel_invocation relation/edge.
	BazelInvocationTable = "metrics"
	// BazelInvocationInverseTable is the table name for the BazelInvocation entity.
	// It exists in this package in order to avoid circular dependency with the "bazelinvocation" package.
	BazelInvocationInverseTable = "bazel_invocations"
	// BazelInvocationColumn is the table column denoting the bazel_invocation relation/edge.
	BazelInvocationColumn = "bazel_invocation_metrics"
	// ActionSummaryTable is the table that holds the action_summary relation/edge.
	ActionSummaryTable = "action_summaries"
	// ActionSummaryInverseTable is the table name for the ActionSummary entity.
	// It exists in this package in order to avoid circular dependency with the "actionsummary" package.
	ActionSummaryInverseTable = "action_summaries"
	// ActionSummaryColumn is the table column denoting the action_summary relation/edge.
	ActionSummaryColumn = "metrics_action_summary"
	// MemoryMetricsTable is the table that holds the memory_metrics relation/edge. The primary key declared below.
	MemoryMetricsTable = "metrics_memory_metrics"
	// MemoryMetricsInverseTable is the table name for the MemoryMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "memorymetrics" package.
	MemoryMetricsInverseTable = "memory_metrics"
	// TargetMetricsTable is the table that holds the target_metrics relation/edge. The primary key declared below.
	TargetMetricsTable = "metrics_target_metrics"
	// TargetMetricsInverseTable is the table name for the TargetMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "targetmetrics" package.
	TargetMetricsInverseTable = "target_metrics"
	// PackageMetricsTable is the table that holds the package_metrics relation/edge. The primary key declared below.
	PackageMetricsTable = "metrics_package_metrics"
	// PackageMetricsInverseTable is the table name for the PackageMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "packagemetrics" package.
	PackageMetricsInverseTable = "package_metrics"
	// TimingMetricsTable is the table that holds the timing_metrics relation/edge. The primary key declared below.
	TimingMetricsTable = "metrics_timing_metrics"
	// TimingMetricsInverseTable is the table name for the TimingMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "timingmetrics" package.
	TimingMetricsInverseTable = "timing_metrics"
	// CumulativeMetricsTable is the table that holds the cumulative_metrics relation/edge. The primary key declared below.
	CumulativeMetricsTable = "metrics_cumulative_metrics"
	// CumulativeMetricsInverseTable is the table name for the CumulativeMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "cumulativemetrics" package.
	CumulativeMetricsInverseTable = "cumulative_metrics"
	// ArtifactMetricsTable is the table that holds the artifact_metrics relation/edge. The primary key declared below.
	ArtifactMetricsTable = "metrics_artifact_metrics"
	// ArtifactMetricsInverseTable is the table name for the ArtifactMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "artifactmetrics" package.
	ArtifactMetricsInverseTable = "artifact_metrics"
	// NetworkMetricsTable is the table that holds the network_metrics relation/edge. The primary key declared below.
	NetworkMetricsTable = "metrics_network_metrics"
	// NetworkMetricsInverseTable is the table name for the NetworkMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "networkmetrics" package.
	NetworkMetricsInverseTable = "network_metrics"
	// DynamicExecutionMetricsTable is the table that holds the dynamic_execution_metrics relation/edge. The primary key declared below.
	DynamicExecutionMetricsTable = "metrics_dynamic_execution_metrics"
	// DynamicExecutionMetricsInverseTable is the table name for the DynamicExecutionMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "dynamicexecutionmetrics" package.
	DynamicExecutionMetricsInverseTable = "dynamic_execution_metrics"
	// BuildGraphMetricsTable is the table that holds the build_graph_metrics relation/edge. The primary key declared below.
	BuildGraphMetricsTable = "metrics_build_graph_metrics"
	// BuildGraphMetricsInverseTable is the table name for the BuildGraphMetrics entity.
	// It exists in this package in order to avoid circular dependency with the "buildgraphmetrics" package.
	BuildGraphMetricsInverseTable = "build_graph_metrics"
)

// Columns holds all SQL columns for metrics fields.
var Columns = []string{
	FieldID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "metrics"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"bazel_invocation_metrics",
}

var (
	// MemoryMetricsPrimaryKey and MemoryMetricsColumn2 are the table columns denoting the
	// primary key for the memory_metrics relation (M2M).
	MemoryMetricsPrimaryKey = []string{"metrics_id", "memory_metrics_id"}
	// TargetMetricsPrimaryKey and TargetMetricsColumn2 are the table columns denoting the
	// primary key for the target_metrics relation (M2M).
	TargetMetricsPrimaryKey = []string{"metrics_id", "target_metrics_id"}
	// PackageMetricsPrimaryKey and PackageMetricsColumn2 are the table columns denoting the
	// primary key for the package_metrics relation (M2M).
	PackageMetricsPrimaryKey = []string{"metrics_id", "package_metrics_id"}
	// TimingMetricsPrimaryKey and TimingMetricsColumn2 are the table columns denoting the
	// primary key for the timing_metrics relation (M2M).
	TimingMetricsPrimaryKey = []string{"metrics_id", "timing_metrics_id"}
	// CumulativeMetricsPrimaryKey and CumulativeMetricsColumn2 are the table columns denoting the
	// primary key for the cumulative_metrics relation (M2M).
	CumulativeMetricsPrimaryKey = []string{"metrics_id", "cumulative_metrics_id"}
	// ArtifactMetricsPrimaryKey and ArtifactMetricsColumn2 are the table columns denoting the
	// primary key for the artifact_metrics relation (M2M).
	ArtifactMetricsPrimaryKey = []string{"metrics_id", "artifact_metrics_id"}
	// NetworkMetricsPrimaryKey and NetworkMetricsColumn2 are the table columns denoting the
	// primary key for the network_metrics relation (M2M).
	NetworkMetricsPrimaryKey = []string{"metrics_id", "network_metrics_id"}
	// DynamicExecutionMetricsPrimaryKey and DynamicExecutionMetricsColumn2 are the table columns denoting the
	// primary key for the dynamic_execution_metrics relation (M2M).
	DynamicExecutionMetricsPrimaryKey = []string{"metrics_id", "dynamic_execution_metrics_id"}
	// BuildGraphMetricsPrimaryKey and BuildGraphMetricsColumn2 are the table columns denoting the
	// primary key for the build_graph_metrics relation (M2M).
	BuildGraphMetricsPrimaryKey = []string{"metrics_id", "build_graph_metrics_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Metrics queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByBazelInvocationField orders the results by bazel_invocation field.
func ByBazelInvocationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBazelInvocationStep(), sql.OrderByField(field, opts...))
	}
}

// ByActionSummaryCount orders the results by action_summary count.
func ByActionSummaryCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newActionSummaryStep(), opts...)
	}
}

// ByActionSummary orders the results by action_summary terms.
func ByActionSummary(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newActionSummaryStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByMemoryMetricsCount orders the results by memory_metrics count.
func ByMemoryMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMemoryMetricsStep(), opts...)
	}
}

// ByMemoryMetrics orders the results by memory_metrics terms.
func ByMemoryMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMemoryMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByTargetMetricsCount orders the results by target_metrics count.
func ByTargetMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTargetMetricsStep(), opts...)
	}
}

// ByTargetMetrics orders the results by target_metrics terms.
func ByTargetMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTargetMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
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

// ByTimingMetricsCount orders the results by timing_metrics count.
func ByTimingMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTimingMetricsStep(), opts...)
	}
}

// ByTimingMetrics orders the results by timing_metrics terms.
func ByTimingMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTimingMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCumulativeMetricsCount orders the results by cumulative_metrics count.
func ByCumulativeMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCumulativeMetricsStep(), opts...)
	}
}

// ByCumulativeMetrics orders the results by cumulative_metrics terms.
func ByCumulativeMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCumulativeMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByArtifactMetricsCount orders the results by artifact_metrics count.
func ByArtifactMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newArtifactMetricsStep(), opts...)
	}
}

// ByArtifactMetrics orders the results by artifact_metrics terms.
func ByArtifactMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newArtifactMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByNetworkMetricsCount orders the results by network_metrics count.
func ByNetworkMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newNetworkMetricsStep(), opts...)
	}
}

// ByNetworkMetrics orders the results by network_metrics terms.
func ByNetworkMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNetworkMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByDynamicExecutionMetricsCount orders the results by dynamic_execution_metrics count.
func ByDynamicExecutionMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDynamicExecutionMetricsStep(), opts...)
	}
}

// ByDynamicExecutionMetrics orders the results by dynamic_execution_metrics terms.
func ByDynamicExecutionMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDynamicExecutionMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByBuildGraphMetricsCount orders the results by build_graph_metrics count.
func ByBuildGraphMetricsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBuildGraphMetricsStep(), opts...)
	}
}

// ByBuildGraphMetrics orders the results by build_graph_metrics terms.
func ByBuildGraphMetrics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBuildGraphMetricsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newBazelInvocationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BazelInvocationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, BazelInvocationTable, BazelInvocationColumn),
	)
}
func newActionSummaryStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ActionSummaryInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ActionSummaryTable, ActionSummaryColumn),
	)
}
func newMemoryMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MemoryMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, MemoryMetricsTable, MemoryMetricsPrimaryKey...),
	)
}
func newTargetMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TargetMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, TargetMetricsTable, TargetMetricsPrimaryKey...),
	)
}
func newPackageMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PackageMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, PackageMetricsTable, PackageMetricsPrimaryKey...),
	)
}
func newTimingMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TimingMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, TimingMetricsTable, TimingMetricsPrimaryKey...),
	)
}
func newCumulativeMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CumulativeMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, CumulativeMetricsTable, CumulativeMetricsPrimaryKey...),
	)
}
func newArtifactMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ArtifactMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, ArtifactMetricsTable, ArtifactMetricsPrimaryKey...),
	)
}
func newNetworkMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NetworkMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, NetworkMetricsTable, NetworkMetricsPrimaryKey...),
	)
}
func newDynamicExecutionMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DynamicExecutionMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, DynamicExecutionMetricsTable, DynamicExecutionMetricsPrimaryKey...),
	)
}
func newBuildGraphMetricsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BuildGraphMetricsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, BuildGraphMetricsTable, BuildGraphMetricsPrimaryKey...),
	)
}