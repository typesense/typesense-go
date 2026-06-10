package typesense

import (
	"context"
)

type MetricsInterface interface {
	// Get current RAM, CPU, Disk & Network usage metrics.
	//
	// Retrieve the metrics.
	//
	// HTTP: GET /metrics.json
	//
	// See: https://typesense.org/docs/latest/api/cluster-operations.html
	Retrieve(ctx context.Context) (map[string]interface{}, error)
}

type metrics struct {
	apiClient APIClientInterface
}

// Get current RAM, CPU, Disk & Network usage metrics.
//
// Retrieve the metrics.
//
// HTTP: GET /metrics.json
//
// See: https://typesense.org/docs/latest/api/cluster-operations.html
func (m *metrics) Retrieve(ctx context.Context) (map[string]interface{}, error) {
	response, err := m.apiClient.RetrieveMetricsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}
