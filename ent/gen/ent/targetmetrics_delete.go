// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetmetrics"
)

// TargetMetricsDelete is the builder for deleting a TargetMetrics entity.
type TargetMetricsDelete struct {
	config
	hooks    []Hook
	mutation *TargetMetricsMutation
}

// Where appends a list predicates to the TargetMetricsDelete builder.
func (tmd *TargetMetricsDelete) Where(ps ...predicate.TargetMetrics) *TargetMetricsDelete {
	tmd.mutation.Where(ps...)
	return tmd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tmd *TargetMetricsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tmd.sqlExec, tmd.mutation, tmd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tmd *TargetMetricsDelete) ExecX(ctx context.Context) int {
	n, err := tmd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tmd *TargetMetricsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(targetmetrics.Table, sqlgraph.NewFieldSpec(targetmetrics.FieldID, field.TypeInt))
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

// TargetMetricsDeleteOne is the builder for deleting a single TargetMetrics entity.
type TargetMetricsDeleteOne struct {
	tmd *TargetMetricsDelete
}

// Where appends a list predicates to the TargetMetricsDelete builder.
func (tmdo *TargetMetricsDeleteOne) Where(ps ...predicate.TargetMetrics) *TargetMetricsDeleteOne {
	tmdo.tmd.mutation.Where(ps...)
	return tmdo
}

// Exec executes the deletion query.
func (tmdo *TargetMetricsDeleteOne) Exec(ctx context.Context) error {
	n, err := tmdo.tmd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{targetmetrics.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tmdo *TargetMetricsDeleteOne) ExecX(ctx context.Context) {
	if err := tmdo.Exec(ctx); err != nil {
		panic(err)
	}
}
