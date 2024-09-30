// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actiondata"
)

// ActionData is the model entity for the ActionData schema.
type ActionData struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Mnemonic holds the value of the "mnemonic" field.
	Mnemonic string `json:"mnemonic,omitempty"`
	// ActionsExecuted holds the value of the "actions_executed" field.
	ActionsExecuted int64 `json:"actions_executed,omitempty"`
	// ActionsCreated holds the value of the "actions_created" field.
	ActionsCreated int64 `json:"actions_created,omitempty"`
	// FirstStartedMs holds the value of the "first_started_ms" field.
	FirstStartedMs int64 `json:"first_started_ms,omitempty"`
	// LastEndedMs holds the value of the "last_ended_ms" field.
	LastEndedMs int64 `json:"last_ended_ms,omitempty"`
	// SystemTime holds the value of the "system_time" field.
	SystemTime int64 `json:"system_time,omitempty"`
	// UserTime holds the value of the "user_time" field.
	UserTime int64 `json:"user_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ActionDataQuery when eager-loading is set.
	Edges        ActionDataEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ActionDataEdges holds the relations/edges for other nodes in the graph.
type ActionDataEdges struct {
	// ActionSummary holds the value of the action_summary edge.
	ActionSummary []*ActionSummary `json:"action_summary,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int

	namedActionSummary map[string][]*ActionSummary
}

// ActionSummaryOrErr returns the ActionSummary value or an error if the edge
// was not loaded in eager-loading.
func (e ActionDataEdges) ActionSummaryOrErr() ([]*ActionSummary, error) {
	if e.loadedTypes[0] {
		return e.ActionSummary, nil
	}
	return nil, &NotLoadedError{edge: "action_summary"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ActionData) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case actiondata.FieldID, actiondata.FieldActionsExecuted, actiondata.FieldActionsCreated, actiondata.FieldFirstStartedMs, actiondata.FieldLastEndedMs, actiondata.FieldSystemTime, actiondata.FieldUserTime:
			values[i] = new(sql.NullInt64)
		case actiondata.FieldMnemonic:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ActionData fields.
func (ad *ActionData) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case actiondata.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ad.ID = int(value.Int64)
		case actiondata.FieldMnemonic:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mnemonic", values[i])
			} else if value.Valid {
				ad.Mnemonic = value.String
			}
		case actiondata.FieldActionsExecuted:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field actions_executed", values[i])
			} else if value.Valid {
				ad.ActionsExecuted = value.Int64
			}
		case actiondata.FieldActionsCreated:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field actions_created", values[i])
			} else if value.Valid {
				ad.ActionsCreated = value.Int64
			}
		case actiondata.FieldFirstStartedMs:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field first_started_ms", values[i])
			} else if value.Valid {
				ad.FirstStartedMs = value.Int64
			}
		case actiondata.FieldLastEndedMs:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field last_ended_ms", values[i])
			} else if value.Valid {
				ad.LastEndedMs = value.Int64
			}
		case actiondata.FieldSystemTime:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field system_time", values[i])
			} else if value.Valid {
				ad.SystemTime = value.Int64
			}
		case actiondata.FieldUserTime:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field user_time", values[i])
			} else if value.Valid {
				ad.UserTime = value.Int64
			}
		default:
			ad.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ActionData.
// This includes values selected through modifiers, order, etc.
func (ad *ActionData) Value(name string) (ent.Value, error) {
	return ad.selectValues.Get(name)
}

// QueryActionSummary queries the "action_summary" edge of the ActionData entity.
func (ad *ActionData) QueryActionSummary() *ActionSummaryQuery {
	return NewActionDataClient(ad.config).QueryActionSummary(ad)
}

// Update returns a builder for updating this ActionData.
// Note that you need to call ActionData.Unwrap() before calling this method if this ActionData
// was returned from a transaction, and the transaction was committed or rolled back.
func (ad *ActionData) Update() *ActionDataUpdateOne {
	return NewActionDataClient(ad.config).UpdateOne(ad)
}

// Unwrap unwraps the ActionData entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ad *ActionData) Unwrap() *ActionData {
	_tx, ok := ad.config.driver.(*txDriver)
	if !ok {
		panic("ent: ActionData is not a transactional entity")
	}
	ad.config.driver = _tx.drv
	return ad
}

// String implements the fmt.Stringer.
func (ad *ActionData) String() string {
	var builder strings.Builder
	builder.WriteString("ActionData(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ad.ID))
	builder.WriteString("mnemonic=")
	builder.WriteString(ad.Mnemonic)
	builder.WriteString(", ")
	builder.WriteString("actions_executed=")
	builder.WriteString(fmt.Sprintf("%v", ad.ActionsExecuted))
	builder.WriteString(", ")
	builder.WriteString("actions_created=")
	builder.WriteString(fmt.Sprintf("%v", ad.ActionsCreated))
	builder.WriteString(", ")
	builder.WriteString("first_started_ms=")
	builder.WriteString(fmt.Sprintf("%v", ad.FirstStartedMs))
	builder.WriteString(", ")
	builder.WriteString("last_ended_ms=")
	builder.WriteString(fmt.Sprintf("%v", ad.LastEndedMs))
	builder.WriteString(", ")
	builder.WriteString("system_time=")
	builder.WriteString(fmt.Sprintf("%v", ad.SystemTime))
	builder.WriteString(", ")
	builder.WriteString("user_time=")
	builder.WriteString(fmt.Sprintf("%v", ad.UserTime))
	builder.WriteByte(')')
	return builder.String()
}

// NamedActionSummary returns the ActionSummary named value or an error if the edge was not
// loaded in eager-loading with this name.
func (ad *ActionData) NamedActionSummary(name string) ([]*ActionSummary, error) {
	if ad.Edges.namedActionSummary == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := ad.Edges.namedActionSummary[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (ad *ActionData) appendNamedActionSummary(name string, edges ...*ActionSummary) {
	if ad.Edges.namedActionSummary == nil {
		ad.Edges.namedActionSummary = make(map[string][]*ActionSummary)
	}
	if len(edges) == 0 {
		ad.Edges.namedActionSummary[name] = []*ActionSummary{}
	} else {
		ad.Edges.namedActionSummary[name] = append(ad.Edges.namedActionSummary[name], edges...)
	}
}

// ActionDataSlice is a parsable slice of ActionData.
type ActionDataSlice []*ActionData
