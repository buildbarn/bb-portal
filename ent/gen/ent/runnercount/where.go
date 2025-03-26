// Code generated by ent, DO NOT EDIT.

package runnercount

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldName, v))
}

// ExecKind applies equality check predicate on the "exec_kind" field. It's identical to ExecKindEQ.
func ExecKind(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldExecKind, v))
}

// ActionsExecuted applies equality check predicate on the "actions_executed" field. It's identical to ActionsExecutedEQ.
func ActionsExecuted(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldActionsExecuted, v))
}

// ActionSummaryID applies equality check predicate on the "action_summary_id" field. It's identical to ActionSummaryIDEQ.
func ActionSummaryID(v int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldActionSummaryID, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldHasSuffix(FieldName, v))
}

// NameIsNil applies the IsNil predicate on the "name" field.
func NameIsNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIsNull(FieldName))
}

// NameNotNil applies the NotNil predicate on the "name" field.
func NameNotNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotNull(FieldName))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldContainsFold(FieldName, v))
}

// ExecKindEQ applies the EQ predicate on the "exec_kind" field.
func ExecKindEQ(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldExecKind, v))
}

// ExecKindNEQ applies the NEQ predicate on the "exec_kind" field.
func ExecKindNEQ(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNEQ(FieldExecKind, v))
}

// ExecKindIn applies the In predicate on the "exec_kind" field.
func ExecKindIn(vs ...string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIn(FieldExecKind, vs...))
}

// ExecKindNotIn applies the NotIn predicate on the "exec_kind" field.
func ExecKindNotIn(vs ...string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotIn(FieldExecKind, vs...))
}

// ExecKindGT applies the GT predicate on the "exec_kind" field.
func ExecKindGT(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGT(FieldExecKind, v))
}

// ExecKindGTE applies the GTE predicate on the "exec_kind" field.
func ExecKindGTE(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGTE(FieldExecKind, v))
}

// ExecKindLT applies the LT predicate on the "exec_kind" field.
func ExecKindLT(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLT(FieldExecKind, v))
}

// ExecKindLTE applies the LTE predicate on the "exec_kind" field.
func ExecKindLTE(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLTE(FieldExecKind, v))
}

// ExecKindContains applies the Contains predicate on the "exec_kind" field.
func ExecKindContains(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldContains(FieldExecKind, v))
}

// ExecKindHasPrefix applies the HasPrefix predicate on the "exec_kind" field.
func ExecKindHasPrefix(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldHasPrefix(FieldExecKind, v))
}

// ExecKindHasSuffix applies the HasSuffix predicate on the "exec_kind" field.
func ExecKindHasSuffix(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldHasSuffix(FieldExecKind, v))
}

// ExecKindIsNil applies the IsNil predicate on the "exec_kind" field.
func ExecKindIsNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIsNull(FieldExecKind))
}

// ExecKindNotNil applies the NotNil predicate on the "exec_kind" field.
func ExecKindNotNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotNull(FieldExecKind))
}

// ExecKindEqualFold applies the EqualFold predicate on the "exec_kind" field.
func ExecKindEqualFold(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEqualFold(FieldExecKind, v))
}

// ExecKindContainsFold applies the ContainsFold predicate on the "exec_kind" field.
func ExecKindContainsFold(v string) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldContainsFold(FieldExecKind, v))
}

// ActionsExecutedEQ applies the EQ predicate on the "actions_executed" field.
func ActionsExecutedEQ(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldActionsExecuted, v))
}

// ActionsExecutedNEQ applies the NEQ predicate on the "actions_executed" field.
func ActionsExecutedNEQ(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNEQ(FieldActionsExecuted, v))
}

// ActionsExecutedIn applies the In predicate on the "actions_executed" field.
func ActionsExecutedIn(vs ...int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIn(FieldActionsExecuted, vs...))
}

// ActionsExecutedNotIn applies the NotIn predicate on the "actions_executed" field.
func ActionsExecutedNotIn(vs ...int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotIn(FieldActionsExecuted, vs...))
}

// ActionsExecutedGT applies the GT predicate on the "actions_executed" field.
func ActionsExecutedGT(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGT(FieldActionsExecuted, v))
}

// ActionsExecutedGTE applies the GTE predicate on the "actions_executed" field.
func ActionsExecutedGTE(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldGTE(FieldActionsExecuted, v))
}

// ActionsExecutedLT applies the LT predicate on the "actions_executed" field.
func ActionsExecutedLT(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLT(FieldActionsExecuted, v))
}

// ActionsExecutedLTE applies the LTE predicate on the "actions_executed" field.
func ActionsExecutedLTE(v int64) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldLTE(FieldActionsExecuted, v))
}

// ActionsExecutedIsNil applies the IsNil predicate on the "actions_executed" field.
func ActionsExecutedIsNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIsNull(FieldActionsExecuted))
}

// ActionsExecutedNotNil applies the NotNil predicate on the "actions_executed" field.
func ActionsExecutedNotNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotNull(FieldActionsExecuted))
}

// ActionSummaryIDEQ applies the EQ predicate on the "action_summary_id" field.
func ActionSummaryIDEQ(v int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldEQ(FieldActionSummaryID, v))
}

// ActionSummaryIDNEQ applies the NEQ predicate on the "action_summary_id" field.
func ActionSummaryIDNEQ(v int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNEQ(FieldActionSummaryID, v))
}

// ActionSummaryIDIn applies the In predicate on the "action_summary_id" field.
func ActionSummaryIDIn(vs ...int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIn(FieldActionSummaryID, vs...))
}

// ActionSummaryIDNotIn applies the NotIn predicate on the "action_summary_id" field.
func ActionSummaryIDNotIn(vs ...int) predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotIn(FieldActionSummaryID, vs...))
}

// ActionSummaryIDIsNil applies the IsNil predicate on the "action_summary_id" field.
func ActionSummaryIDIsNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldIsNull(FieldActionSummaryID))
}

// ActionSummaryIDNotNil applies the NotNil predicate on the "action_summary_id" field.
func ActionSummaryIDNotNil() predicate.RunnerCount {
	return predicate.RunnerCount(sql.FieldNotNull(FieldActionSummaryID))
}

// HasActionSummary applies the HasEdge predicate on the "action_summary" edge.
func HasActionSummary() predicate.RunnerCount {
	return predicate.RunnerCount(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ActionSummaryTable, ActionSummaryColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasActionSummaryWith applies the HasEdge predicate on the "action_summary" edge with a given conditions (other predicates).
func HasActionSummaryWith(preds ...predicate.ActionSummary) predicate.RunnerCount {
	return predicate.RunnerCount(func(s *sql.Selector) {
		step := newActionSummaryStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.RunnerCount) predicate.RunnerCount {
	return predicate.RunnerCount(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.RunnerCount) predicate.RunnerCount {
	return predicate.RunnerCount(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.RunnerCount) predicate.RunnerCount {
	return predicate.RunnerCount(sql.NotPredicates(p))
}
