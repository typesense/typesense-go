package typesense

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

var upsertAction = "upsert"

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
		body, _ := ioutil.ReadAll(response.Body)
		return nil, &httpError{status: response.StatusCode, body: body}
	}
	return response.Body, nil
}

// api.ActionMode
// api.CreateAction
// api.UpsertAction

// TODO client.Collection('name').Documents().Import(documents, WithAction(api.CreateAction))
// TODO client.Collection('name').Documents().ImportJsonlFile(body io.Reader, WithAction(api.UpsertAction))
