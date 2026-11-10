package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CurationSetItemInterface is a type for individual Curation Set Item API operations
type CurationSetItemInterface interface {
	// Retrieve a curation set item.
	//
	// Retrieve a specific curation item by its id.
	//
	// HTTP: GET /curation_sets/{curationSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Retrieve(ctx context.Context) (*api.CurationItemSchema, error)
	// Create or update a curation set item.
	//
	// Create or update a curation set item with the given id.
	//
	// HTTP: PUT /curation_sets/{curationSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Upsert(ctx context.Context, item *api.CurationItemCreateSchema) (*api.CurationItemSchema, error)
	// Delete a curation set item.
	//
	// Delete a specific curation item by its id.
	//
	// HTTP: DELETE /curation_sets/{curationSetName}/items/{itemId}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Delete(ctx context.Context) (*api.CurationItemDeleteSchema, error)
}

// curationSetItem is internal implementation of CurationSetItemInterface
type curationSetItem struct {
	apiClient       APIClientInterface
	curationSetName string
	itemID          string
}

// Retrieve a curation set item.
//
// Retrieve a specific curation item by its id.
//
// HTTP: GET /curation_sets/{curationSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/curation.html
func (c *curationSetItem) Retrieve(ctx context.Context) (*api.CurationItemSchema, error) {
	response, err := c.apiClient.RetrieveCurationSetItemWithResponse(ctx, c.curationSetName, c.itemID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Create or update a curation set item.
//
// Create or update a curation set item with the given id.
//
// HTTP: PUT /curation_sets/{curationSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/curation.html
func (c *curationSetItem) Upsert(ctx context.Context, item *api.CurationItemCreateSchema) (*api.CurationItemSchema, error) {
	response, err := c.apiClient.UpsertCurationSetItemWithResponse(ctx, c.curationSetName, c.itemID, api.UpsertCurationSetItemJSONRequestBody(*item))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Delete a curation set item.
//
// Delete a specific curation item by its id.
//
// HTTP: DELETE /curation_sets/{curationSetName}/items/{itemId}
//
// See: https://typesense.org/docs/latest/api/curation.html
func (c *curationSetItem) Delete(ctx context.Context) (*api.CurationItemDeleteSchema, error) {
	response, err := c.apiClient.DeleteCurationSetItemWithResponse(ctx, c.curationSetName, c.itemID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
