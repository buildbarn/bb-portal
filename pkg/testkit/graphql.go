package testkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

// errOpRegistered
var errOpRegistered = errors.New("operation is already registered")

// registryEntry struct
type registryEntry struct {
	document Document
	used     bool
}

// QueryRegistry struct
type QueryRegistry struct {
	entries map[string]*registryEntry
}

// LoadQueryRegistry function
func LoadQueryRegistry(t *testing.T, queryDir, consumerContractFile string) *QueryRegistry {
	reg := QueryRegistry{entries: make(map[string]*registryEntry)}

	if queryDir != "" {
		queryFiles, err := os.ReadDir(queryDir)
		if err != nil {
			panic(fmt.Sprintf("Could not read query directory %s", queryDir))
		}

		for _, file := range queryFiles {
			fileName := file.Name()
			if strings.HasSuffix(fileName, ".graphql") {
				fullPath := filepath.Join(queryDir, fileName)
				document, err := os.ReadFile(fullPath)
				require.NoError(t, err)
				opName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
				require.NoError(t, reg.Register(opName, Document(document)))
			}
		}
	}

	if consumerContractFile != "" {
		contract := mustLoadContract(consumerContractFile)

		for opName, op := range contract.Operations {
			require.NoError(t, reg.Register(opName, op.Document))
		}
	}
	return &reg
}

// Register function
func (reg QueryRegistry) Register(opName string, document Document) error {
	_, exists := reg.entries[opName]
	if exists {
		return fmt.Errorf("%w: %s", errOpRegistered, opName)
	}
	reg.entries[opName] = &registryEntry{
		document: document,
	}
	return nil
}

// NewRequest fuction
func (reg QueryRegistry) NewRequest(opName string) *graphql.Request {
	document := reg.MustGet(opName)
	return graphql.NewRequest(string(document))
}

// MustGet function
func (reg QueryRegistry) MustGet(opName string) Document {
	if entry, exists := reg.entries[opName]; exists {
		entry.used = true
		return entry.document
	}
	panic(fmt.Sprintf("operation %s is not regisitered", opName))
}

// UnusedOperations function
func (reg QueryRegistry) UnusedOperations() []string {
	unused := []string{}
	for name, entry := range reg.entries {
		if !entry.used {
			unused = append(unused, name)
		}
	}

	return unused
}

// ConsumerContract struct
type ConsumerContract struct {
	GeneratedAt time.Time            `yaml:"generatedAt"`
	Operations  map[string]Operation `yaml:"operations"`
}

type (
	// Document and Variables types
	Document string

	// Variables struct
	Variables map[string]interface{}
)

// Operation struct
type Operation struct {
	Document  Document `yaml:"document"`
	Signature string   `yaml:"signature"`
}

// mustLoadContract
func mustLoadContract(contractFile string) ConsumerContract {
	c := ConsumerContract{
		Operations: map[string]Operation{},
	}

	contract, err := os.ReadFile(contractFile)
	if err != nil {
		panic(err)
	}
	documents := map[string]string{}
	if json.Unmarshal(contract, &documents) != nil {
		panic(err)
	}

	for signature, document := range documents {
		query, err := parser.ParseQuery(&ast.Source{
			Name:  signature,
			Input: document,
		})
		if err != nil {
			panic(err)
		}
		c.Operations[query.Operations[0].Name] = Operation{
			Document:  Document(document),
			Signature: signature,
		}
	}
	return c
}
