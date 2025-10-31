package typesense

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CurationSetsInterface is a type for Curation Sets API operations
type CurationSetsInterface interface {
	// Create or update a curation set
	Upsert(ctx context.Context, curationSetName string, curationSetSchema *api.CurationSetCreateSchema) (*api.CurationSetSchema, error)
	// Retrieve all curation sets
	Retrieve(ctx context.Context) ([]api.CurationSetSchema, error)
}

// curationSets is internal implementation of CurationSetsInterface
type curationSets struct {
	apiClient APIClientInterface
}

func (c *curationSets) Upsert(ctx context.Context, curationSetName string, curationSetSchema *api.CurationSetCreateSchema) (*api.CurationSetSchema, error) {
	response, err := c.apiClient.UpsertCurationSetWithResponse(ctx, curationSetName, api.UpsertCurationSetJSONRequestBody(*curationSetSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *curationSets) Retrieve(ctx context.Context) ([]api.CurationSetSchema, error) {
	response, err := c.apiClient.RetrieveCurationSetsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	// The API returns an array directly as specified in the OpenAPI spec
	var curationSets []api.CurationSetSchema
	if err := json.Unmarshal(response.Body, &curationSets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal curation sets response: %w", err)
	}

	return curationSets, nil
}
