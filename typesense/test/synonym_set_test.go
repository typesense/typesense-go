//go:build integration
// +build integration

package test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func synonymSetCleanUp() {
	synonymSets, _ := typesenseClient.SynonymSets().Retrieve(context.Background())
	for _, ss := range synonymSets {
		typesenseClient.SynonymSet(ss.Name).Delete(context.Background())
	}
}

func TestSynonymSet(t *testing.T) {
	shouldSkipAnalyticsTests(t)
	t.Cleanup(synonymSetCleanUp)

	t.Run("Retrieve", func(t *testing.T) {
		synonymSetName := newUUIDName("test-synonym-set")
		synonymSetData := &api.SynonymSetCreateSchema{
			Items: []api.SynonymItemSchema{
				{
					Id:       "dummy",
					Synonyms: []string{"foo", "bar", "baz"},
				},
			},
		}

		_, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, synonymSetData)
		require.NoError(t, err)

		result, err := typesenseClient.SynonymSet(synonymSetName).Retrieve(context.Background())

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 1, len(result.Items))
		require.Equal(t, "dummy", result.Items[0].Id)
		require.Equal(t, []string{"foo", "bar", "baz"}, result.Items[0].Synonyms)
	})

	t.Run("Update", func(t *testing.T) {
		synonymSetName := newUUIDName("test-synonym-set")
		originalData := &api.SynonymSetCreateSchema{
			Items: []api.SynonymItemSchema{
				{
					Id:       "dummy",
					Synonyms: []string{"foo", "bar"},
				},
			},
		}

		_, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, originalData)
		require.NoError(t, err)

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
	})

	t.Run("Delete", func(t *testing.T) {
		synonymSetName := newUUIDName("test-synonym-set")
		synonymSetData := &api.SynonymSetCreateSchema{
			Items: []api.SynonymItemSchema{
				{
					Id:       "dummy",
					Synonyms: []string{"foo", "bar", "baz"},
				},
			},
		}

		_, err := typesenseClient.SynonymSets().Upsert(context.Background(), synonymSetName, synonymSetData)
		require.NoError(t, err)

		result, err := typesenseClient.SynonymSet(synonymSetName).Delete(context.Background())

		require.NoError(t, err)
		require.Equal(t, synonymSetName, result.Name)

		_, err = typesenseClient.SynonymSet(synonymSetName).Retrieve(context.Background())
		require.Error(t, err)
	})

	t.Run("Items", func(t *testing.T) {
		synonymSetName, _ := createNewSynonymSet(t)
		item := &api.SynonymItemUpsertSchema{
			Root:     pointer.String("smart phone"),
			Synonyms: []string{"iphone", "android"},
		}

		upserted, err := typesenseClient.SynonymSet(synonymSetName).Items().Upsert(context.Background(), "smart-phone", item)
		require.NoError(t, err)
		require.Equal(t, "smart-phone", upserted.Id)
		require.Equal(t, item.Synonyms, upserted.Synonyms)

		items, err := typesenseClient.SynonymSet(synonymSetName).Items().Retrieve(context.Background())
		require.NoError(t, err)
		require.Len(t, items, 2)
		require.ElementsMatch(t, []string{"dummy", "smart-phone"}, []string{items[0].Id, items[1].Id})
	})

	t.Run("Item", func(t *testing.T) {
		synonymSetName, _ := createNewSynonymSet(t)
		item := typesenseClient.SynonymSet(synonymSetName).Item("smart-phone")

		_, err := item.Upsert(context.Background(), &api.SynonymItemUpsertSchema{Synonyms: []string{"iphone", "android"}})
		require.NoError(t, err)

		retrieved, err := item.Retrieve(context.Background())
		require.NoError(t, err)
		require.Equal(t, "smart-phone", retrieved.Id)
		require.Equal(t, []string{"iphone", "android"}, retrieved.Synonyms)

		deleted, err := item.Delete(context.Background())
		require.NoError(t, err)
		require.Equal(t, "smart-phone", deleted.Id)

		_, err = item.Retrieve(context.Background())
		var httpErr *typesense.HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusNotFound, httpErr.Status)
	})
}
