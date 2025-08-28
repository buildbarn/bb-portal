package bes

import (
	"context"
	"strings"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// GetNormalizedIncompleteBuildLogs retrieves all incomplete build logs
// associated with the given Bazel invocation, concatenates them, normalizes
// them by removing ANSI escape sequences that remove previous lines, and
// returns the result.
func GetNormalizedIncompleteBuildLogs(ctx context.Context, invocation *ent.BazelInvocation) (string, error) {
	incompleteLogs, err := invocation.QueryIncompleteBuildLogs().Order(ent.Asc("snippet_id")).All(ctx)
	if err != nil {
		return "", util.StatusWrap(err, "Failed to query incomplete build logs from database")
	}

	var sb strings.Builder
	for _, logSnippet := range incompleteLogs {
		sb.WriteString(logSnippet.LogSnippet)
	}

	result := []string{}
	for _, line := range strings.Split(sb.String(), "\n") {
		for {
			// ANSI escape sequence for removing the previous line
			index := strings.Index(line, "\r\x1b[1A\x1b[K")
			if index == -1 {
				break
			}
			line = line[index+8:]
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n"), nil
}
