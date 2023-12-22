package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// SynonymInterface is a type for Search Synonym API operations
type SynonymInterface interface {
	// Retrieve a single search synonym
	Retrieve(ctx context.Context) (*api.SearchSynonym, error)
	// Delete a synonym associated with a collection
	Delete(ctx context.Context) (*api.SearchSynonym, error)
}

// synonym is internal implementation of SynonymInterface
type synonym struct {
	apiClient      APIClientInterface
	collectionName string
	synonymID      string
}

func (s *synonym) Retrieve(ctx context.Context) (*api.SearchSynonym, error) {
	response, err := s.apiClient.GetSearchSynonymWithResponse(ctx,
		s.collectionName, s.synonymID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (s *synonym) Delete(ctx context.Context) (*api.SearchSynonym, error) {
	response, err := s.apiClient.DeleteSearchSynonymWithResponse(ctx,
		s.collectionName, s.synonymID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
