package typesense

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func createNewDocument(docIDs ...string) interface{} {
	docID := "123"
	if len(docIDs) > 0 {
		docID = docIDs[0]
	}
	document := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"companyName"`
		NumEmployees int    `json:"numEmployees"`
		Country      string `json:"country"`
	}{
		ID:           docID,
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
	document["numEmployees"] = float64(5215)
	document["country"] = "USA"
	return document
}

func TestDocumentCreate(t *testing.T) {
	expectedDocument := createNewDocument()
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{}
	mockAPIClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, expectedDocument).
		Return(&api.IndexDocumentResponse{
			JSON201: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	result, err := client.Collection("companies").Documents().Create(context.Background(), document)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentCreateOnApiClientErrorReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{}
	mockAPIClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Documents().Create(context.Background(), newDocument)
	assert.NotNil(t, err)
}

func TestDocumentCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{}
	mockAPIClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Documents().Create(context.Background(), newDocument)
	assert.NotNil(t, err)
}

func TestDocumentUpsert(t *testing.T) {
	newDocument := createNewDocument()
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{Action: &upsertAction}
	mockAPIClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			JSON201: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Documents().Upsert(context.Background(), newDocument)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentUpsertOnApiClientErrorReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{Action: &upsertAction}
	mockAPIClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Documents().Upsert(context.Background(), newDocument)
	assert.NotNil(t, err)
}

func TestDocumentUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	indexParams := &api.IndexDocumentParams{Action: &upsertAction}
	mockAPIClient.EXPECT().
		IndexDocumentWithResponse(notNill, "companies", indexParams, newDocument).
		Return(&api.IndexDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Documents().Upsert(context.Background(), newDocument)
	assert.NotNil(t, err)
}

func TestDocumentsUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	expectedParams := &api.UpdateDocumentsParams{FilterBy: pointer.String("num_employees:>100")}

	mockedResult := struct {
		NumUpdated int `json:"num_updated"`
	}{27}

	expectedBody := strings.NewReader(`{"country":"USA"}`)

	mockAPIClient.EXPECT().
		UpdateDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedParams, eqReader(expectedBody)).
		Return(&api.UpdateDocumentsResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	updateFields := strings.NewReader(`{"country":"USA"}`)

	client := NewClient(WithAPIClient(mockAPIClient))

	params := &api.UpdateDocumentsParams{FilterBy: pointer.String("num_employees:>100")}
	result, err := client.Collection("companies").Documents().Update(context.Background(), updateFields, params)

	assert.Nil(t, err)
	assert.Equal(t, 27, result)
}

func TestDocumentsDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	expectedFilter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>100"), BatchSize: pointer.Int(100)}

	mockedResult := struct {
		NumDeleted int `json:"num_deleted"`
	}{27}

	mockAPIClient.EXPECT().
		DeleteDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedFilter).
		Return(&api.DeleteDocumentsResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	filter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>100"), BatchSize: pointer.Int(100)}
	result, err := client.Collection("companies").Documents().Delete(context.Background(), filter)

	assert.Nil(t, err)
	assert.Equal(t, 27, result)
}

func TestDocumentsDeleteOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	expectedFilter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>100"), BatchSize: pointer.Int(100)}

	mockAPIClient.EXPECT().
		DeleteDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedFilter).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	filter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>100"), BatchSize: pointer.Int(100)}
	_, err := client.Collection("companies").Documents().Delete(context.Background(), filter)
	assert.NotNil(t, err)
}

func TestDocumentsDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	expectedFilter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>100"), BatchSize: pointer.Int(100)}

	mockAPIClient.EXPECT().
		DeleteDocumentsWithResponse(gomock.Not(gomock.Nil()), "companies", expectedFilter).
		Return(&api.DeleteDocumentsResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	filter := &api.DeleteDocumentsParams{FilterBy: pointer.String("num_employees:>100"), BatchSize: pointer.Int(100)}
	_, err := client.Collection("companies").Documents().Delete(context.Background(), filter)
	assert.NotNil(t, err)
}

func createDocumentStream() io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(`{"id": "125","company_name":"Future Technology","num_employees":1232,"country":"UK"}`))
}

func TestDocumentsExport(t *testing.T) {
	expectedBytes, err := ioutil.ReadAll(createDocumentStream())
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createDocumentStream()

	mockAPIClient.EXPECT().
		ExportDocuments(gomock.Not(gomock.Nil()), "companies", &api.ExportDocumentsParams{}).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Documents().Export(context.Background())
	assert.Nil(t, err)

	resultBytes, err := ioutil.ReadAll(result)
	assert.Nil(t, err)
	assert.Equal(t, string(expectedBytes), string(resultBytes))
}

func TestDocumentsExportOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		ExportDocuments(gomock.Not(gomock.Nil()), "companies", &api.ExportDocumentsParams{}).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Documents().Export(context.Background())
	assert.NotNil(t, err)
}

func TestDocumentsExportOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		ExportDocuments(gomock.Not(gomock.Nil()), "companies", &api.ExportDocumentsParams{}).
		Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader("Internal server error")),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Documents().Export(context.Background())
	assert.NotNil(t, err)
}
