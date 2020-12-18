package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/v-byte-cpu/typesense-go/typesense/api"
	"github.com/v-byte-cpu/typesense-go/typesense/api/mocks"
)

func createNewDocument() interface{} {
	document := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"companyName"`
		NumEmployees int    `json:"numEmployees"`
		Country      string `json:"country"`
	}{
		ID:           "123",
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}
	return &document
}

func createNewDocumentResponse() map[string]interface{} {
	document := map[string]interface{}{}
	document["id"] = "123"
	document["companyName"] = "Stark Industries"
	document["numEmployees"] = 5215
	document["country"] = "USA"
	return document
}

func TestDocumentCreate(t *testing.T) {
	newDocument := createNewDocument()
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{}
	mockApiClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			JSON201: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	result, err := client.Collection("companies").Documents().Create(newDocument)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentCreateOnApiClientErrorReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{}
	mockApiClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Documents().Create(newDocument)
	assert.NotNil(t, err)
}

func TestDocumentCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{}
	mockApiClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Documents().Create(newDocument)
	assert.NotNil(t, err)
}

func TestDocumentUpsert(t *testing.T) {
	newDocument := createNewDocument()
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{Action: &upsertAction}
	mockApiClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			JSON201: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	result, err := client.Collection("companies").Documents().Upsert(newDocument)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentUpsertOnApiClientErrorReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{Action: &upsertAction}
	mockApiClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Documents().Upsert(newDocument)
	assert.NotNil(t, err)
}

func TestDocumentUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{Action: &upsertAction}
	mockApiClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Documents().Upsert(newDocument)
	assert.NotNil(t, err)
}

func TestDocumentsDelete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	expectedFilter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}

	mockedResult := struct {
		NumDeleted int `json:"num_deleted"`
	}{27}

	mockApiClient.EXPECT().
		DeleteDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedFilter).
		Return(&api.DeleteDocumentsResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	filter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}
	result, err := client.Collection("companies").Documents().Delete(filter)

	assert.Nil(t, err)
	assert.Equal(t, 27, result)
}

func TestDocumentsDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	expectedFilter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}

	mockApiClient.EXPECT().
		DeleteDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedFilter).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	filter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}
	_, err := client.Collection("companies").Documents().Delete(filter)
	assert.NotNil(t, err)
}

func TestDocumentsDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	expectedFilter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}

	mockApiClient.EXPECT().
		DeleteDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedFilter).
		Return(&api.DeleteDocumentsResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	filter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}
	_, err := client.Collection("companies").Documents().Delete(filter)
	assert.NotNil(t, err)
}
