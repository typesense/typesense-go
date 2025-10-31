package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CurationSetInterface is a type for individual Curation Set API operations
type CurationSetInterface interface {
	// Retrieve a single curation set
	Retrieve(ctx context.Context) (*api.CurationSetRetrieveSchema, error)
	// Update a curation set
	Upsert(ctx context.Context, curationSetSchema *api.CurationSetCreateSchema) (*api.CurationSetSchema, error)
	// Delete a curation set
	Delete(ctx context.Context) (*api.CurationSetDeleteSchema, error)
}

// curationSet is internal implementation of CurationSetInterface
type curationSet struct {
	apiClient       APIClientInterface
	curationSetName string
}

func (c *curationSet) Retrieve(ctx context.Context) (*api.CurationSetRetrieveSchema, error) {
	response, err := c.apiClient.RetrieveCurationSetWithResponse(ctx, c.curationSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *curationSet) Upsert(ctx context.Context, curationSetSchema *api.CurationSetCreateSchema) (*api.CurationSetSchema, error) {
	response, err := c.apiClient.UpsertCurationSetWithResponse(ctx, c.curationSetName, api.UpsertCurationSetJSONRequestBody(*curationSetSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *curationSet) Delete(ctx context.Context) (*api.CurationSetDeleteSchema, error) {
	response, err := c.apiClient.DeleteCurationSetWithResponse(ctx, c.curationSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
