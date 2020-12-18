package typesense

import (
	"context"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

// client.Collection('name').Retrieve()
// client.Collection('name').Delete()

// CollectionInterface is a type for Collection API operations
type CollectionInterface interface {
	Retrieve() (*api.Collection, error)
	Delete() (*api.Collection, error)
	Documents() DocumentsInterface
}

// collection is internal implementation of CollectionInterface
type collection struct {
	apiClient api.ClientWithResponsesInterface
	name      string
}

func (c *collection) Retrieve() (*api.Collection, error) {
	response, err := c.apiClient.GetCollectionWithResponse(context.Background(), c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}

func (c *collection) Delete() (*api.Collection, error) {
	response, err := c.apiClient.DeleteCollectionWithResponse(context.Background(), c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}

func (c *collection) Documents() DocumentsInterface {
	return &documents{apiClient: c.apiClient, collectionName: c.name}
}
