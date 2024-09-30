// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/racestatistics"
)

// RaceStatisticsDelete is the builder for deleting a RaceStatistics entity.
type RaceStatisticsDelete struct {
	config
	hooks    []Hook
	mutation *RaceStatisticsMutation
}

// Where appends a list predicates to the RaceStatisticsDelete builder.
func (rsd *RaceStatisticsDelete) Where(ps ...predicate.RaceStatistics) *RaceStatisticsDelete {
	rsd.mutation.Where(ps...)
	return rsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (rsd *RaceStatisticsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, rsd.sqlExec, rsd.mutation, rsd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (rsd *RaceStatisticsDelete) ExecX(ctx context.Context) int {
	n, err := rsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (rsd *RaceStatisticsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(racestatistics.Table, sqlgraph.NewFieldSpec(racestatistics.FieldID, field.TypeInt))
	if ps := rsd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, rsd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	rsd.mutation.done = true
	return affected, err
}

// RaceStatisticsDeleteOne is the builder for deleting a single RaceStatistics entity.
type RaceStatisticsDeleteOne struct {
	rsd *RaceStatisticsDelete
}

// Where appends a list predicates to the RaceStatisticsDelete builder.
func (rsdo *RaceStatisticsDeleteOne) Where(ps ...predicate.RaceStatistics) *RaceStatisticsDeleteOne {
	rsdo.rsd.mutation.Where(ps...)
	return rsdo
}

// Exec executes the deletion query.
func (rsdo *RaceStatisticsDeleteOne) Exec(ctx context.Context) error {
	n, err := rsdo.rsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{racestatistics.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (rsdo *RaceStatisticsDeleteOne) ExecX(ctx context.Context) {
	if err := rsdo.Exec(ctx); err != nil {
		panic(err)
	}
}