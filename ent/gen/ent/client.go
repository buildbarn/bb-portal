// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/buildbarn/bb-portal/ent/gen/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocationproblem"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventfile"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// BazelInvocation is the client for interacting with the BazelInvocation builders.
	BazelInvocation *BazelInvocationClient
	// BazelInvocationProblem is the client for interacting with the BazelInvocationProblem builders.
	BazelInvocationProblem *BazelInvocationProblemClient
	// Blob is the client for interacting with the Blob builders.
	Blob *BlobClient
	// Build is the client for interacting with the Build builders.
	Build *BuildClient
	// EventFile is the client for interacting with the EventFile builders.
	EventFile *EventFileClient
	// additional fields for node api
	tables tables
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.BazelInvocation = NewBazelInvocationClient(c.config)
	c.BazelInvocationProblem = NewBazelInvocationProblemClient(c.config)
	c.Blob = NewBlobClient(c.config)
	c.Build = NewBuildClient(c.config)
	c.EventFile = NewEventFileClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:                    ctx,
		config:                 cfg,
		BazelInvocation:        NewBazelInvocationClient(cfg),
		BazelInvocationProblem: NewBazelInvocationProblemClient(cfg),
		Blob:                   NewBlobClient(cfg),
		Build:                  NewBuildClient(cfg),
		EventFile:              NewEventFileClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:                    ctx,
		config:                 cfg,
		BazelInvocation:        NewBazelInvocationClient(cfg),
		BazelInvocationProblem: NewBazelInvocationProblemClient(cfg),
		Blob:                   NewBlobClient(cfg),
		Build:                  NewBuildClient(cfg),
		EventFile:              NewEventFileClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		BazelInvocation.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.BazelInvocation.Use(hooks...)
	c.BazelInvocationProblem.Use(hooks...)
	c.Blob.Use(hooks...)
	c.Build.Use(hooks...)
	c.EventFile.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.BazelInvocation.Intercept(interceptors...)
	c.BazelInvocationProblem.Intercept(interceptors...)
	c.Blob.Intercept(interceptors...)
	c.Build.Intercept(interceptors...)
	c.EventFile.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *BazelInvocationMutation:
		return c.BazelInvocation.mutate(ctx, m)
	case *BazelInvocationProblemMutation:
		return c.BazelInvocationProblem.mutate(ctx, m)
	case *BlobMutation:
		return c.Blob.mutate(ctx, m)
	case *BuildMutation:
		return c.Build.mutate(ctx, m)
	case *EventFileMutation:
		return c.EventFile.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// BazelInvocationClient is a client for the BazelInvocation schema.
type BazelInvocationClient struct {
	config
}

// NewBazelInvocationClient returns a client for the BazelInvocation from the given config.
func NewBazelInvocationClient(c config) *BazelInvocationClient {
	return &BazelInvocationClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `bazelinvocation.Hooks(f(g(h())))`.
func (c *BazelInvocationClient) Use(hooks ...Hook) {
	c.hooks.BazelInvocation = append(c.hooks.BazelInvocation, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `bazelinvocation.Intercept(f(g(h())))`.
func (c *BazelInvocationClient) Intercept(interceptors ...Interceptor) {
	c.inters.BazelInvocation = append(c.inters.BazelInvocation, interceptors...)
}

// Create returns a builder for creating a BazelInvocation entity.
func (c *BazelInvocationClient) Create() *BazelInvocationCreate {
	mutation := newBazelInvocationMutation(c.config, OpCreate)
	return &BazelInvocationCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of BazelInvocation entities.
func (c *BazelInvocationClient) CreateBulk(builders ...*BazelInvocationCreate) *BazelInvocationCreateBulk {
	return &BazelInvocationCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BazelInvocationClient) MapCreateBulk(slice any, setFunc func(*BazelInvocationCreate, int)) *BazelInvocationCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BazelInvocationCreateBulk{err: fmt.Errorf("calling to BazelInvocationClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BazelInvocationCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BazelInvocationCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for BazelInvocation.
func (c *BazelInvocationClient) Update() *BazelInvocationUpdate {
	mutation := newBazelInvocationMutation(c.config, OpUpdate)
	return &BazelInvocationUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BazelInvocationClient) UpdateOne(bi *BazelInvocation) *BazelInvocationUpdateOne {
	mutation := newBazelInvocationMutation(c.config, OpUpdateOne, withBazelInvocation(bi))
	return &BazelInvocationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BazelInvocationClient) UpdateOneID(id int) *BazelInvocationUpdateOne {
	mutation := newBazelInvocationMutation(c.config, OpUpdateOne, withBazelInvocationID(id))
	return &BazelInvocationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for BazelInvocation.
func (c *BazelInvocationClient) Delete() *BazelInvocationDelete {
	mutation := newBazelInvocationMutation(c.config, OpDelete)
	return &BazelInvocationDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BazelInvocationClient) DeleteOne(bi *BazelInvocation) *BazelInvocationDeleteOne {
	return c.DeleteOneID(bi.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BazelInvocationClient) DeleteOneID(id int) *BazelInvocationDeleteOne {
	builder := c.Delete().Where(bazelinvocation.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BazelInvocationDeleteOne{builder}
}

// Query returns a query builder for BazelInvocation.
func (c *BazelInvocationClient) Query() *BazelInvocationQuery {
	return &BazelInvocationQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBazelInvocation},
		inters: c.Interceptors(),
	}
}

// Get returns a BazelInvocation entity by its id.
func (c *BazelInvocationClient) Get(ctx context.Context, id int) (*BazelInvocation, error) {
	return c.Query().Where(bazelinvocation.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BazelInvocationClient) GetX(ctx context.Context, id int) *BazelInvocation {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryEventFile queries the event_file edge of a BazelInvocation.
func (c *BazelInvocationClient) QueryEventFile(bi *BazelInvocation) *EventFileQuery {
	query := (&EventFileClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := bi.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(bazelinvocation.Table, bazelinvocation.FieldID, id),
			sqlgraph.To(eventfile.Table, eventfile.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, bazelinvocation.EventFileTable, bazelinvocation.EventFileColumn),
		)
		fromV = sqlgraph.Neighbors(bi.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryBuild queries the build edge of a BazelInvocation.
func (c *BazelInvocationClient) QueryBuild(bi *BazelInvocation) *BuildQuery {
	query := (&BuildClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := bi.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(bazelinvocation.Table, bazelinvocation.FieldID, id),
			sqlgraph.To(build.Table, build.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, bazelinvocation.BuildTable, bazelinvocation.BuildColumn),
		)
		fromV = sqlgraph.Neighbors(bi.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryProblems queries the problems edge of a BazelInvocation.
func (c *BazelInvocationClient) QueryProblems(bi *BazelInvocation) *BazelInvocationProblemQuery {
	query := (&BazelInvocationProblemClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := bi.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(bazelinvocation.Table, bazelinvocation.FieldID, id),
			sqlgraph.To(bazelinvocationproblem.Table, bazelinvocationproblem.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, bazelinvocation.ProblemsTable, bazelinvocation.ProblemsColumn),
		)
		fromV = sqlgraph.Neighbors(bi.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *BazelInvocationClient) Hooks() []Hook {
	return c.hooks.BazelInvocation
}

// Interceptors returns the client interceptors.
func (c *BazelInvocationClient) Interceptors() []Interceptor {
	return c.inters.BazelInvocation
}

func (c *BazelInvocationClient) mutate(ctx context.Context, m *BazelInvocationMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BazelInvocationCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BazelInvocationUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BazelInvocationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BazelInvocationDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown BazelInvocation mutation op: %q", m.Op())
	}
}

// BazelInvocationProblemClient is a client for the BazelInvocationProblem schema.
type BazelInvocationProblemClient struct {
	config
}

// NewBazelInvocationProblemClient returns a client for the BazelInvocationProblem from the given config.
func NewBazelInvocationProblemClient(c config) *BazelInvocationProblemClient {
	return &BazelInvocationProblemClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `bazelinvocationproblem.Hooks(f(g(h())))`.
func (c *BazelInvocationProblemClient) Use(hooks ...Hook) {
	c.hooks.BazelInvocationProblem = append(c.hooks.BazelInvocationProblem, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `bazelinvocationproblem.Intercept(f(g(h())))`.
func (c *BazelInvocationProblemClient) Intercept(interceptors ...Interceptor) {
	c.inters.BazelInvocationProblem = append(c.inters.BazelInvocationProblem, interceptors...)
}

// Create returns a builder for creating a BazelInvocationProblem entity.
func (c *BazelInvocationProblemClient) Create() *BazelInvocationProblemCreate {
	mutation := newBazelInvocationProblemMutation(c.config, OpCreate)
	return &BazelInvocationProblemCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of BazelInvocationProblem entities.
func (c *BazelInvocationProblemClient) CreateBulk(builders ...*BazelInvocationProblemCreate) *BazelInvocationProblemCreateBulk {
	return &BazelInvocationProblemCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BazelInvocationProblemClient) MapCreateBulk(slice any, setFunc func(*BazelInvocationProblemCreate, int)) *BazelInvocationProblemCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BazelInvocationProblemCreateBulk{err: fmt.Errorf("calling to BazelInvocationProblemClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BazelInvocationProblemCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BazelInvocationProblemCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for BazelInvocationProblem.
func (c *BazelInvocationProblemClient) Update() *BazelInvocationProblemUpdate {
	mutation := newBazelInvocationProblemMutation(c.config, OpUpdate)
	return &BazelInvocationProblemUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BazelInvocationProblemClient) UpdateOne(bip *BazelInvocationProblem) *BazelInvocationProblemUpdateOne {
	mutation := newBazelInvocationProblemMutation(c.config, OpUpdateOne, withBazelInvocationProblem(bip))
	return &BazelInvocationProblemUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BazelInvocationProblemClient) UpdateOneID(id int) *BazelInvocationProblemUpdateOne {
	mutation := newBazelInvocationProblemMutation(c.config, OpUpdateOne, withBazelInvocationProblemID(id))
	return &BazelInvocationProblemUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for BazelInvocationProblem.
func (c *BazelInvocationProblemClient) Delete() *BazelInvocationProblemDelete {
	mutation := newBazelInvocationProblemMutation(c.config, OpDelete)
	return &BazelInvocationProblemDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BazelInvocationProblemClient) DeleteOne(bip *BazelInvocationProblem) *BazelInvocationProblemDeleteOne {
	return c.DeleteOneID(bip.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BazelInvocationProblemClient) DeleteOneID(id int) *BazelInvocationProblemDeleteOne {
	builder := c.Delete().Where(bazelinvocationproblem.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BazelInvocationProblemDeleteOne{builder}
}

// Query returns a query builder for BazelInvocationProblem.
func (c *BazelInvocationProblemClient) Query() *BazelInvocationProblemQuery {
	return &BazelInvocationProblemQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBazelInvocationProblem},
		inters: c.Interceptors(),
	}
}

// Get returns a BazelInvocationProblem entity by its id.
func (c *BazelInvocationProblemClient) Get(ctx context.Context, id int) (*BazelInvocationProblem, error) {
	return c.Query().Where(bazelinvocationproblem.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BazelInvocationProblemClient) GetX(ctx context.Context, id int) *BazelInvocationProblem {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryBazelInvocation queries the bazel_invocation edge of a BazelInvocationProblem.
func (c *BazelInvocationProblemClient) QueryBazelInvocation(bip *BazelInvocationProblem) *BazelInvocationQuery {
	query := (&BazelInvocationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := bip.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(bazelinvocationproblem.Table, bazelinvocationproblem.FieldID, id),
			sqlgraph.To(bazelinvocation.Table, bazelinvocation.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, bazelinvocationproblem.BazelInvocationTable, bazelinvocationproblem.BazelInvocationColumn),
		)
		fromV = sqlgraph.Neighbors(bip.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *BazelInvocationProblemClient) Hooks() []Hook {
	return c.hooks.BazelInvocationProblem
}

// Interceptors returns the client interceptors.
func (c *BazelInvocationProblemClient) Interceptors() []Interceptor {
	return c.inters.BazelInvocationProblem
}

func (c *BazelInvocationProblemClient) mutate(ctx context.Context, m *BazelInvocationProblemMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BazelInvocationProblemCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BazelInvocationProblemUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BazelInvocationProblemUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BazelInvocationProblemDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown BazelInvocationProblem mutation op: %q", m.Op())
	}
}

// BlobClient is a client for the Blob schema.
type BlobClient struct {
	config
}

// NewBlobClient returns a client for the Blob from the given config.
func NewBlobClient(c config) *BlobClient {
	return &BlobClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `blob.Hooks(f(g(h())))`.
func (c *BlobClient) Use(hooks ...Hook) {
	c.hooks.Blob = append(c.hooks.Blob, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `blob.Intercept(f(g(h())))`.
func (c *BlobClient) Intercept(interceptors ...Interceptor) {
	c.inters.Blob = append(c.inters.Blob, interceptors...)
}

// Create returns a builder for creating a Blob entity.
func (c *BlobClient) Create() *BlobCreate {
	mutation := newBlobMutation(c.config, OpCreate)
	return &BlobCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Blob entities.
func (c *BlobClient) CreateBulk(builders ...*BlobCreate) *BlobCreateBulk {
	return &BlobCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BlobClient) MapCreateBulk(slice any, setFunc func(*BlobCreate, int)) *BlobCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BlobCreateBulk{err: fmt.Errorf("calling to BlobClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BlobCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BlobCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Blob.
func (c *BlobClient) Update() *BlobUpdate {
	mutation := newBlobMutation(c.config, OpUpdate)
	return &BlobUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BlobClient) UpdateOne(b *Blob) *BlobUpdateOne {
	mutation := newBlobMutation(c.config, OpUpdateOne, withBlob(b))
	return &BlobUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BlobClient) UpdateOneID(id int) *BlobUpdateOne {
	mutation := newBlobMutation(c.config, OpUpdateOne, withBlobID(id))
	return &BlobUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Blob.
func (c *BlobClient) Delete() *BlobDelete {
	mutation := newBlobMutation(c.config, OpDelete)
	return &BlobDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BlobClient) DeleteOne(b *Blob) *BlobDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BlobClient) DeleteOneID(id int) *BlobDeleteOne {
	builder := c.Delete().Where(blob.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BlobDeleteOne{builder}
}

// Query returns a query builder for Blob.
func (c *BlobClient) Query() *BlobQuery {
	return &BlobQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBlob},
		inters: c.Interceptors(),
	}
}

// Get returns a Blob entity by its id.
func (c *BlobClient) Get(ctx context.Context, id int) (*Blob, error) {
	return c.Query().Where(blob.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BlobClient) GetX(ctx context.Context, id int) *Blob {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *BlobClient) Hooks() []Hook {
	return c.hooks.Blob
}

// Interceptors returns the client interceptors.
func (c *BlobClient) Interceptors() []Interceptor {
	return c.inters.Blob
}

func (c *BlobClient) mutate(ctx context.Context, m *BlobMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BlobCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BlobUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BlobUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BlobDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Blob mutation op: %q", m.Op())
	}
}

// BuildClient is a client for the Build schema.
type BuildClient struct {
	config
}

// NewBuildClient returns a client for the Build from the given config.
func NewBuildClient(c config) *BuildClient {
	return &BuildClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `build.Hooks(f(g(h())))`.
func (c *BuildClient) Use(hooks ...Hook) {
	c.hooks.Build = append(c.hooks.Build, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `build.Intercept(f(g(h())))`.
func (c *BuildClient) Intercept(interceptors ...Interceptor) {
	c.inters.Build = append(c.inters.Build, interceptors...)
}

// Create returns a builder for creating a Build entity.
func (c *BuildClient) Create() *BuildCreate {
	mutation := newBuildMutation(c.config, OpCreate)
	return &BuildCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Build entities.
func (c *BuildClient) CreateBulk(builders ...*BuildCreate) *BuildCreateBulk {
	return &BuildCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BuildClient) MapCreateBulk(slice any, setFunc func(*BuildCreate, int)) *BuildCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BuildCreateBulk{err: fmt.Errorf("calling to BuildClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BuildCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BuildCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Build.
func (c *BuildClient) Update() *BuildUpdate {
	mutation := newBuildMutation(c.config, OpUpdate)
	return &BuildUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BuildClient) UpdateOne(b *Build) *BuildUpdateOne {
	mutation := newBuildMutation(c.config, OpUpdateOne, withBuild(b))
	return &BuildUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BuildClient) UpdateOneID(id int) *BuildUpdateOne {
	mutation := newBuildMutation(c.config, OpUpdateOne, withBuildID(id))
	return &BuildUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Build.
func (c *BuildClient) Delete() *BuildDelete {
	mutation := newBuildMutation(c.config, OpDelete)
	return &BuildDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BuildClient) DeleteOne(b *Build) *BuildDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BuildClient) DeleteOneID(id int) *BuildDeleteOne {
	builder := c.Delete().Where(build.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BuildDeleteOne{builder}
}

// Query returns a query builder for Build.
func (c *BuildClient) Query() *BuildQuery {
	return &BuildQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBuild},
		inters: c.Interceptors(),
	}
}

// Get returns a Build entity by its id.
func (c *BuildClient) Get(ctx context.Context, id int) (*Build, error) {
	return c.Query().Where(build.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BuildClient) GetX(ctx context.Context, id int) *Build {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryInvocations queries the invocations edge of a Build.
func (c *BuildClient) QueryInvocations(b *Build) *BazelInvocationQuery {
	query := (&BazelInvocationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := b.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(build.Table, build.FieldID, id),
			sqlgraph.To(bazelinvocation.Table, bazelinvocation.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, build.InvocationsTable, build.InvocationsColumn),
		)
		fromV = sqlgraph.Neighbors(b.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *BuildClient) Hooks() []Hook {
	return c.hooks.Build
}

// Interceptors returns the client interceptors.
func (c *BuildClient) Interceptors() []Interceptor {
	return c.inters.Build
}

func (c *BuildClient) mutate(ctx context.Context, m *BuildMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BuildCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BuildUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BuildUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BuildDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Build mutation op: %q", m.Op())
	}
}

// EventFileClient is a client for the EventFile schema.
type EventFileClient struct {
	config
}

// NewEventFileClient returns a client for the EventFile from the given config.
func NewEventFileClient(c config) *EventFileClient {
	return &EventFileClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `eventfile.Hooks(f(g(h())))`.
func (c *EventFileClient) Use(hooks ...Hook) {
	c.hooks.EventFile = append(c.hooks.EventFile, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `eventfile.Intercept(f(g(h())))`.
func (c *EventFileClient) Intercept(interceptors ...Interceptor) {
	c.inters.EventFile = append(c.inters.EventFile, interceptors...)
}

// Create returns a builder for creating a EventFile entity.
func (c *EventFileClient) Create() *EventFileCreate {
	mutation := newEventFileMutation(c.config, OpCreate)
	return &EventFileCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of EventFile entities.
func (c *EventFileClient) CreateBulk(builders ...*EventFileCreate) *EventFileCreateBulk {
	return &EventFileCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *EventFileClient) MapCreateBulk(slice any, setFunc func(*EventFileCreate, int)) *EventFileCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &EventFileCreateBulk{err: fmt.Errorf("calling to EventFileClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*EventFileCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &EventFileCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for EventFile.
func (c *EventFileClient) Update() *EventFileUpdate {
	mutation := newEventFileMutation(c.config, OpUpdate)
	return &EventFileUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EventFileClient) UpdateOne(ef *EventFile) *EventFileUpdateOne {
	mutation := newEventFileMutation(c.config, OpUpdateOne, withEventFile(ef))
	return &EventFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *EventFileClient) UpdateOneID(id int) *EventFileUpdateOne {
	mutation := newEventFileMutation(c.config, OpUpdateOne, withEventFileID(id))
	return &EventFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EventFile.
func (c *EventFileClient) Delete() *EventFileDelete {
	mutation := newEventFileMutation(c.config, OpDelete)
	return &EventFileDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *EventFileClient) DeleteOne(ef *EventFile) *EventFileDeleteOne {
	return c.DeleteOneID(ef.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *EventFileClient) DeleteOneID(id int) *EventFileDeleteOne {
	builder := c.Delete().Where(eventfile.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EventFileDeleteOne{builder}
}

// Query returns a query builder for EventFile.
func (c *EventFileClient) Query() *EventFileQuery {
	return &EventFileQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeEventFile},
		inters: c.Interceptors(),
	}
}

// Get returns a EventFile entity by its id.
func (c *EventFileClient) Get(ctx context.Context, id int) (*EventFile, error) {
	return c.Query().Where(eventfile.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EventFileClient) GetX(ctx context.Context, id int) *EventFile {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryBazelInvocation queries the bazel_invocation edge of a EventFile.
func (c *EventFileClient) QueryBazelInvocation(ef *EventFile) *BazelInvocationQuery {
	query := (&BazelInvocationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ef.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(eventfile.Table, eventfile.FieldID, id),
			sqlgraph.To(bazelinvocation.Table, bazelinvocation.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, eventfile.BazelInvocationTable, eventfile.BazelInvocationColumn),
		)
		fromV = sqlgraph.Neighbors(ef.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EventFileClient) Hooks() []Hook {
	return c.hooks.EventFile
}

// Interceptors returns the client interceptors.
func (c *EventFileClient) Interceptors() []Interceptor {
	return c.inters.EventFile
}

func (c *EventFileClient) mutate(ctx context.Context, m *EventFileMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&EventFileCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&EventFileUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&EventFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&EventFileDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown EventFile mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		BazelInvocation, BazelInvocationProblem, Blob, Build, EventFile []ent.Hook
	}
	inters struct {
		BazelInvocation, BazelInvocationProblem, Blob, Build,
		EventFile []ent.Interceptor
	}
)
