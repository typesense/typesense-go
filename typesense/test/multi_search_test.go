//go:build integration
// +build integration

package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func TestMultiSearch(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	documents := []interface{}{
		newDocument("123", withCompanyName("Company 1"), withNumEmployees(50)),
		newDocument("125", withCompanyName("Company 2"), withNumEmployees(150)),
		newDocument("127", withCompanyName("Company 3"), withNumEmployees(250)),
		newDocument("129", withCompanyName("Stark Industries 4"), withNumEmployees(500)),
		newDocument("131", withCompanyName("Stark Industries 5"), withNumEmployees(1000)),
	}

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(documents, params)
	require.NoError(t, err)

	searchParams := &api.MultiSearchParams{
		FilterBy: pointer.String("num_employees:>100"),
		Q:        "Company",
		QueryBy:  "company_name",
	}

	searches := api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: collectionName,
				MultiSearchParameters: api.MultiSearchParameters{
					Q:       pointer.String("Company"),
					QueryBy: pointer.String("company_name"),
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
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v\n", (*result.Results[0].Hits)[0].Document)
	_ = result
	_ = expectedDocs
	_ = err
}
