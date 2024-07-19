package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/typesense/typesense-go/typesense/api/pointer"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func createNewSearchSynonymSchema() *api.SearchSynonymSchema {
	return &api.SearchSynonymSchema{
		Synonyms: []string{"blazer", "coat", "jacket"},
	}
}

func createNewSearchSynonym(synonymID string) *api.SearchSynonym {
	return &api.SearchSynonym{
		Id: pointer.String(synonymID),
	}
}

func TestSearchSynonymUpsert(t *testing.T) {
	newSchema := api.UpsertSearchSynonymJSONRequestBody(
		*createNewSearchSynonymSchema())
	expectedResult := createNewSearchSynonym("coat-synonyms")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewSearchSynonym("coat-synonyms")

	mockAPIClient.EXPECT().
		UpsertSearchSynonymWithResponse(gomock.Not(gomock.Nil()),
			"products", "coat-synonyms", newSchema).
		Return(&api.UpsertSearchSynonymResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := createNewSearchSynonymSchema()
	result, err := client.Collection("products").Synonyms().Upsert(context.Background(), "coat-synonyms", body)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSearchSynonymUpsertOnApiClientErrorReturnsError(t *testing.T) {
	newSchema := api.UpsertSearchSynonymJSONRequestBody(
		*createNewSearchSynonymSchema())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpsertSearchSynonymWithResponse(gomock.Not(gomock.Nil()),
			"products", "coat-synonyms", newSchema).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := createNewSearchSynonymSchema()
	_, err := client.Collection("products").Synonyms().Upsert(context.Background(), "coat-synonyms", body)
	assert.NotNil(t, err)
}

func TestSearchSynonymUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newSchema := api.UpsertSearchSynonymJSONRequestBody(
		*createNewSearchSynonymSchema())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpsertSearchSynonymWithResponse(gomock.Not(gomock.Nil()),
			"products", "coat-synonyms", newSchema).
		Return(&api.UpsertSearchSynonymResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := createNewSearchSynonymSchema()
	_, err := client.Collection("products").Synonyms().Upsert(context.Background(), "coat-synonyms", body)
	assert.NotNil(t, err)
}

func TestSearchSynonymsRetrieve(t *testing.T) {
	expectedResult := []*api.SearchSynonym{
		createNewSearchSynonym("customize1"),
		createNewSearchSynonym("customize2"),
		createNewSearchSynonym("customize3"),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := []*api.SearchSynonym{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetSearchSynonymsWithResponse(gomock.Not(gomock.Nil()), "products").
		Return(&api.GetSearchSynonymsResponse{
			JSON200: &api.SearchSynonymsResponse{
				Synonyms: mockedResult,
			},
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("products").Synonyms().Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSearchSynonymsRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetSearchSynonymsWithResponse(gomock.Not(gomock.Nil()), "products").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("products").Synonyms().Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestSearchSynonymsRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetSearchSynonymsWithResponse(gomock.Not(gomock.Nil()), "products").
		Return(&api.GetSearchSynonymsResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("products").Synonyms().Retrieve(context.Background())
	assert.NotNil(t, err)
}
