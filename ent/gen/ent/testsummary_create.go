// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testcollection"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testfile"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testsummary"
)

// TestSummaryCreate is the builder for creating a TestSummary entity.
type TestSummaryCreate struct {
	config
	mutation *TestSummaryMutation
	hooks    []Hook
}

// SetOverallStatus sets the "overall_status" field.
func (tsc *TestSummaryCreate) SetOverallStatus(ts testsummary.OverallStatus) *TestSummaryCreate {
	tsc.mutation.SetOverallStatus(ts)
	return tsc
}

// SetNillableOverallStatus sets the "overall_status" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableOverallStatus(ts *testsummary.OverallStatus) *TestSummaryCreate {
	if ts != nil {
		tsc.SetOverallStatus(*ts)
	}
	return tsc
}

// SetTotalRunCount sets the "total_run_count" field.
func (tsc *TestSummaryCreate) SetTotalRunCount(i int32) *TestSummaryCreate {
	tsc.mutation.SetTotalRunCount(i)
	return tsc
}

// SetNillableTotalRunCount sets the "total_run_count" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableTotalRunCount(i *int32) *TestSummaryCreate {
	if i != nil {
		tsc.SetTotalRunCount(*i)
	}
	return tsc
}

// SetRunCount sets the "run_count" field.
func (tsc *TestSummaryCreate) SetRunCount(i int32) *TestSummaryCreate {
	tsc.mutation.SetRunCount(i)
	return tsc
}

// SetNillableRunCount sets the "run_count" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableRunCount(i *int32) *TestSummaryCreate {
	if i != nil {
		tsc.SetRunCount(*i)
	}
	return tsc
}

// SetAttemptCount sets the "attempt_count" field.
func (tsc *TestSummaryCreate) SetAttemptCount(i int32) *TestSummaryCreate {
	tsc.mutation.SetAttemptCount(i)
	return tsc
}

// SetNillableAttemptCount sets the "attempt_count" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableAttemptCount(i *int32) *TestSummaryCreate {
	if i != nil {
		tsc.SetAttemptCount(*i)
	}
	return tsc
}

// SetShardCount sets the "shard_count" field.
func (tsc *TestSummaryCreate) SetShardCount(i int32) *TestSummaryCreate {
	tsc.mutation.SetShardCount(i)
	return tsc
}

// SetNillableShardCount sets the "shard_count" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableShardCount(i *int32) *TestSummaryCreate {
	if i != nil {
		tsc.SetShardCount(*i)
	}
	return tsc
}

// SetTotalNumCached sets the "total_num_cached" field.
func (tsc *TestSummaryCreate) SetTotalNumCached(i int32) *TestSummaryCreate {
	tsc.mutation.SetTotalNumCached(i)
	return tsc
}

// SetNillableTotalNumCached sets the "total_num_cached" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableTotalNumCached(i *int32) *TestSummaryCreate {
	if i != nil {
		tsc.SetTotalNumCached(*i)
	}
	return tsc
}

// SetFirstStartTime sets the "first_start_time" field.
func (tsc *TestSummaryCreate) SetFirstStartTime(i int64) *TestSummaryCreate {
	tsc.mutation.SetFirstStartTime(i)
	return tsc
}

// SetNillableFirstStartTime sets the "first_start_time" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableFirstStartTime(i *int64) *TestSummaryCreate {
	if i != nil {
		tsc.SetFirstStartTime(*i)
	}
	return tsc
}

// SetLastStopTime sets the "last_stop_time" field.
func (tsc *TestSummaryCreate) SetLastStopTime(i int64) *TestSummaryCreate {
	tsc.mutation.SetLastStopTime(i)
	return tsc
}

// SetNillableLastStopTime sets the "last_stop_time" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableLastStopTime(i *int64) *TestSummaryCreate {
	if i != nil {
		tsc.SetLastStopTime(*i)
	}
	return tsc
}

// SetTotalRunDuration sets the "total_run_duration" field.
func (tsc *TestSummaryCreate) SetTotalRunDuration(i int64) *TestSummaryCreate {
	tsc.mutation.SetTotalRunDuration(i)
	return tsc
}

// SetNillableTotalRunDuration sets the "total_run_duration" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableTotalRunDuration(i *int64) *TestSummaryCreate {
	if i != nil {
		tsc.SetTotalRunDuration(*i)
	}
	return tsc
}

// SetLabel sets the "label" field.
func (tsc *TestSummaryCreate) SetLabel(s string) *TestSummaryCreate {
	tsc.mutation.SetLabel(s)
	return tsc
}

// SetNillableLabel sets the "label" field if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableLabel(s *string) *TestSummaryCreate {
	if s != nil {
		tsc.SetLabel(*s)
	}
	return tsc
}

// SetTestCollectionID sets the "test_collection" edge to the TestCollection entity by ID.
func (tsc *TestSummaryCreate) SetTestCollectionID(id int) *TestSummaryCreate {
	tsc.mutation.SetTestCollectionID(id)
	return tsc
}

// SetNillableTestCollectionID sets the "test_collection" edge to the TestCollection entity by ID if the given value is not nil.
func (tsc *TestSummaryCreate) SetNillableTestCollectionID(id *int) *TestSummaryCreate {
	if id != nil {
		tsc = tsc.SetTestCollectionID(*id)
	}
	return tsc
}

// SetTestCollection sets the "test_collection" edge to the TestCollection entity.
func (tsc *TestSummaryCreate) SetTestCollection(t *TestCollection) *TestSummaryCreate {
	return tsc.SetTestCollectionID(t.ID)
}

// AddPassedIDs adds the "passed" edge to the TestFile entity by IDs.
func (tsc *TestSummaryCreate) AddPassedIDs(ids ...int) *TestSummaryCreate {
	tsc.mutation.AddPassedIDs(ids...)
	return tsc
}

// AddPassed adds the "passed" edges to the TestFile entity.
func (tsc *TestSummaryCreate) AddPassed(t ...*TestFile) *TestSummaryCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tsc.AddPassedIDs(ids...)
}

// AddFailedIDs adds the "failed" edge to the TestFile entity by IDs.
func (tsc *TestSummaryCreate) AddFailedIDs(ids ...int) *TestSummaryCreate {
	tsc.mutation.AddFailedIDs(ids...)
	return tsc
}

// AddFailed adds the "failed" edges to the TestFile entity.
func (tsc *TestSummaryCreate) AddFailed(t ...*TestFile) *TestSummaryCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tsc.AddFailedIDs(ids...)
}

// Mutation returns the TestSummaryMutation object of the builder.
func (tsc *TestSummaryCreate) Mutation() *TestSummaryMutation {
	return tsc.mutation
}

// Save creates the TestSummary in the database.
func (tsc *TestSummaryCreate) Save(ctx context.Context) (*TestSummary, error) {
	tsc.defaults()
	return withHooks(ctx, tsc.sqlSave, tsc.mutation, tsc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tsc *TestSummaryCreate) SaveX(ctx context.Context) *TestSummary {
	v, err := tsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tsc *TestSummaryCreate) Exec(ctx context.Context) error {
	_, err := tsc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tsc *TestSummaryCreate) ExecX(ctx context.Context) {
	if err := tsc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tsc *TestSummaryCreate) defaults() {
	if _, ok := tsc.mutation.OverallStatus(); !ok {
		v := testsummary.DefaultOverallStatus
		tsc.mutation.SetOverallStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tsc *TestSummaryCreate) check() error {
	if v, ok := tsc.mutation.OverallStatus(); ok {
		if err := testsummary.OverallStatusValidator(v); err != nil {
			return &ValidationError{Name: "overall_status", err: fmt.Errorf(`ent: validator failed for field "TestSummary.overall_status": %w`, err)}
		}
	}
	return nil
}

func (tsc *TestSummaryCreate) sqlSave(ctx context.Context) (*TestSummary, error) {
	if err := tsc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tsc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tsc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tsc.mutation.id = &_node.ID
	tsc.mutation.done = true
	return _node, nil
}

func (tsc *TestSummaryCreate) createSpec() (*TestSummary, *sqlgraph.CreateSpec) {
	var (
		_node = &TestSummary{config: tsc.config}
		_spec = sqlgraph.NewCreateSpec(testsummary.Table, sqlgraph.NewFieldSpec(testsummary.FieldID, field.TypeInt))
	)
	if value, ok := tsc.mutation.OverallStatus(); ok {
		_spec.SetField(testsummary.FieldOverallStatus, field.TypeEnum, value)
		_node.OverallStatus = value
	}
	if value, ok := tsc.mutation.TotalRunCount(); ok {
		_spec.SetField(testsummary.FieldTotalRunCount, field.TypeInt32, value)
		_node.TotalRunCount = value
	}
	if value, ok := tsc.mutation.RunCount(); ok {
		_spec.SetField(testsummary.FieldRunCount, field.TypeInt32, value)
		_node.RunCount = value
	}
	if value, ok := tsc.mutation.AttemptCount(); ok {
		_spec.SetField(testsummary.FieldAttemptCount, field.TypeInt32, value)
		_node.AttemptCount = value
	}
	if value, ok := tsc.mutation.ShardCount(); ok {
		_spec.SetField(testsummary.FieldShardCount, field.TypeInt32, value)
		_node.ShardCount = value
	}
	if value, ok := tsc.mutation.TotalNumCached(); ok {
		_spec.SetField(testsummary.FieldTotalNumCached, field.TypeInt32, value)
		_node.TotalNumCached = value
	}
	if value, ok := tsc.mutation.FirstStartTime(); ok {
		_spec.SetField(testsummary.FieldFirstStartTime, field.TypeInt64, value)
		_node.FirstStartTime = value
	}
	if value, ok := tsc.mutation.LastStopTime(); ok {
		_spec.SetField(testsummary.FieldLastStopTime, field.TypeInt64, value)
		_node.LastStopTime = value
	}
	if value, ok := tsc.mutation.TotalRunDuration(); ok {
		_spec.SetField(testsummary.FieldTotalRunDuration, field.TypeInt64, value)
		_node.TotalRunDuration = value
	}
	if value, ok := tsc.mutation.Label(); ok {
		_spec.SetField(testsummary.FieldLabel, field.TypeString, value)
		_node.Label = value
	}
	if nodes := tsc.mutation.TestCollectionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   testsummary.TestCollectionTable,
			Columns: []string{testsummary.TestCollectionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testcollection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.test_collection_test_summary = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tsc.mutation.PassedIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   testsummary.PassedTable,
			Columns: []string{testsummary.PassedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testfile.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tsc.mutation.FailedIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   testsummary.FailedTable,
			Columns: []string{testsummary.FailedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testfile.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TestSummaryCreateBulk is the builder for creating many TestSummary entities in bulk.
type TestSummaryCreateBulk struct {
	config
	err      error
	builders []*TestSummaryCreate
}

// Save creates the TestSummary entities in the database.
func (tscb *TestSummaryCreateBulk) Save(ctx context.Context) ([]*TestSummary, error) {
	if tscb.err != nil {
		return nil, tscb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tscb.builders))
	nodes := make([]*TestSummary, len(tscb.builders))
	mutators := make([]Mutator, len(tscb.builders))
	for i := range tscb.builders {
		func(i int, root context.Context) {
			builder := tscb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TestSummaryMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tscb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tscb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tscb *TestSummaryCreateBulk) SaveX(ctx context.Context) []*TestSummary {
	v, err := tscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tscb *TestSummaryCreateBulk) Exec(ctx context.Context) error {
	_, err := tscb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tscb *TestSummaryCreateBulk) ExecX(ctx context.Context) {
	if err := tscb.Exec(ctx); err != nil {
		panic(err)
	}
}
