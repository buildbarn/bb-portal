package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// DeleteBuildsBeforeResponse represents the response structure for the deleteBuildsBefore mutation
type DeleteBuildsBeforeResponse struct {
	Data struct {
		DeleteBuildsBefore struct {
			Found      int  `json:"found"`
			Deleted    int  `json:"deleted"`
			Successful bool `json:"successful"`
		} `json:"deleteBuildsBefore"`
	} `json:"data"`
}

// DeleteInvocationsBeforeResponse represents the response structure for the deleteInvocationsBefore mutation
type DeleteInvocationsBeforeResponse struct {
	Data struct {
		DeleteInvocationsBefore struct {
			Found      int  `json:"found"`
			Deleted    int  `json:"deleted"`
			Successful bool `json:"successful"`
		} `json:"deleteInvocationsBefore"`
	} `json:"data"`
}

func main() {
	baseURL := getEnv("BASE_URL", "http://localhost:8081")
	hoursToSubtract := getEnvAsInt("HOURS_TO_SUBTRACT", 1000)
	timezone := getEnv("TIMEZONE", "UTC")
	sslCertFile := getEnv("SSL_CERT_FILE", "/etc/ssl/certs/ca-certificates.crt")

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		return
	}

	now := time.Now().In(loc)
	deleteBeforeTime := now.Add(-time.Duration(hoursToSubtract) * time.Hour)
	formattedTime := deleteBeforeTime.Format("2006-01-02T15:04:05.000Z")

	fmt.Printf("Deleting builds and invocations starting before %s\n", formattedTime)

	apiURL := fmt.Sprintf("%s/graphql", baseURL)

	// First query: Delete builds before the specified time
	queryDeleteBuilds := fmt.Sprintf(`
    mutation deleteBuildsBefore {
        deleteBuildsBefore(time: "%s") {
            found
            deleted
            successful
        }
    }`, formattedTime)

	executeGraphQLQuery(apiURL, queryDeleteBuilds, "DeleteBuildsBeforeResponse", sslCertFile)

	// Second query: Delete invocations before the specified time
	queryDeleteInvocations := fmt.Sprintf(`
    mutation deleteInvocationsBefore {
        deleteInvocationsBefore(time: "%s") {
            found
            deleted
            successful
        }
    }`, formattedTime)

	executeGraphQLQuery(apiURL, queryDeleteInvocations, "DeleteInvocationsBeforeResponse", sslCertFile)
}

func executeGraphQLQuery(apiURL, query, responseType, sslCertFile string) {
	// Determine if the URL is HTTP or HTTPS
	var client *http.Client
	if strings.HasPrefix(apiURL, "https://") {
		// Load custom certificate if sslCertFile is provided
		tlsConfig := &tls.Config{}
		if _, err := os.Stat(sslCertFile); err == nil {
			certPool := x509.NewCertPool()
			certData, err := ioutil.ReadFile(sslCertFile)
			if err != nil {
				fmt.Printf("Error reading SSL certificate file: %v\n", err)
				return
			}
			if !certPool.AppendCertsFromPEM(certData) {
				fmt.Println("Failed to append certificates from SSL certificate file")
				return
			}
			tlsConfig.RootCAs = certPool
			fmt.Printf("Loaded SSL certificate file: %s\n", sslCertFile)
		} else {
			fmt.Println("SSL certificate file not found, using system certificates")
		}

		// Create HTTP client with custom TLS configuration
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	} else {
		// Create HTTP client for HTTP protocol
		client = &http.Client{}
	}

	// Prepare the GraphQL request
	reqBody := map[string]string{"query": query}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshaling request body: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Parse the response
	if responseType == "DeleteBuildsBeforeResponse" {
		var response DeleteBuildsBeforeResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			return
		}
		responseJSON, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println("Delete Builds Response:")
		fmt.Println(string(responseJSON))
	} else if responseType == "DeleteInvocationsBeforeResponse" {
		var response DeleteInvocationsBeforeResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			return
		}
		responseJSON, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println("Delete Invocations Response:")
		fmt.Println(string(responseJSON))
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		var value int
		if _, err := fmt.Sscanf(valueStr, "%d", &value); err == nil {
			return value
		}
	}
	return defaultValue
}
