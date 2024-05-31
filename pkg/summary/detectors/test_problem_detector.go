package detectors

import (
	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

type TestProblemDetector struct {
	testSummaries map[labelKey]*events.BuildEvent
	testResults   map[labelKey]*events.BuildEvent
	outputBlobs   map[labelKey][]BlobURI
}

func NewTestProblemDetector() TestProblemDetector {
	return TestProblemDetector{
		testSummaries: map[labelKey]*events.BuildEvent{},
		testResults:   map[labelKey]*events.BuildEvent{},
		outputBlobs:   map[labelKey][]BlobURI{},
	}
}

func (t TestProblemDetector) ProcessBEPEvent(event *events.BuildEvent) {
	if event.IsTestSummary() {
		// Keep only non-successful test summaries.
		if event.GetTestSummary().GetOverallStatus() == bes.TestStatus_PASSED {
			return
		}
		label := event.GetId().GetTestSummary().GetLabel()
		key := labelKey(label)
		t.testSummaries[key] = event
	} else if event.IsTestResult() {
		// Keep only the latest non-successful test results.
		if event.GetTestResult().GetStatus() == bes.TestStatus_PASSED {
			return
		}
		label := event.GetId().GetTestResult().GetLabel()
		key := labelKey(label)
		t.testResults[key] = event

		// Keep test result's blobs.
		outputs := event.GetTestResult().GetTestActionOutput()
		if outputs == nil {
			return
		}
		t.outputBlobs[key] = getOutputsBlobs(outputs)
	}
}

func (t TestProblemDetector) Problems() ([]Problem, error) {
	problems := make([]Problem, 0, len(t.testResults))
	for label, testSummary := range t.testSummaries {
		// Group test summary and result into a problem's buildEvents.
		buildEvents := []*events.BuildEvent{testSummary}
		testResult, ok := t.testResults[label]
		// If no test result was saved, ignore test summary.
		if !ok {
			continue
		}
		buildEvents = append(buildEvents, testResult)

		// Create problem.
		problem, err := createProblemWithBlobs(BazelInvocationTestProblem, label, buildEvents, t.outputBlobs)
		if err != nil {
			return nil, err
		}
		problems = append(problems, *problem)
	}
	return problems, nil
}
