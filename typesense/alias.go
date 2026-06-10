package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// AliasInterface is a type for Alias API operations
type AliasInterface interface {
	// Retrieve an alias.
	//
	// Find out which collection an alias points to by fetching it
	//
	// HTTP: GET /aliases/{aliasName}
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Retrieve(ctx context.Context) (*api.CollectionAlias, error)
	// Delete an alias.
	//
	// HTTP: DELETE /aliases/{aliasName}
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Delete(ctx context.Context) (*api.CollectionAlias, error)
}

type alias struct {
	apiClient APIClientInterface
	name      string
}

// Retrieve an alias.
//
// # Find out which collection an alias points to by fetching it
//
// HTTP: GET /aliases/{aliasName}
//
// See: https://typesense.org/docs/latest/api/collections.html
func (a *alias) Retrieve(ctx context.Context) (*api.CollectionAlias, error) {
	response, err := a.apiClient.GetAliasWithResponse(ctx, a.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Delete an alias.
//
// HTTP: DELETE /aliases/{aliasName}
//
// See: https://typesense.org/docs/latest/api/collections.html
func (a *alias) Delete(ctx context.Context) (*api.CollectionAlias, error) {
	response, err := a.apiClient.DeleteAliasWithResponse(ctx, a.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
