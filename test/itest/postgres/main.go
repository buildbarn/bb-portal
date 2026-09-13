package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

const (
	postgresBinariesRunfile = "com_github_buildbarn_bb_portal/internal/database/embedded/extraced_embedded_postgres.extracted"
	remoteHostEntry         = "host all all all scram-sha-256 # bb-portal integration test"
)

func main() {
	var err error
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		err = healthcheck()
	} else {
		err = serve()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve() error {
	binariesPath, err := postgresBinariesPath()
	if err != nil {
		return err
	}

	port, err := postgresPort()
	if err != nil {
		return err
	}
	dataPath := environmentOrDefault("PGDATA", "/var/lib/postgresql/data")
	configuration := embeddedpostgres.DefaultConfig().
		Port(port).
		Database(environmentOrDefault("POSTGRES_DB", "postgres")).
		Username(environmentOrDefault("POSTGRES_USER", "postgres")).
		Password(environmentOrDefault("POSTGRES_PASSWORD", "postgres")).
		RuntimePath(environmentOrDefault("POSTGRES_RUNTIME", "/tmp/embedded-postgres")).
		DataPath(dataPath).
		BinariesPath(binariesPath).
		Locale("C").
		Encoding("UTF8").
		StartParameters(map[string]string{
			"listen_addresses":        "*",
			"unix_socket_directories": "/tmp",
		}).
		BinaryRepositoryURL("DISABLED").
		StartTimeout(30 * time.Second).
		Logger(os.Stdout)

	database := embeddedpostgres.NewDatabase(configuration)
	if err := database.Start(); err != nil {
		return fmt.Errorf("start embedded PostgreSQL: %w", err)
	}
	if err := allowRemoteConnections(binariesPath, dataPath); err != nil {
		return errors.Join(err, database.Stop())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)
	received := <-signals
	fmt.Fprintf(os.Stdout, "Received %s; stopping PostgreSQL\n", received)
	if err := database.Stop(); err != nil {
		return fmt.Errorf("stop embedded PostgreSQL: %w", err)
	}
	return nil
}

func postgresBinariesPath() (string, error) {
	if path := os.Getenv("POSTGRES_BINARIES"); path != "" {
		return path, nil
	}
	path, err := runfiles.Rlocation(postgresBinariesRunfile)
	if err != nil {
		return "", fmt.Errorf("locate embedded PostgreSQL binaries: %w", err)
	}
	return path, nil
}

func allowRemoteConnections(binariesPath, dataPath string) error {
	changed, err := ensureRemoteHostEntry(dataPath)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	command := exec.Command(filepath.Join(binariesPath, "bin", "pg_ctl"), "reload", "-D", dataPath)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reload PostgreSQL access rules: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func ensureRemoteHostEntry(dataPath string) (bool, error) {
	hbaPath := filepath.Join(dataPath, "pg_hba.conf")
	contents, err := os.ReadFile(hbaPath)
	if err != nil {
		return false, fmt.Errorf("read %s: %w", hbaPath, err)
	}
	if strings.Contains(string(contents), remoteHostEntry) {
		return false, nil
	}

	file, err := os.OpenFile(hbaPath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return false, fmt.Errorf("open %s: %w", hbaPath, err)
	}
	prefix := ""
	if len(contents) > 0 && contents[len(contents)-1] != '\n' {
		prefix = "\n"
	}
	if _, err := file.WriteString(prefix + remoteHostEntry + "\n"); err != nil {
		_ = file.Close()
		return false, fmt.Errorf("append to %s: %w", hbaPath, err)
	}
	if err := file.Close(); err != nil {
		return false, fmt.Errorf("close %s: %w", hbaPath, err)
	}
	return true, nil
}

func healthcheck() error {
	port, err := postgresPort()
	if err != nil {
		return err
	}
	connection, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.FormatUint(uint64(port), 10)), time.Second)
	if err != nil {
		return fmt.Errorf("PostgreSQL is not ready: %w", err)
	}
	return connection.Close()
}

func postgresPort() (uint32, error) {
	port := environmentOrDefault("POSTGRES_PORT", "5432")
	parsed, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid POSTGRES_PORT %q: %w", port, err)
	}
	return uint32(parsed), nil
}

func environmentOrDefault(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
