package typesense

import (
	"context"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

type KeyInterface interface {
	Retrieve() (*api.ApiKey, error)
	Delete() (*api.ApiKey, error)
}

type key struct {
	apiClient api.ClientWithResponsesInterface
	keyID     int64
}

func (k *key) Retrieve() (*api.ApiKey, error) {
	response, err := k.apiClient.GetKeyWithResponse(context.Background(), k.keyID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}

func (k *key) Delete() (*api.ApiKey, error) {
	response, err := k.apiClient.DeleteKeyWithResponse(context.Background(), k.keyID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}
