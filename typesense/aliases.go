package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// AliasesInterface is a type for Aliases API operations
type AliasesInterface interface {
	// Create or update a collection alias.
	//
	// Create or update a collection alias. An alias is a virtual collection name that points to a real collection. If you're familiar with symbolic links on Linux, it's very similar to that. Aliases are useful when you want to reindex your data in the background on a new collection and switch your application to it without any changes to your code.
	//
	// HTTP: PUT /aliases/{aliasName}
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Upsert(ctx context.Context, aliasName string, aliasSchema *api.CollectionAliasSchema) (*api.CollectionAlias, error)
	// List all aliases.
	//
	// List all aliases and the corresponding collections that they map to.
	//
	// HTTP: GET /aliases
	//
	// See: https://typesense.org/docs/latest/api/collections.html
	Retrieve(ctx context.Context) ([]*api.CollectionAlias, error)
}

// aliases is internal implementation of AliasesInterface
type aliases struct {
	apiClient APIClientInterface
}

// Create or update a collection alias.
//
// Create or update a collection alias. An alias is a virtual collection name that points to a real collection. If you're familiar with symbolic links on Linux, it's very similar to that. Aliases are useful when you want to reindex your data in the background on a new collection and switch your application to it without any changes to your code.
//
// HTTP: PUT /aliases/{aliasName}
//
// See: https://typesense.org/docs/latest/api/collections.html
func (a *aliases) Upsert(ctx context.Context, aliasName string, aliasSchema *api.CollectionAliasSchema) (*api.CollectionAlias, error) {
	response, err := a.apiClient.UpsertAliasWithResponse(ctx,
		aliasName, api.UpsertAliasJSONRequestBody(*aliasSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// List all aliases.
//
// List all aliases and the corresponding collections that they map to.
//
// HTTP: GET /aliases
//
// See: https://typesense.org/docs/latest/api/collections.html
func (a *aliases) Retrieve(ctx context.Context) ([]*api.CollectionAlias, error) {
	response, err := a.apiClient.GetAliasesWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Aliases, nil
}
