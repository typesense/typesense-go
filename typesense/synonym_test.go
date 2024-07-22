package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/typesense/typesense-go/v2/typesense/api/pointer"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"github.com/typesense/typesense-go/v2/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func TestSearchSynonymRetrieve(t *testing.T) {
	expectedResult := createNewSearchSynonym("customize-apple")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewSearchSynonym("customize-apple")

	mockAPIClient.EXPECT().
		GetSearchSynonymWithResponse(gomock.Not(gomock.Nil()), "products", "customize-apple").
		Return(&api.GetSearchSynonymResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("products").Synonym("customize-apple").Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSearchSynonymRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetSearchSynonymWithResponse(gomock.Not(gomock.Nil()), "products", "customize-apple").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("products").Synonym("customize-apple").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestSearchSynonymRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetSearchSynonymWithResponse(gomock.Not(gomock.Nil()), "products", "customize-apple").
		Return(&api.GetSearchSynonymResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("products").Synonym("customize-apple").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestSearchSynonymDelete(t *testing.T) {
	expectedResult := &api.SearchSynonym{Id: pointer.String("customize-apple")}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := &api.SearchSynonym{Id: pointer.String("customize-apple")}

	mockAPIClient.EXPECT().
		DeleteSearchSynonymWithResponse(gomock.Not(gomock.Nil()), "products", "customize-apple").
		Return(&api.DeleteSearchSynonymResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("products").Synonym("customize-apple").Delete(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSearchSynonymDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteSearchSynonymWithResponse(gomock.Not(gomock.Nil()), "products", "customize-apple").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("products").Synonym("customize-apple").Delete(context.Background())
	assert.NotNil(t, err)
}

func TestSearchSynonymDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteSearchSynonymWithResponse(gomock.Not(gomock.Nil()), "products", "customize-apple").
		Return(&api.DeleteSearchSynonymResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("products").Synonym("customize-apple").Delete(context.Background())
	assert.NotNil(t, err)
}
