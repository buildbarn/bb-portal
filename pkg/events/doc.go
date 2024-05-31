// Package events provides high-level functions for reading Bazel Event Protocol files or streams.
//
// An event in this context means a Bazel Event Protocol event. These are defined by the Bazel project
// as protobuf. They are built as a Go library in github.com/buildbarn/bb-portal/third_party/bazel/gen/bes.
//
// This package may provide convenience functions for working with those events, such as:
// - Iterating over events in a line-delimited JSON file (NDJSON file).
// - Converting events to/from a JSON array in order to save them in a DB as JSON.
//
// This package should not contain any code to process or interpret events, and should not be
// aware of other types we define. This package should have very few dependencies. Ideally just
// the standard library & the libraries that are necessary to unmarshall the protobuf messages.
package events
