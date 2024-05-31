// Package summary our model that summarizes a Bazel Invocation or set of invocations.
//
// The summary package contains functions to produce a summary from an event stream, and so it depends
// on pkg/events.
//
// This package contains types to hold a summary in-memory. It is not responsible for persisting
// or viewing summary documents. Therefore, it should not have many dependencies except for pkg/events.
// It should not depend on libraries or frameworks related to databases, HTML, GraphQL, PDF, email, etc. Other
// packages can provide functions that take an in-memory summary & use those libraries to  persist/render/display it.
package summary
