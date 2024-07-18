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

func TestCollectionCreate(t *testing.T) {
	collectionName := newUUIDName("companies")
	schema := newSchema(collectionName)
	expectedResult := expectedNewCollection(collectionName)

	result, err := typesenseClient.Collections().Create(context.Background(), schema)
	result.CreatedAt = pointer.Int64(0)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionsRetrieve(t *testing.T) {
	total := 3
	collectionNames := make([]string, total)
	for i := 0; i < total; i++ {
		collectionNames[i] = newUUIDName("companies")
	}
	schemas := make([]*api.CollectionSchema, total)
	for i := 0; i < total; i++ {
		schemas[i] = newSchema(collectionNames[i])
	}
	expectedResult := map[string]*api.CollectionResponse{}
	for i := 0; i < total; i++ {
		expectedResult[collectionNames[i]] = expectedNewCollection(collectionNames[i])
	}

	for _, schema := range schemas {
		_, err := typesenseClient.Collections().Create(context.Background(), schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Collections().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of collections is invalid")

	resultMap := map[string]*api.CollectionResponse{}
	for _, collection := range result {
		resultMap[collection.Name] = collection
		resultMap[collection.Name].CreatedAt = pointer.Int64(0)
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
