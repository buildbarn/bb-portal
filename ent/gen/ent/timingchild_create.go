// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingbreakdown"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingchild"
)

// TimingChildCreate is the builder for creating a TimingChild entity.
type TimingChildCreate struct {
	config
	mutation *TimingChildMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (tcc *TimingChildCreate) SetName(s string) *TimingChildCreate {
	tcc.mutation.SetName(s)
	return tcc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tcc *TimingChildCreate) SetNillableName(s *string) *TimingChildCreate {
	if s != nil {
		tcc.SetName(*s)
	}
	return tcc
}

// SetTime sets the "time" field.
func (tcc *TimingChildCreate) SetTime(s string) *TimingChildCreate {
	tcc.mutation.SetTime(s)
	return tcc
}

// SetNillableTime sets the "time" field if the given value is not nil.
func (tcc *TimingChildCreate) SetNillableTime(s *string) *TimingChildCreate {
	if s != nil {
		tcc.SetTime(*s)
	}
	return tcc
}

// SetTimingBreakdownID sets the "timing_breakdown" edge to the TimingBreakdown entity by ID.
func (tcc *TimingChildCreate) SetTimingBreakdownID(id int) *TimingChildCreate {
	tcc.mutation.SetTimingBreakdownID(id)
	return tcc
}

// SetNillableTimingBreakdownID sets the "timing_breakdown" edge to the TimingBreakdown entity by ID if the given value is not nil.
func (tcc *TimingChildCreate) SetNillableTimingBreakdownID(id *int) *TimingChildCreate {
	if id != nil {
		tcc = tcc.SetTimingBreakdownID(*id)
	}
	return tcc
}

// SetTimingBreakdown sets the "timing_breakdown" edge to the TimingBreakdown entity.
func (tcc *TimingChildCreate) SetTimingBreakdown(t *TimingBreakdown) *TimingChildCreate {
	return tcc.SetTimingBreakdownID(t.ID)
}

// Mutation returns the TimingChildMutation object of the builder.
func (tcc *TimingChildCreate) Mutation() *TimingChildMutation {
	return tcc.mutation
}

// Save creates the TimingChild in the database.
func (tcc *TimingChildCreate) Save(ctx context.Context) (*TimingChild, error) {
	return withHooks(ctx, tcc.sqlSave, tcc.mutation, tcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tcc *TimingChildCreate) SaveX(ctx context.Context) *TimingChild {
	v, err := tcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcc *TimingChildCreate) Exec(ctx context.Context) error {
	_, err := tcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcc *TimingChildCreate) ExecX(ctx context.Context) {
	if err := tcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcc *TimingChildCreate) check() error {
	return nil
}

func (tcc *TimingChildCreate) sqlSave(ctx context.Context) (*TimingChild, error) {
	if err := tcc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tcc.mutation.id = &_node.ID
	tcc.mutation.done = true
	return _node, nil
}

func (tcc *TimingChildCreate) createSpec() (*TimingChild, *sqlgraph.CreateSpec) {
	var (
		_node = &TimingChild{config: tcc.config}
		_spec = sqlgraph.NewCreateSpec(timingchild.Table, sqlgraph.NewFieldSpec(timingchild.FieldID, field.TypeInt))
	)
	if value, ok := tcc.mutation.Name(); ok {
		_spec.SetField(timingchild.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := tcc.mutation.Time(); ok {
		_spec.SetField(timingchild.FieldTime, field.TypeString, value)
		_node.Time = value
	}
	if nodes := tcc.mutation.TimingBreakdownIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   timingchild.TimingBreakdownTable,
			Columns: []string{timingchild.TimingBreakdownColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timingbreakdown.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.timing_breakdown_child = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TimingChildCreateBulk is the builder for creating many TimingChild entities in bulk.
type TimingChildCreateBulk struct {
	config
	err      error
	builders []*TimingChildCreate
}

// Save creates the TimingChild entities in the database.
func (tccb *TimingChildCreateBulk) Save(ctx context.Context) ([]*TimingChild, error) {
	if tccb.err != nil {
		return nil, tccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tccb.builders))
	nodes := make([]*TimingChild, len(tccb.builders))
	mutators := make([]Mutator, len(tccb.builders))
	for i := range tccb.builders {
		func(i int, root context.Context) {
			builder := tccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TimingChildMutation)
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
					_, err = mutators[i+1].Mutate(root, tccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, tccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tccb *TimingChildCreateBulk) SaveX(ctx context.Context) []*TimingChild {
	v, err := tccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tccb *TimingChildCreateBulk) Exec(ctx context.Context) error {
	_, err := tccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tccb *TimingChildCreateBulk) ExecX(ctx context.Context) {
	if err := tccb.Exec(ctx); err != nil {
		panic(err)
	}
}
