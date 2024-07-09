package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// ConversationModelsInterface is a type for ConversationModels API operations
type ConversationModelsInterface interface {
	Create(ctx context.Context, conversationModelSchema *api.ConversationModelSchema) (*api.ConversationModelResponse, error)
	Retrieve(ctx context.Context) ([]*api.ConversationModelResponse, error)
}

// conversationModels is internal implementation of ConversationModelsInterface
type conversationModels struct {
	apiClient APIClientInterface
}

func (c *conversationModels) Create(ctx context.Context, conversationModelSchema *api.ConversationModelSchema) (*api.ConversationModelResponse, error) {
	response, err := c.apiClient.CreateConversationModelWithResponse(ctx, api.CreateConversationModelJSONRequestBody(*conversationModelSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}

func (c *conversationModels) Retrieve(ctx context.Context) ([]*api.ConversationModelResponse, error) {
	response, err := c.apiClient.RetrieveAllConversationModelsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}
