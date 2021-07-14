// +build integration

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func newUUIDName(namePrefix string) string {
	nameUUID := uuid.New()
	return fmt.Sprintf("%s_%s", namePrefix, nameUUID.String())
}

func newSchema(collectionName string) *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name:  "company_name",
				Type:  "string",
				Index: pointer.True(),
			},
			{
				Name:  "num_employees",
				Type:  "int32",
				Index: pointer.True(),
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: true,
				Index: pointer.True(),
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

func withNumEmployees(numEmployees int) newDocumentOption {
	return func(doc *testDocument) {
		doc.NumEmployees = numEmployees
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

func withResponseNumEmployees(numEmployees int) newDocumentResponseOption {
	return func(doc map[string]interface{}) {
		doc["num_employees"] = float64(numEmployees)
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

type newSearchOverrideSchemaOption func(*api.SearchOverrideSchema)

func withOverrideRuleMatch(match string) newSearchOverrideSchemaOption {
	return func(o *api.SearchOverrideSchema) {
		o.Rule.Match = match
	}
}

func newSearchOverrideSchema(opts ...newSearchOverrideSchemaOption) *api.SearchOverrideSchema {
	schema := &api.SearchOverrideSchema{
		Rule: api.SearchOverrideRule{
			Query: "apple",
			Match: "exact",
		},
		Includes: []api.SearchOverrideInclude{
			{
				Id:       "422",
				Position: 1,
			},
			{
				Id:       "54",
				Position: 2,
			},
		},
		Excludes: []api.SearchOverrideExclude{
			{
				Id: "287",
			},
		},
	}
	for _, opt := range opts {
		opt(schema)
	}
	return schema
}

func newSearchOverride(overrideID string, opts ...newSearchOverrideSchemaOption) *api.SearchOverride {
	return &api.SearchOverride{
		SearchOverrideSchema: *newSearchOverrideSchema(opts...),
		Id:                   overrideID,
	}
}

type newSynonymOption func(*api.SearchSynonymSchema)

func withSynonyms(synonyms ...string) newSynonymOption {
	return func(s *api.SearchSynonymSchema) {
		s.Synonyms = synonyms
	}
}

func newSearchSynonymSchema(opts ...newSynonymOption) *api.SearchSynonymSchema {
	schema := &api.SearchSynonymSchema{
		Synonyms: []string{"blazer", "coat", "jacket"},
	}
	for _, opt := range opts {
		opt(schema)
	}
	return schema
}

func newSearchSynonym(synonymID string, opts ...newSynonymOption) *api.SearchSynonym {
	return &api.SearchSynonym{
		SearchSynonymSchema: *newSearchSynonymSchema(opts...),
		Id:                  synonymID,
	}
}

func newCollectionAlias(collectionName string, name string) *api.CollectionAlias {
	return &api.CollectionAlias{
		CollectionName: collectionName,
		Name:           name,
	}
}

func createNewCollection(t *testing.T, namePrefix string) string {
	t.Helper()
	collectionName := newUUIDName(namePrefix)
	schema := newSchema(collectionName)

	_, err := typesenseClient.Collections().Create(schema)
	require.NoError(t, err)
	return collectionName
}

func createDocument(t *testing.T, collectionName string, document *testDocument) {
	t.Helper()
	_, err := typesenseClient.Collection(collectionName).Documents().Create(document)
	require.NoError(t, err)
}

func createNewKey(t *testing.T) *api.ApiKey {
	t.Helper()
	keySchema := newKeySchema()

	result, err := typesenseClient.Keys().Create(keySchema)

	require.NoError(t, err)
	return result
}

func retrieveDocuments(t *testing.T, collectionName string, docIDs ...string) []map[string]interface{} {
	results := make([]map[string]interface{}, len(docIDs))
	for i, docID := range docIDs {
		doc, err := typesenseClient.Collection(collectionName).Document(docID).Retrieve()
		require.NoError(t, err)
		results[i] = doc
	}
	return results
}
