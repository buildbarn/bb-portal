package detectors

import (
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

type FailedTargetBazelInvocationProblemDetector map[string]*events.BuildEvent

func (f FailedTargetBazelInvocationProblemDetector) ProcessBEPEvent(event *events.BuildEvent) {
	if event == nil || !isFailedTarget(event) {
		return
	}
	label := event.GetTargetCompletedLabel()
	if label == "" {
		return
	}
	f[label] = event
}

func (f FailedTargetBazelInvocationProblemDetector) Problems() ([]Problem, error) {
	if len(f) == 0 {
		return nil, nil
	}
	problems := make([]Problem, 0, len(f))
	for label, event := range f {
		buildEvents := []*events.BuildEvent{event}
		problem, err := createProblem(BazelInvocationProblemFailedTarget, label, buildEvents)
		if err != nil {
			return nil, err
		}
		problems = append(problems, *problem)
	}
	return problems, nil
}

func isFailedTarget(event *events.BuildEvent) bool {
	if !event.IsTargetCompleted() {
		return false
	}
	completed := event.GetCompleted()
	if completed != nil && !completed.GetSuccess() {
		return true
	}
	aborted := event.GetAborted()
	if aborted != nil && aborted.GetReason() == bes.Aborted_ANALYSIS_FAILURE {
		return true
	}
	return false
}
