package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StatsInterface interface {
	Retrieve(ctx context.Context) (*api.APIStatsResponse, error)
}

type stats struct {
	apiClient APIClientInterface
}

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
