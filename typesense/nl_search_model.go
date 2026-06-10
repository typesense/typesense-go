package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type NLSearchModelInterface interface {
	// Retrieve a NL search model.
	//
	// Retrieve a specific NL search model by its ID.
	//
	// HTTP: GET /nl_search_models/{modelId}
	//
	// See: https://typesense.org/docs/latest/api/natural-language-search.html
	Retrieve(ctx context.Context) (*api.NLSearchModelSchema, error)
	// Update a NL search model.
	//
	// Update an existing NL search model.
	//
	// HTTP: PUT /nl_search_models/{modelId}
	//
	// See: https://typesense.org/docs/latest/api/natural-language-search.html
	Update(ctx context.Context, model *api.NLSearchModelUpdateSchema) (*api.NLSearchModelSchema, error)
	// Delete a NL search model.
	//
	// Delete a specific NL search model by its ID.
	//
	// HTTP: DELETE /nl_search_models/{modelId}
	//
	// See: https://typesense.org/docs/latest/api/natural-language-search.html
	Delete(ctx context.Context) (*api.NLSearchModelDeleteSchema, error)
}

type nlSearchModel struct {
	apiClient APIClientInterface
	modelID   string
}

// Retrieve a NL search model.
//
// Retrieve a specific NL search model by its ID.
//
// HTTP: GET /nl_search_models/{modelId}
//
// See: https://typesense.org/docs/latest/api/natural-language-search.html
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

// Update a NL search model.
//
// Update an existing NL search model.
//
// HTTP: PUT /nl_search_models/{modelId}
//
// See: https://typesense.org/docs/latest/api/natural-language-search.html
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

// Delete a NL search model.
//
// Delete a specific NL search model by its ID.
//
// HTTP: DELETE /nl_search_models/{modelId}
//
// See: https://typesense.org/docs/latest/api/natural-language-search.html
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
