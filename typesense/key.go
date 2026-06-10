package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type KeyInterface interface {
	// Retrieve (metadata about) a key.
	//
	// Retrieve (metadata about) a key. Only the key prefix is returned when you retrieve a key. Due to security reasons, only the create endpoint returns the full API key.
	//
	// HTTP: GET /keys/{keyId}
	//
	// See: https://typesense.org/docs/latest/api/api-keys.html
	Retrieve(ctx context.Context) (*api.ApiKey, error)
	// Delete an API key given its ID.
	//
	// HTTP: DELETE /keys/{keyId}
	//
	// See: https://typesense.org/docs/latest/api/api-keys.html
	Delete(ctx context.Context) (*api.ApiKeyDeleteResponse, error)
}

type key struct {
	apiClient APIClientInterface
	keyID     int64
}

// Retrieve (metadata about) a key.
//
// Retrieve (metadata about) a key. Only the key prefix is returned when you retrieve a key. Due to security reasons, only the create endpoint returns the full API key.
//
// HTTP: GET /keys/{keyId}
//
// See: https://typesense.org/docs/latest/api/api-keys.html
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

// Delete an API key given its ID.
//
// HTTP: DELETE /keys/{keyId}
//
// See: https://typesense.org/docs/latest/api/api-keys.html
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
