// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actioncachestatistics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actiondata"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actionsummary"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/runnercount"
)

// ActionSummaryCreate is the builder for creating a ActionSummary entity.
type ActionSummaryCreate struct {
	config
	mutation *ActionSummaryMutation
	hooks    []Hook
}

// SetActionsCreated sets the "actions_created" field.
func (asc *ActionSummaryCreate) SetActionsCreated(i int64) *ActionSummaryCreate {
	asc.mutation.SetActionsCreated(i)
	return asc
}

// SetNillableActionsCreated sets the "actions_created" field if the given value is not nil.
func (asc *ActionSummaryCreate) SetNillableActionsCreated(i *int64) *ActionSummaryCreate {
	if i != nil {
		asc.SetActionsCreated(*i)
	}
	return asc
}

// SetActionsCreatedNotIncludingAspects sets the "actions_created_not_including_aspects" field.
func (asc *ActionSummaryCreate) SetActionsCreatedNotIncludingAspects(i int64) *ActionSummaryCreate {
	asc.mutation.SetActionsCreatedNotIncludingAspects(i)
	return asc
}

// SetNillableActionsCreatedNotIncludingAspects sets the "actions_created_not_including_aspects" field if the given value is not nil.
func (asc *ActionSummaryCreate) SetNillableActionsCreatedNotIncludingAspects(i *int64) *ActionSummaryCreate {
	if i != nil {
		asc.SetActionsCreatedNotIncludingAspects(*i)
	}
	return asc
}

// SetActionsExecuted sets the "actions_executed" field.
func (asc *ActionSummaryCreate) SetActionsExecuted(i int64) *ActionSummaryCreate {
	asc.mutation.SetActionsExecuted(i)
	return asc
}

// SetNillableActionsExecuted sets the "actions_executed" field if the given value is not nil.
func (asc *ActionSummaryCreate) SetNillableActionsExecuted(i *int64) *ActionSummaryCreate {
	if i != nil {
		asc.SetActionsExecuted(*i)
	}
	return asc
}

// SetRemoteCacheHits sets the "remote_cache_hits" field.
func (asc *ActionSummaryCreate) SetRemoteCacheHits(i int64) *ActionSummaryCreate {
	asc.mutation.SetRemoteCacheHits(i)
	return asc
}

// SetNillableRemoteCacheHits sets the "remote_cache_hits" field if the given value is not nil.
func (asc *ActionSummaryCreate) SetNillableRemoteCacheHits(i *int64) *ActionSummaryCreate {
	if i != nil {
		asc.SetRemoteCacheHits(*i)
	}
	return asc
}

// SetMetricsID sets the "metrics" edge to the Metrics entity by ID.
func (asc *ActionSummaryCreate) SetMetricsID(id int) *ActionSummaryCreate {
	asc.mutation.SetMetricsID(id)
	return asc
}

// SetNillableMetricsID sets the "metrics" edge to the Metrics entity by ID if the given value is not nil.
func (asc *ActionSummaryCreate) SetNillableMetricsID(id *int) *ActionSummaryCreate {
	if id != nil {
		asc = asc.SetMetricsID(*id)
	}
	return asc
}

// SetMetrics sets the "metrics" edge to the Metrics entity.
func (asc *ActionSummaryCreate) SetMetrics(m *Metrics) *ActionSummaryCreate {
	return asc.SetMetricsID(m.ID)
}

// AddActionDatumIDs adds the "action_data" edge to the ActionData entity by IDs.
func (asc *ActionSummaryCreate) AddActionDatumIDs(ids ...int) *ActionSummaryCreate {
	asc.mutation.AddActionDatumIDs(ids...)
	return asc
}

// AddActionData adds the "action_data" edges to the ActionData entity.
func (asc *ActionSummaryCreate) AddActionData(a ...*ActionData) *ActionSummaryCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return asc.AddActionDatumIDs(ids...)
}

// AddRunnerCountIDs adds the "runner_count" edge to the RunnerCount entity by IDs.
func (asc *ActionSummaryCreate) AddRunnerCountIDs(ids ...int) *ActionSummaryCreate {
	asc.mutation.AddRunnerCountIDs(ids...)
	return asc
}

// AddRunnerCount adds the "runner_count" edges to the RunnerCount entity.
func (asc *ActionSummaryCreate) AddRunnerCount(r ...*RunnerCount) *ActionSummaryCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return asc.AddRunnerCountIDs(ids...)
}

// AddActionCacheStatisticIDs adds the "action_cache_statistics" edge to the ActionCacheStatistics entity by IDs.
func (asc *ActionSummaryCreate) AddActionCacheStatisticIDs(ids ...int) *ActionSummaryCreate {
	asc.mutation.AddActionCacheStatisticIDs(ids...)
	return asc
}

// AddActionCacheStatistics adds the "action_cache_statistics" edges to the ActionCacheStatistics entity.
func (asc *ActionSummaryCreate) AddActionCacheStatistics(a ...*ActionCacheStatistics) *ActionSummaryCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return asc.AddActionCacheStatisticIDs(ids...)
}

// Mutation returns the ActionSummaryMutation object of the builder.
func (asc *ActionSummaryCreate) Mutation() *ActionSummaryMutation {
	return asc.mutation
}

// Save creates the ActionSummary in the database.
func (asc *ActionSummaryCreate) Save(ctx context.Context) (*ActionSummary, error) {
	return withHooks(ctx, asc.sqlSave, asc.mutation, asc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (asc *ActionSummaryCreate) SaveX(ctx context.Context) *ActionSummary {
	v, err := asc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (asc *ActionSummaryCreate) Exec(ctx context.Context) error {
	_, err := asc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (asc *ActionSummaryCreate) ExecX(ctx context.Context) {
	if err := asc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (asc *ActionSummaryCreate) check() error {
	return nil
}

func (asc *ActionSummaryCreate) sqlSave(ctx context.Context) (*ActionSummary, error) {
	if err := asc.check(); err != nil {
		return nil, err
	}
	_node, _spec := asc.createSpec()
	if err := sqlgraph.CreateNode(ctx, asc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	asc.mutation.id = &_node.ID
	asc.mutation.done = true
	return _node, nil
}

func (asc *ActionSummaryCreate) createSpec() (*ActionSummary, *sqlgraph.CreateSpec) {
	var (
		_node = &ActionSummary{config: asc.config}
		_spec = sqlgraph.NewCreateSpec(actionsummary.Table, sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt))
	)
	if value, ok := asc.mutation.ActionsCreated(); ok {
		_spec.SetField(actionsummary.FieldActionsCreated, field.TypeInt64, value)
		_node.ActionsCreated = value
	}
	if value, ok := asc.mutation.ActionsCreatedNotIncludingAspects(); ok {
		_spec.SetField(actionsummary.FieldActionsCreatedNotIncludingAspects, field.TypeInt64, value)
		_node.ActionsCreatedNotIncludingAspects = value
	}
	if value, ok := asc.mutation.ActionsExecuted(); ok {
		_spec.SetField(actionsummary.FieldActionsExecuted, field.TypeInt64, value)
		_node.ActionsExecuted = value
	}
	if value, ok := asc.mutation.RemoteCacheHits(); ok {
		_spec.SetField(actionsummary.FieldRemoteCacheHits, field.TypeInt64, value)
		_node.RemoteCacheHits = value
	}
	if nodes := asc.mutation.MetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   actionsummary.MetricsTable,
			Columns: []string{actionsummary.MetricsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.metrics_action_summary = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := asc.mutation.ActionDataIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   actionsummary.ActionDataTable,
			Columns: actionsummary.ActionDataPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actiondata.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := asc.mutation.RunnerCountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   actionsummary.RunnerCountTable,
			Columns: actionsummary.RunnerCountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(runnercount.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := asc.mutation.ActionCacheStatisticsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   actionsummary.ActionCacheStatisticsTable,
			Columns: actionsummary.ActionCacheStatisticsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actioncachestatistics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ActionSummaryCreateBulk is the builder for creating many ActionSummary entities in bulk.
type ActionSummaryCreateBulk struct {
	config
	err      error
	builders []*ActionSummaryCreate
}

// Save creates the ActionSummary entities in the database.
func (ascb *ActionSummaryCreateBulk) Save(ctx context.Context) ([]*ActionSummary, error) {
	if ascb.err != nil {
		return nil, ascb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ascb.builders))
	nodes := make([]*ActionSummary, len(ascb.builders))
	mutators := make([]Mutator, len(ascb.builders))
	for i := range ascb.builders {
		func(i int, root context.Context) {
			builder := ascb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ActionSummaryMutation)
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
					_, err = mutators[i+1].Mutate(root, ascb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ascb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ascb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ascb *ActionSummaryCreateBulk) SaveX(ctx context.Context) []*ActionSummary {
	v, err := ascb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ascb *ActionSummaryCreateBulk) Exec(ctx context.Context) error {
	_, err := ascb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ascb *ActionSummaryCreateBulk) ExecX(ctx context.Context) {
	if err := ascb.Exec(ctx); err != nil {
		panic(err)
	}
}