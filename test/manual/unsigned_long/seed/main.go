package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

const invocationID = "00000000-0000-4000-8000-000000000053"

var expectedUnsignedLongValues = map[string]string{
	"sizeInBytes":                     "9007199254740993",
	"saveTimeInMs":                    "9007199254740995",
	"loadTimeInMs":                    "9223372036854775807",
	"cacheCheckSemaphoreWaitTimeInMs": "9223372036854775808",
	"bytesSent":                       "9007199254740997",
	"bytesRecv":                       "9007199254740999",
	"packetsSent":                     "9223372036854775805",
	"packetsRecv":                     "9223372036854775807",
	"peakBytesSentPerSec":             "9223372036854775808",
	"peakBytesRecvPerSec":             "9223372036854775810",
	"peakPacketsSentPerSec":           "18446744073709551613",
	"peakPacketsRecvPerSec":           "18446744073709551615",
}

type graphqlResponse struct {
	Data struct {
		Invocation struct {
			Metrics struct {
				ActionSummary struct {
					ActionCacheStatistics map[string]json.Number `json:"actionCacheStatistics"`
				} `json:"actionSummary"`
				NetworkMetrics struct {
					SystemNetworkStats map[string]json.Number `json:"systemNetworkStats"`
				} `json:"networkMetrics"`
			} `json:"metrics"`
		} `json:"getBazelInvocation"`
	} `json:"data"`
	Errors []json.RawMessage `json:"errors"`
}

func uploadBEP(client *http.Client, portalURL, bepPath string) error {
	fixture, err := os.Open(bepPath)
	if err != nil {
		return fmt.Errorf("open BEP fixture: %w", err)
	}
	defer fixture.Close()

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	part, err := multipartWriter.CreateFormFile("file", filepath.Base(bepPath))
	if err != nil {
		return fmt.Errorf("create BEP multipart field: %w", err)
	}
	if _, err := io.Copy(part, fixture); err != nil {
		return fmt.Errorf("copy BEP fixture into request: %w", err)
	}
	if err := multipartWriter.Close(); err != nil {
		return fmt.Errorf("finish BEP multipart request: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, portalURL+"/api/v1/bep/upload", &requestBody)
	if err != nil {
		return fmt.Errorf("create BEP upload request: %w", err)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("upload BEP fixture: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read BEP upload response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("upload BEP fixture: status %s: %s", response.Status, responseBody)
	}

	var locationResponse struct {
		Location string
	}
	if err := json.Unmarshal(responseBody, &locationResponse); err != nil {
		return fmt.Errorf("decode BEP upload response: %w", err)
	}
	expectedLocation := "/bazel-invocations/" + invocationID
	if locationResponse.Location != expectedLocation {
		return fmt.Errorf("BEP upload location is %q, expected %q", locationResponse.Location, expectedLocation)
	}
	return nil
}

func queryUnsignedLongValues(client *http.Client, portalURL string) (*graphqlResponse, error) {
	requestBody, err := json.Marshal(map[string]any{
		"query": `
			query UnsignedLongManualValidation($invocationID: UUID!) {
				getBazelInvocation(invocationID: $invocationID) {
					metrics {
						actionSummary {
							actionCacheStatistics {
								sizeInBytes
								saveTimeInMs
								loadTimeInMs
								cacheCheckSemaphoreWaitTimeInMs
							}
						}
						networkMetrics {
							systemNetworkStats {
								bytesSent
								bytesRecv
								packetsSent
								packetsRecv
								peakBytesSentPerSec
								peakBytesRecvPerSec
								peakPacketsSentPerSec
								peakPacketsRecvPerSec
							}
						}
					}
				}
			}
		`,
		"variables": map[string]string{"invocationID": invocationID},
	})
	if err != nil {
		return nil, fmt.Errorf("encode GraphQL request: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, portalURL+"/graphql", bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("create GraphQL request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("query UnsignedLong values: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("query UnsignedLong values: status %s: %s", response.Status, responseBody)
	}

	var graphqlResponse graphqlResponse
	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&graphqlResponse); err != nil {
		return nil, fmt.Errorf("decode GraphQL response: %w", err)
	}
	if len(graphqlResponse.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL returned errors: %v", graphqlResponse.Errors)
	}
	return &graphqlResponse, nil
}

func validateUnsignedLongValues(response *graphqlResponse) error {
	actualValues := map[string]json.Number{}
	for key, value := range response.Data.Invocation.Metrics.ActionSummary.ActionCacheStatistics {
		actualValues[key] = value
	}
	for key, value := range response.Data.Invocation.Metrics.NetworkMetrics.SystemNetworkStats {
		actualValues[key] = value
	}
	for key, expected := range expectedUnsignedLongValues {
		actual, ok := actualValues[key]
		if !ok {
			return fmt.Errorf("GraphQL response is missing %s", key)
		}
		if actual.String() != expected {
			return fmt.Errorf("GraphQL response %s is %s, expected %s", key, actual, expected)
		}
	}
	return nil
}

func run() error {
	portalURL := flag.String("portal-url", "", "base URL of bb-portal")
	bepPath := flag.String("bep-file", "", "path to the UnsignedLong BEP fixture")
	flag.Parse()

	if *portalURL == "" {
		return errors.New("--portal-url is required")
	}
	if *bepPath == "" {
		return errors.New("--bep-file is required")
	}
	resolvedBEPPath, err := runfiles.Rlocation(*bepPath)
	if err != nil {
		return fmt.Errorf("resolve BEP fixture runfile: %w", err)
	}
	baseURL := strings.TrimRight(*portalURL, "/")
	client := &http.Client{Timeout: 30 * time.Second}
	if err := uploadBEP(client, baseURL, resolvedBEPPath); err != nil {
		return err
	}
	response, err := queryUnsignedLongValues(client, baseURL)
	if err != nil {
		return err
	}
	if err := validateUnsignedLongValues(response); err != nil {
		return err
	}

	log.Printf("UnsignedLong fixture loaded and verified")
	log.Printf("Open %s/bazel-invocations/%s/metrics", baseURL, invocationID)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
