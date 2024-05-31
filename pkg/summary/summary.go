package summary

import (
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

const (
	// stepLabelKey is used in buildMetadata events to provide a human-readable label for build steps.
	stepLabelKey = "BUILD_STEP_LABEL"
)

const (
	ExitCodeSuccess     = 0
	ExitCodeInterrupted = 8
)

type Summary struct {
	*InvocationSummary
	Problems       []detectors.Problem
	RelatedFiles   map[string]string
	EventFileURL   string
	BEPCompleted   bool
	StartedAt      time.Time
	InvocationID   string
	StepLabel      string
	EndedAt        *time.Time
	ChangeNumber   int
	PatchsetNumber int
	BuildURL       string
	BuildUUID      uuid.UUID
}

type InvocationSummary struct {
	EnvVars          map[string]string
	ExitCode         *ExitCode
	BazelVersion     string
	BazelCommandLine BazelCommandLine
}

type ExitCode struct {
	Code int
	Name string
}

type BazelCommandLine struct {
	Executable string
	Command    string
	Residual   string
	Options    []string
}

// Blob holds information about a blob in the CAS. Should be easily converted to/from the one in the
// cas package. Copied into here so this package does not have *any* dependencies except standard
// libraries.
type Blob struct {
	BlobURI  url.URL
	Size     int
	Contents string
	Name     string
}
