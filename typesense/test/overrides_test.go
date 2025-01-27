//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v3/typesense/api"
	"github.com/typesense/typesense-go/v3/typesense/api/pointer"
)

func TestSearchOverrideUpsertNewOverride(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)

	body := newSearchOverrideSchema()
	result, err := typesenseClient.Collection(collectionName).Overrides().Upsert(context.Background(), overrideID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Override(overrideID).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSearchOverrideUpsertExistingOverride(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)
	expectedResult.Rule.Match = pointer.Any(api.Contains)

	body := newSearchOverrideSchema()
	body.Rule.Match = pointer.Any(api.Exact)
	_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(context.Background(), overrideID, body)
	require.NoError(t, err)

	body.Rule.Match = pointer.Any(api.Contains)

	result, err := typesenseClient.Collection(collectionName).Overrides().Upsert(context.Background(), overrideID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Collection(collectionName).Override(overrideID).Retrieve(context.Background())

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
		_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(context.Background(), overrideIDs[i], schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Collection(collectionName).Overrides().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of overrides is invalid")

	resultMap := map[string]*api.SearchOverride{}
	for _, override := range result {
		resultMap[*override.Id] = override
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
