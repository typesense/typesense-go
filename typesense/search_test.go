package typesense

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/typesense/mocks"
)

func newSearchParams() *api.SearchCollectionParams {
	return &api.SearchCollectionParams{
		Q:              "text",
		QueryBy:        []string{"company_name"},
		MaxHits:        pointer.Interface("all"),
		Prefix:         pointer.True(),
		FilterBy:       pointer.String("num_employees:=100"),
		SortBy:         &([]string{"num_employees:desc"}),
		FacetBy:        &([]string{"year_started"}),
		MaxFacetValues: pointer.Int(10),
		FacetQuery:     pointer.String("facetQuery"),
		NumTypos:       pointer.Int(2),
		Page:           pointer.Int(1),
		PerPage:        pointer.Int(10),
		GroupBy:        &([]string{"country"}),
		GroupLimit:     pointer.Int(3),
		IncludeFields:  &([]string{"company_name"}),
	}
}

func newSearchResult() *api.SearchResult {
	return &api.SearchResult{
		Found:        1,
		SearchTimeMs: 1,
		FacetCounts:  []int{},
		Hits: []api.SearchResultHit{
			{
				Highlights: []api.SearchHighlight{
					{
						Field:         "company_name",
						Snippet:       "<mark>Stark</mark> Industries",
						MatchedTokens: []interface{}{"Stark"},
					},
				},
				Document: map[string]interface{}{
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
		Found:        1,
		SearchTimeMs: 1,
		FacetCounts:  []int{},
		Hits: []api.SearchResultHit{
			{
				Highlights: []api.SearchHighlight{
					{
						Field:         "company_name",
						Snippet:       "<mark>Stark</mark> Industries",
						MatchedTokens: []interface{}{"Stark"},
					},
				},
				Document: map[string]interface{}{
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
	result, err := client.Collection("companies").Documents().Search(params)

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
	_, err := client.Collection("companies").Documents().Search(params)
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
	_, err := client.Collection("companies").Documents().Search(params)
	assert.NotNil(t, err)
}
