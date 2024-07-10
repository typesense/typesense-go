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

func validateRequestMetadata(t *testing.T, r *http.Request, expectedEndpoint string, expectedMethod string) {
	if r.RequestURI != expectedEndpoint {
		t.Fatal("Invalid request endpoint!")
	}
	if r.Method != expectedMethod {
		t.Fatal("Invalid HTTP method!")
	}
}

func jsonEncode(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	return data
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
		validateRequestMetadata(t, r, "/conversations", http.MethodGet)
		data := jsonEncode(t, expectedData)
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
		validateRequestMetadata(t, r, "/conversations", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Retrieve(context.Background())
	assert.NotNil(t, err)
}
