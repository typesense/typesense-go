package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// CollectionInterface is a type for Collection API operations
type CollectionInterface interface {
	Retrieve(ctx context.Context) (*api.CollectionResponse, error)
	Delete(ctx context.Context) (*api.CollectionResponse, error)
	Documents() DocumentsInterface
	Document(documentID string) DocumentInterface
	Overrides() OverridesInterface
	Override(overrideID string) OverrideInterface
	Synonyms() SynonymsInterface
	Synonym(synonymID string) SynonymInterface
	Update(context.Context, *api.CollectionUpdateSchema) (*api.CollectionUpdateSchema, error)
}

// collection is internal implementation of CollectionInterface
type collection struct {
	apiClient APIClientInterface
	name      string
}

func (c *collection) Retrieve(ctx context.Context) (*api.CollectionResponse, error) {
	response, err := c.apiClient.GetCollectionWithResponse(ctx, c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *collection) Delete(ctx context.Context) (*api.CollectionResponse, error) {
	response, err := c.apiClient.DeleteCollectionWithResponse(ctx, c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *collection) Documents() DocumentsInterface {
	return &documents{apiClient: c.apiClient, collectionName: c.name}
}

func (c *collection) Document(documentID string) DocumentInterface {
	return &document{apiClient: c.apiClient, collectionName: c.name, documentID: documentID}
}

func (c *collection) Overrides() OverridesInterface {
	return &overrides{apiClient: c.apiClient, collectionName: c.name}
}

func (c *collection) Override(overrideID string) OverrideInterface {
	return &override{apiClient: c.apiClient, collectionName: c.name, overrideID: overrideID}
}

func (c *collection) Synonyms() SynonymsInterface {
	return &synonyms{apiClient: c.apiClient, collectionName: c.name}
}

func (c *collection) Synonym(synonymID string) SynonymInterface {
	return &synonym{apiClient: c.apiClient, collectionName: c.name, synonymID: synonymID}
}

func (c *collection) Update(ctx context.Context, schema *api.CollectionUpdateSchema) (*api.CollectionUpdateSchema, error) {
	response, err := c.apiClient.UpdateCollectionWithResponse(ctx, c.name,
		api.UpdateCollectionJSONRequestBody(*schema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
