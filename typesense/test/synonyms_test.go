//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"github.com/typesense/typesense-go/v2/typesense/api/pointer"
)

func TestSearchSynonymUpsertNewSynonym(t *testing.T) {
	collectionName := createNewCollection(t, "product")
	synonymID := newUUIDName("customize-apple")
	expectedResult := newSearchSynonym(synonymID)

	body := newSearchSynonymSchema()
	result, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(context.Background(), synonymID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Synonym(synonymID).Retrieve(context.Background())

	require.NoError(t, err)
	expectedResult.Root = pointer.String("")
	require.Equal(t, expectedResult, result)
}

func TestSearchSynonymUpsertExistingSynonym(t *testing.T) {
	collectionName := createNewCollection(t, "product")
	synonymID := newUUIDName("customize-apple")
	expectedResult := newSearchSynonym(synonymID)
	expectedResult.Synonyms = []string{"blazer", "coat", "jacket"}

	body := newSearchSynonymSchema(withSynonyms("blazer", "coat"))
	_, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(context.Background(), synonymID, body)
	require.NoError(t, err)

	body.Synonyms = []string{"blazer", "coat", "jacket"}

	result, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(context.Background(), synonymID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Synonym(synonymID).Retrieve(context.Background())

	require.NoError(t, err)
	expectedResult.Root = pointer.String("")
	require.Equal(t, expectedResult, result)
}

func TestSearchSynonymsRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "products")
	total := 3
	synonymIDs := make([]string, total)
	for i := 0; i < total; i++ {
		synonymIDs[i] = newUUIDName("customize-apple")
	}
	schema := newSearchSynonymSchema()
	expectedResult := map[string]*api.SearchSynonym{}
	for i := 0; i < total; i++ {
		expectedResult[synonymIDs[i]] = newSearchSynonym(synonymIDs[i])
	}

	for i := 0; i < total; i++ {
		_, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(context.Background(), synonymIDs[i], schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Collection(collectionName).Synonyms().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of overrides is invalid")

	resultMap := map[string]*api.SearchSynonym{}
	for _, synonym := range result {
		resultMap[*synonym.Id] = synonym
	}

	for k, v := range expectedResult {
		v.Root = pointer.String("")
		assert.Equal(t, v, resultMap[k])
	}
}
