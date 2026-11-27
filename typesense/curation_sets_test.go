package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestCurationSetsRetrieve(t *testing.T) {
	expectedData := []api.CurationSetSchema{
		{Name: "promos", Items: []api.CurationItemCreateSchema{
			{Id: pointer.String("apple-promo"), Rule: api.CurationRule{Query: pointer.String("apple")}, Includes: &[]api.CurationInclude{{Id: "422", Position: 1}}},
		}},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.CurationSets().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestCurationSetsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/curation_sets", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Forbidden - a valid x-typesense-api-key header must be sent."}`))
	})
	defer server.Close()

	_, err := client.CurationSets().Retrieve(context.Background())
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
}
