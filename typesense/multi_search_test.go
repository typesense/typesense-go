package typesense

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"bytes"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/typesense/mocks"
)

func newMultiSearchParams() *api.MultiSearchParams {
	return &api.MultiSearchParams{
		Q:              pointer.String("text"),
		QueryBy:        pointer.String("company_name"),
		Prefix:         pointer.String("true"),
		FilterBy:       pointer.String("num_employees:=100"),
		SortBy:         pointer.String("num_employees:desc"),
		FacetBy:        pointer.String("year_started"),
		MaxFacetValues: pointer.Int(10),
		FacetQuery:     pointer.String("facetQuery"),
		NumTypos:       pointer.String("2"),
		Page:           pointer.Int(1),
		PerPage:        pointer.Int(10),
		GroupBy:        pointer.String("country"),
		GroupLimit:     pointer.Int(3),
		IncludeFields:  pointer.String("company_name"),
	}
}

func newMultiSearchBodyParams() api.MultiSearchSearchesParameter {
	return api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: "companies",
				Q:          pointer.String("text"),
				QueryBy:    pointer.String("company_name"),
			},
			{
				Collection: "companies",
				Q:          pointer.String("text"),
				QueryBy:    pointer.String("company_name"),
			},
		},
	}
}

func newMultiSearchResult() *api.MultiSearchResult {
	return &api.MultiSearchResult{
		Results: []api.SearchResult{
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]api.FacetCounts{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]api.FacetCounts{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
		},
	}
}

func TestMultiSearchResultDeserialization(t *testing.T) {
	inputJSON := `{
			"results": [
				{
					"facet_counts": [],
					"found": 1,
					"search_time_ms": 1,
					"hits": [
						{
						"highlights": [
							{
								"field": "company_name",
								"snippet": "<mark>Stark</mark> Industries",
								"matched_tokens": ["Stark"]
							}
						],
						"document": {
								"id": "124",
								"company_name": "Stark Industries",
								"num_employees": 5215,
								"country": "USA"
							}
						}
					]
				},
				{
					"facet_counts": [],
					"found": 1,
					"search_time_ms": 1,
					"hits": [
						{
						"highlights": [
							{
								"field": "company_name",
								"snippet": "<mark>Stark</mark> Industries",
								"matched_tokens": ["Stark"]
							}
						],
						"document": {
								"id": "124",
								"company_name": "Stark Industries",
								"num_employees": 5215,
								"country": "USA"
							}
						}
					]
				}
		]
	}`
	expected := &api.MultiSearchResult{
		Results: []api.SearchResult{
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]api.FacetCounts{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
			{
				Found:        pointer.Int(1),
				SearchTimeMs: pointer.Int(1),
				FacetCounts:  &[]api.FacetCounts{},
				Hits: &[]api.SearchResultHit{
					{
						Highlights: &[]api.SearchHighlight{
							{
								Field:         pointer.String("company_name"),
								Snippet:       pointer.String("<mark>Stark</mark> Industries"),
								MatchedTokens: &[]interface{}{"Stark"},
							},
						},
						Document: &map[string]interface{}{
							"id":            "124",
							"company_name":  "Stark Industries",
							"num_employees": float64(5215),
							"country":       "USA",
						},
					},
				},
			},
		},
	}
	result := &api.MultiSearchResult{}
	err := json.Unmarshal([]byte(inputJSON), result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestMultiSearch(t *testing.T) {
	expectedParams := newMultiSearchParams()
	expectedResult := newMultiSearchResult()
	expectedBody := newMultiSearchBodyParams()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	t.Run("Perform JSON search request", func(t *testing.T) {
		mockedResult := newMultiSearchResult()

		mockAPIClient.EXPECT().
			MultiSearchWithResponse(gomock.Not(gomock.Nil()), expectedParams, api.MultiSearchJSONRequestBody(expectedBody)).Return(&api.MultiSearchResponse{
			JSON200: mockedResult,
		}, nil).Times(1)

		client := NewClient(WithAPIClient(mockAPIClient))
		params := newMultiSearchParams()
		body := newMultiSearchBodyParams()
		result, err := client.MultiSearch.Perform(context.Background(), params, body)

		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Perform with content type", func(t *testing.T) {
		expectedContentType := "application/x-json-stream"
		expectedResponseBytes, err := json.Marshal(expectedResult)
		assert.Nil(t, err)

		expectedReqBody, err := json.Marshal(expectedBody)
		assert.Nil(t, err)
		reqReader := bytes.NewReader(expectedReqBody)
		mockAPIClient.EXPECT().
			MultiSearchWithBodyWithResponse(gomock.Not(gomock.Nil()), expectedParams, expectedContentType, reqReader).
			Return(&api.MultiSearchResponse{
				Body: expectedResponseBytes,
			}, nil).Times(1)

		client := NewClient(WithAPIClient(mockAPIClient))
		params := newMultiSearchParams()
		reqBody := newMultiSearchBodyParams()
		result, err := client.MultiSearch.PerformWithContentType(context.Background(), params, reqBody, expectedContentType)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponseBytes, result.Body)
	})
}

func TestMultiSearchOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	expectedParams := newMultiSearchParams()
	expectedBody := newMultiSearchBodyParams()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		MultiSearchWithResponse(gomock.Not(gomock.Nil()), expectedParams, api.MultiSearchJSONRequestBody(expectedBody)).
		Return(&api.MultiSearchResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newMultiSearchParams()
	_, err := client.MultiSearch.Perform(context.Background(), params, newMultiSearchBodyParams())
	assert.NotNil(t, err)
}

func TestMultiSearchOnApiClientError(t *testing.T) {
	expectedParams := newMultiSearchParams()
	expectedBody := newMultiSearchBodyParams()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		MultiSearchWithResponse(gomock.Not(gomock.Nil()), expectedParams, api.MultiSearchJSONRequestBody(expectedBody)).
		Return(nil, errors.New("failed request")).Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newMultiSearchParams()
	_, err := client.MultiSearch.Perform(context.Background(), params, newMultiSearchBodyParams())
	assert.NotNil(t, err)
}
