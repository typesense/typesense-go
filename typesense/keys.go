package typesense

import (
	"context"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

type KeysInterface interface {
	Create(key *api.ApiKeySchema) (*api.ApiKey, error)
	Retrieve() ([]*api.ApiKey, error)
}

type keys struct {
	apiClient api.ClientWithResponsesInterface
}

func (k *keys) Create(key *api.ApiKeySchema) (*api.ApiKey, error) {
	response, err := k.apiClient.CreateKeyWithResponse(context.Background(),
		api.CreateKeyJSONRequestBody(*key))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON201, nil
}

func (k *keys) Retrieve() ([]*api.ApiKey, error) {
	response, err := k.apiClient.GetKeysWithResponse(context.Background())
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200.Keys, nil
}
