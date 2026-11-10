package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// SynonymSetItemsInterface is a type for Synonym Set Items API operations
type SynonymSetItemsInterface interface {
	// List items in a synonym set.
	//
	// Retrieve all synonym items in a set.
	//
	// HTTP: GET /synonym_sets/{synonymSetName}/items
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Retrieve(ctx context.Context) ([]api.SynonymItemSchema, error)
	// Create or update a synonym set item.
	//
	// Create or update a synonym set item with the given id.
	//
	// HTTP: PUT /synonym_sets/{synonymSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/synonyms.html
	Upsert(ctx context.Context, itemID string, item *api.SynonymItemUpsertSchema) (*api.SynonymItemSchema, error)
}

// synonymSetItems is internal implementation of SynonymSetItemsInterface
type synonymSetItems struct {
	apiClient      APIClientInterface
	synonymSetName string
}

// List items in a synonym set.
//
// Retrieve all synonym items in a set.
//
// HTTP: GET /synonym_sets/{synonymSetName}/items
//
// See: https://typesense.org/docs/latest/api/synonyms.html
func (s *synonymSetItems) Retrieve(ctx context.Context) ([]api.SynonymItemSchema, error) {
	response, err := s.apiClient.RetrieveSynonymSetItemsWithResponse(ctx, s.synonymSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}

// Create or update a synonym set item.
//
// Create or update a synonym set item with the given id.
//
// HTTP: PUT /synonym_sets/{synonymSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/synonyms.html
func (s *synonymSetItems) Upsert(ctx context.Context, itemID string, item *api.SynonymItemUpsertSchema) (*api.SynonymItemSchema, error) {
	response, err := s.apiClient.UpsertSynonymSetItemWithResponse(ctx, s.synonymSetName, itemID, api.UpsertSynonymSetItemJSONRequestBody(*item))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
