//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchOverrideRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)

	body := newSearchOverrideSchema()
	_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(context.Background(), overrideID, body)
	require.NoError(t, err)

	result, err := typesenseClient.Collection(collectionName).Override(overrideID).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSearchOverrideDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	overrideID := newUUIDName("customize-apple")
	expectedResult := newSearchOverride(overrideID)

	body := newSearchOverrideSchema()
	_, err := typesenseClient.Collection(collectionName).Overrides().Upsert(context.Background(), overrideID, body)
	require.NoError(t, err)

	result, err := typesenseClient.Collection(collectionName).Override(overrideID).Delete(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult.Id, result.Id)

	_, err = typesenseClient.Collection(collectionName).Override(overrideID).Retrieve(context.Background())
	require.Error(t, err)
}
