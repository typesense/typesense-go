//go:build integration
// +build integration

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func TestStemmingDictionary(t *testing.T) {
    dictionaryId := fmt.Sprintf("dictionary_%d", time.Now().UnixNano())
    words := []api.StemmingDictionaryWord{
        {
            Root: "exampleRoot1",
            Word: "exampleWord1",
        },
        {
            Root: "exampleRoot2",
            Word: "exampleWord2",
        },
    }

    t.Run("Upsert", func(t *testing.T) {
        result, err := typesenseClient.Stemming().Dictionaries().Upsert(
            context.Background(),
            dictionaryId,
            words,
        )
        require.NoError(t, err)
        require.Len(t, result, len(words))
    })

    t.Run("Retrieve Single", func(t *testing.T) {
        result, err := typesenseClient.Stemming().Dictionary(dictionaryId).Retrieve(context.Background())
        require.NoError(t, err)
        require.Equal(t, dictionaryId, result.Id)
        // Convert result.Words to the same type for comparison
        retrievedWords := make([]api.StemmingDictionaryWord, len(result.Words))
        for i, w := range result.Words {
            retrievedWords[i] = api.StemmingDictionaryWord{
                Root: w.Root,
                Word: w.Word,
            }
        }
        require.Equal(t, words, retrievedWords)
    })

    t.Run("List All", func(t *testing.T) {
        result, err := typesenseClient.Stemming().Dictionaries().Retrieve(context.Background())
        require.NoError(t, err)
        require.Contains(t, *result.JSON200.Dictionaries, dictionaryId)
    })
}