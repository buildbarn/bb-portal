// Code generated by ent, DO NOT EDIT.

package bazelinvocationproblem

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the bazelinvocationproblem type in the database.
	Label = "bazel_invocation_problem"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldProblemType holds the string denoting the problem_type field in the database.
	FieldProblemType = "problem_type"
	// FieldLabel holds the string denoting the label field in the database.
	FieldLabel = "label"
	// FieldBepEvents holds the string denoting the bep_events field in the database.
	FieldBepEvents = "bep_events"
	// EdgeBazelInvocation holds the string denoting the bazel_invocation edge name in mutations.
	EdgeBazelInvocation = "bazel_invocation"
	// Table holds the table name of the bazelinvocationproblem in the database.
	Table = "bazel_invocation_problems"
	// BazelInvocationTable is the table that holds the bazel_invocation relation/edge.
	BazelInvocationTable = "bazel_invocation_problems"
	// BazelInvocationInverseTable is the table name for the BazelInvocation entity.
	// It exists in this package in order to avoid circular dependency with the "bazelinvocation" package.
	BazelInvocationInverseTable = "bazel_invocations"
	// BazelInvocationColumn is the table column denoting the bazel_invocation relation/edge.
	BazelInvocationColumn = "bazel_invocation_problems"
)

// Columns holds all SQL columns for bazelinvocationproblem fields.
var Columns = []string{
	FieldID,
	FieldProblemType,
	FieldLabel,
	FieldBepEvents,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "bazel_invocation_problems"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"bazel_invocation_problems",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the BazelInvocationProblem queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByProblemType orders the results by the problem_type field.
func ByProblemType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProblemType, opts...).ToFunc()
}

// ByLabel orders the results by the label field.
func ByLabel(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLabel, opts...).ToFunc()
}

// ByBazelInvocationField orders the results by bazel_invocation field.
func ByBazelInvocationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBazelInvocationStep(), sql.OrderByField(field, opts...))
	}
}
func newBazelInvocationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BazelInvocationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, BazelInvocationTable, BazelInvocationColumn),
	)
}
