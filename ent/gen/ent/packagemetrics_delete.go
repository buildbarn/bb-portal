// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/packagemetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// PackageMetricsDelete is the builder for deleting a PackageMetrics entity.
type PackageMetricsDelete struct {
	config
	hooks    []Hook
	mutation *PackageMetricsMutation
}

// Where appends a list predicates to the PackageMetricsDelete builder.
func (pmd *PackageMetricsDelete) Where(ps ...predicate.PackageMetrics) *PackageMetricsDelete {
	pmd.mutation.Where(ps...)
	return pmd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pmd *PackageMetricsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pmd.sqlExec, pmd.mutation, pmd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pmd *PackageMetricsDelete) ExecX(ctx context.Context) int {
	n, err := pmd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pmd *PackageMetricsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(packagemetrics.Table, sqlgraph.NewFieldSpec(packagemetrics.FieldID, field.TypeInt))
	if ps := pmd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pmd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pmd.mutation.done = true
	return affected, err
}

// PackageMetricsDeleteOne is the builder for deleting a single PackageMetrics entity.
type PackageMetricsDeleteOne struct {
	pmd *PackageMetricsDelete
}

// Where appends a list predicates to the PackageMetricsDelete builder.
func (pmdo *PackageMetricsDeleteOne) Where(ps ...predicate.PackageMetrics) *PackageMetricsDeleteOne {
	pmdo.pmd.mutation.Where(ps...)
	return pmdo
}

// Exec executes the deletion query.
func (pmdo *PackageMetricsDeleteOne) Exec(ctx context.Context) error {
	n, err := pmdo.pmd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{packagemetrics.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pmdo *PackageMetricsDeleteOne) ExecX(ctx context.Context) {
	if err := pmdo.Exec(ctx); err != nil {
		panic(err)
	}
}
