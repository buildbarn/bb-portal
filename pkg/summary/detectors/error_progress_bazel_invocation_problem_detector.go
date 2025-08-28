package detectors

import (
	"strings"

	"github.com/buildbarn/bb-portal/pkg/events"
)

// ErrorProgressBazelInvocationProblemDetector type
type ErrorProgressBazelInvocationProblemDetector struct{}

// GetProblems implementation for ErrorProgressBazelInvocationProblemDetector
func (ErrorProgressBazelInvocationProblemDetector) GetProblems(event *events.BuildEvent) ([]Problem, error) {
	if event == nil || event.GetProgress() == nil {
		return nil, nil
	}
	stderr := event.GetProgress().GetStderr()
	if strings.HasPrefix(stderr, "ERROR: ") || strings.Contains(stderr, "\nERROR: ") {
		problem, err := createProblem(BazelInvocationProblemErrorProgress, "", []*events.BuildEvent{event})
		if err != nil {
			return nil, err
		}
		return []Problem{*problem}, nil
	}
	return nil, nil
}
