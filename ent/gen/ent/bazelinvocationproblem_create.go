// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocationproblem"
)

// BazelInvocationProblemCreate is the builder for creating a BazelInvocationProblem entity.
type BazelInvocationProblemCreate struct {
	config
	mutation *BazelInvocationProblemMutation
	hooks    []Hook
}

// SetProblemType sets the "problem_type" field.
func (bipc *BazelInvocationProblemCreate) SetProblemType(s string) *BazelInvocationProblemCreate {
	bipc.mutation.SetProblemType(s)
	return bipc
}

// SetLabel sets the "label" field.
func (bipc *BazelInvocationProblemCreate) SetLabel(s string) *BazelInvocationProblemCreate {
	bipc.mutation.SetLabel(s)
	return bipc
}

// SetBepEvents sets the "bep_events" field.
func (bipc *BazelInvocationProblemCreate) SetBepEvents(jm json.RawMessage) *BazelInvocationProblemCreate {
	bipc.mutation.SetBepEvents(jm)
	return bipc
}

// SetBazelInvocationID sets the "bazel_invocation_id" field.
func (bipc *BazelInvocationProblemCreate) SetBazelInvocationID(i int) *BazelInvocationProblemCreate {
	bipc.mutation.SetBazelInvocationID(i)
	return bipc
}

// SetNillableBazelInvocationID sets the "bazel_invocation_id" field if the given value is not nil.
func (bipc *BazelInvocationProblemCreate) SetNillableBazelInvocationID(i *int) *BazelInvocationProblemCreate {
	if i != nil {
		bipc.SetBazelInvocationID(*i)
	}
	return bipc
}

// SetBazelInvocation sets the "bazel_invocation" edge to the BazelInvocation entity.
func (bipc *BazelInvocationProblemCreate) SetBazelInvocation(b *BazelInvocation) *BazelInvocationProblemCreate {
	return bipc.SetBazelInvocationID(b.ID)
}

// Mutation returns the BazelInvocationProblemMutation object of the builder.
func (bipc *BazelInvocationProblemCreate) Mutation() *BazelInvocationProblemMutation {
	return bipc.mutation
}

// Save creates the BazelInvocationProblem in the database.
func (bipc *BazelInvocationProblemCreate) Save(ctx context.Context) (*BazelInvocationProblem, error) {
	return withHooks(ctx, bipc.sqlSave, bipc.mutation, bipc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (bipc *BazelInvocationProblemCreate) SaveX(ctx context.Context) *BazelInvocationProblem {
	v, err := bipc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bipc *BazelInvocationProblemCreate) Exec(ctx context.Context) error {
	_, err := bipc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bipc *BazelInvocationProblemCreate) ExecX(ctx context.Context) {
	if err := bipc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bipc *BazelInvocationProblemCreate) check() error {
	if _, ok := bipc.mutation.ProblemType(); !ok {
		return &ValidationError{Name: "problem_type", err: errors.New(`ent: missing required field "BazelInvocationProblem.problem_type"`)}
	}
	if _, ok := bipc.mutation.Label(); !ok {
		return &ValidationError{Name: "label", err: errors.New(`ent: missing required field "BazelInvocationProblem.label"`)}
	}
	if _, ok := bipc.mutation.BepEvents(); !ok {
		return &ValidationError{Name: "bep_events", err: errors.New(`ent: missing required field "BazelInvocationProblem.bep_events"`)}
	}
	return nil
}

func (bipc *BazelInvocationProblemCreate) sqlSave(ctx context.Context) (*BazelInvocationProblem, error) {
	if err := bipc.check(); err != nil {
		return nil, err
	}
	_node, _spec := bipc.createSpec()
	if err := sqlgraph.CreateNode(ctx, bipc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	bipc.mutation.id = &_node.ID
	bipc.mutation.done = true
	return _node, nil
}

func (bipc *BazelInvocationProblemCreate) createSpec() (*BazelInvocationProblem, *sqlgraph.CreateSpec) {
	var (
		_node = &BazelInvocationProblem{config: bipc.config}
		_spec = sqlgraph.NewCreateSpec(bazelinvocationproblem.Table, sqlgraph.NewFieldSpec(bazelinvocationproblem.FieldID, field.TypeInt))
	)
	if value, ok := bipc.mutation.ProblemType(); ok {
		_spec.SetField(bazelinvocationproblem.FieldProblemType, field.TypeString, value)
		_node.ProblemType = value
	}
	if value, ok := bipc.mutation.Label(); ok {
		_spec.SetField(bazelinvocationproblem.FieldLabel, field.TypeString, value)
		_node.Label = value
	}
	if value, ok := bipc.mutation.BepEvents(); ok {
		_spec.SetField(bazelinvocationproblem.FieldBepEvents, field.TypeJSON, value)
		_node.BepEvents = value
	}
	if nodes := bipc.mutation.BazelInvocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bazelinvocationproblem.BazelInvocationTable,
			Columns: []string{bazelinvocationproblem.BazelInvocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bazelinvocation.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.BazelInvocationID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// BazelInvocationProblemCreateBulk is the builder for creating many BazelInvocationProblem entities in bulk.
type BazelInvocationProblemCreateBulk struct {
	config
	err      error
	builders []*BazelInvocationProblemCreate
}

// Save creates the BazelInvocationProblem entities in the database.
func (bipcb *BazelInvocationProblemCreateBulk) Save(ctx context.Context) ([]*BazelInvocationProblem, error) {
	if bipcb.err != nil {
		return nil, bipcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(bipcb.builders))
	nodes := make([]*BazelInvocationProblem, len(bipcb.builders))
	mutators := make([]Mutator, len(bipcb.builders))
	for i := range bipcb.builders {
		func(i int, root context.Context) {
			builder := bipcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*BazelInvocationProblemMutation)
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
					_, err = mutators[i+1].Mutate(root, bipcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, bipcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, bipcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (bipcb *BazelInvocationProblemCreateBulk) SaveX(ctx context.Context) []*BazelInvocationProblem {
	v, err := bipcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bipcb *BazelInvocationProblemCreateBulk) Exec(ctx context.Context) error {
	_, err := bipcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bipcb *BazelInvocationProblemCreateBulk) ExecX(ctx context.Context) {
	if err := bipcb.Exec(ctx); err != nil {
		panic(err)
	}
}
