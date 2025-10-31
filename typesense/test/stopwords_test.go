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

func TestStopwordsUpsert(t *testing.T) {
	stopwordsSetID := newUUIDName("stopwordsSet-test")
	upsertData := &api.StopwordsSetUpsertSchema{
		Locale:    pointer.String("en"),
		Stopwords: []string{"Germany", "France", "Italy", "United States"},
	}

	expectedData := &api.StopwordsSetSchema{
		Id:        stopwordsSetID,
		Locale:    upsertData.Locale,
		Stopwords: upsertData.Stopwords,
	}

	result, err := typesenseClient.Stopwords().Upsert(context.Background(), stopwordsSetID, upsertData)

	require.NoError(t, err)
	require.Equal(t, expectedData, result)
}

func TestStopwordsRetrieve(t *testing.T) {
	total := 3
	stopwordsSetIDs := make([]string, total)
	for i := 0; i < total; i++ {
		stopwordsSetIDs[i] = newUUIDName("stopwordsSet-test")
	}

	upsertData := &api.StopwordsSetUpsertSchema{
		Locale:    pointer.String("en"),
		Stopwords: []string{"Germany", "France", "Italy", "United States"},
	}

	expectedResult := map[string]api.StopwordsSetSchema{}
	for i := 0; i < total; i++ {
		expectedResult[stopwordsSetIDs[i]] = api.StopwordsSetSchema{
			Id:        stopwordsSetIDs[i],
			Locale:    upsertData.Locale,
			Stopwords: []string{"states", "united", "france", "germany", "italy"},
		}
	}

	for i := 0; i < total; i++ {
		_, err := typesenseClient.Stopwords().Upsert(context.Background(), stopwordsSetIDs[i], upsertData)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Stopwords().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of stopwordsSets is invalid")

	resultMap := map[string]api.StopwordsSetSchema{}
	for _, stopwordsSet := range result {
		resultMap[stopwordsSet.Id] = stopwordsSet
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
