package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StopwordsInterface interface {
	// Retrieves all stopwords sets.
	//
	// Retrieve the details of all stopwords sets.
	//
	// HTTP: GET /stopwords
	//
	// See: https://typesense.org/docs/latest/api/stopwords.html
	Retrieve(ctx context.Context) ([]api.StopwordsSetSchema, error)
	// Upserts a stopwords set.
	//
	// Create the stopwords set with the given id, or overwrite it if it exists, and return the stored set.
	//
	// HTTP: PUT /stopwords/{setId}
	//
	// See: https://typesense.org/docs/latest/api/stopwords.html
	//
	//docgen:keep
	Upsert(ctx context.Context, stopwordsSetId string, stopwordsSetUpsertSchema *api.StopwordsSetUpsertSchema) (*api.StopwordsSetSchema, error)
}

type stopwords struct {
	apiClient APIClientInterface
}

// Retrieves all stopwords sets.
//
// Retrieve the details of all stopwords sets.
//
// HTTP: GET /stopwords
//
// See: https://typesense.org/docs/latest/api/stopwords.html
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

// Upserts a stopwords set.
//
// Create the stopwords set with the given id, or overwrite it if it exists, and return the stored set.
//
// HTTP: PUT /stopwords/{setId}
//
// See: https://typesense.org/docs/latest/api/stopwords.html
//
//docgen:keep
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
