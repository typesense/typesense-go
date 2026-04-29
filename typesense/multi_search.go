package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type MultiSearchInterface interface {
	Perform(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error)
	PerformWithContentType(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter, contentType string) (*api.MultiSearchResponse, error)
}

type multiSearch struct {
	apiClient APIClientInterface
}

func (m *multiSearch) Perform(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error) {
	response, err := m.apiClient.MultiSearchWithResponse(ctx, commonSearchParams, api.MultiSearchJSONRequestBody(searchParams))
	if err != nil {
		return nil, err
	}
	if err := multiSearchTopLevelError(response); err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (m *multiSearch) PerformWithContentType(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter, contentType string) (*api.MultiSearchResponse, error) {
	body := api.MultiSearchJSONRequestBody(searchParams)
	var requestReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	requestReader = bytes.NewReader(buf)
	response, err := m.apiClient.MultiSearchWithBodyWithResponse(ctx, commonSearchParams, contentType, requestReader)
	if err != nil {
		return nil, err
	}
	if err := multiSearchTopLevelError(response); err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response, nil
}

func multiSearchTopLevelError(response *api.MultiSearchResponse) error {
	if response == nil || len(response.Body) == 0 {
		return nil
	}

	var errorResponse struct {
		Code  *int    `json:"code"`
		Error *string `json:"error"`
	}
	if err := json.Unmarshal(response.Body, &errorResponse); err != nil {
		return nil
	}
	if errorResponse.Code == nil || errorResponse.Error == nil {
		return nil
	}

	return &HTTPError{Status: *errorResponse.Code, Body: response.Body}
}
