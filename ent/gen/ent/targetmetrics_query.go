// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/metrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/targetmetrics"
)

// TargetMetricsQuery is the builder for querying TargetMetrics entities.
type TargetMetricsQuery struct {
	config
	ctx         *QueryContext
	order       []targetmetrics.OrderOption
	inters      []Interceptor
	predicates  []predicate.TargetMetrics
	withMetrics *MetricsQuery
	modifiers   []func(*sql.Selector)
	loadTotal   []func(context.Context, []*TargetMetrics) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TargetMetricsQuery builder.
func (tmq *TargetMetricsQuery) Where(ps ...predicate.TargetMetrics) *TargetMetricsQuery {
	tmq.predicates = append(tmq.predicates, ps...)
	return tmq
}

// Limit the number of records to be returned by this query.
func (tmq *TargetMetricsQuery) Limit(limit int) *TargetMetricsQuery {
	tmq.ctx.Limit = &limit
	return tmq
}

// Offset to start from.
func (tmq *TargetMetricsQuery) Offset(offset int) *TargetMetricsQuery {
	tmq.ctx.Offset = &offset
	return tmq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tmq *TargetMetricsQuery) Unique(unique bool) *TargetMetricsQuery {
	tmq.ctx.Unique = &unique
	return tmq
}

// Order specifies how the records should be ordered.
func (tmq *TargetMetricsQuery) Order(o ...targetmetrics.OrderOption) *TargetMetricsQuery {
	tmq.order = append(tmq.order, o...)
	return tmq
}

// QueryMetrics chains the current query on the "metrics" edge.
func (tmq *TargetMetricsQuery) QueryMetrics() *MetricsQuery {
	query := (&MetricsClient{config: tmq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tmq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(targetmetrics.Table, targetmetrics.FieldID, selector),
			sqlgraph.To(metrics.Table, metrics.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, targetmetrics.MetricsTable, targetmetrics.MetricsColumn),
		)
		fromU = sqlgraph.SetNeighbors(tmq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first TargetMetrics entity from the query.
// Returns a *NotFoundError when no TargetMetrics was found.
func (tmq *TargetMetricsQuery) First(ctx context.Context) (*TargetMetrics, error) {
	nodes, err := tmq.Limit(1).All(setContextOp(ctx, tmq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{targetmetrics.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tmq *TargetMetricsQuery) FirstX(ctx context.Context) *TargetMetrics {
	node, err := tmq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TargetMetrics ID from the query.
// Returns a *NotFoundError when no TargetMetrics ID was found.
func (tmq *TargetMetricsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tmq.Limit(1).IDs(setContextOp(ctx, tmq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{targetmetrics.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tmq *TargetMetricsQuery) FirstIDX(ctx context.Context) int {
	id, err := tmq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TargetMetrics entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one TargetMetrics entity is found.
// Returns a *NotFoundError when no TargetMetrics entities are found.
func (tmq *TargetMetricsQuery) Only(ctx context.Context) (*TargetMetrics, error) {
	nodes, err := tmq.Limit(2).All(setContextOp(ctx, tmq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{targetmetrics.Label}
	default:
		return nil, &NotSingularError{targetmetrics.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tmq *TargetMetricsQuery) OnlyX(ctx context.Context) *TargetMetrics {
	node, err := tmq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TargetMetrics ID in the query.
// Returns a *NotSingularError when more than one TargetMetrics ID is found.
// Returns a *NotFoundError when no entities are found.
func (tmq *TargetMetricsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tmq.Limit(2).IDs(setContextOp(ctx, tmq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{targetmetrics.Label}
	default:
		err = &NotSingularError{targetmetrics.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tmq *TargetMetricsQuery) OnlyIDX(ctx context.Context) int {
	id, err := tmq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TargetMetricsSlice.
func (tmq *TargetMetricsQuery) All(ctx context.Context) ([]*TargetMetrics, error) {
	ctx = setContextOp(ctx, tmq.ctx, "All")
	if err := tmq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*TargetMetrics, *TargetMetricsQuery]()
	return withInterceptors[[]*TargetMetrics](ctx, tmq, qr, tmq.inters)
}

// AllX is like All, but panics if an error occurs.
func (tmq *TargetMetricsQuery) AllX(ctx context.Context) []*TargetMetrics {
	nodes, err := tmq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TargetMetrics IDs.
func (tmq *TargetMetricsQuery) IDs(ctx context.Context) (ids []int, err error) {
	if tmq.ctx.Unique == nil && tmq.path != nil {
		tmq.Unique(true)
	}
	ctx = setContextOp(ctx, tmq.ctx, "IDs")
	if err = tmq.Select(targetmetrics.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tmq *TargetMetricsQuery) IDsX(ctx context.Context) []int {
	ids, err := tmq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tmq *TargetMetricsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, tmq.ctx, "Count")
	if err := tmq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, tmq, querierCount[*TargetMetricsQuery](), tmq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (tmq *TargetMetricsQuery) CountX(ctx context.Context) int {
	count, err := tmq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tmq *TargetMetricsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, tmq.ctx, "Exist")
	switch _, err := tmq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (tmq *TargetMetricsQuery) ExistX(ctx context.Context) bool {
	exist, err := tmq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TargetMetricsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tmq *TargetMetricsQuery) Clone() *TargetMetricsQuery {
	if tmq == nil {
		return nil
	}
	return &TargetMetricsQuery{
		config:      tmq.config,
		ctx:         tmq.ctx.Clone(),
		order:       append([]targetmetrics.OrderOption{}, tmq.order...),
		inters:      append([]Interceptor{}, tmq.inters...),
		predicates:  append([]predicate.TargetMetrics{}, tmq.predicates...),
		withMetrics: tmq.withMetrics.Clone(),
		// clone intermediate query.
		sql:  tmq.sql.Clone(),
		path: tmq.path,
	}
}

// WithMetrics tells the query-builder to eager-load the nodes that are connected to
// the "metrics" edge. The optional arguments are used to configure the query builder of the edge.
func (tmq *TargetMetricsQuery) WithMetrics(opts ...func(*MetricsQuery)) *TargetMetricsQuery {
	query := (&MetricsClient{config: tmq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tmq.withMetrics = query
	return tmq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		TargetsLoaded int64 `json:"targets_loaded,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.TargetMetrics.Query().
//		GroupBy(targetmetrics.FieldTargetsLoaded).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (tmq *TargetMetricsQuery) GroupBy(field string, fields ...string) *TargetMetricsGroupBy {
	tmq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &TargetMetricsGroupBy{build: tmq}
	grbuild.flds = &tmq.ctx.Fields
	grbuild.label = targetmetrics.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		TargetsLoaded int64 `json:"targets_loaded,omitempty"`
//	}
//
//	client.TargetMetrics.Query().
//		Select(targetmetrics.FieldTargetsLoaded).
//		Scan(ctx, &v)
func (tmq *TargetMetricsQuery) Select(fields ...string) *TargetMetricsSelect {
	tmq.ctx.Fields = append(tmq.ctx.Fields, fields...)
	sbuild := &TargetMetricsSelect{TargetMetricsQuery: tmq}
	sbuild.label = targetmetrics.Label
	sbuild.flds, sbuild.scan = &tmq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a TargetMetricsSelect configured with the given aggregations.
func (tmq *TargetMetricsQuery) Aggregate(fns ...AggregateFunc) *TargetMetricsSelect {
	return tmq.Select().Aggregate(fns...)
}

func (tmq *TargetMetricsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range tmq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, tmq); err != nil {
				return err
			}
		}
	}
	for _, f := range tmq.ctx.Fields {
		if !targetmetrics.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tmq.path != nil {
		prev, err := tmq.path(ctx)
		if err != nil {
			return err
		}
		tmq.sql = prev
	}
	return nil
}

func (tmq *TargetMetricsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*TargetMetrics, error) {
	var (
		nodes       = []*TargetMetrics{}
		_spec       = tmq.querySpec()
		loadedTypes = [1]bool{
			tmq.withMetrics != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*TargetMetrics).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &TargetMetrics{config: tmq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(tmq.modifiers) > 0 {
		_spec.Modifiers = tmq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tmq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := tmq.withMetrics; query != nil {
		if err := tmq.loadMetrics(ctx, query, nodes, nil,
			func(n *TargetMetrics, e *Metrics) { n.Edges.Metrics = e }); err != nil {
			return nil, err
		}
	}
	for i := range tmq.loadTotal {
		if err := tmq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (tmq *TargetMetricsQuery) loadMetrics(ctx context.Context, query *MetricsQuery, nodes []*TargetMetrics, init func(*TargetMetrics), assign func(*TargetMetrics, *Metrics)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*TargetMetrics)
	for i := range nodes {
		fk := nodes[i].MetricsID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(metrics.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "metrics_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (tmq *TargetMetricsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tmq.querySpec()
	if len(tmq.modifiers) > 0 {
		_spec.Modifiers = tmq.modifiers
	}
	_spec.Node.Columns = tmq.ctx.Fields
	if len(tmq.ctx.Fields) > 0 {
		_spec.Unique = tmq.ctx.Unique != nil && *tmq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, tmq.driver, _spec)
}

func (tmq *TargetMetricsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(targetmetrics.Table, targetmetrics.Columns, sqlgraph.NewFieldSpec(targetmetrics.FieldID, field.TypeInt))
	_spec.From = tmq.sql
	if unique := tmq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if tmq.path != nil {
		_spec.Unique = true
	}
	if fields := tmq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, targetmetrics.FieldID)
		for i := range fields {
			if fields[i] != targetmetrics.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if tmq.withMetrics != nil {
			_spec.Node.AddColumnOnce(targetmetrics.FieldMetricsID)
		}
	}
	if ps := tmq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tmq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tmq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tmq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tmq *TargetMetricsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tmq.driver.Dialect())
	t1 := builder.Table(targetmetrics.Table)
	columns := tmq.ctx.Fields
	if len(columns) == 0 {
		columns = targetmetrics.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tmq.sql != nil {
		selector = tmq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tmq.ctx.Unique != nil && *tmq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range tmq.predicates {
		p(selector)
	}
	for _, p := range tmq.order {
		p(selector)
	}
	if offset := tmq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tmq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TargetMetricsGroupBy is the group-by builder for TargetMetrics entities.
type TargetMetricsGroupBy struct {
	selector
	build *TargetMetricsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tmgb *TargetMetricsGroupBy) Aggregate(fns ...AggregateFunc) *TargetMetricsGroupBy {
	tmgb.fns = append(tmgb.fns, fns...)
	return tmgb
}

// Scan applies the selector query and scans the result into the given value.
func (tmgb *TargetMetricsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tmgb.build.ctx, "GroupBy")
	if err := tmgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TargetMetricsQuery, *TargetMetricsGroupBy](ctx, tmgb.build, tmgb, tmgb.build.inters, v)
}

func (tmgb *TargetMetricsGroupBy) sqlScan(ctx context.Context, root *TargetMetricsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(tmgb.fns))
	for _, fn := range tmgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*tmgb.flds)+len(tmgb.fns))
		for _, f := range *tmgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*tmgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tmgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// TargetMetricsSelect is the builder for selecting fields of TargetMetrics entities.
type TargetMetricsSelect struct {
	*TargetMetricsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (tms *TargetMetricsSelect) Aggregate(fns ...AggregateFunc) *TargetMetricsSelect {
	tms.fns = append(tms.fns, fns...)
	return tms
}

// Scan applies the selector query and scans the result into the given value.
func (tms *TargetMetricsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tms.ctx, "Select")
	if err := tms.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TargetMetricsQuery, *TargetMetricsSelect](ctx, tms.TargetMetricsQuery, tms, tms.inters, v)
}

func (tms *TargetMetricsSelect) sqlScan(ctx context.Context, root *TargetMetricsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(tms.fns))
	for _, fn := range tms.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*tms.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
