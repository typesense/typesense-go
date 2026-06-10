package typesense

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type KeysInterface interface {
	// Create an API Key.
	//
	// Create an API Key with fine-grain access control. You can restrict access on both a per-collection and per-action level. The generated key is returned only during creation. You want to store this key carefully in a secure place.
	//
	// HTTP: POST /keys
	//
	// See: https://typesense.org/docs/latest/api/api-keys.html
	Create(context.Context, *api.ApiKeySchema) (*api.ApiKey, error)
	// Retrieve (metadata about) all keys.
	//
	// HTTP: GET /keys
	//
	// See: https://typesense.org/docs/latest/api/api-keys.html
	Retrieve(context.Context) ([]*api.ApiKey, error)
	GenerateScopedSearchKey(searchKey string, params map[string]interface{}) (string, error)
}

type keys struct {
	apiClient APIClientInterface
}

// Create an API Key.
//
// Create an API Key with fine-grain access control. You can restrict access on both a per-collection and per-action level. The generated key is returned only during creation. You want to store this key carefully in a secure place.
//
// HTTP: POST /keys
//
// See: https://typesense.org/docs/latest/api/api-keys.html
func (k *keys) Create(ctx context.Context, key *api.ApiKeySchema) (*api.ApiKey, error) {
	response, err := k.apiClient.CreateKeyWithResponse(ctx,
		api.CreateKeyJSONRequestBody(*key))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}

// Retrieve (metadata about) all keys.
//
// HTTP: GET /keys
//
// See: https://typesense.org/docs/latest/api/api-keys.html
func (k *keys) Retrieve(ctx context.Context) ([]*api.ApiKey, error) {
	response, err := k.apiClient.GetKeysWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Keys, nil
}

func (k *keys) GenerateScopedSearchKey(searchKey string, params map[string]interface{}) (string, error) {
	paramsStr, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(searchKey))
	mac.Write(paramsStr)

	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	rawScopedKey := fmt.Sprintf("%s%s%s", digest, searchKey[0:4], paramsStr)
	return base64.StdEncoding.EncodeToString([]byte(rawScopedKey)), nil
}
