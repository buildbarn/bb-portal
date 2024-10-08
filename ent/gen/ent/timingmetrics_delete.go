// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingmetrics"
)

// TimingMetricsDelete is the builder for deleting a TimingMetrics entity.
type TimingMetricsDelete struct {
	config
	hooks    []Hook
	mutation *TimingMetricsMutation
}

// Where appends a list predicates to the TimingMetricsDelete builder.
func (tmd *TimingMetricsDelete) Where(ps ...predicate.TimingMetrics) *TimingMetricsDelete {
	tmd.mutation.Where(ps...)
	return tmd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tmd *TimingMetricsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tmd.sqlExec, tmd.mutation, tmd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tmd *TimingMetricsDelete) ExecX(ctx context.Context) int {
	n, err := tmd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tmd *TimingMetricsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(timingmetrics.Table, sqlgraph.NewFieldSpec(timingmetrics.FieldID, field.TypeInt))
	if ps := tmd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tmd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tmd.mutation.done = true
	return affected, err
}

// TimingMetricsDeleteOne is the builder for deleting a single TimingMetrics entity.
type TimingMetricsDeleteOne struct {
	tmd *TimingMetricsDelete
}

// Where appends a list predicates to the TimingMetricsDelete builder.
func (tmdo *TimingMetricsDeleteOne) Where(ps ...predicate.TimingMetrics) *TimingMetricsDeleteOne {
	tmdo.tmd.mutation.Where(ps...)
	return tmdo
}

// Exec executes the deletion query.
func (tmdo *TimingMetricsDeleteOne) Exec(ctx context.Context) error {
	n, err := tmdo.tmd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{timingmetrics.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tmdo *TimingMetricsDeleteOne) ExecX(ctx context.Context) {
	if err := tmdo.Exec(ctx); err != nil {
		panic(err)
	}
}
