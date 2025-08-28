package detectors

import (
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/pkg/events"
)

// ActionProblemDetector struct
type ActionProblemDetector struct{}

// GetProblems implementation for ActionProblemDetector
func (a ActionProblemDetector) GetProblems(event *events.BuildEvent) ([]Problem, error) {
	if event == nil || !isFailedAction(event) {
		return nil, nil
	}
	label := event.GetActionCompletedLabel()
	if label == "" {
		return nil, nil
	}

	// Ignore actions that don't have any output.
	output := getActionOutputs(event)
	if len(output) == 0 {
		return nil, nil
	}
	outputs := map[labelKey][]BlobURI{}
	outputs[labelKey(label)] = getOutputsBlobs(output)

	problem, err := createProblemWithBlobs(BazelInvocationActionProblem, labelKey(label), []*events.BuildEvent{event}, outputs)
	if err != nil {
		return nil, err
	}
	return []Problem{*problem}, nil
}

// getActionOutputs
func getActionOutputs(event *events.BuildEvent) []*bes.File {
	action := event.GetAction()
	if action == nil {
		return nil
	}
	outputs := []*bes.File{}
	if action.GetStderr() != nil {
		outputs = append(outputs, action.GetStderr())
	}
	if action.GetStdout() != nil {
		outputs = append(outputs, action.GetStdout())
	}
	return outputs
}

// getOutputsBlobs
func getOutputsBlobs(outputs []*bes.File) []BlobURI {
	if outputs == nil {
		return nil
	}
	blobs := []BlobURI{}
	for _, output := range outputs {
		blobs = append(blobs, BlobURI(output.GetUri()))
	}
	return blobs
}

// isFailedAction
func isFailedAction(event *events.BuildEvent) bool {
	if !event.IsActionCompleted() {
		return false
	}
	action := event.GetAction()
	return action != nil && !action.GetSuccess()
}
