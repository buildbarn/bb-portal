// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/artifactmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/filesmetric"
)

// FilesMetricCreate is the builder for creating a FilesMetric entity.
type FilesMetricCreate struct {
	config
	mutation *FilesMetricMutation
	hooks    []Hook
}

// SetSizeInBytes sets the "size_in_bytes" field.
func (fmc *FilesMetricCreate) SetSizeInBytes(i int64) *FilesMetricCreate {
	fmc.mutation.SetSizeInBytes(i)
	return fmc
}

// SetNillableSizeInBytes sets the "size_in_bytes" field if the given value is not nil.
func (fmc *FilesMetricCreate) SetNillableSizeInBytes(i *int64) *FilesMetricCreate {
	if i != nil {
		fmc.SetSizeInBytes(*i)
	}
	return fmc
}

// SetCount sets the "count" field.
func (fmc *FilesMetricCreate) SetCount(i int32) *FilesMetricCreate {
	fmc.mutation.SetCount(i)
	return fmc
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (fmc *FilesMetricCreate) SetNillableCount(i *int32) *FilesMetricCreate {
	if i != nil {
		fmc.SetCount(*i)
	}
	return fmc
}

// AddArtifactMetricIDs adds the "artifact_metrics" edge to the ArtifactMetrics entity by IDs.
func (fmc *FilesMetricCreate) AddArtifactMetricIDs(ids ...int) *FilesMetricCreate {
	fmc.mutation.AddArtifactMetricIDs(ids...)
	return fmc
}

// AddArtifactMetrics adds the "artifact_metrics" edges to the ArtifactMetrics entity.
func (fmc *FilesMetricCreate) AddArtifactMetrics(a ...*ArtifactMetrics) *FilesMetricCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return fmc.AddArtifactMetricIDs(ids...)
}

// Mutation returns the FilesMetricMutation object of the builder.
func (fmc *FilesMetricCreate) Mutation() *FilesMetricMutation {
	return fmc.mutation
}

// Save creates the FilesMetric in the database.
func (fmc *FilesMetricCreate) Save(ctx context.Context) (*FilesMetric, error) {
	return withHooks(ctx, fmc.sqlSave, fmc.mutation, fmc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fmc *FilesMetricCreate) SaveX(ctx context.Context) *FilesMetric {
	v, err := fmc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fmc *FilesMetricCreate) Exec(ctx context.Context) error {
	_, err := fmc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fmc *FilesMetricCreate) ExecX(ctx context.Context) {
	if err := fmc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fmc *FilesMetricCreate) check() error {
	return nil
}

func (fmc *FilesMetricCreate) sqlSave(ctx context.Context) (*FilesMetric, error) {
	if err := fmc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fmc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fmc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	fmc.mutation.id = &_node.ID
	fmc.mutation.done = true
	return _node, nil
}

func (fmc *FilesMetricCreate) createSpec() (*FilesMetric, *sqlgraph.CreateSpec) {
	var (
		_node = &FilesMetric{config: fmc.config}
		_spec = sqlgraph.NewCreateSpec(filesmetric.Table, sqlgraph.NewFieldSpec(filesmetric.FieldID, field.TypeInt))
	)
	if value, ok := fmc.mutation.SizeInBytes(); ok {
		_spec.SetField(filesmetric.FieldSizeInBytes, field.TypeInt64, value)
		_node.SizeInBytes = value
	}
	if value, ok := fmc.mutation.Count(); ok {
		_spec.SetField(filesmetric.FieldCount, field.TypeInt32, value)
		_node.Count = value
	}
	if nodes := fmc.mutation.ArtifactMetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   filesmetric.ArtifactMetricsTable,
			Columns: filesmetric.ArtifactMetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(artifactmetrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// FilesMetricCreateBulk is the builder for creating many FilesMetric entities in bulk.
type FilesMetricCreateBulk struct {
	config
	err      error
	builders []*FilesMetricCreate
}

// Save creates the FilesMetric entities in the database.
func (fmcb *FilesMetricCreateBulk) Save(ctx context.Context) ([]*FilesMetric, error) {
	if fmcb.err != nil {
		return nil, fmcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(fmcb.builders))
	nodes := make([]*FilesMetric, len(fmcb.builders))
	mutators := make([]Mutator, len(fmcb.builders))
	for i := range fmcb.builders {
		func(i int, root context.Context) {
			builder := fmcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FilesMetricMutation)
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
					_, err = mutators[i+1].Mutate(root, fmcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fmcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, fmcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fmcb *FilesMetricCreateBulk) SaveX(ctx context.Context) []*FilesMetric {
	v, err := fmcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fmcb *FilesMetricCreateBulk) Exec(ctx context.Context) error {
	_, err := fmcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fmcb *FilesMetricCreateBulk) ExecX(ctx context.Context) {
	if err := fmcb.Exec(ctx); err != nil {
		panic(err)
	}
}