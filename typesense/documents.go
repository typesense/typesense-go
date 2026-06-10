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
	// Update documents with conditional query.
	//
	// The filter_by query parameter is used to filter to specify a condition against which the documents are matched. The request body contains the fields that should be updated for any documents that match the filter condition. This endpoint is only available if the Typesense server is version `0.25.0.rc12` or later.
	//
	// HTTP: PATCH /collections/{collectionName}/documents
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Update(ctx context.Context, updateFields interface{}, params *api.UpdateDocumentsParams) (int, error)
	// Upsert returns indexed/updated document
	Upsert(ctx context.Context, document interface{}, params *api.DocumentIndexParameters) (map[string]interface{}, error)
	// Delete a bunch of documents.
	//
	// Delete a bunch of documents that match a specific filter condition. Use the `batch_size` parameter to control the number of documents that should deleted at a time. A larger value will speed up deletions, but will impact performance of other operations running on the server.
	//
	// HTTP: DELETE /collections/{collectionName}/documents
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Delete(ctx context.Context, filter *api.DeleteDocumentsParams) (int, error)
	// Search for documents in a collection.
	//
	// Search for documents in a collection that match the search criteria.
	//
	// HTTP: GET /collections/{collectionName}/documents/search
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Search(ctx context.Context, params *api.SearchCollectionParams) (*api.SearchResult, error)
	// Export all documents in a collection.
	//
	// Export all documents in a collection in JSON lines format.
	//
	// HTTP: GET /collections/{collectionName}/documents/export
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Export(ctx context.Context, params *api.ExportDocumentsParams) (io.ReadCloser, error)
	// Import returns json array. Each item of the response indicates
	// the result of each document present in the request body (in the same order).
	Import(ctx context.Context, documents []interface{}, params *api.ImportDocumentsParams) ([]*api.ImportDocumentResponse, error)
	// Import documents into a collection.
	//
	// The documents to be imported must be formatted in a newline delimited JSON structure. You can feed the output file from a Typesense export operation directly as import.
	//
	// HTTP: POST /collections/{collectionName}/documents/import
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	ImportJsonl(ctx context.Context, body io.Reader, params *api.ImportDocumentsParams) (io.ReadCloser, error)
}

// documents is internal implementation of DocumentsInterface
type documents struct {
	apiClient      APIClientInterface
	collectionName string
}

// Index a document.
//
// A document to be indexed in a given collection must conform to the schema of the collection.
//
// HTTP: POST /collections/{collectionName}/documents
//
// See: https://typesense.org/docs/latest/api/documents.html
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

// Update documents with conditional query.
//
// The filter_by query parameter is used to filter to specify a condition against which the documents are matched. The request body contains the fields that should be updated for any documents that match the filter condition. This endpoint is only available if the Typesense server is version `0.25.0.rc12` or later.
//
// HTTP: PATCH /collections/{collectionName}/documents
//
// See: https://typesense.org/docs/latest/api/documents.html
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

// Delete a bunch of documents.
//
// Delete a bunch of documents that match a specific filter condition. Use the `batch_size` parameter to control the number of documents that should deleted at a time. A larger value will speed up deletions, but will impact performance of other operations running on the server.
//
// HTTP: DELETE /collections/{collectionName}/documents
//
// See: https://typesense.org/docs/latest/api/documents.html
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

// Search for documents in a collection.
//
// Search for documents in a collection that match the search criteria.
//
// HTTP: GET /collections/{collectionName}/documents/search
//
// See: https://typesense.org/docs/latest/api/documents.html
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

// Export all documents in a collection.
//
// Export all documents in a collection in JSON lines format.
//
// HTTP: GET /collections/{collectionName}/documents/export
//
// See: https://typesense.org/docs/latest/api/documents.html
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

// Import documents into a collection.
//
// The documents to be imported must be formatted in a newline delimited JSON structure. You can feed the output file from a Typesense export operation directly as import.
//
// HTTP: POST /collections/{collectionName}/documents/import
//
// See: https://typesense.org/docs/latest/api/documents.html
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
