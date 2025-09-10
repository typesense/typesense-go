package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/typesense/typesense-go/v4/typesense/api/pointer"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func createNewCollectionAlias(collectionName string, name string) *api.CollectionAlias {
	return &api.CollectionAlias{
		CollectionName: collectionName,
		Name:           pointer.String(name),
	}
}

func TestCollectionAliasUpsert(t *testing.T) {
	newSchema := api.UpsertAliasJSONRequestBody(
		api.CollectionAliasSchema{CollectionName: "companies"})
	expectedResult := createNewCollectionAlias("companies", "companies_alias")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollectionAlias("companies", "companies_alias")

	mockAPIClient.EXPECT().
		UpsertAliasWithResponse(gomock.Not(gomock.Nil()), "companies_alias", newSchema).
		Return(&api.UpsertAliasResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := &api.CollectionAliasSchema{CollectionName: "companies"}
	result, err := client.Aliases().Upsert(context.Background(), "companies_alias", body)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionAliasUpsertOnApiClientErrorReturnsError(t *testing.T) {
	newSchema := api.UpsertAliasJSONRequestBody(
		api.CollectionAliasSchema{CollectionName: "companies"})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpsertAliasWithResponse(gomock.Not(gomock.Nil()), "companies_alias", newSchema).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := &api.CollectionAliasSchema{CollectionName: "companies"}
	_, err := client.Aliases().Upsert(context.Background(), "companies_alias", body)
	assert.NotNil(t, err)
}

func TestCollectionAliasUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newSchema := api.UpsertAliasJSONRequestBody(
		api.CollectionAliasSchema{CollectionName: "companies"})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpsertAliasWithResponse(gomock.Not(gomock.Nil()), "companies_alias", newSchema).
		Return(&api.UpsertAliasResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := &api.CollectionAliasSchema{CollectionName: "companies"}
	_, err := client.Aliases().Upsert(context.Background(), "companies_alias", body)
	assert.NotNil(t, err)
}

func TestCollectionAliasesRetrieve(t *testing.T) {
	expectedResult := []*api.CollectionAlias{
		createNewCollectionAlias("collection", "collection_alias1"),
		createNewCollectionAlias("collection", "collection_alias2"),
		createNewCollectionAlias("collection", "collection_alias3"),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := []*api.CollectionAlias{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetAliasesWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetAliasesResponse{
			JSON200: &api.CollectionAliasesResponse{
				Aliases: mockedResult,
			},
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Aliases().Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionAliasesRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetAliasesWithResponse(gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Aliases().Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionAliasesRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetAliasesWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetAliasesResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Aliases().Retrieve(context.Background())
	assert.NotNil(t, err)
}
