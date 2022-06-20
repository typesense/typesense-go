package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/mocks"
)

func TestDocumentRetrieve(t *testing.T) {
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	mockAPIClient.EXPECT().
		GetDocumentWithResponse(gomock.Not(gomock.Nil()), "companies", "123").
		Return(&api.GetDocumentResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Document("123").Retrieve()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetDocumentWithResponse(gomock.Not(gomock.Nil()), "companies", "123").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Retrieve()
	assert.NotNil(t, err)
}

func TestDocumentRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetDocumentWithResponse(gomock.Not(gomock.Nil()), "companies", "123").
		Return(&api.GetDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Retrieve()
	assert.NotNil(t, err)
}

func TestDocumentUpdate(t *testing.T) {
	expectedDocument := createNewDocument()
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	notNill := gomock.Not(gomock.Nil())
	mockAPIClient.EXPECT().
		UpdateDocumentWithResponse(notNill, "companies", "123", expectedDocument).
		Return(&api.UpdateDocumentResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	result, err := client.Collection("companies").Document("123").Update(document)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentUpdateOnApiClientErrorReturnsError(t *testing.T) {
	expectedDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	mockAPIClient.EXPECT().
		UpdateDocumentWithResponse(notNill, "companies", "123", expectedDocument).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	_, err := client.Collection("companies").Document("123").Update(document)
	assert.NotNil(t, err)
}

func TestDocumentUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	expectedDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	mockAPIClient.EXPECT().
		UpdateDocumentWithResponse(notNill, "companies", "123", expectedDocument).
		Return(&api.UpdateDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	_, err := client.Collection("companies").Document("123").Update(document)
	assert.NotNil(t, err)
}

func TestDocumentDelete(t *testing.T) {
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	mockAPIClient.EXPECT().
		DeleteDocumentWithResponse(gomock.Not(gomock.Nil()), "companies", "123").
		Return(&api.DeleteDocumentResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Document("123").Delete()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentDeleteOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteDocumentWithResponse(gomock.Not(gomock.Nil()), "companies", "123").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Delete()
	assert.NotNil(t, err)
}

func TestDocumentDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteDocumentWithResponse(gomock.Not(gomock.Nil()), "companies", "123").
		Return(&api.DeleteDocumentResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Delete()
	assert.NotNil(t, err)
}
