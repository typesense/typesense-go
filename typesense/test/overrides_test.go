// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestSearchOverrideUpsertNewOverride(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)

	body := newSearchOverrideSchema()
	result, err := typesenseClient.Collection(collectionName).Overrides().Upsert(overrideID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Override(overrideID).Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSearchOverrideUpsertExistingOverride(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID, withOverrideRuleMatch("contains"))

	body := newSearchOverrideSchema(withOverrideRuleMatch("exact"))
	_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(overrideID, body)
	require.NoError(t, err)

	body.Rule.Match = "contains"

	result, err := typesenseClient.Collection(collectionName).Overrides().Upsert(overrideID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Override(overrideID).Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSearchOverridesRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	total := 3
	overrideIDs := make([]string, total)
	for i := 0; i < total; i++ {
		overrideIDs[i] = newUUIDName("customize-apple")
	}
	schema := newSearchOverrideSchema()
	expectedResult := map[string]*api.SearchOverride{}
	for i := 0; i < total; i++ {
		expectedResult[overrideIDs[i]] = newSearchOverride(overrideIDs[i])
	}

	for i := 0; i < total; i++ {
		_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(overrideIDs[i], schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Collection(collectionName).Overrides().Retrieve()

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of overrides is invalid")

	resultMap := map[string]*api.SearchOverride{}
	for _, override := range result {
		resultMap[override.Id] = override
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
