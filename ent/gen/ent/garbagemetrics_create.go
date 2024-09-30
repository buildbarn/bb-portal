// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/garbagemetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/memorymetrics"
)

// GarbageMetricsCreate is the builder for creating a GarbageMetrics entity.
type GarbageMetricsCreate struct {
	config
	mutation *GarbageMetricsMutation
	hooks    []Hook
}

// SetType sets the "type" field.
func (gmc *GarbageMetricsCreate) SetType(s string) *GarbageMetricsCreate {
	gmc.mutation.SetType(s)
	return gmc
}

// SetNillableType sets the "type" field if the given value is not nil.
func (gmc *GarbageMetricsCreate) SetNillableType(s *string) *GarbageMetricsCreate {
	if s != nil {
		gmc.SetType(*s)
	}
	return gmc
}

// SetGarbageCollected sets the "garbage_collected" field.
func (gmc *GarbageMetricsCreate) SetGarbageCollected(i int64) *GarbageMetricsCreate {
	gmc.mutation.SetGarbageCollected(i)
	return gmc
}

// SetNillableGarbageCollected sets the "garbage_collected" field if the given value is not nil.
func (gmc *GarbageMetricsCreate) SetNillableGarbageCollected(i *int64) *GarbageMetricsCreate {
	if i != nil {
		gmc.SetGarbageCollected(*i)
	}
	return gmc
}

// AddMemoryMetricIDs adds the "memory_metrics" edge to the MemoryMetrics entity by IDs.
func (gmc *GarbageMetricsCreate) AddMemoryMetricIDs(ids ...int) *GarbageMetricsCreate {
	gmc.mutation.AddMemoryMetricIDs(ids...)
	return gmc
}

// AddMemoryMetrics adds the "memory_metrics" edges to the MemoryMetrics entity.
func (gmc *GarbageMetricsCreate) AddMemoryMetrics(m ...*MemoryMetrics) *GarbageMetricsCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return gmc.AddMemoryMetricIDs(ids...)
}

// Mutation returns the GarbageMetricsMutation object of the builder.
func (gmc *GarbageMetricsCreate) Mutation() *GarbageMetricsMutation {
	return gmc.mutation
}

// Save creates the GarbageMetrics in the database.
func (gmc *GarbageMetricsCreate) Save(ctx context.Context) (*GarbageMetrics, error) {
	return withHooks(ctx, gmc.sqlSave, gmc.mutation, gmc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (gmc *GarbageMetricsCreate) SaveX(ctx context.Context) *GarbageMetrics {
	v, err := gmc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gmc *GarbageMetricsCreate) Exec(ctx context.Context) error {
	_, err := gmc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gmc *GarbageMetricsCreate) ExecX(ctx context.Context) {
	if err := gmc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gmc *GarbageMetricsCreate) check() error {
	return nil
}

func (gmc *GarbageMetricsCreate) sqlSave(ctx context.Context) (*GarbageMetrics, error) {
	if err := gmc.check(); err != nil {
		return nil, err
	}
	_node, _spec := gmc.createSpec()
	if err := sqlgraph.CreateNode(ctx, gmc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	gmc.mutation.id = &_node.ID
	gmc.mutation.done = true
	return _node, nil
}

func (gmc *GarbageMetricsCreate) createSpec() (*GarbageMetrics, *sqlgraph.CreateSpec) {
	var (
		_node = &GarbageMetrics{config: gmc.config}
		_spec = sqlgraph.NewCreateSpec(garbagemetrics.Table, sqlgraph.NewFieldSpec(garbagemetrics.FieldID, field.TypeInt))
	)
	if value, ok := gmc.mutation.GetType(); ok {
		_spec.SetField(garbagemetrics.FieldType, field.TypeString, value)
		_node.Type = value
	}
	if value, ok := gmc.mutation.GarbageCollected(); ok {
		_spec.SetField(garbagemetrics.FieldGarbageCollected, field.TypeInt64, value)
		_node.GarbageCollected = value
	}
	if nodes := gmc.mutation.MemoryMetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   garbagemetrics.MemoryMetricsTable,
			Columns: garbagemetrics.MemoryMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(memorymetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// GarbageMetricsCreateBulk is the builder for creating many GarbageMetrics entities in bulk.
type GarbageMetricsCreateBulk struct {
	config
	err      error
	builders []*GarbageMetricsCreate
}

// Save creates the GarbageMetrics entities in the database.
func (gmcb *GarbageMetricsCreateBulk) Save(ctx context.Context) ([]*GarbageMetrics, error) {
	if gmcb.err != nil {
		return nil, gmcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(gmcb.builders))
	nodes := make([]*GarbageMetrics, len(gmcb.builders))
	mutators := make([]Mutator, len(gmcb.builders))
	for i := range gmcb.builders {
		func(i int, root context.Context) {
			builder := gmcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GarbageMetricsMutation)
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
					_, err = mutators[i+1].Mutate(root, gmcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gmcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, gmcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gmcb *GarbageMetricsCreateBulk) SaveX(ctx context.Context) []*GarbageMetrics {
	v, err := gmcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gmcb *GarbageMetricsCreateBulk) Exec(ctx context.Context) error {
	_, err := gmcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gmcb *GarbageMetricsCreateBulk) ExecX(ctx context.Context) {
	if err := gmcb.Exec(ctx); err != nil {
		panic(err)
	}
}