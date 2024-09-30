// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/networkmetrics"
)

// NetworkMetrics is the model entity for the NetworkMetrics schema.
type NetworkMetrics struct {
	config
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the NetworkMetricsQuery when eager-loading is set.
	Edges        NetworkMetricsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// NetworkMetricsEdges holds the relations/edges for other nodes in the graph.
type NetworkMetricsEdges struct {
	// Metrics holds the value of the metrics edge.
	Metrics []*Metrics `json:"metrics,omitempty"`
	// SystemNetworkStats holds the value of the system_network_stats edge.
	SystemNetworkStats []*SystemNetworkStats `json:"system_network_stats,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int

	namedMetrics            map[string][]*Metrics
	namedSystemNetworkStats map[string][]*SystemNetworkStats
}

// MetricsOrErr returns the Metrics value or an error if the edge
// was not loaded in eager-loading.
func (e NetworkMetricsEdges) MetricsOrErr() ([]*Metrics, error) {
	if e.loadedTypes[0] {
		return e.Metrics, nil
	}
	return nil, &NotLoadedError{edge: "metrics"}
}

// SystemNetworkStatsOrErr returns the SystemNetworkStats value or an error if the edge
// was not loaded in eager-loading.
func (e NetworkMetricsEdges) SystemNetworkStatsOrErr() ([]*SystemNetworkStats, error) {
	if e.loadedTypes[1] {
		return e.SystemNetworkStats, nil
	}
	return nil, &NotLoadedError{edge: "system_network_stats"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*NetworkMetrics) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case networkmetrics.FieldID:
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the NetworkMetrics fields.
func (nm *NetworkMetrics) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case networkmetrics.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			nm.ID = int(value.Int64)
		default:
			nm.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the NetworkMetrics.
// This includes values selected through modifiers, order, etc.
func (nm *NetworkMetrics) Value(name string) (ent.Value, error) {
	return nm.selectValues.Get(name)
}

// QueryMetrics queries the "metrics" edge of the NetworkMetrics entity.
func (nm *NetworkMetrics) QueryMetrics() *MetricsQuery {
	return NewNetworkMetricsClient(nm.config).QueryMetrics(nm)
}

// QuerySystemNetworkStats queries the "system_network_stats" edge of the NetworkMetrics entity.
func (nm *NetworkMetrics) QuerySystemNetworkStats() *SystemNetworkStatsQuery {
	return NewNetworkMetricsClient(nm.config).QuerySystemNetworkStats(nm)
}

// Update returns a builder for updating this NetworkMetrics.
// Note that you need to call NetworkMetrics.Unwrap() before calling this method if this NetworkMetrics
// was returned from a transaction, and the transaction was committed or rolled back.
func (nm *NetworkMetrics) Update() *NetworkMetricsUpdateOne {
	return NewNetworkMetricsClient(nm.config).UpdateOne(nm)
}

// Unwrap unwraps the NetworkMetrics entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (nm *NetworkMetrics) Unwrap() *NetworkMetrics {
	_tx, ok := nm.config.driver.(*txDriver)
	if !ok {
		panic("ent: NetworkMetrics is not a transactional entity")
	}
	nm.config.driver = _tx.drv
	return nm
}

// String implements the fmt.Stringer.
func (nm *NetworkMetrics) String() string {
	var builder strings.Builder
	builder.WriteString("NetworkMetrics(")
	builder.WriteString(fmt.Sprintf("id=%v", nm.ID))
	builder.WriteByte(')')
	return builder.String()
}

// NamedMetrics returns the Metrics named value or an error if the edge was not
// loaded in eager-loading with this name.
func (nm *NetworkMetrics) NamedMetrics(name string) ([]*Metrics, error) {
	if nm.Edges.namedMetrics == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := nm.Edges.namedMetrics[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (nm *NetworkMetrics) appendNamedMetrics(name string, edges ...*Metrics) {
	if nm.Edges.namedMetrics == nil {
		nm.Edges.namedMetrics = make(map[string][]*Metrics)
	}
	if len(edges) == 0 {
		nm.Edges.namedMetrics[name] = []*Metrics{}
	} else {
		nm.Edges.namedMetrics[name] = append(nm.Edges.namedMetrics[name], edges...)
	}
}

// NamedSystemNetworkStats returns the SystemNetworkStats named value or an error if the edge was not
// loaded in eager-loading with this name.
func (nm *NetworkMetrics) NamedSystemNetworkStats(name string) ([]*SystemNetworkStats, error) {
	if nm.Edges.namedSystemNetworkStats == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := nm.Edges.namedSystemNetworkStats[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (nm *NetworkMetrics) appendNamedSystemNetworkStats(name string, edges ...*SystemNetworkStats) {
	if nm.Edges.namedSystemNetworkStats == nil {
		nm.Edges.namedSystemNetworkStats = make(map[string][]*SystemNetworkStats)
	}
	if len(edges) == 0 {
		nm.Edges.namedSystemNetworkStats[name] = []*SystemNetworkStats{}
	} else {
		nm.Edges.namedSystemNetworkStats[name] = append(nm.Edges.namedSystemNetworkStats[name], edges...)
	}
}

// NetworkMetricsSlice is a parsable slice of NetworkMetrics.
type NetworkMetricsSlice []*NetworkMetrics
