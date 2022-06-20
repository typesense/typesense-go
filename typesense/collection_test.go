package typesense

import (
	"errors"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/mocks"
)

func updateExistingSchema() *api.CollectionUpdateSchema {
	return &api.CollectionUpdateSchema{
		Fields: []api.Field{
			{
				Name: "url",
				Drop: pointer.True(),
			},
			{
				Name:  "url",
				Type:  "string",
				Index: pointer.False(),
			},
		},
	}
}

func updateExistingCollection() *api.CollectionUpdateSchema {
	return updateExistingSchema()
}

func TestCollectionRetrieve(t *testing.T) {
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollection("companies")

	mockAPIClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Retrieve()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionDelete(t *testing.T) {
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollection("companies")

	mockAPIClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.DeleteCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Delete()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Delete()
	assert.NotNil(t, err)
}

func TestCollectionDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.DeleteCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Delete()
	assert.NotNil(t, err)
}

func TestCollectionUpdate(t *testing.T) {
	updateSchema := updateExistingSchema()
	expectedResult := updateExistingCollection()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := updateExistingCollection()

	mockAPIClient.EXPECT().
		UpdateCollectionWithResponse(gomock.Not(gomock.Nil()), "companies",
			api.UpdateCollectionJSONRequestBody(*updateSchema)).
		Return(&api.UpdateCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Update(updateSchema)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionUpdateOnApiClientErrorReturnsError(t *testing.T) {
	updateSchema := updateExistingSchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpdateCollectionWithResponse(gomock.Not(gomock.Nil()), "companies",
			api.UpdateCollectionJSONRequestBody(*updateSchema)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Update(updateSchema)
	assert.Error(t, err)
}

func TestCollectionUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	updateSchema := updateExistingSchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpdateCollectionWithResponse(gomock.Not(gomock.Nil()), "non_existent",
			api.UpdateCollectionJSONRequestBody(*updateSchema)).
		Return(&api.UpdateCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 404,
			},
			Body: []byte("Collection not found"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("non_existent").Update(updateSchema)
	assert.Error(t, err)
}
