// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actioncachestatistics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/missdetail"
)

// MissDetail is the model entity for the MissDetail schema.
type MissDetail struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Reason holds the value of the "reason" field.
	Reason missdetail.Reason `json:"reason,omitempty"`
	// Count holds the value of the "count" field.
	Count int32 `json:"count,omitempty"`
	// ActionCacheStatisticsID holds the value of the "action_cache_statistics_id" field.
	ActionCacheStatisticsID int `json:"action_cache_statistics_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MissDetailQuery when eager-loading is set.
	Edges        MissDetailEdges `json:"edges"`
	selectValues sql.SelectValues
}

// MissDetailEdges holds the relations/edges for other nodes in the graph.
type MissDetailEdges struct {
	// ActionCacheStatistics holds the value of the action_cache_statistics edge.
	ActionCacheStatistics *ActionCacheStatistics `json:"action_cache_statistics,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// ActionCacheStatisticsOrErr returns the ActionCacheStatistics value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MissDetailEdges) ActionCacheStatisticsOrErr() (*ActionCacheStatistics, error) {
	if e.ActionCacheStatistics != nil {
		return e.ActionCacheStatistics, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: actioncachestatistics.Label}
	}
	return nil, &NotLoadedError{edge: "action_cache_statistics"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*MissDetail) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case missdetail.FieldID, missdetail.FieldCount, missdetail.FieldActionCacheStatisticsID:
			values[i] = new(sql.NullInt64)
		case missdetail.FieldReason:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the MissDetail fields.
func (md *MissDetail) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case missdetail.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			md.ID = int(value.Int64)
		case missdetail.FieldReason:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field reason", values[i])
			} else if value.Valid {
				md.Reason = missdetail.Reason(value.String)
			}
		case missdetail.FieldCount:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field count", values[i])
			} else if value.Valid {
				md.Count = int32(value.Int64)
			}
		case missdetail.FieldActionCacheStatisticsID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field action_cache_statistics_id", values[i])
			} else if value.Valid {
				md.ActionCacheStatisticsID = int(value.Int64)
			}
		default:
			md.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the MissDetail.
// This includes values selected through modifiers, order, etc.
func (md *MissDetail) Value(name string) (ent.Value, error) {
	return md.selectValues.Get(name)
}

// QueryActionCacheStatistics queries the "action_cache_statistics" edge of the MissDetail entity.
func (md *MissDetail) QueryActionCacheStatistics() *ActionCacheStatisticsQuery {
	return NewMissDetailClient(md.config).QueryActionCacheStatistics(md)
}

// Update returns a builder for updating this MissDetail.
// Note that you need to call MissDetail.Unwrap() before calling this method if this MissDetail
// was returned from a transaction, and the transaction was committed or rolled back.
func (md *MissDetail) Update() *MissDetailUpdateOne {
	return NewMissDetailClient(md.config).UpdateOne(md)
}

// Unwrap unwraps the MissDetail entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (md *MissDetail) Unwrap() *MissDetail {
	_tx, ok := md.config.driver.(*txDriver)
	if !ok {
		panic("ent: MissDetail is not a transactional entity")
	}
	md.config.driver = _tx.drv
	return md
}

// String implements the fmt.Stringer.
func (md *MissDetail) String() string {
	var builder strings.Builder
	builder.WriteString("MissDetail(")
	builder.WriteString(fmt.Sprintf("id=%v, ", md.ID))
	builder.WriteString("reason=")
	builder.WriteString(fmt.Sprintf("%v", md.Reason))
	builder.WriteString(", ")
	builder.WriteString("count=")
	builder.WriteString(fmt.Sprintf("%v", md.Count))
	builder.WriteString(", ")
	builder.WriteString("action_cache_statistics_id=")
	builder.WriteString(fmt.Sprintf("%v", md.ActionCacheStatisticsID))
	builder.WriteByte(')')
	return builder.String()
}

// MissDetails is a parsable slice of MissDetail.
type MissDetails []*MissDetail
