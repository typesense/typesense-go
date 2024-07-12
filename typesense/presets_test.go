package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestPresetsRetrieveAsSearchParameters(t *testing.T) {
	expectedData := &api.PresetsRetrieveSchema{
		Presets: []api.PresetSchema{
			{
				Name: "test",
			},
		},
	}
	expectedData.Presets[0].Value.FromSearchParameters(api.SearchParameters{Q: "Hello"})

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Retrieve(context.Background())
	assert.NoError(t, err)

	parsedRes, err := res.Presets[0].Value.AsSearchParameters()
	assert.NoError(t, err)

	parsedData, err := expectedData.Presets[0].Value.AsSearchParameters()
	assert.NoError(t, err)

	assert.Equal(t, parsedData, parsedRes)
}
func TestPresetsRetrieveAsMultiSearchSearchesParameter(t *testing.T) {
	expectedData := &api.PresetsRetrieveSchema{
		Presets: []api.PresetSchema{
			{
				Name: "test",
			},
		},
	}
	expectedData.Presets[0].Value.FromMultiSearchSearchesParameter(api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: "test",
			},
		},
	})

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Retrieve(context.Background())
	assert.NoError(t, err)

	parsedRes, err := res.Presets[0].Value.AsMultiSearchSearchesParameter()
	assert.NoError(t, err)

	parsedData, err := expectedData.Presets[0].Value.AsMultiSearchSearchesParameter()
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

func TestPresetsUpsert(t *testing.T) {
	var expectedData api.PresetUpsertSchema

	expectedData.Value.FromSearchParameters(api.SearchParameters{Q: "Xin chao"})

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/presets/test", http.MethodPut)

		var reqBody api.PresetUpsertSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, expectedData, reqBody)

		data := jsonEncode(t, map[string]any{
			"name": "test",
			"value": api.SearchParameters{
				Q: "Xin chao",
			},
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Presets().Upsert(context.Background(), "test", &expectedData)
	assert.NoError(t, err)
	assert.Equal(t, res.Name, "test")

	parsedRes, err := res.Value.AsSearchParameters()
	assert.NoError(t, err)

	parsedData, err := expectedData.Value.AsSearchParameters()
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
