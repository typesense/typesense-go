package typesense

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// SynonymSetsInterface is a type for Synonym Sets API operations
type SynonymSetsInterface interface {
	// Create or update a synonym set
	Upsert(ctx context.Context, synonymSetName string, synonymSetSchema *api.SynonymSetCreateSchema) (*api.SynonymSetSchema, error)
	// Retrieve all synonym sets
	Retrieve(ctx context.Context) ([]api.SynonymSetSchema, error)
}

// synonymSets is internal implementation of SynonymSetsInterface
type synonymSets struct {
	apiClient APIClientInterface
}

func (s *synonymSets) Upsert(ctx context.Context, synonymSetName string, synonymSetSchema *api.SynonymSetCreateSchema) (*api.SynonymSetSchema, error) {
	response, err := s.apiClient.UpsertSynonymSetWithResponse(ctx, synonymSetName, api.UpsertSynonymSetJSONRequestBody(*synonymSetSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (s *synonymSets) Retrieve(ctx context.Context) ([]api.SynonymSetSchema, error) {
	response, err := s.apiClient.RetrieveSynonymSetsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	// The API returns an array directly as specified in the OpenAPI spec
	var synonymSets []api.SynonymSetSchema
	if err := json.Unmarshal(response.Body, &synonymSets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal synonym sets response: %w", err)
	}

	return synonymSets, nil
}
