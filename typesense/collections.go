package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CollectionsInterface is a type for Collections API operations
type CollectionsInterface interface {
	Create(ctx context.Context, schema *api.CollectionSchema) (*api.CollectionResponse, error)
	Retrieve(ctx context.Context, params *api.GetCollectionsParams) ([]*api.CollectionResponse, error)
}

// collections is internal implementation of CollectionsInterface
type collections struct {
	apiClient APIClientInterface
}

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
