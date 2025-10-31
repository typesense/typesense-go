package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type NLSearchModelsInterface interface {
	Retrieve(ctx context.Context) ([]*api.NLSearchModelSchema, error)
	Create(ctx context.Context, model *api.NLSearchModelCreateSchema) (*api.NLSearchModelSchema, error)
}

type nlSearchModels struct {
	apiClient APIClientInterface
}

func (n *nlSearchModels) Retrieve(ctx context.Context) ([]*api.NLSearchModelSchema, error) {
	response, err := n.apiClient.RetrieveAllNLSearchModelsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}

	// Convert []NLSearchModelSchema to []*NLSearchModelSchema
	result := make([]*api.NLSearchModelSchema, len(*response.JSON200))
	for i, model := range *response.JSON200 {
		modelCopy := model // Create a copy to get address
		result[i] = &modelCopy
	}
	return result, nil
}

func (n *nlSearchModels) Create(ctx context.Context, model *api.NLSearchModelCreateSchema) (*api.NLSearchModelSchema, error) {
	response, err := n.apiClient.CreateNLSearchModelWithResponse(ctx, *model)
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}
