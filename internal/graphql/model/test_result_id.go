package model

// TestResultID A Test Result ID struct.
type TestResultID struct {
	// Fields need to be exported since GraphQL relay (un)marshaling is using JSON.
	ProblemID uint64 `json:"problem"`
	Run       int32  `json:"run"`
	Shard     int32  `json:"shard"`
	Attempt   int32  `json:"attempt"`
}
