//go:build integration
// +build integration

package test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyRetrieve(t *testing.T) {
	expectedKey := createNewKey(t)

	result, err := typesenseClient.Key(*expectedKey.Id).Retrieve(context.Background())

	require.NoError(t, err)

	require.Equal(t, expectedKey.Description, result.Description)
	require.Equal(t, expectedKey.Actions, result.Actions)
	require.Equal(t, expectedKey.Collections, result.Collections)
	require.True(t, strings.HasPrefix(*expectedKey.Value, *result.ValuePrefix),
		"value_prefix is invalid")
}

func TestKeyDelete(t *testing.T) {
	expectedKey := createNewKey(t)

	result, err := typesenseClient.Key(*expectedKey.Id).Delete(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedKey.Id, result.Id)

	_, err = typesenseClient.Key(*expectedKey.Id).Retrieve(context.Background())
	require.Error(t, err)
}
