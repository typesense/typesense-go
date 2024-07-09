package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

// ConversationsInterface is a type for Conversations API operations
type ConversationsInterface interface {
	Retrieve(ctx context.Context) ([]*api.ConversationSchema, error)
	Models() ConversationModelsInterface
}

// conversations is internal implementation of ConversationsInterface
type conversations struct {
	apiClient APIClientInterface
}

func (c *conversations) Retrieve(ctx context.Context) ([]*api.ConversationSchema, error) {
	response, err := c.apiClient.RetrieveAllConversationsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Conversations, nil
}

func (c *conversations) Models() ConversationModelsInterface {
	return &conversationModels{apiClient: c.apiClient}
}
