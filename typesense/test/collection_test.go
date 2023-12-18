//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func TestCollectionRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := expectedNewCollection(collectionName)

	result, err := typesenseClient.Collection(collectionName).Retrieve(context.Background())
	result.CreatedAt = pointer.Int64(0)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	expectedResult := expectedNewCollection(collectionName)

	result, err := typesenseClient.Collection(collectionName).Delete(context.Background())
	result.CreatedAt = pointer.Int64(0)
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	_, err = typesenseClient.Collection(collectionName).Retrieve(context.Background())
	require.Error(t, err)
}

func TestCollectionUpdate(t *testing.T) {
	collectionName := createNewCollection(t, "companies")

	updateSchema := &api.CollectionUpdateSchema{
		Fields: []api.Field{
			{
				Name: "country",
				Drop: pointer.True(),
			},
		},
	}

	result, err := typesenseClient.Collection(collectionName).Update(context.Background(), updateSchema)
	require.NoError(t, err)
	require.Equal(t, 1, len(result.Fields))
	require.Equal(t, "country", result.Fields[0].Name)
	require.Equal(t, pointer.True(), result.Fields[0].Drop)
}
