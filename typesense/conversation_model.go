package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type ConversationModelInterface interface {
	Retrieve(ctx context.Context) (*api.ConversationModelSchema, error)
	Update(ctx context.Context, schema *api.ConversationModelUpdateSchema) (*api.ConversationModelSchema, error)
	Delete(ctx context.Context) (*api.ConversationModelSchema, error)
}

type conversationModel struct {
	apiClient APIClientInterface
	modelId   string
}

func (c *conversationModel) Retrieve(ctx context.Context) (*api.ConversationModelSchema, error) {
	response, err := c.apiClient.RetrieveConversationModelWithResponse(ctx, c.modelId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *conversationModel) Update(ctx context.Context, schema *api.ConversationModelUpdateSchema) (*api.ConversationModelSchema, error) {
	response, err := c.apiClient.UpdateConversationModelWithResponse(ctx, c.modelId, api.UpdateConversationModelJSONRequestBody(*schema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *conversationModel) Delete(ctx context.Context) (*api.ConversationModelSchema, error) {
	response, err := c.apiClient.DeleteConversationModelWithResponse(ctx, c.modelId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
