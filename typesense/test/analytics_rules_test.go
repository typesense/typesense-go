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

	t.Run("RetrieveWithParams", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		createNewAnalyticsRule(t, collectionName, sourceCollectionName, newUUIDName("event"))

		ruleTag := newUUIDName("homepage")
		taggedRule := &api.AnalyticsRuleCreate{
			Name:       newUUIDName("test-rule"),
			Type:       api.Counter,
			Collection: collectionName,
			EventType:  "click",
			RuleTag:    pointer.String(ruleTag),
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("num_employees"),
				Weight:       pointer.Int(1),
			},
		}
		created, err := typesenseClient.Analytics().Rules().Create(context.Background(), []*api.AnalyticsRuleCreate{taggedRule})
		require.NoError(t, err)
		require.Equal(t, taggedRule.Name, created[0].Name)

		results, err := typesenseClient.Analytics().Rules().RetrieveWithParams(context.Background(), &api.RetrieveAnalyticsRulesParams{RuleTag: pointer.String(ruleTag)})
		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Equal(t, taggedRule.Name, results[0].Name)

		all, err := typesenseClient.Analytics().Rules().RetrieveWithParams(context.Background(), nil)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(all), 2)
	})

	t.Run("CreateMany", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		existingRule := createNewAnalyticsRule(t, collectionName, sourceCollectionName, newUUIDName("event"))

		newRule := &api.AnalyticsRuleCreate{
			Name:       newUUIDName("test-rule"),
			Type:       api.Counter,
			Collection: collectionName,
			EventType:  "click",
			Params:     existingRule.Params,
		}
		duplicateRule := &api.AnalyticsRuleCreate{
			Name:       existingRule.Name,
			Type:       existingRule.Type,
			Collection: existingRule.Collection,
			EventType:  existingRule.EventType,
			Params:     existingRule.Params,
		}

		responses, err := typesenseClient.Analytics().Rules().CreateMany(context.Background(), []*api.AnalyticsRuleCreate{newRule, duplicateRule})
		require.NoError(t, err)
		require.Len(t, responses, 2)
		require.Empty(t, responses[0].Error)
		require.Equal(t, newRule.Name, responses[0].Name)
		require.NotEmpty(t, responses[1].Error)
		require.Empty(t, responses[1].Name)
	})

	t.Run("Upsert", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		ruleName := newUUIDName("test-rule")
		rule := &api.AnalyticsRuleCreate{
			Name:       ruleName,
			Type:       api.Counter,
			Collection: collectionName,
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("num_employees"),
				Weight:       pointer.Int(1),
			},
		}

		created, err := typesenseClient.Analytics().Rules().Upsert(context.Background(), ruleName, rule)
		require.NoError(t, err)
		require.Equal(t, ruleName, created.Name)
		require.Equal(t, 1, *created.Params.Weight)

		rule.Params.Weight = pointer.Int(5)
		_, err = typesenseClient.Analytics().Rules().Upsert(context.Background(), ruleName, rule)
		require.NoError(t, err)

		updated, err := typesenseClient.Analytics().Rule(ruleName).Retrieve(context.Background())
		require.NoError(t, err)
		require.Equal(t, 5, *updated.Params.Weight)
	})
}
