//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

func TestSynonymSetRetrieve(t *testing.T) {
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

	// Retrieve the specific synonym set
	result, err := typesenseClient.SynonymSet(synonymSetName).Retrieve(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result.Items))
	require.Equal(t, "dummy", result.Items[0].Id)
	require.Equal(t, []string{"foo", "bar", "baz"}, result.Items[0].Synonyms)

	// Cleanup
	_, err = typesenseClient.SynonymSet(synonymSetName).Delete(context.Background())
	require.NoError(t, err)
}

func TestSynonymSetUpdate(t *testing.T) {
	synonymSetName := newUUIDName("test-synonym-set")
	originalData := &api.SynonymSetCreateSchema{
		Items: []api.SynonymItemSchema{
			{
				Id:       "dummy",
				Synonyms: []string{"foo", "bar"},
			},
		},
	}

	// Create a synonym set first
	_, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, originalData)
	require.NoError(t, err)

	// Update the synonym set
	updatedData := &api.SynonymSetCreateSchema{
		Items: []api.SynonymItemSchema{
			{
				Id:       "dummy",
				Synonyms: []string{"foo", "bar", "baz", "qux"},
			},
		},
	}

	result, err := typesenseClient.SynonymSet(synonymSetName).Upsert(context.Background(), updatedData)

	require.NoError(t, err)
	require.Equal(t, synonymSetName, result.Name)
	require.Equal(t, 1, len(result.Items))
	require.Equal(t, "dummy", result.Items[0].Id)
	require.Equal(t, []string{"foo", "bar", "baz", "qux"}, result.Items[0].Synonyms)

	// Cleanup
	_, err = typesenseClient.SynonymSet(synonymSetName).Delete(context.Background())
	require.NoError(t, err)
}

func TestSynonymSetDelete(t *testing.T) {
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

	// Delete the synonym set
	result, err := typesenseClient.SynonymSet(synonymSetName).Delete(context.Background())

	require.NoError(t, err)
	require.Equal(t, synonymSetName, result.Name)

	// Verify it's deleted by trying to retrieve it
	_, err = typesenseClient.SynonymSet(synonymSetName).Retrieve(context.Background())
	require.Error(t, err)
} 