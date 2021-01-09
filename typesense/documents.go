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

var upsertAction = "upsert"

const (
	defaultImportBatchSize = 40
	defaultImportAction    = "create"
)

// DocumentsInterface is a type for Documents API operations
type DocumentsInterface interface {
	// Create returns indexed document
	Create(document interface{}) (map[string]interface{}, error)
	// Upsert returns indexed/updated document
	Upsert(document interface{}) (map[string]interface{}, error)
	// Delete returns number of deleted documents
	Delete(filter *api.DeleteDocumentsParams) (int, error)
	// Search performs document search in collection
	Search(params *api.SearchCollectionParams) (*api.SearchResult, error)
	// Export returns all documents from index in jsonl format
	Export() (io.ReadCloser, error)
	// Import returns json array. Each item of the response indicates
	// the result of each document present in the request body (in the same order).
	Import(documents []interface{}, params *api.ImportDocumentsParams) ([]*api.ImportDocumentResponse, error)
	// ImportJsonl accepts documents and returns result in jsonl format. Each line of the
	// response indicates the result of each document present in the
	// request body (in the same order).
	ImportJsonl(body io.Reader, params *api.ImportDocumentsParams) (io.ReadCloser, error)
}

// documents is internal implementation of DocumentsInterface
type documents struct {
	apiClient      APIClientInterface
	collectionName string
}

func (d *documents) indexDocument(document interface{}, params *api.IndexDocumentParams) (map[string]interface{}, error) {
	response, err := d.apiClient.IndexDocumentWithResponse(context.Background(),
		d.collectionName, params, document)
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return *response.JSON201, nil
}

func (d *documents) Create(document interface{}) (map[string]interface{}, error) {
	return d.indexDocument(document, &api.IndexDocumentParams{})
}

func (d *documents) Upsert(document interface{}) (map[string]interface{}, error) {
	return d.indexDocument(document, &api.IndexDocumentParams{Action: &upsertAction})
}

func (d *documents) Delete(filter *api.DeleteDocumentsParams) (int, error) {
	response, err := d.apiClient.DeleteDocumentsWithResponse(context.Background(),
		d.collectionName, filter)
	if err != nil {
		return 0, err
	}
	if response.JSON200 == nil {
		return 0, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200.NumDeleted, nil
}

func (d *documents) Search(params *api.SearchCollectionParams) (*api.SearchResult, error) {
	response, err := d.apiClient.SearchCollectionWithResponse(context.Background(),
		d.collectionName, params)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200, nil
}

func (d *documents) Export() (io.ReadCloser, error) {
	response, err := d.apiClient.ExportDocuments(context.Background(), d.collectionName)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &httpError{status: response.StatusCode, body: body}
	}
	return response.Body, nil
}

func initImportParams(params *api.ImportDocumentsParams) {
	if params.BatchSize == 0 {
		params.BatchSize = defaultImportBatchSize
	}
	if len(params.Action) == 0 {
		params.Action = defaultImportAction
	}
}

func (d *documents) Import(documents []interface{}, params *api.ImportDocumentsParams) ([]*api.ImportDocumentResponse, error) {
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

	response, err := d.ImportJsonl(&buf, params)
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

func (d *documents) ImportJsonl(body io.Reader, params *api.ImportDocumentsParams) (io.ReadCloser, error) {
	initImportParams(params)
	response, err := d.apiClient.ImportDocumentsWithBody(context.Background(),
		d.collectionName, params, "application/octet-stream", body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &httpError{status: response.StatusCode, body: body}
	}
	return response.Body, nil
}
