package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StopwordInterface interface {
	Retrieve(ctx context.Context) (*api.StopwordsSetSchema, error)
	Delete(ctx context.Context) (*struct {
		Id string "json:\"id\""
	}, error)
}

type stopword struct {
	apiClient      APIClientInterface
	stopwordsSetId string
}

func (s *stopword) Retrieve(ctx context.Context) (*api.StopwordsSetSchema, error) {
	response, err := s.apiClient.RetrieveStopwordsSetWithResponse(ctx, s.stopwordsSetId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return &response.JSON200.Stopwords, nil
}

func (s *stopword) Delete(ctx context.Context) (*struct {
	Id string "json:\"id\""
}, error) {
	response, err := s.apiClient.DeleteStopwordsSetWithResponse(ctx, s.stopwordsSetId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
