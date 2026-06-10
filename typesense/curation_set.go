package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CurationSetInterface is a type for individual Curation Set API operations
type CurationSetInterface interface {
	// Retrieve a curation set.
	//
	// Retrieve a specific curation set by its name
	//
	// HTTP: GET /curation_sets/{curationSetName}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Retrieve(ctx context.Context) (*api.CurationSetSchema, error)
	// Create or update a curation set.
	//
	// Create or update a curation set with the given name
	//
	// HTTP: PUT /curation_sets/{curationSetName}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Upsert(ctx context.Context, curationSetSchema *api.CurationSetCreateSchema) (*api.CurationSetSchema, error)
	// Delete a curation set.
	//
	// Delete a specific curation set by its name
	//
	// HTTP: DELETE /curation_sets/{curationSetName}
	//
	// See: https://typesense.org/docs/latest/api/curation.html
	Delete(ctx context.Context) (*api.CurationSetDeleteSchema, error)
}

// curationSet is internal implementation of CurationSetInterface
type curationSet struct {
	apiClient       APIClientInterface
	curationSetName string
}

// Retrieve a curation set.
//
// # Retrieve a specific curation set by its name
//
// HTTP: GET /curation_sets/{curationSetName}
//
// See: https://typesense.org/docs/latest/api/curation.html
func (c *curationSet) Retrieve(ctx context.Context) (*api.CurationSetSchema, error) {
	response, err := c.apiClient.RetrieveCurationSetWithResponse(ctx, c.curationSetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Create or update a curation set.
//
// # Create or update a curation set with the given name
//
// HTTP: PUT /curation_sets/{curationSetName}
//
// See: https://typesense.org/docs/latest/api/curation.html
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

// Delete a curation set.
//
// # Delete a specific curation set by its name
//
// HTTP: DELETE /curation_sets/{curationSetName}
//
// See: https://typesense.org/docs/latest/api/curation.html
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
