package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/typesense/typesense-go/typesense/api"
)

var upsertAction api.IndexDocumentParamsAction = "upsert"

const (
	defaultImportBatchSize = 40
	defaultImportAction    = "create"
)

// DocumentsInterface is a type for Documents API operations
type DocumentsInterface interface {
	// Create returns indexed document
	Create(ctx context.Context, document interface{}) (map[string]interface{}, error)
	// Update updates documents matching the filter_by condition
	Update(ctx context.Context, updateFields interface{}, params *api.UpdateDocumentsParams) (int, error)
	// Upsert returns indexed/updated document
	Upsert(context.Context, interface{}) (map[string]interface{}, error)
	// Delete returns number of deleted documents
	Delete(ctx context.Context, filter *api.DeleteDocumentsParams) (int, error)
	// Search performs document search in collection
	Search(ctx context.Context, params *api.SearchCollectionParams) (*api.SearchResult, error)
	// Export returns all documents from index in jsonl format
	Export(ctx context.Context) (io.ReadCloser, error)
	// Import returns json array. Each item of the response indicates
	// the result of each document present in the request body (in the same order).
	Import(ctx context.Context, documents []interface{}, params *api.ImportDocumentsParams) ([]*api.ImportDocumentResponse, error)
	// ImportJsonl accepts documents and returns result in jsonl format. Each line of the
	// response indicates the result of each document present in the
	// request body (in the same order).
	ImportJsonl(ctx context.Context, body io.Reader, params *api.ImportDocumentsParams) (io.ReadCloser, error)
}

// documents is internal implementation of DocumentsInterface
type documents struct {
	apiClient      APIClientInterface
	collectionName string
}

func (d *documents) indexDocument(ctx context.Context, document interface{}, params *api.IndexDocumentParams) (map[string]interface{}, error) {
	response, err := d.apiClient.IndexDocumentWithResponse(ctx,
		d.collectionName, params, document)
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON201, nil
}

func (d *documents) Create(ctx context.Context, document interface{}) (map[string]interface{}, error) {
	return d.indexDocument(ctx, document, &api.IndexDocumentParams{})
}

func (d *documents) Update(ctx context.Context, updateFields interface{}, params *api.UpdateDocumentsParams) (int, error) {
	response, err := d.apiClient.UpdateDocumentsWithResponse(ctx,
		d.collectionName, params, updateFields)
	if err != nil {
		return 0, err
	}
	if response.JSON200 == nil {
		return 0, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.NumUpdated, nil
}

func (d *documents) Upsert(ctx context.Context, document interface{}) (map[string]interface{}, error) {
	return d.indexDocument(ctx, document, &api.IndexDocumentParams{Action: &upsertAction})
}

func (d *documents) Delete(ctx context.Context, filter *api.DeleteDocumentsParams) (int, error) {
	response, err := d.apiClient.DeleteDocumentsWithResponse(ctx,
		d.collectionName, filter)
	if err != nil {
		return 0, err
	}
	if response.JSON200 == nil {
		return 0, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.NumDeleted, nil
}

func (d *documents) Search(ctx context.Context, params *api.SearchCollectionParams) (*api.SearchResult, error) {
	response, err := d.apiClient.SearchCollectionWithResponse(ctx,
		d.collectionName, params)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (d *documents) Export(ctx context.Context) (io.ReadCloser, error) {
	response, err := d.apiClient.ExportDocuments(ctx, d.collectionName, &api.ExportDocumentsParams{})
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &HTTPError{Status: response.StatusCode, Body: body}
	}
	return response.Body, nil
}

func initImportParams(params *api.ImportDocumentsParams) {
	if params.BatchSize == nil {
		params.BatchSize = new(int)
		*params.BatchSize = defaultImportBatchSize
	}
	if params.Action == nil {
		params.Action = new(string)
		*params.Action = defaultImportAction
	}
}

func (d *documents) Import(ctx context.Context, documents []interface{}, params *api.ImportDocumentsParams) ([]*api.ImportDocumentResponse, error) {
	if len(documents) == 0 {
		return nil, errors.New("documents list is empty")
	}

	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)
	for _, doc := range documents {
		if err := jsonEncoder.Encode(doc); err != nil {
			return nil, err
		}
	}

	response, err := d.ImportJsonl(ctx, &buf, params)
	if err != nil {
		return nil, err
	}

	var result []*api.ImportDocumentResponse
	jsonDecoder := json.NewDecoder(response)
	for jsonDecoder.More() {
		var docResult *api.ImportDocumentResponse
		if err := jsonDecoder.Decode(&docResult); err != nil {
			return result, errors.New("failed to decode result")
		}
		result = append(result, docResult)
	}

	return result, nil
}

func (d *documents) ImportJsonl(ctx context.Context, body io.Reader, params *api.ImportDocumentsParams) (io.ReadCloser, error) {
	initImportParams(params)
	response, err := d.apiClient.ImportDocumentsWithBody(ctx,
		d.collectionName, params, "application/octet-stream", body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &HTTPError{Status: response.StatusCode, Body: body}
	}
	return response.Body, nil
}
