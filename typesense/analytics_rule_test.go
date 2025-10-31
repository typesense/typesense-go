package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestAnalyticsRuleRetrieve(t *testing.T) {
	expectedData := &api.AnalyticsRule{
		Name:       "test_rule",
		Type:       api.AnalyticsRuleTypeCounter,
		Collection: "test_collection",
		EventType:  "click",
		Params: &api.AnalyticsRuleCreateParams{
			CounterField: pointer.String("popularity"),
			Weight:       pointer.Int(10),
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules/test_rule", http.MethodGet)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Rule(expectedData.Name).Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRuleRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules/test_rule", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Analytics().Rule("test_rule").Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}

func TestAnalyticsRuleDelete(t *testing.T) {
	expectedData := &api.AnalyticsRule{
		Name: "test_rule",
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules/test_rule", http.MethodDelete)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Rule("test_rule").Delete(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRuleDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules/test_rule", http.MethodDelete)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Analytics().Rule("test_rule").Delete(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}
