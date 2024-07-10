package typesense

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestConversationRetrieveConversation(t *testing.T) {
	expectedData := []*api.ConversationSchema{{
		Id: 1,
		Conversation: []map[string]any{
			{"any": "any"},
		},
		LastUpdated: 12,
		Ttl:         34,
	}}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodGet)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversation(123).Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversation(123).Retrieve(context.Background())
	assert.Error(t, err)
}
func TestConversationUpdateConversation(t *testing.T) {
	expectedData := &api.ConversationUpdateSchema{
		Ttl: 3000,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodPut)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversation(123).Update(context.Background(), expectedData)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodPut)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversation(123).Update(context.Background(), &api.ConversationUpdateSchema{
		Ttl: 0,
	})
	assert.Error(t, err)
}

func TestConversationDeleteConversation(t *testing.T) {
	expectedData := &api.ConversationDeleteSchema{
		Id: 123,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodDelete)
		data := jsonEncode(t, expectedData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversation(123).Delete(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/123", http.MethodDelete)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversation(123).Delete(context.Background())
	assert.Error(t, err)
}
