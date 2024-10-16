// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/dynamicexecutionmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/racestatistics"
)

// RaceStatistics is the model entity for the RaceStatistics schema.
type RaceStatistics struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Mnemonic holds the value of the "mnemonic" field.
	Mnemonic string `json:"mnemonic,omitempty"`
	// LocalRunner holds the value of the "local_runner" field.
	LocalRunner string `json:"local_runner,omitempty"`
	// RemoteRunner holds the value of the "remote_runner" field.
	RemoteRunner string `json:"remote_runner,omitempty"`
	// LocalWins holds the value of the "local_wins" field.
	LocalWins int64 `json:"local_wins,omitempty"`
	// RenoteWins holds the value of the "renote_wins" field.
	RenoteWins int64 `json:"renote_wins,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the RaceStatisticsQuery when eager-loading is set.
	Edges                                     RaceStatisticsEdges `json:"edges"`
	dynamic_execution_metrics_race_statistics *int
	selectValues                              sql.SelectValues
}

// RaceStatisticsEdges holds the relations/edges for other nodes in the graph.
type RaceStatisticsEdges struct {
	// DynamicExecutionMetrics holds the value of the dynamic_execution_metrics edge.
	DynamicExecutionMetrics *DynamicExecutionMetrics `json:"dynamic_execution_metrics,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// DynamicExecutionMetricsOrErr returns the DynamicExecutionMetrics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e RaceStatisticsEdges) DynamicExecutionMetricsOrErr() (*DynamicExecutionMetrics, error) {
	if e.DynamicExecutionMetrics != nil {
		return e.DynamicExecutionMetrics, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: dynamicexecutionmetrics.Label}
	}
	return nil, &NotLoadedError{edge: "dynamic_execution_metrics"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*RaceStatistics) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case racestatistics.FieldID, racestatistics.FieldLocalWins, racestatistics.FieldRenoteWins:
			values[i] = new(sql.NullInt64)
		case racestatistics.FieldMnemonic, racestatistics.FieldLocalRunner, racestatistics.FieldRemoteRunner:
			values[i] = new(sql.NullString)
		case racestatistics.ForeignKeys[0]: // dynamic_execution_metrics_race_statistics
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the RaceStatistics fields.
func (rs *RaceStatistics) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case racestatistics.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			rs.ID = int(value.Int64)
		case racestatistics.FieldMnemonic:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mnemonic", values[i])
			} else if value.Valid {
				rs.Mnemonic = value.String
			}
		case racestatistics.FieldLocalRunner:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field local_runner", values[i])
			} else if value.Valid {
				rs.LocalRunner = value.String
			}
		case racestatistics.FieldRemoteRunner:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remote_runner", values[i])
			} else if value.Valid {
				rs.RemoteRunner = value.String
			}
		case racestatistics.FieldLocalWins:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field local_wins", values[i])
			} else if value.Valid {
				rs.LocalWins = value.Int64
			}
		case racestatistics.FieldRenoteWins:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field renote_wins", values[i])
			} else if value.Valid {
				rs.RenoteWins = value.Int64
			}
		case racestatistics.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field dynamic_execution_metrics_race_statistics", value)
			} else if value.Valid {
				rs.dynamic_execution_metrics_race_statistics = new(int)
				*rs.dynamic_execution_metrics_race_statistics = int(value.Int64)
			}
		default:
			rs.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the RaceStatistics.
// This includes values selected through modifiers, order, etc.
func (rs *RaceStatistics) Value(name string) (ent.Value, error) {
	return rs.selectValues.Get(name)
}

// QueryDynamicExecutionMetrics queries the "dynamic_execution_metrics" edge of the RaceStatistics entity.
func (rs *RaceStatistics) QueryDynamicExecutionMetrics() *DynamicExecutionMetricsQuery {
	return NewRaceStatisticsClient(rs.config).QueryDynamicExecutionMetrics(rs)
}

// Update returns a builder for updating this RaceStatistics.
// Note that you need to call RaceStatistics.Unwrap() before calling this method if this RaceStatistics
// was returned from a transaction, and the transaction was committed or rolled back.
func (rs *RaceStatistics) Update() *RaceStatisticsUpdateOne {
	return NewRaceStatisticsClient(rs.config).UpdateOne(rs)
}

// Unwrap unwraps the RaceStatistics entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (rs *RaceStatistics) Unwrap() *RaceStatistics {
	_tx, ok := rs.config.driver.(*txDriver)
	if !ok {
		panic("ent: RaceStatistics is not a transactional entity")
	}
	rs.config.driver = _tx.drv
	return rs
}

// String implements the fmt.Stringer.
func (rs *RaceStatistics) String() string {
	var builder strings.Builder
	builder.WriteString("RaceStatistics(")
	builder.WriteString(fmt.Sprintf("id=%v, ", rs.ID))
	builder.WriteString("mnemonic=")
	builder.WriteString(rs.Mnemonic)
	builder.WriteString(", ")
	builder.WriteString("local_runner=")
	builder.WriteString(rs.LocalRunner)
	builder.WriteString(", ")
	builder.WriteString("remote_runner=")
	builder.WriteString(rs.RemoteRunner)
	builder.WriteString(", ")
	builder.WriteString("local_wins=")
	builder.WriteString(fmt.Sprintf("%v", rs.LocalWins))
	builder.WriteString(", ")
	builder.WriteString("renote_wins=")
	builder.WriteString(fmt.Sprintf("%v", rs.RenoteWins))
	builder.WriteByte(')')
	return builder.String()
}

// RaceStatisticsSlice is a parsable slice of RaceStatistics.
type RaceStatisticsSlice []*RaceStatistics
