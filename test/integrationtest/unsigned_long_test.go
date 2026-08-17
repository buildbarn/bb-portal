package integrationtest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/testkit"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGraphQLUnsignedLongValuesAboveJavaScriptSafeInteger(t *testing.T) {
	const (
		firstUnsafeInteger    uint64 = 1<<53 + 1
		maxSignedLong         uint64 = 1<<63 - 1
		firstUnsignedLongOnly uint64 = 1 << 63
		maxUnsignedLong       uint64 = 1<<64 - 1
	)

	actionCacheValues := struct {
		sizeInBytes                     uint64
		saveTimeInMs                    uint64
		loadTimeInMs                    uint64
		cacheCheckSemaphoreWaitTimeInMs uint64
	}{
		sizeInBytes:                     firstUnsafeInteger,
		saveTimeInMs:                    firstUnsafeInteger + 2,
		loadTimeInMs:                    maxSignedLong,
		cacheCheckSemaphoreWaitTimeInMs: firstUnsignedLongOnly,
	}
	networkValues := struct {
		bytesSent             uint64
		bytesRecv             uint64
		packetsSent           uint64
		packetsRecv           uint64
		peakBytesSentPerSec   uint64
		peakBytesRecvPerSec   uint64
		peakPacketsSentPerSec uint64
		peakPacketsRecvPerSec uint64
	}{
		bytesSent:             firstUnsafeInteger + 4,
		bytesRecv:             firstUnsafeInteger + 6,
		packetsSent:           maxSignedLong - 2,
		packetsRecv:           maxSignedLong,
		peakBytesSentPerSec:   firstUnsignedLongOnly,
		peakBytesRecvPerSec:   firstUnsignedLongOnly + 2,
		peakPacketsSentPerSec: maxUnsignedLong - 2,
		peakPacketsRecvPerSec: maxUnsignedLong,
	}

	ctx := context.Background()
	db := testutils.SetupTestDB(t, dbProvider)
	testCase := testCase{
		saveDataLevel: &bb_portal.BuildEventStreamService_SaveDataLevel{
			Level: &bb_portal.BuildEventStreamService_SaveDataLevel_BasicAndTarget{
				BasicAndTarget: &emptypb.Empty{},
			},
		},
	}
	bepUploader := setupTestBepUploader(t, db, testCase)
	bepFile, err := os.Open(bepFolderPath + "/" + abortedTests.filename)
	require.NoError(t, err)
	_, _, err = bepUploader.RecordEventNdjsonFile(ctx, bepFile)
	require.NoError(t, err)
	require.NoError(t, bepFile.Close())

	actionCacheStatistics, err := db.Ent().ActionCacheStatistics.Query().Only(ctx)
	require.NoError(t, err)
	actionSummary, err := actionCacheStatistics.QueryActionSummary().Only(ctx)
	require.NoError(t, err)
	require.NoError(t, db.Ent().ActionCacheStatistics.DeleteOne(actionCacheStatistics).Exec(ctx))
	_, err = db.Ent().ActionCacheStatistics.Create().
		SetActionSummaryID(actionSummary.ID).
		SetSizeInBytes(actionCacheValues.sizeInBytes).
		SetSaveTimeInMs(actionCacheValues.saveTimeInMs).
		SetLoadTimeInMs(actionCacheValues.loadTimeInMs).
		SetCacheCheckSemaphoreWaitTimeInMs(actionCacheValues.cacheCheckSemaphoreWaitTimeInMs).
		Save(ctx)
	require.NoError(t, err)

	systemNetworkStats, err := db.Ent().SystemNetworkStats.Query().Only(ctx)
	require.NoError(t, err)
	networkMetrics, err := systemNetworkStats.QueryNetworkMetrics().Only(ctx)
	require.NoError(t, err)
	require.NoError(t, db.Ent().SystemNetworkStats.DeleteOne(systemNetworkStats).Exec(ctx))
	_, err = db.Ent().SystemNetworkStats.Create().
		SetNetworkMetricsID(networkMetrics.ID).
		SetBytesSent(networkValues.bytesSent).
		SetBytesRecv(networkValues.bytesRecv).
		SetPacketsSent(networkValues.packetsSent).
		SetPacketsRecv(networkValues.packetsRecv).
		SetPeakBytesSentPerSec(networkValues.peakBytesSentPerSec).
		SetPeakBytesRecvPerSec(networkValues.peakBytesRecvPerSec).
		SetPeakPacketsSentPerSec(networkValues.peakPacketsSentPerSec).
		SetPeakPacketsRecvPerSec(networkValues.peakPacketsRecvPerSec).
		Save(ctx)
	require.NoError(t, err)

	graphqlServer := startGraphqlHTTPServer(t, db)
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": `
			query GetUnsignedLongValues($invocationID: UUID!) {
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
		"variables": map[string]interface{}{
			"invocationID": abortedTests.invocationID,
		},
	})
	require.NoError(t, err)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, graphqlServer.URL, bytes.NewReader(requestBody))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	response, err := graphqlServer.Client().Do(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	var graphqlResponse struct {
		Data   map[string]interface{}   `json:"data"`
		Errors []map[string]interface{} `json:"errors"`
	}
	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()
	require.NoError(t, decoder.Decode(&graphqlResponse))
	require.NoError(t, response.Body.Close())
	require.Empty(t, graphqlResponse.Errors)

	jsonNumber := func(value uint64) json.Number {
		return json.Number(strconv.FormatUint(value, 10))
	}
	want := map[string]interface{}{
		"getBazelInvocation": map[string]interface{}{
			"metrics": map[string]interface{}{
				"actionSummary": map[string]interface{}{
					"actionCacheStatistics": map[string]interface{}{
						"sizeInBytes":                     jsonNumber(actionCacheValues.sizeInBytes),
						"saveTimeInMs":                    jsonNumber(actionCacheValues.saveTimeInMs),
						"loadTimeInMs":                    jsonNumber(actionCacheValues.loadTimeInMs),
						"cacheCheckSemaphoreWaitTimeInMs": jsonNumber(actionCacheValues.cacheCheckSemaphoreWaitTimeInMs),
					},
				},
				"networkMetrics": map[string]interface{}{
					"systemNetworkStats": map[string]interface{}{
						"bytesSent":             jsonNumber(networkValues.bytesSent),
						"bytesRecv":             jsonNumber(networkValues.bytesRecv),
						"packetsSent":           jsonNumber(networkValues.packetsSent),
						"packetsRecv":           jsonNumber(networkValues.packetsRecv),
						"peakBytesSentPerSec":   jsonNumber(networkValues.peakBytesSentPerSec),
						"peakBytesRecvPerSec":   jsonNumber(networkValues.peakBytesRecvPerSec),
						"peakPacketsSentPerSec": jsonNumber(networkValues.peakPacketsSentPerSec),
						"peakPacketsRecvPerSec": jsonNumber(networkValues.peakPacketsRecvPerSec),
					},
				},
			},
		},
	}
	require.Equal(t, want, graphqlResponse.Data)
	testkit.CheckAgainstGoldenFile(t, graphqlResponse.Data, goldenFolderPath, t.Name(), updateGoldenFiles, nil)
}
