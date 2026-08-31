package typesense

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// SynonymSetsInterface is a type for Synonym Sets API operations
type SynonymSetsInterface interface {
	// Create or update a synonym set.
	//
	// Create or update a synonym set with the given name
	//
	// HTTP: PUT /synonym_sets/{synonymSetName}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Upsert(ctx context.Context, synonymSetName string, synonymSetSchema *api.SynonymSetCreateSchema) (*api.SynonymSetSchema, error)
	// List all synonym sets.
	//
	// Retrieve all synonym sets
	//
	// HTTP: GET /synonym_sets
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Retrieve(ctx context.Context) ([]api.SynonymSetSchema, error)
}

// synonymSets is internal implementation of SynonymSetsInterface
type synonymSets struct {
	apiClient APIClientInterface
}

// Create or update a synonym set.
//
// # Create or update a synonym set with the given name
//
// HTTP: PUT /synonym_sets/{synonymSetName}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
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

// List all synonym sets.
//
// # Retrieve all synonym sets
//
// HTTP: GET /synonym_sets
//
// See: https://typesense.org/docs/latest/api/synonyms.html
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
