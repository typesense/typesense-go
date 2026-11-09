package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CurationSetItemsInterface is a type for Curation Set Items API operations
type CurationSetItemsInterface interface {
	// List items in a curation set.
	//
	// Retrieve all curation items in a set.
	//
	// HTTP: GET /curation_sets/{curationSetName}/items
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Retrieve(ctx context.Context) ([]api.CurationItemSchema, error)
	// Create or update a curation set item.
	//
	// Create or update a curation set item with the given id.
	//
	// HTTP: PUT /curation_sets/{curationSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Upsert(ctx context.Context, itemID string, item *api.CurationItemCreateSchema) (*api.CurationItemSchema, error)
}

// curationSetItems is internal implementation of CurationSetItemsInterface
type curationSetItems struct {
	apiClient       APIClientInterface
	curationSetName string
}

// List items in a curation set.
//
// Retrieve all curation items in a set.
//
// HTTP: GET /curation_sets/{curationSetName}/items
//
// See: https://typesense.org/docs/latest/api/curation.html
func (c *curationSetItems) Retrieve(ctx context.Context) ([]api.CurationItemSchema, error) {
	response, err := c.apiClient.RetrieveCurationSetItemsWithResponse(ctx, c.curationSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}

// Create or update a curation set item.
//
// Create or update a curation set item with the given id.
//
// HTTP: PUT /curation_sets/{curationSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/curation.html
func (c *curationSetItems) Upsert(ctx context.Context, itemID string, item *api.CurationItemCreateSchema) (*api.CurationItemSchema, error) {
	response, err := c.apiClient.UpsertCurationSetItemWithResponse(ctx, c.curationSetName, itemID, api.UpsertCurationSetItemJSONRequestBody(*item))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
