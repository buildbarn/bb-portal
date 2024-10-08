// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/artifactmetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/filesmetric"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
)

// ArtifactMetricsCreate is the builder for creating a ArtifactMetrics entity.
type ArtifactMetricsCreate struct {
	config
	mutation *ArtifactMetricsMutation
	hooks    []Hook
}

// AddMetricIDs adds the "metrics" edge to the Metrics entity by IDs.
func (amc *ArtifactMetricsCreate) AddMetricIDs(ids ...int) *ArtifactMetricsCreate {
	amc.mutation.AddMetricIDs(ids...)
	return amc
}

// AddMetrics adds the "metrics" edges to the Metrics entity.
func (amc *ArtifactMetricsCreate) AddMetrics(m ...*Metrics) *ArtifactMetricsCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return amc.AddMetricIDs(ids...)
}

// AddSourceArtifactsReadIDs adds the "source_artifacts_read" edge to the FilesMetric entity by IDs.
func (amc *ArtifactMetricsCreate) AddSourceArtifactsReadIDs(ids ...int) *ArtifactMetricsCreate {
	amc.mutation.AddSourceArtifactsReadIDs(ids...)
	return amc
}

// AddSourceArtifactsRead adds the "source_artifacts_read" edges to the FilesMetric entity.
func (amc *ArtifactMetricsCreate) AddSourceArtifactsRead(f ...*FilesMetric) *ArtifactMetricsCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return amc.AddSourceArtifactsReadIDs(ids...)
}

// AddOutputArtifactsSeenIDs adds the "output_artifacts_seen" edge to the FilesMetric entity by IDs.
func (amc *ArtifactMetricsCreate) AddOutputArtifactsSeenIDs(ids ...int) *ArtifactMetricsCreate {
	amc.mutation.AddOutputArtifactsSeenIDs(ids...)
	return amc
}

// AddOutputArtifactsSeen adds the "output_artifacts_seen" edges to the FilesMetric entity.
func (amc *ArtifactMetricsCreate) AddOutputArtifactsSeen(f ...*FilesMetric) *ArtifactMetricsCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return amc.AddOutputArtifactsSeenIDs(ids...)
}

// AddOutputArtifactsFromActionCacheIDs adds the "output_artifacts_from_action_cache" edge to the FilesMetric entity by IDs.
func (amc *ArtifactMetricsCreate) AddOutputArtifactsFromActionCacheIDs(ids ...int) *ArtifactMetricsCreate {
	amc.mutation.AddOutputArtifactsFromActionCacheIDs(ids...)
	return amc
}

// AddOutputArtifactsFromActionCache adds the "output_artifacts_from_action_cache" edges to the FilesMetric entity.
func (amc *ArtifactMetricsCreate) AddOutputArtifactsFromActionCache(f ...*FilesMetric) *ArtifactMetricsCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return amc.AddOutputArtifactsFromActionCacheIDs(ids...)
}

// AddTopLevelArtifactIDs adds the "top_level_artifacts" edge to the FilesMetric entity by IDs.
func (amc *ArtifactMetricsCreate) AddTopLevelArtifactIDs(ids ...int) *ArtifactMetricsCreate {
	amc.mutation.AddTopLevelArtifactIDs(ids...)
	return amc
}

// AddTopLevelArtifacts adds the "top_level_artifacts" edges to the FilesMetric entity.
func (amc *ArtifactMetricsCreate) AddTopLevelArtifacts(f ...*FilesMetric) *ArtifactMetricsCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return amc.AddTopLevelArtifactIDs(ids...)
}

// Mutation returns the ArtifactMetricsMutation object of the builder.
func (amc *ArtifactMetricsCreate) Mutation() *ArtifactMetricsMutation {
	return amc.mutation
}

// Save creates the ArtifactMetrics in the database.
func (amc *ArtifactMetricsCreate) Save(ctx context.Context) (*ArtifactMetrics, error) {
	return withHooks(ctx, amc.sqlSave, amc.mutation, amc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (amc *ArtifactMetricsCreate) SaveX(ctx context.Context) *ArtifactMetrics {
	v, err := amc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amc *ArtifactMetricsCreate) Exec(ctx context.Context) error {
	_, err := amc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amc *ArtifactMetricsCreate) ExecX(ctx context.Context) {
	if err := amc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (amc *ArtifactMetricsCreate) check() error {
	return nil
}

func (amc *ArtifactMetricsCreate) sqlSave(ctx context.Context) (*ArtifactMetrics, error) {
	if err := amc.check(); err != nil {
		return nil, err
	}
	_node, _spec := amc.createSpec()
	if err := sqlgraph.CreateNode(ctx, amc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	amc.mutation.id = &_node.ID
	amc.mutation.done = true
	return _node, nil
}

func (amc *ArtifactMetricsCreate) createSpec() (*ArtifactMetrics, *sqlgraph.CreateSpec) {
	var (
		_node = &ArtifactMetrics{config: amc.config}
		_spec = sqlgraph.NewCreateSpec(artifactmetrics.Table, sqlgraph.NewFieldSpec(artifactmetrics.FieldID, field.TypeInt))
	)
	if nodes := amc.mutation.MetricsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   artifactmetrics.MetricsTable,
			Columns: artifactmetrics.MetricsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metrics.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amc.mutation.SourceArtifactsReadIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   artifactmetrics.SourceArtifactsReadTable,
			Columns: []string{artifactmetrics.SourceArtifactsReadColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(filesmetric.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amc.mutation.OutputArtifactsSeenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   artifactmetrics.OutputArtifactsSeenTable,
			Columns: []string{artifactmetrics.OutputArtifactsSeenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(filesmetric.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amc.mutation.OutputArtifactsFromActionCacheIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   artifactmetrics.OutputArtifactsFromActionCacheTable,
			Columns: []string{artifactmetrics.OutputArtifactsFromActionCacheColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(filesmetric.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amc.mutation.TopLevelArtifactsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   artifactmetrics.TopLevelArtifactsTable,
			Columns: artifactmetrics.TopLevelArtifactsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(filesmetric.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ArtifactMetricsCreateBulk is the builder for creating many ArtifactMetrics entities in bulk.
type ArtifactMetricsCreateBulk struct {
	config
	err      error
	builders []*ArtifactMetricsCreate
}

// Save creates the ArtifactMetrics entities in the database.
func (amcb *ArtifactMetricsCreateBulk) Save(ctx context.Context) ([]*ArtifactMetrics, error) {
	if amcb.err != nil {
		return nil, amcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(amcb.builders))
	nodes := make([]*ArtifactMetrics, len(amcb.builders))
	mutators := make([]Mutator, len(amcb.builders))
	for i := range amcb.builders {
		func(i int, root context.Context) {
			builder := amcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ArtifactMetricsMutation)
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
					_, err = mutators[i+1].Mutate(root, amcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, amcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, amcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (amcb *ArtifactMetricsCreateBulk) SaveX(ctx context.Context) []*ArtifactMetrics {
	v, err := amcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amcb *ArtifactMetricsCreateBulk) Exec(ctx context.Context) error {
	_, err := amcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amcb *ArtifactMetricsCreateBulk) ExecX(ctx context.Context) {
	if err := amcb.Exec(ctx); err != nil {
		panic(err)
	}
}
