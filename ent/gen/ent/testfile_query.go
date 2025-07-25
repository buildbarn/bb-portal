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
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testfile"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testresultbes"
)

// TestFileQuery is the builder for querying TestFile entities.
type TestFileQuery struct {
	config
	ctx            *QueryContext
	order          []testfile.OrderOption
	inters         []Interceptor
	predicates     []predicate.TestFile
	withTestResult *TestResultBESQuery
	withFKs        bool
	modifiers      []func(*sql.Selector)
	loadTotal      []func(context.Context, []*TestFile) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TestFileQuery builder.
func (tfq *TestFileQuery) Where(ps ...predicate.TestFile) *TestFileQuery {
	tfq.predicates = append(tfq.predicates, ps...)
	return tfq
}

// Limit the number of records to be returned by this query.
func (tfq *TestFileQuery) Limit(limit int) *TestFileQuery {
	tfq.ctx.Limit = &limit
	return tfq
}

// Offset to start from.
func (tfq *TestFileQuery) Offset(offset int) *TestFileQuery {
	tfq.ctx.Offset = &offset
	return tfq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tfq *TestFileQuery) Unique(unique bool) *TestFileQuery {
	tfq.ctx.Unique = &unique
	return tfq
}

// Order specifies how the records should be ordered.
func (tfq *TestFileQuery) Order(o ...testfile.OrderOption) *TestFileQuery {
	tfq.order = append(tfq.order, o...)
	return tfq
}

// QueryTestResult chains the current query on the "test_result" edge.
func (tfq *TestFileQuery) QueryTestResult() *TestResultBESQuery {
	query := (&TestResultBESClient{config: tfq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := tfq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := tfq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(testfile.Table, testfile.FieldID, selector),
			sqlgraph.To(testresultbes.Table, testresultbes.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, testfile.TestResultTable, testfile.TestResultColumn),
		)
		fromU = sqlgraph.SetNeighbors(tfq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first TestFile entity from the query.
// Returns a *NotFoundError when no TestFile was found.
func (tfq *TestFileQuery) First(ctx context.Context) (*TestFile, error) {
	nodes, err := tfq.Limit(1).All(setContextOp(ctx, tfq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{testfile.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tfq *TestFileQuery) FirstX(ctx context.Context) *TestFile {
	node, err := tfq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TestFile ID from the query.
// Returns a *NotFoundError when no TestFile ID was found.
func (tfq *TestFileQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tfq.Limit(1).IDs(setContextOp(ctx, tfq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{testfile.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tfq *TestFileQuery) FirstIDX(ctx context.Context) int {
	id, err := tfq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TestFile entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one TestFile entity is found.
// Returns a *NotFoundError when no TestFile entities are found.
func (tfq *TestFileQuery) Only(ctx context.Context) (*TestFile, error) {
	nodes, err := tfq.Limit(2).All(setContextOp(ctx, tfq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{testfile.Label}
	default:
		return nil, &NotSingularError{testfile.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tfq *TestFileQuery) OnlyX(ctx context.Context) *TestFile {
	node, err := tfq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TestFile ID in the query.
// Returns a *NotSingularError when more than one TestFile ID is found.
// Returns a *NotFoundError when no entities are found.
func (tfq *TestFileQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tfq.Limit(2).IDs(setContextOp(ctx, tfq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{testfile.Label}
	default:
		err = &NotSingularError{testfile.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tfq *TestFileQuery) OnlyIDX(ctx context.Context) int {
	id, err := tfq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TestFiles.
func (tfq *TestFileQuery) All(ctx context.Context) ([]*TestFile, error) {
	ctx = setContextOp(ctx, tfq.ctx, ent.OpQueryAll)
	if err := tfq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*TestFile, *TestFileQuery]()
	return withInterceptors[[]*TestFile](ctx, tfq, qr, tfq.inters)
}

// AllX is like All, but panics if an error occurs.
func (tfq *TestFileQuery) AllX(ctx context.Context) []*TestFile {
	nodes, err := tfq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TestFile IDs.
func (tfq *TestFileQuery) IDs(ctx context.Context) (ids []int, err error) {
	if tfq.ctx.Unique == nil && tfq.path != nil {
		tfq.Unique(true)
	}
	ctx = setContextOp(ctx, tfq.ctx, ent.OpQueryIDs)
	if err = tfq.Select(testfile.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tfq *TestFileQuery) IDsX(ctx context.Context) []int {
	ids, err := tfq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tfq *TestFileQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, tfq.ctx, ent.OpQueryCount)
	if err := tfq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, tfq, querierCount[*TestFileQuery](), tfq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (tfq *TestFileQuery) CountX(ctx context.Context) int {
	count, err := tfq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tfq *TestFileQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, tfq.ctx, ent.OpQueryExist)
	switch _, err := tfq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (tfq *TestFileQuery) ExistX(ctx context.Context) bool {
	exist, err := tfq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TestFileQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tfq *TestFileQuery) Clone() *TestFileQuery {
	if tfq == nil {
		return nil
	}
	return &TestFileQuery{
		config:         tfq.config,
		ctx:            tfq.ctx.Clone(),
		order:          append([]testfile.OrderOption{}, tfq.order...),
		inters:         append([]Interceptor{}, tfq.inters...),
		predicates:     append([]predicate.TestFile{}, tfq.predicates...),
		withTestResult: tfq.withTestResult.Clone(),
		// clone intermediate query.
		sql:  tfq.sql.Clone(),
		path: tfq.path,
	}
}

// WithTestResult tells the query-builder to eager-load the nodes that are connected to
// the "test_result" edge. The optional arguments are used to configure the query builder of the edge.
func (tfq *TestFileQuery) WithTestResult(opts ...func(*TestResultBESQuery)) *TestFileQuery {
	query := (&TestResultBESClient{config: tfq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	tfq.withTestResult = query
	return tfq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Digest string `json:"digest,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.TestFile.Query().
//		GroupBy(testfile.FieldDigest).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (tfq *TestFileQuery) GroupBy(field string, fields ...string) *TestFileGroupBy {
	tfq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &TestFileGroupBy{build: tfq}
	grbuild.flds = &tfq.ctx.Fields
	grbuild.label = testfile.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Digest string `json:"digest,omitempty"`
//	}
//
//	client.TestFile.Query().
//		Select(testfile.FieldDigest).
//		Scan(ctx, &v)
func (tfq *TestFileQuery) Select(fields ...string) *TestFileSelect {
	tfq.ctx.Fields = append(tfq.ctx.Fields, fields...)
	sbuild := &TestFileSelect{TestFileQuery: tfq}
	sbuild.label = testfile.Label
	sbuild.flds, sbuild.scan = &tfq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a TestFileSelect configured with the given aggregations.
func (tfq *TestFileQuery) Aggregate(fns ...AggregateFunc) *TestFileSelect {
	return tfq.Select().Aggregate(fns...)
}

func (tfq *TestFileQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range tfq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, tfq); err != nil {
				return err
			}
		}
	}
	for _, f := range tfq.ctx.Fields {
		if !testfile.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tfq.path != nil {
		prev, err := tfq.path(ctx)
		if err != nil {
			return err
		}
		tfq.sql = prev
	}
	return nil
}

func (tfq *TestFileQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*TestFile, error) {
	var (
		nodes       = []*TestFile{}
		withFKs     = tfq.withFKs
		_spec       = tfq.querySpec()
		loadedTypes = [1]bool{
			tfq.withTestResult != nil,
		}
	)
	if tfq.withTestResult != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, testfile.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*TestFile).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &TestFile{config: tfq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(tfq.modifiers) > 0 {
		_spec.Modifiers = tfq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tfq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := tfq.withTestResult; query != nil {
		if err := tfq.loadTestResult(ctx, query, nodes, nil,
			func(n *TestFile, e *TestResultBES) { n.Edges.TestResult = e }); err != nil {
			return nil, err
		}
	}
	for i := range tfq.loadTotal {
		if err := tfq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (tfq *TestFileQuery) loadTestResult(ctx context.Context, query *TestResultBESQuery, nodes []*TestFile, init func(*TestFile), assign func(*TestFile, *TestResultBES)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*TestFile)
	for i := range nodes {
		if nodes[i].test_result_bes_test_action_output == nil {
			continue
		}
		fk := *nodes[i].test_result_bes_test_action_output
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(testresultbes.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "test_result_bes_test_action_output" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (tfq *TestFileQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tfq.querySpec()
	if len(tfq.modifiers) > 0 {
		_spec.Modifiers = tfq.modifiers
	}
	_spec.Node.Columns = tfq.ctx.Fields
	if len(tfq.ctx.Fields) > 0 {
		_spec.Unique = tfq.ctx.Unique != nil && *tfq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, tfq.driver, _spec)
}

func (tfq *TestFileQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(testfile.Table, testfile.Columns, sqlgraph.NewFieldSpec(testfile.FieldID, field.TypeInt))
	_spec.From = tfq.sql
	if unique := tfq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if tfq.path != nil {
		_spec.Unique = true
	}
	if fields := tfq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, testfile.FieldID)
		for i := range fields {
			if fields[i] != testfile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tfq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tfq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tfq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tfq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tfq *TestFileQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tfq.driver.Dialect())
	t1 := builder.Table(testfile.Table)
	columns := tfq.ctx.Fields
	if len(columns) == 0 {
		columns = testfile.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tfq.sql != nil {
		selector = tfq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tfq.ctx.Unique != nil && *tfq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range tfq.predicates {
		p(selector)
	}
	for _, p := range tfq.order {
		p(selector)
	}
	if offset := tfq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tfq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TestFileGroupBy is the group-by builder for TestFile entities.
type TestFileGroupBy struct {
	selector
	build *TestFileQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tfgb *TestFileGroupBy) Aggregate(fns ...AggregateFunc) *TestFileGroupBy {
	tfgb.fns = append(tfgb.fns, fns...)
	return tfgb
}

// Scan applies the selector query and scans the result into the given value.
func (tfgb *TestFileGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tfgb.build.ctx, ent.OpQueryGroupBy)
	if err := tfgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TestFileQuery, *TestFileGroupBy](ctx, tfgb.build, tfgb, tfgb.build.inters, v)
}

func (tfgb *TestFileGroupBy) sqlScan(ctx context.Context, root *TestFileQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(tfgb.fns))
	for _, fn := range tfgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*tfgb.flds)+len(tfgb.fns))
		for _, f := range *tfgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*tfgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tfgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// TestFileSelect is the builder for selecting fields of TestFile entities.
type TestFileSelect struct {
	*TestFileQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (tfs *TestFileSelect) Aggregate(fns ...AggregateFunc) *TestFileSelect {
	tfs.fns = append(tfs.fns, fns...)
	return tfs
}

// Scan applies the selector query and scans the result into the given value.
func (tfs *TestFileSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, tfs.ctx, ent.OpQuerySelect)
	if err := tfs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*TestFileQuery, *TestFileSelect](ctx, tfs.TestFileQuery, tfs, tfs.inters, v)
}

func (tfs *TestFileSelect) sqlScan(ctx context.Context, root *TestFileQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(tfs.fns))
	for _, fn := range tfs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*tfs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tfs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
