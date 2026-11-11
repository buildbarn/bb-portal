package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// bbPortalProcess is a bb_portal binary started by startBbPortal.
type bbPortalProcess struct {
	// done is closed when the process has exited, after which err holds
	// the result of waiting for it.
	done        chan struct{}
	err         error
	countTables func() int
}

// startBbPortal starts the bb_portal binary with a database using the
// given schema management mode. An HTTP server is configured as well,
// so that bb_portal has something to serve: without any server, it
// exits as soon as it has started.
func startBbPortal(t *testing.T, schemaManagement string) *bbPortalProcess {
	owner := newTestConnection(t)
	var databaseName string
	var port int
	require.NoError(t, owner.QueryRow("SELECT current_database(), inet_server_port()").Scan(&databaseName, &port))
	// The embedded database does not expose the password of its
	// superuser, so the process connects as the owner of the database.
	role := "portal_" + uuid.NewString()[:8]
	execSQL(
		t, owner,
		fmt.Sprintf(`CREATE ROLE %q LOGIN PASSWORD 'password'`, role),
		fmt.Sprintf(`ALTER DATABASE %q OWNER TO %q`, databaseName, role),
		fmt.Sprintf(`ALTER SCHEMA public OWNER TO %q`, role),
	)

	configurationPath := filepath.Join(t.TempDir(), "bb_portal.jsonnet")
	require.NoError(t, os.WriteFile(configurationPath, []byte(fmt.Sprintf(`{
		httpServers: [{
			listenAddresses: ['127.0.0.1:0'],
			authenticationPolicy: { allow: {} },
		}],
		database: {
			postgres: { connectionString: 'postgres://%s:password@localhost:%d/%s?sslmode=disable' },
			schemaManagement: %s,
		},
	}`, role, port, databaseName, schemaManagement)), 0o600))

	binary, err := runfiles.Rlocation(os.Getenv("BB_PORTAL_BINARY"))
	require.NoError(t, err)
	cmd := exec.Command(binary, configurationPath)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Start())

	process := &bbPortalProcess{
		done:        make(chan struct{}),
		countTables: func() int { return countPublicTables(t, owner) },
	}
	go func() {
		process.err = cmd.Wait()
		close(process.done)
	}()
	t.Cleanup(func() {
		cmd.Process.Signal(syscall.SIGTERM)
		select {
		case <-process.done:
		case <-time.After(10 * time.Second):
			cmd.Process.Kill()
			<-process.done
		}
	})
	return process
}

func TestBbPortalApplyAndExitTerminates(t *testing.T) {
	process := startBbPortal(t, `{ applyAndExit: {} }`)
	select {
	case <-process.done:
		require.NoError(t, process.err, "bb_portal should exit with status 0")
	case <-time.After(60 * time.Second):
		t.Fatal("bb_portal did not exit")
	}
	require.NotZero(t, process.countTables(), "the schema should have been applied")
}

func TestBbPortalApplyKeepsRunning(t *testing.T) {
	process := startBbPortal(t, `{ apply: {} }`)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for process.countTables() == 0 {
		select {
		case <-process.done:
			t.Fatalf("bb_portal exited without applying the schema: %v", process.err)
		case <-ctx.Done():
			t.Fatal("bb_portal did not apply the schema")
		case <-time.After(100 * time.Millisecond):
		}
	}
	// Past schema setup, bb_portal must keep serving.
	select {
	case <-process.done:
		t.Fatalf("bb_portal exited after applying the schema: %v", process.err)
	case <-time.After(2 * time.Second):
	}
}
