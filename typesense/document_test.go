package typesense

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func TestDocumentRetrieve(t *testing.T) {
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	mockAPIClient.EXPECT().
		GetDocument(gomock.Not(gomock.Nil()), "companies", "123").
		Return(createResponse(200, "", mockedResult), nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Document("123").Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetDocument(gomock.Not(gomock.Nil()), "companies", "123").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestDocumentRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetDocument(gomock.Not(gomock.Nil()), "companies", "123").
		Return(createResponse(500, "Internal server error", nil), nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Retrieve(context.Background())
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
		UpdateDocument(notNill, "companies", "123", &api.UpdateDocumentParams{DirtyValues: pointer.Any(api.CoerceOrDrop)}, expectedDocument).
		Return(createResponse(200, "", mockedResult), nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	result, err := client.Collection("companies").Document("123").Update(context.Background(), document, &api.DocumentIndexParameters{DirtyValues: pointer.Any(api.CoerceOrDrop)})

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
		UpdateDocument(notNill, "companies", "123", &api.UpdateDocumentParams{}, expectedDocument).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	_, err := client.Collection("companies").Document("123").Update(context.Background(), document, &api.DocumentIndexParameters{})
	assert.NotNil(t, err)
}

func TestDocumentUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	expectedDocument := createNewDocument()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	notNill := gomock.Not(gomock.Nil())
	mockAPIClient.EXPECT().
		UpdateDocument(notNill, "companies", "123", &api.UpdateDocumentParams{}, expectedDocument).
		Return(createResponse(500, "Internal server error", nil), nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	document := createNewDocument()
	_, err := client.Collection("companies").Document("123").Update(context.Background(), document, &api.DocumentIndexParameters{})
	assert.NotNil(t, err)
}

func TestDocumentDelete(t *testing.T) {
	expectedResult := createNewDocumentResponse()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewDocumentResponse()

	mockAPIClient.EXPECT().
		DeleteDocument(gomock.Not(gomock.Nil()), "companies", "123").
		Return(createResponse(200, "", mockedResult), nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Document("123").Delete(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDocumentDeleteOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteDocument(gomock.Not(gomock.Nil()), "companies", "123").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Delete(context.Background())
	assert.NotNil(t, err)
}

func TestDocumentDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteDocument(gomock.Not(gomock.Nil()), "companies", "123").
		Return(createResponse(500, "Internal server error", nil), nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Document("123").Delete(context.Background())
	assert.NotNil(t, err)
}
