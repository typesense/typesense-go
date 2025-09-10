//go:build integration
// +build integration

package test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func TestKeyCreate(t *testing.T) {
	keySchema := newKeySchema()
	expectedResult := newKey()

	result, err := typesenseClient.Keys().Create(context.Background(), keySchema)

	require.NoError(t, err)
	require.Equal(t, expectedResult.Description, result.Description)
	require.Equal(t, expectedResult.Actions, result.Actions)
	require.Equal(t, expectedResult.Collections, result.Collections)
	require.NotEmpty(t, result.Value)
}

func TestKeysRetrieve(t *testing.T) {
	expectedKey := createNewKey(t)

	results, err := typesenseClient.Keys().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(results) >= 1, "number of keys is invalid")
	var result *api.ApiKey
	for _, key := range results {
		if *key.Id == *expectedKey.Id {
			result = key
			break
		}
	}

	require.NotNil(t, result, "key not found")
	require.Equal(t, expectedKey.Description, result.Description)
	require.Equal(t, expectedKey.Actions, result.Actions)
	require.Equal(t, expectedKey.Collections, result.Collections)
	require.True(t, strings.HasPrefix(*expectedKey.Value, *result.ValuePrefix),
		"value_prefix is invalid")
}
