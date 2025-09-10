package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type KeyInterface interface {
	Retrieve(ctx context.Context) (*api.ApiKey, error)
	Delete(ctx context.Context) (*api.ApiKeyDeleteResponse, error)
}

type key struct {
	apiClient APIClientInterface
	keyID     int64
}

func (k *key) Retrieve(ctx context.Context) (*api.ApiKey, error) {
	response, err := k.apiClient.GetKeyWithResponse(ctx, k.keyID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (k *key) Delete(ctx context.Context) (*api.ApiKeyDeleteResponse, error) {
	response, err := k.apiClient.DeleteKeyWithResponse(ctx, k.keyID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
