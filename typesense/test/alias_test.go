// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
)

func TestCollectionAliasRetrieve(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	aliasName := newUUIDName("companies-alias")
	expectedResult := newCollectionAlias(collectionName, aliasName)

	body := &api.CollectionAliasSchema{CollectionName: collectionName}
	_, err := typesenseClient.Aliases().Upsert(aliasName, body)
	require.NoError(t, err)

	result, err := typesenseClient.Alias(aliasName).Retrieve()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestCollectionAliasDelete(t *testing.T) {
	collectionName := createNewCollection(t, "companies")
	aliasName := newUUIDName("companies-alias")
	expectedResult := newCollectionAlias(collectionName, aliasName)

	body := &api.CollectionAliasSchema{CollectionName: collectionName}
	_, err := typesenseClient.Aliases().Upsert(aliasName, body)
	require.NoError(t, err)

	result, err := typesenseClient.Alias(aliasName).Delete()

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	_, err = typesenseClient.Alias(aliasName).Retrieve()
	require.Error(t, err)
}
