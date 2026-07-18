package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StatsInterface interface {
	// Get stats about API endpoints.
	//
	// Retrieve the stats about API endpoints.
	//
	// HTTP: GET /stats.json
	//
	// See: https://typesense.org/docs/latest/api/cluster-operations.html
	Retrieve(ctx context.Context) (*api.APIStatsResponse, error)
}

type stats struct {
	apiClient APIClientInterface
}

// Get stats about API endpoints.
//
// Retrieve the stats about API endpoints.
//
// HTTP: GET /stats.json
//
// See: https://typesense.org/docs/latest/api/cluster-operations.html
func (s *stats) Retrieve(ctx context.Context) (*api.APIStatsResponse, error) {
	response, err := s.apiClient.RetrieveAPIStatsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
