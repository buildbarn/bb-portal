package detectors

import (
	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/pkg/events"
)

// TestProblemDetector struct
type TestProblemDetector struct{}

// GetProblems implementation for TestProblemDetector
func (t TestProblemDetector) GetProblems(event *events.BuildEvent) ([]Problem, error) {
	if event == nil || !event.IsTestResult() {
		return nil, nil
	}
	if event.GetTestResult().GetStatus() == bes.TestStatus_PASSED {
		return nil, nil
	}

	label := event.GetId().GetTestResult().GetLabel()
	outputs := map[labelKey][]BlobURI{}
	output := event.GetTestResult().GetTestActionOutput()
	if output == nil {
		return nil, nil
	}
	outputs[labelKey(label)] = getOutputsBlobs(output)

	problem, err := createProblemWithBlobs(BazelInvocationTestProblem, labelKey(label), []*events.BuildEvent{event}, outputs)
	if err != nil {
		return nil, err
	}
	return []Problem{*problem}, nil
}
