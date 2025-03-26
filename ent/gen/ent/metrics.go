// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actionsummary"
	"github.com/buildbarn/bb-portal/ent/gen/ent/artifactmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildgraphmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/cumulativemetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/dynamicexecutionmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/memorymetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/networkmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/packagemetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingmetrics"
)

// Metrics is the model entity for the Metrics schema.
type Metrics struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// BazelInvocationID holds the value of the "bazel_invocation_id" field.
	BazelInvocationID int `json:"bazel_invocation_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MetricsQuery when eager-loading is set.
	Edges        MetricsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// MetricsEdges holds the relations/edges for other nodes in the graph.
type MetricsEdges struct {
	// BazelInvocation holds the value of the bazel_invocation edge.
	BazelInvocation *BazelInvocation `json:"bazel_invocation,omitempty"`
	// ActionSummary holds the value of the action_summary edge.
	ActionSummary *ActionSummary `json:"action_summary,omitempty"`
	// MemoryMetrics holds the value of the memory_metrics edge.
	MemoryMetrics *MemoryMetrics `json:"memory_metrics,omitempty"`
	// TargetMetrics holds the value of the target_metrics edge.
	TargetMetrics *TargetMetrics `json:"target_metrics,omitempty"`
	// PackageMetrics holds the value of the package_metrics edge.
	PackageMetrics *PackageMetrics `json:"package_metrics,omitempty"`
	// TimingMetrics holds the value of the timing_metrics edge.
	TimingMetrics *TimingMetrics `json:"timing_metrics,omitempty"`
	// CumulativeMetrics holds the value of the cumulative_metrics edge.
	CumulativeMetrics *CumulativeMetrics `json:"cumulative_metrics,omitempty"`
	// ArtifactMetrics holds the value of the artifact_metrics edge.
	ArtifactMetrics *ArtifactMetrics `json:"artifact_metrics,omitempty"`
	// NetworkMetrics holds the value of the network_metrics edge.
	NetworkMetrics *NetworkMetrics `json:"network_metrics,omitempty"`
	// DynamicExecutionMetrics holds the value of the dynamic_execution_metrics edge.
	DynamicExecutionMetrics *DynamicExecutionMetrics `json:"dynamic_execution_metrics,omitempty"`
	// BuildGraphMetrics holds the value of the build_graph_metrics edge.
	BuildGraphMetrics *BuildGraphMetrics `json:"build_graph_metrics,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [11]bool
	// totalCount holds the count of the edges above.
	totalCount [11]map[string]int
}

// BazelInvocationOrErr returns the BazelInvocation value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) BazelInvocationOrErr() (*BazelInvocation, error) {
	if e.BazelInvocation != nil {
		return e.BazelInvocation, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: bazelinvocation.Label}
	}
	return nil, &NotLoadedError{edge: "bazel_invocation"}
}

// ActionSummaryOrErr returns the ActionSummary value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) ActionSummaryOrErr() (*ActionSummary, error) {
	if e.ActionSummary != nil {
		return e.ActionSummary, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: actionsummary.Label}
	}
	return nil, &NotLoadedError{edge: "action_summary"}
}

// MemoryMetricsOrErr returns the MemoryMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) MemoryMetricsOrErr() (*MemoryMetrics, error) {
	if e.MemoryMetrics != nil {
		return e.MemoryMetrics, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: memorymetrics.Label}
	}
	return nil, &NotLoadedError{edge: "memory_metrics"}
}

// TargetMetricsOrErr returns the TargetMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) TargetMetricsOrErr() (*TargetMetrics, error) {
	if e.TargetMetrics != nil {
		return e.TargetMetrics, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: targetmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "target_metrics"}
}

// PackageMetricsOrErr returns the PackageMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) PackageMetricsOrErr() (*PackageMetrics, error) {
	if e.PackageMetrics != nil {
		return e.PackageMetrics, nil
	} else if e.loadedTypes[4] {
		return nil, &NotFoundError{label: packagemetrics.Label}
	}
	return nil, &NotLoadedError{edge: "package_metrics"}
}

// TimingMetricsOrErr returns the TimingMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) TimingMetricsOrErr() (*TimingMetrics, error) {
	if e.TimingMetrics != nil {
		return e.TimingMetrics, nil
	} else if e.loadedTypes[5] {
		return nil, &NotFoundError{label: timingmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "timing_metrics"}
}

// CumulativeMetricsOrErr returns the CumulativeMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) CumulativeMetricsOrErr() (*CumulativeMetrics, error) {
	if e.CumulativeMetrics != nil {
		return e.CumulativeMetrics, nil
	} else if e.loadedTypes[6] {
		return nil, &NotFoundError{label: cumulativemetrics.Label}
	}
	return nil, &NotLoadedError{edge: "cumulative_metrics"}
}

// ArtifactMetricsOrErr returns the ArtifactMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) ArtifactMetricsOrErr() (*ArtifactMetrics, error) {
	if e.ArtifactMetrics != nil {
		return e.ArtifactMetrics, nil
	} else if e.loadedTypes[7] {
		return nil, &NotFoundError{label: artifactmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "artifact_metrics"}
}

// NetworkMetricsOrErr returns the NetworkMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) NetworkMetricsOrErr() (*NetworkMetrics, error) {
	if e.NetworkMetrics != nil {
		return e.NetworkMetrics, nil
	} else if e.loadedTypes[8] {
		return nil, &NotFoundError{label: networkmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "network_metrics"}
}

// DynamicExecutionMetricsOrErr returns the DynamicExecutionMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) DynamicExecutionMetricsOrErr() (*DynamicExecutionMetrics, error) {
	if e.DynamicExecutionMetrics != nil {
		return e.DynamicExecutionMetrics, nil
	} else if e.loadedTypes[9] {
		return nil, &NotFoundError{label: dynamicexecutionmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "dynamic_execution_metrics"}
}

// BuildGraphMetricsOrErr returns the BuildGraphMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MetricsEdges) BuildGraphMetricsOrErr() (*BuildGraphMetrics, error) {
	if e.BuildGraphMetrics != nil {
		return e.BuildGraphMetrics, nil
	} else if e.loadedTypes[10] {
		return nil, &NotFoundError{label: buildgraphmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "build_graph_metrics"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Metrics) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case metrics.FieldID, metrics.FieldBazelInvocationID:
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Metrics fields.
func (m *Metrics) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case metrics.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			m.ID = int(value.Int64)
		case metrics.FieldBazelInvocationID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field bazel_invocation_id", values[i])
			} else if value.Valid {
				m.BazelInvocationID = int(value.Int64)
			}
		default:
			m.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Metrics.
// This includes values selected through modifiers, order, etc.
func (m *Metrics) Value(name string) (ent.Value, error) {
	return m.selectValues.Get(name)
}

// QueryBazelInvocation queries the "bazel_invocation" edge of the Metrics entity.
func (m *Metrics) QueryBazelInvocation() *BazelInvocationQuery {
	return NewMetricsClient(m.config).QueryBazelInvocation(m)
}

// QueryActionSummary queries the "action_summary" edge of the Metrics entity.
func (m *Metrics) QueryActionSummary() *ActionSummaryQuery {
	return NewMetricsClient(m.config).QueryActionSummary(m)
}

// QueryMemoryMetrics queries the "memory_metrics" edge of the Metrics entity.
func (m *Metrics) QueryMemoryMetrics() *MemoryMetricsQuery {
	return NewMetricsClient(m.config).QueryMemoryMetrics(m)
}

// QueryTargetMetrics queries the "target_metrics" edge of the Metrics entity.
func (m *Metrics) QueryTargetMetrics() *TargetMetricsQuery {
	return NewMetricsClient(m.config).QueryTargetMetrics(m)
}

// QueryPackageMetrics queries the "package_metrics" edge of the Metrics entity.
func (m *Metrics) QueryPackageMetrics() *PackageMetricsQuery {
	return NewMetricsClient(m.config).QueryPackageMetrics(m)
}

// QueryTimingMetrics queries the "timing_metrics" edge of the Metrics entity.
func (m *Metrics) QueryTimingMetrics() *TimingMetricsQuery {
	return NewMetricsClient(m.config).QueryTimingMetrics(m)
}

// QueryCumulativeMetrics queries the "cumulative_metrics" edge of the Metrics entity.
func (m *Metrics) QueryCumulativeMetrics() *CumulativeMetricsQuery {
	return NewMetricsClient(m.config).QueryCumulativeMetrics(m)
}

// QueryArtifactMetrics queries the "artifact_metrics" edge of the Metrics entity.
func (m *Metrics) QueryArtifactMetrics() *ArtifactMetricsQuery {
	return NewMetricsClient(m.config).QueryArtifactMetrics(m)
}

// QueryNetworkMetrics queries the "network_metrics" edge of the Metrics entity.
func (m *Metrics) QueryNetworkMetrics() *NetworkMetricsQuery {
	return NewMetricsClient(m.config).QueryNetworkMetrics(m)
}

// QueryDynamicExecutionMetrics queries the "dynamic_execution_metrics" edge of the Metrics entity.
func (m *Metrics) QueryDynamicExecutionMetrics() *DynamicExecutionMetricsQuery {
	return NewMetricsClient(m.config).QueryDynamicExecutionMetrics(m)
}

// QueryBuildGraphMetrics queries the "build_graph_metrics" edge of the Metrics entity.
func (m *Metrics) QueryBuildGraphMetrics() *BuildGraphMetricsQuery {
	return NewMetricsClient(m.config).QueryBuildGraphMetrics(m)
}

// Update returns a builder for updating this Metrics.
// Note that you need to call Metrics.Unwrap() before calling this method if this Metrics
// was returned from a transaction, and the transaction was committed or rolled back.
func (m *Metrics) Update() *MetricsUpdateOne {
	return NewMetricsClient(m.config).UpdateOne(m)
}

// Unwrap unwraps the Metrics entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (m *Metrics) Unwrap() *Metrics {
	_tx, ok := m.config.driver.(*txDriver)
	if !ok {
		panic("ent: Metrics is not a transactional entity")
	}
	m.config.driver = _tx.drv
	return m
}

// String implements the fmt.Stringer.
func (m *Metrics) String() string {
	var builder strings.Builder
	builder.WriteString("Metrics(")
	builder.WriteString(fmt.Sprintf("id=%v, ", m.ID))
	builder.WriteString("bazel_invocation_id=")
	builder.WriteString(fmt.Sprintf("%v", m.BazelInvocationID))
	builder.WriteByte(')')
	return builder.String()
}

// MetricsSlice is a parsable slice of Metrics.
type MetricsSlice []*Metrics
