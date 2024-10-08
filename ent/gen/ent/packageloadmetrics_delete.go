// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/packageloadmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// PackageLoadMetricsDelete is the builder for deleting a PackageLoadMetrics entity.
type PackageLoadMetricsDelete struct {
	config
	hooks    []Hook
	mutation *PackageLoadMetricsMutation
}

// Where appends a list predicates to the PackageLoadMetricsDelete builder.
func (plmd *PackageLoadMetricsDelete) Where(ps ...predicate.PackageLoadMetrics) *PackageLoadMetricsDelete {
	plmd.mutation.Where(ps...)
	return plmd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (plmd *PackageLoadMetricsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, plmd.sqlExec, plmd.mutation, plmd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (plmd *PackageLoadMetricsDelete) ExecX(ctx context.Context) int {
	n, err := plmd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (plmd *PackageLoadMetricsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(packageloadmetrics.Table, sqlgraph.NewFieldSpec(packageloadmetrics.FieldID, field.TypeInt))
	if ps := plmd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, plmd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	plmd.mutation.done = true
	return affected, err
}

// PackageLoadMetricsDeleteOne is the builder for deleting a single PackageLoadMetrics entity.
type PackageLoadMetricsDeleteOne struct {
	plmd *PackageLoadMetricsDelete
}

// Where appends a list predicates to the PackageLoadMetricsDelete builder.
func (plmdo *PackageLoadMetricsDeleteOne) Where(ps ...predicate.PackageLoadMetrics) *PackageLoadMetricsDeleteOne {
	plmdo.plmd.mutation.Where(ps...)
	return plmdo
}

// Exec executes the deletion query.
func (plmdo *PackageLoadMetricsDeleteOne) Exec(ctx context.Context) error {
	n, err := plmdo.plmd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{packageloadmetrics.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (plmdo *PackageLoadMetricsDeleteOne) ExecX(ctx context.Context) {
	if err := plmdo.Exec(ctx); err != nil {
		panic(err)
	}
}
