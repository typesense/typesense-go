// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestSearchSynonymUpsertNewSynonym(t *testing.T) {
	collectionName := createNewCollection(t, "product")
	synonymID := newUUIDName("customize-apple")
	expectedResult := newSearchSynonym(synonymID)

	body := newSearchSynonymSchema()
	result, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(synonymID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Synonym(synonymID).Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSearchSynonymUpsertExistingSynonym(t *testing.T) {
	collectionName := createNewCollection(t, "product")
	synonymID := newUUIDName("customize-apple")
	expectedResult := newSearchSynonym(synonymID, withSynonyms("blazer", "coat", "jacket"))

	body := newSearchSynonymSchema(withSynonyms("blazer", "coat"))
	_, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(synonymID, body)
	require.NoError(t, err)

	body.Synonyms = []string{"blazer", "coat", "jacket"}

	result, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(synonymID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Synonym(synonymID).Retrieve()

	require.NoError(t, err)
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
		_, err := typesenseClient.Collection(collectionName).Synonyms().Upsert(synonymIDs[i], schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Collection(collectionName).Synonyms().Retrieve()

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of overrides is invalid")

	resultMap := map[string]*api.SearchSynonym{}
	for _, synonym := range result {
		resultMap[synonym.Id] = synonym
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
