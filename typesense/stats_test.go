package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestStatsRetrieve(t *testing.T) {
	expectedData := &api.APIStatsResponse{
		DeleteLatencyMs:         pointer.Float64(10.5),
		DeleteRequestsPerSecond: pointer.Float64(5.0),
		ImportLatencyMs:         pointer.Float64(3.7142857142857144),
		ImportRequestsPerSecond: pointer.Float64(9.5),
		LatencyMs: &map[string]float64{
			"GET /stats.json": *pointer.Float64(32.5),
		},
		OverloadedRequestsPerSecond: pointer.Float64(9.5),
		PendingWriteBatches:         pointer.Float64(9.5),
		RequestsPerSecond: &map[string]float64{
			"GET /stats.json": *pointer.Float64(0.1111111111111111),
		},
		SearchLatencyMs:         pointer.Float64(9.5),
		SearchRequestsPerSecond: pointer.Float64(9.5),
		TotalRequestsPerSecond:  pointer.Float64(3.3),
		WriteLatencyMs:          pointer.Float64(9.5),
		WriteRequestsPerSecond:  pointer.Float64(9.5),
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
