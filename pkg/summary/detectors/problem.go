package detectors

import (
	"github.com/buildbarn/bb-portal/pkg/events"
)

const (
	BazelInvocationProblemFailedTarget  = "FAILED_TARGET"
	BazelInvocationProblemErrorProgress = "ERROR_PROGRESS"
	BazelInvocationTestProblem          = "TEST_PROBLEM"
	BazelInvocationActionProblem        = "ACTION_PROBLEM"
)

func createProblem(problemType BazelInvocationProblemType, label string, buildEvents []*events.BuildEvent) (*Problem, error) {
	bepEvents, err := events.AsJSONArray(buildEvents)
	if err != nil {
		return nil, err
	}

	return &Problem{
		ProblemType: problemType,
		Label:       label,
		BEPEvents:   bepEvents,
	}, nil
}

func createProblemWithBlobs(problemType BazelInvocationProblemType, key labelKey, events []*events.BuildEvent, outputBlobs map[labelKey][]BlobURI) (*Problem, error) {
	problem, err := createProblem(problemType, string(key), events)
	if err != nil {
		return nil, err
	}

	// Set blobs referenced by the problem.
	problem.DetectedBlobs = outputBlobs[key]
	return problem, nil
}
