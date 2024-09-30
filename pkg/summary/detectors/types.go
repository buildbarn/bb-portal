package detectors

import (
	"encoding/json"
)

// labelKey string
type labelKey string

// BazelInvocationProblemType string
type BazelInvocationProblemType string

// BlobURI string
type BlobURI string

// NamedBlob struct
type NamedBlob struct {
	BlobURI
	Name string
}

// Problem struct
type Problem struct {
	//*ent.BazelInvocationProblem
	DetectedBlobs []BlobURI
	ProblemType   BazelInvocationProblemType
	Label         string
	BEPEvents     json.RawMessage
}
