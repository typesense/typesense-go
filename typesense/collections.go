package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CollectionsInterface is a type for Collections API operations
type CollectionsInterface interface {
	// Create a new collection.
	//
	// When a collection is created, we give it a name and describe the fields that will be indexed from the documents added to the collection.
	//
	// HTTP: POST /collections
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Create(ctx context.Context, schema *api.CollectionSchema) (*api.CollectionResponse, error)
	// List all collections.
	//
	// Returns a summary of all your collections. The collections are returned sorted by creation date, with the most recent collections appearing first.
	//
	// HTTP: GET /collections
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Retrieve(ctx context.Context, params *api.GetCollectionsParams) ([]*api.CollectionResponse, error)
}

// collections is internal implementation of CollectionsInterface
type collections struct {
	apiClient APIClientInterface
}

// Create a new collection.
//
// When a collection is created, we give it a name and describe the fields that will be indexed from the documents added to the collection.
//
// HTTP: POST /collections
//
// See: https://typesense.org/docs/latest/api/collections.html
func (c *collections) Create(ctx context.Context, schema *api.CollectionSchema) (*api.CollectionResponse, error) {
	response, err := c.apiClient.CreateCollectionWithResponse(ctx,
		api.CreateCollectionJSONRequestBody(*schema))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}

// List all collections.
//
// Returns a summary of all your collections. The collections are returned sorted by creation date, with the most recent collections appearing first.
//
// HTTP: GET /collections
//
// See: https://typesense.org/docs/latest/api/collections.html
func (c *collections) Retrieve(ctx context.Context, params *api.GetCollectionsParams) ([]*api.CollectionResponse, error) {
	response, err := c.apiClient.GetCollectionsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}
