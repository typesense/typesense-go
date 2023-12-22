//go:build integration
// +build integration

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func TestDocumentsImport(t *testing.T) {
	collectionName := createNewCollection(t, "companies")

	expectedResults := []map[string]interface{}{
		newDocumentResponse("123"),
		newDocumentResponse("125", withResponseCompanyName("Company2")),
		newDocumentResponse("127", withResponseCompanyName("Company3")),
	}

	documents := []interface{}{
		newDocument("123"),
		newDocument("125", withCompanyName("Company2")),
		newDocument("127", withCompanyName("Company3")),
	}

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	responses, err := typesenseClient.Collection(collectionName).Documents().Import(context.Background(), documents, params)

	require.NoError(t, err)
	for _, response := range responses {
		require.True(t, response.Success, "document import failed")
	}

	results := retrieveDocuments(t, collectionName, "123", "125", "127")
	require.Equal(t, expectedResults, results)
}

func TestDocumentsImportJsonl(t *testing.T) {
	collectionName := createNewCollection(t, "companies")

	expectedResults := []map[string]interface{}{
		newDocumentResponse("123"),
		newDocumentResponse("125", withResponseCompanyName("Company2")),
		newDocumentResponse("127", withResponseCompanyName("Company3")),
	}

	var buffer bytes.Buffer
	je := json.NewEncoder(&buffer)
	require.NoError(t, je.Encode(newDocument("123")))
	require.NoError(t, buffer.WriteByte('\n'))
	require.NoError(t, je.Encode(newDocument("125", withCompanyName("Company2"))))
	require.NoError(t, buffer.WriteByte('\n'))
	require.NoError(t, je.Encode(newDocument("127", withCompanyName("Company3"))))

	params := &api.ImportDocumentsParams{Action: pointer.String("create")}
	responses, err := typesenseClient.Collection(collectionName).Documents().ImportJsonl(context.Background(), &buffer, params)

	require.NoError(t, err)
	defer responses.Close()

	jd := json.NewDecoder(responses)
	for i := 0; i < 3; i++ {
		require.True(t, jd.More(), "no json element")
		response := &api.ImportDocumentResponse{}
		require.NoError(t, jd.Decode(&response))
		require.True(t, response.Success, "document import failed")
	}

	results := retrieveDocuments(t, collectionName, "123", "125", "127")
	require.Equal(t, expectedResults, results)
}
