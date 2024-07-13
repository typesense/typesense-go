//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPresetValueFromSearchParametersRetrieve(t *testing.T) {
	t.Cleanup(presetsCleanUp)
	presetName, expectedResult := createNewPreset(t, true)

	result, err := typesenseClient.Preset(presetName).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	parsedExpected, err := expectedResult.Value.AsSearchParameters()
	require.NoError(t, err)

	parsedResult, err := expectedResult.Value.AsSearchParameters()
	require.NoError(t, err)

	require.Equal(t, parsedExpected, parsedResult)
}

func TestPresetsFromMultiSearchSearchesParameterRetrieve(t *testing.T) {
	t.Cleanup(presetsCleanUp)
	presetName, expectedResult := createNewPreset(t)

	result, err := typesenseClient.Preset(presetName).Retrieve(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	parsedExpected, err := expectedResult.Value.AsMultiSearchSearchesParameter()
	require.NoError(t, err)

	parsedResult, err := expectedResult.Value.AsMultiSearchSearchesParameter()
	require.NoError(t, err)

	require.Equal(t, parsedExpected, parsedResult)
}

func TestPresetDelete(t *testing.T) {
	t.Cleanup(presetsCleanUp)
	presetName, expectedResult := createNewPreset(t)

	result, err := typesenseClient.Preset(presetName).Delete(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResult.Name, result.Name)

	_, err = typesenseClient.Preset(presetName).Retrieve(context.Background())
	require.Error(t, err)
}
