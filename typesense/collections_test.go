package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/typesense/mocks"
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

type MockStruct struct {
	Field1 string `typesense:"string"`
}

func (m MockStruct) CollectionName() string {
	return "custom_collection_name"
}

func TestCreateSchemaFromGoStruct(t *testing.T) {
	mockStruct := MockStruct{Field1: "Test"}

	schema, err := CreateSchemaFromGoStruct(mockStruct)
	assert.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Equal(t, "custom_collection_name", schema.Name)
}

func TestCreateCollectionFromStruct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	client := NewClient(WithAPIClient(mockAPIClient))
	mockStruct := MockStruct{Field1: "Test"}

	schema, _ := CreateSchemaFromGoStruct(mockStruct)
	expectedResult := &api.CollectionResponse{Name: "custom_collection_name", NumDocuments: pointer.Int64(0)}

	mockAPIClient.EXPECT().
		CreateCollectionWithResponse(gomock.Not(gomock.Nil()), api.CreateCollectionJSONRequestBody(*schema)).
		Return(&api.CreateCollectionResponse{JSON201: expectedResult}, nil).
		Times(1)

	result, err := client.Collections().CreateCollectionFromStruct(mockStruct)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
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
	result, err := client.Collections().Create(newSchema)

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
	_, err := client.Collections().Create(newSchema)
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
	_, err := client.Collections().Create(newSchema)
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
		GetCollectionsWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetCollectionsResponse{
			JSON200: &mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collections().Retrieve()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionsRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionsWithResponse(gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collections().Retrieve()
	assert.Error(t, err)
}

func TestCollectionsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

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
	assert.Error(t, err)
}
