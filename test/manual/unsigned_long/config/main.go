package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

func resolveRunfile(path string) (string, error) {
	if strings.HasPrefix(path, "/") {
		return path, nil
	}
	resolvedPath, err := runfiles.Rlocation(path)
	if err != nil {
		return "", fmt.Errorf("resolve runfile %s: %w", path, err)
	}
	return resolvedPath, nil
}

func run() error {
	templatePath := flag.String("template", "", "path to the portal configuration template")
	outputPath := flag.String("output", "", "path for the generated portal configuration")
	postgresPort := flag.Int("postgres-port", 0, "embedded Postgres port")
	httpPort := flag.Int("http-port", 0, "bb-portal HTTP port")
	flag.Parse()

	if *templatePath == "" {
		return errors.New("--template is required")
	}
	if *outputPath == "" {
		return errors.New("--output is required")
	}
	if *postgresPort == 0 {
		return errors.New("--postgres-port is required")
	}
	if *httpPort == 0 {
		return errors.New("--http-port is required")
	}

	resolvedTemplatePath, err := resolveRunfile(*templatePath)
	if err != nil {
		return err
	}
	template, err := os.ReadFile(resolvedTemplatePath)
	if err != nil {
		return fmt.Errorf("read portal configuration template: %w", err)
	}

	configuration := string(template)
	replacements := map[string]string{
		"@@POSTGRES_PORT@@": strconv.Itoa(*postgresPort),
		"@@HTTP_PORT@@":     strconv.Itoa(*httpPort),
	}
	for placeholder, value := range replacements {
		if !strings.Contains(configuration, placeholder) {
			return fmt.Errorf("portal configuration template does not contain %s", placeholder)
		}
		configuration = strings.ReplaceAll(configuration, placeholder, value)
	}
	if strings.Contains(configuration, "@@") {
		return errors.New("portal configuration contains an unresolved placeholder")
	}

	if err := os.WriteFile(*outputPath, []byte(configuration), 0o600); err != nil {
		return fmt.Errorf("write portal configuration: %w", err)
	}
	log.Printf("Generated bb-portal configuration at %s", *outputPath)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
