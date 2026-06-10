package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type NLSearchModelsInterface interface {
	// List all NL search models.
	//
	// Retrieve all NL search models.
	//
	// HTTP: GET /nl_search_models
	//
	// See: https://typesense.org/docs/latest/api/natural-language-search.html
	Retrieve(ctx context.Context) ([]*api.NLSearchModelSchema, error)
	// Create a NL search model.
	//
	// Create a new NL search model.
	//
	// HTTP: POST /nl_search_models
	//
	// See: https://typesense.org/docs/latest/api/natural-language-search.html
	Create(ctx context.Context, model *api.NLSearchModelCreateSchema) (*api.NLSearchModelSchema, error)
}

type nlSearchModels struct {
	apiClient APIClientInterface
}

// List all NL search models.
//
// Retrieve all NL search models.
//
// HTTP: GET /nl_search_models
//
// See: https://typesense.org/docs/latest/api/natural-language-search.html
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

// Create a NL search model.
//
// Create a new NL search model.
//
// HTTP: POST /nl_search_models
//
// See: https://typesense.org/docs/latest/api/natural-language-search.html
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
