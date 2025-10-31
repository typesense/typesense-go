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

func TestConversationModelsRetrieve(t *testing.T) {
	expectedData := []*api.ConversationModelSchema{{
		Id:                "123",
		ModelName:         "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:            pointer.String("CLOUDFLARE_API_KEY"),
		AccountId:         pointer.String("CLOUDFLARE_ACCOUNT_ID"),
		SystemPrompt:      pointer.String("..."),
		MaxBytes:          16384,
		HistoryCollection: "conversation-store",
		Ttl:               pointer.Int(10000),
	}}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodGet)
		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
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
	assert.ErrorContains(t, err, "status: 409")
}

func TestConversationModelsCreate(t *testing.T) {
	model := &api.ConversationModelCreateSchema{
		ModelName:         "cf/mistral/mistral-7b-instruct-v0.1",
		ApiKey:            pointer.String("CLOUDFLARE_API_KEY"),
		AccountId:         pointer.String("CLOUDFLARE_ACCOUNT_ID"),
		SystemPrompt:      pointer.String("..."),
		MaxBytes:          16384,
		HistoryCollection: "conversation-store",
		Ttl:               pointer.Int(10000),
	}
	expectedData := &api.ConversationModelSchema{
		Id:                "123",
		ModelName:         model.ModelName,
		ApiKey:            model.ApiKey,
		AccountId:         model.AccountId,
		SystemPrompt:      model.SystemPrompt,
		MaxBytes:          model.MaxBytes,
		HistoryCollection: model.HistoryCollection,
		Ttl:               pointer.Int(10000),
	}

	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodPost)

		var reqBody api.ConversationModelCreateSchema
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		assert.NoError(t, err)
		assert.Equal(t, model, &reqBody)

		data := jsonEncode(t, expectedData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	})
	defer server.Close()

	res, err := client.Conversations().Models().Create(context.Background(), model)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}

func TestConversationModelsCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	server, client := newTestServerAndClient(func(w http.ResponseWriter, r *http.Request) {
		validateRequestMetadata(t, r, "/conversations/models", http.MethodPost)
		w.WriteHeader(http.StatusConflict)
	})
	defer server.Close()

	_, err := client.Conversations().Models().Create(context.Background(), &api.ConversationModelCreateSchema{})
	assert.ErrorContains(t, err, "status: 409")
}
