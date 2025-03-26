// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/dynamicexecutionmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/racestatistics"
)

// DynamicExecutionMetricsUpdate is the builder for updating DynamicExecutionMetrics entities.
type DynamicExecutionMetricsUpdate struct {
	config
	hooks    []Hook
	mutation *DynamicExecutionMetricsMutation
}

// Where appends a list predicates to the DynamicExecutionMetricsUpdate builder.
func (demu *DynamicExecutionMetricsUpdate) Where(ps ...predicate.DynamicExecutionMetrics) *DynamicExecutionMetricsUpdate {
	demu.mutation.Where(ps...)
	return demu
}

// SetMetricsID sets the "metrics_id" field.
func (demu *DynamicExecutionMetricsUpdate) SetMetricsID(i int) *DynamicExecutionMetricsUpdate {
	demu.mutation.SetMetricsID(i)
	return demu
}

// SetNillableMetricsID sets the "metrics_id" field if the given value is not nil.
func (demu *DynamicExecutionMetricsUpdate) SetNillableMetricsID(i *int) *DynamicExecutionMetricsUpdate {
	if i != nil {
		demu.SetMetricsID(*i)
	}
	return demu
}

// ClearMetricsID clears the value of the "metrics_id" field.
func (demu *DynamicExecutionMetricsUpdate) ClearMetricsID() *DynamicExecutionMetricsUpdate {
	demu.mutation.ClearMetricsID()
	return demu
}

// SetMetrics sets the "metrics" edge to the Metrics entity.
func (demu *DynamicExecutionMetricsUpdate) SetMetrics(m *Metrics) *DynamicExecutionMetricsUpdate {
	return demu.SetMetricsID(m.ID)
}

// AddRaceStatisticIDs adds the "race_statistics" edge to the RaceStatistics entity by IDs.
func (demu *DynamicExecutionMetricsUpdate) AddRaceStatisticIDs(ids ...int) *DynamicExecutionMetricsUpdate {
	demu.mutation.AddRaceStatisticIDs(ids...)
	return demu
}

// AddRaceStatistics adds the "race_statistics" edges to the RaceStatistics entity.
func (demu *DynamicExecutionMetricsUpdate) AddRaceStatistics(r ...*RaceStatistics) *DynamicExecutionMetricsUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return demu.AddRaceStatisticIDs(ids...)
}

// Mutation returns the DynamicExecutionMetricsMutation object of the builder.
func (demu *DynamicExecutionMetricsUpdate) Mutation() *DynamicExecutionMetricsMutation {
	return demu.mutation
}

// ClearMetrics clears the "metrics" edge to the Metrics entity.
func (demu *DynamicExecutionMetricsUpdate) ClearMetrics() *DynamicExecutionMetricsUpdate {
	demu.mutation.ClearMetrics()
	return demu
}

// ClearRaceStatistics clears all "race_statistics" edges to the RaceStatistics entity.
func (demu *DynamicExecutionMetricsUpdate) ClearRaceStatistics() *DynamicExecutionMetricsUpdate {
	demu.mutation.ClearRaceStatistics()
	return demu
}

// RemoveRaceStatisticIDs removes the "race_statistics" edge to RaceStatistics entities by IDs.
func (demu *DynamicExecutionMetricsUpdate) RemoveRaceStatisticIDs(ids ...int) *DynamicExecutionMetricsUpdate {
	demu.mutation.RemoveRaceStatisticIDs(ids...)
	return demu
}

// RemoveRaceStatistics removes "race_statistics" edges to RaceStatistics entities.
func (demu *DynamicExecutionMetricsUpdate) RemoveRaceStatistics(r ...*RaceStatistics) *DynamicExecutionMetricsUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return demu.RemoveRaceStatisticIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (demu *DynamicExecutionMetricsUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, demu.sqlSave, demu.mutation, demu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (demu *DynamicExecutionMetricsUpdate) SaveX(ctx context.Context) int {
	affected, err := demu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (demu *DynamicExecutionMetricsUpdate) Exec(ctx context.Context) error {
	_, err := demu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (demu *DynamicExecutionMetricsUpdate) ExecX(ctx context.Context) {
	if err := demu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (demu *DynamicExecutionMetricsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(dynamicexecutionmetrics.Table, dynamicexecutionmetrics.Columns, sqlgraph.NewFieldSpec(dynamicexecutionmetrics.FieldID, field.TypeInt))
	if ps := demu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if demu.mutation.MetricsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   dynamicexecutionmetrics.MetricsTable,
			Columns: []string{dynamicexecutionmetrics.MetricsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := demu.mutation.MetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   dynamicexecutionmetrics.MetricsTable,
			Columns: []string{dynamicexecutionmetrics.MetricsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if demu.mutation.RaceStatisticsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   dynamicexecutionmetrics.RaceStatisticsTable,
			Columns: []string{dynamicexecutionmetrics.RaceStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := demu.mutation.RemovedRaceStatisticsIDs(); len(nodes) > 0 && !demu.mutation.RaceStatisticsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   dynamicexecutionmetrics.RaceStatisticsTable,
			Columns: []string{dynamicexecutionmetrics.RaceStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := demu.mutation.RaceStatisticsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   dynamicexecutionmetrics.RaceStatisticsTable,
			Columns: []string{dynamicexecutionmetrics.RaceStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, demu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{dynamicexecutionmetrics.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	demu.mutation.done = true
	return n, nil
}

// DynamicExecutionMetricsUpdateOne is the builder for updating a single DynamicExecutionMetrics entity.
type DynamicExecutionMetricsUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *DynamicExecutionMetricsMutation
}

// SetMetricsID sets the "metrics_id" field.
func (demuo *DynamicExecutionMetricsUpdateOne) SetMetricsID(i int) *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.SetMetricsID(i)
	return demuo
}

// SetNillableMetricsID sets the "metrics_id" field if the given value is not nil.
func (demuo *DynamicExecutionMetricsUpdateOne) SetNillableMetricsID(i *int) *DynamicExecutionMetricsUpdateOne {
	if i != nil {
		demuo.SetMetricsID(*i)
	}
	return demuo
}

// ClearMetricsID clears the value of the "metrics_id" field.
func (demuo *DynamicExecutionMetricsUpdateOne) ClearMetricsID() *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.ClearMetricsID()
	return demuo
}

// SetMetrics sets the "metrics" edge to the Metrics entity.
func (demuo *DynamicExecutionMetricsUpdateOne) SetMetrics(m *Metrics) *DynamicExecutionMetricsUpdateOne {
	return demuo.SetMetricsID(m.ID)
}

// AddRaceStatisticIDs adds the "race_statistics" edge to the RaceStatistics entity by IDs.
func (demuo *DynamicExecutionMetricsUpdateOne) AddRaceStatisticIDs(ids ...int) *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.AddRaceStatisticIDs(ids...)
	return demuo
}

// AddRaceStatistics adds the "race_statistics" edges to the RaceStatistics entity.
func (demuo *DynamicExecutionMetricsUpdateOne) AddRaceStatistics(r ...*RaceStatistics) *DynamicExecutionMetricsUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return demuo.AddRaceStatisticIDs(ids...)
}

// Mutation returns the DynamicExecutionMetricsMutation object of the builder.
func (demuo *DynamicExecutionMetricsUpdateOne) Mutation() *DynamicExecutionMetricsMutation {
	return demuo.mutation
}

// ClearMetrics clears the "metrics" edge to the Metrics entity.
func (demuo *DynamicExecutionMetricsUpdateOne) ClearMetrics() *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.ClearMetrics()
	return demuo
}

// ClearRaceStatistics clears all "race_statistics" edges to the RaceStatistics entity.
func (demuo *DynamicExecutionMetricsUpdateOne) ClearRaceStatistics() *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.ClearRaceStatistics()
	return demuo
}

// RemoveRaceStatisticIDs removes the "race_statistics" edge to RaceStatistics entities by IDs.
func (demuo *DynamicExecutionMetricsUpdateOne) RemoveRaceStatisticIDs(ids ...int) *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.RemoveRaceStatisticIDs(ids...)
	return demuo
}

// RemoveRaceStatistics removes "race_statistics" edges to RaceStatistics entities.
func (demuo *DynamicExecutionMetricsUpdateOne) RemoveRaceStatistics(r ...*RaceStatistics) *DynamicExecutionMetricsUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return demuo.RemoveRaceStatisticIDs(ids...)
}

// Where appends a list predicates to the DynamicExecutionMetricsUpdate builder.
func (demuo *DynamicExecutionMetricsUpdateOne) Where(ps ...predicate.DynamicExecutionMetrics) *DynamicExecutionMetricsUpdateOne {
	demuo.mutation.Where(ps...)
	return demuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (demuo *DynamicExecutionMetricsUpdateOne) Select(field string, fields ...string) *DynamicExecutionMetricsUpdateOne {
	demuo.fields = append([]string{field}, fields...)
	return demuo
}

// Save executes the query and returns the updated DynamicExecutionMetrics entity.
func (demuo *DynamicExecutionMetricsUpdateOne) Save(ctx context.Context) (*DynamicExecutionMetrics, error) {
	return withHooks(ctx, demuo.sqlSave, demuo.mutation, demuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (demuo *DynamicExecutionMetricsUpdateOne) SaveX(ctx context.Context) *DynamicExecutionMetrics {
	node, err := demuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (demuo *DynamicExecutionMetricsUpdateOne) Exec(ctx context.Context) error {
	_, err := demuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (demuo *DynamicExecutionMetricsUpdateOne) ExecX(ctx context.Context) {
	if err := demuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (demuo *DynamicExecutionMetricsUpdateOne) sqlSave(ctx context.Context) (_node *DynamicExecutionMetrics, err error) {
	_spec := sqlgraph.NewUpdateSpec(dynamicexecutionmetrics.Table, dynamicexecutionmetrics.Columns, sqlgraph.NewFieldSpec(dynamicexecutionmetrics.FieldID, field.TypeInt))
	id, ok := demuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "DynamicExecutionMetrics.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := demuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, dynamicexecutionmetrics.FieldID)
		for _, f := range fields {
			if !dynamicexecutionmetrics.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != dynamicexecutionmetrics.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := demuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if demuo.mutation.MetricsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   dynamicexecutionmetrics.MetricsTable,
			Columns: []string{dynamicexecutionmetrics.MetricsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := demuo.mutation.MetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   dynamicexecutionmetrics.MetricsTable,
			Columns: []string{dynamicexecutionmetrics.MetricsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if demuo.mutation.RaceStatisticsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   dynamicexecutionmetrics.RaceStatisticsTable,
			Columns: []string{dynamicexecutionmetrics.RaceStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := demuo.mutation.RemovedRaceStatisticsIDs(); len(nodes) > 0 && !demuo.mutation.RaceStatisticsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   dynamicexecutionmetrics.RaceStatisticsTable,
			Columns: []string{dynamicexecutionmetrics.RaceStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := demuo.mutation.RaceStatisticsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   dynamicexecutionmetrics.RaceStatisticsTable,
			Columns: []string{dynamicexecutionmetrics.RaceStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &DynamicExecutionMetrics{config: demuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, demuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{dynamicexecutionmetrics.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	demuo.mutation.done = true
	return _node, nil
}
