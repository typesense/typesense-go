//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func TestCollectionAliasUpsertNewAlias(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	aliasName := newUUIDName("companies-alias")
	expectedResult := newCollectionAlias(collectionName, aliasName)

	body := &api.CollectionAliasSchema{CollectionName: collectionName}
	result, err := typesenseClient.Aliases().Upsert(context.Background(), aliasName, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Alias(aliasName).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionAliasUpsertExistingAlias(t *testing.T) {
	collectionName1 := createNewCollection(t, "companies")
	collectionName2 := createNewCollection(t, "companies")
	aliasName := newUUIDName("companies-alias")
	expectedResult := newCollectionAlias(collectionName2, aliasName)

	body := &api.CollectionAliasSchema{CollectionName: collectionName1}
	_, err := typesenseClient.Aliases().Upsert(context.Background(), aliasName, body)
	require.NoError(t, err)

	body = &api.CollectionAliasSchema{CollectionName: collectionName2}
	result, err := typesenseClient.Aliases().Upsert(context.Background(), aliasName, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	result, err = typesenseClient.Alias(aliasName).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionAliasesRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	total := 3
	aliasNames := make([]string, total)
	for i := 0; i < total; i++ {
		aliasNames[i] = newUUIDName("companies-alias")
	}
	body := &api.CollectionAliasSchema{CollectionName: collectionName}
	expectedResult := map[string]*api.CollectionAlias{}
	for i := 0; i < total; i++ {
		expectedResult[aliasNames[i]] = newCollectionAlias(collectionName, aliasNames[i])
	}

	for i := 0; i < total; i++ {
		_, err := typesenseClient.Aliases().Upsert(context.Background(), aliasNames[i], body)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Aliases().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of aliases is invalid")

	resultMap := map[string]*api.CollectionAlias{}
	for _, alias := range result {
		resultMap[*alias.Name] = alias
	}

	for k, v := range expectedResult {
		assert.Equal(t, v, resultMap[k])
	}
}
