package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

func TestAnalyticsEventsCreate(t *testing.T) {
	expectedData := &api.AnalyticsEventCreateSchema{
		Name: "products_click_event",
		Type: "click",
		Data: map[string]interface{}{
			"hello": "hi",
		},
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/events", http.MethodPost)

		var reqBody api.AnalyticsEventCreateSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, *expectedData, reqBody)

		data := jsonEncode(t, api.AnalyticsEventCreateResponse{
			Ok: true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Analytics().Events().Create(context.Background(), expectedData)
	assert.NoError(t, err)
	assert.True(t, res.Ok)
}

func TestAnalyticsEventsCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/analytics/events", http.MethodPost)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Analytics().Events().Create(context.Background(), &api.AnalyticsEventCreateSchema{})
	assert.ErrorContains(t, err, "status: 409")
}
