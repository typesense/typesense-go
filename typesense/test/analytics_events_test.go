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

func analyticsEventsCleanUp() {
	// Clean up analytics rules
	result, _ := typesenseClient.Analytics().Rules().Retrieve(context.Background())
	for _, rule := range result {
		typesenseClient.Analytics().Rule(rule.Name).Delete(context.Background())
	}
	
	// Clean up collections
	collections, _ := typesenseClient.Collections().Retrieve(context.Background(), nil)
	for _, collection := range collections {
		typesenseClient.Collection(collection.Name).Delete(context.Background())
	}
}

func TestAnalyticsEvents(t *testing.T) {
	shouldSkipAnalyticsTests(t)
	t.Cleanup(analyticsEventsCleanUp)

	t.Run("Create", func(t *testing.T) {
		eventName := newUUIDName("event")
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		
		expectedRule := createNewAnalyticsRule(t, collectionName, sourceCollectionName, eventName)

		result, err := typesenseClient.Analytics().Events().Create(context.Background(), &api.AnalyticsEvent{
			Name:      expectedRule.Name,
			EventType: "click",
			Data: api.AnalyticsEventData{
				Q:       pointer.String("nike shoes"),
				DocId:   pointer.String("1024"),
				UserId:  pointer.String("111112"),
			},
		})

		require.NoError(t, err)
		require.True(t, result.Ok)
	})

	t.Run("Retrieve", func(t *testing.T) {
		eventName := newUUIDName("event")
		collectionName := createNewCollection(t, "analytics-events-collection")
		
		expectedRule := createNewAnalyticsRule(t, collectionName, collectionName, eventName)

		_, err := typesenseClient.Analytics().Events().Create(context.Background(), &api.AnalyticsEvent{
			Name:      expectedRule.Name,
			EventType: "click",
			Data: api.AnalyticsEventData{
				Q:       pointer.String("nike shoes"),
				DocId:   pointer.String("1024"),
				UserId:  pointer.String("111112"),
			},
		})

		result, err := typesenseClient.Analytics().Events().Retrieve(context.Background(), &api.GetAnalyticsEventsParams{
			UserId: "111112",
			Name:   expectedRule.Name,
			N:      1000,
		})
		
		require.NoError(t, err)
		require.NotNil(t, result)
		require.IsType(t, &api.AnalyticsEventsResponse{}, result)
	})
}