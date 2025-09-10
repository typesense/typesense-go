package typesense

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
)

const (
	defaultImportBatchSize = 40
)

// DocumentsInterface is a type for Documents API operations
type DocumentsInterface interface {
	// Create returns indexed document
	Create(ctx context.Context, document interface{}, params *api.DocumentIndexParameters) (map[string]interface{}, error)
	// Update updates documents matching the filter_by condition
	Update(ctx context.Context, updateFields interface{}, params *api.UpdateDocumentsParams) (int, error)
	// Upsert returns indexed/updated document
	Upsert(ctx context.Context, document interface{}, params *api.DocumentIndexParameters) (map[string]interface{}, error)
	// Delete returns number of deleted documents
	Delete(ctx context.Context, filter *api.DeleteDocumentsParams) (int, error)
	// Search performs document search in collection
	Search(ctx context.Context, params *api.SearchCollectionParams) (*api.SearchResult, error)
	// Export returns all documents from index in jsonl format
	Export(ctx context.Context, params *api.ExportDocumentsParams) (io.ReadCloser, error)
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

func (d *documents) Create(ctx context.Context, document interface{}, params *api.DocumentIndexParameters) (map[string]interface{}, error) {
	return d.indexDocument(ctx, document, &api.IndexDocumentParams{DirtyValues: params.DirtyValues})
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

func (d *documents) Upsert(ctx context.Context, document interface{}, params *api.DocumentIndexParameters) (map[string]interface{}, error) {
	return d.indexDocument(ctx, document, &api.IndexDocumentParams{Action: pointer.Any(api.Upsert), DirtyValues: params.DirtyValues})
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

func (d *documents) Export(ctx context.Context, params *api.ExportDocumentsParams) (io.ReadCloser, error) {
	response, err := d.apiClient.ExportDocuments(ctx, d.collectionName, params)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
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
		params.Action = pointer.Any(api.Create)
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
	defer response.Close()

	var result []*api.ImportDocumentResponse
	scanner := bufio.NewScanner(response)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var docResult *api.ImportDocumentResponse
		if err := json.Unmarshal(scanner.Bytes(), &docResult); err != nil {
			return result, fmt.Errorf("failed to decode result: %w", err)
		}
		result = append(result, docResult)
	}

	return result, scanner.Err()
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
		body, _ := io.ReadAll(response.Body)
		return nil, &HTTPError{Status: response.StatusCode, Body: body}
	}
	return response.Body, nil
}
