package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsRetrieve(t *testing.T) {
	expectedData := map[string]interface{}{
		"system_cpu1_active_percentage":        "0.00",
		"system_cpu2_active_percentage":        "0.00",
		"system_cpu3_active_percentage":        "0.00",
		"system_cpu4_active_percentage":        "0.00",
		"system_cpu_active_percentage":         "0.00",
		"system_disk_total_bytes":              "1043447808",
		"system_disk_used_bytes":               "561152",
		"system_memory_total_bytes":            "2086899712",
		"system_memory_used_bytes":             "1004507136",
		"system_memory_total_swap_bytes":       "1004507136",
		"system_memory_used_swap_bytes":        "0.00",
		"system_network_received_bytes":        "1466",
		"system_network_sent_bytes":            "182",
		"typesense_memory_active_bytes":        "29630464",
		"typesense_memory_allocated_bytes":     "27886840",
		"typesense_memory_fragmentation_ratio": "0.06",
		"typesense_memory_mapped_bytes":        "69701632",
		"typesense_memory_metadata_bytes":      "4588768",
		"typesense_memory_resident_bytes":      "29630464",
		"typesense_memory_retained_bytes":      "25718784",
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/metrics.json", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Metrics().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestMetricsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/metrics.json", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Metrics().Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}
