package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

type MultiSearchInterface interface {
	Perform(searchParams *api.MultiSearchParams, commonSearchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error)
}

type multiSearch struct {
	apiClient APIClientInterface
}

func (m *multiSearch) Perform(searchParams *api.MultiSearchParams, commonSearchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error) {
	response, err := m.apiClient.MultiSearchWithResponse(context.Background(), searchParams, api.MultiSearchJSONRequestBody(commonSearchParams))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
