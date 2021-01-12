package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// OverrideInterface is a type for Search Override API operations
type OverrideInterface interface {
	Retrieve() (*api.SearchOverride, error)
	Delete() (*api.SearchOverride, error)
}

// override is internal implementation of OverrideInterface
type override struct {
	apiClient      APIClientInterface
	collectionName string
	overrideID     string
}

func (o *override) Retrieve() (*api.SearchOverride, error) {
	response, err := o.apiClient.GetSearchOverrideWithResponse(context.Background(),
		o.collectionName, o.overrideID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (o *override) Delete() (*api.SearchOverride, error) {
	response, err := o.apiClient.DeleteSearchOverrideWithResponse(context.Background(),
		o.collectionName, o.overrideID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
