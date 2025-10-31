//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

func nlSearchModelsCleanUp() {
	result, _ := typesenseClient.NLSearchModels().Retrieve(context.Background())
	for _, model := range result {
		typesenseClient.NLSearchModel(model.Id).Delete(context.Background())
	}
}

func TestNLSearchModel(t *testing.T) {
	shouldSkipNLSearchModelTests(t)
	t.Cleanup(nlSearchModelsCleanUp)

	t.Run("Retrieve", func(t *testing.T) {
		modelID, expectedResult := createNewNLSearchModel(t)

		result, err := typesenseClient.NLSearchModel(modelID).Retrieve(context.Background())

		require.NoError(t, err)
		require.Equal(t, expectedResult, result)
	})

	t.Run("Update", func(t *testing.T) {
		modelID, originalModel := createNewNLSearchModel(t)

		updateSchema := newNLSearchModelUpdateSchema()
		updateSchema.Temperature = pointer.Float32(0.8)

		result, err := typesenseClient.NLSearchModel(modelID).Update(context.Background(), updateSchema)

		require.NoError(t, err)
		require.Equal(t, "openai/gpt-4", *result.ModelName)
		require.Equal(t, float32(0.8), *result.Temperature)
		require.Equal(t, originalModel.Id, result.Id)
	})

	t.Run("Delete", func(t *testing.T) {
		modelID, expectedResult := createNewNLSearchModel(t)

		result, err := typesenseClient.NLSearchModel(modelID).Delete(context.Background())

		require.NoError(t, err)
		require.Equal(t, expectedResult.Id, result.Id)

		_, err = typesenseClient.NLSearchModel(modelID).Retrieve(context.Background())
		require.Error(t, err)
	})
} 