package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// AliasesInterface is a type for Aliases API operations
type AliasesInterface interface {
	Upsert(ctx context.Context, aliasName string, aliasSchema *api.CollectionAliasSchema) (*api.CollectionAlias, error)
	Retrieve(ctx context.Context) ([]*api.CollectionAlias, error)
}

// aliases is internal implementation of AliasesInterface
type aliases struct {
	apiClient APIClientInterface
}

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
