package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// SynonymSetInterface is a type for individual Synonym Set API operations
type SynonymSetInterface interface {
	// Retrieve a single synonym set
	Retrieve(ctx context.Context) (*api.SynonymSetRetrieveSchema, error)
	// Update a synonym set
	Upsert(ctx context.Context, synonymSetSchema *api.SynonymSetCreateSchema) (*api.SynonymSetSchema, error)
	// Delete a synonym set
	Delete(ctx context.Context) (*api.SynonymSetDeleteSchema, error)
}

// synonymSet is internal implementation of SynonymSetInterface
type synonymSet struct {
	apiClient      APIClientInterface
	synonymSetName string
}

func (s *synonymSet) Retrieve(ctx context.Context) (*api.SynonymSetRetrieveSchema, error) {
	response, err := s.apiClient.RetrieveSynonymSetWithResponse(ctx, s.synonymSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (s *synonymSet) Upsert(ctx context.Context, synonymSetSchema *api.SynonymSetCreateSchema) (*api.SynonymSetSchema, error) {
	response, err := s.apiClient.UpsertSynonymSetWithResponse(ctx, s.synonymSetName, api.UpsertSynonymSetJSONRequestBody(*synonymSetSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (s *synonymSet) Delete(ctx context.Context) (*api.SynonymSetDeleteSchema, error) {
	response, err := s.apiClient.DeleteSynonymSetWithResponse(ctx, s.synonymSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
