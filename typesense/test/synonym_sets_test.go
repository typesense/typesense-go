//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

func TestSynonymSetsUpsert(t *testing.T) {
	synonymSetName := newUUIDName("test-synonym-set")
	synonymSetData := &api.SynonymSetCreateSchema{
		Items: []api.SynonymItemSchema{
			{
				Id:       "dummy",
				Synonyms: []string{"foo", "bar", "baz"},
			},
		},
	}

	result, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, synonymSetData)

	require.NoError(t, err)
	require.Equal(t, synonymSetName, result.Name)
	require.Equal(t, 1, len(result.Items))
	require.Equal(t, "dummy", result.Items[0].Id)
	require.Equal(t, []string{"foo", "bar", "baz"}, result.Items[0].Synonyms)

	// Cleanup
	_, err = typesenseClient.SynonymSet(synonymSetName).Delete(context.Background())
	require.NoError(t, err)
}

func TestSynonymSetsRetrieve(t *testing.T) {
	synonymSetName := newUUIDName("test-synonym-set")
	synonymSetData := &api.SynonymSetCreateSchema{
		Items: []api.SynonymItemSchema{
			{
				Id:       "dummy",
				Synonyms: []string{"foo", "bar", "baz"},
			},
		},
	}

	// Create a synonym set first
	_, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, synonymSetData)
	require.NoError(t, err)

	// Retrieve all synonym sets
	result, err := typesenseClient.SynonymSets().Retrieve(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.GreaterOrEqual(t, len(result), 1)

	// Find our created synonym set
	var foundSynonymSet *api.SynonymSetSchema
	for _, ss := range result {
		if ss.Name == synonymSetName {
			foundSynonymSet = &ss
			break
		}
	}

	require.NotNil(t, foundSynonymSet)
	require.Equal(t, synonymSetName, foundSynonymSet.Name)
	require.Equal(t, 1, len(foundSynonymSet.Items))
	require.Equal(t, "dummy", foundSynonymSet.Items[0].Id)
	require.Equal(t, []string{"foo", "bar", "baz"}, foundSynonymSet.Items[0].Synonyms)

	// Cleanup
	_, err = typesenseClient.SynonymSet(synonymSetName).Delete(context.Background())
	require.NoError(t, err)
} 