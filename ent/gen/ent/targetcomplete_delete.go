// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetcomplete"
)

// TargetCompleteDelete is the builder for deleting a TargetComplete entity.
type TargetCompleteDelete struct {
	config
	hooks    []Hook
	mutation *TargetCompleteMutation
}

// Where appends a list predicates to the TargetCompleteDelete builder.
func (tcd *TargetCompleteDelete) Where(ps ...predicate.TargetComplete) *TargetCompleteDelete {
	tcd.mutation.Where(ps...)
	return tcd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tcd *TargetCompleteDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tcd.sqlExec, tcd.mutation, tcd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tcd *TargetCompleteDelete) ExecX(ctx context.Context) int {
	n, err := tcd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tcd *TargetCompleteDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(targetcomplete.Table, sqlgraph.NewFieldSpec(targetcomplete.FieldID, field.TypeInt))
	if ps := tcd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tcd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tcd.mutation.done = true
	return affected, err
}

// TargetCompleteDeleteOne is the builder for deleting a single TargetComplete entity.
type TargetCompleteDeleteOne struct {
	tcd *TargetCompleteDelete
}

// Where appends a list predicates to the TargetCompleteDelete builder.
func (tcdo *TargetCompleteDeleteOne) Where(ps ...predicate.TargetComplete) *TargetCompleteDeleteOne {
	tcdo.tcd.mutation.Where(ps...)
	return tcdo
}

// Exec executes the deletion query.
func (tcdo *TargetCompleteDeleteOne) Exec(ctx context.Context) error {
	n, err := tcdo.tcd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{targetcomplete.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tcdo *TargetCompleteDeleteOne) ExecX(ctx context.Context) {
	if err := tcdo.Exec(ctx); err != nil {
		panic(err)
	}
}