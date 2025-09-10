package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func createNewSchema(collectionName string) *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name:  "company_name",
				Type:  "string",
				Index: pointer.True(),
			},
			{
				Name:  "num_employees",
				Type:  "int32",
				Index: pointer.True(),
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: pointer.True(),
				Index: pointer.True(),
			},
			{
				Name:  "url",
				Type:  "string",
				Index: pointer.False(),
			},
		},
		DefaultSortingField: pointer.String("num_employees"),
	}
}

func createNewCollection(name string) *api.CollectionResponse {
	return &api.CollectionResponse{
		Name:         name,
		NumDocuments: pointer.Int64(0),
	}
}

func TestCollectionCreate(t *testing.T) {
	newSchema := createNewSchema("companies")
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollection("companies")

	mockAPIClient.EXPECT().
		CreateCollectionWithResponse(gomock.Not(gomock.Nil()),
			api.CreateCollectionJSONRequestBody(*newSchema)).
		Return(&api.CreateCollectionResponse{
			JSON201: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collections().Create(context.Background(), newSchema)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionCreateOnApiClientErrorReturnsError(t *testing.T) {
	newSchema := createNewSchema("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateCollectionWithResponse(gomock.Not(gomock.Nil()),
			api.CreateCollectionJSONRequestBody(*newSchema)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Create(context.Background(), newSchema)
	assert.Error(t, err)
}

func TestCollectionCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newSchema := createNewSchema("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateCollectionWithResponse(gomock.Not(gomock.Nil()),
			api.CreateCollectionJSONRequestBody(*newSchema)).
		Return(&api.CreateCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 409,
			},
			Body: []byte("Collection already exists"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Create(context.Background(), newSchema)
	assert.Error(t, err)
}

func TestCollectionsRetrieve(t *testing.T) {
	expectedResult := []*api.CollectionResponse{
		createNewCollection("collection1"),
		createNewCollection("collection2"),
		createNewCollection("collection3"),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := []*api.CollectionResponse{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
		Return(&api.GetCollectionsResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collections().Retrieve(context.Background(), &api.GetCollectionsParams{})

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionsRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Retrieve(context.Background(), &api.GetCollectionsParams{})
	assert.Error(t, err)
}

func TestCollectionsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).
		Return(&api.GetCollectionsResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Retrieve(context.Background(), &api.GetCollectionsParams{})
	assert.Error(t, err)
}
