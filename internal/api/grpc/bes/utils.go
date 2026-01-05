package bes

import (
	"bytes"
	"context"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/incompletebuildlog"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
)

// GetNormalizedIncompleteBuildLogs retrieves all incomplete build logs
// associated with the given Bazel invocation, concatenates them, normalizes
// them by removing ANSI escape sequences that remove previous lines, and
// returns the result.
func GetNormalizedIncompleteBuildLogs(ctx context.Context, client *ent.Client, invocationID uuid.UUID) ([]byte, error) {
	incompleteLogs, err := client.IncompleteBuildLog.Query().Where(
		incompletebuildlog.HasBazelInvocationWith(
			bazelinvocation.InvocationIDEQ(invocationID),
		),
	).Order(ent.Asc("snippet_id")).All(ctx)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to query incomplete build logs from database")
	}
	// Read all snippets into one big buffer.
	var buf bytes.Buffer
	for _, logSnippet := range incompleteLogs {
		_, err := buf.Write(logSnippet.LogSnippet)
		if err != nil {
			return nil, util.StatusWrap(err, "Could not write to an in memory buffer, is your computer on fire?")
		}

	}
	// Determine which lines to keep and where they start/end.
	var lineSpans [][2]int
	data := buf.Bytes()
	size := len(data)
	current := 0
	ansi := []byte("\r\x1b[1A\x1b[K")
	for current < size {
		start, end := current, 0
		endOffset := bytes.IndexByte(data[current:], '\n')
		var line []byte
		if endOffset == -1 {
			line = data[start:]
			end = size
		} else {
			// Include the \n in the line.
			end = start + endOffset + 1
			line = data[start:end]
		}
		current = end
		// Remove any earlier lines
		count := bytes.Count(line, ansi)
		if count > 0 {
			if len(lineSpans) >= count {
				lineSpans = lineSpans[:len(lineSpans)-count]
			} else {
				lineSpans = lineSpans[:0]
			}
			start = start + bytes.LastIndex(line, ansi) + len(ansi)
		}
		lineSpans = append(lineSpans, [2]int{start, end})
	}
	// Write all lines to keep to the beginning of the buffer.
	cursor := 0
	for _, span := range lineSpans {
		start, end := span[0], span[1]
		copy(data[cursor:], data[start:end])
		cursor += end - start
	}
	// Return resulting string.
	return data[:cursor], nil
}
