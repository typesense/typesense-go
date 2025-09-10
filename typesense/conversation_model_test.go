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

func TestConversationModelRetrieve(t *testing.T) {
	expectedData := &api.ConversationModelSchema{
		Id:           "123",
		ModelName:    "cloudflare/@cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:       pointer.String("CLOUDFLARE_API_KEY"),
		AccountId:    pointer.String("CLOUDFLARE_ACCOUNT_ID"),
		SystemPrompt: pointer.String("..."),
		MaxBytes:     16384,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
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
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Model("123").Retrieve(context.Background())
	assert.ErrorContains(t, err, "status: 409")
}

func TestConversationModelUpdate(t *testing.T) {
	model := &api.ConversationModelUpdateSchema{
		ModelName:    pointer.String("cf/mistral/mistral-7b-instruct-v0.1"),
		ApiKey:       pointer.String("CLOUDFLARE_API_KEY"),
		AccountId:    pointer.String("CLOUDFLARE_ACCOUNT_ID"),
		SystemPrompt: pointer.String("..."),
		MaxBytes:     pointer.Int(16384),
	}
	expectedData := &api.ConversationModelSchema{
		Id:           "123",
		ModelName:    *model.ModelName,
		ApiKey:       model.ApiKey,
		AccountId:    model.AccountId,
		SystemPrompt: model.SystemPrompt,
		MaxBytes:     *model.MaxBytes,
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodPut)

		var reqBody api.ConversationModelUpdateSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, model, &reqBody)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Model("123").Update(context.Background(), model)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodPut)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Model("123").Update(context.Background(), &api.ConversationModelUpdateSchema{})
	assert.ErrorContains(t, err, "status: 409")
}

func TestConversationModelDelete(t *testing.T) {
	expectedData := &api.ConversationModelSchema{
		Id:                "123",
		ModelName:         "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:            pointer.String("CLOUDFLARE_API_KEY"),
		AccountId:         pointer.String("CLOUDFLARE_ACCOUNT_ID"),
		SystemPrompt:      pointer.String("..."),
		MaxBytes:          16384,
		HistoryCollection: "conversation-store",
		Ttl:               pointer.Int(10000),
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models/123", http.MethodDelete)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
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
	assert.ErrorContains(t, err, "status: 409")
}
