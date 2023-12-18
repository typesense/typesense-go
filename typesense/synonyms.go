package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// SynonymsInterface is a type for Search Synonyms API operations
type SynonymsInterface interface {
	// Create or update a synonym
	Upsert(ctx context.Context, synonymID string, synonymSchema *api.SearchSynonymSchema) (*api.SearchSynonym, error)
	// List all collection synonyms
	Retrieve(ctx context.Context) ([]*api.SearchSynonym, error)
}

// synonyms is internal implementation of SynonymsInterface
type synonyms struct {
	apiClient      APIClientInterface
	collectionName string
}

func (s *synonyms) Upsert(ctx context.Context, synonymID string, synonymSchema *api.SearchSynonymSchema) (*api.SearchSynonym, error) {
	response, err := s.apiClient.UpsertSearchSynonymWithResponse(ctx,
		s.collectionName, synonymID, api.UpsertSearchSynonymJSONRequestBody(*synonymSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (s *synonyms) Retrieve(ctx context.Context) ([]*api.SearchSynonym, error) {
	response, err := s.apiClient.GetSearchSynonymsWithResponse(ctx, s.collectionName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Synonyms, nil
}
