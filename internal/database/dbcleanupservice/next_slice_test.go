package dbcleanupservice_test

import (
	"testing"

	"github.com/buildbarn/bb-portal/internal/database/dbcleanupservice"
	"github.com/stretchr/testify/require"
)

func TestCalculateNextSlice(t *testing.T) {
	tests := []struct {
		name          string
		counter       int64
		totalPages    int64
		runsPerHour   int64
		expectedFrom  int64
		expectedCount int64
	}{
		// name, counter, totalPages, runsPerHour, expectedFrom, expectedCount
		{"counter 0", 0, 100, 10, 0, 10},
		{"runs > pages", 0, 5, 10, 0, 1},
		{"pages 1", 0, 1, 10, 0, 1},
		{"mid cycle step 1", 1, 100, 10, 10, 10},
		{"mid cycle step 2", 2, 100, 10, 20, 10},
		{"perfect wrap", 10, 100, 10, 0, 10},
		{"post perfect wrap", 11, 100, 10, 10, 10},
		{"pre stretch wrap", 3, 10, 3, 9, 3},
		{"stretch wrap trigger", 4, 10, 3, 0, 5},
		{"post stretch wrap", 5, 10, 3, 5, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			from, count := dbcleanupservice.CalculateNextSlice(test.counter, test.totalPages, test.runsPerHour)
			require.Equal(t, test.expectedFrom, from)
			require.Equal(t, test.expectedCount, count)
		})
	}
}
