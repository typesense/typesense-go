package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func TestAnalyticsEventsCreate(t *testing.T) {
	eventData := &api.AnalyticsEvent{
		Name:      "test_rule",
		EventType: "click",
		Data: struct {
			AnalyticsTag *string   `json:"analytics_tag,omitempty"`
			DocId        *string   `json:"doc_id,omitempty"`
			DocIds       *[]string `json:"doc_ids,omitempty"` //nolint:revive // matches API type
			Q            *string   `json:"q,omitempty"`
			UserId       *string   `json:"user_id,omitempty"`
		}{
			DocId:  stringPtr("123"),
			UserId: stringPtr("user_123"),
		},
	}

	expectedResponse := &api.AnalyticsEventCreateResponse{
		Ok: true,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/events", http.MethodPost)

		var reqBody api.AnalyticsEvent
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, eventData.Name, reqBody.Name)
		assert.Equal(t, eventData.EventType, reqBody.EventType)

		data := jsonEncode(t, expectedResponse)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Events().Create(context.Background(), eventData)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, res)
}

func TestAnalyticsEventsRetrieve(t *testing.T) {
	params := &api.GetAnalyticsEventsParams{
		Name:   "test_rule",
		UserId: "user_123",
		N:      10,
	}

	expectedResponse := &api.AnalyticsEventsResponse{
		Events: []struct {
			Collection *string   `json:"collection,omitempty"`
			DocId      *string   `json:"doc_id,omitempty"`
			DocIds     *[]string `json:"doc_ids,omitempty"` //nolint:revive // matches API type
			EventType  *string   `json:"event_type,omitempty"`
			Name       *string   `json:"name,omitempty"`
			Query      *string   `json:"query,omitempty"`
			Timestamp  *int64    `json:"timestamp,omitempty"`
			UserId     *string   `json:"user_id,omitempty"`
		}{
			{
				Name:      stringPtr("test_rule"),
				EventType: stringPtr("click"),
				DocId:     stringPtr("123"),
				UserId:    stringPtr("user_123"),
			},
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/events", http.MethodGet)

		// Check query parameters
		assert.Equal(t, "test_rule", r.URL.Query().Get("name"))
		assert.Equal(t, "user_123", r.URL.Query().Get("user_id"))
		assert.Equal(t, "10", r.URL.Query().Get("n"))

		data := jsonEncode(t, expectedResponse)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Events().Retrieve(context.Background(), params)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, res)
}

func TestAnalyticsEventsCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/events", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
	})
	defer server.Close()

	_, err := client.Analytics().Events().Create(context.Background(), &api.AnalyticsEvent{})
	assert.ErrorContains(t, err, "status: 400")
}

func TestAnalyticsEventsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/events", http.MethodGet)
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := client.Analytics().Events().Retrieve(context.Background(), &api.GetAnalyticsEventsParams{})
	assert.ErrorContains(t, err, "status: 500")
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
