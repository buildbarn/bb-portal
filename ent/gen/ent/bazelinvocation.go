// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventfile"
	"github.com/buildbarn/bb-portal/pkg/summary"
	"github.com/google/uuid"
)

// BazelInvocation is the model entity for the BazelInvocation schema.
type BazelInvocation struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// InvocationID holds the value of the "invocation_id" field.
	InvocationID uuid.UUID `json:"invocation_id,omitempty"`
	// StartedAt holds the value of the "started_at" field.
	StartedAt time.Time `json:"started_at,omitempty"`
	// EndedAt holds the value of the "ended_at" field.
	EndedAt time.Time `json:"ended_at,omitempty"`
	// ChangeNumber holds the value of the "change_number" field.
	ChangeNumber int32 `json:"change_number,omitempty"`
	// PatchsetNumber holds the value of the "patchset_number" field.
	PatchsetNumber int32 `json:"patchset_number,omitempty"`
	// Summary holds the value of the "summary" field.
	Summary summary.InvocationSummary `json:"summary,omitempty"`
	// BepCompleted holds the value of the "bep_completed" field.
	BepCompleted bool `json:"bep_completed,omitempty"`
	// StepLabel holds the value of the "step_label" field.
	StepLabel string `json:"step_label,omitempty"`
	// RelatedFiles holds the value of the "related_files" field.
	RelatedFiles map[string]string `json:"related_files,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BazelInvocationQuery when eager-loading is set.
	Edges                       BazelInvocationEdges `json:"edges"`
	build_invocations           *int
	event_file_bazel_invocation *int
	selectValues                sql.SelectValues
}

// BazelInvocationEdges holds the relations/edges for other nodes in the graph.
type BazelInvocationEdges struct {
	// EventFile holds the value of the event_file edge.
	EventFile *EventFile `json:"event_file,omitempty"`
	// Build holds the value of the build edge.
	Build *Build `json:"build,omitempty"`
	// Problems holds the value of the problems edge.
	Problems []*BazelInvocationProblem `json:"problems,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int

	namedProblems map[string][]*BazelInvocationProblem
}

// EventFileOrErr returns the EventFile value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BazelInvocationEdges) EventFileOrErr() (*EventFile, error) {
	if e.EventFile != nil {
		return e.EventFile, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: eventfile.Label}
	}
	return nil, &NotLoadedError{edge: "event_file"}
}

// BuildOrErr returns the Build value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BazelInvocationEdges) BuildOrErr() (*Build, error) {
	if e.Build != nil {
		return e.Build, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: build.Label}
	}
	return nil, &NotLoadedError{edge: "build"}
}

// ProblemsOrErr returns the Problems value or an error if the edge
// was not loaded in eager-loading.
func (e BazelInvocationEdges) ProblemsOrErr() ([]*BazelInvocationProblem, error) {
	if e.loadedTypes[2] {
		return e.Problems, nil
	}
	return nil, &NotLoadedError{edge: "problems"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*BazelInvocation) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case bazelinvocation.FieldSummary, bazelinvocation.FieldRelatedFiles:
			values[i] = new([]byte)
		case bazelinvocation.FieldBepCompleted:
			values[i] = new(sql.NullBool)
		case bazelinvocation.FieldID, bazelinvocation.FieldChangeNumber, bazelinvocation.FieldPatchsetNumber:
			values[i] = new(sql.NullInt64)
		case bazelinvocation.FieldStepLabel:
			values[i] = new(sql.NullString)
		case bazelinvocation.FieldStartedAt, bazelinvocation.FieldEndedAt:
			values[i] = new(sql.NullTime)
		case bazelinvocation.FieldInvocationID:
			values[i] = new(uuid.UUID)
		case bazelinvocation.ForeignKeys[0]: // build_invocations
			values[i] = new(sql.NullInt64)
		case bazelinvocation.ForeignKeys[1]: // event_file_bazel_invocation
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the BazelInvocation fields.
func (bi *BazelInvocation) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case bazelinvocation.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			bi.ID = int(value.Int64)
		case bazelinvocation.FieldInvocationID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field invocation_id", values[i])
			} else if value != nil {
				bi.InvocationID = *value
			}
		case bazelinvocation.FieldStartedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field started_at", values[i])
			} else if value.Valid {
				bi.StartedAt = value.Time
			}
		case bazelinvocation.FieldEndedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field ended_at", values[i])
			} else if value.Valid {
				bi.EndedAt = value.Time
			}
		case bazelinvocation.FieldChangeNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field change_number", values[i])
			} else if value.Valid {
				bi.ChangeNumber = int32(value.Int64)
			}
		case bazelinvocation.FieldPatchsetNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field patchset_number", values[i])
			} else if value.Valid {
				bi.PatchsetNumber = int32(value.Int64)
			}
		case bazelinvocation.FieldSummary:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field summary", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &bi.Summary); err != nil {
					return fmt.Errorf("unmarshal field summary: %w", err)
				}
			}
		case bazelinvocation.FieldBepCompleted:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field bep_completed", values[i])
			} else if value.Valid {
				bi.BepCompleted = value.Bool
			}
		case bazelinvocation.FieldStepLabel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field step_label", values[i])
			} else if value.Valid {
				bi.StepLabel = value.String
			}
		case bazelinvocation.FieldRelatedFiles:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field related_files", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &bi.RelatedFiles); err != nil {
					return fmt.Errorf("unmarshal field related_files: %w", err)
				}
			}
		case bazelinvocation.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field build_invocations", value)
			} else if value.Valid {
				bi.build_invocations = new(int)
				*bi.build_invocations = int(value.Int64)
			}
		case bazelinvocation.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field event_file_bazel_invocation", value)
			} else if value.Valid {
				bi.event_file_bazel_invocation = new(int)
				*bi.event_file_bazel_invocation = int(value.Int64)
			}
		default:
			bi.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the BazelInvocation.
// This includes values selected through modifiers, order, etc.
func (bi *BazelInvocation) Value(name string) (ent.Value, error) {
	return bi.selectValues.Get(name)
}

// QueryEventFile queries the "event_file" edge of the BazelInvocation entity.
func (bi *BazelInvocation) QueryEventFile() *EventFileQuery {
	return NewBazelInvocationClient(bi.config).QueryEventFile(bi)
}

// QueryBuild queries the "build" edge of the BazelInvocation entity.
func (bi *BazelInvocation) QueryBuild() *BuildQuery {
	return NewBazelInvocationClient(bi.config).QueryBuild(bi)
}

// QueryProblems queries the "problems" edge of the BazelInvocation entity.
func (bi *BazelInvocation) QueryProblems() *BazelInvocationProblemQuery {
	return NewBazelInvocationClient(bi.config).QueryProblems(bi)
}

// Update returns a builder for updating this BazelInvocation.
// Note that you need to call BazelInvocation.Unwrap() before calling this method if this BazelInvocation
// was returned from a transaction, and the transaction was committed or rolled back.
func (bi *BazelInvocation) Update() *BazelInvocationUpdateOne {
	return NewBazelInvocationClient(bi.config).UpdateOne(bi)
}

// Unwrap unwraps the BazelInvocation entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (bi *BazelInvocation) Unwrap() *BazelInvocation {
	_tx, ok := bi.config.driver.(*txDriver)
	if !ok {
		panic("ent: BazelInvocation is not a transactional entity")
	}
	bi.config.driver = _tx.drv
	return bi
}

// String implements the fmt.Stringer.
func (bi *BazelInvocation) String() string {
	var builder strings.Builder
	builder.WriteString("BazelInvocation(")
	builder.WriteString(fmt.Sprintf("id=%v, ", bi.ID))
	builder.WriteString("invocation_id=")
	builder.WriteString(fmt.Sprintf("%v", bi.InvocationID))
	builder.WriteString(", ")
	builder.WriteString("started_at=")
	builder.WriteString(bi.StartedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("ended_at=")
	builder.WriteString(bi.EndedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("change_number=")
	builder.WriteString(fmt.Sprintf("%v", bi.ChangeNumber))
	builder.WriteString(", ")
	builder.WriteString("patchset_number=")
	builder.WriteString(fmt.Sprintf("%v", bi.PatchsetNumber))
	builder.WriteString(", ")
	builder.WriteString("summary=")
	builder.WriteString(fmt.Sprintf("%v", bi.Summary))
	builder.WriteString(", ")
	builder.WriteString("bep_completed=")
	builder.WriteString(fmt.Sprintf("%v", bi.BepCompleted))
	builder.WriteString(", ")
	builder.WriteString("step_label=")
	builder.WriteString(bi.StepLabel)
	builder.WriteString(", ")
	builder.WriteString("related_files=")
	builder.WriteString(fmt.Sprintf("%v", bi.RelatedFiles))
	builder.WriteByte(')')
	return builder.String()
}

// NamedProblems returns the Problems named value or an error if the edge was not
// loaded in eager-loading with this name.
func (bi *BazelInvocation) NamedProblems(name string) ([]*BazelInvocationProblem, error) {
	if bi.Edges.namedProblems == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := bi.Edges.namedProblems[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (bi *BazelInvocation) appendNamedProblems(name string, edges ...*BazelInvocationProblem) {
	if bi.Edges.namedProblems == nil {
		bi.Edges.namedProblems = make(map[string][]*BazelInvocationProblem)
	}
	if len(edges) == 0 {
		bi.Edges.namedProblems[name] = []*BazelInvocationProblem{}
	} else {
		bi.Edges.namedProblems[name] = append(bi.Edges.namedProblems[name], edges...)
	}
}

// BazelInvocations is a parsable slice of BazelInvocation.
type BazelInvocations []*BazelInvocation
