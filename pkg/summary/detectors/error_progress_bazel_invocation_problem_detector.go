package detectors

import (
	"strings"

	"github.com/buildbarn/bb-portal/pkg/events"
)

type ErrorProgressBazelInvocationProblemDetector []*events.BuildEvent

func (e *ErrorProgressBazelInvocationProblemDetector) ProcessBEPEvent(event *events.BuildEvent) {
	if event == nil || event.GetProgress() == nil {
		return
	}
	stderr := event.GetProgress().GetStderr()
	if strings.HasPrefix(stderr, "ERROR: ") || strings.Contains(stderr, "\nERROR: ") {
		*e = append(*e, event)
	}
}

func (e *ErrorProgressBazelInvocationProblemDetector) Problems() ([]Problem, error) {
	if len(*e) == 0 {
		return nil, nil
	}
	// Progress is not a labeled event.
	problem, err := createProblem(BazelInvocationProblemErrorProgress, "", *e)
	if err != nil {
		return nil, err
	}
	return []Problem{*problem}, nil
}
