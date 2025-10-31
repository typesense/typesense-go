package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StemmingDictionaryInterface interface {
	Retrieve(ctx context.Context) (*api.StemmingDictionary, error)
}

type stemmingDictionary struct {
	apiClient    APIClientInterface
	dictionaryId string
}

func (s *stemmingDictionary) Retrieve(ctx context.Context) (*api.StemmingDictionary, error) {
	response, err := s.apiClient.GetStemmingDictionaryWithResponse(ctx, s.dictionaryId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
