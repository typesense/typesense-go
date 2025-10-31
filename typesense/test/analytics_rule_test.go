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

func analyticsRuleCleanUp() {
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

func TestAnalyticsRule(t *testing.T) {
	shouldSkipAnalyticsTests(t)
	t.Cleanup(analyticsRuleCleanUp)

	t.Run("Create", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		
		// Create the rule directly using the Create API
		ruleName := newUUIDName("test-rule")
		ruleCreate := &api.AnalyticsRuleCreate{
			Name:       ruleName,
			Type:       api.AnalyticsRuleCreateTypeCounter,
			Collection: collectionName,
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("num_employees"),
				Weight:       pointer.Int(1),
			},
		}

		result, err := typesenseClient.Analytics().Rules().Create(context.Background(), []*api.AnalyticsRuleCreate{ruleCreate})
		require.NoError(t, err)
		require.Len(t, result, 1)
		
		createdRule := result[0]
		require.Equal(t, ruleName, createdRule.Name)
		require.Equal(t, api.AnalyticsRuleTypeCounter, createdRule.Type)
		require.Equal(t, collectionName, createdRule.Collection)
		require.Equal(t, "click", createdRule.EventType)
	})

	t.Run("Retrieve", func(t *testing.T) {
		eventName := newUUIDName("event")
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		expectedRule := createNewAnalyticsRule(t, collectionName, sourceCollectionName, eventName)

		result, err := typesenseClient.Analytics().Rule(expectedRule.Name).Retrieve(context.Background())
		require.NoError(t, err)
		require.Equal(t, expectedRule.Name, result.Name)
		require.Equal(t, expectedRule.Type, result.Type)
		require.Equal(t, expectedRule.Collection, result.Collection)
		require.Equal(t, expectedRule.EventType, result.EventType)
	})

	t.Run("Delete", func(t *testing.T) {
		eventName := newUUIDName("event")
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		expectedRule := createNewAnalyticsRule(t, collectionName, sourceCollectionName, eventName)

		result, err := typesenseClient.Analytics().Rule(expectedRule.Name).Delete(context.Background())
		require.NoError(t, err)
		require.Equal(t, expectedRule.Name, result.Name)
	})
}
