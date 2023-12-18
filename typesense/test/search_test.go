//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func TestCollectionSearch(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	documents := []interface{}{
		newDocument("123", withCompanyName("Company 1"), withNumEmployees(50)),
		newDocument("125", withCompanyName("Company 2"), withNumEmployees(150)),
		newDocument("127", withCompanyName("Company 3"), withNumEmployees(250)),
		newDocument("129", withCompanyName("Stark Industries 4"), withNumEmployees(500)),
		newDocument("131", withCompanyName("Stark Industries 5"), withNumEmployees(1000)),
	}

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:              "Company",
		QueryBy:        "company_name, company_name",
		QueryByWeights: pointer.String("2, 1"),
		FilterBy:       pointer.String("num_employees:>=100"),
		SortBy:         pointer.String("num_employees:desc"),
		NumTypos:       pointer.String("2"),
		Page:           pointer.Int(1),
		PerPage:        pointer.Int(10),
	}

	expectedDocs := []map[string]interface{}{
		newDocumentResponse("127", withResponseCompanyName("Company 3"),
			withResponseNumEmployees(250)),
		newDocumentResponse("125", withResponseCompanyName("Company 2"),
			withResponseNumEmployees(150)),
	}

	result, err := typesenseClient.Collection(collectionName).Documents().Search(context.Background(), searchParams)

	require.NoError(t, err)
	require.Equal(t, 2, *result.Found, "found documents number is invalid")
	require.Equal(t, 2, len(*result.Hits), "number of hits is invalid")

	docs := make([]map[string]interface{}, len(*result.Hits))
	for i, hit := range *result.Hits {
		docs[i] = *hit.Document
	}

	require.Equal(t, expectedDocs, docs)
}

func TestCollectionSearchRange(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	documents := []interface{}{
		newDocument("123", withCompanyName("Company 1"), withNumEmployees(50)),
		newDocument("125", withCompanyName("Company 2"), withNumEmployees(150)),
		newDocument("127", withCompanyName("Company 3"), withNumEmployees(250)),
		newDocument("129", withCompanyName("Stark Industries 4"), withNumEmployees(500)),
	}

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:        "*",
		FilterBy: pointer.String("num_employees:>=100&&num_employees:<=300"),
		SortBy:   pointer.String("num_employees:asc"),
		Page:     pointer.Int(1),
		PerPage:  pointer.Int(10),
		QueryBy:  "company_name, country",
	}

	expectedDocs := []map[string]interface{}{
		newDocumentResponse("125", withResponseCompanyName("Company 2"),
			withResponseNumEmployees(150)),
		newDocumentResponse("127", withResponseCompanyName("Company 3"),
			withResponseNumEmployees(250)),
	}

	result, err := typesenseClient.Collection(collectionName).Documents().Search(context.Background(), searchParams)

	require.NoError(t, err)
	require.Equal(t, 2, *result.Found, "found documents number is invalid")
	require.Equal(t, 2, len(*result.Hits), "number of hits is invalid")

	docs := make([]map[string]interface{}, len(*result.Hits))
	for i, hit := range *result.Hits {
		docs[i] = *hit.Document
	}

	require.Equal(t, expectedDocs, docs)
}

func TestCollectionGroupByStringArray(t *testing.T) {
	collectionName := "tags"
	_, err := typesenseClient.Collection(collectionName).Delete(context.Background())

	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name:  "tags",
				Type:  "string[]",
				Facet: pointer.True(),
			},
		},
	}

	_, err = typesenseClient.Collections().Create(context.Background(), schema)
	require.NoError(t, err)

	type docWithArray struct {
		ID   string   `json:"id"`
		Tags []string `json:"tags"`
	}

	documents := []interface{}{
		&docWithArray{
			ID:   "1",
			Tags: []string{"go", "programming", "example"},
		},
	}

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	_, err = typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:       "*",
		GroupBy: pointer.String("tags"),
	}

	result, err := typesenseClient.Collection(collectionName).Documents().Search(context.Background(), searchParams)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, 1, *result.Found, "found documents number is invalid")
	require.Equal(t, 1, len(*result.GroupedHits), "number of grouped hits is invalid")
}
