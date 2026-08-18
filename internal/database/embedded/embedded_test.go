package embedded

import (
	"errors"
	"os"
	"testing"
)

func TestCleanupRemovesTemporaryDirectories(t *testing.T) {
	runtimePath := t.TempDir()
	dataPath := t.TempDir()
	databaseProvider := &DatabaseProvider{
		runtimePath: runtimePath,
		dataPath:    dataPath,
	}

	if err := databaseProvider.Cleanup(); err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}
	for _, path := range []string{runtimePath, dataPath} {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Temporary directory %q still exists after cleanup", path)
		}
	}
}
