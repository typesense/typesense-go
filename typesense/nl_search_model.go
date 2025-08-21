package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v3/typesense/api"
)

type NLSearchModelInterface interface {
	Retrieve(ctx context.Context) (*api.NLSearchModelSchema, error)
	Update(ctx context.Context, model *api.NLSearchModelUpdateSchema) (*api.NLSearchModelSchema, error)
	Delete(ctx context.Context) (*api.NLSearchModelDeleteSchema, error)
}

type nlSearchModel struct {
	apiClient APIClientInterface
	modelID   string
}

func (n *nlSearchModel) Retrieve(ctx context.Context) (*api.NLSearchModelSchema, error) {
	response, err := n.apiClient.RetrieveNLSearchModelWithResponse(ctx, n.modelID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (n *nlSearchModel) Update(ctx context.Context, model *api.NLSearchModelUpdateSchema) (*api.NLSearchModelSchema, error) {
	response, err := n.apiClient.UpdateNLSearchModelWithResponse(ctx, n.modelID, *model)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (n *nlSearchModel) Delete(ctx context.Context) (*api.NLSearchModelDeleteSchema, error) {
	response, err := n.apiClient.DeleteNLSearchModelWithResponse(ctx, n.modelID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
