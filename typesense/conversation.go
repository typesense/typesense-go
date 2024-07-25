package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v2/typesense/api"
)

type ConversationInterface interface {
	Retrieve(ctx context.Context) (*api.ConversationSchema, error)
	Update(ctx context.Context, conversation *api.ConversationUpdateSchema) (*api.ConversationSchema, error)
	Delete(ctx context.Context) (*api.ConversationSchema, error)
}

type conversation struct {
	apiClient      APIClientInterface
	conversationId string
}

func (c *conversation) Retrieve(ctx context.Context) (*api.ConversationSchema, error) {
	response, err := c.apiClient.RetrieveConversationWithResponse(ctx, c.conversationId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *conversation) Update(ctx context.Context, conversation *api.ConversationUpdateSchema) (*api.ConversationSchema, error) {
	response, err := c.apiClient.UpdateConversationWithResponse(ctx, c.conversationId, api.UpdateConversationJSONRequestBody(*conversation))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *conversation) Delete(ctx context.Context) (*api.ConversationSchema, error) {
	response, err := c.apiClient.DeleteConversationWithResponse(ctx, c.conversationId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
