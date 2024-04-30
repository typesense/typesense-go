//go:build integration
// +build integration

package test

import (
	"context"
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
				Name: "company_name",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:     "country",
				Type:     "string",
				Facet:    pointer.True(),
				Optional: pointer.True(),
			},
		},
	}
}

func expectedNewCollection(name string) *api.CollectionResponse {
	return &api.CollectionResponse{
		Name: name,
		Fields: []api.Field{
			{
				Name:     "company_name",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
				Index:    pointer.True(),
				Infix:    pointer.False(),
				Locale:   pointer.String(""),
				Sort:     pointer.False(),
				Drop:     nil,
			},
			{
				Name:     "num_employees",
				Type:     "int32",
				Facet:    pointer.False(),
				Optional: pointer.False(),
				Index:    pointer.True(),
				Infix:    pointer.False(),
				Locale:   pointer.String(""),
				Sort:     pointer.True(),
				Drop:     nil,
			},
			{
				Name:     "country",
				Type:     "string",
				Facet:    pointer.True(),
				Optional: pointer.True(),
				Index:    pointer.True(),
				Infix:    pointer.False(),
				Locale:   pointer.String(""),
				Sort:     pointer.False(),
				Drop:     nil,
			},
		},
		EnableNestedFields:  pointer.False(),
		DefaultSortingField: pointer.String(""),
		TokenSeparators:     &[]string{},
		SymbolsToIndex:      &[]string{},
		NumDocuments:        pointer.Int64(0),
		CreatedAt:           pointer.Int64(0),
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

func withCountry(country string) newDocumentOption {
	return func(doc *testDocument) {
		doc.Country = country
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

func newStructResponse(docID string, opts ...func(*testDocument)) *testDocument {
	document := &testDocument{}
	document.ID = docID
	document.CompanyName = "Stark Industries"
	document.NumEmployees = 5215
	document.Country = "USA"
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
		ExpiresAt:   pointer.Int64(time.Now().Add(1 * time.Hour).Unix()),
	}
}

func newKey() *api.ApiKey {
	return &api.ApiKey{
		Description: "Search-only key.",
		Actions:     []string{"documents:search"},
		Collections: []string{"*"},
		ExpiresAt:   pointer.Int64(time.Now().Add(1 * time.Hour).Unix()),
	}
}

type newSearchOverrideSchemaOption func(*api.SearchOverrideSchema)

func withOverrideRuleMatch(match api.SearchOverrideRuleMatch) newSearchOverrideSchemaOption {
	return func(o *api.SearchOverrideSchema) {
		o.Rule.Match = match
	}
}

func newSearchOverrideSchema() *api.SearchOverrideSchema {
	schema := &api.SearchOverrideSchema{
		Rule: api.SearchOverrideRule{
			Query: "apple",
			Match: "exact",
		},
		Includes: &[]api.SearchOverrideInclude{
			{
				Id:       "422",
				Position: 1,
			},
			{
				Id:       "54",
				Position: 2,
			},
		},
		Excludes: &[]api.SearchOverrideExclude{
			{
				Id: "287",
			},
		},
		RemoveMatchedTokens: pointer.True(),
	}

	return schema
}

func newSearchOverride(overrideID string) *api.SearchOverride {
	return &api.SearchOverride{
		Id: pointer.String(overrideID),
		Rule: api.SearchOverrideRule{
			Query: "apple",
			Match: "exact",
		},
		Includes: &[]api.SearchOverrideInclude{
			{
				Id:       "422",
				Position: 1,
			},
			{
				Id:       "54",
				Position: 2,
			},
		},
		Excludes: &[]api.SearchOverrideExclude{
			{
				Id: "287",
			},
		},
		RemoveMatchedTokens: pointer.True(),
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

func newSearchSynonym(synonymID string) *api.SearchSynonym {
	return &api.SearchSynonym{
		Id:       pointer.String(synonymID),
		Synonyms: []string{"blazer", "coat", "jacket"},
	}
}

func newCollectionAlias(collectionName string, name string) *api.CollectionAlias {
	return &api.CollectionAlias{
		CollectionName: collectionName,
		Name:           pointer.String(name),
	}
}

func createNewCollection(t *testing.T, namePrefix string) string {
	t.Helper()
	collectionName := newUUIDName(namePrefix)
	schema := newSchema(collectionName)

	_, err := typesenseClient.Collections().Create(context.Background(), schema)
	require.NoError(t, err)
	return collectionName
}

func createDocument(t *testing.T, collectionName string, document *testDocument) {
	t.Helper()
	_, err := typesenseClient.Collection(collectionName).Documents().Create(context.Background(), document)
	require.NoError(t, err)
}

func createNewKey(t *testing.T) *api.ApiKey {
	t.Helper()
	keySchema := newKeySchema()

	result, err := typesenseClient.Keys().Create(context.Background(), keySchema)

	require.NoError(t, err)
	return result
}

func retrieveDocuments(t *testing.T, collectionName string, docIDs ...string) []map[string]interface{} {
	results := make([]map[string]interface{}, len(docIDs))
	for i, docID := range docIDs {
		doc, err := typesenseClient.Collection(collectionName).Document(docID).Retrieve(context.Background())
		require.NoError(t, err)
		results[i] = doc
	}
	return results
}
