package typesense

// ConversationsInterface is a type for Conversations API operations
type ConversationsInterface interface {
	// Models manages conversation models, which configure the LLM used for conversational search.
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Models() ConversationModelsInterface
	// Model returns a handle for the conversation model with the given id.
	//
	// See: https://typesense.org/docs/latest/api/conversational-search-rag.html
	Model(modelId string) ConversationModelInterface
}

// conversations is internal implementation of ConversationsInterface
type conversations struct {
	apiClient APIClientInterface
}

func (c *conversations) Models() ConversationModelsInterface {
	return &conversationModels{apiClient: c.apiClient}
}

func (c *conversations) Model(modelId string) ConversationModelInterface {
	return &conversationModel{apiClient: c.apiClient, modelId: modelId}
}
