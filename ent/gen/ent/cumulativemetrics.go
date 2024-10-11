// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/cumulativemetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
)

// CumulativeMetrics is the model entity for the CumulativeMetrics schema.
type CumulativeMetrics struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// NumAnalyses holds the value of the "num_analyses" field.
	NumAnalyses int32 `json:"num_analyses,omitempty"`
	// NumBuilds holds the value of the "num_builds" field.
	NumBuilds int32 `json:"num_builds,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CumulativeMetricsQuery when eager-loading is set.
	Edges                      CumulativeMetricsEdges `json:"edges"`
	metrics_cumulative_metrics *int
	selectValues               sql.SelectValues
}

// CumulativeMetricsEdges holds the relations/edges for other nodes in the graph.
type CumulativeMetricsEdges struct {
	// Metrics holds the value of the metrics edge.
	Metrics *Metrics `json:"metrics,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// MetricsOrErr returns the Metrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CumulativeMetricsEdges) MetricsOrErr() (*Metrics, error) {
	if e.Metrics != nil {
		return e.Metrics, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: metrics.Label}
	}
	return nil, &NotLoadedError{edge: "metrics"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CumulativeMetrics) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case cumulativemetrics.FieldID, cumulativemetrics.FieldNumAnalyses, cumulativemetrics.FieldNumBuilds:
			values[i] = new(sql.NullInt64)
		case cumulativemetrics.ForeignKeys[0]: // metrics_cumulative_metrics
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CumulativeMetrics fields.
func (cm *CumulativeMetrics) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case cumulativemetrics.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			cm.ID = int(value.Int64)
		case cumulativemetrics.FieldNumAnalyses:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field num_analyses", values[i])
			} else if value.Valid {
				cm.NumAnalyses = int32(value.Int64)
			}
		case cumulativemetrics.FieldNumBuilds:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field num_builds", values[i])
			} else if value.Valid {
				cm.NumBuilds = int32(value.Int64)
			}
		case cumulativemetrics.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field metrics_cumulative_metrics", value)
			} else if value.Valid {
				cm.metrics_cumulative_metrics = new(int)
				*cm.metrics_cumulative_metrics = int(value.Int64)
			}
		default:
			cm.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the CumulativeMetrics.
// This includes values selected through modifiers, order, etc.
func (cm *CumulativeMetrics) Value(name string) (ent.Value, error) {
	return cm.selectValues.Get(name)
}

// QueryMetrics queries the "metrics" edge of the CumulativeMetrics entity.
func (cm *CumulativeMetrics) QueryMetrics() *MetricsQuery {
	return NewCumulativeMetricsClient(cm.config).QueryMetrics(cm)
}

// Update returns a builder for updating this CumulativeMetrics.
// Note that you need to call CumulativeMetrics.Unwrap() before calling this method if this CumulativeMetrics
// was returned from a transaction, and the transaction was committed or rolled back.
func (cm *CumulativeMetrics) Update() *CumulativeMetricsUpdateOne {
	return NewCumulativeMetricsClient(cm.config).UpdateOne(cm)
}

// Unwrap unwraps the CumulativeMetrics entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (cm *CumulativeMetrics) Unwrap() *CumulativeMetrics {
	_tx, ok := cm.config.driver.(*txDriver)
	if !ok {
		panic("ent: CumulativeMetrics is not a transactional entity")
	}
	cm.config.driver = _tx.drv
	return cm
}

// String implements the fmt.Stringer.
func (cm *CumulativeMetrics) String() string {
	var builder strings.Builder
	builder.WriteString("CumulativeMetrics(")
	builder.WriteString(fmt.Sprintf("id=%v, ", cm.ID))
	builder.WriteString("num_analyses=")
	builder.WriteString(fmt.Sprintf("%v", cm.NumAnalyses))
	builder.WriteString(", ")
	builder.WriteString("num_builds=")
	builder.WriteString(fmt.Sprintf("%v", cm.NumBuilds))
	builder.WriteByte(')')
	return builder.String()
}

// CumulativeMetricsSlice is a parsable slice of CumulativeMetrics.
type CumulativeMetricsSlice []*CumulativeMetrics
