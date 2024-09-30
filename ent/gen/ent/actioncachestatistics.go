// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actioncachestatistics"
)

// ActionCacheStatistics is the model entity for the ActionCacheStatistics schema.
type ActionCacheStatistics struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// SizeInBytes holds the value of the "size_in_bytes" field.
	SizeInBytes uint64 `json:"size_in_bytes,omitempty"`
	// SaveTimeInMs holds the value of the "save_time_in_ms" field.
	SaveTimeInMs uint64 `json:"save_time_in_ms,omitempty"`
	// LoadTimeInMs holds the value of the "load_time_in_ms" field.
	LoadTimeInMs int64 `json:"load_time_in_ms,omitempty"`
	// Hits holds the value of the "hits" field.
	Hits int32 `json:"hits,omitempty"`
	// Misses holds the value of the "misses" field.
	Misses int32 `json:"misses,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ActionCacheStatisticsQuery when eager-loading is set.
	Edges        ActionCacheStatisticsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ActionCacheStatisticsEdges holds the relations/edges for other nodes in the graph.
type ActionCacheStatisticsEdges struct {
	// ActionSummary holds the value of the action_summary edge.
	ActionSummary []*ActionSummary `json:"action_summary,omitempty"`
	// MissDetails holds the value of the miss_details edge.
	MissDetails []*MissDetail `json:"miss_details,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int

	namedActionSummary map[string][]*ActionSummary
	namedMissDetails   map[string][]*MissDetail
}

// ActionSummaryOrErr returns the ActionSummary value or an error if the edge
// was not loaded in eager-loading.
func (e ActionCacheStatisticsEdges) ActionSummaryOrErr() ([]*ActionSummary, error) {
	if e.loadedTypes[0] {
		return e.ActionSummary, nil
	}
	return nil, &NotLoadedError{edge: "action_summary"}
}

// MissDetailsOrErr returns the MissDetails value or an error if the edge
// was not loaded in eager-loading.
func (e ActionCacheStatisticsEdges) MissDetailsOrErr() ([]*MissDetail, error) {
	if e.loadedTypes[1] {
		return e.MissDetails, nil
	}
	return nil, &NotLoadedError{edge: "miss_details"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ActionCacheStatistics) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case actioncachestatistics.FieldID, actioncachestatistics.FieldSizeInBytes, actioncachestatistics.FieldSaveTimeInMs, actioncachestatistics.FieldLoadTimeInMs, actioncachestatistics.FieldHits, actioncachestatistics.FieldMisses:
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ActionCacheStatistics fields.
func (acs *ActionCacheStatistics) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case actioncachestatistics.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			acs.ID = int(value.Int64)
		case actioncachestatistics.FieldSizeInBytes:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field size_in_bytes", values[i])
			} else if value.Valid {
				acs.SizeInBytes = uint64(value.Int64)
			}
		case actioncachestatistics.FieldSaveTimeInMs:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field save_time_in_ms", values[i])
			} else if value.Valid {
				acs.SaveTimeInMs = uint64(value.Int64)
			}
		case actioncachestatistics.FieldLoadTimeInMs:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field load_time_in_ms", values[i])
			} else if value.Valid {
				acs.LoadTimeInMs = value.Int64
			}
		case actioncachestatistics.FieldHits:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field hits", values[i])
			} else if value.Valid {
				acs.Hits = int32(value.Int64)
			}
		case actioncachestatistics.FieldMisses:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field misses", values[i])
			} else if value.Valid {
				acs.Misses = int32(value.Int64)
			}
		default:
			acs.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ActionCacheStatistics.
// This includes values selected through modifiers, order, etc.
func (acs *ActionCacheStatistics) Value(name string) (ent.Value, error) {
	return acs.selectValues.Get(name)
}

// QueryActionSummary queries the "action_summary" edge of the ActionCacheStatistics entity.
func (acs *ActionCacheStatistics) QueryActionSummary() *ActionSummaryQuery {
	return NewActionCacheStatisticsClient(acs.config).QueryActionSummary(acs)
}

// QueryMissDetails queries the "miss_details" edge of the ActionCacheStatistics entity.
func (acs *ActionCacheStatistics) QueryMissDetails() *MissDetailQuery {
	return NewActionCacheStatisticsClient(acs.config).QueryMissDetails(acs)
}

// Update returns a builder for updating this ActionCacheStatistics.
// Note that you need to call ActionCacheStatistics.Unwrap() before calling this method if this ActionCacheStatistics
// was returned from a transaction, and the transaction was committed or rolled back.
func (acs *ActionCacheStatistics) Update() *ActionCacheStatisticsUpdateOne {
	return NewActionCacheStatisticsClient(acs.config).UpdateOne(acs)
}

// Unwrap unwraps the ActionCacheStatistics entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (acs *ActionCacheStatistics) Unwrap() *ActionCacheStatistics {
	_tx, ok := acs.config.driver.(*txDriver)
	if !ok {
		panic("ent: ActionCacheStatistics is not a transactional entity")
	}
	acs.config.driver = _tx.drv
	return acs
}

// String implements the fmt.Stringer.
func (acs *ActionCacheStatistics) String() string {
	var builder strings.Builder
	builder.WriteString("ActionCacheStatistics(")
	builder.WriteString(fmt.Sprintf("id=%v, ", acs.ID))
	builder.WriteString("size_in_bytes=")
	builder.WriteString(fmt.Sprintf("%v", acs.SizeInBytes))
	builder.WriteString(", ")
	builder.WriteString("save_time_in_ms=")
	builder.WriteString(fmt.Sprintf("%v", acs.SaveTimeInMs))
	builder.WriteString(", ")
	builder.WriteString("load_time_in_ms=")
	builder.WriteString(fmt.Sprintf("%v", acs.LoadTimeInMs))
	builder.WriteString(", ")
	builder.WriteString("hits=")
	builder.WriteString(fmt.Sprintf("%v", acs.Hits))
	builder.WriteString(", ")
	builder.WriteString("misses=")
	builder.WriteString(fmt.Sprintf("%v", acs.Misses))
	builder.WriteByte(')')
	return builder.String()
}

// NamedActionSummary returns the ActionSummary named value or an error if the edge was not
// loaded in eager-loading with this name.
func (acs *ActionCacheStatistics) NamedActionSummary(name string) ([]*ActionSummary, error) {
	if acs.Edges.namedActionSummary == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := acs.Edges.namedActionSummary[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (acs *ActionCacheStatistics) appendNamedActionSummary(name string, edges ...*ActionSummary) {
	if acs.Edges.namedActionSummary == nil {
		acs.Edges.namedActionSummary = make(map[string][]*ActionSummary)
	}
	if len(edges) == 0 {
		acs.Edges.namedActionSummary[name] = []*ActionSummary{}
	} else {
		acs.Edges.namedActionSummary[name] = append(acs.Edges.namedActionSummary[name], edges...)
	}
}

// NamedMissDetails returns the MissDetails named value or an error if the edge was not
// loaded in eager-loading with this name.
func (acs *ActionCacheStatistics) NamedMissDetails(name string) ([]*MissDetail, error) {
	if acs.Edges.namedMissDetails == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := acs.Edges.namedMissDetails[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (acs *ActionCacheStatistics) appendNamedMissDetails(name string, edges ...*MissDetail) {
	if acs.Edges.namedMissDetails == nil {
		acs.Edges.namedMissDetails = make(map[string][]*MissDetail)
	}
	if len(edges) == 0 {
		acs.Edges.namedMissDetails[name] = []*MissDetail{}
	} else {
		acs.Edges.namedMissDetails[name] = append(acs.Edges.namedMissDetails[name], edges...)
	}
}

// ActionCacheStatisticsSlice is a parsable slice of ActionCacheStatistics.
type ActionCacheStatisticsSlice []*ActionCacheStatistics
