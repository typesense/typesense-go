package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// CollectionsInterface is a type for Collections API operations
type CollectionsInterface interface {
	Create(schema *api.CollectionSchema) (*api.CollectionResponse, error)
	CreateCollectionFromStruct(structData interface{}) (*api.CollectionResponse, error)
	Retrieve() ([]*api.CollectionResponse, error)
}

// collections is internal implementation of CollectionsInterface
type collections struct {
	apiClient APIClientInterface
}

func (c *collections) Create(schema *api.CollectionSchema) (*api.CollectionResponse, error) {
	response, err := c.apiClient.CreateCollectionWithResponse(context.Background(),
		api.CreateCollectionJSONRequestBody(*schema))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}

func (c *collections) Retrieve() ([]*api.CollectionResponse, error) {
	response, err := c.apiClient.GetCollectionsWithResponse(context.Background())
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}

// CreateCollectionFromStruct creates a Typesense collection from a Go struct.
func (c *collections) CreateCollectionFromStruct(structData interface{}) (*api.CollectionResponse, error) {
	// Generate Typesense schema from the Go struct
	schema, err := CreateSchemaFromGoStruct(structData)
	if err != nil {
		return nil, err
	}

	// Use the generated schema to create a collection in Typesense
	response, err := c.apiClient.CreateCollectionWithResponse(context.Background(),
		api.CreateCollectionJSONRequestBody(*schema))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}
