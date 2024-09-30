// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actionsummary"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ActionSummaryDelete is the builder for deleting a ActionSummary entity.
type ActionSummaryDelete struct {
	config
	hooks    []Hook
	mutation *ActionSummaryMutation
}

// Where appends a list predicates to the ActionSummaryDelete builder.
func (asd *ActionSummaryDelete) Where(ps ...predicate.ActionSummary) *ActionSummaryDelete {
	asd.mutation.Where(ps...)
	return asd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (asd *ActionSummaryDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, asd.sqlExec, asd.mutation, asd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (asd *ActionSummaryDelete) ExecX(ctx context.Context) int {
	n, err := asd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (asd *ActionSummaryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(actionsummary.Table, sqlgraph.NewFieldSpec(actionsummary.FieldID, field.TypeInt))
	if ps := asd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, asd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	asd.mutation.done = true
	return affected, err
}

// ActionSummaryDeleteOne is the builder for deleting a single ActionSummary entity.
type ActionSummaryDeleteOne struct {
	asd *ActionSummaryDelete
}

// Where appends a list predicates to the ActionSummaryDelete builder.
func (asdo *ActionSummaryDeleteOne) Where(ps ...predicate.ActionSummary) *ActionSummaryDeleteOne {
	asdo.asd.mutation.Where(ps...)
	return asdo
}

// Exec executes the deletion query.
func (asdo *ActionSummaryDeleteOne) Exec(ctx context.Context) error {
	n, err := asdo.asd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{actionsummary.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (asdo *ActionSummaryDeleteOne) ExecX(ctx context.Context) {
	if err := asdo.Exec(ctx); err != nil {
		panic(err)
	}
}