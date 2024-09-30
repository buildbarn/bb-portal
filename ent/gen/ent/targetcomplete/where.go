// Code generated by ent, DO NOT EDIT.

package targetcomplete

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLTE(FieldID, id))
}

// Success applies equality check predicate on the "success" field. It's identical to SuccessEQ.
func Success(v bool) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldSuccess, v))
}

// TargetKind applies equality check predicate on the "target_kind" field. It's identical to TargetKindEQ.
func TargetKind(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTargetKind, v))
}

// EndTimeInMs applies equality check predicate on the "end_time_in_ms" field. It's identical to EndTimeInMsEQ.
func EndTimeInMs(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldEndTimeInMs, v))
}

// TestTimeoutSeconds applies equality check predicate on the "test_timeout_seconds" field. It's identical to TestTimeoutSecondsEQ.
func TestTimeoutSeconds(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTestTimeoutSeconds, v))
}

// TestTimeout applies equality check predicate on the "test_timeout" field. It's identical to TestTimeoutEQ.
func TestTimeout(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTestTimeout, v))
}

// SuccessEQ applies the EQ predicate on the "success" field.
func SuccessEQ(v bool) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldSuccess, v))
}

// SuccessNEQ applies the NEQ predicate on the "success" field.
func SuccessNEQ(v bool) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldSuccess, v))
}

// SuccessIsNil applies the IsNil predicate on the "success" field.
func SuccessIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldSuccess))
}

// SuccessNotNil applies the NotNil predicate on the "success" field.
func SuccessNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldSuccess))
}

// TagIsNil applies the IsNil predicate on the "tag" field.
func TagIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldTag))
}

// TagNotNil applies the NotNil predicate on the "tag" field.
func TagNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldTag))
}

// TargetKindEQ applies the EQ predicate on the "target_kind" field.
func TargetKindEQ(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTargetKind, v))
}

// TargetKindNEQ applies the NEQ predicate on the "target_kind" field.
func TargetKindNEQ(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldTargetKind, v))
}

// TargetKindIn applies the In predicate on the "target_kind" field.
func TargetKindIn(vs ...string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIn(FieldTargetKind, vs...))
}

// TargetKindNotIn applies the NotIn predicate on the "target_kind" field.
func TargetKindNotIn(vs ...string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotIn(FieldTargetKind, vs...))
}

// TargetKindGT applies the GT predicate on the "target_kind" field.
func TargetKindGT(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGT(FieldTargetKind, v))
}

// TargetKindGTE applies the GTE predicate on the "target_kind" field.
func TargetKindGTE(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGTE(FieldTargetKind, v))
}

// TargetKindLT applies the LT predicate on the "target_kind" field.
func TargetKindLT(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLT(FieldTargetKind, v))
}

// TargetKindLTE applies the LTE predicate on the "target_kind" field.
func TargetKindLTE(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLTE(FieldTargetKind, v))
}

// TargetKindContains applies the Contains predicate on the "target_kind" field.
func TargetKindContains(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldContains(FieldTargetKind, v))
}

// TargetKindHasPrefix applies the HasPrefix predicate on the "target_kind" field.
func TargetKindHasPrefix(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldHasPrefix(FieldTargetKind, v))
}

// TargetKindHasSuffix applies the HasSuffix predicate on the "target_kind" field.
func TargetKindHasSuffix(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldHasSuffix(FieldTargetKind, v))
}

// TargetKindIsNil applies the IsNil predicate on the "target_kind" field.
func TargetKindIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldTargetKind))
}

// TargetKindNotNil applies the NotNil predicate on the "target_kind" field.
func TargetKindNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldTargetKind))
}

// TargetKindEqualFold applies the EqualFold predicate on the "target_kind" field.
func TargetKindEqualFold(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEqualFold(FieldTargetKind, v))
}

// TargetKindContainsFold applies the ContainsFold predicate on the "target_kind" field.
func TargetKindContainsFold(v string) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldContainsFold(FieldTargetKind, v))
}

// EndTimeInMsEQ applies the EQ predicate on the "end_time_in_ms" field.
func EndTimeInMsEQ(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldEndTimeInMs, v))
}

// EndTimeInMsNEQ applies the NEQ predicate on the "end_time_in_ms" field.
func EndTimeInMsNEQ(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldEndTimeInMs, v))
}

// EndTimeInMsIn applies the In predicate on the "end_time_in_ms" field.
func EndTimeInMsIn(vs ...int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIn(FieldEndTimeInMs, vs...))
}

// EndTimeInMsNotIn applies the NotIn predicate on the "end_time_in_ms" field.
func EndTimeInMsNotIn(vs ...int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotIn(FieldEndTimeInMs, vs...))
}

// EndTimeInMsGT applies the GT predicate on the "end_time_in_ms" field.
func EndTimeInMsGT(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGT(FieldEndTimeInMs, v))
}

// EndTimeInMsGTE applies the GTE predicate on the "end_time_in_ms" field.
func EndTimeInMsGTE(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGTE(FieldEndTimeInMs, v))
}

// EndTimeInMsLT applies the LT predicate on the "end_time_in_ms" field.
func EndTimeInMsLT(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLT(FieldEndTimeInMs, v))
}

// EndTimeInMsLTE applies the LTE predicate on the "end_time_in_ms" field.
func EndTimeInMsLTE(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLTE(FieldEndTimeInMs, v))
}

// EndTimeInMsIsNil applies the IsNil predicate on the "end_time_in_ms" field.
func EndTimeInMsIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldEndTimeInMs))
}

// EndTimeInMsNotNil applies the NotNil predicate on the "end_time_in_ms" field.
func EndTimeInMsNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldEndTimeInMs))
}

// TestTimeoutSecondsEQ applies the EQ predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsEQ(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTestTimeoutSeconds, v))
}

// TestTimeoutSecondsNEQ applies the NEQ predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsNEQ(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldTestTimeoutSeconds, v))
}

// TestTimeoutSecondsIn applies the In predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsIn(vs ...int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIn(FieldTestTimeoutSeconds, vs...))
}

// TestTimeoutSecondsNotIn applies the NotIn predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsNotIn(vs ...int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotIn(FieldTestTimeoutSeconds, vs...))
}

// TestTimeoutSecondsGT applies the GT predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsGT(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGT(FieldTestTimeoutSeconds, v))
}

// TestTimeoutSecondsGTE applies the GTE predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsGTE(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGTE(FieldTestTimeoutSeconds, v))
}

// TestTimeoutSecondsLT applies the LT predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsLT(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLT(FieldTestTimeoutSeconds, v))
}

// TestTimeoutSecondsLTE applies the LTE predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsLTE(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLTE(FieldTestTimeoutSeconds, v))
}

// TestTimeoutSecondsIsNil applies the IsNil predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldTestTimeoutSeconds))
}

// TestTimeoutSecondsNotNil applies the NotNil predicate on the "test_timeout_seconds" field.
func TestTimeoutSecondsNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldTestTimeoutSeconds))
}

// TestTimeoutEQ applies the EQ predicate on the "test_timeout" field.
func TestTimeoutEQ(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTestTimeout, v))
}

// TestTimeoutNEQ applies the NEQ predicate on the "test_timeout" field.
func TestTimeoutNEQ(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldTestTimeout, v))
}

// TestTimeoutIn applies the In predicate on the "test_timeout" field.
func TestTimeoutIn(vs ...int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIn(FieldTestTimeout, vs...))
}

// TestTimeoutNotIn applies the NotIn predicate on the "test_timeout" field.
func TestTimeoutNotIn(vs ...int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotIn(FieldTestTimeout, vs...))
}

// TestTimeoutGT applies the GT predicate on the "test_timeout" field.
func TestTimeoutGT(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGT(FieldTestTimeout, v))
}

// TestTimeoutGTE applies the GTE predicate on the "test_timeout" field.
func TestTimeoutGTE(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldGTE(FieldTestTimeout, v))
}

// TestTimeoutLT applies the LT predicate on the "test_timeout" field.
func TestTimeoutLT(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLT(FieldTestTimeout, v))
}

// TestTimeoutLTE applies the LTE predicate on the "test_timeout" field.
func TestTimeoutLTE(v int64) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldLTE(FieldTestTimeout, v))
}

// TestTimeoutIsNil applies the IsNil predicate on the "test_timeout" field.
func TestTimeoutIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldTestTimeout))
}

// TestTimeoutNotNil applies the NotNil predicate on the "test_timeout" field.
func TestTimeoutNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldTestTimeout))
}

// TestSizeEQ applies the EQ predicate on the "test_size" field.
func TestSizeEQ(v TestSize) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldEQ(FieldTestSize, v))
}

// TestSizeNEQ applies the NEQ predicate on the "test_size" field.
func TestSizeNEQ(v TestSize) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNEQ(FieldTestSize, v))
}

// TestSizeIn applies the In predicate on the "test_size" field.
func TestSizeIn(vs ...TestSize) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIn(FieldTestSize, vs...))
}

// TestSizeNotIn applies the NotIn predicate on the "test_size" field.
func TestSizeNotIn(vs ...TestSize) predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotIn(FieldTestSize, vs...))
}

// TestSizeIsNil applies the IsNil predicate on the "test_size" field.
func TestSizeIsNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldIsNull(FieldTestSize))
}

// TestSizeNotNil applies the NotNil predicate on the "test_size" field.
func TestSizeNotNil() predicate.TargetComplete {
	return predicate.TargetComplete(sql.FieldNotNull(FieldTestSize))
}

// HasTargetPair applies the HasEdge predicate on the "target_pair" edge.
func HasTargetPair() predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, TargetPairTable, TargetPairColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTargetPairWith applies the HasEdge predicate on the "target_pair" edge with a given conditions (other predicates).
func HasTargetPairWith(preds ...predicate.TargetPair) predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := newTargetPairStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasImportantOutput applies the HasEdge predicate on the "important_output" edge.
func HasImportantOutput() predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ImportantOutputTable, ImportantOutputColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasImportantOutputWith applies the HasEdge predicate on the "important_output" edge with a given conditions (other predicates).
func HasImportantOutputWith(preds ...predicate.TestFile) predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := newImportantOutputStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDirectoryOutput applies the HasEdge predicate on the "directory_output" edge.
func HasDirectoryOutput() predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, DirectoryOutputTable, DirectoryOutputColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDirectoryOutputWith applies the HasEdge predicate on the "directory_output" edge with a given conditions (other predicates).
func HasDirectoryOutputWith(preds ...predicate.TestFile) predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := newDirectoryOutputStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOutputGroup applies the HasEdge predicate on the "output_group" edge.
func HasOutputGroup() predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, OutputGroupTable, OutputGroupColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOutputGroupWith applies the HasEdge predicate on the "output_group" edge with a given conditions (other predicates).
func HasOutputGroupWith(preds ...predicate.OutputGroup) predicate.TargetComplete {
	return predicate.TargetComplete(func(s *sql.Selector) {
		step := newOutputGroupStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TargetComplete) predicate.TargetComplete {
	return predicate.TargetComplete(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TargetComplete) predicate.TargetComplete {
	return predicate.TargetComplete(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TargetComplete) predicate.TargetComplete {
	return predicate.TargetComplete(sql.NotPredicates(p))
}