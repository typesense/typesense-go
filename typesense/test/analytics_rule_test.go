//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyticsRuleRetrieve(t *testing.T) {
	eventName := newUUIDName("event")
	collectionName := createNewCollection(t, "analytics-rules-collection")
	expectedRule := createNewAnalyticsRule(t, collectionName, eventName)

	result, err := typesenseClient.Analytics().Rule(expectedRule.Name).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedRule, result)
}

func TestAnalyticsRuleDelete(t *testing.T) {
	eventName := newUUIDName("event")
	collectionName := createNewCollection(t, "analytics-rules-collection")
	expectedRule := createNewAnalyticsRule(t, collectionName, eventName)

	result, err := typesenseClient.Analytics().Rule(expectedRule.Name).Delete(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedRule.Name, result.Name)
}
