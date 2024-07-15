package typesense

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestConversationModelsRetrieve(t *testing.T) {
	accountId, systemPrompt := "CLOUDFLARE_ACCOUNT_ID", "..."
	expectedData := []*api.ConversationModelSchema{{
		Id:           "123",
		ModelName:    "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:       "CLOUDFLARE_API_KEY",
		AccountId:    &accountId,
		SystemPrompt: &systemPrompt,
		MaxBytes:     16384,
	}}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Models().Retrieve(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodGet)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Models().Retrieve(context.Background())
	assert.Error(t, err)
}

func TestConversationModelsCreate(t *testing.T) {
	accountId, systemPrompt := "CLOUDFLARE_ACCOUNT_ID", "..."
	expectedData := &api.ConversationModelCreateAndUpdateSchema{
		ModelName:    "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:       "CLOUDFLARE_API_KEY",
		AccountId:    &accountId,
		SystemPrompt: &systemPrompt,
		MaxBytes:     16384,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodPost)

		var reqBody api.ConversationModelCreateAndUpdateSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, expectedData, &reqBody)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Models().Create(context.Background(), expectedData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelsCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodPost)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Models().Create(context.Background(), &api.ConversationModelCreateAndUpdateSchema{
		ModelName: "cf/mistral/mistral-7b-instruct-v0.1",
	})
	assert.Error(t, err)
}
