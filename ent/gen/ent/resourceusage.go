// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/exectioninfo"
	"github.com/buildbarn/bb-portal/ent/gen/ent/resourceusage"
)

// ResourceUsage is the model entity for the ResourceUsage schema.
type ResourceUsage struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
	// ExecutionInfoID holds the value of the "execution_info_id" field.
	ExecutionInfoID int `json:"execution_info_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ResourceUsageQuery when eager-loading is set.
	Edges        ResourceUsageEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ResourceUsageEdges holds the relations/edges for other nodes in the graph.
type ResourceUsageEdges struct {
	// ExecutionInfo holds the value of the execution_info edge.
	ExecutionInfo *ExectionInfo `json:"execution_info,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// ExecutionInfoOrErr returns the ExecutionInfo value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ResourceUsageEdges) ExecutionInfoOrErr() (*ExectionInfo, error) {
	if e.ExecutionInfo != nil {
		return e.ExecutionInfo, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: exectioninfo.Label}
	}
	return nil, &NotLoadedError{edge: "execution_info"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ResourceUsage) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case resourceusage.FieldID, resourceusage.FieldExecutionInfoID:
			values[i] = new(sql.NullInt64)
		case resourceusage.FieldName, resourceusage.FieldValue:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ResourceUsage fields.
func (ru *ResourceUsage) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case resourceusage.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ru.ID = int(value.Int64)
		case resourceusage.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ru.Name = value.String
			}
		case resourceusage.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				ru.Value = value.String
			}
		case resourceusage.FieldExecutionInfoID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field execution_info_id", values[i])
			} else if value.Valid {
				ru.ExecutionInfoID = int(value.Int64)
			}
		default:
			ru.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// GetValue returns the ent.Value that was dynamically selected and assigned to the ResourceUsage.
// This includes values selected through modifiers, order, etc.
func (ru *ResourceUsage) GetValue(name string) (ent.Value, error) {
	return ru.selectValues.Get(name)
}

// QueryExecutionInfo queries the "execution_info" edge of the ResourceUsage entity.
func (ru *ResourceUsage) QueryExecutionInfo() *ExectionInfoQuery {
	return NewResourceUsageClient(ru.config).QueryExecutionInfo(ru)
}

// Update returns a builder for updating this ResourceUsage.
// Note that you need to call ResourceUsage.Unwrap() before calling this method if this ResourceUsage
// was returned from a transaction, and the transaction was committed or rolled back.
func (ru *ResourceUsage) Update() *ResourceUsageUpdateOne {
	return NewResourceUsageClient(ru.config).UpdateOne(ru)
}

// Unwrap unwraps the ResourceUsage entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ru *ResourceUsage) Unwrap() *ResourceUsage {
	_tx, ok := ru.config.driver.(*txDriver)
	if !ok {
		panic("ent: ResourceUsage is not a transactional entity")
	}
	ru.config.driver = _tx.drv
	return ru
}

// String implements the fmt.Stringer.
func (ru *ResourceUsage) String() string {
	var builder strings.Builder
	builder.WriteString("ResourceUsage(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ru.ID))
	builder.WriteString("name=")
	builder.WriteString(ru.Name)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(ru.Value)
	builder.WriteString(", ")
	builder.WriteString("execution_info_id=")
	builder.WriteString(fmt.Sprintf("%v", ru.ExecutionInfoID))
	builder.WriteByte(')')
	return builder.String()
}

// ResourceUsages is a parsable slice of ResourceUsage.
type ResourceUsages []*ResourceUsage
