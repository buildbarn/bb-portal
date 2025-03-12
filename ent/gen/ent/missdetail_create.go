// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actioncachestatistics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/missdetail"
)

// MissDetailCreate is the builder for creating a MissDetail entity.
type MissDetailCreate struct {
	config
	mutation *MissDetailMutation
	hooks    []Hook
}

// SetReason sets the "reason" field.
func (mdc *MissDetailCreate) SetReason(m missdetail.Reason) *MissDetailCreate {
	mdc.mutation.SetReason(m)
	return mdc
}

// SetNillableReason sets the "reason" field if the given value is not nil.
func (mdc *MissDetailCreate) SetNillableReason(m *missdetail.Reason) *MissDetailCreate {
	if m != nil {
		mdc.SetReason(*m)
	}
	return mdc
}

// SetCount sets the "count" field.
func (mdc *MissDetailCreate) SetCount(i int32) *MissDetailCreate {
	mdc.mutation.SetCount(i)
	return mdc
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (mdc *MissDetailCreate) SetNillableCount(i *int32) *MissDetailCreate {
	if i != nil {
		mdc.SetCount(*i)
	}
	return mdc
}

// SetActionCacheStatisticsID sets the "action_cache_statistics_id" field.
func (mdc *MissDetailCreate) SetActionCacheStatisticsID(i int) *MissDetailCreate {
	mdc.mutation.SetActionCacheStatisticsID(i)
	return mdc
}

// SetNillableActionCacheStatisticsID sets the "action_cache_statistics_id" field if the given value is not nil.
func (mdc *MissDetailCreate) SetNillableActionCacheStatisticsID(i *int) *MissDetailCreate {
	if i != nil {
		mdc.SetActionCacheStatisticsID(*i)
	}
	return mdc
}

// SetActionCacheStatistics sets the "action_cache_statistics" edge to the ActionCacheStatistics entity.
func (mdc *MissDetailCreate) SetActionCacheStatistics(a *ActionCacheStatistics) *MissDetailCreate {
	return mdc.SetActionCacheStatisticsID(a.ID)
}

// Mutation returns the MissDetailMutation object of the builder.
func (mdc *MissDetailCreate) Mutation() *MissDetailMutation {
	return mdc.mutation
}

// Save creates the MissDetail in the database.
func (mdc *MissDetailCreate) Save(ctx context.Context) (*MissDetail, error) {
	mdc.defaults()
	return withHooks(ctx, mdc.sqlSave, mdc.mutation, mdc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mdc *MissDetailCreate) SaveX(ctx context.Context) *MissDetail {
	v, err := mdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mdc *MissDetailCreate) Exec(ctx context.Context) error {
	_, err := mdc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mdc *MissDetailCreate) ExecX(ctx context.Context) {
	if err := mdc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mdc *MissDetailCreate) defaults() {
	if _, ok := mdc.mutation.Reason(); !ok {
		v := missdetail.DefaultReason
		mdc.mutation.SetReason(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mdc *MissDetailCreate) check() error {
	if v, ok := mdc.mutation.Reason(); ok {
		if err := missdetail.ReasonValidator(v); err != nil {
			return &ValidationError{Name: "reason", err: fmt.Errorf(`ent: validator failed for field "MissDetail.reason": %w`, err)}
		}
	}
	return nil
}

func (mdc *MissDetailCreate) sqlSave(ctx context.Context) (*MissDetail, error) {
	if err := mdc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mdc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mdc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	mdc.mutation.id = &_node.ID
	mdc.mutation.done = true
	return _node, nil
}

func (mdc *MissDetailCreate) createSpec() (*MissDetail, *sqlgraph.CreateSpec) {
	var (
		_node = &MissDetail{config: mdc.config}
		_spec = sqlgraph.NewCreateSpec(missdetail.Table, sqlgraph.NewFieldSpec(missdetail.FieldID, field.TypeInt))
	)
	if value, ok := mdc.mutation.Reason(); ok {
		_spec.SetField(missdetail.FieldReason, field.TypeEnum, value)
		_node.Reason = value
	}
	if value, ok := mdc.mutation.Count(); ok {
		_spec.SetField(missdetail.FieldCount, field.TypeInt32, value)
		_node.Count = value
	}
	if nodes := mdc.mutation.ActionCacheStatisticsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   missdetail.ActionCacheStatisticsTable,
			Columns: []string{missdetail.ActionCacheStatisticsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(actioncachestatistics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ActionCacheStatisticsID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MissDetailCreateBulk is the builder for creating many MissDetail entities in bulk.
type MissDetailCreateBulk struct {
	config
	err      error
	builders []*MissDetailCreate
}

// Save creates the MissDetail entities in the database.
func (mdcb *MissDetailCreateBulk) Save(ctx context.Context) ([]*MissDetail, error) {
	if mdcb.err != nil {
		return nil, mdcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mdcb.builders))
	nodes := make([]*MissDetail, len(mdcb.builders))
	mutators := make([]Mutator, len(mdcb.builders))
	for i := range mdcb.builders {
		func(i int, root context.Context) {
			builder := mdcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MissDetailMutation)
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
					_, err = mutators[i+1].Mutate(root, mdcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mdcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, mdcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mdcb *MissDetailCreateBulk) SaveX(ctx context.Context) []*MissDetail {
	v, err := mdcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mdcb *MissDetailCreateBulk) Exec(ctx context.Context) error {
	_, err := mdcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mdcb *MissDetailCreateBulk) ExecX(ctx context.Context) {
	if err := mdcb.Exec(ctx); err != nil {
		panic(err)
	}
}
