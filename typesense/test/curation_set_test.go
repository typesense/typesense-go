//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func curationSetCleanUp() {
	curationSets, _ := typesenseClient.CurationSets().Retrieve(context.Background())
	for _, cs := range curationSets {
		typesenseClient.CurationSet(cs.Name).Delete(context.Background())
	}
}

func TestCurationSet(t *testing.T) {
	shouldSkipAnalyticsTests(t)
	t.Cleanup(curationSetCleanUp)

	t.Run("Retrieve", func(t *testing.T) {
		curationSetName := newUUIDName("test-curation-set")
		curationSetData := newCurationSetCreateSchema()

		_, err := typesenseClient.CurationSets().Upsert(context.Background(), curationSetName, curationSetData)
		require.NoError(t, err)

		result, err := typesenseClient.CurationSet(curationSetName).Retrieve(context.Background())

		require.NoError(t, err)
		require.NotNil(t, result)
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

	t.Run("Update", func(t *testing.T) {
		curationSetName := newUUIDName("test-curation-set")
		originalData := newCurationSetCreateSchema()

		_, err := typesenseClient.CurationSets().Upsert(context.Background(), curationSetName, originalData)
		require.NoError(t, err)

		updatedData := &api.CurationSetCreateSchema{
			Items: []api.CurationItemCreateSchema{
				{
					Rule: api.CurationRule{
						Query: pointer.String("updated query"),
						Match: pointer.Any(api.Exact),
					},
					Id: pointer.String("dummy"),
					Includes: &[]api.CurationInclude{
						{
							Id: "422",
						},
						{
							Id: "54",
						},
						{
							Id: "999",
						},
					},
					Excludes: &[]api.CurationExclude{
						{
							Id: "287",
						},
					},
					RemoveMatchedTokens: pointer.True(),
					FilterBy:            pointer.String("category:=Electronics"),
					StopProcessing:      pointer.True(),
				},
			},
			Description: pointer.String("Updated test curation set"),
		}

		result, err := typesenseClient.CurationSet(curationSetName).Upsert(context.Background(), updatedData)

		require.NoError(t, err)
		require.Equal(t, curationSetName, result.Name)
		require.Equal(t, 1, len(result.Items))
		require.Equal(t, "dummy", *result.Items[0].Id)
		require.Equal(t, "updated query", *result.Items[0].Rule.Query)
		require.NotNil(t, result.Items[0].Includes)
		require.Equal(t, 3, len(*result.Items[0].Includes))
		require.Equal(t, "422", (*result.Items[0].Includes)[0].Id)
		require.Equal(t, "54", (*result.Items[0].Includes)[1].Id)
		require.Equal(t, "999", (*result.Items[0].Includes)[2].Id)
		require.NotNil(t, result.Items[0].Excludes)
		require.Equal(t, 1, len(*result.Items[0].Excludes))
		require.Equal(t, "287", (*result.Items[0].Excludes)[0].Id)
	})

	t.Run("Delete", func(t *testing.T) {
		curationSetName := newUUIDName("test-curation-set")
		curationSetData := newCurationSetCreateSchema()

		_, err := typesenseClient.CurationSets().Upsert(context.Background(), curationSetName, curationSetData)
		require.NoError(t, err)

		result, err := typesenseClient.CurationSet(curationSetName).Delete(context.Background())

		require.NoError(t, err)
		require.Equal(t, curationSetName, result.Name)

		_, err = typesenseClient.CurationSet(curationSetName).Retrieve(context.Background())
		require.Error(t, err)
	})
}
