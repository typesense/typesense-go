package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestConversationModelRetrieve(t *testing.T) {
	accountId, systemPrompt := "CLOUDFLARE_ACCOUNT_ID", "..."
	expectedData := &api.ConversationModelSchema{
		Id:           "123",
		ModelName:    "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:       "CLOUDFLARE_API_KEY",
		AccountId:    &accountId,
		SystemPrompt: &systemPrompt,
		MaxBytes:     16384,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Model("123").Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodGet)
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := client.Conversations().Model("123").Retrieve(context.Background())
	assert.Error(t, err)
}

func TestConversationModelUpdate(t *testing.T) {
	accountId, systemPrompt := "CLOUDFLARE_ACCOUNT_ID", "..."
	expectedData := &api.ConversationModelCreateAndUpdateSchema{
		ModelName:    "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:       "CLOUDFLARE_API_KEY",
		AccountId:    &accountId,
		SystemPrompt: &systemPrompt,
		MaxBytes:     16384,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodPut)

		var reqBody api.ConversationModelCreateAndUpdateSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, expectedData, &reqBody)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Model("123").Update(context.Background(), expectedData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodPut)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Model("123").Update(context.Background(), &api.ConversationModelCreateAndUpdateSchema{
		ModelName: "cf/mistral/mistral-7b-instruct-v0.1",
	})
	assert.Error(t, err)
}

func TestConversationModelDelete(t *testing.T) {
	expectedData := &api.ConversationModelDeleteSchema{
		Id: "123",
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodDelete)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Model("123").Delete(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodDelete)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Model("123").Delete(context.Background())
	assert.Error(t, err)
}
