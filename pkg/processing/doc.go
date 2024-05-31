// Package processing contains the high-level workflows an actions that are used to analyze Bazel invocations.
//
// This package what wires together the backend processes. This where orchestration and integration happens
// between the filesystem, summarization functions and database.
//
// This package is responsible for detecting, processing & storing new data. It is not responsible for queries
// against data that has already been stored.
package processing
