//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
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

	params := &api.ImportDocumentsParams{Action: pointer.Any(api.Create)}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:              pointer.String("Company"),
		QueryBy:        pointer.String("company_name, company_name"),
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

	params := &api.ImportDocumentsParams{Action: pointer.Any(api.Create)}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:        pointer.String("*"),
		FilterBy: pointer.String("num_employees:>=100&&num_employees:<=300"),
		SortBy:   pointer.String("num_employees:asc"),
		Page:     pointer.Int(1),
		PerPage:  pointer.Int(10),
		QueryBy:  pointer.String("company_name, country"),
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

	params := &api.ImportDocumentsParams{Action: pointer.Any(api.Create)}
	_, err = typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:       pointer.String("*"),
		GroupBy: pointer.String("tags"),
	}

	result, err := typesenseClient.Collection(collectionName).Documents().Search(context.Background(), searchParams)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, 1, *result.Found, "found documents number is invalid")
	require.Equal(t, 1, len(*result.GroupedHits), "number of grouped hits is invalid")
}

func TestCollectionSearchWithPreset(t *testing.T) {
	t.Cleanup(presetsCleanUp)

	collectionName := createNewCollection(t, "companies")
	documents := []interface{}{
		newDocument("123", withCompanyName("Company 1"), withNumEmployees(50)),
		newDocument("125", withCompanyName("Company 2"), withNumEmployees(150)),
		newDocument("127", withCompanyName("Company 3"), withNumEmployees(250)),
	}

	params := &api.ImportDocumentsParams{Action: pointer.Any(api.Create)}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	searchParams := api.SearchParameters{
		Q:              pointer.Any("Company"),
		QueryBy:        pointer.Any("company_name, company_name"),
		QueryByWeights: pointer.String("2, 1"),
		FilterBy:       pointer.String("num_employees:>=100"),
		SortBy:         pointer.String("num_employees:desc"),
		NumTypos:       pointer.String("2"),
		Page:           pointer.Int(1),
		PerPage:        pointer.Int(10),
	}

	presetName := newUUIDName("preset-single-collection-search")
	preset := &api.PresetUpsertSchema{}
	preset.Value.FromSearchParameters(searchParams)

	_, err = typesenseClient.Presets().Upsert(context.Background(), presetName, preset)
	require.NoError(t, err)

	expectedDocs := []map[string]interface{}{
		newDocumentResponse("127", withResponseCompanyName("Company 3"),
			withResponseNumEmployees(250)),
		newDocumentResponse("125", withResponseCompanyName("Company 2"),
			withResponseNumEmployees(150)),
	}

	result, err := typesenseClient.Collection(collectionName).Documents().Search(context.Background(), &api.SearchCollectionParams{
		Preset: &presetName,
	})

	require.NoError(t, err)
	require.Equal(t, 2, *result.Found, "found documents number is invalid")
	require.Equal(t, 2, len(*result.Hits), "number of hits is invalid")

	docs := make([]map[string]interface{}, len(*result.Hits))
	for i, hit := range *result.Hits {
		docs[i] = *hit.Document
	}

	require.Equal(t, expectedDocs, docs)
}

func TestCollectionSearchWithStopwords(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	documents := []interface{}{
		newDocument("123", withCompanyName("Company 1"), withNumEmployees(50)),
		newDocument("125", withCompanyName("Company 2"), withNumEmployees(150)),
		newDocument("127", withCompanyName("Company Stark Industries 3"), withNumEmployees(1000)),
		newDocument("129", withCompanyName("Stark Industries 4"), withNumEmployees(2000)),
	}

	params := &api.ImportDocumentsParams{Action: pointer.Any(api.Create)}
	_, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)
	require.NoError(t, err)

	stopwordsSetID := newUUIDName("stopwordsSet-test")
	upsertData := &api.StopwordsSetUpsertSchema{
		Locale:    pointer.String("en"),
		Stopwords: []string{"Stark Industries"},
	}

	_, err = typesenseClient.Stopwords().Upsert(context.Background(), stopwordsSetID, upsertData)
	require.NoError(t, err)

	searchParams := &api.SearchCollectionParams{
		Q:         pointer.String("Company Stark"),
		QueryBy:   pointer.String("company_name"),
		SortBy:    pointer.String("num_employees:desc"),
		Stopwords: pointer.String(stopwordsSetID),
	}

	expectedDocs := []map[string]interface{}{
		newDocumentResponse("127", withResponseCompanyName("Company Stark Industries 3"),
			withResponseNumEmployees(1000)),
		newDocumentResponse("125", withResponseCompanyName("Company 2"),
			withResponseNumEmployees(150)),
		newDocumentResponse("123", withResponseCompanyName("Company 1"),
			withResponseNumEmployees(50)),
	}

	result, err := typesenseClient.Collection(collectionName).Documents().Search(context.Background(), searchParams)

	require.NoError(t, err)
	require.Equal(t, 3, *result.Found, "found documents number is invalid")
	require.Equal(t, 3, len(*result.Hits), "number of hits is invalid")

	docs := make([]map[string]interface{}, len(*result.Hits))
	for i, hit := range *result.Hits {
		docs[i] = *hit.Document
	}

	require.Equal(t, expectedDocs, docs)
}
