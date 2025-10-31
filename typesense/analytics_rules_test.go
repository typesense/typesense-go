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

func TestAnalyticsRulesRetrieve(t *testing.T) {
	expectedData := []*api.AnalyticsRule{
		{
			Name:       "test_rule_1",
			Type:       api.AnalyticsRuleTypeCounter,
			Collection: "test_collection",
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("popularity"),
				Weight:       pointer.Int(10),
			},
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodGet)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRulesRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Analytics().Rules().Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}

func TestAnalyticsRulesCreate(t *testing.T) {
	createData := []*api.AnalyticsRuleCreate{
		{
			Name:       "test_rule",
			Type:       api.AnalyticsRuleCreateTypeCounter,
			Collection: "test_collection",
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("popularity"),
				Weight:       pointer.Int(10),
			},
		},
	}

	expectedData := []*api.AnalyticsRule{
		{
			Name:       "test_rule",
			Type:       api.AnalyticsRuleTypeCounter,
			Collection: "test_collection",
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("popularity"),
				Weight:       pointer.Int(10),
			},
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)

		var reqBody []api.AnalyticsRuleCreate
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, len(createData), len(reqBody))
		assert.Equal(t, createData[0].Name, reqBody[0].Name)

		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Create(context.Background(), createData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRulesCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
	})
	defer server.Close()

	_, err := client.Analytics().Rules().Create(context.Background(), []*api.AnalyticsRuleCreate{})
	assert.ErrorContains(t, err, "status: 400")
}
