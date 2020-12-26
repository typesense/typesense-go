package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/v-byte-cpu/typesense-go/typesense/api"
	"github.com/v-byte-cpu/typesense-go/typesense/mocks"
)

func createNewSchema(collectionName string) *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name: "company_name",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: true,
			},
		},
		DefaultSortingField: "num_employees",
	}
}

func createNewCollection(name string) *api.Collection {
	return &api.Collection{
		CollectionSchema: *createNewSchema(name),
		NumDocuments:     0,
	}
}

func TestCollectionCreate(t *testing.T) {
	newSchema := createNewSchema("companies")
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)
	mockedResult := createNewCollection("companies")

	mockAPIClient.EXPECT().
		CreateCollectionWithResponse(gomock.Not(gomock.Nil()),
			api.CreateCollectionJSONRequestBody(*newSchema)).
		Return(&api.CreateCollectionResponse{
			JSON201: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collections().Create(newSchema)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionCreateOnApiClientErrorReturnsError(t *testing.T) {
	newSchema := createNewSchema("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateCollectionWithResponse(gomock.Not(gomock.Nil()),
			api.CreateCollectionJSONRequestBody(*newSchema)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Create(newSchema)
	assert.NotNil(t, err)
}

func TestCollectionCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newSchema := createNewSchema("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

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
	_, err := client.Collections().Create(newSchema)
	assert.NotNil(t, err)
}

func TestCollectionsRetrieve(t *testing.T) {
	expectedResult := []*api.Collection{
		createNewCollection("collection1"),
		createNewCollection("collection2"),
		createNewCollection("collection3"),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)
	mockedResult := []*api.Collection{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetCollectionsResponse{
			JSON200: &api.CollectionsResponse{
				Collections: mockedResult,
			},
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collections().Retrieve()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionsRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Retrieve()
	assert.NotNil(t, err)
}

func TestCollectionsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockApiClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetCollectionsResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Retrieve()
	assert.NotNil(t, err)
}
