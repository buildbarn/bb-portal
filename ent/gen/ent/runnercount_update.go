// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actionsummary"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/runnercount"
)

// RunnerCountUpdate is the builder for updating RunnerCount entities.
type RunnerCountUpdate struct {
	config
	hooks    []Hook
	mutation *RunnerCountMutation
}

// Where appends a list predicates to the RunnerCountUpdate builder.
func (rcu *RunnerCountUpdate) Where(ps ...predicate.RunnerCount) *RunnerCountUpdate {
	rcu.mutation.Where(ps...)
	return rcu
}

// SetName sets the "name" field.
func (rcu *RunnerCountUpdate) SetName(s string) *RunnerCountUpdate {
	rcu.mutation.SetName(s)
	return rcu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (rcu *RunnerCountUpdate) SetNillableName(s *string) *RunnerCountUpdate {
	if s != nil {
		rcu.SetName(*s)
	}
	return rcu
}

// ClearName clears the value of the "name" field.
func (rcu *RunnerCountUpdate) ClearName() *RunnerCountUpdate {
	rcu.mutation.ClearName()
	return rcu
}

// SetExecKind sets the "exec_kind" field.
func (rcu *RunnerCountUpdate) SetExecKind(s string) *RunnerCountUpdate {
	rcu.mutation.SetExecKind(s)
	return rcu
}

// SetNillableExecKind sets the "exec_kind" field if the given value is not nil.
func (rcu *RunnerCountUpdate) SetNillableExecKind(s *string) *RunnerCountUpdate {
	if s != nil {
		rcu.SetExecKind(*s)
	}
	return rcu
}

// ClearExecKind clears the value of the "exec_kind" field.
func (rcu *RunnerCountUpdate) ClearExecKind() *RunnerCountUpdate {
	rcu.mutation.ClearExecKind()
	return rcu
}

// SetActionsExecuted sets the "actions_executed" field.
func (rcu *RunnerCountUpdate) SetActionsExecuted(i int64) *RunnerCountUpdate {
	rcu.mutation.ResetActionsExecuted()
	rcu.mutation.SetActionsExecuted(i)
	return rcu
}

// SetNillableActionsExecuted sets the "actions_executed" field if the given value is not nil.
func (rcu *RunnerCountUpdate) SetNillableActionsExecuted(i *int64) *RunnerCountUpdate {
	if i != nil {
		rcu.SetActionsExecuted(*i)
	}
	return rcu
}

// AddActionsExecuted adds i to the "actions_executed" field.
func (rcu *RunnerCountUpdate) AddActionsExecuted(i int64) *RunnerCountUpdate {
	rcu.mutation.AddActionsExecuted(i)
	return rcu
}

// ClearActionsExecuted clears the value of the "actions_executed" field.
func (rcu *RunnerCountUpdate) ClearActionsExecuted() *RunnerCountUpdate {
	rcu.mutation.ClearActionsExecuted()
	return rcu
}

// AddActionSummaryIDs adds the "action_summary" edge to the ActionSummary entity by IDs.
func (rcu *RunnerCountUpdate) AddActionSummaryIDs(ids ...int) *RunnerCountUpdate {
	rcu.mutation.AddActionSummaryIDs(ids...)
	return rcu
}

// AddActionSummary adds the "action_summary" edges to the ActionSummary entity.
func (rcu *RunnerCountUpdate) AddActionSummary(a ...*ActionSummary) *RunnerCountUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return rcu.AddActionSummaryIDs(ids...)
}

// Mutation returns the RunnerCountMutation object of the builder.
func (rcu *RunnerCountUpdate) Mutation() *RunnerCountMutation {
	return rcu.mutation
}

// ClearActionSummary clears all "action_summary" edges to the ActionSummary entity.
func (rcu *RunnerCountUpdate) ClearActionSummary() *RunnerCountUpdate {
	rcu.mutation.ClearActionSummary()
	return rcu
}

// RemoveActionSummaryIDs removes the "action_summary" edge to ActionSummary entities by IDs.
func (rcu *RunnerCountUpdate) RemoveActionSummaryIDs(ids ...int) *RunnerCountUpdate {
	rcu.mutation.RemoveActionSummaryIDs(ids...)
	return rcu
}

// RemoveActionSummary removes "action_summary" edges to ActionSummary entities.
func (rcu *RunnerCountUpdate) RemoveActionSummary(a ...*ActionSummary) *RunnerCountUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return rcu.RemoveActionSummaryIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (rcu *RunnerCountUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, rcu.sqlSave, rcu.mutation, rcu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (rcu *RunnerCountUpdate) SaveX(ctx context.Context) int {
	affected, err := rcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (rcu *RunnerCountUpdate) Exec(ctx context.Context) error {
	_, err := rcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcu *RunnerCountUpdate) ExecX(ctx context.Context) {
	if err := rcu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (rcu *RunnerCountUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(runnercount.Table, runnercount.Columns, sqlgraph.NewFieldSpec(runnercount.FieldID, field.TypeInt))
	if ps := rcu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := rcu.mutation.Name(); ok {
		_spec.SetField(runnercount.FieldName, field.TypeString, value)
	}
	if rcu.mutation.NameCleared() {
		_spec.ClearField(runnercount.FieldName, field.TypeString)
	}
	if value, ok := rcu.mutation.ExecKind(); ok {
		_spec.SetField(runnercount.FieldExecKind, field.TypeString, value)
	}
	if rcu.mutation.ExecKindCleared() {
		_spec.ClearField(runnercount.FieldExecKind, field.TypeString)
	}
	if value, ok := rcu.mutation.ActionsExecuted(); ok {
		_spec.SetField(runnercount.FieldActionsExecuted, field.TypeInt64, value)
	}
	if value, ok := rcu.mutation.AddedActionsExecuted(); ok {
		_spec.AddField(runnercount.FieldActionsExecuted, field.TypeInt64, value)
	}
	if rcu.mutation.ActionsExecutedCleared() {
		_spec.ClearField(runnercount.FieldActionsExecuted, field.TypeInt64)
	}
	if rcu.mutation.ActionSummaryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   runnercount.ActionSummaryTable,
			Columns: runnercount.ActionSummaryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rcu.mutation.RemovedActionSummaryIDs(); len(nodes) > 0 && !rcu.mutation.ActionSummaryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   runnercount.ActionSummaryTable,
			Columns: runnercount.ActionSummaryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rcu.mutation.ActionSummaryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   runnercount.ActionSummaryTable,
			Columns: runnercount.ActionSummaryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, rcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{runnercount.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	rcu.mutation.done = true
	return n, nil
}

// RunnerCountUpdateOne is the builder for updating a single RunnerCount entity.
type RunnerCountUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *RunnerCountMutation
}

// SetName sets the "name" field.
func (rcuo *RunnerCountUpdateOne) SetName(s string) *RunnerCountUpdateOne {
	rcuo.mutation.SetName(s)
	return rcuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (rcuo *RunnerCountUpdateOne) SetNillableName(s *string) *RunnerCountUpdateOne {
	if s != nil {
		rcuo.SetName(*s)
	}
	return rcuo
}

// ClearName clears the value of the "name" field.
func (rcuo *RunnerCountUpdateOne) ClearName() *RunnerCountUpdateOne {
	rcuo.mutation.ClearName()
	return rcuo
}

// SetExecKind sets the "exec_kind" field.
func (rcuo *RunnerCountUpdateOne) SetExecKind(s string) *RunnerCountUpdateOne {
	rcuo.mutation.SetExecKind(s)
	return rcuo
}

// SetNillableExecKind sets the "exec_kind" field if the given value is not nil.
func (rcuo *RunnerCountUpdateOne) SetNillableExecKind(s *string) *RunnerCountUpdateOne {
	if s != nil {
		rcuo.SetExecKind(*s)
	}
	return rcuo
}

// ClearExecKind clears the value of the "exec_kind" field.
func (rcuo *RunnerCountUpdateOne) ClearExecKind() *RunnerCountUpdateOne {
	rcuo.mutation.ClearExecKind()
	return rcuo
}

// SetActionsExecuted sets the "actions_executed" field.
func (rcuo *RunnerCountUpdateOne) SetActionsExecuted(i int64) *RunnerCountUpdateOne {
	rcuo.mutation.ResetActionsExecuted()
	rcuo.mutation.SetActionsExecuted(i)
	return rcuo
}

// SetNillableActionsExecuted sets the "actions_executed" field if the given value is not nil.
func (rcuo *RunnerCountUpdateOne) SetNillableActionsExecuted(i *int64) *RunnerCountUpdateOne {
	if i != nil {
		rcuo.SetActionsExecuted(*i)
	}
	return rcuo
}

// AddActionsExecuted adds i to the "actions_executed" field.
func (rcuo *RunnerCountUpdateOne) AddActionsExecuted(i int64) *RunnerCountUpdateOne {
	rcuo.mutation.AddActionsExecuted(i)
	return rcuo
}

// ClearActionsExecuted clears the value of the "actions_executed" field.
func (rcuo *RunnerCountUpdateOne) ClearActionsExecuted() *RunnerCountUpdateOne {
	rcuo.mutation.ClearActionsExecuted()
	return rcuo
}

// AddActionSummaryIDs adds the "action_summary" edge to the ActionSummary entity by IDs.
func (rcuo *RunnerCountUpdateOne) AddActionSummaryIDs(ids ...int) *RunnerCountUpdateOne {
	rcuo.mutation.AddActionSummaryIDs(ids...)
	return rcuo
}

// AddActionSummary adds the "action_summary" edges to the ActionSummary entity.
func (rcuo *RunnerCountUpdateOne) AddActionSummary(a ...*ActionSummary) *RunnerCountUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return rcuo.AddActionSummaryIDs(ids...)
}

// Mutation returns the RunnerCountMutation object of the builder.
func (rcuo *RunnerCountUpdateOne) Mutation() *RunnerCountMutation {
	return rcuo.mutation
}

// ClearActionSummary clears all "action_summary" edges to the ActionSummary entity.
func (rcuo *RunnerCountUpdateOne) ClearActionSummary() *RunnerCountUpdateOne {
	rcuo.mutation.ClearActionSummary()
	return rcuo
}

// RemoveActionSummaryIDs removes the "action_summary" edge to ActionSummary entities by IDs.
func (rcuo *RunnerCountUpdateOne) RemoveActionSummaryIDs(ids ...int) *RunnerCountUpdateOne {
	rcuo.mutation.RemoveActionSummaryIDs(ids...)
	return rcuo
}

// RemoveActionSummary removes "action_summary" edges to ActionSummary entities.
func (rcuo *RunnerCountUpdateOne) RemoveActionSummary(a ...*ActionSummary) *RunnerCountUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return rcuo.RemoveActionSummaryIDs(ids...)
}

// Where appends a list predicates to the RunnerCountUpdate builder.
func (rcuo *RunnerCountUpdateOne) Where(ps ...predicate.RunnerCount) *RunnerCountUpdateOne {
	rcuo.mutation.Where(ps...)
	return rcuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (rcuo *RunnerCountUpdateOne) Select(field string, fields ...string) *RunnerCountUpdateOne {
	rcuo.fields = append([]string{field}, fields...)
	return rcuo
}

// Save executes the query and returns the updated RunnerCount entity.
func (rcuo *RunnerCountUpdateOne) Save(ctx context.Context) (*RunnerCount, error) {
	return withHooks(ctx, rcuo.sqlSave, rcuo.mutation, rcuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (rcuo *RunnerCountUpdateOne) SaveX(ctx context.Context) *RunnerCount {
	node, err := rcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (rcuo *RunnerCountUpdateOne) Exec(ctx context.Context) error {
	_, err := rcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcuo *RunnerCountUpdateOne) ExecX(ctx context.Context) {
	if err := rcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (rcuo *RunnerCountUpdateOne) sqlSave(ctx context.Context) (_node *RunnerCount, err error) {
	_spec := sqlgraph.NewUpdateSpec(runnercount.Table, runnercount.Columns, sqlgraph.NewFieldSpec(runnercount.FieldID, field.TypeInt))
	id, ok := rcuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "RunnerCount.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := rcuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, runnercount.FieldID)
		for _, f := range fields {
			if !runnercount.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != runnercount.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := rcuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := rcuo.mutation.Name(); ok {
		_spec.SetField(runnercount.FieldName, field.TypeString, value)
	}
	if rcuo.mutation.NameCleared() {
		_spec.ClearField(runnercount.FieldName, field.TypeString)
	}
	if value, ok := rcuo.mutation.ExecKind(); ok {
		_spec.SetField(runnercount.FieldExecKind, field.TypeString, value)
	}
	if rcuo.mutation.ExecKindCleared() {
		_spec.ClearField(runnercount.FieldExecKind, field.TypeString)
	}
	if value, ok := rcuo.mutation.ActionsExecuted(); ok {
		_spec.SetField(runnercount.FieldActionsExecuted, field.TypeInt64, value)
	}
	if value, ok := rcuo.mutation.AddedActionsExecuted(); ok {
		_spec.AddField(runnercount.FieldActionsExecuted, field.TypeInt64, value)
	}
	if rcuo.mutation.ActionsExecutedCleared() {
		_spec.ClearField(runnercount.FieldActionsExecuted, field.TypeInt64)
	}
	if rcuo.mutation.ActionSummaryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   runnercount.ActionSummaryTable,
			Columns: runnercount.ActionSummaryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rcuo.mutation.RemovedActionSummaryIDs(); len(nodes) > 0 && !rcuo.mutation.ActionSummaryCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   runnercount.ActionSummaryTable,
			Columns: runnercount.ActionSummaryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rcuo.mutation.ActionSummaryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   runnercount.ActionSummaryTable,
			Columns: runnercount.ActionSummaryPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &RunnerCount{config: rcuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, rcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{runnercount.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	rcuo.mutation.done = true
	return _node, nil
}
