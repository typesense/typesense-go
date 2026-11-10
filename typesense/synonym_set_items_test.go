package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestSynonymSetItemsRetrieve(t *testing.T) {
	expectedData := []api.SynonymItemSchema{
		{Id: "smart-phone", Root: pointer.String("smart phone"), Synonyms: []string{"iphone", "android"}},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.SynonymSet("tech").Items().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestSynonymSetItemsUpsert(t *testing.T) {
	item := &api.SynonymItemUpsertSchema{Root: pointer.String("smart phone"), Synonyms: []string{"iphone", "android"}}
	expectedData := &api.SynonymItemSchema{Id: "smart-phone", Root: item.Root, Synonyms: item.Synonyms}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items/smart-phone", http.MethodPut)

		var reqBody api.SynonymItemUpsertSchema
		require.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, *item, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.SynonymSet("tech").Items().Upsert(context.Background(), "smart-phone", item)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestSynonymSetItemsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items", http.MethodGet)
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := client.SynonymSet("tech").Items().Retrieve(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Status)
}

func TestSynonymSetItemRetrieve(t *testing.T) {
	expectedData := &api.SynonymItemSchema{Id: "smart-phone", Synonyms: []string{"iphone", "android"}}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items/smart-phone", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.SynonymSet("tech").Item("smart-phone").Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestSynonymSetItemRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items/missing", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not found."}`))
	})
	defer server.Close()

	_, err := client.SynonymSet("tech").Item("missing").Retrieve(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Status)
}

func TestSynonymSetItemUpsert(t *testing.T) {
	item := &api.SynonymItemUpsertSchema{Synonyms: []string{"iphone", "android"}}
	expectedData := &api.SynonymItemSchema{Id: "smart-phone", Synonyms: item.Synonyms}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items/smart-phone", http.MethodPut)

		var reqBody api.SynonymItemUpsertSchema
		require.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, *item, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.SynonymSet("tech").Item("smart-phone").Upsert(context.Background(), item)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestSynonymSetItemDelete(t *testing.T) {
	expectedData := &api.SynonymItemDeleteSchema{Id: "smart-phone"}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items/smart-phone", http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.SynonymSet("tech").Item("smart-phone").Delete(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestSynonymSetItemDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/synonym_sets/tech/items/missing", http.MethodDelete)
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := client.SynonymSet("tech").Item("missing").Delete(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Status)
}
