package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	// Mock environment variables
	os.Setenv("BASE_URL", "http://localhost:8081")
	os.Setenv("HOURS_TO_SUBTRACT", "120")
	os.Setenv("TIMEZONE", "UTC")
	os.Setenv("SSL_CERT_FILE", "/path/to/mock/cert.pem")

	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/graphql", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Mock response for builds
		if r.Body != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
                "data": {
                    "deleteBuildsBefore": {
                        "found": 10,
                        "deleted": 8,
                        "successful": true
                    }
                }
            }`))
		}
	}))
	defer mockServer.Close()

	// Override BASE_URL with the mock server URL
	os.Setenv("BASE_URL", mockServer.URL)

	// Run the main function
	main()

	// Reset environment variables
	os.Unsetenv("BASE_URL")
	os.Unsetenv("HOURS_TO_SUBTRACT")
	os.Unsetenv("TIMEZONE")
	os.Unsetenv("SSL_CERT_FILE")
}

func TestExecuteGraphQLQuery(t *testing.T) {
	// Mock environment variables
	os.Setenv("SSL_CERT_FILE", "/path/to/mock/cert.pem")

	// Mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/graphql", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Mock response for invocations
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "data": {
                "deleteInvocationsBefore": {
                    "found": 5,
                    "deleted": 4,
                    "successful": true
                }
            }
        }`))
	}))
	defer mockServer.Close()

	// Test the executeGraphQLQuery function
	apiURL := mockServer.URL + "/graphql"
	query := `
        mutation deleteInvocationsBefore {
            deleteInvocationsBefore(time: "2025-05-01T00:00:00.000Z") {
                found
                deleted
                successful
            }
        }`
	sslCertFile := "/path/to/mock/cert.pem" // Add the SSL certificate file argument
	executeGraphQLQuery(apiURL, query, "DeleteInvocationsBeforeResponse", sslCertFile)

	// Reset environment variables
	os.Unsetenv("SSL_CERT_FILE")
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	assert.Equal(t, "test_value", getEnv("TEST_KEY", "default_value"))
	assert.Equal(t, "default_value", getEnv("NON_EXISTENT_KEY", "default_value"))
	os.Unsetenv("TEST_KEY")
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT_KEY", "42")
	assert.Equal(t, 42, getEnvAsInt("TEST_INT_KEY", 0))
	assert.Equal(t, 0, getEnvAsInt("NON_EXISTENT_INT_KEY", 0))
	os.Unsetenv("TEST_INT_KEY")
}
