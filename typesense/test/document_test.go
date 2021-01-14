// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")
	createDocument(t, collectionName, newDocument("123"))

	result, err := typesenseClient.Collection(collectionName).Document("123").Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentUpdate(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	newCompanyName := "HighTech Inc."
	expectedResult := newDocumentResponse("123", withResponseCompanyName(newCompanyName))

	document := newDocument("123")
	createDocument(t, collectionName, document)

	document.CompanyName = newCompanyName
	typesenseClient.Collection(collectionName).Document("123").Update(document)

	result, err := typesenseClient.Collection(collectionName).Document("123").Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")
	createDocument(t, collectionName, newDocument("123"))

	result, err := typesenseClient.Collection(collectionName).Document("123").Delete()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	_, err = typesenseClient.Collection(collectionName).Document("123").Retrieve()
	require.Error(t, err)
}
