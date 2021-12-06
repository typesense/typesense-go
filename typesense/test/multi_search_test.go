//go:build integration
// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func TestMultiSearch(t *testing.T) {
	collectionName1 := createNewCollection(t, "companies")
	collectionName2 := createNewCollection(t, "companies")
	documents := []interface{}{
		newDocument("123", withCompanyName("Company 1"), withNumEmployees(50)),
		newDocument("125", withCompanyName("Company 2"), withNumEmployees(150)),
		newDocument("127", withCompanyName("Company 3"), withNumEmployees(250)),
		newDocument("129", withCompanyName("Stark Industries 4"), withNumEmployees(500)),
		newDocument("131", withCompanyName("Stark Industries 5"), withNumEmployees(1000)),
	}

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	_, err := typesenseClient.Collection(collectionName1).Documents().Import(documents, params)
	require.NoError(t, err)

	_, err = typesenseClient.Collection(collectionName1).Documents().Import(documents, params)
	require.NoError(t, err)

	searchParams := &api.MultiSearchParams{
		FilterBy: pointer.String("num_employees:>100"),
		Q:        "Company",
		QueryBy:  "company_name",
	}

	searches := api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: collectionName1,
				MultiSearchParameters: api.MultiSearchParameters{
					FilterBy: pointer.String("num_employees:>100"),
					SortBy:   pointer.String("num_employees:desc"),
				},
			},
			{
				Collection: collectionName1,
				MultiSearchParameters: api.MultiSearchParameters{
					FilterBy: pointer.String("num_employees:>1000"),
				},
			},
			{
				Collection: collectionName2,
				MultiSearchParameters: api.MultiSearchParameters{
					FilterBy: pointer.String("num_employees:>100"),
				},
			},
		},
	}

	expectedDocs := []map[string]interface{}{
		newDocumentResponse("127", withResponseCompanyName("Company 3"),
			withResponseNumEmployees(250)),
		newDocumentResponse("125", withResponseCompanyName("Company 2"),
			withResponseNumEmployees(150)),
	}

	result, err := typesenseClient.MultiSearch.Perform(searchParams, searches)
	require.NoError(t, err)

	require.Equal(t, 1, len(result.Results))

	// Check first result
	for i, doc := range *result.Results[0].Hits {
		require.Equal(t, *doc.Document, expectedDocs[i])
	}
}
