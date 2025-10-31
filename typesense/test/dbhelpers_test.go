//go:build integration
// +build integration

package test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func isV30OrAbove(t *testing.T) bool {
	t.Helper()

	debug, err := typesenseClient.Debug(context.Background())
	if err != nil {
		t.Logf("Failed to get debug info: %v", err)
		return false
	}

	if debug.JSON200 == nil || debug.JSON200.Version == nil {
		t.Log("Debug response or version is nil")
		return false
	}

	version := *debug.JSON200.Version
	if version == "nightly" {
		return true
	}

	var numberedVersion string
	if strings.HasPrefix(version, "v") {
		numberedVersion = strings.Split(version, "v")[1]
	} else {
		numberedVersion = version
	}
	parts := strings.Split(numberedVersion, ".")
	if len(parts) == 0 {
		t.Logf("Version parts empty: %s", numberedVersion)
		return false
	}

	majorVersion, err := strconv.Atoi(parts[0])
	if err != nil {
		t.Logf("Failed to parse major version: %v", err)
		return false
	}

	return majorVersion >= 30
}

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

func expectedNewCollection(t *testing.T, name string) *api.CollectionResponse {
	if !isV30OrAbove(t) {
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
		}
	}
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
		SynonymSets:         &[]string{},
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



func newCurationSetCreateSchema() *api.CurationSetCreateSchema {
	return &api.CurationSetCreateSchema{
		Items: []api.CurationItemCreateSchema{
			{
				Id: pointer.String("dummy"),
				Rule: api.CurationRule{
					Query: pointer.String("apple"),
					Match: pointer.Any(api.Exact),
				},
				Includes: &[]api.CurationInclude{
					{
						Id: "422",
					},
					{
						Id: "54",
					},
				},
				Excludes: &[]api.CurationExclude{
					{
						Id: "287",
					},
				},
				RemoveMatchedTokens: pointer.True(),
				FilterBy:            pointer.String("category:=Electronics"),
				StopProcessing:      pointer.True(),
			},
		},
		Description: pointer.String("Test curation set"),
	}
}

func newCurationSetSchema(curationSetName string) *api.CurationSetSchema {
	return &api.CurationSetSchema{
		Name: curationSetName,
		Items: []api.CurationItemCreateSchema{
			{
				Id: pointer.String("dummy"),
				Rule: api.CurationRule{
					Query: pointer.String("apple"),
					Match: pointer.Any(api.Exact),
				},
				Includes: &[]api.CurationInclude{
					{
						Id: "422",
					},
					{
						Id: "54",
					},
				},
				Excludes: &[]api.CurationExclude{
					{
						Id: "287",
					},
				},
				RemoveMatchedTokens: pointer.True(),
				FilterBy:            pointer.String("category:=Electronics"),
				StopProcessing:      pointer.True(),
			},
		},
		Description: pointer.String("Test curation set"),
	}
}

func newSynonymSetCreateSchema() *api.SynonymSetCreateSchema {
	return &api.SynonymSetCreateSchema{
		Items: []api.SynonymItemSchema{
			{
				Id:       "dummy",
				Synonyms: []string{"foo", "bar", "baz"},
			},
		},
	}
}

func newSynonymSetSchema(synonymSetName string) *api.SynonymSetSchema {
	return &api.SynonymSetSchema{
		Name: synonymSetName,
		Items: []api.SynonymItemSchema{
			{
				Id:       "dummy",
				Synonyms: []string{"foo", "bar", "baz"},
			},
		},
	}
}

func createNewSynonymSet(t *testing.T) (string, *api.SynonymSetSchema) {
	t.Helper()
	synonymSetName := newUUIDName("synonym-set-test")
	synonymSetSchema := newSynonymSetCreateSchema()

	result, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, synonymSetSchema)

	require.NoError(t, err)
	return synonymSetName, result
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

func newAnalyticsRule(ruleName string, collectionName string, sourceCollectionName string, eventName string) *api.AnalyticsRule {
	return &api.AnalyticsRule{
		Name:       ruleName,
		Type:       api.AnalyticsRuleTypeCounter,
		Collection: collectionName,
		EventType:  "click",
		Params: &api.AnalyticsRuleCreateParams{
			CounterField: pointer.String("num_employees"),
			Weight:       pointer.Int(1),
		},
	}
}

func createNewAnalyticsRule(t *testing.T, collectionName string, sourceCollectionName string, eventName string) *api.AnalyticsRule {
	t.Helper()
	ruleName := newUUIDName("test-rule")

	// Create the rule using the new API
	ruleCreate := &api.AnalyticsRuleCreate{
		Name:       ruleName,
		Type:       api.AnalyticsRuleCreateTypeCounter,
		Collection: collectionName,
		EventType:  "click",
		Params: &api.AnalyticsRuleCreateParams{
			CounterField: pointer.String("num_employees"),
			Weight:       pointer.Int(1),
		},
	}

	// Create the rule via the API
	_, err := typesenseClient.Analytics().Rules().Create(context.Background(), []*api.AnalyticsRuleCreate{ruleCreate})
	require.NoError(t, err)

	// Return the expected rule structure
	return newAnalyticsRule(ruleName, collectionName, sourceCollectionName, eventName)
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
		ModelName:     pointer.String("openai/gpt-3.5-turbo"),
		ApiKey:        pointer.String(apiKey),
		MaxBytes:      pointer.Int(1000),
		Temperature:   pointer.Float32(0.7),
		SystemPrompt:  pointer.String("You are a helpful assistant."),
		TopP:          pointer.Float32(0.9),
		TopK:          pointer.Int(40),
		StopSequences: &[]string{"END", "STOP"},
		ApiVersion:    pointer.String("v1"),
	}
}

func newNLSearchModelSchema(modelID string) *api.NLSearchModelSchema {
	apiKey := os.Getenv("NL_SEARCH_MODEL_API_KEY")

	return &api.NLSearchModelSchema{
		Id:            modelID,
		ModelName:     pointer.String("openai/gpt-3.5-turbo"),
		ApiKey:        pointer.String(apiKey),
		MaxBytes:      pointer.Int(1000),
		Temperature:   pointer.Float32(0.7),
		SystemPrompt:  pointer.String("You are a helpful assistant."),
		TopP:          pointer.Float32(0.9),
		TopK:          pointer.Int(40),
		StopSequences: &[]string{"END", "STOP"},
		ApiVersion:    pointer.String("v1"),
	}
}

func newNLSearchModelUpdateSchema() *api.NLSearchModelUpdateSchema {
	apiKey := os.Getenv("NL_SEARCH_MODEL_API_KEY")

	return &api.NLSearchModelUpdateSchema{
		ModelName:     pointer.String("openai/gpt-4"),
		ApiKey:        pointer.String(apiKey),
		MaxBytes:      pointer.Int(2000),
		Temperature:   pointer.Float32(0.5),
		SystemPrompt:  pointer.String("You are an expert assistant."),
		TopP:          pointer.Float32(0.8),
		TopK:          pointer.Int(50),
		StopSequences: &[]string{"END", "STOP", "QUIT"},
		ApiVersion:    pointer.String("v1"),
	}
}

func shouldSkipNLSearchModelTests(t *testing.T) {
	if os.Getenv("NL_SEARCH_MODEL_API_KEY") == "" {

		t.Skip("Skipping NL search model test: NL_SEARCH_MODEL_API_KEY not set")
	}
}

func shouldSkipAnalyticsTests(t *testing.T) {
	if !isV30OrAbove(t) {
		t.Skip("Skipping analytics tests: requires Typesense v30 or above")
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
