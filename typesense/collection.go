package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// CollectionInterface is a type for Collection API operations
type CollectionInterface interface {
	Retrieve() (*api.CollectionResponse, error)
	Delete() (*api.CollectionResponse, error)
	Documents() DocumentsInterface
	Document(documentID string) DocumentInterface
	Overrides() OverridesInterface
	Override(overrideID string) OverrideInterface
	Synonyms() SynonymsInterface
	Synonym(synonymID string) SynonymInterface
}

// collection is internal implementation of CollectionInterface
type collection struct {
	apiClient APIClientInterface
	name      string
}

func (c *collection) Retrieve() (*api.CollectionResponse, error) {
	response, err := c.apiClient.GetCollectionWithResponse(context.Background(), c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *collection) Delete() (*api.CollectionResponse, error) {
	response, err := c.apiClient.DeleteCollectionWithResponse(context.Background(), c.name)
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
