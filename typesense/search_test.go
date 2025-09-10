package typesense

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func newSearchParams() *api.SearchCollectionParams {
	return &api.SearchCollectionParams{
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

func newSearchResult() *api.SearchResult {
	return &api.SearchResult{
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
	}
}

func TestSearchResultDeserialization(t *testing.T) {
	inputJSON := `{
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
	  }`
	expected := &api.SearchResult{
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
	}

	result := &api.SearchResult{}
	err := json.Unmarshal([]byte(inputJSON), &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectionSearch(t *testing.T) {
	expectedParams := newSearchParams()
	expectedResult := newSearchResult()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := newSearchResult()

	mockAPIClient.EXPECT().
		SearchCollectionWithResponse(gomock.Not(gomock.Nil()), "companies", expectedParams).
		Return(&api.SearchCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newSearchParams()
	result, err := client.Collection("companies").Documents().Search(context.Background(), params)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionSearchOnApiClientErrorReturnsError(t *testing.T) {
	expectedParams := newSearchParams()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		SearchCollectionWithResponse(gomock.Not(gomock.Nil()), "companies", expectedParams).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newSearchParams()
	_, err := client.Collection("companies").Documents().Search(context.Background(), params)
	assert.NotNil(t, err)
}

func TestCollectionSearchOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	expectedParams := newSearchParams()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		SearchCollectionWithResponse(gomock.Not(gomock.Nil()), "companies", expectedParams).
		Return(&api.SearchCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	params := newSearchParams()
	_, err := client.Collection("companies").Documents().Search(context.Background(), params)
	assert.NotNil(t, err)
}
