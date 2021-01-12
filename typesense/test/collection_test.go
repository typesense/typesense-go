// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectionRetrieve(t *testing.T) {
	collectionName := getNewCollectionName("companies")
	newSchema := createNewSchema(collectionName)
	expectedResult := expectedNewCollection(collectionName)

	_, err := typesenseClient.Collections().Create(newSchema)
	require.NoError(t, err)

	result, err := typesenseClient.Collection(collectionName).Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionDelete(t *testing.T) {
	collectionName := getNewCollectionName("companies")
	newSchema := createNewSchema(collectionName)
	expectedResult := expectedNewCollection(collectionName)

	_, err := typesenseClient.Collections().Create(newSchema)
	require.NoError(t, err)

	result, err := typesenseClient.Collection(collectionName).Delete()
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	_, err = typesenseClient.Collection(collectionName).Retrieve()
	require.Error(t, err)
}
