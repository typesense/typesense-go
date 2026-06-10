package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

// ConversationModelsInterface is a type for ConversationModels API operations
type ConversationModelsInterface interface {
	// Create a conversation model.
	//
	// HTTP: POST /conversations/models
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Create(ctx context.Context, schema *api.ConversationModelCreateSchema) (*api.ConversationModelSchema, error)
	// List all conversation models.
	//
	// Retrieve all conversation models
	//
	// HTTP: GET /conversations/models
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Retrieve(ctx context.Context) ([]*api.ConversationModelSchema, error)
}

// conversationModels is internal implementation of ConversationModelsInterface
type conversationModels struct {
	apiClient APIClientInterface
}

// Create a conversation model.
//
// HTTP: POST /conversations/models
//
// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
func (c *conversationModels) Create(ctx context.Context, schema *api.ConversationModelCreateSchema) (*api.ConversationModelSchema, error) {
	response, err := c.apiClient.CreateConversationModelWithResponse(ctx, api.CreateConversationModelJSONRequestBody(*schema))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}

// List all conversation models.
//
// # Retrieve all conversation models
//
// HTTP: GET /conversations/models
//
// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
func (c *conversationModels) Retrieve(ctx context.Context) ([]*api.ConversationModelSchema, error) {
	response, err := c.apiClient.RetrieveAllConversationModelsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}
