package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type ConversationModelInterface interface {
	// Retrieve a conversation model.
	//
	// HTTP: GET /conversations/models/{modelId}
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Retrieve(ctx context.Context) (*api.ConversationModelSchema, error)
	// Update a conversation model.
	//
	// HTTP: PUT /conversations/models/{modelId}
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Update(ctx context.Context, schema *api.ConversationModelUpdateSchema) (*api.ConversationModelSchema, error)
	// Delete a conversation model.
	//
	// HTTP: DELETE /conversations/models/{modelId}
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Delete(ctx context.Context) (*api.ConversationModelSchema, error)
}

type conversationModel struct {
	apiClient APIClientInterface
	modelId   string
}

// Retrieve a conversation model.
//
// HTTP: GET /conversations/models/{modelId}
//
// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
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

// Update a conversation model.
//
// HTTP: PUT /conversations/models/{modelId}
//
// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
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

// Delete a conversation model.
//
// HTTP: DELETE /conversations/models/{modelId}
//
// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
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
