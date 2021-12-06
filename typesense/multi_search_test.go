package typesense

import (
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func newMultiSearchParams() *api.MultiSearchParameters {
	return &api.MultiSearchParameters{
		Q:              pointer.String("text"),
		QueryBy:        pointer.String("company_name"),
		MaxHits:        pointer.Interface("all"),
		Prefix:         pointer.String("true"),
		FilterBy:       pointer.String("num_employees:=100"),
		SortBy:         pointer.String("num_employees:desc"),
		FacetBy:        pointer.String("year_started"),
		MaxFacetValues: pointer.Int(10),
		FacetQuery:     pointer.String("facetQuery"),
		NumTypos:       pointer.Int(2),
		Page:           pointer.Int(1),
		PerPage:        pointer.Int(10),
		GroupBy:        pointer.String("country"),
		GroupLimit:     pointer.Int(3),
		IncludeFields:  pointer.String("company_name"),
	}
}

func newMultiSearchResult() *api.SearchResult {
	return &api.SearchResult{
		Found:        pointer.Int(1),
		SearchTimeMs: pointer.Int(1),
		FacetCounts:  &[]int{},
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
