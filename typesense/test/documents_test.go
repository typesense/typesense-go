//go:build integration
// +build integration

package test

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestDocumentCreate(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")

	document := newDocument("123")
	result, err := typesenseClient.Collection(collectionName).Documents().Create(context.Background(), document, &api.DocumentIndexParameters{})

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentUpsertNewDocument(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")

	document := newDocument("123")
	result, err := typesenseClient.Collection(collectionName).Documents().Upsert(context.Background(), document, &api.DocumentIndexParameters{})

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentUpsertExistingDocument(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	newCompanyName := "HighTech Inc."
	expectedResult := newDocumentResponse("123", withResponseCompanyName(newCompanyName))

	document := newDocument("123")
	_, err := typesenseClient.Collection(collectionName).Documents().Create(context.Background(), document, &api.DocumentIndexParameters{})
	require.NoError(t, err)

	document.CompanyName = newCompanyName

	result, err := typesenseClient.Collection(collectionName).Documents().Upsert(context.Background(), document, &api.DocumentIndexParameters{})

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentsDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")

	document := newDocument("123")
	document.NumEmployees = 5000
	_, err := typesenseClient.Collection(collectionName).Documents().Create(context.Background(), document, &api.DocumentIndexParameters{})
	require.NoError(t, err)

	document = newDocument("124")
	document.NumEmployees = 7000
	_, err = typesenseClient.Collection(collectionName).Documents().Create(context.Background(), document, &api.DocumentIndexParameters{})
	require.NoError(t, err)

	filter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>6500"), BatchSize: pointer.Int(100)}
	result, err := typesenseClient.Collection(collectionName).Documents().Delete(context.Background(), filter)

	require.NoError(t, err)
	require.Equal(t, 1, result)

	_, err = typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())
	require.NoError(t, err)
	_, err = typesenseClient.Collection(collectionName).Document("124").Retrieve(context.Background())
	require.Error(t, err)
}

func TestDocumentsExport(t *testing.T) {
	collectionName := createNewCollection(t, "companies")

	expectedResults := []map[string]interface{}{
		newDocumentResponse("123"),
		newDocumentResponse("125", withResponseCompanyName("Company2")),
		newDocumentResponse("127", withResponseCompanyName("Company3")),
	}

	createDocument(t, collectionName, newDocument("123"))
	createDocument(t, collectionName, newDocument("125", withCompanyName("Company2")))
	createDocument(t, collectionName, newDocument("127", withCompanyName("Company3")))

	body, err := typesenseClient.Collection(collectionName).Documents().Export(context.Background(), &api.ExportDocumentsParams{})
	require.NoError(t, err)
	defer body.Close()

	jd := json.NewDecoder(body)
	results := make([]map[string]interface{}, 3)
	for i := 0; i < 3; i++ {
		require.True(t, jd.More(), "no json element")
		doc := map[string]interface{}{}
		require.NoError(t, jd.Decode(&doc))
		results[i] = doc
	}
	sort.Slice(results, func(i, j int) bool {
		id1 := results[i]["id"].(string)
		id2 := results[j]["id"].(string)
		return id1 < id2
	})

	require.Equal(t, expectedResults, results)
}
