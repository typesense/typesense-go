package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StopwordsInterface interface {
	Retrieve(ctx context.Context) ([]api.StopwordsSetSchema, error)
	Upsert(ctx context.Context, stopwordsSetId string, stopwordsSetUpsertSchema *api.StopwordsSetUpsertSchema) (*api.StopwordsSetSchema, error)
}

type stopwords struct {
	apiClient APIClientInterface
}

func (s *stopwords) Retrieve(ctx context.Context) ([]api.StopwordsSetSchema, error) {
	response, err := s.apiClient.RetrieveStopwordsSetsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Stopwords, nil
}

func (s *stopwords) Upsert(ctx context.Context, stopwordsSetId string, stopwordsSetUpsertSchema *api.StopwordsSetUpsertSchema) (*api.StopwordsSetSchema, error) {
	response, err := s.apiClient.UpsertStopwordsSetWithResponse(ctx, stopwordsSetId, *stopwordsSetUpsertSchema)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
