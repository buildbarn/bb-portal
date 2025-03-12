// Code generated by ent, DO NOT EDIT.

package runnercount

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the runnercount type in the database.
	Label = "runner_count"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldExecKind holds the string denoting the exec_kind field in the database.
	FieldExecKind = "exec_kind"
	// FieldActionsExecuted holds the string denoting the actions_executed field in the database.
	FieldActionsExecuted = "actions_executed"
	// FieldActionSummaryID holds the string denoting the action_summary_id field in the database.
	FieldActionSummaryID = "action_summary_id"
	// EdgeActionSummary holds the string denoting the action_summary edge name in mutations.
	EdgeActionSummary = "action_summary"
	// Table holds the table name of the runnercount in the database.
	Table = "runner_counts"
	// ActionSummaryTable is the table that holds the action_summary relation/edge.
	ActionSummaryTable = "runner_counts"
	// ActionSummaryInverseTable is the table name for the ActionSummary entity.
	// It exists in this package in order to avoid circular dependency with the "actionsummary" package.
	ActionSummaryInverseTable = "action_summaries"
	// ActionSummaryColumn is the table column denoting the action_summary relation/edge.
	ActionSummaryColumn = "action_summary_id"
)

// Columns holds all SQL columns for runnercount fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldExecKind,
	FieldActionsExecuted,
	FieldActionSummaryID,
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

// OrderOption defines the ordering options for the RunnerCount queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByExecKind orders the results by the exec_kind field.
func ByExecKind(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExecKind, opts...).ToFunc()
}

// ByActionsExecuted orders the results by the actions_executed field.
func ByActionsExecuted(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActionsExecuted, opts...).ToFunc()
}

// ByActionSummaryID orders the results by the action_summary_id field.
func ByActionSummaryID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActionSummaryID, opts...).ToFunc()
}

// ByActionSummaryField orders the results by action_summary field.
func ByActionSummaryField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newActionSummaryStep(), sql.OrderByField(field, opts...))
	}
}
func newActionSummaryStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ActionSummaryInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ActionSummaryTable, ActionSummaryColumn),
	)
}
