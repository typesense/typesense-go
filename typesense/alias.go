package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// AliasInterface is a type for Alias API operations
type AliasInterface interface {
	Retrieve(ctx context.Context) (*api.CollectionAlias, error)
	Delete(ctx context.Context) (*api.CollectionAlias, error)
}

type alias struct {
	apiClient APIClientInterface
	name      string
}

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
