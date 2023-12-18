package typesense

import (
	"context"
)

type DocumentInterface interface {
	Retrieve(ctx context.Context) (map[string]interface{}, error)
	Update(ctx context.Context, document interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context) (map[string]interface{}, error)
}

type document struct {
	apiClient      APIClientInterface
	collectionName string
	documentID     string
}

func (d *document) Retrieve(ctx context.Context) (map[string]interface{}, error) {
	response, err := d.apiClient.GetDocumentWithResponse(ctx,
		d.collectionName, d.documentID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}

func (d *document) Update(ctx context.Context, document interface{}) (map[string]interface{}, error) {
	response, err := d.apiClient.UpdateDocumentWithResponse(ctx,
		d.collectionName, d.documentID, document)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}

func (d *document) Delete(ctx context.Context) (map[string]interface{}, error) {
	response, err := d.apiClient.DeleteDocumentWithResponse(ctx,
		d.collectionName, d.documentID)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200, nil
}
