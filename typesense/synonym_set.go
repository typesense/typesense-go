package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// SynonymSetInterface is a type for individual Synonym Set API operations
type SynonymSetInterface interface {
	// Retrieve a synonym set.
	//
	// Retrieve a specific synonym set by its name
	//
	// HTTP: GET /synonym_sets/{synonymSetName}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Retrieve(ctx context.Context) (*api.SynonymSetSchema, error)
	// Create or update a synonym set.
	//
	// Create or update a synonym set with the given name
	//
	// HTTP: PUT /synonym_sets/{synonymSetName}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Upsert(ctx context.Context, synonymSetSchema *api.SynonymSetCreateSchema) (*api.SynonymSetSchema, error)
	// Delete a synonym set.
	//
	// Delete a specific synonym set by its name
	//
	// HTTP: DELETE /synonym_sets/{synonymSetName}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Delete(ctx context.Context) (*api.SynonymSetDeleteSchema, error)
}

// synonymSet is internal implementation of SynonymSetInterface
type synonymSet struct {
	apiClient      APIClientInterface
	synonymSetName string
}

// Retrieve a synonym set.
//
// # Retrieve a specific synonym set by its name
//
// HTTP: GET /synonym_sets/{synonymSetName}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
func (s *synonymSet) Retrieve(ctx context.Context) (*api.SynonymSetSchema, error) {
	response, err := s.apiClient.RetrieveSynonymSetWithResponse(ctx, s.synonymSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Create or update a synonym set.
//
// # Create or update a synonym set with the given name
//
// HTTP: PUT /synonym_sets/{synonymSetName}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
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

// Delete a synonym set.
//
// # Delete a specific synonym set by its name
//
// HTTP: DELETE /synonym_sets/{synonymSetName}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
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
