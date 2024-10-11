// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/exectioninfo"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingbreakdown"
)

// TimingBreakdown is the model entity for the TimingBreakdown schema.
type TimingBreakdown struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Time holds the value of the "time" field.
	Time string `json:"time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TimingBreakdownQuery when eager-loading is set.
	Edges                          TimingBreakdownEdges `json:"edges"`
	exection_info_timing_breakdown *int
	selectValues                   sql.SelectValues
}

// TimingBreakdownEdges holds the relations/edges for other nodes in the graph.
type TimingBreakdownEdges struct {
	// ExecutionInfo holds the value of the execution_info edge.
	ExecutionInfo *ExectionInfo `json:"execution_info,omitempty"`
	// Child holds the value of the child edge.
	Child []*TimingChild `json:"child,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int

	namedChild map[string][]*TimingChild
}

// ExecutionInfoOrErr returns the ExecutionInfo value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TimingBreakdownEdges) ExecutionInfoOrErr() (*ExectionInfo, error) {
	if e.ExecutionInfo != nil {
		return e.ExecutionInfo, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: exectioninfo.Label}
	}
	return nil, &NotLoadedError{edge: "execution_info"}
}

// ChildOrErr returns the Child value or an error if the edge
// was not loaded in eager-loading.
func (e TimingBreakdownEdges) ChildOrErr() ([]*TimingChild, error) {
	if e.loadedTypes[1] {
		return e.Child, nil
	}
	return nil, &NotLoadedError{edge: "child"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TimingBreakdown) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case timingbreakdown.FieldID:
			values[i] = new(sql.NullInt64)
		case timingbreakdown.FieldName, timingbreakdown.FieldTime:
			values[i] = new(sql.NullString)
		case timingbreakdown.ForeignKeys[0]: // exection_info_timing_breakdown
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TimingBreakdown fields.
func (tb *TimingBreakdown) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case timingbreakdown.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tb.ID = int(value.Int64)
		case timingbreakdown.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				tb.Name = value.String
			}
		case timingbreakdown.FieldTime:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				tb.Time = value.String
			}
		case timingbreakdown.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field exection_info_timing_breakdown", value)
			} else if value.Valid {
				tb.exection_info_timing_breakdown = new(int)
				*tb.exection_info_timing_breakdown = int(value.Int64)
			}
		default:
			tb.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TimingBreakdown.
// This includes values selected through modifiers, order, etc.
func (tb *TimingBreakdown) Value(name string) (ent.Value, error) {
	return tb.selectValues.Get(name)
}

// QueryExecutionInfo queries the "execution_info" edge of the TimingBreakdown entity.
func (tb *TimingBreakdown) QueryExecutionInfo() *ExectionInfoQuery {
	return NewTimingBreakdownClient(tb.config).QueryExecutionInfo(tb)
}

// QueryChild queries the "child" edge of the TimingBreakdown entity.
func (tb *TimingBreakdown) QueryChild() *TimingChildQuery {
	return NewTimingBreakdownClient(tb.config).QueryChild(tb)
}

// Update returns a builder for updating this TimingBreakdown.
// Note that you need to call TimingBreakdown.Unwrap() before calling this method if this TimingBreakdown
// was returned from a transaction, and the transaction was committed or rolled back.
func (tb *TimingBreakdown) Update() *TimingBreakdownUpdateOne {
	return NewTimingBreakdownClient(tb.config).UpdateOne(tb)
}

// Unwrap unwraps the TimingBreakdown entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tb *TimingBreakdown) Unwrap() *TimingBreakdown {
	_tx, ok := tb.config.driver.(*txDriver)
	if !ok {
		panic("ent: TimingBreakdown is not a transactional entity")
	}
	tb.config.driver = _tx.drv
	return tb
}

// String implements the fmt.Stringer.
func (tb *TimingBreakdown) String() string {
	var builder strings.Builder
	builder.WriteString("TimingBreakdown(")
	builder.WriteString(fmt.Sprintf("id=%v, ", tb.ID))
	builder.WriteString("name=")
	builder.WriteString(tb.Name)
	builder.WriteString(", ")
	builder.WriteString("time=")
	builder.WriteString(tb.Time)
	builder.WriteByte(')')
	return builder.String()
}

// NamedChild returns the Child named value or an error if the edge was not
// loaded in eager-loading with this name.
func (tb *TimingBreakdown) NamedChild(name string) ([]*TimingChild, error) {
	if tb.Edges.namedChild == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := tb.Edges.namedChild[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (tb *TimingBreakdown) appendNamedChild(name string, edges ...*TimingChild) {
	if tb.Edges.namedChild == nil {
		tb.Edges.namedChild = make(map[string][]*TimingChild)
	}
	if len(edges) == 0 {
		tb.Edges.namedChild[name] = []*TimingChild{}
	} else {
		tb.Edges.namedChild[name] = append(tb.Edges.namedChild[name], edges...)
	}
}

// TimingBreakdowns is a parsable slice of TimingBreakdown.
type TimingBreakdowns []*TimingBreakdown
