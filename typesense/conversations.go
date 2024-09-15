package typesense

// ConversationsInterface is a type for Conversations API operations
type ConversationsInterface interface {
	Models() ConversationModelsInterface
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
