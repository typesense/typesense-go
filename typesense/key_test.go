package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func TestKeyRetrieve(t *testing.T) {
	expectedResult := createNewKey(1)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewKey(1)

	mockAPIClient.EXPECT().
		GetKeyWithResponse(gomock.Not(gomock.Nil()), int64(1)).
		Return(&api.GetKeyResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Key(1).Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestKeyRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetKeyWithResponse(gomock.Not(gomock.Nil()), int64(1)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Key(1).Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestKeyRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetKeyWithResponse(gomock.Not(gomock.Nil()), int64(1)).
		Return(&api.GetKeyResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Key(1).Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestKeyDelete(t *testing.T) {
	expectedResult := &api.ApiKeyDeleteResponse{Id: 1}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := &api.ApiKeyDeleteResponse{Id: 1}

	mockAPIClient.EXPECT().
		DeleteKeyWithResponse(gomock.Not(gomock.Nil()), int64(1)).
		Return(&api.DeleteKeyResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Key(1).Delete(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestKeyDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteKeyWithResponse(gomock.Not(gomock.Nil()), int64(1)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Key(1).Delete(context.Background())
	assert.NotNil(t, err)
}

func TestKeyDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteKeyWithResponse(gomock.Not(gomock.Nil()), int64(1)).
		Return(&api.DeleteKeyResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Key(1).Delete(context.Background())
	assert.NotNil(t, err)
}
