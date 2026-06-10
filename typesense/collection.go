package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// CollectionInterface is a type for Collection API operations
type CollectionInterface[T any] interface {
	// Retrieve a single collection.
	//
	// Retrieve the details of a collection, given its name.
	//
	// HTTP: GET /collections/{collectionName}
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Retrieve(ctx context.Context) (*api.CollectionResponse, error)
	// Delete a collection.
	//
	// Permanently drops a collection. This action cannot be undone. For large collections, this might have an impact on read latencies.
	//
	// HTTP: DELETE /collections/{collectionName}
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Delete(ctx context.Context) (*api.CollectionResponse, error)
	Documents() DocumentsInterface
	Document(documentID string) DocumentInterface[T]

	// Update a collection.
	//
	// Update a collection's schema to modify the fields and their types.
	//
	// HTTP: PATCH /collections/{collectionName}
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Update(context.Context, *api.CollectionUpdateSchema) (*api.CollectionUpdateSchema, error)
}

var _ CollectionInterface[any] = (*collection[any])(nil)

// collection is internal implementation of CollectionInterface
type collection[T any] struct {
	apiClient APIClientInterface
	name      string
}

// Retrieve a single collection.
//
// Retrieve the details of a collection, given its name.
//
// HTTP: GET /collections/{collectionName}
//
// See: https://typesense.org/docs/latest/api/collections.html
func (c *collection[T]) Retrieve(ctx context.Context) (*api.CollectionResponse, error) {
	response, err := c.apiClient.GetCollectionWithResponse(ctx, c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Delete a collection.
//
// Permanently drops a collection. This action cannot be undone. For large collections, this might have an impact on read latencies.
//
// HTTP: DELETE /collections/{collectionName}
//
// See: https://typesense.org/docs/latest/api/collections.html
func (c *collection[T]) Delete(ctx context.Context) (*api.CollectionResponse, error) {
	response, err := c.apiClient.DeleteCollectionWithResponse(ctx, c.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *collection[T]) Documents() DocumentsInterface {
	return &documents{apiClient: c.apiClient, collectionName: c.name}
}

func (c *collection[T]) Document(documentID string) DocumentInterface[T] {
	return &document[T]{apiClient: c.apiClient, collectionName: c.name, documentID: documentID}
}

// Update a collection.
//
// Update a collection's schema to modify the fields and their types.
//
// HTTP: PATCH /collections/{collectionName}
//
// See: https://typesense.org/docs/latest/api/collections.html
func (c *collection[T]) Update(ctx context.Context, schema *api.CollectionUpdateSchema) (*api.CollectionUpdateSchema, error) {
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
