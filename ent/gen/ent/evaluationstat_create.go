// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/buildgraphmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/evaluationstat"
)

// EvaluationStatCreate is the builder for creating a EvaluationStat entity.
type EvaluationStatCreate struct {
	config
	mutation *EvaluationStatMutation
	hooks    []Hook
}

// SetSkyfunctionName sets the "skyfunction_name" field.
func (esc *EvaluationStatCreate) SetSkyfunctionName(s string) *EvaluationStatCreate {
	esc.mutation.SetSkyfunctionName(s)
	return esc
}

// SetNillableSkyfunctionName sets the "skyfunction_name" field if the given value is not nil.
func (esc *EvaluationStatCreate) SetNillableSkyfunctionName(s *string) *EvaluationStatCreate {
	if s != nil {
		esc.SetSkyfunctionName(*s)
	}
	return esc
}

// SetCount sets the "count" field.
func (esc *EvaluationStatCreate) SetCount(i int64) *EvaluationStatCreate {
	esc.mutation.SetCount(i)
	return esc
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (esc *EvaluationStatCreate) SetNillableCount(i *int64) *EvaluationStatCreate {
	if i != nil {
		esc.SetCount(*i)
	}
	return esc
}

// AddBuildGraphMetricIDs adds the "build_graph_metrics" edge to the BuildGraphMetrics entity by IDs.
func (esc *EvaluationStatCreate) AddBuildGraphMetricIDs(ids ...int) *EvaluationStatCreate {
	esc.mutation.AddBuildGraphMetricIDs(ids...)
	return esc
}

// AddBuildGraphMetrics adds the "build_graph_metrics" edges to the BuildGraphMetrics entity.
func (esc *EvaluationStatCreate) AddBuildGraphMetrics(b ...*BuildGraphMetrics) *EvaluationStatCreate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return esc.AddBuildGraphMetricIDs(ids...)
}

// Mutation returns the EvaluationStatMutation object of the builder.
func (esc *EvaluationStatCreate) Mutation() *EvaluationStatMutation {
	return esc.mutation
}

// Save creates the EvaluationStat in the database.
func (esc *EvaluationStatCreate) Save(ctx context.Context) (*EvaluationStat, error) {
	return withHooks(ctx, esc.sqlSave, esc.mutation, esc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (esc *EvaluationStatCreate) SaveX(ctx context.Context) *EvaluationStat {
	v, err := esc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (esc *EvaluationStatCreate) Exec(ctx context.Context) error {
	_, err := esc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (esc *EvaluationStatCreate) ExecX(ctx context.Context) {
	if err := esc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (esc *EvaluationStatCreate) check() error {
	return nil
}

func (esc *EvaluationStatCreate) sqlSave(ctx context.Context) (*EvaluationStat, error) {
	if err := esc.check(); err != nil {
		return nil, err
	}
	_node, _spec := esc.createSpec()
	if err := sqlgraph.CreateNode(ctx, esc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	esc.mutation.id = &_node.ID
	esc.mutation.done = true
	return _node, nil
}

func (esc *EvaluationStatCreate) createSpec() (*EvaluationStat, *sqlgraph.CreateSpec) {
	var (
		_node = &EvaluationStat{config: esc.config}
		_spec = sqlgraph.NewCreateSpec(evaluationstat.Table, sqlgraph.NewFieldSpec(evaluationstat.FieldID, field.TypeInt))
	)
	if value, ok := esc.mutation.SkyfunctionName(); ok {
		_spec.SetField(evaluationstat.FieldSkyfunctionName, field.TypeString, value)
		_node.SkyfunctionName = value
	}
	if value, ok := esc.mutation.Count(); ok {
		_spec.SetField(evaluationstat.FieldCount, field.TypeInt64, value)
		_node.Count = value
	}
	if nodes := esc.mutation.BuildGraphMetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   evaluationstat.BuildGraphMetricsTable,
			Columns: evaluationstat.BuildGraphMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildgraphmetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// EvaluationStatCreateBulk is the builder for creating many EvaluationStat entities in bulk.
type EvaluationStatCreateBulk struct {
	config
	err      error
	builders []*EvaluationStatCreate
}

// Save creates the EvaluationStat entities in the database.
func (escb *EvaluationStatCreateBulk) Save(ctx context.Context) ([]*EvaluationStat, error) {
	if escb.err != nil {
		return nil, escb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(escb.builders))
	nodes := make([]*EvaluationStat, len(escb.builders))
	mutators := make([]Mutator, len(escb.builders))
	for i := range escb.builders {
		func(i int, root context.Context) {
			builder := escb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EvaluationStatMutation)
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
					_, err = mutators[i+1].Mutate(root, escb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, escb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, escb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (escb *EvaluationStatCreateBulk) SaveX(ctx context.Context) []*EvaluationStat {
	v, err := escb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (escb *EvaluationStatCreateBulk) Exec(ctx context.Context) error {
	_, err := escb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (escb *EvaluationStatCreateBulk) ExecX(ctx context.Context) {
	if err := escb.Exec(ctx); err != nil {
		panic(err)
	}
}
