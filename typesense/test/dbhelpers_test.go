//go:build integration
// +build integration

package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"github.com/typesense/typesense-go/v3/typesense/api/pointer"
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
		Metadata: &map[string]interface{}{
			"revision": "1",
		},
	}
}

func expectedNewCollection(name string) *api.CollectionResponse {
	return &api.CollectionResponse{
		Name: name,
		Fields: []api.Field{
			{
				Name:           "company_name",
				Type:           "string",
				Facet:          pointer.False(),
				Optional:       pointer.False(),
				Index:          pointer.True(),
				Infix:          pointer.False(),
				Locale:         pointer.String(""),
				Sort:           pointer.False(),
				Drop:           nil,
				Store:          pointer.True(),
				Stem:           pointer.False(),
				StemDictionary: pointer.String(""),
			},
			{
				Name:           "num_employees",
				Type:           "int32",
				Facet:          pointer.False(),
				Optional:       pointer.False(),
				Index:          pointer.True(),
				Infix:          pointer.False(),
				Locale:         pointer.String(""),
				Sort:           pointer.True(),
				Drop:           nil,
				Store:          pointer.True(),
				Stem:           pointer.False(),
				StemDictionary: pointer.String(""),
			},
			{
				Name:           "country",
				Type:           "string",
				Facet:          pointer.True(),
				Optional:       pointer.True(),
				Index:          pointer.True(),
				Infix:          pointer.False(),
				Locale:         pointer.String(""),
				Sort:           pointer.False(),
				Drop:           nil,
				Store:          pointer.True(),
				Stem:           pointer.False(),
				StemDictionary: pointer.String(""),
			},
		},
		EnableNestedFields:  pointer.False(),
		DefaultSortingField: pointer.String(""),
		TokenSeparators:     &[]string{},
		SymbolsToIndex:      &[]string{},
		NumDocuments:        pointer.Int64(0),
		CreatedAt:           pointer.Int64(0),
		Metadata: &map[string]interface{}{
			"revision": "1",
		},
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

func newSearchOverrideSchema() *api.SearchOverrideSchema {
	schema := &api.SearchOverrideSchema{
		Rule: api.SearchOverrideRule{
			Query: pointer.String("apple"),
			Match: pointer.Any(api.Exact),
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
		FilterCuratedHits:   pointer.False(),
		StopProcessing:      pointer.True(),
	}

	return schema
}

func newSearchOverride(overrideID string) *api.SearchOverride {
	return &api.SearchOverride{
		Id: pointer.String(overrideID),
		Rule: api.SearchOverrideRule{
			Query: pointer.String("apple"),
			Match: pointer.Any(api.Exact),
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
		FilterCuratedHits:   pointer.False(),
		StopProcessing:      pointer.True(),
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

func newPresetFromSearchParametersUpsertSchema() *api.PresetUpsertSchema {
	preset := &api.PresetUpsertSchema{}
	preset.Value.FromSearchParameters(api.SearchParameters{
		Q: pointer.Any("hello"),
	})
	return preset
}

func newPresetFromSearchParameters(presetName string) *api.PresetSchema {
	preset := &api.PresetSchema{
		Name: presetName,
	}
	preset.Value.FromSearchParameters(api.SearchParameters{
		Q: pointer.Any("hello"),
	})
	return preset
}

func newPresetFromMultiSearchSearchesParameterUpsertSchema() *api.PresetUpsertSchema {
	preset := &api.PresetUpsertSchema{}
	preset.Value.FromMultiSearchSearchesParameter(api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: pointer.Any("test"),
			},
		},
	})
	return preset
}

func newPresetFromMultiSearchSearchesParameter(presetName string) *api.PresetSchema {
	preset := &api.PresetSchema{
		Name: presetName,
	}
	preset.Value.FromMultiSearchSearchesParameter(api.MultiSearchSearchesParameter{
		Searches: []api.MultiSearchCollectionParameters{
			{
				Collection: pointer.Any("test"),
			},
		},
	})
	return preset
}

func newAnalyticsRuleUpsertSchema(collectionName string, sourceCollectionName string, eventName string) *api.AnalyticsRuleUpsertSchema {
	return &api.AnalyticsRuleUpsertSchema{
		Type: "counter",
		Params: api.AnalyticsRuleParameters{
			Source: api.AnalyticsRuleParametersSource{
				Collections: []string{sourceCollectionName},
				Events: &[]struct {
					Name   string  "json:\"name\""
					Type   string  "json:\"type\""
					Weight float32 "json:\"weight\""
				}{
					{Type: "click", Weight: 1, Name: eventName},
				},
			},
			Destination: api.AnalyticsRuleParametersDestination{
				Collection:   collectionName,
				CounterField: pointer.String("num_employees"),
			},
			Limit: pointer.Int(9999),
		},
	}
}

func newAnalyticsRule(ruleName string, collectionName string, sourceCollectionName string, eventName string) *api.AnalyticsRuleSchema {
	return &api.AnalyticsRuleSchema{
		Name: ruleName,
		Type: "counter",
		Params: api.AnalyticsRuleParameters{
			Source: api.AnalyticsRuleParametersSource{
				Collections: []string{sourceCollectionName},
				Events: &[]struct {
					Name   string  "json:\"name\""
					Type   string  "json:\"type\""
					Weight float32 "json:\"weight\""
				}{
					{Type: "click", Weight: 1, Name: eventName},
				},
			},
			Destination: api.AnalyticsRuleParametersDestination{
				Collection:   collectionName,
				CounterField: pointer.String("num_employees"),
			},
			Limit: pointer.Int(9999),
		},
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
	_, err := typesenseClient.Collection(collectionName).Documents().Create(context.Background(), document, &api.DocumentIndexParameters{})
	require.NoError(t, err)
}

func createNewKey(t *testing.T) *api.ApiKey {
	t.Helper()
	keySchema := newKeySchema()

	result, err := typesenseClient.Keys().Create(context.Background(), keySchema)

	require.NoError(t, err)
	return result
}

func createNewPreset(t *testing.T, presetValueIsFromSearchParameters ...bool) (string, *api.PresetSchema) {
	t.Helper()
	presetName := newUUIDName("preset-test")
	presetSchema := newPresetFromMultiSearchSearchesParameterUpsertSchema()

	if len(presetValueIsFromSearchParameters) > 0 {
		presetSchema = newPresetFromSearchParametersUpsertSchema()
	}

	result, err := typesenseClient.Presets().Upsert(context.Background(), presetName, presetSchema)

	require.NoError(t, err)
	return presetName, result
}

func createNewAnalyticsRule(t *testing.T, collectionName string, sourceCollectionName string, eventName string) *api.AnalyticsRuleSchema {
	t.Helper()
	ruleSchema := newAnalyticsRuleUpsertSchema(collectionName, sourceCollectionName, eventName)
	ruleName := newUUIDName("test-rule")

	result, err := typesenseClient.Analytics().Rules().Upsert(context.Background(), ruleName, ruleSchema)

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

func newNLSearchModelCreateSchema() *api.NLSearchModelCreateSchema {
	apiKey := os.Getenv("NL_SEARCH_MODEL_API_KEY")
	
	return &api.NLSearchModelCreateSchema{
		ModelName:    pointer.String("openai/gpt-3.5-turbo"),
		ApiKey:       pointer.String(apiKey),
		MaxBytes:     pointer.Int(1000),
		Temperature:  pointer.Float32(0.7),
		SystemPrompt: pointer.String("You are a helpful assistant."),
		TopP:         pointer.Float32(0.9),
		TopK:         pointer.Int(40),
		StopSequences: &[]string{"END", "STOP"},
		ApiVersion:   pointer.String("v1"),
	}
}

func newNLSearchModelSchema(modelID string) *api.NLSearchModelSchema {
	apiKey := os.Getenv("NL_SEARCH_MODEL_API_KEY")
	
	return &api.NLSearchModelSchema{
		Id:           modelID,
		ModelName:    pointer.String("openai/gpt-3.5-turbo"),
		ApiKey:       pointer.String(apiKey),
		MaxBytes:     pointer.Int(1000),
		Temperature:  pointer.Float32(0.7),
		SystemPrompt: pointer.String("You are a helpful assistant."),
		TopP:         pointer.Float32(0.9),
		TopK:         pointer.Int(40),
		StopSequences: &[]string{"END", "STOP"},
		ApiVersion:   pointer.String("v1"),
	}
}

func newNLSearchModelUpdateSchema() *api.NLSearchModelUpdateSchema {
	apiKey := os.Getenv("NL_SEARCH_MODEL_API_KEY")
	
	return &api.NLSearchModelUpdateSchema{
		ModelName:    pointer.String("openai/gpt-4"),
		ApiKey:       pointer.String(apiKey),
		MaxBytes:     pointer.Int(2000),
		Temperature:  pointer.Float32(0.5),
		SystemPrompt: pointer.String("You are an expert assistant."),
		TopP:         pointer.Float32(0.8),
		TopK:         pointer.Int(50),
		StopSequences: &[]string{"END", "STOP", "QUIT"},
		ApiVersion:   pointer.String("v1"),
	}
}

func shouldSkipNLSearchModelTests(t *testing.T)  {
	if os.Getenv("NL_SEARCH_MODEL_API_KEY") == "" {
		t.Skip("Skipping NL search model test: NL_SEARCH_MODEL_API_KEY not set")
	}
}

func createNewNLSearchModel(t *testing.T) (string, *api.NLSearchModelSchema) {
	t.Helper()
	modelID := newUUIDName("nl-model-test")
	modelSchema := newNLSearchModelCreateSchema()
	modelSchema.Id = pointer.String(modelID)

	result, err := typesenseClient.NLSearchModels().Create(context.Background(), modelSchema)

	require.NoError(t, err)
	return modelID, result
}
