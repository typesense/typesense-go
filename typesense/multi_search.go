package typesense

import (
	"context"
	"fmt"

	"github.com/typesense/typesense-go/typesense/api"
)

type MultiSearchInterface interface {
	Perform(searchRequest *api.MultiSearchParams, commonSearchParams api.MultiSearchParameters) (*struct {
		Results *[]api.SearchResult "json:\"results,omitempty\""
	}, error)
}

type multiSearch struct {
	apiClient APIClientInterface
}

func (m *multiSearch) Perform(searchParams *api.MultiSearchParams, commonSearchParams api.MultiSearchParameters) (*struct {
	Results *[]api.SearchResult "json:\"results,omitempty\""
}, error) {
	fmt.Printf("Hello")
	fmt.Printf("%v\n", *searchParams)
	fmt.Printf("%v\n", commonSearchParams)
	response, err := m.apiClient.MultiSearchWithResponse(context.Background(), searchParams, api.MultiSearchJSONRequestBody(commonSearchParams))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
