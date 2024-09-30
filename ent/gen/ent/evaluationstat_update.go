// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildgraphmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/evaluationstat"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// EvaluationStatUpdate is the builder for updating EvaluationStat entities.
type EvaluationStatUpdate struct {
	config
	hooks    []Hook
	mutation *EvaluationStatMutation
}

// Where appends a list predicates to the EvaluationStatUpdate builder.
func (esu *EvaluationStatUpdate) Where(ps ...predicate.EvaluationStat) *EvaluationStatUpdate {
	esu.mutation.Where(ps...)
	return esu
}

// SetSkyfunctionName sets the "skyfunction_name" field.
func (esu *EvaluationStatUpdate) SetSkyfunctionName(s string) *EvaluationStatUpdate {
	esu.mutation.SetSkyfunctionName(s)
	return esu
}

// SetNillableSkyfunctionName sets the "skyfunction_name" field if the given value is not nil.
func (esu *EvaluationStatUpdate) SetNillableSkyfunctionName(s *string) *EvaluationStatUpdate {
	if s != nil {
		esu.SetSkyfunctionName(*s)
	}
	return esu
}

// ClearSkyfunctionName clears the value of the "skyfunction_name" field.
func (esu *EvaluationStatUpdate) ClearSkyfunctionName() *EvaluationStatUpdate {
	esu.mutation.ClearSkyfunctionName()
	return esu
}

// SetCount sets the "count" field.
func (esu *EvaluationStatUpdate) SetCount(i int64) *EvaluationStatUpdate {
	esu.mutation.ResetCount()
	esu.mutation.SetCount(i)
	return esu
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (esu *EvaluationStatUpdate) SetNillableCount(i *int64) *EvaluationStatUpdate {
	if i != nil {
		esu.SetCount(*i)
	}
	return esu
}

// AddCount adds i to the "count" field.
func (esu *EvaluationStatUpdate) AddCount(i int64) *EvaluationStatUpdate {
	esu.mutation.AddCount(i)
	return esu
}

// ClearCount clears the value of the "count" field.
func (esu *EvaluationStatUpdate) ClearCount() *EvaluationStatUpdate {
	esu.mutation.ClearCount()
	return esu
}

// AddBuildGraphMetricIDs adds the "build_graph_metrics" edge to the BuildGraphMetrics entity by IDs.
func (esu *EvaluationStatUpdate) AddBuildGraphMetricIDs(ids ...int) *EvaluationStatUpdate {
	esu.mutation.AddBuildGraphMetricIDs(ids...)
	return esu
}

// AddBuildGraphMetrics adds the "build_graph_metrics" edges to the BuildGraphMetrics entity.
func (esu *EvaluationStatUpdate) AddBuildGraphMetrics(b ...*BuildGraphMetrics) *EvaluationStatUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return esu.AddBuildGraphMetricIDs(ids...)
}

// Mutation returns the EvaluationStatMutation object of the builder.
func (esu *EvaluationStatUpdate) Mutation() *EvaluationStatMutation {
	return esu.mutation
}

// ClearBuildGraphMetrics clears all "build_graph_metrics" edges to the BuildGraphMetrics entity.
func (esu *EvaluationStatUpdate) ClearBuildGraphMetrics() *EvaluationStatUpdate {
	esu.mutation.ClearBuildGraphMetrics()
	return esu
}

// RemoveBuildGraphMetricIDs removes the "build_graph_metrics" edge to BuildGraphMetrics entities by IDs.
func (esu *EvaluationStatUpdate) RemoveBuildGraphMetricIDs(ids ...int) *EvaluationStatUpdate {
	esu.mutation.RemoveBuildGraphMetricIDs(ids...)
	return esu
}

// RemoveBuildGraphMetrics removes "build_graph_metrics" edges to BuildGraphMetrics entities.
func (esu *EvaluationStatUpdate) RemoveBuildGraphMetrics(b ...*BuildGraphMetrics) *EvaluationStatUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return esu.RemoveBuildGraphMetricIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (esu *EvaluationStatUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, esu.sqlSave, esu.mutation, esu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (esu *EvaluationStatUpdate) SaveX(ctx context.Context) int {
	affected, err := esu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (esu *EvaluationStatUpdate) Exec(ctx context.Context) error {
	_, err := esu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (esu *EvaluationStatUpdate) ExecX(ctx context.Context) {
	if err := esu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (esu *EvaluationStatUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(evaluationstat.Table, evaluationstat.Columns, sqlgraph.NewFieldSpec(evaluationstat.FieldID, field.TypeInt))
	if ps := esu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := esu.mutation.SkyfunctionName(); ok {
		_spec.SetField(evaluationstat.FieldSkyfunctionName, field.TypeString, value)
	}
	if esu.mutation.SkyfunctionNameCleared() {
		_spec.ClearField(evaluationstat.FieldSkyfunctionName, field.TypeString)
	}
	if value, ok := esu.mutation.Count(); ok {
		_spec.SetField(evaluationstat.FieldCount, field.TypeInt64, value)
	}
	if value, ok := esu.mutation.AddedCount(); ok {
		_spec.AddField(evaluationstat.FieldCount, field.TypeInt64, value)
	}
	if esu.mutation.CountCleared() {
		_spec.ClearField(evaluationstat.FieldCount, field.TypeInt64)
	}
	if esu.mutation.BuildGraphMetricsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := esu.mutation.RemovedBuildGraphMetricsIDs(); len(nodes) > 0 && !esu.mutation.BuildGraphMetricsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := esu.mutation.BuildGraphMetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, esu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{evaluationstat.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	esu.mutation.done = true
	return n, nil
}

// EvaluationStatUpdateOne is the builder for updating a single EvaluationStat entity.
type EvaluationStatUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EvaluationStatMutation
}

// SetSkyfunctionName sets the "skyfunction_name" field.
func (esuo *EvaluationStatUpdateOne) SetSkyfunctionName(s string) *EvaluationStatUpdateOne {
	esuo.mutation.SetSkyfunctionName(s)
	return esuo
}

// SetNillableSkyfunctionName sets the "skyfunction_name" field if the given value is not nil.
func (esuo *EvaluationStatUpdateOne) SetNillableSkyfunctionName(s *string) *EvaluationStatUpdateOne {
	if s != nil {
		esuo.SetSkyfunctionName(*s)
	}
	return esuo
}

// ClearSkyfunctionName clears the value of the "skyfunction_name" field.
func (esuo *EvaluationStatUpdateOne) ClearSkyfunctionName() *EvaluationStatUpdateOne {
	esuo.mutation.ClearSkyfunctionName()
	return esuo
}

// SetCount sets the "count" field.
func (esuo *EvaluationStatUpdateOne) SetCount(i int64) *EvaluationStatUpdateOne {
	esuo.mutation.ResetCount()
	esuo.mutation.SetCount(i)
	return esuo
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (esuo *EvaluationStatUpdateOne) SetNillableCount(i *int64) *EvaluationStatUpdateOne {
	if i != nil {
		esuo.SetCount(*i)
	}
	return esuo
}

// AddCount adds i to the "count" field.
func (esuo *EvaluationStatUpdateOne) AddCount(i int64) *EvaluationStatUpdateOne {
	esuo.mutation.AddCount(i)
	return esuo
}

// ClearCount clears the value of the "count" field.
func (esuo *EvaluationStatUpdateOne) ClearCount() *EvaluationStatUpdateOne {
	esuo.mutation.ClearCount()
	return esuo
}

// AddBuildGraphMetricIDs adds the "build_graph_metrics" edge to the BuildGraphMetrics entity by IDs.
func (esuo *EvaluationStatUpdateOne) AddBuildGraphMetricIDs(ids ...int) *EvaluationStatUpdateOne {
	esuo.mutation.AddBuildGraphMetricIDs(ids...)
	return esuo
}

// AddBuildGraphMetrics adds the "build_graph_metrics" edges to the BuildGraphMetrics entity.
func (esuo *EvaluationStatUpdateOne) AddBuildGraphMetrics(b ...*BuildGraphMetrics) *EvaluationStatUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return esuo.AddBuildGraphMetricIDs(ids...)
}

// Mutation returns the EvaluationStatMutation object of the builder.
func (esuo *EvaluationStatUpdateOne) Mutation() *EvaluationStatMutation {
	return esuo.mutation
}

// ClearBuildGraphMetrics clears all "build_graph_metrics" edges to the BuildGraphMetrics entity.
func (esuo *EvaluationStatUpdateOne) ClearBuildGraphMetrics() *EvaluationStatUpdateOne {
	esuo.mutation.ClearBuildGraphMetrics()
	return esuo
}

// RemoveBuildGraphMetricIDs removes the "build_graph_metrics" edge to BuildGraphMetrics entities by IDs.
func (esuo *EvaluationStatUpdateOne) RemoveBuildGraphMetricIDs(ids ...int) *EvaluationStatUpdateOne {
	esuo.mutation.RemoveBuildGraphMetricIDs(ids...)
	return esuo
}

// RemoveBuildGraphMetrics removes "build_graph_metrics" edges to BuildGraphMetrics entities.
func (esuo *EvaluationStatUpdateOne) RemoveBuildGraphMetrics(b ...*BuildGraphMetrics) *EvaluationStatUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return esuo.RemoveBuildGraphMetricIDs(ids...)
}

// Where appends a list predicates to the EvaluationStatUpdate builder.
func (esuo *EvaluationStatUpdateOne) Where(ps ...predicate.EvaluationStat) *EvaluationStatUpdateOne {
	esuo.mutation.Where(ps...)
	return esuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (esuo *EvaluationStatUpdateOne) Select(field string, fields ...string) *EvaluationStatUpdateOne {
	esuo.fields = append([]string{field}, fields...)
	return esuo
}

// Save executes the query and returns the updated EvaluationStat entity.
func (esuo *EvaluationStatUpdateOne) Save(ctx context.Context) (*EvaluationStat, error) {
	return withHooks(ctx, esuo.sqlSave, esuo.mutation, esuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (esuo *EvaluationStatUpdateOne) SaveX(ctx context.Context) *EvaluationStat {
	node, err := esuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (esuo *EvaluationStatUpdateOne) Exec(ctx context.Context) error {
	_, err := esuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (esuo *EvaluationStatUpdateOne) ExecX(ctx context.Context) {
	if err := esuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (esuo *EvaluationStatUpdateOne) sqlSave(ctx context.Context) (_node *EvaluationStat, err error) {
	_spec := sqlgraph.NewUpdateSpec(evaluationstat.Table, evaluationstat.Columns, sqlgraph.NewFieldSpec(evaluationstat.FieldID, field.TypeInt))
	id, ok := esuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "EvaluationStat.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := esuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, evaluationstat.FieldID)
		for _, f := range fields {
			if !evaluationstat.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != evaluationstat.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := esuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := esuo.mutation.SkyfunctionName(); ok {
		_spec.SetField(evaluationstat.FieldSkyfunctionName, field.TypeString, value)
	}
	if esuo.mutation.SkyfunctionNameCleared() {
		_spec.ClearField(evaluationstat.FieldSkyfunctionName, field.TypeString)
	}
	if value, ok := esuo.mutation.Count(); ok {
		_spec.SetField(evaluationstat.FieldCount, field.TypeInt64, value)
	}
	if value, ok := esuo.mutation.AddedCount(); ok {
		_spec.AddField(evaluationstat.FieldCount, field.TypeInt64, value)
	}
	if esuo.mutation.CountCleared() {
		_spec.ClearField(evaluationstat.FieldCount, field.TypeInt64)
	}
	if esuo.mutation.BuildGraphMetricsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := esuo.mutation.RemovedBuildGraphMetricsIDs(); len(nodes) > 0 && !esuo.mutation.BuildGraphMetricsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := esuo.mutation.BuildGraphMetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &EvaluationStat{config: esuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, esuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{evaluationstat.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	esuo.mutation.done = true
	return _node, nil
}
