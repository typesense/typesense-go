package typesense

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type DocumentInterface[T any] interface {
	// Retrieve a document.
	//
	// Fetch an individual document from a collection by using its ID.
	//
	// HTTP: GET /collections/{collectionName}/documents/{documentId}
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Retrieve(ctx context.Context) (T, error)
	// Update a document.
	//
	// Update an individual document from a collection by using its ID. The update can be partial.
	//
	// HTTP: PATCH /collections/{collectionName}/documents/{documentId}
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Update(ctx context.Context, document any, params *api.DocumentIndexParameters) (T, error)
	// Delete a document.
	//
	// Delete an individual document from a collection by using its ID.
	//
	// HTTP: DELETE /collections/{collectionName}/documents/{documentId}
	//
	// See: https://typesense.org/docs/latest/api/documents.html
	Delete(ctx context.Context) (T, error)
}

var _ DocumentInterface[any] = (*document[any])(nil)

type document[T any] struct {
	apiClient      APIClientInterface
	collectionName string
	documentID     string
}

// Retrieve a document.
//
// Fetch an individual document from a collection by using its ID.
//
// HTTP: GET /collections/{collectionName}/documents/{documentId}
//
// See: https://typesense.org/docs/latest/api/documents.html
func (d *document[T]) Retrieve(ctx context.Context) (resp T, err error) {
	response, err := d.apiClient.GetDocument(ctx,
		d.collectionName, d.documentID)
	if err != nil {
		return resp, err
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "json") || response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		return resp, &HTTPError{Status: response.StatusCode, Body: body}
	}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Update a document.
//
// Update an individual document from a collection by using its ID. The update can be partial.
//
// HTTP: PATCH /collections/{collectionName}/documents/{documentId}
//
// See: https://typesense.org/docs/latest/api/documents.html
func (d *document[T]) Update(ctx context.Context, document any, params *api.DocumentIndexParameters) (resp T, err error) {
	response, err := d.apiClient.UpdateDocument(ctx,
		d.collectionName, d.documentID, &api.UpdateDocumentParams{DirtyValues: params.DirtyValues}, document)
	if err != nil {
		return resp, err
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "json") || response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		return resp, &HTTPError{Status: response.StatusCode, Body: body}
	}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Delete a document.
//
// Delete an individual document from a collection by using its ID.
//
// HTTP: DELETE /collections/{collectionName}/documents/{documentId}
//
// See: https://typesense.org/docs/latest/api/documents.html
func (d *document[T]) Delete(ctx context.Context) (resp T, err error) {
	response, err := d.apiClient.DeleteDocument(ctx,
		d.collectionName, d.documentID)
	if err != nil {
		return resp, err
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "json") || response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		return resp, &HTTPError{Status: response.StatusCode, Body: body}
	}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
