package detectors

import (
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

// ActionProblemDetector struct
type ActionProblemDetector struct {
	actionsCompleted map[labelKey]*events.BuildEvent
	outputBlobs      map[labelKey][]BlobURI
}

// NewActionProblemDetector function
func NewActionProblemDetector() ActionProblemDetector {
	return ActionProblemDetector{
		actionsCompleted: map[labelKey]*events.BuildEvent{},
		outputBlobs:      map[labelKey][]BlobURI{},
	}
}

// ProcessBEPEvent function
func (a ActionProblemDetector) ProcessBEPEvent(event *events.BuildEvent) {
	if event == nil || !isFailedAction(event) {
		return
	}
	label := event.GetActionCompletedLabel()
	if label == "" {
		return
	}

	// Ignore actions that don't have any output.
	outputs := getActionOutputs(event)
	if len(outputs) == 0 {
		return
	}

	key := labelKey(label)
	// Save action.
	a.actionsCompleted[key] = event

	// Save action's blobs.
	a.outputBlobs[key] = getOutputsBlobs(outputs)
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

// Problems function
func (a ActionProblemDetector) Problems() ([]Problem, error) {
	if len(a.actionsCompleted) == 0 {
		return nil, nil
	}
	problems := make([]Problem, 0, len(a.actionsCompleted))
	for labelKey, event := range a.actionsCompleted {
		// Create problem.
		buildEvents := []*events.BuildEvent{event}
		problem, err := createProblemWithBlobs(BazelInvocationActionProblem, labelKey, buildEvents, a.outputBlobs)
		if err != nil {
			return nil, err
		}

		problems = append(problems, *problem)
	}
	return problems, nil
}

// isFailedAction
func isFailedAction(event *events.BuildEvent) bool {
	if !event.IsActionCompleted() {
		return false
	}
	action := event.GetAction()
	return action != nil && !action.GetSuccess()
}
