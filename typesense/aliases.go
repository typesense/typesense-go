package typesense

import (
	"context"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

// AliasesInterface is a type for Aliases API operations
type AliasesInterface interface {
	Upsert(aliasName string, aliasSchema *api.CollectionAliasSchema) (*api.CollectionAlias, error)
	Retrieve() ([]*api.CollectionAlias, error)
}

// aliases is internal implementation of AliasesInterface
type aliases struct {
	apiClient api.ClientWithResponsesInterface
}

func (a *aliases) Upsert(aliasName string, aliasSchema *api.CollectionAliasSchema) (*api.CollectionAlias, error) {
	response, err := a.apiClient.UpsertAliasWithResponse(context.Background(),
		aliasName, api.UpsertAliasJSONRequestBody(*aliasSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}

func (a *aliases) Retrieve() ([]*api.CollectionAlias, error) {
	response, err := a.apiClient.GetAliasesWithResponse(context.Background())
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200.Aliases, nil
}
