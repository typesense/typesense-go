package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/v-byte-cpu/typesense-go/typesense/api"
	"github.com/v-byte-cpu/typesense-go/typesense/api/mocks"
)

func TestCollectionRetrieve(t *testing.T) {
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	mockedResult := &api.Collection{}
	copier.Copy(mockedResult, expectedResult)

	mockApiClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	result, err := client.Collection("companies").Retrieve()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockApiClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockApiClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionDelete(t *testing.T) {
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	mockedResult := &api.Collection{}
	copier.Copy(mockedResult, expectedResult)

	mockApiClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.DeleteCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	result, err := client.Collection("companies").Delete()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockApiClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Delete()
	assert.NotNil(t, err)
}

func TestCollectionDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApiClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockApiClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.DeleteCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithApiClient(mockApiClient))
	_, err := client.Collection("companies").Delete()
	assert.NotNil(t, err)
}
