// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestCollectionCreate(t *testing.T) {
	collectionName := getNewCollectionName("companies")
	newSchema := createNewSchema(collectionName)
	expectedResult := expectedNewCollection(collectionName)

	result, err := typesenseClient.Collections().Create(newSchema)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionsRetrieve(t *testing.T) {
	total := 3
	collectionNames := make([]string, total)
	for i := 0; i < total; i++ {
		collectionNames[i] = getNewCollectionName("companies")
	}
	newSchemas := make([]*api.CollectionSchema, total)
	for i := 0; i < total; i++ {
		newSchemas[i] = createNewSchema(collectionNames[i])
	}
	expectedResult := map[string]*api.Collection{}
	for i := 0; i < total; i++ {
		expectedResult[collectionNames[i]] = expectedNewCollection(collectionNames[i])
	}

	for _, newSchema := range newSchemas {
		_, err := typesenseClient.Collections().Create(newSchema)
		assert.NoError(t, err)
	}

	result, err := typesenseClient.Collections().Retrieve()

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of collections is invalid")

	resultMap := map[string]*api.Collection{}
	for _, collection := range result {
		resultMap[collection.Name] = collection
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
