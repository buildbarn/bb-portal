// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/exectioninfo"
	"github.com/buildbarn/bb-portal/ent/gen/ent/resourceusage"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testresultbes"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingbreakdown"
)

// ExectionInfoCreate is the builder for creating a ExectionInfo entity.
type ExectionInfoCreate struct {
	config
	mutation *ExectionInfoMutation
	hooks    []Hook
}

// SetTimeoutSeconds sets the "timeout_seconds" field.
func (eic *ExectionInfoCreate) SetTimeoutSeconds(i int32) *ExectionInfoCreate {
	eic.mutation.SetTimeoutSeconds(i)
	return eic
}

// SetNillableTimeoutSeconds sets the "timeout_seconds" field if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableTimeoutSeconds(i *int32) *ExectionInfoCreate {
	if i != nil {
		eic.SetTimeoutSeconds(*i)
	}
	return eic
}

// SetStrategy sets the "strategy" field.
func (eic *ExectionInfoCreate) SetStrategy(s string) *ExectionInfoCreate {
	eic.mutation.SetStrategy(s)
	return eic
}

// SetNillableStrategy sets the "strategy" field if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableStrategy(s *string) *ExectionInfoCreate {
	if s != nil {
		eic.SetStrategy(*s)
	}
	return eic
}

// SetCachedRemotely sets the "cached_remotely" field.
func (eic *ExectionInfoCreate) SetCachedRemotely(b bool) *ExectionInfoCreate {
	eic.mutation.SetCachedRemotely(b)
	return eic
}

// SetNillableCachedRemotely sets the "cached_remotely" field if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableCachedRemotely(b *bool) *ExectionInfoCreate {
	if b != nil {
		eic.SetCachedRemotely(*b)
	}
	return eic
}

// SetExitCode sets the "exit_code" field.
func (eic *ExectionInfoCreate) SetExitCode(i int32) *ExectionInfoCreate {
	eic.mutation.SetExitCode(i)
	return eic
}

// SetNillableExitCode sets the "exit_code" field if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableExitCode(i *int32) *ExectionInfoCreate {
	if i != nil {
		eic.SetExitCode(*i)
	}
	return eic
}

// SetHostname sets the "hostname" field.
func (eic *ExectionInfoCreate) SetHostname(s string) *ExectionInfoCreate {
	eic.mutation.SetHostname(s)
	return eic
}

// SetNillableHostname sets the "hostname" field if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableHostname(s *string) *ExectionInfoCreate {
	if s != nil {
		eic.SetHostname(*s)
	}
	return eic
}

// SetExecutionInfoID sets the "execution_info_id" field.
func (eic *ExectionInfoCreate) SetExecutionInfoID(i int) *ExectionInfoCreate {
	eic.mutation.SetExecutionInfoID(i)
	return eic
}

// SetNillableExecutionInfoID sets the "execution_info_id" field if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableExecutionInfoID(i *int) *ExectionInfoCreate {
	if i != nil {
		eic.SetExecutionInfoID(*i)
	}
	return eic
}

// SetTestResultID sets the "test_result" edge to the TestResultBES entity by ID.
func (eic *ExectionInfoCreate) SetTestResultID(id int) *ExectionInfoCreate {
	eic.mutation.SetTestResultID(id)
	return eic
}

// SetNillableTestResultID sets the "test_result" edge to the TestResultBES entity by ID if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableTestResultID(id *int) *ExectionInfoCreate {
	if id != nil {
		eic = eic.SetTestResultID(*id)
	}
	return eic
}

// SetTestResult sets the "test_result" edge to the TestResultBES entity.
func (eic *ExectionInfoCreate) SetTestResult(t *TestResultBES) *ExectionInfoCreate {
	return eic.SetTestResultID(t.ID)
}

// SetTimingBreakdownID sets the "timing_breakdown" edge to the TimingBreakdown entity by ID.
func (eic *ExectionInfoCreate) SetTimingBreakdownID(id int) *ExectionInfoCreate {
	eic.mutation.SetTimingBreakdownID(id)
	return eic
}

// SetNillableTimingBreakdownID sets the "timing_breakdown" edge to the TimingBreakdown entity by ID if the given value is not nil.
func (eic *ExectionInfoCreate) SetNillableTimingBreakdownID(id *int) *ExectionInfoCreate {
	if id != nil {
		eic = eic.SetTimingBreakdownID(*id)
	}
	return eic
}

// SetTimingBreakdown sets the "timing_breakdown" edge to the TimingBreakdown entity.
func (eic *ExectionInfoCreate) SetTimingBreakdown(t *TimingBreakdown) *ExectionInfoCreate {
	return eic.SetTimingBreakdownID(t.ID)
}

// AddResourceUsageIDs adds the "resource_usage" edge to the ResourceUsage entity by IDs.
func (eic *ExectionInfoCreate) AddResourceUsageIDs(ids ...int) *ExectionInfoCreate {
	eic.mutation.AddResourceUsageIDs(ids...)
	return eic
}

// AddResourceUsage adds the "resource_usage" edges to the ResourceUsage entity.
func (eic *ExectionInfoCreate) AddResourceUsage(r ...*ResourceUsage) *ExectionInfoCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return eic.AddResourceUsageIDs(ids...)
}

// Mutation returns the ExectionInfoMutation object of the builder.
func (eic *ExectionInfoCreate) Mutation() *ExectionInfoMutation {
	return eic.mutation
}

// Save creates the ExectionInfo in the database.
func (eic *ExectionInfoCreate) Save(ctx context.Context) (*ExectionInfo, error) {
	return withHooks(ctx, eic.sqlSave, eic.mutation, eic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (eic *ExectionInfoCreate) SaveX(ctx context.Context) *ExectionInfo {
	v, err := eic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (eic *ExectionInfoCreate) Exec(ctx context.Context) error {
	_, err := eic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eic *ExectionInfoCreate) ExecX(ctx context.Context) {
	if err := eic.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eic *ExectionInfoCreate) check() error {
	return nil
}

func (eic *ExectionInfoCreate) sqlSave(ctx context.Context) (*ExectionInfo, error) {
	if err := eic.check(); err != nil {
		return nil, err
	}
	_node, _spec := eic.createSpec()
	if err := sqlgraph.CreateNode(ctx, eic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	eic.mutation.id = &_node.ID
	eic.mutation.done = true
	return _node, nil
}

func (eic *ExectionInfoCreate) createSpec() (*ExectionInfo, *sqlgraph.CreateSpec) {
	var (
		_node = &ExectionInfo{config: eic.config}
		_spec = sqlgraph.NewCreateSpec(exectioninfo.Table, sqlgraph.NewFieldSpec(exectioninfo.FieldID, field.TypeInt))
	)
	if value, ok := eic.mutation.TimeoutSeconds(); ok {
		_spec.SetField(exectioninfo.FieldTimeoutSeconds, field.TypeInt32, value)
		_node.TimeoutSeconds = value
	}
	if value, ok := eic.mutation.Strategy(); ok {
		_spec.SetField(exectioninfo.FieldStrategy, field.TypeString, value)
		_node.Strategy = value
	}
	if value, ok := eic.mutation.CachedRemotely(); ok {
		_spec.SetField(exectioninfo.FieldCachedRemotely, field.TypeBool, value)
		_node.CachedRemotely = value
	}
	if value, ok := eic.mutation.ExitCode(); ok {
		_spec.SetField(exectioninfo.FieldExitCode, field.TypeInt32, value)
		_node.ExitCode = value
	}
	if value, ok := eic.mutation.Hostname(); ok {
		_spec.SetField(exectioninfo.FieldHostname, field.TypeString, value)
		_node.Hostname = value
	}
	if nodes := eic.mutation.TestResultIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   exectioninfo.TestResultTable,
			Columns: []string{exectioninfo.TestResultColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testresultbes.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ExecutionInfoID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := eic.mutation.TimingBreakdownIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   exectioninfo.TimingBreakdownTable,
			Columns: []string{exectioninfo.TimingBreakdownColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timingbreakdown.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := eic.mutation.ResourceUsageIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   exectioninfo.ResourceUsageTable,
			Columns: []string{exectioninfo.ResourceUsageColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resourceusage.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ExectionInfoCreateBulk is the builder for creating many ExectionInfo entities in bulk.
type ExectionInfoCreateBulk struct {
	config
	err      error
	builders []*ExectionInfoCreate
}

// Save creates the ExectionInfo entities in the database.
func (eicb *ExectionInfoCreateBulk) Save(ctx context.Context) ([]*ExectionInfo, error) {
	if eicb.err != nil {
		return nil, eicb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(eicb.builders))
	nodes := make([]*ExectionInfo, len(eicb.builders))
	mutators := make([]Mutator, len(eicb.builders))
	for i := range eicb.builders {
		func(i int, root context.Context) {
			builder := eicb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ExectionInfoMutation)
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
					_, err = mutators[i+1].Mutate(root, eicb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, eicb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, eicb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (eicb *ExectionInfoCreateBulk) SaveX(ctx context.Context) []*ExectionInfo {
	v, err := eicb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (eicb *ExectionInfoCreateBulk) Exec(ctx context.Context) error {
	_, err := eicb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eicb *ExectionInfoCreateBulk) ExecX(ctx context.Context) {
	if err := eicb.Exec(ctx); err != nil {
		panic(err)
	}
}
