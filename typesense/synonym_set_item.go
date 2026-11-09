package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// SynonymSetItemInterface is a type for individual Synonym Set Item API operations
type SynonymSetItemInterface interface {
	// Retrieve a synonym set item.
	//
	// Retrieve a specific synonym item by its id.
	//
	// HTTP: GET /synonym_sets/{synonymSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Retrieve(ctx context.Context) (*api.SynonymItemSchema, error)
	// Create or update a synonym set item.
	//
	// Create or update a synonym set item with the given id.
	//
	// HTTP: PUT /synonym_sets/{synonymSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Upsert(ctx context.Context, item *api.SynonymItemUpsertSchema) (*api.SynonymItemSchema, error)
	// Delete a synonym set item.
	//
	// Delete a specific synonym item by its id.
	//
	// HTTP: DELETE /synonym_sets/{synonymSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Delete(ctx context.Context) (*api.SynonymItemDeleteSchema, error)
}

// synonymSetItem is internal implementation of SynonymSetItemInterface
type synonymSetItem struct {
	apiClient      APIClientInterface
	synonymSetName string
	itemID         string
}

// Retrieve a synonym set item.
//
// Retrieve a specific synonym item by its id.
//
// HTTP: GET /synonym_sets/{synonymSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
func (s *synonymSetItem) Retrieve(ctx context.Context) (*api.SynonymItemSchema, error) {
	response, err := s.apiClient.RetrieveSynonymSetItemWithResponse(ctx, s.synonymSetName, s.itemID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Create or update a synonym set item.
//
// Create or update a synonym set item with the given id.
//
// HTTP: PUT /synonym_sets/{synonymSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
func (s *synonymSetItem) Upsert(ctx context.Context, item *api.SynonymItemUpsertSchema) (*api.SynonymItemSchema, error) {
	response, err := s.apiClient.UpsertSynonymSetItemWithResponse(ctx, s.synonymSetName, s.itemID, api.UpsertSynonymSetItemJSONRequestBody(*item))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Delete a synonym set item.
//
// Delete a specific synonym item by its id.
//
// HTTP: DELETE /synonym_sets/{synonymSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
func (s *synonymSetItem) Delete(ctx context.Context) (*api.SynonymItemDeleteSchema, error) {
	response, err := s.apiClient.DeleteSynonymSetItemWithResponse(ctx, s.synonymSetName, s.itemID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
