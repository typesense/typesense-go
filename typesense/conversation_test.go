package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

func TestConversationRetrieveConversation(t *testing.T) {
	expectedData := &api.ConversationSchema{
		Id: "123",
		Conversation: []map[string]any{
			{
				"user": "can you suggest an action series",
			},
			{
				"assistant": "...",
			}},
		LastUpdated: 12,
		Ttl:         34,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodGet)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversation(expectedData.Id).Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversation("123").Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}
func TestConversationUpdateConversation(t *testing.T) {
	expectedData := &api.ConversationSchema{
		Id: "123",
	}
	updateData := &api.ConversationUpdateSchema{
		Ttl: 3000,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodPut)

		var reqBody api.ConversationUpdateSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, updateData, &reqBody)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversation("123").Update(context.Background(), updateData)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodPut)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversation("123").Update(context.Background(), &api.ConversationUpdateSchema{
		Ttl: 0,
	})
	assert.ErrorContains(t, err, "status: 409")
}

func TestConversationDeleteConversation(t *testing.T) {
	expectedData := &api.ConversationSchema{
		Id: "123",
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodDelete)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversation("123").Delete(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodDelete)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversation("123").Delete(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}
