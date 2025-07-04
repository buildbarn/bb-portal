// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/actionsummary"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/runnercount"
)

// RunnerCountQuery is the builder for querying RunnerCount entities.
type RunnerCountQuery struct {
	config
	ctx               *QueryContext
	order             []runnercount.OrderOption
	inters            []Interceptor
	predicates        []predicate.RunnerCount
	withActionSummary *ActionSummaryQuery
	withFKs           bool
	modifiers         []func(*sql.Selector)
	loadTotal         []func(context.Context, []*RunnerCount) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RunnerCountQuery builder.
func (rcq *RunnerCountQuery) Where(ps ...predicate.RunnerCount) *RunnerCountQuery {
	rcq.predicates = append(rcq.predicates, ps...)
	return rcq
}

// Limit the number of records to be returned by this query.
func (rcq *RunnerCountQuery) Limit(limit int) *RunnerCountQuery {
	rcq.ctx.Limit = &limit
	return rcq
}

// Offset to start from.
func (rcq *RunnerCountQuery) Offset(offset int) *RunnerCountQuery {
	rcq.ctx.Offset = &offset
	return rcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rcq *RunnerCountQuery) Unique(unique bool) *RunnerCountQuery {
	rcq.ctx.Unique = &unique
	return rcq
}

// Order specifies how the records should be ordered.
func (rcq *RunnerCountQuery) Order(o ...runnercount.OrderOption) *RunnerCountQuery {
	rcq.order = append(rcq.order, o...)
	return rcq
}

// QueryActionSummary chains the current query on the "action_summary" edge.
func (rcq *RunnerCountQuery) QueryActionSummary() *ActionSummaryQuery {
	query := (&ActionSummaryClient{config: rcq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(runnercount.Table, runnercount.FieldID, selector),
			sqlgraph.To(actionsummary.Table, actionsummary.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, runnercount.ActionSummaryTable, runnercount.ActionSummaryColumn),
		)
		fromU = sqlgraph.SetNeighbors(rcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first RunnerCount entity from the query.
// Returns a *NotFoundError when no RunnerCount was found.
func (rcq *RunnerCountQuery) First(ctx context.Context) (*RunnerCount, error) {
	nodes, err := rcq.Limit(1).All(setContextOp(ctx, rcq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{runnercount.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rcq *RunnerCountQuery) FirstX(ctx context.Context) *RunnerCount {
	node, err := rcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first RunnerCount ID from the query.
// Returns a *NotFoundError when no RunnerCount ID was found.
func (rcq *RunnerCountQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rcq.Limit(1).IDs(setContextOp(ctx, rcq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{runnercount.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rcq *RunnerCountQuery) FirstIDX(ctx context.Context) int {
	id, err := rcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single RunnerCount entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one RunnerCount entity is found.
// Returns a *NotFoundError when no RunnerCount entities are found.
func (rcq *RunnerCountQuery) Only(ctx context.Context) (*RunnerCount, error) {
	nodes, err := rcq.Limit(2).All(setContextOp(ctx, rcq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{runnercount.Label}
	default:
		return nil, &NotSingularError{runnercount.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rcq *RunnerCountQuery) OnlyX(ctx context.Context) *RunnerCount {
	node, err := rcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only RunnerCount ID in the query.
// Returns a *NotSingularError when more than one RunnerCount ID is found.
// Returns a *NotFoundError when no entities are found.
func (rcq *RunnerCountQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = rcq.Limit(2).IDs(setContextOp(ctx, rcq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{runnercount.Label}
	default:
		err = &NotSingularError{runnercount.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rcq *RunnerCountQuery) OnlyIDX(ctx context.Context) int {
	id, err := rcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of RunnerCounts.
func (rcq *RunnerCountQuery) All(ctx context.Context) ([]*RunnerCount, error) {
	ctx = setContextOp(ctx, rcq.ctx, ent.OpQueryAll)
	if err := rcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*RunnerCount, *RunnerCountQuery]()
	return withInterceptors[[]*RunnerCount](ctx, rcq, qr, rcq.inters)
}

// AllX is like All, but panics if an error occurs.
func (rcq *RunnerCountQuery) AllX(ctx context.Context) []*RunnerCount {
	nodes, err := rcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of RunnerCount IDs.
func (rcq *RunnerCountQuery) IDs(ctx context.Context) (ids []int, err error) {
	if rcq.ctx.Unique == nil && rcq.path != nil {
		rcq.Unique(true)
	}
	ctx = setContextOp(ctx, rcq.ctx, ent.OpQueryIDs)
	if err = rcq.Select(runnercount.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rcq *RunnerCountQuery) IDsX(ctx context.Context) []int {
	ids, err := rcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rcq *RunnerCountQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, rcq.ctx, ent.OpQueryCount)
	if err := rcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, rcq, querierCount[*RunnerCountQuery](), rcq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (rcq *RunnerCountQuery) CountX(ctx context.Context) int {
	count, err := rcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rcq *RunnerCountQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, rcq.ctx, ent.OpQueryExist)
	switch _, err := rcq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (rcq *RunnerCountQuery) ExistX(ctx context.Context) bool {
	exist, err := rcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RunnerCountQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rcq *RunnerCountQuery) Clone() *RunnerCountQuery {
	if rcq == nil {
		return nil
	}
	return &RunnerCountQuery{
		config:            rcq.config,
		ctx:               rcq.ctx.Clone(),
		order:             append([]runnercount.OrderOption{}, rcq.order...),
		inters:            append([]Interceptor{}, rcq.inters...),
		predicates:        append([]predicate.RunnerCount{}, rcq.predicates...),
		withActionSummary: rcq.withActionSummary.Clone(),
		// clone intermediate query.
		sql:  rcq.sql.Clone(),
		path: rcq.path,
	}
}

// WithActionSummary tells the query-builder to eager-load the nodes that are connected to
// the "action_summary" edge. The optional arguments are used to configure the query builder of the edge.
func (rcq *RunnerCountQuery) WithActionSummary(opts ...func(*ActionSummaryQuery)) *RunnerCountQuery {
	query := (&ActionSummaryClient{config: rcq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rcq.withActionSummary = query
	return rcq
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
//	client.RunnerCount.Query().
//		GroupBy(runnercount.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (rcq *RunnerCountQuery) GroupBy(field string, fields ...string) *RunnerCountGroupBy {
	rcq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &RunnerCountGroupBy{build: rcq}
	grbuild.flds = &rcq.ctx.Fields
	grbuild.label = runnercount.Label
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
//	client.RunnerCount.Query().
//		Select(runnercount.FieldName).
//		Scan(ctx, &v)
func (rcq *RunnerCountQuery) Select(fields ...string) *RunnerCountSelect {
	rcq.ctx.Fields = append(rcq.ctx.Fields, fields...)
	sbuild := &RunnerCountSelect{RunnerCountQuery: rcq}
	sbuild.label = runnercount.Label
	sbuild.flds, sbuild.scan = &rcq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a RunnerCountSelect configured with the given aggregations.
func (rcq *RunnerCountQuery) Aggregate(fns ...AggregateFunc) *RunnerCountSelect {
	return rcq.Select().Aggregate(fns...)
}

func (rcq *RunnerCountQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range rcq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, rcq); err != nil {
				return err
			}
		}
	}
	for _, f := range rcq.ctx.Fields {
		if !runnercount.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rcq.path != nil {
		prev, err := rcq.path(ctx)
		if err != nil {
			return err
		}
		rcq.sql = prev
	}
	return nil
}

func (rcq *RunnerCountQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*RunnerCount, error) {
	var (
		nodes       = []*RunnerCount{}
		withFKs     = rcq.withFKs
		_spec       = rcq.querySpec()
		loadedTypes = [1]bool{
			rcq.withActionSummary != nil,
		}
	)
	if rcq.withActionSummary != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, runnercount.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*RunnerCount).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &RunnerCount{config: rcq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(rcq.modifiers) > 0 {
		_spec.Modifiers = rcq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rcq.withActionSummary; query != nil {
		if err := rcq.loadActionSummary(ctx, query, nodes, nil,
			func(n *RunnerCount, e *ActionSummary) { n.Edges.ActionSummary = e }); err != nil {
			return nil, err
		}
	}
	for i := range rcq.loadTotal {
		if err := rcq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rcq *RunnerCountQuery) loadActionSummary(ctx context.Context, query *ActionSummaryQuery, nodes []*RunnerCount, init func(*RunnerCount), assign func(*RunnerCount, *ActionSummary)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*RunnerCount)
	for i := range nodes {
		if nodes[i].action_summary_runner_count == nil {
			continue
		}
		fk := *nodes[i].action_summary_runner_count
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(actionsummary.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "action_summary_runner_count" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (rcq *RunnerCountQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rcq.querySpec()
	if len(rcq.modifiers) > 0 {
		_spec.Modifiers = rcq.modifiers
	}
	_spec.Node.Columns = rcq.ctx.Fields
	if len(rcq.ctx.Fields) > 0 {
		_spec.Unique = rcq.ctx.Unique != nil && *rcq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, rcq.driver, _spec)
}

func (rcq *RunnerCountQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(runnercount.Table, runnercount.Columns, sqlgraph.NewFieldSpec(runnercount.FieldID, field.TypeInt))
	_spec.From = rcq.sql
	if unique := rcq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if rcq.path != nil {
		_spec.Unique = true
	}
	if fields := rcq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, runnercount.FieldID)
		for i := range fields {
			if fields[i] != runnercount.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rcq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rcq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rcq *RunnerCountQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rcq.driver.Dialect())
	t1 := builder.Table(runnercount.Table)
	columns := rcq.ctx.Fields
	if len(columns) == 0 {
		columns = runnercount.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rcq.sql != nil {
		selector = rcq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rcq.ctx.Unique != nil && *rcq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range rcq.predicates {
		p(selector)
	}
	for _, p := range rcq.order {
		p(selector)
	}
	if offset := rcq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rcq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// RunnerCountGroupBy is the group-by builder for RunnerCount entities.
type RunnerCountGroupBy struct {
	selector
	build *RunnerCountQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rcgb *RunnerCountGroupBy) Aggregate(fns ...AggregateFunc) *RunnerCountGroupBy {
	rcgb.fns = append(rcgb.fns, fns...)
	return rcgb
}

// Scan applies the selector query and scans the result into the given value.
func (rcgb *RunnerCountGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rcgb.build.ctx, ent.OpQueryGroupBy)
	if err := rcgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RunnerCountQuery, *RunnerCountGroupBy](ctx, rcgb.build, rcgb, rcgb.build.inters, v)
}

func (rcgb *RunnerCountGroupBy) sqlScan(ctx context.Context, root *RunnerCountQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(rcgb.fns))
	for _, fn := range rcgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*rcgb.flds)+len(rcgb.fns))
		for _, f := range *rcgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*rcgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rcgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// RunnerCountSelect is the builder for selecting fields of RunnerCount entities.
type RunnerCountSelect struct {
	*RunnerCountQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rcs *RunnerCountSelect) Aggregate(fns ...AggregateFunc) *RunnerCountSelect {
	rcs.fns = append(rcs.fns, fns...)
	return rcs
}

// Scan applies the selector query and scans the result into the given value.
func (rcs *RunnerCountSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rcs.ctx, ent.OpQuerySelect)
	if err := rcs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*RunnerCountQuery, *RunnerCountSelect](ctx, rcs.RunnerCountQuery, rcs, rcs.inters, v)
}

func (rcs *RunnerCountSelect) sqlScan(ctx context.Context, root *RunnerCountQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(rcs.fns))
	for _, fn := range rcs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*rcs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
