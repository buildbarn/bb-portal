// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingchild"
)

// TimingChild is the model entity for the TimingChild schema.
type TimingChild struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Time holds the value of the "time" field.
	Time string `json:"time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TimingChildQuery when eager-loading is set.
	Edges        TimingChildEdges `json:"edges"`
	selectValues sql.SelectValues
}

// TimingChildEdges holds the relations/edges for other nodes in the graph.
type TimingChildEdges struct {
	// TimingBreakdown holds the value of the timing_breakdown edge.
	TimingBreakdown []*TimingBreakdown `json:"timing_breakdown,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int

	namedTimingBreakdown map[string][]*TimingBreakdown
}

// TimingBreakdownOrErr returns the TimingBreakdown value or an error if the edge
// was not loaded in eager-loading.
func (e TimingChildEdges) TimingBreakdownOrErr() ([]*TimingBreakdown, error) {
	if e.loadedTypes[0] {
		return e.TimingBreakdown, nil
	}
	return nil, &NotLoadedError{edge: "timing_breakdown"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TimingChild) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case timingchild.FieldID:
			values[i] = new(sql.NullInt64)
		case timingchild.FieldName, timingchild.FieldTime:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TimingChild fields.
func (tc *TimingChild) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case timingchild.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tc.ID = int(value.Int64)
		case timingchild.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				tc.Name = value.String
			}
		case timingchild.FieldTime:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				tc.Time = value.String
			}
		default:
			tc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TimingChild.
// This includes values selected through modifiers, order, etc.
func (tc *TimingChild) Value(name string) (ent.Value, error) {
	return tc.selectValues.Get(name)
}

// QueryTimingBreakdown queries the "timing_breakdown" edge of the TimingChild entity.
func (tc *TimingChild) QueryTimingBreakdown() *TimingBreakdownQuery {
	return NewTimingChildClient(tc.config).QueryTimingBreakdown(tc)
}

// Update returns a builder for updating this TimingChild.
// Note that you need to call TimingChild.Unwrap() before calling this method if this TimingChild
// was returned from a transaction, and the transaction was committed or rolled back.
func (tc *TimingChild) Update() *TimingChildUpdateOne {
	return NewTimingChildClient(tc.config).UpdateOne(tc)
}

// Unwrap unwraps the TimingChild entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tc *TimingChild) Unwrap() *TimingChild {
	_tx, ok := tc.config.driver.(*txDriver)
	if !ok {
		panic("ent: TimingChild is not a transactional entity")
	}
	tc.config.driver = _tx.drv
	return tc
}

// String implements the fmt.Stringer.
func (tc *TimingChild) String() string {
	var builder strings.Builder
	builder.WriteString("TimingChild(")
	builder.WriteString(fmt.Sprintf("id=%v, ", tc.ID))
	builder.WriteString("name=")
	builder.WriteString(tc.Name)
	builder.WriteString(", ")
	builder.WriteString("time=")
	builder.WriteString(tc.Time)
	builder.WriteByte(')')
	return builder.String()
}

// NamedTimingBreakdown returns the TimingBreakdown named value or an error if the edge was not
// loaded in eager-loading with this name.
func (tc *TimingChild) NamedTimingBreakdown(name string) ([]*TimingBreakdown, error) {
	if tc.Edges.namedTimingBreakdown == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := tc.Edges.namedTimingBreakdown[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (tc *TimingChild) appendNamedTimingBreakdown(name string, edges ...*TimingBreakdown) {
	if tc.Edges.namedTimingBreakdown == nil {
		tc.Edges.namedTimingBreakdown = make(map[string][]*TimingBreakdown)
	}
	if len(edges) == 0 {
		tc.Edges.namedTimingBreakdown[name] = []*TimingBreakdown{}
	} else {
		tc.Edges.namedTimingBreakdown[name] = append(tc.Edges.namedTimingBreakdown[name], edges...)
	}
}

// TimingChilds is a parsable slice of TimingChild.
type TimingChilds []*TimingChild