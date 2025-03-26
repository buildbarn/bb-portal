// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/cumulativemetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
)

// CumulativeMetricsCreate is the builder for creating a CumulativeMetrics entity.
type CumulativeMetricsCreate struct {
	config
	mutation *CumulativeMetricsMutation
	hooks    []Hook
}

// SetNumAnalyses sets the "num_analyses" field.
func (cmc *CumulativeMetricsCreate) SetNumAnalyses(i int32) *CumulativeMetricsCreate {
	cmc.mutation.SetNumAnalyses(i)
	return cmc
}

// SetNillableNumAnalyses sets the "num_analyses" field if the given value is not nil.
func (cmc *CumulativeMetricsCreate) SetNillableNumAnalyses(i *int32) *CumulativeMetricsCreate {
	if i != nil {
		cmc.SetNumAnalyses(*i)
	}
	return cmc
}

// SetNumBuilds sets the "num_builds" field.
func (cmc *CumulativeMetricsCreate) SetNumBuilds(i int32) *CumulativeMetricsCreate {
	cmc.mutation.SetNumBuilds(i)
	return cmc
}

// SetNillableNumBuilds sets the "num_builds" field if the given value is not nil.
func (cmc *CumulativeMetricsCreate) SetNillableNumBuilds(i *int32) *CumulativeMetricsCreate {
	if i != nil {
		cmc.SetNumBuilds(*i)
	}
	return cmc
}

// SetMetricsID sets the "metrics_id" field.
func (cmc *CumulativeMetricsCreate) SetMetricsID(i int) *CumulativeMetricsCreate {
	cmc.mutation.SetMetricsID(i)
	return cmc
}

// SetNillableMetricsID sets the "metrics_id" field if the given value is not nil.
func (cmc *CumulativeMetricsCreate) SetNillableMetricsID(i *int) *CumulativeMetricsCreate {
	if i != nil {
		cmc.SetMetricsID(*i)
	}
	return cmc
}

// SetMetrics sets the "metrics" edge to the Metrics entity.
func (cmc *CumulativeMetricsCreate) SetMetrics(m *Metrics) *CumulativeMetricsCreate {
	return cmc.SetMetricsID(m.ID)
}

// Mutation returns the CumulativeMetricsMutation object of the builder.
func (cmc *CumulativeMetricsCreate) Mutation() *CumulativeMetricsMutation {
	return cmc.mutation
}

// Save creates the CumulativeMetrics in the database.
func (cmc *CumulativeMetricsCreate) Save(ctx context.Context) (*CumulativeMetrics, error) {
	return withHooks(ctx, cmc.sqlSave, cmc.mutation, cmc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cmc *CumulativeMetricsCreate) SaveX(ctx context.Context) *CumulativeMetrics {
	v, err := cmc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cmc *CumulativeMetricsCreate) Exec(ctx context.Context) error {
	_, err := cmc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cmc *CumulativeMetricsCreate) ExecX(ctx context.Context) {
	if err := cmc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cmc *CumulativeMetricsCreate) check() error {
	return nil
}

func (cmc *CumulativeMetricsCreate) sqlSave(ctx context.Context) (*CumulativeMetrics, error) {
	if err := cmc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cmc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cmc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	cmc.mutation.id = &_node.ID
	cmc.mutation.done = true
	return _node, nil
}

func (cmc *CumulativeMetricsCreate) createSpec() (*CumulativeMetrics, *sqlgraph.CreateSpec) {
	var (
		_node = &CumulativeMetrics{config: cmc.config}
		_spec = sqlgraph.NewCreateSpec(cumulativemetrics.Table, sqlgraph.NewFieldSpec(cumulativemetrics.FieldID, field.TypeInt))
	)
	if value, ok := cmc.mutation.NumAnalyses(); ok {
		_spec.SetField(cumulativemetrics.FieldNumAnalyses, field.TypeInt32, value)
		_node.NumAnalyses = value
	}
	if value, ok := cmc.mutation.NumBuilds(); ok {
		_spec.SetField(cumulativemetrics.FieldNumBuilds, field.TypeInt32, value)
		_node.NumBuilds = value
	}
	if nodes := cmc.mutation.MetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   cumulativemetrics.MetricsTable,
			Columns: []string{cumulativemetrics.MetricsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.MetricsID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CumulativeMetricsCreateBulk is the builder for creating many CumulativeMetrics entities in bulk.
type CumulativeMetricsCreateBulk struct {
	config
	err      error
	builders []*CumulativeMetricsCreate
}

// Save creates the CumulativeMetrics entities in the database.
func (cmcb *CumulativeMetricsCreateBulk) Save(ctx context.Context) ([]*CumulativeMetrics, error) {
	if cmcb.err != nil {
		return nil, cmcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(cmcb.builders))
	nodes := make([]*CumulativeMetrics, len(cmcb.builders))
	mutators := make([]Mutator, len(cmcb.builders))
	for i := range cmcb.builders {
		func(i int, root context.Context) {
			builder := cmcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CumulativeMetricsMutation)
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
					_, err = mutators[i+1].Mutate(root, cmcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cmcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, cmcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cmcb *CumulativeMetricsCreateBulk) SaveX(ctx context.Context) []*CumulativeMetrics {
	v, err := cmcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cmcb *CumulativeMetricsCreateBulk) Exec(ctx context.Context) error {
	_, err := cmcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cmcb *CumulativeMetricsCreateBulk) ExecX(ctx context.Context) {
	if err := cmcb.Exec(ctx); err != nil {
		panic(err)
	}
}
