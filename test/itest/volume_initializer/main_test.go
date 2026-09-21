package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitializeDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "directory")
	require.NoError(t, initializeDirectories([]directory{{
		path: path,
		mode: 0o750,
		uid:  -1,
		gid:  -1,
	}}))

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.True(t, info.IsDir())
	require.Equal(t, os.FileMode(0o750), info.Mode().Perm())
}
