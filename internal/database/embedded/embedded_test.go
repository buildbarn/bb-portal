package embedded

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanupRemovesTemporaryDirectories(t *testing.T) {
	runtimePath := t.TempDir()
	dataPath := t.TempDir()
	databaseProvider := &DatabaseProvider{
		runtimePath: runtimePath,
		dataPath:    dataPath,
	}

	require.NoError(t, databaseProvider.Cleanup())
	for _, path := range []string{runtimePath, dataPath} {
		_, err := os.Stat(path)
		require.Equal(t, true, errors.Is(err, os.ErrNotExist), "Temporary directory %q still exists after cleanup", path)
	}
}
