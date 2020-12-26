package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/v-byte-cpu/typesense-go/typesense/api"
	"github.com/v-byte-cpu/typesense-go/typesense/mocks"
)

func TestCollectionAliasRetrieve(t *testing.T) {
	expectedResult := createNewCollectionAlias("collection", "collection_alias")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)
	mockedResult := createNewCollectionAlias("collection", "collection_alias")

	mockAPIClient.EXPECT().
		GetAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(&api.GetAliasResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Alias("collection_alias").Retrieve()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionAliasRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Alias("collection_alias").Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionAliasRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

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
	_, err := client.Alias("collection_alias").Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionAliasDelete(t *testing.T) {
	expectedResult := createNewCollectionAlias("collection", "collection_alias")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)
	mockedResult := createNewCollectionAlias("collection", "collection_alias")

	mockAPIClient.EXPECT().
		DeleteAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(&api.DeleteAliasResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Alias("collection_alias").Delete()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionAliasDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteAliasWithResponse(gomock.Not(gomock.Nil()), "collection_alias").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Alias("collection_alias").Delete()
	assert.NotNil(t, err)
}

func TestCollectionAliasDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

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
	_, err := client.Alias("collection_alias").Delete()
	assert.NotNil(t, err)
}
