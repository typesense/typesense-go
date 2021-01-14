// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchOverrideRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)

	body := newSearchOverrideSchema()
	_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(overrideID, body)
	require.NoError(t, err)

	result, err := typesenseClient.Collection(collectionName).Override(overrideID).Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSearchOverrideDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)

	body := newSearchOverrideSchema()
	_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(overrideID, body)
	require.NoError(t, err)

	result, err := typesenseClient.Collection(collectionName).Override(overrideID).Delete()

	require.NoError(t, err)
	require.Equal(t, expectedResult.Id, result.Id)

	_, err = typesenseClient.Collection(collectionName).Override(overrideID).Retrieve()
	require.Error(t, err)
}
