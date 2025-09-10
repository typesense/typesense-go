package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestPresetsRetrieveAsSearchParameters(t *testing.T) {
	expectedData := []*api.PresetSchema{
		{
			Name: "test",
		},
	}

	presetValue := api.SearchParameters{Q: pointer.Any("Hello")}

	expectedData[0].Value.FromSearchParameters(presetValue)

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets", http.MethodGet)
		data := jsonEncode(t, map[string][]map[string]any{
			"presets": {{
				"name":  expectedData[0].Name,
				"value": presetValue,
			}},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)

	parsedRes, err := res[0].Value.AsSearchParameters()
	assert.NoError(t, err)

	parsedData, err := expectedData[0].Value.AsSearchParameters()
	assert.NoError(t, err)

	assert.Equal(t, parsedData, parsedRes)
}
func TestPresetsRetrieveAsMultiSearchSearchesParameter(t *testing.T) {
	expectedData := []*api.PresetSchema{
		{
			Name: "test",
		},
	}

	presetValue := api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: pointer.String("test"),
			},
		},
	}

	expectedData[0].Value.FromMultiSearchSearchesParameter(presetValue)

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets", http.MethodGet)
		data := jsonEncode(t, map[string][]map[string]any{
			"presets": {{
				"name":  expectedData[0].Name,
				"value": presetValue,
			}},
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)

	parsedRes, err := res[0].Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	parsedData, err := expectedData[0].Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	assert.Equal(t, parsedData.Searches[0], parsedRes.Searches[0])
}

func TestPresetsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Presets().Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestPresetsFromSearchParametersUpsert(t *testing.T) {
	expectedData := &api.PresetUpsertSchema{}

	presetValue := api.SearchParameters{Q: pointer.Any("Xin chao")}

	expectedData.Value.FromSearchParameters(presetValue)

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodPut)

		var reqBody api.PresetUpsertSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, *expectedData, reqBody)

		data := jsonEncode(t, map[string]any{
			"name":  "test",
			"value": presetValue,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Upsert(context.Background(), "test", expectedData)
	assert.NoError(t, err)
	assert.Equal(t, "test", res.Name)
	assert.EqualValues(t, expectedData.Value, res.Value)

	parsedRes, err := res.Value.AsSearchParameters()
	assert.NoError(t, err)

	parsedData, err := expectedData.Value.AsSearchParameters()
	assert.NoError(t, err)

	assert.Equal(t, parsedData, parsedRes)
}

func TestPresetsFromMultiSearchSearchesParameterUpsert(t *testing.T) {
	expectedData := &api.PresetUpsertSchema{}

	presetValue := api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: pointer.String("test"),
			},
		},
	}

	expectedData.Value.FromMultiSearchSearchesParameter(presetValue)

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodPut)

		var reqBody api.PresetUpsertSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, *expectedData, reqBody)

		data := jsonEncode(t, map[string]any{
			"name":  "test",
			"value": presetValue,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Upsert(context.Background(), "test", expectedData)
	assert.NoError(t, err)
	assert.Equal(t, "test", res.Name)
	assert.EqualValues(t, expectedData.Value, res.Value)

	parsedRes, err := res.Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	parsedData, err := expectedData.Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	assert.Equal(t, parsedData, parsedRes)
}

func TestPresetsUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/123", http.MethodPut)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Presets().Upsert(context.Background(), "123", &api.PresetUpsertSchema{})
	assert.NotNil(t, err)
}
