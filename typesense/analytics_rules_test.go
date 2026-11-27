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

	res, err := client.Analytics().Rules().Retrieve(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRulesRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Analytics().Rules().Retrieve(context.Background(), nil)
	assert.ErrorContains(t, err, "status: 409")
}

func TestAnalyticsRulesRetrieveSendsRuleTag(t *testing.T) {
	expectedData := []*api.AnalyticsRule{
		{
			Name:       "test_rule_1",
			Type:       api.AnalyticsRuleTypeCounter,
			Collection: "test_collection",
			EventType:  "click",
			RuleTag:    pointer.String("homepage"),
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules?rule_tag=homepage", http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Retrieve(context.Background(), &api.RetrieveAnalyticsRulesParams{RuleTag: pointer.String("homepage")})
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRulesRetrieveWithNilParamsSendsNoQuery(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodGet)
		assert.Empty(t, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Retrieve(context.Background(), nil)
	assert.NoError(t, err)
	assert.Empty(t, res)
}

func TestAnalyticsRulesCreate(t *testing.T) {
	createData := &api.AnalyticsRuleCreate{
		Name:       "test_rule",
		Type:       api.AnalyticsRuleTypeCounter,
		Collection: "test_collection",
		EventType:  "click",
		Params: &api.AnalyticsRuleCreateParams{
			CounterField: pointer.String("popularity"),
			Weight:       pointer.Int(10),
		},
	}
	expectedData := &api.AnalyticsRule{
		Name:       createData.Name,
		Type:       createData.Type,
		Collection: createData.Collection,
		EventType:  createData.EventType,
		Params:     createData.Params,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)

		var reqBody api.AnalyticsRuleCreate
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, *createData, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name": "test_rule", "type": "counter", "collection": "test_collection", "event_type": "click", "params": {"counter_field": "popularity", "weight": 10}}`))
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Create(context.Background(), createData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRulesCreateOnRejectedRuleReturnsHTTPError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Rule already exists"}`))
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Create(context.Background(), &api.AnalyticsRuleCreate{Name: "test_rule"})
	assert.Nil(t, res)
	var httpErr *HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.Status)
	assert.JSONEq(t, `{"message": "Rule already exists"}`, string(httpErr.Body))
}

func TestAnalyticsRulesCreateMany(t *testing.T) {
	newRule := &api.AnalyticsRuleCreate{
		Name:       "test_rule",
		Type:       api.AnalyticsRuleTypeCounter,
		Collection: "test_collection",
		EventType:  "click",
		Params: &api.AnalyticsRuleCreateParams{
			CounterField: pointer.String("popularity"),
			Weight:       pointer.Int(10),
		},
	}
	existingRule := &api.AnalyticsRuleCreate{Name: "existing_rule", Type: api.AnalyticsRuleTypeCounter, Collection: "test_collection", EventType: "click"}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)

		var reqBody []api.AnalyticsRuleCreate
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, []api.AnalyticsRuleCreate{*newRule, *existingRule}, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"name": "test_rule", "type": "counter", "collection": "test_collection", "event_type": "click", "params": {"counter_field": "popularity", "weight": 10}},
			{"error": "Rule already exists"}
		]`))
	})
	defer server.Close()

	res, err := client.Analytics().Rules().CreateMany(context.Background(), []*api.AnalyticsRuleCreate{newRule, existingRule})
	assert.NoError(t, err)
	assert.Equal(t, []*api.AnalyticsRuleCreateResponse{
		{AnalyticsRule: api.AnalyticsRule{
			Name:       newRule.Name,
			Type:       newRule.Type,
			Collection: newRule.Collection,
			EventType:  newRule.EventType,
			Params:     newRule.Params,
		}},
		{Error: "Rule already exists"},
	}, res)
}

func TestAnalyticsRulesCreateManyOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
	})
	defer server.Close()

	_, err := client.Analytics().Rules().CreateMany(context.Background(), []*api.AnalyticsRuleCreate{{Name: "test_rule"}})
	assert.ErrorContains(t, err, "status: 400")
}

func TestAnalyticsRulesCreateManyWithNoRulesSendsNoRequest(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	defer server.Close()

	for _, rules := range [][]*api.AnalyticsRuleCreate{nil, {}} {
		res, err := client.Analytics().Rules().CreateMany(context.Background(), rules)
		assert.NoError(t, err)
		assert.Empty(t, res)
	}
}

func TestAnalyticsRulesCreateManyOnObjectResponseReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules", http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name": "test_rule", "type": "counter", "collection": "test_collection", "event_type": "click"}`))
	})
	defer server.Close()

	_, err := client.Analytics().Rules().CreateMany(context.Background(), []*api.AnalyticsRuleCreate{{Name: "test_rule"}})
	assert.Error(t, err)
}

func TestAnalyticsRulesUpsert(t *testing.T) {
	upsertData := &api.AnalyticsRuleCreate{
		Name:       "test_rule",
		Type:       api.AnalyticsRuleTypeCounter,
		Collection: "test_collection",
		EventType:  "click",
		Params: &api.AnalyticsRuleCreateParams{
			CounterField: pointer.String("popularity"),
			Weight:       pointer.Int(10),
		},
	}
	expectedData := &api.AnalyticsRule{
		Name:       upsertData.Name,
		Type:       upsertData.Type,
		Collection: upsertData.Collection,
		EventType:  upsertData.EventType,
		Params:     upsertData.Params,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules/test_rule", http.MethodPut)

		var reqBody api.AnalyticsRuleCreate
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
		assert.Equal(t, *upsertData, reqBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEncode(t, expectedData))
	})
	defer server.Close()

	res, err := client.Analytics().Rules().Upsert(context.Background(), "test_rule", upsertData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestAnalyticsRulesUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/rules/test_rule", http.MethodPut)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Rule type cannot be changed"}`))
	})
	defer server.Close()

	_, err := client.Analytics().Rules().Upsert(context.Background(), "test_rule", &api.AnalyticsRuleCreate{Type: api.AnalyticsRuleTypeLog})
	var httpErr *HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.Status)
}
