// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testcollection"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testresultbes"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testsummary"
)

// TestCollectionCreate is the builder for creating a TestCollection entity.
type TestCollectionCreate struct {
	config
	mutation *TestCollectionMutation
	hooks    []Hook
}

// SetLabel sets the "label" field.
func (tcc *TestCollectionCreate) SetLabel(s string) *TestCollectionCreate {
	tcc.mutation.SetLabel(s)
	return tcc
}

// SetNillableLabel sets the "label" field if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableLabel(s *string) *TestCollectionCreate {
	if s != nil {
		tcc.SetLabel(*s)
	}
	return tcc
}

// SetOverallStatus sets the "overall_status" field.
func (tcc *TestCollectionCreate) SetOverallStatus(ts testcollection.OverallStatus) *TestCollectionCreate {
	tcc.mutation.SetOverallStatus(ts)
	return tcc
}

// SetNillableOverallStatus sets the "overall_status" field if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableOverallStatus(ts *testcollection.OverallStatus) *TestCollectionCreate {
	if ts != nil {
		tcc.SetOverallStatus(*ts)
	}
	return tcc
}

// SetStrategy sets the "strategy" field.
func (tcc *TestCollectionCreate) SetStrategy(s string) *TestCollectionCreate {
	tcc.mutation.SetStrategy(s)
	return tcc
}

// SetNillableStrategy sets the "strategy" field if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableStrategy(s *string) *TestCollectionCreate {
	if s != nil {
		tcc.SetStrategy(*s)
	}
	return tcc
}

// SetCachedLocally sets the "cached_locally" field.
func (tcc *TestCollectionCreate) SetCachedLocally(b bool) *TestCollectionCreate {
	tcc.mutation.SetCachedLocally(b)
	return tcc
}

// SetNillableCachedLocally sets the "cached_locally" field if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableCachedLocally(b *bool) *TestCollectionCreate {
	if b != nil {
		tcc.SetCachedLocally(*b)
	}
	return tcc
}

// SetCachedRemotely sets the "cached_remotely" field.
func (tcc *TestCollectionCreate) SetCachedRemotely(b bool) *TestCollectionCreate {
	tcc.mutation.SetCachedRemotely(b)
	return tcc
}

// SetNillableCachedRemotely sets the "cached_remotely" field if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableCachedRemotely(b *bool) *TestCollectionCreate {
	if b != nil {
		tcc.SetCachedRemotely(*b)
	}
	return tcc
}

// SetDurationMs sets the "duration_ms" field.
func (tcc *TestCollectionCreate) SetDurationMs(i int64) *TestCollectionCreate {
	tcc.mutation.SetDurationMs(i)
	return tcc
}

// SetNillableDurationMs sets the "duration_ms" field if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableDurationMs(i *int64) *TestCollectionCreate {
	if i != nil {
		tcc.SetDurationMs(*i)
	}
	return tcc
}

// AddBazelInvocationIDs adds the "bazel_invocation" edge to the BazelInvocation entity by IDs.
func (tcc *TestCollectionCreate) AddBazelInvocationIDs(ids ...int) *TestCollectionCreate {
	tcc.mutation.AddBazelInvocationIDs(ids...)
	return tcc
}

// AddBazelInvocation adds the "bazel_invocation" edges to the BazelInvocation entity.
func (tcc *TestCollectionCreate) AddBazelInvocation(b ...*BazelInvocation) *TestCollectionCreate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return tcc.AddBazelInvocationIDs(ids...)
}

// SetTestSummaryID sets the "test_summary" edge to the TestSummary entity by ID.
func (tcc *TestCollectionCreate) SetTestSummaryID(id int) *TestCollectionCreate {
	tcc.mutation.SetTestSummaryID(id)
	return tcc
}

// SetNillableTestSummaryID sets the "test_summary" edge to the TestSummary entity by ID if the given value is not nil.
func (tcc *TestCollectionCreate) SetNillableTestSummaryID(id *int) *TestCollectionCreate {
	if id != nil {
		tcc = tcc.SetTestSummaryID(*id)
	}
	return tcc
}

// SetTestSummary sets the "test_summary" edge to the TestSummary entity.
func (tcc *TestCollectionCreate) SetTestSummary(t *TestSummary) *TestCollectionCreate {
	return tcc.SetTestSummaryID(t.ID)
}

// AddTestResultIDs adds the "test_results" edge to the TestResultBES entity by IDs.
func (tcc *TestCollectionCreate) AddTestResultIDs(ids ...int) *TestCollectionCreate {
	tcc.mutation.AddTestResultIDs(ids...)
	return tcc
}

// AddTestResults adds the "test_results" edges to the TestResultBES entity.
func (tcc *TestCollectionCreate) AddTestResults(t ...*TestResultBES) *TestCollectionCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcc.AddTestResultIDs(ids...)
}

// Mutation returns the TestCollectionMutation object of the builder.
func (tcc *TestCollectionCreate) Mutation() *TestCollectionMutation {
	return tcc.mutation
}

// Save creates the TestCollection in the database.
func (tcc *TestCollectionCreate) Save(ctx context.Context) (*TestCollection, error) {
	tcc.defaults()
	return withHooks(ctx, tcc.sqlSave, tcc.mutation, tcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tcc *TestCollectionCreate) SaveX(ctx context.Context) *TestCollection {
	v, err := tcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcc *TestCollectionCreate) Exec(ctx context.Context) error {
	_, err := tcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcc *TestCollectionCreate) ExecX(ctx context.Context) {
	if err := tcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tcc *TestCollectionCreate) defaults() {
	if _, ok := tcc.mutation.OverallStatus(); !ok {
		v := testcollection.DefaultOverallStatus
		tcc.mutation.SetOverallStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcc *TestCollectionCreate) check() error {
	if v, ok := tcc.mutation.OverallStatus(); ok {
		if err := testcollection.OverallStatusValidator(v); err != nil {
			return &ValidationError{Name: "overall_status", err: fmt.Errorf(`ent: validator failed for field "TestCollection.overall_status": %w`, err)}
		}
	}
	return nil
}

func (tcc *TestCollectionCreate) sqlSave(ctx context.Context) (*TestCollection, error) {
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

func (tcc *TestCollectionCreate) createSpec() (*TestCollection, *sqlgraph.CreateSpec) {
	var (
		_node = &TestCollection{config: tcc.config}
		_spec = sqlgraph.NewCreateSpec(testcollection.Table, sqlgraph.NewFieldSpec(testcollection.FieldID, field.TypeInt))
	)
	if value, ok := tcc.mutation.Label(); ok {
		_spec.SetField(testcollection.FieldLabel, field.TypeString, value)
		_node.Label = value
	}
	if value, ok := tcc.mutation.OverallStatus(); ok {
		_spec.SetField(testcollection.FieldOverallStatus, field.TypeEnum, value)
		_node.OverallStatus = value
	}
	if value, ok := tcc.mutation.Strategy(); ok {
		_spec.SetField(testcollection.FieldStrategy, field.TypeString, value)
		_node.Strategy = value
	}
	if value, ok := tcc.mutation.CachedLocally(); ok {
		_spec.SetField(testcollection.FieldCachedLocally, field.TypeBool, value)
		_node.CachedLocally = value
	}
	if value, ok := tcc.mutation.CachedRemotely(); ok {
		_spec.SetField(testcollection.FieldCachedRemotely, field.TypeBool, value)
		_node.CachedRemotely = value
	}
	if value, ok := tcc.mutation.DurationMs(); ok {
		_spec.SetField(testcollection.FieldDurationMs, field.TypeInt64, value)
		_node.DurationMs = value
	}
	if nodes := tcc.mutation.BazelInvocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   testcollection.BazelInvocationTable,
			Columns: testcollection.BazelInvocationPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bazelinvocation.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tcc.mutation.TestSummaryIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   testcollection.TestSummaryTable,
			Columns: []string{testcollection.TestSummaryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testsummary.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.test_collection_test_summary = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tcc.mutation.TestResultsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   testcollection.TestResultsTable,
			Columns: []string{testcollection.TestResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testresultbes.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TestCollectionCreateBulk is the builder for creating many TestCollection entities in bulk.
type TestCollectionCreateBulk struct {
	config
	err      error
	builders []*TestCollectionCreate
}

// Save creates the TestCollection entities in the database.
func (tccb *TestCollectionCreateBulk) Save(ctx context.Context) ([]*TestCollection, error) {
	if tccb.err != nil {
		return nil, tccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tccb.builders))
	nodes := make([]*TestCollection, len(tccb.builders))
	mutators := make([]Mutator, len(tccb.builders))
	for i := range tccb.builders {
		func(i int, root context.Context) {
			builder := tccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TestCollectionMutation)
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
func (tccb *TestCollectionCreateBulk) SaveX(ctx context.Context) []*TestCollection {
	v, err := tccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tccb *TestCollectionCreateBulk) Exec(ctx context.Context) error {
	_, err := tccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tccb *TestCollectionCreateBulk) ExecX(ctx context.Context) {
	if err := tccb.Exec(ctx); err != nil {
		panic(err)
	}
}
