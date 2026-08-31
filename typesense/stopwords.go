package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StopwordsInterface interface {
	// Retrieves all stopwords sets.
	//
	// Retrieve the details of all stopwords sets
	//
	// HTTP: GET /stopwords
	//
	// See: https://typesense.org/docs/latest/api/stopwords.html
	Retrieve(ctx context.Context) ([]api.StopwordsSetSchema, error)
	// Upserts a stopwords set.
	//
	// When an analytics rule is created, we give it a name and describe the type, the source collections and the destination collection.
	//
	// HTTP: PUT /stopwords/{setId}
	//
	// See: https://typesense.org/docs/latest/api/stopwords.html
	Upsert(ctx context.Context, stopwordsSetId string, stopwordsSetUpsertSchema *api.StopwordsSetUpsertSchema) (*api.StopwordsSetSchema, error)
}

type stopwords struct {
	apiClient APIClientInterface
}

// Retrieves all stopwords sets.
//
// # Retrieve the details of all stopwords sets
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
// When an analytics rule is created, we give it a name and describe the type, the source collections and the destination collection.
//
// HTTP: PUT /stopwords/{setId}
//
// See: https://typesense.org/docs/latest/api/stopwords.html
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
