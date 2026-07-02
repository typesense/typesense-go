package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type MultiSearchInterface interface {
	Perform(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error)
	PerformWithContentType(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter, contentType string) (*api.MultiSearchResponse, error)
	// PerformUnion performs a multi-search and merges the results into a single `SearchResult`.
	// The `Union` field in searchParams is automatically set to `true`. If it is explicitly
	// passed as `false`, this method will return an error.
	PerformUnion(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.SearchResult, error)
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

func (m *multiSearch) PerformUnion(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.SearchResult, error) {
	if searchParams.Union != nil && !*searchParams.Union {
		return nil, errors.New("Invalid parameter: cannot set `Union` to `false` when calling PerformUnion")
	}

	// Force the Union parameter to be true
	unionTrue := true
	searchParams.Union = &unionTrue

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

	// Unmarshal the raw JSON body into SearchResult instead of MultiSearchResult
	var searchResult api.SearchResult
	if err := json.Unmarshal(response.Body, &searchResult); err != nil {
		return nil, err
	}

	return &searchResult, nil
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
