package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/typesense/typesense-go/v2/typesense/api/pointer"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"github.com/typesense/typesense-go/v2/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func createNewSearchOverrideSchema() *api.SearchOverrideSchema {
	return &api.SearchOverrideSchema{
		Rule: api.SearchOverrideRule{
			Query: "apple",
			Match: "exact",
		},
		Includes: &[]api.SearchOverrideInclude{
			{
				Id:       "422",
				Position: 1,
			},
			{
				Id:       "54",
				Position: 2,
			},
		},
		Excludes: &[]api.SearchOverrideExclude{
			{
				Id: "287",
			},
		},
	}
}

func createNewSearchOverride(overrideID string) *api.SearchOverride {
	return &api.SearchOverride{
		Id: pointer.String(overrideID),
	}
}

func TestSearchOverrideUpsert(t *testing.T) {
	newSchema := api.UpsertSearchOverrideJSONRequestBody(
		*createNewSearchOverrideSchema())
	expectedResult := createNewSearchOverride("customize-apple")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewSearchOverride("customize-apple")

	mockAPIClient.EXPECT().
		UpsertSearchOverrideWithResponse(gomock.Not(gomock.Nil()),
			"companies", "customize-apple", newSchema).
		Return(&api.UpsertSearchOverrideResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := createNewSearchOverrideSchema()
	result, err := client.Collection("companies").Overrides().Upsert(context.Background(), "customize-apple", body)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSearchOverrideUpsertOnApiClientErrorReturnsError(t *testing.T) {
	newSchema := api.UpsertSearchOverrideJSONRequestBody(
		*createNewSearchOverrideSchema())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpsertSearchOverrideWithResponse(gomock.Not(gomock.Nil()),
			"companies", "customize-apple", newSchema).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := createNewSearchOverrideSchema()
	_, err := client.Collection("companies").Overrides().Upsert(context.Background(), "customize-apple", body)
	assert.NotNil(t, err)
}

func TestSearchOverrideUpsertOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newSchema := api.UpsertSearchOverrideJSONRequestBody(
		*createNewSearchOverrideSchema())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpsertSearchOverrideWithResponse(gomock.Not(gomock.Nil()),
			"companies", "customize-apple", newSchema).
		Return(&api.UpsertSearchOverrideResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	body := createNewSearchOverrideSchema()
	_, err := client.Collection("companies").Overrides().Upsert(context.Background(), "customize-apple", body)
	assert.NotNil(t, err)
}

func TestSearchOverridesRetrieve(t *testing.T) {
	expectedResult := []*api.SearchOverride{
		createNewSearchOverride("customize1"),
		createNewSearchOverride("customize2"),
		createNewSearchOverride("customize3"),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := []*api.SearchOverride{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetSearchOverridesWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetSearchOverridesResponse{
			JSON200: &api.SearchOverridesResponse{
				Overrides: mockedResult,
			},
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Overrides().Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSearchOverridesRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetSearchOverridesWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Overrides().Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestSearchOverridesRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetSearchOverridesWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetSearchOverridesResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Overrides().Retrieve(context.Background())
	assert.NotNil(t, err)
}
