package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
)

func newTestServerAndClient(handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	return server, NewClient(WithServer(server.URL))
}
func TestConversationsRetrieveAllConversations(t *testing.T) {
	expectedData := api.ConversationsRetrieveSchema{
		Conversations: []*api.ConversationSchema{
			{
				Id: 1,
				Conversation: []map[string]any{
					{"any": "any"},
				},
				LastUpdated: 12,
				Ttl:         34,
			},
		},
	}
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/conversations" {
			t.Fatal("Invalid request endpoint!")
		}
		data, err := json.Marshal(expectedData)
		assert.NoError(t, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData.Conversations, res)
}

func TestConversationsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/conversations" {
			t.Fatal("Invalid request endpoint!")
		}
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Retrieve(context.Background())
	assert.NotNil(t, err)
}
