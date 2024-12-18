// Code generated by ent, DO NOT EDIT.

package timingchild

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the timingchild type in the database.
	Label = "timing_child"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldTime holds the string denoting the time field in the database.
	FieldTime = "time"
	// EdgeTimingBreakdown holds the string denoting the timing_breakdown edge name in mutations.
	EdgeTimingBreakdown = "timing_breakdown"
	// Table holds the table name of the timingchild in the database.
	Table = "timing_childs"
	// TimingBreakdownTable is the table that holds the timing_breakdown relation/edge.
	TimingBreakdownTable = "timing_childs"
	// TimingBreakdownInverseTable is the table name for the TimingBreakdown entity.
	// It exists in this package in order to avoid circular dependency with the "timingbreakdown" package.
	TimingBreakdownInverseTable = "timing_breakdowns"
	// TimingBreakdownColumn is the table column denoting the timing_breakdown relation/edge.
	TimingBreakdownColumn = "timing_breakdown_child"
)

// Columns holds all SQL columns for timingchild fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "timing_childs"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"timing_breakdown_child",
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

// OrderOption defines the ordering options for the TimingChild queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByTime orders the results by the time field.
func ByTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTime, opts...).ToFunc()
}

// ByTimingBreakdownField orders the results by timing_breakdown field.
func ByTimingBreakdownField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTimingBreakdownStep(), sql.OrderByField(field, opts...))
	}
}
func newTimingBreakdownStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TimingBreakdownInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, TimingBreakdownTable, TimingBreakdownColumn),
	)
}
