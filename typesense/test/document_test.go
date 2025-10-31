//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func TestDocumentRetrieveGeneric(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newStructResponse("123")
	createDocument(t, collectionName, newDocument("123"))
	result, err := typesense.GenericCollection[*testDocument](typesenseClient, collectionName).Document("123").Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")
	createDocument(t, collectionName, newDocument("123"))

	result, err := typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())

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
	typesenseClient.Collection(collectionName).Document("123").Update(context.Background(), document, &api.DocumentIndexParameters{})

	result, err := typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestDocumentDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := newDocumentResponse("123")
	createDocument(t, collectionName, newDocument("123"))

	result, err := typesenseClient.Collection(collectionName).Document("123").Delete(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	_, err = typesenseClient.Collection(collectionName).Document("123").Retrieve(context.Background())
	require.Error(t, err)
}
