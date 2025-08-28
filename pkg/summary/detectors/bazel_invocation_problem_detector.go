package detectors

import "github.com/buildbarn/bb-portal/pkg/events"

// BazelInvocationProblemDetector interface
type BazelInvocationProblemDetector interface {
	GetProblems(*events.BuildEvent) ([]Problem, error)
}

// ProblemDetector struct
type ProblemDetector struct {
	detectors         []BazelInvocationProblemDetector
	fallbackDetectors []BazelInvocationProblemDetector
}

// NewProblemDetector function
func NewProblemDetector() ProblemDetector {
	return ProblemDetector{
		detectors: []BazelInvocationProblemDetector{
			ActionProblemDetector{},
			TestProblemDetector{},
		},
		fallbackDetectors: []BazelInvocationProblemDetector{
			FailedTargetBazelInvocationProblemDetector{},
			ErrorProgressBazelInvocationProblemDetector{},
		},
	}
}

// GetProblems implementation for ProblemDetector
func (p ProblemDetector) GetProblems(event *events.BuildEvent) ([]Problem, error) {
	var problems []Problem
	for _, detector := range p.detectors {
		detectorProblems, err := detector.GetProblems(event)
		if err != nil {
			return nil, err
		}
		problems = append(problems, detectorProblems...)
	}
	if len(problems) == 0 {
		for _, detector := range p.fallbackDetectors {
			detectorProblems, err := detector.GetProblems(event)
			if err != nil {
				return nil, err
			}
			problems = append(problems, detectorProblems...)
		}
	}
	return problems, nil
}
