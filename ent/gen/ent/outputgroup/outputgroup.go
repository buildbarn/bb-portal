// Code generated by ent, DO NOT EDIT.

package outputgroup

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the outputgroup type in the database.
	Label = "output_group"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldIncomplete holds the string denoting the incomplete field in the database.
	FieldIncomplete = "incomplete"
	// EdgeTargetComplete holds the string denoting the target_complete edge name in mutations.
	EdgeTargetComplete = "target_complete"
	// EdgeInlineFiles holds the string denoting the inline_files edge name in mutations.
	EdgeInlineFiles = "inline_files"
	// EdgeFileSets holds the string denoting the file_sets edge name in mutations.
	EdgeFileSets = "file_sets"
	// Table holds the table name of the outputgroup in the database.
	Table = "output_groups"
	// TargetCompleteTable is the table that holds the target_complete relation/edge.
	TargetCompleteTable = "output_groups"
	// TargetCompleteInverseTable is the table name for the TargetComplete entity.
	// It exists in this package in order to avoid circular dependency with the "targetcomplete" package.
	TargetCompleteInverseTable = "target_completes"
	// TargetCompleteColumn is the table column denoting the target_complete relation/edge.
	TargetCompleteColumn = "target_complete_output_group"
	// InlineFilesTable is the table that holds the inline_files relation/edge.
	InlineFilesTable = "test_files"
	// InlineFilesInverseTable is the table name for the TestFile entity.
	// It exists in this package in order to avoid circular dependency with the "testfile" package.
	InlineFilesInverseTable = "test_files"
	// InlineFilesColumn is the table column denoting the inline_files relation/edge.
	InlineFilesColumn = "output_group_inline_files"
	// FileSetsTable is the table that holds the file_sets relation/edge.
	FileSetsTable = "named_set_of_files"
	// FileSetsInverseTable is the table name for the NamedSetOfFiles entity.
	// It exists in this package in order to avoid circular dependency with the "namedsetoffiles" package.
	FileSetsInverseTable = "named_set_of_files"
	// FileSetsColumn is the table column denoting the file_sets relation/edge.
	FileSetsColumn = "output_group_file_sets"
)

// Columns holds all SQL columns for outputgroup fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldIncomplete,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "output_groups"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"target_complete_output_group",
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

// OrderOption defines the ordering options for the OutputGroup queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByIncomplete orders the results by the incomplete field.
func ByIncomplete(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIncomplete, opts...).ToFunc()
}

// ByTargetCompleteField orders the results by target_complete field.
func ByTargetCompleteField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTargetCompleteStep(), sql.OrderByField(field, opts...))
	}
}

// ByInlineFilesCount orders the results by inline_files count.
func ByInlineFilesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newInlineFilesStep(), opts...)
	}
}

// ByInlineFiles orders the results by inline_files terms.
func ByInlineFiles(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInlineFilesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByFileSetsField orders the results by file_sets field.
func ByFileSetsField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFileSetsStep(), sql.OrderByField(field, opts...))
	}
}
func newTargetCompleteStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TargetCompleteInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, TargetCompleteTable, TargetCompleteColumn),
	)
}
func newInlineFilesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InlineFilesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, InlineFilesTable, InlineFilesColumn),
	)
}
func newFileSetsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FileSetsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, FileSetsTable, FileSetsColumn),
	)
}
