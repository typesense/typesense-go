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

var appleRule = api.CurationRule{Query: pointer.String("apple"), Match: pointer.Any(api.Exact)}

func TestCurationSetItemsRetrieve(t *testing.T) {
	expectedData := []api.CurationItemSchema{
		{Id: "apple-promo", Rule: appleRule, Includes: &[]api.CurationInclude{{Id: "422", Position: 1}}},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.CurationSet("products").Items().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestCurationSetItemsUpsert(t *testing.T) {
	item := &api.CurationItemCreateSchema{Rule: appleRule, Excludes: &[]api.CurationExclude{{Id: "287"}}}
	expectedData := &api.CurationItemSchema{Id: "apple-promo", Rule: item.Rule, Excludes: item.Excludes}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items/apple-promo", http.MethodPut)

		var reqBody api.CurationItemCreateSchema
		require.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, *item, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.CurationSet("products").Items().Upsert(context.Background(), "apple-promo", item)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestCurationSetItemsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items", http.MethodGet)
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := client.CurationSet("products").Items().Retrieve(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Status)
}

func TestCurationSetItemRetrieve(t *testing.T) {
	expectedData := &api.CurationItemSchema{Id: "apple-promo", Rule: appleRule}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items/apple-promo", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.CurationSet("products").Item("apple-promo").Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestCurationSetItemRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items/missing", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not found."}`))
	})
	defer server.Close()

	_, err := client.CurationSet("products").Item("missing").Retrieve(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Status)
}

func TestCurationSetItemUpsert(t *testing.T) {
	item := &api.CurationItemCreateSchema{Rule: appleRule, ReplaceQuery: pointer.String("iphone")}
	expectedData := &api.CurationItemSchema{Id: "apple-promo", Rule: item.Rule, ReplaceQuery: item.ReplaceQuery}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items/apple-promo", http.MethodPut)

		var reqBody api.CurationItemCreateSchema
		require.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, *item, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.CurationSet("products").Item("apple-promo").Upsert(context.Background(), item)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestCurationSetItemDelete(t *testing.T) {
	expectedData := &api.CurationItemDeleteSchema{Id: "apple-promo"}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items/apple-promo", http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.CurationSet("products").Item("apple-promo").Delete(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestCurationSetItemDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets/products/items/missing", http.MethodDelete)
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := client.CurationSet("products").Item("missing").Delete(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.Status)
}
