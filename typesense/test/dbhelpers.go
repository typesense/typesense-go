// +build integration

package test

import (
	"fmt"
	"testing"
	"time"

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

type newDocumentOption func(*testDocument)

func withCompanyName(companyName string) newDocumentOption {
	return func(doc *testDocument) {
		doc.CompanyName = companyName
	}
}

func newDocument(docID string, opts ...newDocumentOption) *testDocument {
	doc := &testDocument{
		ID:           docID,
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}
	for _, opt := range opts {
		opt(doc)
	}
	return doc
}

type newDocumentResponseOption func(map[string]interface{})

func withResponseCompanyName(companyName string) newDocumentResponseOption {
	return func(doc map[string]interface{}) {
		doc["company_name"] = companyName
	}
}

func newDocumentResponse(docID string, opts ...newDocumentResponseOption) map[string]interface{} {
	document := map[string]interface{}{}
	document["id"] = docID
	document["company_name"] = "Stark Industries"
	document["num_employees"] = float64(5215)
	document["country"] = "USA"
	for _, opt := range opts {
		opt(document)
	}
	return document
}

func newKeySchema() *api.ApiKeySchema {
	return &api.ApiKeySchema{
		Description: "Search-only key.",
		Actions:     []string{"documents:search"},
		Collections: []string{"*"},
		ExpiresAt:   time.Now().Add(1 * time.Hour).Unix(),
	}
}

func newKey() *api.ApiKey {
	return &api.ApiKey{
		ApiKeySchema: *newKeySchema(),
	}
}

func createNewCollection(t *testing.T, namePrefix string) string {
	t.Helper()
	collectionName := newCollectionName(namePrefix)
	schema := newSchema(collectionName)

	_, err := typesenseClient.Collections().Create(schema)
	require.NoError(t, err)
	return collectionName
}

func createDocument(t *testing.T, collectionName string, document *testDocument) {
	_, err := typesenseClient.Collection(collectionName).Documents().Create(document)
	require.NoError(t, err)
}

func createNewKey(t *testing.T) *api.ApiKey {
	keySchema := newKeySchema()

	result, err := typesenseClient.Keys().Create(keySchema)

	require.NoError(t, err)
	return result
}
