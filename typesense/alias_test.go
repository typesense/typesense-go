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

func TestCollectionAliasRetrieve(t *testing.T) {
	expectedResult := createNewCollectionAlias("collection", "collection_alias")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollectionAlias("collection", "collection_alias")

	mockAPIClient.EXPECT().
		GetAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(&api.GetAliasResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Alias("collection_alias").Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionAliasRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Alias("collection_alias").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionAliasRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(&api.GetAliasResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Alias("collection_alias").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionAliasDelete(t *testing.T) {
	expectedResult := createNewCollectionAlias("collection", "collection_alias")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollectionAlias("collection", "collection_alias")

	mockAPIClient.EXPECT().
		DeleteAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(&api.DeleteAliasResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Alias("collection_alias").Delete(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionAliasDeleteOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Alias("collection_alias").Delete(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionAliasDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(&api.DeleteAliasResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Alias("collection_alias").Delete(context.Background())
	assert.NotNil(t, err)
}
