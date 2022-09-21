package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/typesense/typesense-go/typesense/api"
)

type MultiSearchInterface interface {
	Perform(searchParams *api.MultiSearchParams, commonSearchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error)
	PerformWithContentType(searchParams *api.MultiSearchParams, commonSearchParams api.MultiSearchSearchesParameter, contentType string) (*api.MultiSearchResponse, error)
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

func (m *multiSearch) PerformWithContentType(searchParams *api.MultiSearchParams, commonSearchParams api.MultiSearchSearchesParameter, contentType string) (*api.MultiSearchResponse, error) {
	body := api.MultiSearchJSONRequestBody(commonSearchParams)
	var requestReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	requestReader = bytes.NewReader(buf)
	response, err := m.apiClient.MultiSearchWithBodyWithResponse(context.Background(), searchParams, contentType, requestReader)
	if err != nil {
		return nil, err
	}
	if response.HTTPResponse == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response, nil
}
