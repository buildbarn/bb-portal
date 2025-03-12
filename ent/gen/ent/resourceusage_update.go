// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/exectioninfo"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/resourceusage"
)

// ResourceUsageUpdate is the builder for updating ResourceUsage entities.
type ResourceUsageUpdate struct {
	config
	hooks    []Hook
	mutation *ResourceUsageMutation
}

// Where appends a list predicates to the ResourceUsageUpdate builder.
func (ruu *ResourceUsageUpdate) Where(ps ...predicate.ResourceUsage) *ResourceUsageUpdate {
	ruu.mutation.Where(ps...)
	return ruu
}

// SetName sets the "name" field.
func (ruu *ResourceUsageUpdate) SetName(s string) *ResourceUsageUpdate {
	ruu.mutation.SetName(s)
	return ruu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ruu *ResourceUsageUpdate) SetNillableName(s *string) *ResourceUsageUpdate {
	if s != nil {
		ruu.SetName(*s)
	}
	return ruu
}

// ClearName clears the value of the "name" field.
func (ruu *ResourceUsageUpdate) ClearName() *ResourceUsageUpdate {
	ruu.mutation.ClearName()
	return ruu
}

// SetValue sets the "value" field.
func (ruu *ResourceUsageUpdate) SetValue(s string) *ResourceUsageUpdate {
	ruu.mutation.SetValue(s)
	return ruu
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (ruu *ResourceUsageUpdate) SetNillableValue(s *string) *ResourceUsageUpdate {
	if s != nil {
		ruu.SetValue(*s)
	}
	return ruu
}

// ClearValue clears the value of the "value" field.
func (ruu *ResourceUsageUpdate) ClearValue() *ResourceUsageUpdate {
	ruu.mutation.ClearValue()
	return ruu
}

// SetExecutionInfoID sets the "execution_info_id" field.
func (ruu *ResourceUsageUpdate) SetExecutionInfoID(i int) *ResourceUsageUpdate {
	ruu.mutation.SetExecutionInfoID(i)
	return ruu
}

// SetNillableExecutionInfoID sets the "execution_info_id" field if the given value is not nil.
func (ruu *ResourceUsageUpdate) SetNillableExecutionInfoID(i *int) *ResourceUsageUpdate {
	if i != nil {
		ruu.SetExecutionInfoID(*i)
	}
	return ruu
}

// ClearExecutionInfoID clears the value of the "execution_info_id" field.
func (ruu *ResourceUsageUpdate) ClearExecutionInfoID() *ResourceUsageUpdate {
	ruu.mutation.ClearExecutionInfoID()
	return ruu
}

// SetExecutionInfo sets the "execution_info" edge to the ExectionInfo entity.
func (ruu *ResourceUsageUpdate) SetExecutionInfo(e *ExectionInfo) *ResourceUsageUpdate {
	return ruu.SetExecutionInfoID(e.ID)
}

// Mutation returns the ResourceUsageMutation object of the builder.
func (ruu *ResourceUsageUpdate) Mutation() *ResourceUsageMutation {
	return ruu.mutation
}

// ClearExecutionInfo clears the "execution_info" edge to the ExectionInfo entity.
func (ruu *ResourceUsageUpdate) ClearExecutionInfo() *ResourceUsageUpdate {
	ruu.mutation.ClearExecutionInfo()
	return ruu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ruu *ResourceUsageUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ruu.sqlSave, ruu.mutation, ruu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruu *ResourceUsageUpdate) SaveX(ctx context.Context) int {
	affected, err := ruu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ruu *ResourceUsageUpdate) Exec(ctx context.Context) error {
	_, err := ruu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruu *ResourceUsageUpdate) ExecX(ctx context.Context) {
	if err := ruu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ruu *ResourceUsageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(resourceusage.Table, resourceusage.Columns, sqlgraph.NewFieldSpec(resourceusage.FieldID, field.TypeInt))
	if ps := ruu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruu.mutation.Name(); ok {
		_spec.SetField(resourceusage.FieldName, field.TypeString, value)
	}
	if ruu.mutation.NameCleared() {
		_spec.ClearField(resourceusage.FieldName, field.TypeString)
	}
	if value, ok := ruu.mutation.Value(); ok {
		_spec.SetField(resourceusage.FieldValue, field.TypeString, value)
	}
	if ruu.mutation.ValueCleared() {
		_spec.ClearField(resourceusage.FieldValue, field.TypeString)
	}
	if ruu.mutation.ExecutionInfoCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   resourceusage.ExecutionInfoTable,
			Columns: []string{resourceusage.ExecutionInfoColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exectioninfo.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruu.mutation.ExecutionInfoIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   resourceusage.ExecutionInfoTable,
			Columns: []string{resourceusage.ExecutionInfoColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exectioninfo.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ruu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resourceusage.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ruu.mutation.done = true
	return n, nil
}

// ResourceUsageUpdateOne is the builder for updating a single ResourceUsage entity.
type ResourceUsageUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ResourceUsageMutation
}

// SetName sets the "name" field.
func (ruuo *ResourceUsageUpdateOne) SetName(s string) *ResourceUsageUpdateOne {
	ruuo.mutation.SetName(s)
	return ruuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ruuo *ResourceUsageUpdateOne) SetNillableName(s *string) *ResourceUsageUpdateOne {
	if s != nil {
		ruuo.SetName(*s)
	}
	return ruuo
}

// ClearName clears the value of the "name" field.
func (ruuo *ResourceUsageUpdateOne) ClearName() *ResourceUsageUpdateOne {
	ruuo.mutation.ClearName()
	return ruuo
}

// SetValue sets the "value" field.
func (ruuo *ResourceUsageUpdateOne) SetValue(s string) *ResourceUsageUpdateOne {
	ruuo.mutation.SetValue(s)
	return ruuo
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (ruuo *ResourceUsageUpdateOne) SetNillableValue(s *string) *ResourceUsageUpdateOne {
	if s != nil {
		ruuo.SetValue(*s)
	}
	return ruuo
}

// ClearValue clears the value of the "value" field.
func (ruuo *ResourceUsageUpdateOne) ClearValue() *ResourceUsageUpdateOne {
	ruuo.mutation.ClearValue()
	return ruuo
}

// SetExecutionInfoID sets the "execution_info_id" field.
func (ruuo *ResourceUsageUpdateOne) SetExecutionInfoID(i int) *ResourceUsageUpdateOne {
	ruuo.mutation.SetExecutionInfoID(i)
	return ruuo
}

// SetNillableExecutionInfoID sets the "execution_info_id" field if the given value is not nil.
func (ruuo *ResourceUsageUpdateOne) SetNillableExecutionInfoID(i *int) *ResourceUsageUpdateOne {
	if i != nil {
		ruuo.SetExecutionInfoID(*i)
	}
	return ruuo
}

// ClearExecutionInfoID clears the value of the "execution_info_id" field.
func (ruuo *ResourceUsageUpdateOne) ClearExecutionInfoID() *ResourceUsageUpdateOne {
	ruuo.mutation.ClearExecutionInfoID()
	return ruuo
}

// SetExecutionInfo sets the "execution_info" edge to the ExectionInfo entity.
func (ruuo *ResourceUsageUpdateOne) SetExecutionInfo(e *ExectionInfo) *ResourceUsageUpdateOne {
	return ruuo.SetExecutionInfoID(e.ID)
}

// Mutation returns the ResourceUsageMutation object of the builder.
func (ruuo *ResourceUsageUpdateOne) Mutation() *ResourceUsageMutation {
	return ruuo.mutation
}

// ClearExecutionInfo clears the "execution_info" edge to the ExectionInfo entity.
func (ruuo *ResourceUsageUpdateOne) ClearExecutionInfo() *ResourceUsageUpdateOne {
	ruuo.mutation.ClearExecutionInfo()
	return ruuo
}

// Where appends a list predicates to the ResourceUsageUpdate builder.
func (ruuo *ResourceUsageUpdateOne) Where(ps ...predicate.ResourceUsage) *ResourceUsageUpdateOne {
	ruuo.mutation.Where(ps...)
	return ruuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruuo *ResourceUsageUpdateOne) Select(field string, fields ...string) *ResourceUsageUpdateOne {
	ruuo.fields = append([]string{field}, fields...)
	return ruuo
}

// Save executes the query and returns the updated ResourceUsage entity.
func (ruuo *ResourceUsageUpdateOne) Save(ctx context.Context) (*ResourceUsage, error) {
	return withHooks(ctx, ruuo.sqlSave, ruuo.mutation, ruuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruuo *ResourceUsageUpdateOne) SaveX(ctx context.Context) *ResourceUsage {
	node, err := ruuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruuo *ResourceUsageUpdateOne) Exec(ctx context.Context) error {
	_, err := ruuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruuo *ResourceUsageUpdateOne) ExecX(ctx context.Context) {
	if err := ruuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ruuo *ResourceUsageUpdateOne) sqlSave(ctx context.Context) (_node *ResourceUsage, err error) {
	_spec := sqlgraph.NewUpdateSpec(resourceusage.Table, resourceusage.Columns, sqlgraph.NewFieldSpec(resourceusage.FieldID, field.TypeInt))
	id, ok := ruuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ResourceUsage.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, resourceusage.FieldID)
		for _, f := range fields {
			if !resourceusage.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != resourceusage.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruuo.mutation.Name(); ok {
		_spec.SetField(resourceusage.FieldName, field.TypeString, value)
	}
	if ruuo.mutation.NameCleared() {
		_spec.ClearField(resourceusage.FieldName, field.TypeString)
	}
	if value, ok := ruuo.mutation.Value(); ok {
		_spec.SetField(resourceusage.FieldValue, field.TypeString, value)
	}
	if ruuo.mutation.ValueCleared() {
		_spec.ClearField(resourceusage.FieldValue, field.TypeString)
	}
	if ruuo.mutation.ExecutionInfoCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   resourceusage.ExecutionInfoTable,
			Columns: []string{resourceusage.ExecutionInfoColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exectioninfo.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruuo.mutation.ExecutionInfoIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   resourceusage.ExecutionInfoTable,
			Columns: []string{resourceusage.ExecutionInfoColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exectioninfo.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ResourceUsage{config: ruuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resourceusage.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruuo.mutation.done = true
	return _node, nil
}
