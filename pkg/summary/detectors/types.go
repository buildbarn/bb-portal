package detectors

import (
	"encoding/json"
)

type labelKey string

type BazelInvocationProblemType string

type BlobURI string

type NamedBlob struct {
	BlobURI
	Name string
}

type Problem struct {
	//*ent.BazelInvocationProblem
	DetectedBlobs []BlobURI
	ProblemType   BazelInvocationProblemType
	Label         string
	BEPEvents     json.RawMessage
}
