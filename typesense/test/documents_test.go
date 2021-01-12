// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestDocumentCreate(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")

	document := newDocument("123")
	result, err := typesenseClient.Collection(collectionName).Documents().Create(document)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Document("123").Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentUpsertNewDocument(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")

	document := newDocument("123")
	result, err := typesenseClient.Collection(collectionName).Documents().Upsert(document)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Document("123").Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentUpsertExistingDocument(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")
	newCompanyName := "HighTech Inc."
	expectedResult["company_name"] = newCompanyName

	document := newDocument("123")
	_, err := typesenseClient.Collection(collectionName).Documents().Create(document)
	require.NoError(t, err)

	document.CompanyName = newCompanyName

	result, err := typesenseClient.Collection(collectionName).Documents().Upsert(document)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Document("123").Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentsDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")

	document := newDocument("123")
	document.NumEmployees = 5000
	_, err := typesenseClient.Collection(collectionName).Documents().Create(document)
	require.NoError(t, err)

	document = newDocument("124")
	document.NumEmployees = 7000
	_, err = typesenseClient.Collection(collectionName).Documents().Create(document)
	require.NoError(t, err)

	filter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>6500", BatchSize: 100}
	result, err := typesenseClient.Collection(collectionName).Documents().Delete(filter)

	require.NoError(t, err)
	require.Equal(t, 1, result)

	_, err = typesenseClient.Collection(collectionName).Document("123").Retrieve()
	require.NoError(t, err)
	_, err = typesenseClient.Collection(collectionName).Document("124").Retrieve()
	require.Error(t, err)
}
