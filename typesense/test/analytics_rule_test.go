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

func analyticsRuleCleanUp() {
	// Clean up analytics rules
	result, _ := typesenseClient.Analytics().Rules().Retrieve(context.Background(), nil)
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
			Type:       api.AnalyticsRuleTypeCounter,
			Collection: collectionName,
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("num_employees"),
				Weight:       pointer.Int(1),
			},
		}

		createdRule, err := typesenseClient.Analytics().Rules().Create(context.Background(), ruleCreate)
		require.NoError(t, err)
		require.Equal(t, ruleName, createdRule.Name)
		require.Equal(t, api.AnalyticsRuleTypeCounter, createdRule.Type)
		require.Equal(t, collectionName, createdRule.Collection)
		require.Equal(t, "click", createdRule.EventType)
	})

	t.Run("CreateExistingReturnsHTTPError", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		ruleCreate := &api.AnalyticsRuleCreate{
			Name:       newUUIDName("test-rule"),
			Type:       api.AnalyticsRuleTypeCounter,
			Collection: collectionName,
			EventType:  "click",
			Params: &api.AnalyticsRuleCreateParams{
				CounterField: pointer.String("num_employees"),
				Weight:       pointer.Int(1),
			},
		}
		_, err := typesenseClient.Analytics().Rules().Create(context.Background(), ruleCreate)
		require.NoError(t, err)

		result, err := typesenseClient.Analytics().Rules().Create(context.Background(), ruleCreate)
		require.Nil(t, result)
		var httpErr *typesense.HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusBadRequest, httpErr.Status)
		require.Contains(t, string(httpErr.Body), "Rule already exists")
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

	t.Run("Update", func(t *testing.T) {
		collectionName := createNewCollection(t, "analytics-rules-collection")
		sourceCollectionName := createNewCollection(t, "analytics-rules-source-collection")
		expectedRule := createNewAnalyticsRule(t, collectionName, sourceCollectionName, newUUIDName("event"))

		result, err := typesenseClient.Analytics().Rule(expectedRule.Name).Update(context.Background(), &api.AnalyticsRuleUpdate{
			Params: &api.AnalyticsRuleUpdateParams{Weight: pointer.Int(5)},
		})
		require.NoError(t, err)
		require.Equal(t, expectedRule.Name, result.Name)
		require.Equal(t, expectedRule.Type, result.Type)
		require.Equal(t, expectedRule.Collection, result.Collection)
		require.Equal(t, expectedRule.EventType, result.EventType)
		require.Equal(t, *expectedRule.Params.CounterField, *result.Params.CounterField)
		require.Equal(t, 5, *result.Params.Weight)
	})

	t.Run("UpdateMissingReturnsHTTPError", func(t *testing.T) {
		result, err := typesenseClient.Analytics().Rule(newUUIDName("missing-rule")).Update(context.Background(), &api.AnalyticsRuleUpdate{
			Params: &api.AnalyticsRuleUpdateParams{Weight: pointer.Int(5)},
		})
		require.Nil(t, result)
		var httpErr *typesense.HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, http.StatusBadRequest, httpErr.Status)
	})
}
