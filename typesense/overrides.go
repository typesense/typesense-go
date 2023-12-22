package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// OverridesInterface is a type for Search Overrides API operations
type OverridesInterface interface {
	Upsert(ctx context.Context, overrideID string, overrideSchema *api.SearchOverrideSchema) (*api.SearchOverride, error)
	Retrieve(ctx context.Context) ([]*api.SearchOverride, error)
}

// overrides is internal implementation of OverridesInterface
type overrides struct {
	apiClient      APIClientInterface
	collectionName string
}

func (o *overrides) Upsert(ctx context.Context, overrideID string, overrideSchema *api.SearchOverrideSchema) (*api.SearchOverride, error) {
	response, err := o.apiClient.UpsertSearchOverrideWithResponse(ctx,
		o.collectionName, overrideID, api.UpsertSearchOverrideJSONRequestBody(*overrideSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (o *overrides) Retrieve(ctx context.Context) ([]*api.SearchOverride, error) {
	response, err := o.apiClient.GetSearchOverridesWithResponse(ctx, o.collectionName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Overrides, nil
}
