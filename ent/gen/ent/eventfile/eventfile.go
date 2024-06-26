// Code generated by ent, DO NOT EDIT.

package eventfile

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the eventfile type in the database.
	Label = "event_file"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldModTime holds the string denoting the mod_time field in the database.
	FieldModTime = "mod_time"
	// FieldProtocol holds the string denoting the protocol field in the database.
	FieldProtocol = "protocol"
	// FieldMimeType holds the string denoting the mime_type field in the database.
	FieldMimeType = "mime_type"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldReason holds the string denoting the reason field in the database.
	FieldReason = "reason"
	// EdgeBazelInvocation holds the string denoting the bazel_invocation edge name in mutations.
	EdgeBazelInvocation = "bazel_invocation"
	// Table holds the table name of the eventfile in the database.
	Table = "event_files"
	// BazelInvocationTable is the table that holds the bazel_invocation relation/edge.
	BazelInvocationTable = "bazel_invocations"
	// BazelInvocationInverseTable is the table name for the BazelInvocation entity.
	// It exists in this package in order to avoid circular dependency with the "bazelinvocation" package.
	BazelInvocationInverseTable = "bazel_invocations"
	// BazelInvocationColumn is the table column denoting the bazel_invocation relation/edge.
	BazelInvocationColumn = "event_file_bazel_invocation"
)

// Columns holds all SQL columns for eventfile fields.
var Columns = []string{
	FieldID,
	FieldURL,
	FieldModTime,
	FieldProtocol,
	FieldMimeType,
	FieldStatus,
	FieldReason,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus string
)

// OrderOption defines the ordering options for the EventFile queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByURL orders the results by the url field.
func ByURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURL, opts...).ToFunc()
}

// ByModTime orders the results by the mod_time field.
func ByModTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModTime, opts...).ToFunc()
}

// ByProtocol orders the results by the protocol field.
func ByProtocol(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProtocol, opts...).ToFunc()
}

// ByMimeType orders the results by the mime_type field.
func ByMimeType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMimeType, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByReason orders the results by the reason field.
func ByReason(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReason, opts...).ToFunc()
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
		sqlgraph.Edge(sqlgraph.O2O, false, BazelInvocationTable, BazelInvocationColumn),
	)
}
