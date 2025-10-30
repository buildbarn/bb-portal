package detectors

import (
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/pkg/events"
)

// FailedTargetBazelInvocationProblemDetector struct
type FailedTargetBazelInvocationProblemDetector struct{}

// GetProblems implementation for FailedTargetBazelInvocationProblemDetector
func (FailedTargetBazelInvocationProblemDetector) GetProblems(event *events.BuildEvent) ([]Problem, error) {
	if event == nil || !isFailedTarget(event) {
		return nil, nil
	}
	label := event.GetTargetCompletedLabel()
	if label == "" {
		return nil, nil
	}
	problems, err := createProblem(BazelInvocationProblemFailedTarget, label, []*events.BuildEvent{event})
	if err != nil {
		return nil, err
	}
	return []Problem{*problems}, nil
}

// isFailedTarget
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
