package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func TestStatsRetrieve(t *testing.T) {
	expectedData := &api.APIStatsResponse{
		DeleteLatencyMs:         pointer.Float32(10.5),
		DeleteRequestsPerSecond: pointer.Float32(5.0),
		ImportLatencyMs:         pointer.Float32(9.5),
		ImportRequestsPerSecond: pointer.Float32(9.5),
		LatencyMs: &map[string]interface{}{
			"GET /stats.json": 0.0,
		},
		OverloadedRequestsPerSecond: pointer.Float32(9.5),
		PendingWriteBatches:         pointer.Float32(9.5),
		RequestsPerSecond: &map[string]interface{}{
			"GET /stats.json": pointer.Float32(3.3),
		},
		SearchLatencyMs:         pointer.Float32(9.5),
		SearchRequestsPerSecond: pointer.Float32(9.5),
		TotalRequestsPerSecond:  pointer.Float32(3.3),
		WriteLatencyMs:          pointer.Float32(9.5),
		WriteRequestsPerSecond:  pointer.Float32(9.5),
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/stats.json", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Stats().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestStatsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/stats.json", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Stats().Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}
