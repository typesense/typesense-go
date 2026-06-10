package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type MultiSearchInterface interface {
	// Send multiple search requests in a single HTTP request.
	//
	// This is especially useful to avoid round-trip network latencies incurred otherwise if each of these requests are sent in separate HTTP requests. You can also use this feature to do a federated search across multiple collections in a single HTTP request.
	//
	// HTTP: POST /multi_search
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Perform(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error)
	// Send multiple search requests in a single HTTP request.
	//
	// This is especially useful to avoid round-trip network latencies incurred otherwise if each of these requests are sent in separate HTTP requests. You can also use this feature to do a federated search across multiple collections in a single HTTP request.
	//
	// HTTP: POST /multi_search
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	PerformWithContentType(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter, contentType string) (*api.MultiSearchResponse, error)
}

type multiSearch struct {
	apiClient APIClientInterface
}

// Send multiple search requests in a single HTTP request.
//
// This is especially useful to avoid round-trip network latencies incurred otherwise if each of these requests are sent in separate HTTP requests. You can also use this feature to do a federated search across multiple collections in a single HTTP request.
//
// HTTP: POST /multi_search
//
// See: https://typesense.org/docs/latest/api/documents.html
func (m *multiSearch) Perform(ctx context.Context, commonSearchParams *api.MultiSearchParams, searchParams api.MultiSearchSearchesParameter) (*api.MultiSearchResult, error) {
	response, err := m.apiClient.MultiSearchWithResponse(ctx, commonSearchParams, api.MultiSearchJSONRequestBody(searchParams))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

// Send multiple search requests in a single HTTP request.
//
// This is especially useful to avoid round-trip network latencies incurred otherwise if each of these requests are sent in separate HTTP requests. You can also use this feature to do a federated search across multiple collections in a single HTTP request.
//
// HTTP: POST /multi_search
//
// See: https://typesense.org/docs/latest/api/documents.html
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
	if response.Body == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response, nil
}
