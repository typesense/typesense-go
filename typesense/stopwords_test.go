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

func TestStopwordsRetrieve(t *testing.T) {
	expectedData := []api.StopwordsSetSchema{
		{
			Id:        "stopwords_set1",
			Locale:    pointer.String("en"),
			Stopwords: []string{"Germany", "France", "Italy", "United States"},
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/stopwords", http.MethodGet)
		data := jsonEncode(t, map[string][]api.StopwordsSetSchema{
			"stopwords": expectedData,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Stopwords().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestStopwordsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/stopwords", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Stopwords().Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}

func TestStopwordsUpsert(t *testing.T) {
	upsertData := &api.StopwordsSetUpsertSchema{
		Locale:    pointer.String("en"),
		Stopwords: []string{"Germany", "France", "Italy", "United States"},
	}

	expectedData := &api.StopwordsSetSchema{
		Id:        "stopwords_set1",
		Locale:    upsertData.Locale,
		Stopwords: upsertData.Stopwords,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/stopwords/stopwords_set1", http.MethodPut)

		var reqBody api.StopwordsSetUpsertSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, *upsertData, reqBody)

		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Stopwords().Upsert(context.Background(), "stopwords_set1", upsertData)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestStopwordsUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/stopwords/stopwords_set1", http.MethodPut)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Stopwords().Upsert(context.Background(), "stopwords_set1", &api.StopwordsSetUpsertSchema{})
	assert.ErrorContains(t, err, "status: 409")
}
