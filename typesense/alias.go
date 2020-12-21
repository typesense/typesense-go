package typesense

import (
	"context"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

// AliasInterface is a type for Alias API operations
type AliasInterface interface {
	Retrieve() (*api.CollectionAlias, error)
	Delete() (*api.CollectionAlias, error)
}

type alias struct {
	apiClient api.ClientWithResponsesInterface
	name      string
}

func (a *alias) Retrieve() (*api.CollectionAlias, error) {
	response, err := a.apiClient.GetAliasWithResponse(context.Background(), a.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}

func (a *alias) Delete() (*api.CollectionAlias, error) {
	response, err := a.apiClient.DeleteAliasWithResponse(context.Background(), a.name)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}
