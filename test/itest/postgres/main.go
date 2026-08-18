package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

const (
	postgresPort     = 15432
	postgresDatabase = "postgres"
	postgresUser     = "app"
	postgresPassword = "password"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "PostgreSQL service failed: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	binariesPath, err := runfiles.Rlocation("com_github_buildbarn_bb_portal/internal/database/embedded/extraced_embedded_postgres.extracted")
	if err != nil {
		return fmt.Errorf("locate embedded PostgreSQL binaries: %w", err)
	}

	temporaryDirectory := os.Getenv("TEST_TMPDIR")
	if temporaryDirectory == "" {
		return fmt.Errorf("TEST_TMPDIR is not set; run this binary through rules_itest")
	}
	dataPath := filepath.Join(temporaryDirectory, "postgres_data")
	if err := os.RemoveAll(dataPath); err != nil {
		return fmt.Errorf("remove old PostgreSQL data directory: %w", err)
	}
	passwordFile := filepath.Join(temporaryDirectory, "postgres_password")
	if err := os.WriteFile(passwordFile, []byte(postgresPassword), 0o600); err != nil {
		return fmt.Errorf("write PostgreSQL password file: %w", err)
	}

	initdbPath := filepath.Join(binariesPath, "bin", "initdb")
	initdb := exec.Command(
		initdbPath,
		"-A", "password",
		"-U", postgresUser,
		"-D", dataPath,
		"--pwfile="+passwordFile,
	)
	initdb.Stdout = os.Stdout
	initdb.Stderr = os.Stderr
	if err := initdb.Run(); err != nil {
		return fmt.Errorf("initialize PostgreSQL: %w", err)
	}
	if err := os.Remove(passwordFile); err != nil {
		return fmt.Errorf("remove PostgreSQL password file: %w", err)
	}

	socketDirectory := os.Getenv("SOCKET_DIR")
	if socketDirectory == "" {
		socketDirectory = temporaryDirectory
	}
	postgresPath := filepath.Join(binariesPath, "bin", "postgres")
	arguments := []string{
		postgresPath,
		"-D", dataPath,
		"-p", strconv.Itoa(postgresPort),
		"-k", socketDirectory,
		"-c", "listen_addresses=127.0.0.1",
	}
	fmt.Printf("Starting PostgreSQL at postgresql://%s@127.0.0.1:%d/%s\n", postgresUser, postgresPort, postgresDatabase)

	// Replace the launcher with PostgreSQL instead of starting it through
	// pg_ctl. This keeps PostgreSQL in rules_itest's process group, ensuring it
	// cannot survive when the service manager is interrupted or killed.
	return syscall.Exec(postgresPath, arguments, os.Environ())
}
