package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/graphql/model"
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

// Error helpers.
var (
	ErrOnlyURLOrUUID      = errors.New("either buildURL or buildUUID variable must be used, but not both")
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

// GetInt64Pointer returns a pointer to an int64
func GetInt64Pointer(i *int64) *int64 {
	return i
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

		problems = append(problems, problem)
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
		status, err := helper.Status()
		if err != nil {
			return nil, fmt.Errorf("could not get status: %w", err)
		}
		results, err := helper.Results()
		if err != nil {
			return nil, fmt.Errorf("could not get results: %w", err)
		}
		return &model.TestProblem{
			ID:      GraphQLIDFromTypeAndID("TestProblem", problem.ID),
			Label:   problem.Label,
			Status:  status,
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
func (problem testProblemHelper) Status() (string, error) {
	bepEvents, err := events.FromJSONArray(problem.BepEvents)
	if err != nil {
		return "", fmt.Errorf("failed to create test problem results: %w", err)
	}
	for _, event := range bepEvents {
		if event.IsTestSummary() {
			return event.GetTestSummary().GetOverallStatus().String(), nil
		}
	}
	return "", errStatusNotFound
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
