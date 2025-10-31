//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func curationSetsCleanUp() {
	curationSets, _ := typesenseClient.CurationSets().Retrieve(context.Background())
	for _, cs := range curationSets {
		typesenseClient.CurationSet(cs.Name).Delete(context.Background())
	}
}

func TestCurationSets(t *testing.T) {
	shouldSkipAnalyticsTests(t)
	t.Cleanup(curationSetsCleanUp)

	t.Run("Upsert", func(t *testing.T) {
		curationSetName := newUUIDName("test-curation-set")
		curationSetData := newCurationSetCreateSchema()

		result, err := typesenseClient.CurationSets().Upsert(context.Background(), curationSetName, curationSetData)

		require.NoError(t, err)
		require.Equal(t, curationSetName, result.Name)
		require.Equal(t, 1, len(result.Items))
		require.Equal(t, "dummy", *result.Items[0].Id)
		require.Equal(t, "apple", *result.Items[0].Rule.Query)
		require.NotNil(t, result.Items[0].Includes)
		require.Equal(t, 2, len(*result.Items[0].Includes))
		require.Equal(t, "422", (*result.Items[0].Includes)[0].Id)
		require.Equal(t, "54", (*result.Items[0].Includes)[1].Id)
		require.NotNil(t, result.Items[0].Excludes)
		require.Equal(t, 1, len(*result.Items[0].Excludes))
		require.Equal(t, "287", (*result.Items[0].Excludes)[0].Id)
	})

	t.Run("Retrieve", func(t *testing.T) {
		curationSetName := newUUIDName("test-curation-set")
		curationSetData := newCurationSetCreateSchema()

		_, err := typesenseClient.CurationSets().Upsert(context.Background(), curationSetName, curationSetData)
		require.NoError(t, err)

		result, err := typesenseClient.CurationSets().Retrieve(context.Background())

		require.NoError(t, err)
		require.NotNil(t, result)
		require.GreaterOrEqual(t, len(result), 1)

		var foundCurationSet *api.CurationSetSchema
		for _, cs := range result {
			if cs.Name == curationSetName {
				foundCurationSet = &cs
				break
			}
		}

		require.NotNil(t, foundCurationSet)
		require.Equal(t, curationSetName, foundCurationSet.Name)
		require.Equal(t, 1, len(foundCurationSet.Items))
		require.Equal(t, "dummy", *foundCurationSet.Items[0].Id)
		require.Equal(t, "apple", *foundCurationSet.Items[0].Rule.Query)
		require.NotNil(t, foundCurationSet.Items[0].Includes)
		require.Equal(t, 2, len(*foundCurationSet.Items[0].Includes))
		require.Equal(t, "422", (*foundCurationSet.Items[0].Includes)[0].Id)
		require.Equal(t, "54", (*foundCurationSet.Items[0].Includes)[1].Id)
		require.NotNil(t, foundCurationSet.Items[0].Excludes)
		require.Equal(t, 1, len(*foundCurationSet.Items[0].Excludes))
		require.Equal(t, "287", (*foundCurationSet.Items[0].Excludes)[0].Id)
	})
}
