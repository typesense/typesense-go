//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func TestCollectionCreate(t *testing.T) {
	collectionName := newUUIDName("companies")
	schema := newSchema(collectionName)
	expectedResult := expectedNewCollection(t, collectionName)

	result, err := typesenseClient.Collections().Create(context.Background(), schema)
	require.NoError(t, err)
	result.CreatedAt = pointer.Int64(0)
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
		expectedResult[collectionNames[i]] = expectedNewCollection(t, collectionNames[i])
	}

	for _, schema := range schemas {
		_, err := typesenseClient.Collections().Create(context.Background(), schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Collections().Retrieve(context.Background(), &api.GetCollectionsParams{})

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

	exclude := "fields"
	excluded, err := typesenseClient.Collections().Retrieve(context.Background(), &api.GetCollectionsParams{
		ExcludeFields: &exclude,
	})
	require.NoError(t, err)
	for _, collection := range excluded {
		assert.Empty(t, collection.Fields)
	}
}
