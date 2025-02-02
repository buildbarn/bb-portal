// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/namedsetoffiles"
	"github.com/buildbarn/bb-portal/ent/gen/ent/outputgroup"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetcomplete"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testfile"
)

// OutputGroupCreate is the builder for creating a OutputGroup entity.
type OutputGroupCreate struct {
	config
	mutation *OutputGroupMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (ogc *OutputGroupCreate) SetName(s string) *OutputGroupCreate {
	ogc.mutation.SetName(s)
	return ogc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ogc *OutputGroupCreate) SetNillableName(s *string) *OutputGroupCreate {
	if s != nil {
		ogc.SetName(*s)
	}
	return ogc
}

// SetIncomplete sets the "incomplete" field.
func (ogc *OutputGroupCreate) SetIncomplete(b bool) *OutputGroupCreate {
	ogc.mutation.SetIncomplete(b)
	return ogc
}

// SetNillableIncomplete sets the "incomplete" field if the given value is not nil.
func (ogc *OutputGroupCreate) SetNillableIncomplete(b *bool) *OutputGroupCreate {
	if b != nil {
		ogc.SetIncomplete(*b)
	}
	return ogc
}

// SetTargetCompleteID sets the "target_complete" edge to the TargetComplete entity by ID.
func (ogc *OutputGroupCreate) SetTargetCompleteID(id int) *OutputGroupCreate {
	ogc.mutation.SetTargetCompleteID(id)
	return ogc
}

// SetNillableTargetCompleteID sets the "target_complete" edge to the TargetComplete entity by ID if the given value is not nil.
func (ogc *OutputGroupCreate) SetNillableTargetCompleteID(id *int) *OutputGroupCreate {
	if id != nil {
		ogc = ogc.SetTargetCompleteID(*id)
	}
	return ogc
}

// SetTargetComplete sets the "target_complete" edge to the TargetComplete entity.
func (ogc *OutputGroupCreate) SetTargetComplete(t *TargetComplete) *OutputGroupCreate {
	return ogc.SetTargetCompleteID(t.ID)
}

// AddInlineFileIDs adds the "inline_files" edge to the TestFile entity by IDs.
func (ogc *OutputGroupCreate) AddInlineFileIDs(ids ...int) *OutputGroupCreate {
	ogc.mutation.AddInlineFileIDs(ids...)
	return ogc
}

// AddInlineFiles adds the "inline_files" edges to the TestFile entity.
func (ogc *OutputGroupCreate) AddInlineFiles(t ...*TestFile) *OutputGroupCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ogc.AddInlineFileIDs(ids...)
}

// SetFileSetsID sets the "file_sets" edge to the NamedSetOfFiles entity by ID.
func (ogc *OutputGroupCreate) SetFileSetsID(id int) *OutputGroupCreate {
	ogc.mutation.SetFileSetsID(id)
	return ogc
}

// SetNillableFileSetsID sets the "file_sets" edge to the NamedSetOfFiles entity by ID if the given value is not nil.
func (ogc *OutputGroupCreate) SetNillableFileSetsID(id *int) *OutputGroupCreate {
	if id != nil {
		ogc = ogc.SetFileSetsID(*id)
	}
	return ogc
}

// SetFileSets sets the "file_sets" edge to the NamedSetOfFiles entity.
func (ogc *OutputGroupCreate) SetFileSets(n *NamedSetOfFiles) *OutputGroupCreate {
	return ogc.SetFileSetsID(n.ID)
}

// Mutation returns the OutputGroupMutation object of the builder.
func (ogc *OutputGroupCreate) Mutation() *OutputGroupMutation {
	return ogc.mutation
}

// Save creates the OutputGroup in the database.
func (ogc *OutputGroupCreate) Save(ctx context.Context) (*OutputGroup, error) {
	return withHooks(ctx, ogc.sqlSave, ogc.mutation, ogc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ogc *OutputGroupCreate) SaveX(ctx context.Context) *OutputGroup {
	v, err := ogc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ogc *OutputGroupCreate) Exec(ctx context.Context) error {
	_, err := ogc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ogc *OutputGroupCreate) ExecX(ctx context.Context) {
	if err := ogc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ogc *OutputGroupCreate) check() error {
	return nil
}

func (ogc *OutputGroupCreate) sqlSave(ctx context.Context) (*OutputGroup, error) {
	if err := ogc.check(); err != nil {
		return nil, err
	}
	_node, _spec := ogc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ogc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ogc.mutation.id = &_node.ID
	ogc.mutation.done = true
	return _node, nil
}

func (ogc *OutputGroupCreate) createSpec() (*OutputGroup, *sqlgraph.CreateSpec) {
	var (
		_node = &OutputGroup{config: ogc.config}
		_spec = sqlgraph.NewCreateSpec(outputgroup.Table, sqlgraph.NewFieldSpec(outputgroup.FieldID, field.TypeInt))
	)
	if value, ok := ogc.mutation.Name(); ok {
		_spec.SetField(outputgroup.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := ogc.mutation.Incomplete(); ok {
		_spec.SetField(outputgroup.FieldIncomplete, field.TypeBool, value)
		_node.Incomplete = value
	}
	if nodes := ogc.mutation.TargetCompleteIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   outputgroup.TargetCompleteTable,
			Columns: []string{outputgroup.TargetCompleteColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(targetcomplete.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.target_complete_output_group = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ogc.mutation.InlineFilesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   outputgroup.InlineFilesTable,
			Columns: []string{outputgroup.InlineFilesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(testfile.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ogc.mutation.FileSetsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   outputgroup.FileSetsTable,
			Columns: []string{outputgroup.FileSetsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(namedsetoffiles.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OutputGroupCreateBulk is the builder for creating many OutputGroup entities in bulk.
type OutputGroupCreateBulk struct {
	config
	err      error
	builders []*OutputGroupCreate
}

// Save creates the OutputGroup entities in the database.
func (ogcb *OutputGroupCreateBulk) Save(ctx context.Context) ([]*OutputGroup, error) {
	if ogcb.err != nil {
		return nil, ogcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ogcb.builders))
	nodes := make([]*OutputGroup, len(ogcb.builders))
	mutators := make([]Mutator, len(ogcb.builders))
	for i := range ogcb.builders {
		func(i int, root context.Context) {
			builder := ogcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OutputGroupMutation)
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
					_, err = mutators[i+1].Mutate(root, ogcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ogcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ogcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ogcb *OutputGroupCreateBulk) SaveX(ctx context.Context) []*OutputGroup {
	v, err := ogcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ogcb *OutputGroupCreateBulk) Exec(ctx context.Context) error {
	_, err := ogcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ogcb *OutputGroupCreateBulk) ExecX(ctx context.Context) {
	if err := ogcb.Exec(ctx); err != nil {
		panic(err)
	}
}
