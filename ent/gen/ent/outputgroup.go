// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/namedsetoffiles"
	"github.com/buildbarn/bb-portal/ent/gen/ent/outputgroup"
)

// OutputGroup is the model entity for the OutputGroup schema.
type OutputGroup struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Incomplete holds the value of the "incomplete" field.
	Incomplete bool `json:"incomplete,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OutputGroupQuery when eager-loading is set.
	Edges                  OutputGroupEdges `json:"edges"`
	output_group_file_sets *int
	selectValues           sql.SelectValues
}

// OutputGroupEdges holds the relations/edges for other nodes in the graph.
type OutputGroupEdges struct {
	// TargetComplete holds the value of the target_complete edge.
	TargetComplete []*TargetComplete `json:"target_complete,omitempty"`
	// InlineFiles holds the value of the inline_files edge.
	InlineFiles []*TestFile `json:"inline_files,omitempty"`
	// FileSets holds the value of the file_sets edge.
	FileSets *NamedSetOfFiles `json:"file_sets,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
	// totalCount holds the count of the edges above.
	totalCount [3]map[string]int

	namedTargetComplete map[string][]*TargetComplete
	namedInlineFiles    map[string][]*TestFile
}

// TargetCompleteOrErr returns the TargetComplete value or an error if the edge
// was not loaded in eager-loading.
func (e OutputGroupEdges) TargetCompleteOrErr() ([]*TargetComplete, error) {
	if e.loadedTypes[0] {
		return e.TargetComplete, nil
	}
	return nil, &NotLoadedError{edge: "target_complete"}
}

// InlineFilesOrErr returns the InlineFiles value or an error if the edge
// was not loaded in eager-loading.
func (e OutputGroupEdges) InlineFilesOrErr() ([]*TestFile, error) {
	if e.loadedTypes[1] {
		return e.InlineFiles, nil
	}
	return nil, &NotLoadedError{edge: "inline_files"}
}

// FileSetsOrErr returns the FileSets value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OutputGroupEdges) FileSetsOrErr() (*NamedSetOfFiles, error) {
	if e.FileSets != nil {
		return e.FileSets, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: namedsetoffiles.Label}
	}
	return nil, &NotLoadedError{edge: "file_sets"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*OutputGroup) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case outputgroup.FieldIncomplete:
			values[i] = new(sql.NullBool)
		case outputgroup.FieldID:
			values[i] = new(sql.NullInt64)
		case outputgroup.FieldName:
			values[i] = new(sql.NullString)
		case outputgroup.ForeignKeys[0]: // output_group_file_sets
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the OutputGroup fields.
func (og *OutputGroup) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case outputgroup.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			og.ID = int(value.Int64)
		case outputgroup.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				og.Name = value.String
			}
		case outputgroup.FieldIncomplete:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field incomplete", values[i])
			} else if value.Valid {
				og.Incomplete = value.Bool
			}
		case outputgroup.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field output_group_file_sets", value)
			} else if value.Valid {
				og.output_group_file_sets = new(int)
				*og.output_group_file_sets = int(value.Int64)
			}
		default:
			og.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the OutputGroup.
// This includes values selected through modifiers, order, etc.
func (og *OutputGroup) Value(name string) (ent.Value, error) {
	return og.selectValues.Get(name)
}

// QueryTargetComplete queries the "target_complete" edge of the OutputGroup entity.
func (og *OutputGroup) QueryTargetComplete() *TargetCompleteQuery {
	return NewOutputGroupClient(og.config).QueryTargetComplete(og)
}

// QueryInlineFiles queries the "inline_files" edge of the OutputGroup entity.
func (og *OutputGroup) QueryInlineFiles() *TestFileQuery {
	return NewOutputGroupClient(og.config).QueryInlineFiles(og)
}

// QueryFileSets queries the "file_sets" edge of the OutputGroup entity.
func (og *OutputGroup) QueryFileSets() *NamedSetOfFilesQuery {
	return NewOutputGroupClient(og.config).QueryFileSets(og)
}

// Update returns a builder for updating this OutputGroup.
// Note that you need to call OutputGroup.Unwrap() before calling this method if this OutputGroup
// was returned from a transaction, and the transaction was committed or rolled back.
func (og *OutputGroup) Update() *OutputGroupUpdateOne {
	return NewOutputGroupClient(og.config).UpdateOne(og)
}

// Unwrap unwraps the OutputGroup entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (og *OutputGroup) Unwrap() *OutputGroup {
	_tx, ok := og.config.driver.(*txDriver)
	if !ok {
		panic("ent: OutputGroup is not a transactional entity")
	}
	og.config.driver = _tx.drv
	return og
}

// String implements the fmt.Stringer.
func (og *OutputGroup) String() string {
	var builder strings.Builder
	builder.WriteString("OutputGroup(")
	builder.WriteString(fmt.Sprintf("id=%v, ", og.ID))
	builder.WriteString("name=")
	builder.WriteString(og.Name)
	builder.WriteString(", ")
	builder.WriteString("incomplete=")
	builder.WriteString(fmt.Sprintf("%v", og.Incomplete))
	builder.WriteByte(')')
	return builder.String()
}

// NamedTargetComplete returns the TargetComplete named value or an error if the edge was not
// loaded in eager-loading with this name.
func (og *OutputGroup) NamedTargetComplete(name string) ([]*TargetComplete, error) {
	if og.Edges.namedTargetComplete == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := og.Edges.namedTargetComplete[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (og *OutputGroup) appendNamedTargetComplete(name string, edges ...*TargetComplete) {
	if og.Edges.namedTargetComplete == nil {
		og.Edges.namedTargetComplete = make(map[string][]*TargetComplete)
	}
	if len(edges) == 0 {
		og.Edges.namedTargetComplete[name] = []*TargetComplete{}
	} else {
		og.Edges.namedTargetComplete[name] = append(og.Edges.namedTargetComplete[name], edges...)
	}
}

// NamedInlineFiles returns the InlineFiles named value or an error if the edge was not
// loaded in eager-loading with this name.
func (og *OutputGroup) NamedInlineFiles(name string) ([]*TestFile, error) {
	if og.Edges.namedInlineFiles == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := og.Edges.namedInlineFiles[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (og *OutputGroup) appendNamedInlineFiles(name string, edges ...*TestFile) {
	if og.Edges.namedInlineFiles == nil {
		og.Edges.namedInlineFiles = make(map[string][]*TestFile)
	}
	if len(edges) == 0 {
		og.Edges.namedInlineFiles[name] = []*TestFile{}
	} else {
		og.Edges.namedInlineFiles[name] = append(og.Edges.namedInlineFiles[name], edges...)
	}
}

// OutputGroups is a parsable slice of OutputGroup.
type OutputGroups []*OutputGroup
