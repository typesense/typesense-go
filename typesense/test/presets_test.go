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

func presetsCleanUp() {
	result, _ := typesenseClient.Presets().Retrieve(context.Background())
	for _, preset := range result {
		typesenseClient.Preset(preset.Name).Delete(context.Background())
	}
}

func TestPresetsUpsertValueFromSearchParameters(t *testing.T) {
	t.Cleanup(presetsCleanUp)

	presetID := newUUIDName("preset-test")
	expectedResult := newPresetFromSearchParameters(presetID)

	body := newPresetFromSearchParametersUpsertSchema()
	result, err := typesenseClient.Presets().Upsert(context.Background(), presetID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	parsedExpected, err := expectedResult.Value.AsSearchParameters()
	require.NoError(t, err)

	parsedResult, err := expectedResult.Value.AsSearchParameters()
	require.NoError(t, err)

	require.Equal(t, parsedExpected, parsedResult)
}

func TestPresetsUpsertValueFromMultiSearchSearchesParameter(t *testing.T) {
	t.Cleanup(presetsCleanUp)

	presetID := newUUIDName("preset-test")
	expectedResult := newPresetFromMultiSearchSearchesParameter(presetID)

	body := newPresetFromMultiSearchSearchesParameterUpsertSchema()
	result, err := typesenseClient.Presets().Upsert(context.Background(), presetID, body)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

	parsedExpected, err := expectedResult.Value.AsMultiSearchSearchesParameter()
	require.NoError(t, err)

	parsedResult, err := expectedResult.Value.AsMultiSearchSearchesParameter()
	require.NoError(t, err)

	require.Equal(t, parsedExpected, parsedResult)
}

func TestPresetsRetrievePresetFromSearchParameters(t *testing.T) {
	t.Cleanup(presetsCleanUp)

	total := 3
	presetIDs := make([]string, total)
	for i := 0; i < total; i++ {
		presetIDs[i] = newUUIDName("preset-test")
	}
	schema := newPresetFromSearchParametersUpsertSchema()

	expectedResult := map[string]*api.PresetSchema{}
	for i := 0; i < total; i++ {
		expectedResult[presetIDs[i]] = newPresetFromSearchParameters(presetIDs[i])
	}

	for i := 0; i < total; i++ {
		_, err := typesenseClient.Presets().Upsert(context.Background(), presetIDs[i], schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Presets().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of presets is invalid")

	resultMap := map[string]*api.PresetSchema{}
	for _, preset := range result {
		resultMap[preset.Name] = preset
	}

	for k, v := range expectedResult {
		assert.EqualValues(t, v, resultMap[k])

		parsedExpected, err := v.Value.AsSearchParameters()
		require.NoError(t, err)

		parsedResult, err := resultMap[k].Value.AsSearchParameters()
		require.NoError(t, err)

		require.Equal(t, parsedExpected, parsedResult)
	}
}
func TestPresetsRetrieveValueFromMultiSearchSearchesParameter(t *testing.T) {
	t.Cleanup(presetsCleanUp)

	total := 3
	presetIDs := make([]string, total)
	for i := 0; i < total; i++ {
		presetIDs[i] = newUUIDName("preset-test")
	}
	schema := newPresetFromMultiSearchSearchesParameterUpsertSchema()

	expectedResult := map[string]*api.PresetSchema{}
	for i := 0; i < total; i++ {
		expectedResult[presetIDs[i]] = newPresetFromMultiSearchSearchesParameter(presetIDs[i])
	}

	for i := 0; i < total; i++ {
		_, err := typesenseClient.Presets().Upsert(context.Background(), presetIDs[i], schema)
		require.NoError(t, err)
	}

	result, err := typesenseClient.Presets().Retrieve(context.Background())

	require.NoError(t, err)
	require.True(t, len(result) >= total, "number of presets is invalid")

	resultMap := map[string]*api.PresetSchema{}
	for _, preset := range result {
		resultMap[preset.Name] = preset
	}

	for k, v := range expectedResult {
		assert.EqualValues(t, v, resultMap[k])

		parsedExpected, err := v.Value.AsMultiSearchSearchesParameter()
		require.NoError(t, err)

		parsedResult, err := resultMap[k].Value.AsMultiSearchSearchesParameter()
		require.NoError(t, err)

		require.Equal(t, parsedExpected, parsedResult)
	}
}
