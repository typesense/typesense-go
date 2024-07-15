package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

type ConversationModelInterface interface {
	Retrieve(ctx context.Context) (*api.ConversationModelSchema, error)
	Update(ctx context.Context, model *api.ConversationModelCreateAndUpdateSchema) (*api.ConversationModelCreateAndUpdateSchema, error)
	Delete(ctx context.Context) (*api.ConversationModelDeleteSchema, error)
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

func (c *conversationModel) Update(ctx context.Context, conversationModelCreateAndUpdateSchema *api.ConversationModelCreateAndUpdateSchema) (*api.ConversationModelCreateAndUpdateSchema, error) {
	response, err := c.apiClient.UpdateConversationModelWithResponse(ctx, c.modelId, api.UpdateConversationModelJSONRequestBody(*conversationModelCreateAndUpdateSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (c *conversationModel) Delete(ctx context.Context) (*api.ConversationModelDeleteSchema, error) {
	response, err := c.apiClient.DeleteConversationModelWithResponse(ctx, c.modelId)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
