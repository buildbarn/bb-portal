// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/exectioninfo"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingbreakdown"
	"github.com/buildbarn/bb-portal/ent/gen/ent/timingchild"
)

// TimingBreakdownQuery is the builder for querying TimingBreakdown entities.
type TimingBreakdownQuery struct {
	config
	ctx                    *QueryContext
	order                  []timingbreakdown.OrderOption
	inters                 []Interceptor
	predicates             []predicate.TimingBreakdown
	withExecutionInfo      *ExectionInfoQuery
	withChild              *TimingChildQuery
	modifiers              []func(*sql.Selector)
	loadTotal              []func(context.Context, []*TimingBreakdown) error
	withNamedExecutionInfo map[string]*ExectionInfoQuery
	withNamedChild         map[string]*TimingChildQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TimingBreakdownQuery builder.
func (tbq *TimingBreakdownQuery) Where(ps ...predicate.TimingBreakdown) *TimingBreakdownQuery {
	tbq.predicates = append(tbq.predicates, ps...)
	return tbq
}

// Limit the number of records to be returned by this query.
func (tbq *TimingBreakdownQuery) Limit(limit int) *TimingBreakdownQuery {
	tbq.ctx.Limit = &limit
	return tbq
}

// Offset to start from.
func (tbq *TimingBreakdownQuery) Offset(offset int) *TimingBreakdownQuery {
	tbq.ctx.Offset = &offset
	return tbq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tbq *TimingBreakdownQuery) Unique(unique bool) *TimingBreakdownQuery {
	tbq.ctx.Unique = &unique
	return tbq
}

// Order specifies how the records should be ordered.
func (tbq *TimingBreakdownQuery) Order(o ...timingbreakdown.OrderOption) *TimingBreakdownQuery {
	tbq.order = append(tbq.order, o...)
	return tbq
}

// QueryExecutionInfo chains the current query on the "execution_info" edge.
func (tbq *TimingBreakdownQuery) QueryExecutionInfo() *ExectionInfoQuery {
	query := (&ExectionInfoClient{config: tbq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tbq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tbq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(timingbreakdown.Table, timingbreakdown.FieldID, selector),
			sqlgraph.To(exectioninfo.Table, exectioninfo.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, timingbreakdown.ExecutionInfoTable, timingbreakdown.ExecutionInfoColumn),
		)
		fromU = sqlgraph.SetNeighbors(tbq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryChild chains the current query on the "child" edge.
func (tbq *TimingBreakdownQuery) QueryChild() *TimingChildQuery {
	query := (&TimingChildClient{config: tbq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tbq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tbq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(timingbreakdown.Table, timingbreakdown.FieldID, selector),
			sqlgraph.To(timingchild.Table, timingchild.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, timingbreakdown.ChildTable, timingbreakdown.ChildPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(tbq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first TimingBreakdown entity from the query.
// Returns a *NotFoundError when no TimingBreakdown was found.
func (tbq *TimingBreakdownQuery) First(ctx context.Context) (*TimingBreakdown, error) {
	nodes, err := tbq.Limit(1).All(setContextOp(ctx, tbq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{timingbreakdown.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) FirstX(ctx context.Context) *TimingBreakdown {
	node, err := tbq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TimingBreakdown ID from the query.
// Returns a *NotFoundError when no TimingBreakdown ID was found.
func (tbq *TimingBreakdownQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tbq.Limit(1).IDs(setContextOp(ctx, tbq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{timingbreakdown.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) FirstIDX(ctx context.Context) int {
	id, err := tbq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TimingBreakdown entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one TimingBreakdown entity is found.
// Returns a *NotFoundError when no TimingBreakdown entities are found.
func (tbq *TimingBreakdownQuery) Only(ctx context.Context) (*TimingBreakdown, error) {
	nodes, err := tbq.Limit(2).All(setContextOp(ctx, tbq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{timingbreakdown.Label}
	default:
		return nil, &NotSingularError{timingbreakdown.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) OnlyX(ctx context.Context) *TimingBreakdown {
	node, err := tbq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TimingBreakdown ID in the query.
// Returns a *NotSingularError when more than one TimingBreakdown ID is found.
// Returns a *NotFoundError when no entities are found.
func (tbq *TimingBreakdownQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tbq.Limit(2).IDs(setContextOp(ctx, tbq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{timingbreakdown.Label}
	default:
		err = &NotSingularError{timingbreakdown.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) OnlyIDX(ctx context.Context) int {
	id, err := tbq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TimingBreakdowns.
func (tbq *TimingBreakdownQuery) All(ctx context.Context) ([]*TimingBreakdown, error) {
	ctx = setContextOp(ctx, tbq.ctx, "All")
	if err := tbq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*TimingBreakdown, *TimingBreakdownQuery]()
	return withInterceptors[[]*TimingBreakdown](ctx, tbq, qr, tbq.inters)
}

// AllX is like All, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) AllX(ctx context.Context) []*TimingBreakdown {
	nodes, err := tbq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TimingBreakdown IDs.
func (tbq *TimingBreakdownQuery) IDs(ctx context.Context) (ids []int, err error) {
	if tbq.ctx.Unique == nil && tbq.path != nil {
		tbq.Unique(true)
	}
	ctx = setContextOp(ctx, tbq.ctx, "IDs")
	if err = tbq.Select(timingbreakdown.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) IDsX(ctx context.Context) []int {
	ids, err := tbq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tbq *TimingBreakdownQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, tbq.ctx, "Count")
	if err := tbq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, tbq, querierCount[*TimingBreakdownQuery](), tbq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) CountX(ctx context.Context) int {
	count, err := tbq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tbq *TimingBreakdownQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, tbq.ctx, "Exist")
	switch _, err := tbq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (tbq *TimingBreakdownQuery) ExistX(ctx context.Context) bool {
	exist, err := tbq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TimingBreakdownQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tbq *TimingBreakdownQuery) Clone() *TimingBreakdownQuery {
	if tbq == nil {
		return nil
	}
	return &TimingBreakdownQuery{
		config:            tbq.config,
		ctx:               tbq.ctx.Clone(),
		order:             append([]timingbreakdown.OrderOption{}, tbq.order...),
		inters:            append([]Interceptor{}, tbq.inters...),
		predicates:        append([]predicate.TimingBreakdown{}, tbq.predicates...),
		withExecutionInfo: tbq.withExecutionInfo.Clone(),
		withChild:         tbq.withChild.Clone(),
		// clone intermediate query.
		sql:  tbq.sql.Clone(),
		path: tbq.path,
	}
}

// WithExecutionInfo tells the query-builder to eager-load the nodes that are connected to
// the "execution_info" edge. The optional arguments are used to configure the query builder of the edge.
func (tbq *TimingBreakdownQuery) WithExecutionInfo(opts ...func(*ExectionInfoQuery)) *TimingBreakdownQuery {
	query := (&ExectionInfoClient{config: tbq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tbq.withExecutionInfo = query
	return tbq
}

// WithChild tells the query-builder to eager-load the nodes that are connected to
// the "child" edge. The optional arguments are used to configure the query builder of the edge.
func (tbq *TimingBreakdownQuery) WithChild(opts ...func(*TimingChildQuery)) *TimingBreakdownQuery {
	query := (&TimingChildClient{config: tbq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tbq.withChild = query
	return tbq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.TimingBreakdown.Query().
//		GroupBy(timingbreakdown.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (tbq *TimingBreakdownQuery) GroupBy(field string, fields ...string) *TimingBreakdownGroupBy {
	tbq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &TimingBreakdownGroupBy{build: tbq}
	grbuild.flds = &tbq.ctx.Fields
	grbuild.label = timingbreakdown.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.TimingBreakdown.Query().
//		Select(timingbreakdown.FieldName).
//		Scan(ctx, &v)
func (tbq *TimingBreakdownQuery) Select(fields ...string) *TimingBreakdownSelect {
	tbq.ctx.Fields = append(tbq.ctx.Fields, fields...)
	sbuild := &TimingBreakdownSelect{TimingBreakdownQuery: tbq}
	sbuild.label = timingbreakdown.Label
	sbuild.flds, sbuild.scan = &tbq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a TimingBreakdownSelect configured with the given aggregations.
func (tbq *TimingBreakdownQuery) Aggregate(fns ...AggregateFunc) *TimingBreakdownSelect {
	return tbq.Select().Aggregate(fns...)
}

func (tbq *TimingBreakdownQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range tbq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, tbq); err != nil {
				return err
			}
		}
	}
	for _, f := range tbq.ctx.Fields {
		if !timingbreakdown.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tbq.path != nil {
		prev, err := tbq.path(ctx)
		if err != nil {
			return err
		}
		tbq.sql = prev
	}
	return nil
}

func (tbq *TimingBreakdownQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*TimingBreakdown, error) {
	var (
		nodes       = []*TimingBreakdown{}
		_spec       = tbq.querySpec()
		loadedTypes = [2]bool{
			tbq.withExecutionInfo != nil,
			tbq.withChild != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*TimingBreakdown).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &TimingBreakdown{config: tbq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(tbq.modifiers) > 0 {
		_spec.Modifiers = tbq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tbq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := tbq.withExecutionInfo; query != nil {
		if err := tbq.loadExecutionInfo(ctx, query, nodes,
			func(n *TimingBreakdown) { n.Edges.ExecutionInfo = []*ExectionInfo{} },
			func(n *TimingBreakdown, e *ExectionInfo) { n.Edges.ExecutionInfo = append(n.Edges.ExecutionInfo, e) }); err != nil {
			return nil, err
		}
	}
	if query := tbq.withChild; query != nil {
		if err := tbq.loadChild(ctx, query, nodes,
			func(n *TimingBreakdown) { n.Edges.Child = []*TimingChild{} },
			func(n *TimingBreakdown, e *TimingChild) { n.Edges.Child = append(n.Edges.Child, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range tbq.withNamedExecutionInfo {
		if err := tbq.loadExecutionInfo(ctx, query, nodes,
			func(n *TimingBreakdown) { n.appendNamedExecutionInfo(name) },
			func(n *TimingBreakdown, e *ExectionInfo) { n.appendNamedExecutionInfo(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range tbq.withNamedChild {
		if err := tbq.loadChild(ctx, query, nodes,
			func(n *TimingBreakdown) { n.appendNamedChild(name) },
			func(n *TimingBreakdown, e *TimingChild) { n.appendNamedChild(name, e) }); err != nil {
			return nil, err
		}
	}
	for i := range tbq.loadTotal {
		if err := tbq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (tbq *TimingBreakdownQuery) loadExecutionInfo(ctx context.Context, query *ExectionInfoQuery, nodes []*TimingBreakdown, init func(*TimingBreakdown), assign func(*TimingBreakdown, *ExectionInfo)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*TimingBreakdown)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.ExectionInfo(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(timingbreakdown.ExecutionInfoColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.exection_info_timing_breakdown
		if fk == nil {
			return fmt.Errorf(`foreign-key "exection_info_timing_breakdown" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "exection_info_timing_breakdown" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (tbq *TimingBreakdownQuery) loadChild(ctx context.Context, query *TimingChildQuery, nodes []*TimingBreakdown, init func(*TimingBreakdown), assign func(*TimingBreakdown, *TimingChild)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*TimingBreakdown)
	nids := make(map[int]map[*TimingBreakdown]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(timingbreakdown.ChildTable)
		s.Join(joinT).On(s.C(timingchild.FieldID), joinT.C(timingbreakdown.ChildPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(timingbreakdown.ChildPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(timingbreakdown.ChildPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*TimingBreakdown]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*TimingChild](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "child" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (tbq *TimingBreakdownQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tbq.querySpec()
	if len(tbq.modifiers) > 0 {
		_spec.Modifiers = tbq.modifiers
	}
	_spec.Node.Columns = tbq.ctx.Fields
	if len(tbq.ctx.Fields) > 0 {
		_spec.Unique = tbq.ctx.Unique != nil && *tbq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, tbq.driver, _spec)
}

func (tbq *TimingBreakdownQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(timingbreakdown.Table, timingbreakdown.Columns, sqlgraph.NewFieldSpec(timingbreakdown.FieldID, field.TypeInt))
	_spec.From = tbq.sql
	if unique := tbq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if tbq.path != nil {
		_spec.Unique = true
	}
	if fields := tbq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, timingbreakdown.FieldID)
		for i := range fields {
			if fields[i] != timingbreakdown.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tbq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tbq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tbq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tbq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tbq *TimingBreakdownQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tbq.driver.Dialect())
	t1 := builder.Table(timingbreakdown.Table)
	columns := tbq.ctx.Fields
	if len(columns) == 0 {
		columns = timingbreakdown.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tbq.sql != nil {
		selector = tbq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tbq.ctx.Unique != nil && *tbq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range tbq.predicates {
		p(selector)
	}
	for _, p := range tbq.order {
		p(selector)
	}
	if offset := tbq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tbq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WithNamedExecutionInfo tells the query-builder to eager-load the nodes that are connected to the "execution_info"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (tbq *TimingBreakdownQuery) WithNamedExecutionInfo(name string, opts ...func(*ExectionInfoQuery)) *TimingBreakdownQuery {
	query := (&ExectionInfoClient{config: tbq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if tbq.withNamedExecutionInfo == nil {
		tbq.withNamedExecutionInfo = make(map[string]*ExectionInfoQuery)
	}
	tbq.withNamedExecutionInfo[name] = query
	return tbq
}

// WithNamedChild tells the query-builder to eager-load the nodes that are connected to the "child"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (tbq *TimingBreakdownQuery) WithNamedChild(name string, opts ...func(*TimingChildQuery)) *TimingBreakdownQuery {
	query := (&TimingChildClient{config: tbq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if tbq.withNamedChild == nil {
		tbq.withNamedChild = make(map[string]*TimingChildQuery)
	}
	tbq.withNamedChild[name] = query
	return tbq
}

// TimingBreakdownGroupBy is the group-by builder for TimingBreakdown entities.
type TimingBreakdownGroupBy struct {
	selector
	build *TimingBreakdownQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tbgb *TimingBreakdownGroupBy) Aggregate(fns ...AggregateFunc) *TimingBreakdownGroupBy {
	tbgb.fns = append(tbgb.fns, fns...)
	return tbgb
}

// Scan applies the selector query and scans the result into the given value.
func (tbgb *TimingBreakdownGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tbgb.build.ctx, "GroupBy")
	if err := tbgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TimingBreakdownQuery, *TimingBreakdownGroupBy](ctx, tbgb.build, tbgb, tbgb.build.inters, v)
}

func (tbgb *TimingBreakdownGroupBy) sqlScan(ctx context.Context, root *TimingBreakdownQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(tbgb.fns))
	for _, fn := range tbgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*tbgb.flds)+len(tbgb.fns))
		for _, f := range *tbgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*tbgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tbgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// TimingBreakdownSelect is the builder for selecting fields of TimingBreakdown entities.
type TimingBreakdownSelect struct {
	*TimingBreakdownQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (tbs *TimingBreakdownSelect) Aggregate(fns ...AggregateFunc) *TimingBreakdownSelect {
	tbs.fns = append(tbs.fns, fns...)
	return tbs
}

// Scan applies the selector query and scans the result into the given value.
func (tbs *TimingBreakdownSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tbs.ctx, "Select")
	if err := tbs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TimingBreakdownQuery, *TimingBreakdownSelect](ctx, tbs.TimingBreakdownQuery, tbs, tbs.inters, v)
}

func (tbs *TimingBreakdownSelect) sqlScan(ctx context.Context, root *TimingBreakdownQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(tbs.fns))
	for _, fn := range tbs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*tbs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tbs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}