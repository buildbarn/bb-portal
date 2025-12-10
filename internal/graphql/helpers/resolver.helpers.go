package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"entgo.io/contrib/entgql"
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testsummary"
	"github.com/buildbarn/bb-portal/internal/graphql/model"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

// Error helpers.
var (
	ErrWrongType          = errors.New("received unexpected type while trying to convert node to *ent.BazelInvocationProblem")
	errUnknownProblemType = errors.New("unknown problem type")
	errActionNotFound     = errors.New("action not found")
	errStatusNotFound     = errors.New("status not found")
)

// StringSliceArrayToPointerArray takes an array of strings and returns an array of string pointers
func StringSliceArrayToPointerArray(strings []string) []*string {
	result := make([]*string, len(strings))
	for i, str := range strings {
		result[i] = &str
	}
	return result
}

// GetFloatPointer returns a pointer to a float
func GetFloatPointer(f *float64) *float64 {
	return f
}

// Helper A Helper struct.
type Helper struct {
	*problemHelper
}

// NewHelper Initializer for helper
func NewHelper() *Helper {
	return &Helper{
		problemHelper: &problemHelper{},
	}
}

// A problem helper.
type problemHelper struct{}

// DBProblemsToAPIProblems Convert db problem to api problem.
func (ph problemHelper) DBProblemsToAPIProblems(ctx context.Context, dbProblems []*ent.BazelInvocationProblem) ([]model.Problem, error) {
	problems := make([]model.Problem, 0, len(dbProblems))
	for _, dbProblem := range dbProblems {
		problem, err := ph.DBProblemToAPIProblem(ctx, dbProblem)
		if err != nil {
			return nil, err
		}
		if problem != nil {
			problems = append(problems, problem)
		}
	}
	return problems, nil
}

// DBProblemToAPIProblem Convert a DB problem to an API problem.
func (ph problemHelper) DBProblemToAPIProblem(ctx context.Context, problem *ent.BazelInvocationProblem) (model.Problem, error) {
	switch problem.ProblemType {
	case detectors.BazelInvocationActionProblem:
		actionType, err := ph.getActionType(ctx, problem)
		if err != nil {
			return nil, fmt.Errorf("could not get action type: %w", err)
		}
		return ph.actionProblemFromDBModel(problem, actionType), nil

	case detectors.BazelInvocationTestProblem:
		helper := testProblemHelper{BazelInvocationProblem: problem}
		status, err := helper.Status(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get status: %w", err)
		}
		if status == testsummary.OverallStatusPASSED {
			return nil, nil
		}
		results, err := helper.Results()
		if err != nil {
			return nil, fmt.Errorf("could not get results: %w", err)
		}
		return &model.TestProblem{
			ID:      GraphQLIDFromTypeAndID("TestProblem", problem.ID),
			Label:   problem.Label,
			Status:  status.String(),
			Results: results,
		}, nil

	case detectors.BazelInvocationProblemFailedTarget:
		return &model.TargetProblem{
			ID:    GraphQLIDFromTypeAndID("TargetProblem", problem.ID),
			Label: problem.Label,
		}, nil

	case detectors.BazelInvocationProblemErrorProgress:
		helper := progressProblemHelper{problem}
		output, err := helper.Output()
		if err != nil {
			return nil, fmt.Errorf("could not get output: %w", err)
		}
		return &model.ProgressProblem{
			ID:     GraphQLIDFromTypeAndID("ProgressProblem", problem.ID),
			Output: output,
		}, nil

	default:
		return nil, fmt.Errorf("unknown type: %s: %w", problem.ProblemType, errUnknownProblemType)
	}
}

// Get an action type.
func (ph problemHelper) getActionType(ctx context.Context, problem *ent.BazelInvocationProblem) (string, error) {
	action, err := ph.getAction(ctx, problem)
	if err != nil {
		return "", err
	}
	return action.GetType(), nil
}

// Get an action.
func (ph problemHelper) getAction(ctx context.Context, problem *ent.BazelInvocationProblem) (*bes.ActionExecuted, error) {
	bepEvents, err := events.FromJSONArray(problem.BepEvents)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal action problem events", "problem", problem)
		return nil, fmt.Errorf("failed to parse action problem: %w", err)
	}
	for _, event := range bepEvents {
		if event.IsActionCompleted() {
			return event.GetAction(), nil
		}
	}
	return nil, errActionNotFound
}

// Get an action problem form a database model.
func (ph problemHelper) actionProblemFromDBModel(problem *ent.BazelInvocationProblem, actionType string) model.Problem {
	return &model.ActionProblem{
		ID:      GraphQLIDFromTypeAndID("ActionProblem", problem.ID),
		Label:   problem.Label,
		Type:    actionType,
		Problem: problem,
	}
}

// A test problem helper struct
type testProblemHelper struct {
	*ent.BazelInvocationProblem
}

// Get the graphql id.
func (problem testProblemHelper) GraphQLID() string {
	// TODO: scalars.GraphQLIDFromString
	return fmt.Sprintf("testProblem:%d", problem.ID)
}

// get the status of the problm helper.
func (problem testProblemHelper) Status(ctx context.Context) (testsummary.OverallStatus, error) {
	testSummary, err := problem.QueryBazelInvocation().QueryTestCollection().QueryTestSummary().Where(testsummary.LabelEQ(problem.Label)).Select(testsummary.FieldOverallStatus).Only(ctx)
	return testSummary.OverallStatus, err
}

// The results.
func (problem testProblemHelper) Results() ([]*model.TestResult, error) {
	bepEvents, err := events.FromJSONArray(problem.BepEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to create test problem results: %w", err)
	}
	var results []*model.TestResult
	for _, event := range bepEvents {
		if event.IsTestResult() {
			testResultEventID := event.GetId().GetTestResult()
			helper := testResultOverviewHelper{
				TestResult: event.GetTestResult(),
				testResultID: model.TestResultID{
					ProblemID: uint64(problem.ID), //nolint:gosec
					Run:       testResultEventID.GetRun(),
					Shard:     testResultEventID.GetShard(),
					Attempt:   testResultEventID.GetAttempt(),
				},
			}
			result := model.TestResult{
				ID:            GraphQLIDFromTypeAndID("TestResult", problem.ID),
				Run:           int(helper.Run()),
				Shard:         int(helper.Shard()),
				Attempt:       int(helper.Attempt()),
				Status:        helper.Status(),
				BESTestResult: event.GetTestResult(),
			}

			results = append(results, &result)
		}
	}
	return results, nil
}

// The Progress problem helper.
type progressProblemHelper struct {
	*ent.BazelInvocationProblem
}

// The Output.
func (e progressProblemHelper) Output() (string, error) {
	bepEvents, err := events.FromJSONArray(e.BepEvents)
	if err != nil {
		return "", fmt.Errorf("failed to create error progress output: %w", err)
	}
	output := strings.Builder{}
	for _, event := range bepEvents {
		stderr := event.GetProgress().GetStderr()
		output.WriteString(stderr)
	}
	return output.String(), nil
}

// Test Result Overview Helper.
type testResultOverviewHelper struct {
	*bes.TestResult
	testResultID model.TestResultID
}

// The Run property.
func (helper testResultOverviewHelper) Run() int32 {
	return helper.testResultID.Run
}

// The Shard property.
func (helper testResultOverviewHelper) Shard() int32 {
	return helper.testResultID.Shard
}

// The Attempt.
func (helper testResultOverviewHelper) Attempt() int32 {
	return helper.testResultID.Attempt
}

// The Status.
func (helper testResultOverviewHelper) Status() string {
	return helper.GetStatus().String()
}

func paginationCursorToUTC(cursor *entgql.Cursor[int]) {
	if cursor == nil || cursor.Value == nil {
		return
	}
	switch v := cursor.Value.(type) {
	case time.Time:
		cursor.Value = v.UTC()
	case *time.Time:
		if v != nil {
			ut := v.UTC()
			cursor.Value = &ut
		}
	}
}

// PaginationCursorsToUTC converts pagination cursors that consist of
// timestamps to UTC instead of local time. When the backend sends the cursors
// to the frontend, they are in UTC. However, when the frontend sends them
// back, they are interpreted as local time. This causes issues since Sqlite
// cannot handle comparisons between timestamps in different timezones.
func PaginationCursorsToUTC(after, before *entgql.Cursor[int]) {
	paginationCursorToUTC(after)
	paginationCursorToUTC(before)
	return
}
