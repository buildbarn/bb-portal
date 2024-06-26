// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventfile"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// EventFileUpdate is the builder for updating EventFile entities.
type EventFileUpdate struct {
	config
	hooks    []Hook
	mutation *EventFileMutation
}

// Where appends a list predicates to the EventFileUpdate builder.
func (efu *EventFileUpdate) Where(ps ...predicate.EventFile) *EventFileUpdate {
	efu.mutation.Where(ps...)
	return efu
}

// SetModTime sets the "mod_time" field.
func (efu *EventFileUpdate) SetModTime(t time.Time) *EventFileUpdate {
	efu.mutation.SetModTime(t)
	return efu
}

// SetNillableModTime sets the "mod_time" field if the given value is not nil.
func (efu *EventFileUpdate) SetNillableModTime(t *time.Time) *EventFileUpdate {
	if t != nil {
		efu.SetModTime(*t)
	}
	return efu
}

// SetProtocol sets the "protocol" field.
func (efu *EventFileUpdate) SetProtocol(s string) *EventFileUpdate {
	efu.mutation.SetProtocol(s)
	return efu
}

// SetNillableProtocol sets the "protocol" field if the given value is not nil.
func (efu *EventFileUpdate) SetNillableProtocol(s *string) *EventFileUpdate {
	if s != nil {
		efu.SetProtocol(*s)
	}
	return efu
}

// SetMimeType sets the "mime_type" field.
func (efu *EventFileUpdate) SetMimeType(s string) *EventFileUpdate {
	efu.mutation.SetMimeType(s)
	return efu
}

// SetNillableMimeType sets the "mime_type" field if the given value is not nil.
func (efu *EventFileUpdate) SetNillableMimeType(s *string) *EventFileUpdate {
	if s != nil {
		efu.SetMimeType(*s)
	}
	return efu
}

// SetStatus sets the "status" field.
func (efu *EventFileUpdate) SetStatus(s string) *EventFileUpdate {
	efu.mutation.SetStatus(s)
	return efu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (efu *EventFileUpdate) SetNillableStatus(s *string) *EventFileUpdate {
	if s != nil {
		efu.SetStatus(*s)
	}
	return efu
}

// SetReason sets the "reason" field.
func (efu *EventFileUpdate) SetReason(s string) *EventFileUpdate {
	efu.mutation.SetReason(s)
	return efu
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (efu *EventFileUpdate) SetNillableReason(s *string) *EventFileUpdate {
	if s != nil {
		efu.SetReason(*s)
	}
	return efu
}

// ClearReason clears the value of the "reason" field.
func (efu *EventFileUpdate) ClearReason() *EventFileUpdate {
	efu.mutation.ClearReason()
	return efu
}

// SetBazelInvocationID sets the "bazel_invocation" edge to the BazelInvocation entity by ID.
func (efu *EventFileUpdate) SetBazelInvocationID(id int) *EventFileUpdate {
	efu.mutation.SetBazelInvocationID(id)
	return efu
}

// SetNillableBazelInvocationID sets the "bazel_invocation" edge to the BazelInvocation entity by ID if the given value is not nil.
func (efu *EventFileUpdate) SetNillableBazelInvocationID(id *int) *EventFileUpdate {
	if id != nil {
		efu = efu.SetBazelInvocationID(*id)
	}
	return efu
}

// SetBazelInvocation sets the "bazel_invocation" edge to the BazelInvocation entity.
func (efu *EventFileUpdate) SetBazelInvocation(b *BazelInvocation) *EventFileUpdate {
	return efu.SetBazelInvocationID(b.ID)
}

// Mutation returns the EventFileMutation object of the builder.
func (efu *EventFileUpdate) Mutation() *EventFileMutation {
	return efu.mutation
}

// ClearBazelInvocation clears the "bazel_invocation" edge to the BazelInvocation entity.
func (efu *EventFileUpdate) ClearBazelInvocation() *EventFileUpdate {
	efu.mutation.ClearBazelInvocation()
	return efu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (efu *EventFileUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, efu.sqlSave, efu.mutation, efu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (efu *EventFileUpdate) SaveX(ctx context.Context) int {
	affected, err := efu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (efu *EventFileUpdate) Exec(ctx context.Context) error {
	_, err := efu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (efu *EventFileUpdate) ExecX(ctx context.Context) {
	if err := efu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (efu *EventFileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(eventfile.Table, eventfile.Columns, sqlgraph.NewFieldSpec(eventfile.FieldID, field.TypeInt))
	if ps := efu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := efu.mutation.ModTime(); ok {
		_spec.SetField(eventfile.FieldModTime, field.TypeTime, value)
	}
	if value, ok := efu.mutation.Protocol(); ok {
		_spec.SetField(eventfile.FieldProtocol, field.TypeString, value)
	}
	if value, ok := efu.mutation.MimeType(); ok {
		_spec.SetField(eventfile.FieldMimeType, field.TypeString, value)
	}
	if value, ok := efu.mutation.Status(); ok {
		_spec.SetField(eventfile.FieldStatus, field.TypeString, value)
	}
	if value, ok := efu.mutation.Reason(); ok {
		_spec.SetField(eventfile.FieldReason, field.TypeString, value)
	}
	if efu.mutation.ReasonCleared() {
		_spec.ClearField(eventfile.FieldReason, field.TypeString)
	}
	if efu.mutation.BazelInvocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   eventfile.BazelInvocationTable,
			Columns: []string{eventfile.BazelInvocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bazelinvocation.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := efu.mutation.BazelInvocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   eventfile.BazelInvocationTable,
			Columns: []string{eventfile.BazelInvocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bazelinvocation.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, efu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{eventfile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	efu.mutation.done = true
	return n, nil
}

// EventFileUpdateOne is the builder for updating a single EventFile entity.
type EventFileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EventFileMutation
}

// SetModTime sets the "mod_time" field.
func (efuo *EventFileUpdateOne) SetModTime(t time.Time) *EventFileUpdateOne {
	efuo.mutation.SetModTime(t)
	return efuo
}

// SetNillableModTime sets the "mod_time" field if the given value is not nil.
func (efuo *EventFileUpdateOne) SetNillableModTime(t *time.Time) *EventFileUpdateOne {
	if t != nil {
		efuo.SetModTime(*t)
	}
	return efuo
}

// SetProtocol sets the "protocol" field.
func (efuo *EventFileUpdateOne) SetProtocol(s string) *EventFileUpdateOne {
	efuo.mutation.SetProtocol(s)
	return efuo
}

// SetNillableProtocol sets the "protocol" field if the given value is not nil.
func (efuo *EventFileUpdateOne) SetNillableProtocol(s *string) *EventFileUpdateOne {
	if s != nil {
		efuo.SetProtocol(*s)
	}
	return efuo
}

// SetMimeType sets the "mime_type" field.
func (efuo *EventFileUpdateOne) SetMimeType(s string) *EventFileUpdateOne {
	efuo.mutation.SetMimeType(s)
	return efuo
}

// SetNillableMimeType sets the "mime_type" field if the given value is not nil.
func (efuo *EventFileUpdateOne) SetNillableMimeType(s *string) *EventFileUpdateOne {
	if s != nil {
		efuo.SetMimeType(*s)
	}
	return efuo
}

// SetStatus sets the "status" field.
func (efuo *EventFileUpdateOne) SetStatus(s string) *EventFileUpdateOne {
	efuo.mutation.SetStatus(s)
	return efuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (efuo *EventFileUpdateOne) SetNillableStatus(s *string) *EventFileUpdateOne {
	if s != nil {
		efuo.SetStatus(*s)
	}
	return efuo
}

// SetReason sets the "reason" field.
func (efuo *EventFileUpdateOne) SetReason(s string) *EventFileUpdateOne {
	efuo.mutation.SetReason(s)
	return efuo
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (efuo *EventFileUpdateOne) SetNillableReason(s *string) *EventFileUpdateOne {
	if s != nil {
		efuo.SetReason(*s)
	}
	return efuo
}

// ClearReason clears the value of the "reason" field.
func (efuo *EventFileUpdateOne) ClearReason() *EventFileUpdateOne {
	efuo.mutation.ClearReason()
	return efuo
}

// SetBazelInvocationID sets the "bazel_invocation" edge to the BazelInvocation entity by ID.
func (efuo *EventFileUpdateOne) SetBazelInvocationID(id int) *EventFileUpdateOne {
	efuo.mutation.SetBazelInvocationID(id)
	return efuo
}

// SetNillableBazelInvocationID sets the "bazel_invocation" edge to the BazelInvocation entity by ID if the given value is not nil.
func (efuo *EventFileUpdateOne) SetNillableBazelInvocationID(id *int) *EventFileUpdateOne {
	if id != nil {
		efuo = efuo.SetBazelInvocationID(*id)
	}
	return efuo
}

// SetBazelInvocation sets the "bazel_invocation" edge to the BazelInvocation entity.
func (efuo *EventFileUpdateOne) SetBazelInvocation(b *BazelInvocation) *EventFileUpdateOne {
	return efuo.SetBazelInvocationID(b.ID)
}

// Mutation returns the EventFileMutation object of the builder.
func (efuo *EventFileUpdateOne) Mutation() *EventFileMutation {
	return efuo.mutation
}

// ClearBazelInvocation clears the "bazel_invocation" edge to the BazelInvocation entity.
func (efuo *EventFileUpdateOne) ClearBazelInvocation() *EventFileUpdateOne {
	efuo.mutation.ClearBazelInvocation()
	return efuo
}

// Where appends a list predicates to the EventFileUpdate builder.
func (efuo *EventFileUpdateOne) Where(ps ...predicate.EventFile) *EventFileUpdateOne {
	efuo.mutation.Where(ps...)
	return efuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (efuo *EventFileUpdateOne) Select(field string, fields ...string) *EventFileUpdateOne {
	efuo.fields = append([]string{field}, fields...)
	return efuo
}

// Save executes the query and returns the updated EventFile entity.
func (efuo *EventFileUpdateOne) Save(ctx context.Context) (*EventFile, error) {
	return withHooks(ctx, efuo.sqlSave, efuo.mutation, efuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (efuo *EventFileUpdateOne) SaveX(ctx context.Context) *EventFile {
	node, err := efuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (efuo *EventFileUpdateOne) Exec(ctx context.Context) error {
	_, err := efuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (efuo *EventFileUpdateOne) ExecX(ctx context.Context) {
	if err := efuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (efuo *EventFileUpdateOne) sqlSave(ctx context.Context) (_node *EventFile, err error) {
	_spec := sqlgraph.NewUpdateSpec(eventfile.Table, eventfile.Columns, sqlgraph.NewFieldSpec(eventfile.FieldID, field.TypeInt))
	id, ok := efuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "EventFile.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := efuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, eventfile.FieldID)
		for _, f := range fields {
			if !eventfile.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != eventfile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := efuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := efuo.mutation.ModTime(); ok {
		_spec.SetField(eventfile.FieldModTime, field.TypeTime, value)
	}
	if value, ok := efuo.mutation.Protocol(); ok {
		_spec.SetField(eventfile.FieldProtocol, field.TypeString, value)
	}
	if value, ok := efuo.mutation.MimeType(); ok {
		_spec.SetField(eventfile.FieldMimeType, field.TypeString, value)
	}
	if value, ok := efuo.mutation.Status(); ok {
		_spec.SetField(eventfile.FieldStatus, field.TypeString, value)
	}
	if value, ok := efuo.mutation.Reason(); ok {
		_spec.SetField(eventfile.FieldReason, field.TypeString, value)
	}
	if efuo.mutation.ReasonCleared() {
		_spec.ClearField(eventfile.FieldReason, field.TypeString)
	}
	if efuo.mutation.BazelInvocationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   eventfile.BazelInvocationTable,
			Columns: []string{eventfile.BazelInvocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bazelinvocation.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := efuo.mutation.BazelInvocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   eventfile.BazelInvocationTable,
			Columns: []string{eventfile.BazelInvocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bazelinvocation.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &EventFile{config: efuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, efuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{eventfile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	efuo.mutation.done = true
	return _node, nil
}
