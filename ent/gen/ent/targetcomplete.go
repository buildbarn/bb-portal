// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/outputgroup"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetcomplete"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetpair"
)

// TargetComplete is the model entity for the TargetComplete schema.
type TargetComplete struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Success holds the value of the "success" field.
	Success bool `json:"success,omitempty"`
	// Tag holds the value of the "tag" field.
	Tag []string `json:"tag,omitempty"`
	// TargetKind holds the value of the "target_kind" field.
	TargetKind string `json:"target_kind,omitempty"`
	// EndTimeInMs holds the value of the "end_time_in_ms" field.
	EndTimeInMs int64 `json:"end_time_in_ms,omitempty"`
	// TestTimeoutSeconds holds the value of the "test_timeout_seconds" field.
	TestTimeoutSeconds int64 `json:"test_timeout_seconds,omitempty"`
	// TestTimeout holds the value of the "test_timeout" field.
	TestTimeout int64 `json:"test_timeout,omitempty"`
	// TestSize holds the value of the "test_size" field.
	TestSize targetcomplete.TestSize `json:"test_size,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TargetCompleteQuery when eager-loading is set.
	Edges                  TargetCompleteEdges `json:"edges"`
	target_pair_completion *int
	selectValues           sql.SelectValues
}

// TargetCompleteEdges holds the relations/edges for other nodes in the graph.
type TargetCompleteEdges struct {
	// TargetPair holds the value of the target_pair edge.
	TargetPair *TargetPair `json:"target_pair,omitempty"`
	// ImportantOutput holds the value of the important_output edge.
	ImportantOutput []*TestFile `json:"important_output,omitempty"`
	// DirectoryOutput holds the value of the directory_output edge.
	DirectoryOutput []*TestFile `json:"directory_output,omitempty"`
	// OutputGroup holds the value of the output_group edge.
	OutputGroup *OutputGroup `json:"output_group,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
	// totalCount holds the count of the edges above.
	totalCount [4]map[string]int

	namedImportantOutput map[string][]*TestFile
	namedDirectoryOutput map[string][]*TestFile
}

// TargetPairOrErr returns the TargetPair value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TargetCompleteEdges) TargetPairOrErr() (*TargetPair, error) {
	if e.TargetPair != nil {
		return e.TargetPair, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: targetpair.Label}
	}
	return nil, &NotLoadedError{edge: "target_pair"}
}

// ImportantOutputOrErr returns the ImportantOutput value or an error if the edge
// was not loaded in eager-loading.
func (e TargetCompleteEdges) ImportantOutputOrErr() ([]*TestFile, error) {
	if e.loadedTypes[1] {
		return e.ImportantOutput, nil
	}
	return nil, &NotLoadedError{edge: "important_output"}
}

// DirectoryOutputOrErr returns the DirectoryOutput value or an error if the edge
// was not loaded in eager-loading.
func (e TargetCompleteEdges) DirectoryOutputOrErr() ([]*TestFile, error) {
	if e.loadedTypes[2] {
		return e.DirectoryOutput, nil
	}
	return nil, &NotLoadedError{edge: "directory_output"}
}

// OutputGroupOrErr returns the OutputGroup value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TargetCompleteEdges) OutputGroupOrErr() (*OutputGroup, error) {
	if e.OutputGroup != nil {
		return e.OutputGroup, nil
	} else if e.loadedTypes[3] {
		return nil, &NotFoundError{label: outputgroup.Label}
	}
	return nil, &NotLoadedError{edge: "output_group"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TargetComplete) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case targetcomplete.FieldTag:
			values[i] = new([]byte)
		case targetcomplete.FieldSuccess:
			values[i] = new(sql.NullBool)
		case targetcomplete.FieldID, targetcomplete.FieldEndTimeInMs, targetcomplete.FieldTestTimeoutSeconds, targetcomplete.FieldTestTimeout:
			values[i] = new(sql.NullInt64)
		case targetcomplete.FieldTargetKind, targetcomplete.FieldTestSize:
			values[i] = new(sql.NullString)
		case targetcomplete.ForeignKeys[0]: // target_pair_completion
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TargetComplete fields.
func (tc *TargetComplete) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case targetcomplete.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tc.ID = int(value.Int64)
		case targetcomplete.FieldSuccess:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field success", values[i])
			} else if value.Valid {
				tc.Success = value.Bool
			}
		case targetcomplete.FieldTag:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field tag", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &tc.Tag); err != nil {
					return fmt.Errorf("unmarshal field tag: %w", err)
				}
			}
		case targetcomplete.FieldTargetKind:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field target_kind", values[i])
			} else if value.Valid {
				tc.TargetKind = value.String
			}
		case targetcomplete.FieldEndTimeInMs:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field end_time_in_ms", values[i])
			} else if value.Valid {
				tc.EndTimeInMs = value.Int64
			}
		case targetcomplete.FieldTestTimeoutSeconds:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field test_timeout_seconds", values[i])
			} else if value.Valid {
				tc.TestTimeoutSeconds = value.Int64
			}
		case targetcomplete.FieldTestTimeout:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field test_timeout", values[i])
			} else if value.Valid {
				tc.TestTimeout = value.Int64
			}
		case targetcomplete.FieldTestSize:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field test_size", values[i])
			} else if value.Valid {
				tc.TestSize = targetcomplete.TestSize(value.String)
			}
		case targetcomplete.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field target_pair_completion", value)
			} else if value.Valid {
				tc.target_pair_completion = new(int)
				*tc.target_pair_completion = int(value.Int64)
			}
		default:
			tc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TargetComplete.
// This includes values selected through modifiers, order, etc.
func (tc *TargetComplete) Value(name string) (ent.Value, error) {
	return tc.selectValues.Get(name)
}

// QueryTargetPair queries the "target_pair" edge of the TargetComplete entity.
func (tc *TargetComplete) QueryTargetPair() *TargetPairQuery {
	return NewTargetCompleteClient(tc.config).QueryTargetPair(tc)
}

// QueryImportantOutput queries the "important_output" edge of the TargetComplete entity.
func (tc *TargetComplete) QueryImportantOutput() *TestFileQuery {
	return NewTargetCompleteClient(tc.config).QueryImportantOutput(tc)
}

// QueryDirectoryOutput queries the "directory_output" edge of the TargetComplete entity.
func (tc *TargetComplete) QueryDirectoryOutput() *TestFileQuery {
	return NewTargetCompleteClient(tc.config).QueryDirectoryOutput(tc)
}

// QueryOutputGroup queries the "output_group" edge of the TargetComplete entity.
func (tc *TargetComplete) QueryOutputGroup() *OutputGroupQuery {
	return NewTargetCompleteClient(tc.config).QueryOutputGroup(tc)
}

// Update returns a builder for updating this TargetComplete.
// Note that you need to call TargetComplete.Unwrap() before calling this method if this TargetComplete
// was returned from a transaction, and the transaction was committed or rolled back.
func (tc *TargetComplete) Update() *TargetCompleteUpdateOne {
	return NewTargetCompleteClient(tc.config).UpdateOne(tc)
}

// Unwrap unwraps the TargetComplete entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tc *TargetComplete) Unwrap() *TargetComplete {
	_tx, ok := tc.config.driver.(*txDriver)
	if !ok {
		panic("ent: TargetComplete is not a transactional entity")
	}
	tc.config.driver = _tx.drv
	return tc
}

// String implements the fmt.Stringer.
func (tc *TargetComplete) String() string {
	var builder strings.Builder
	builder.WriteString("TargetComplete(")
	builder.WriteString(fmt.Sprintf("id=%v, ", tc.ID))
	builder.WriteString("success=")
	builder.WriteString(fmt.Sprintf("%v", tc.Success))
	builder.WriteString(", ")
	builder.WriteString("tag=")
	builder.WriteString(fmt.Sprintf("%v", tc.Tag))
	builder.WriteString(", ")
	builder.WriteString("target_kind=")
	builder.WriteString(tc.TargetKind)
	builder.WriteString(", ")
	builder.WriteString("end_time_in_ms=")
	builder.WriteString(fmt.Sprintf("%v", tc.EndTimeInMs))
	builder.WriteString(", ")
	builder.WriteString("test_timeout_seconds=")
	builder.WriteString(fmt.Sprintf("%v", tc.TestTimeoutSeconds))
	builder.WriteString(", ")
	builder.WriteString("test_timeout=")
	builder.WriteString(fmt.Sprintf("%v", tc.TestTimeout))
	builder.WriteString(", ")
	builder.WriteString("test_size=")
	builder.WriteString(fmt.Sprintf("%v", tc.TestSize))
	builder.WriteByte(')')
	return builder.String()
}

// NamedImportantOutput returns the ImportantOutput named value or an error if the edge was not
// loaded in eager-loading with this name.
func (tc *TargetComplete) NamedImportantOutput(name string) ([]*TestFile, error) {
	if tc.Edges.namedImportantOutput == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := tc.Edges.namedImportantOutput[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (tc *TargetComplete) appendNamedImportantOutput(name string, edges ...*TestFile) {
	if tc.Edges.namedImportantOutput == nil {
		tc.Edges.namedImportantOutput = make(map[string][]*TestFile)
	}
	if len(edges) == 0 {
		tc.Edges.namedImportantOutput[name] = []*TestFile{}
	} else {
		tc.Edges.namedImportantOutput[name] = append(tc.Edges.namedImportantOutput[name], edges...)
	}
}

// NamedDirectoryOutput returns the DirectoryOutput named value or an error if the edge was not
// loaded in eager-loading with this name.
func (tc *TargetComplete) NamedDirectoryOutput(name string) ([]*TestFile, error) {
	if tc.Edges.namedDirectoryOutput == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := tc.Edges.namedDirectoryOutput[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (tc *TargetComplete) appendNamedDirectoryOutput(name string, edges ...*TestFile) {
	if tc.Edges.namedDirectoryOutput == nil {
		tc.Edges.namedDirectoryOutput = make(map[string][]*TestFile)
	}
	if len(edges) == 0 {
		tc.Edges.namedDirectoryOutput[name] = []*TestFile{}
	} else {
		tc.Edges.namedDirectoryOutput[name] = append(tc.Edges.namedDirectoryOutput[name], edges...)
	}
}

// TargetCompletes is a parsable slice of TargetComplete.
type TargetCompletes []*TargetComplete
