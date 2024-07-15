package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

type ConversationInterface interface {
	Retrieve(ctx context.Context) ([]*api.ConversationSchema, error)
	Update(ctx context.Context, conversation *api.ConversationUpdateSchema) (*api.ConversationUpdateSchema, error)
	Delete(ctx context.Context) (*api.ConversationDeleteSchema, error)
}

type conversation struct {
	apiClient      APIClientInterface
	conversationId int64
}

func (c *conversation) Retrieve(ctx context.Context) ([]*api.ConversationSchema, error) {
	response, err := c.apiClient.RetrieveConversationWithResponse(ctx, c.conversationId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}

func (c *conversation) Update(ctx context.Context, conversation *api.ConversationUpdateSchema) (*api.ConversationUpdateSchema, error) {
	response, err := c.apiClient.UpdateConversationWithResponse(ctx, c.conversationId, api.UpdateConversationJSONRequestBody(*conversation))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *conversation) Delete(ctx context.Context) (*api.ConversationDeleteSchema, error) {
	response, err := c.apiClient.DeleteConversationWithResponse(ctx, c.conversationId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
