// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetcomplete"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetconfigured"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetpair"
)

// TargetPairQuery is the builder for querying TargetPair entities.
type TargetPairQuery struct {
	config
	ctx                 *QueryContext
	order               []targetpair.OrderOption
	inters              []Interceptor
	predicates          []predicate.TargetPair
	withBazelInvocation *BazelInvocationQuery
	withConfiguration   *TargetConfiguredQuery
	withCompletion      *TargetCompleteQuery
	withFKs             bool
	modifiers           []func(*sql.Selector)
	loadTotal           []func(context.Context, []*TargetPair) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TargetPairQuery builder.
func (tpq *TargetPairQuery) Where(ps ...predicate.TargetPair) *TargetPairQuery {
	tpq.predicates = append(tpq.predicates, ps...)
	return tpq
}

// Limit the number of records to be returned by this query.
func (tpq *TargetPairQuery) Limit(limit int) *TargetPairQuery {
	tpq.ctx.Limit = &limit
	return tpq
}

// Offset to start from.
func (tpq *TargetPairQuery) Offset(offset int) *TargetPairQuery {
	tpq.ctx.Offset = &offset
	return tpq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tpq *TargetPairQuery) Unique(unique bool) *TargetPairQuery {
	tpq.ctx.Unique = &unique
	return tpq
}

// Order specifies how the records should be ordered.
func (tpq *TargetPairQuery) Order(o ...targetpair.OrderOption) *TargetPairQuery {
	tpq.order = append(tpq.order, o...)
	return tpq
}

// QueryBazelInvocation chains the current query on the "bazel_invocation" edge.
func (tpq *TargetPairQuery) QueryBazelInvocation() *BazelInvocationQuery {
	query := (&BazelInvocationClient{config: tpq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(targetpair.Table, targetpair.FieldID, selector),
			sqlgraph.To(bazelinvocation.Table, bazelinvocation.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, targetpair.BazelInvocationTable, targetpair.BazelInvocationColumn),
		)
		fromU = sqlgraph.SetNeighbors(tpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryConfiguration chains the current query on the "configuration" edge.
func (tpq *TargetPairQuery) QueryConfiguration() *TargetConfiguredQuery {
	query := (&TargetConfiguredClient{config: tpq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(targetpair.Table, targetpair.FieldID, selector),
			sqlgraph.To(targetconfigured.Table, targetconfigured.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, targetpair.ConfigurationTable, targetpair.ConfigurationColumn),
		)
		fromU = sqlgraph.SetNeighbors(tpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryCompletion chains the current query on the "completion" edge.
func (tpq *TargetPairQuery) QueryCompletion() *TargetCompleteQuery {
	query := (&TargetCompleteClient{config: tpq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tpq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tpq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(targetpair.Table, targetpair.FieldID, selector),
			sqlgraph.To(targetcomplete.Table, targetcomplete.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, targetpair.CompletionTable, targetpair.CompletionColumn),
		)
		fromU = sqlgraph.SetNeighbors(tpq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first TargetPair entity from the query.
// Returns a *NotFoundError when no TargetPair was found.
func (tpq *TargetPairQuery) First(ctx context.Context) (*TargetPair, error) {
	nodes, err := tpq.Limit(1).All(setContextOp(ctx, tpq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{targetpair.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tpq *TargetPairQuery) FirstX(ctx context.Context) *TargetPair {
	node, err := tpq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TargetPair ID from the query.
// Returns a *NotFoundError when no TargetPair ID was found.
func (tpq *TargetPairQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tpq.Limit(1).IDs(setContextOp(ctx, tpq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{targetpair.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tpq *TargetPairQuery) FirstIDX(ctx context.Context) int {
	id, err := tpq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TargetPair entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one TargetPair entity is found.
// Returns a *NotFoundError when no TargetPair entities are found.
func (tpq *TargetPairQuery) Only(ctx context.Context) (*TargetPair, error) {
	nodes, err := tpq.Limit(2).All(setContextOp(ctx, tpq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{targetpair.Label}
	default:
		return nil, &NotSingularError{targetpair.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tpq *TargetPairQuery) OnlyX(ctx context.Context) *TargetPair {
	node, err := tpq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TargetPair ID in the query.
// Returns a *NotSingularError when more than one TargetPair ID is found.
// Returns a *NotFoundError when no entities are found.
func (tpq *TargetPairQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tpq.Limit(2).IDs(setContextOp(ctx, tpq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{targetpair.Label}
	default:
		err = &NotSingularError{targetpair.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tpq *TargetPairQuery) OnlyIDX(ctx context.Context) int {
	id, err := tpq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TargetPairs.
func (tpq *TargetPairQuery) All(ctx context.Context) ([]*TargetPair, error) {
	ctx = setContextOp(ctx, tpq.ctx, ent.OpQueryAll)
	if err := tpq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*TargetPair, *TargetPairQuery]()
	return withInterceptors[[]*TargetPair](ctx, tpq, qr, tpq.inters)
}

// AllX is like All, but panics if an error occurs.
func (tpq *TargetPairQuery) AllX(ctx context.Context) []*TargetPair {
	nodes, err := tpq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TargetPair IDs.
func (tpq *TargetPairQuery) IDs(ctx context.Context) (ids []int, err error) {
	if tpq.ctx.Unique == nil && tpq.path != nil {
		tpq.Unique(true)
	}
	ctx = setContextOp(ctx, tpq.ctx, ent.OpQueryIDs)
	if err = tpq.Select(targetpair.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tpq *TargetPairQuery) IDsX(ctx context.Context) []int {
	ids, err := tpq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tpq *TargetPairQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, tpq.ctx, ent.OpQueryCount)
	if err := tpq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, tpq, querierCount[*TargetPairQuery](), tpq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (tpq *TargetPairQuery) CountX(ctx context.Context) int {
	count, err := tpq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tpq *TargetPairQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, tpq.ctx, ent.OpQueryExist)
	switch _, err := tpq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (tpq *TargetPairQuery) ExistX(ctx context.Context) bool {
	exist, err := tpq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TargetPairQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tpq *TargetPairQuery) Clone() *TargetPairQuery {
	if tpq == nil {
		return nil
	}
	return &TargetPairQuery{
		config:              tpq.config,
		ctx:                 tpq.ctx.Clone(),
		order:               append([]targetpair.OrderOption{}, tpq.order...),
		inters:              append([]Interceptor{}, tpq.inters...),
		predicates:          append([]predicate.TargetPair{}, tpq.predicates...),
		withBazelInvocation: tpq.withBazelInvocation.Clone(),
		withConfiguration:   tpq.withConfiguration.Clone(),
		withCompletion:      tpq.withCompletion.Clone(),
		// clone intermediate query.
		sql:  tpq.sql.Clone(),
		path: tpq.path,
	}
}

// WithBazelInvocation tells the query-builder to eager-load the nodes that are connected to
// the "bazel_invocation" edge. The optional arguments are used to configure the query builder of the edge.
func (tpq *TargetPairQuery) WithBazelInvocation(opts ...func(*BazelInvocationQuery)) *TargetPairQuery {
	query := (&BazelInvocationClient{config: tpq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tpq.withBazelInvocation = query
	return tpq
}

// WithConfiguration tells the query-builder to eager-load the nodes that are connected to
// the "configuration" edge. The optional arguments are used to configure the query builder of the edge.
func (tpq *TargetPairQuery) WithConfiguration(opts ...func(*TargetConfiguredQuery)) *TargetPairQuery {
	query := (&TargetConfiguredClient{config: tpq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tpq.withConfiguration = query
	return tpq
}

// WithCompletion tells the query-builder to eager-load the nodes that are connected to
// the "completion" edge. The optional arguments are used to configure the query builder of the edge.
func (tpq *TargetPairQuery) WithCompletion(opts ...func(*TargetCompleteQuery)) *TargetPairQuery {
	query := (&TargetCompleteClient{config: tpq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tpq.withCompletion = query
	return tpq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Label string `json:"label,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.TargetPair.Query().
//		GroupBy(targetpair.FieldLabel).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (tpq *TargetPairQuery) GroupBy(field string, fields ...string) *TargetPairGroupBy {
	tpq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &TargetPairGroupBy{build: tpq}
	grbuild.flds = &tpq.ctx.Fields
	grbuild.label = targetpair.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Label string `json:"label,omitempty"`
//	}
//
//	client.TargetPair.Query().
//		Select(targetpair.FieldLabel).
//		Scan(ctx, &v)
func (tpq *TargetPairQuery) Select(fields ...string) *TargetPairSelect {
	tpq.ctx.Fields = append(tpq.ctx.Fields, fields...)
	sbuild := &TargetPairSelect{TargetPairQuery: tpq}
	sbuild.label = targetpair.Label
	sbuild.flds, sbuild.scan = &tpq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a TargetPairSelect configured with the given aggregations.
func (tpq *TargetPairQuery) Aggregate(fns ...AggregateFunc) *TargetPairSelect {
	return tpq.Select().Aggregate(fns...)
}

func (tpq *TargetPairQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range tpq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, tpq); err != nil {
				return err
			}
		}
	}
	for _, f := range tpq.ctx.Fields {
		if !targetpair.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tpq.path != nil {
		prev, err := tpq.path(ctx)
		if err != nil {
			return err
		}
		tpq.sql = prev
	}
	return nil
}

func (tpq *TargetPairQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*TargetPair, error) {
	var (
		nodes       = []*TargetPair{}
		withFKs     = tpq.withFKs
		_spec       = tpq.querySpec()
		loadedTypes = [3]bool{
			tpq.withBazelInvocation != nil,
			tpq.withConfiguration != nil,
			tpq.withCompletion != nil,
		}
	)
	if tpq.withBazelInvocation != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, targetpair.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*TargetPair).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &TargetPair{config: tpq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(tpq.modifiers) > 0 {
		_spec.Modifiers = tpq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tpq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := tpq.withBazelInvocation; query != nil {
		if err := tpq.loadBazelInvocation(ctx, query, nodes, nil,
			func(n *TargetPair, e *BazelInvocation) { n.Edges.BazelInvocation = e }); err != nil {
			return nil, err
		}
	}
	if query := tpq.withConfiguration; query != nil {
		if err := tpq.loadConfiguration(ctx, query, nodes, nil,
			func(n *TargetPair, e *TargetConfigured) { n.Edges.Configuration = e }); err != nil {
			return nil, err
		}
	}
	if query := tpq.withCompletion; query != nil {
		if err := tpq.loadCompletion(ctx, query, nodes, nil,
			func(n *TargetPair, e *TargetComplete) { n.Edges.Completion = e }); err != nil {
			return nil, err
		}
	}
	for i := range tpq.loadTotal {
		if err := tpq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (tpq *TargetPairQuery) loadBazelInvocation(ctx context.Context, query *BazelInvocationQuery, nodes []*TargetPair, init func(*TargetPair), assign func(*TargetPair, *BazelInvocation)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*TargetPair)
	for i := range nodes {
		if nodes[i].bazel_invocation_targets == nil {
			continue
		}
		fk := *nodes[i].bazel_invocation_targets
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(bazelinvocation.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "bazel_invocation_targets" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (tpq *TargetPairQuery) loadConfiguration(ctx context.Context, query *TargetConfiguredQuery, nodes []*TargetPair, init func(*TargetPair), assign func(*TargetPair, *TargetConfigured)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*TargetPair)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.TargetConfigured(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(targetpair.ConfigurationColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.target_pair_configuration
		if fk == nil {
			return fmt.Errorf(`foreign-key "target_pair_configuration" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "target_pair_configuration" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (tpq *TargetPairQuery) loadCompletion(ctx context.Context, query *TargetCompleteQuery, nodes []*TargetPair, init func(*TargetPair), assign func(*TargetPair, *TargetComplete)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*TargetPair)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.TargetComplete(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(targetpair.CompletionColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.target_pair_completion
		if fk == nil {
			return fmt.Errorf(`foreign-key "target_pair_completion" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "target_pair_completion" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (tpq *TargetPairQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tpq.querySpec()
	if len(tpq.modifiers) > 0 {
		_spec.Modifiers = tpq.modifiers
	}
	_spec.Node.Columns = tpq.ctx.Fields
	if len(tpq.ctx.Fields) > 0 {
		_spec.Unique = tpq.ctx.Unique != nil && *tpq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, tpq.driver, _spec)
}

func (tpq *TargetPairQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(targetpair.Table, targetpair.Columns, sqlgraph.NewFieldSpec(targetpair.FieldID, field.TypeInt))
	_spec.From = tpq.sql
	if unique := tpq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if tpq.path != nil {
		_spec.Unique = true
	}
	if fields := tpq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, targetpair.FieldID)
		for i := range fields {
			if fields[i] != targetpair.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tpq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tpq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tpq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tpq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tpq *TargetPairQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tpq.driver.Dialect())
	t1 := builder.Table(targetpair.Table)
	columns := tpq.ctx.Fields
	if len(columns) == 0 {
		columns = targetpair.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tpq.sql != nil {
		selector = tpq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tpq.ctx.Unique != nil && *tpq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range tpq.predicates {
		p(selector)
	}
	for _, p := range tpq.order {
		p(selector)
	}
	if offset := tpq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tpq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TargetPairGroupBy is the group-by builder for TargetPair entities.
type TargetPairGroupBy struct {
	selector
	build *TargetPairQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tpgb *TargetPairGroupBy) Aggregate(fns ...AggregateFunc) *TargetPairGroupBy {
	tpgb.fns = append(tpgb.fns, fns...)
	return tpgb
}

// Scan applies the selector query and scans the result into the given value.
func (tpgb *TargetPairGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tpgb.build.ctx, ent.OpQueryGroupBy)
	if err := tpgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TargetPairQuery, *TargetPairGroupBy](ctx, tpgb.build, tpgb, tpgb.build.inters, v)
}

func (tpgb *TargetPairGroupBy) sqlScan(ctx context.Context, root *TargetPairQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(tpgb.fns))
	for _, fn := range tpgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*tpgb.flds)+len(tpgb.fns))
		for _, f := range *tpgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*tpgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tpgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// TargetPairSelect is the builder for selecting fields of TargetPair entities.
type TargetPairSelect struct {
	*TargetPairQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (tps *TargetPairSelect) Aggregate(fns ...AggregateFunc) *TargetPairSelect {
	tps.fns = append(tps.fns, fns...)
	return tps
}

// Scan applies the selector query and scans the result into the given value.
func (tps *TargetPairSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tps.ctx, ent.OpQuerySelect)
	if err := tps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TargetPairQuery, *TargetPairSelect](ctx, tps.TargetPairQuery, tps, tps.inters, v)
}

func (tps *TargetPairSelect) sqlScan(ctx context.Context, root *TargetPairQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(tps.fns))
	for _, fn := range tps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*tps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
