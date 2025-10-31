//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func analyticsRulesCleanUp() {
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

func TestAnalyticsRules(t *testing.T) {
	shouldSkipAnalyticsTests(t)
	t.Cleanup(analyticsRulesCleanUp)

	t.Run("Retrieve", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		eventName := newUUIDName("event")
		expectedRule := createNewAnalyticsRule(t, collectionName, sourceCollectionName, eventName)

		results, err := typesenseClient.Analytics().Rules().Retrieve(context.Background())
		require.NoError(t, err)
		require.True(t, len(results) >= 1, "number of rules is invalid")

		var result *api.AnalyticsRule
		for _, rule := range results {
			if rule.Name == expectedRule.Name {
				result = rule
				break
			}
		}

		require.NotNil(t, result, "rule not found")
		require.Equal(t, expectedRule.Name, result.Name)
		require.Equal(t, expectedRule.Type, result.Type)
		require.Equal(t, expectedRule.Collection, result.Collection)
		require.Equal(t, expectedRule.EventType, result.EventType)
	})
}
