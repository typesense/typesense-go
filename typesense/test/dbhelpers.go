// +build integration

package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
)

func newCollectionName(namePrefix string) string {
	nameUUID := uuid.New()
	return fmt.Sprintf("%s_%s", namePrefix, nameUUID.String())
}

func newSchema(collectionName string) *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name: "company_name",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: true,
			},
		},
		DefaultSortingField: "num_employees",
	}
}

func expectedNewCollection(name string) *api.Collection {
	return &api.Collection{
		CollectionSchema: *newSchema(name),
		NumDocuments:     0,
	}
}

type testDocument struct {
	ID           string `json:"id"`
	CompanyName  string `json:"company_name"`
	NumEmployees int    `json:"num_employees"`
	Country      string `json:"country"`
}

func newDocument(docID string) *testDocument {
	return &testDocument{
		ID:           docID,
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}
}

func newDocumentResponse(docID string) map[string]interface{} {
	document := map[string]interface{}{}
	document["id"] = docID
	document["company_name"] = "Stark Industries"
	document["num_employees"] = float64(5215)
	document["country"] = "USA"
	return document
}

func createNewCollection(t *testing.T, namePrefix string) string {
	t.Helper()
	collectionName := newCollectionName(namePrefix)
	schema := newSchema(collectionName)

	_, err := typesenseClient.Collections().Create(schema)
	require.NoError(t, err)
	return collectionName
}

func createNewDocument(t *testing.T, collectionName string, docID string) *testDocument {
	document := newDocument(docID)
	_, err := typesenseClient.Collection(collectionName).Documents().Create(document)
	require.NoError(t, err)
	return document
}
