package detectors

import "github.com/buildbarn/bb-portal/pkg/events"

// BazelInvocationProblemDetector interface
type BazelInvocationProblemDetector interface {
	ProcessBEPEvent(*events.BuildEvent)
	Problems() ([]Problem, error)
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
			NewActionProblemDetector(),
			NewTestProblemDetector(),
		},
		fallbackDetectors: []BazelInvocationProblemDetector{
			FailedTargetBazelInvocationProblemDetector{},
			&ErrorProgressBazelInvocationProblemDetector{},
		},
	}
}

// ProcessBEPEvent function
func (p ProblemDetector) ProcessBEPEvent(event *events.BuildEvent) {
	for _, detector := range p.detectors {
		detector.ProcessBEPEvent(event)
	}
	for _, detector := range p.fallbackDetectors {
		detector.ProcessBEPEvent(event)
	}
}

// Problems function
func (p ProblemDetector) Problems() ([]Problem, error) {
	problems, err := p.detectorsProblems(p.detectors)
	if err != nil {
		return nil, err
	}
	if !p.detectorsFoundProblems(p.detectors) {
		// It's OK if we add error progress even when no failed target (than it's a sole fallback, but very unlikely).
		fallbackProblems, err := p.detectorsProblems(p.fallbackDetectors)
		if err != nil {
			return nil, err
		}
		problems = append(problems, fallbackProblems...)
	}
	return problems, nil
}

// detectorsProblems
func (p ProblemDetector) detectorsProblems(detectors []BazelInvocationProblemDetector) ([]Problem, error) {
	var problems []Problem
	for _, detector := range detectors {
		detectorProblems, err := detector.Problems()
		if err != nil {
			return nil, err
		}
		problems = append(problems, detectorProblems...)
	}
	return problems, nil
}

// detectorsFoundProblems
func (p ProblemDetector) detectorsFoundProblems(detectors []BazelInvocationProblemDetector) bool {
	for _, detector := range detectors {
		problems, _ := detector.Problems()
		if len(problems) > 0 {
			return true
		}
	}
	return false
}
