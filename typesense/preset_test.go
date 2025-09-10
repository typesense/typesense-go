package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestPresetRetrieveAsSearchParameters(t *testing.T) {
	expectedData := &api.PresetSchema{
		Name: "test",
	}
	presetValue := api.SearchParameters{Q: pointer.Any("Hello")}

	expectedData.Value.FromSearchParameters(presetValue)

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodGet)
		data := jsonEncode(t, map[string]any{
			"name":  "test",
			"value": presetValue,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Preset("test").Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)

	parsedRes, err := res.Value.AsSearchParameters()
	assert.NoError(t, err)

	parsedData, err := expectedData.Value.AsSearchParameters()
	assert.NoError(t, err)

	assert.Equal(t, parsedData, parsedRes)
}
func TestPresetRetrieveAsMultiSearchSearchesParameter(t *testing.T) {
	expectedData := &api.PresetSchema{
		Name: "test",
	}
	presetValue := api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: pointer.String("test"),
			},
		},
	}

	expectedData.Value.FromMultiSearchSearchesParameter(presetValue)

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodGet)
		data := jsonEncode(t, map[string]any{
			"name":  expectedData.Name,
			"value": presetValue,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Preset("test").Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)

	parsedRes, err := res.Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	parsedData, err := expectedData.Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	assert.Equal(t, parsedData, parsedRes)
}

func TestPresetRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Preset("test").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestPresetDelete(t *testing.T) {
	expectedData := &api.PresetDeleteSchema{
		Name: "test",
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodDelete)

		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Preset("test").Delete(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestPresetDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/123", http.MethodDelete)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Preset("123").Delete(context.Background())
	assert.NotNil(t, err)
}
