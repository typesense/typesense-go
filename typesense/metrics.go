package typesense

import (
	"context"
)

type MetricsInterface interface {
	Retrieve(ctx context.Context) (map[string]interface{}, error)
}

type metrics struct {
	apiClient APIClientInterface
}

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
