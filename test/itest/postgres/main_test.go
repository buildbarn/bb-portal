package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureRemoteHostEntry(t *testing.T) {
	dataPath := t.TempDir()
	hbaPath := filepath.Join(dataPath, "pg_hba.conf")
	require.NoError(t, os.WriteFile(hbaPath, []byte("local all all trust"), 0o600))

	changed, err := ensureRemoteHostEntry(dataPath)
	require.NoError(t, err)
	require.True(t, changed)

	changed, err = ensureRemoteHostEntry(dataPath)
	require.NoError(t, err)
	require.False(t, changed)

	contents, err := os.ReadFile(hbaPath)
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(string(contents), remoteHostEntry))
}
