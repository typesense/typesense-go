package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func createNewKeySchema() *api.ApiKeySchema {
	return &api.ApiKeySchema{
		Description: "Search-only key.",
		Actions:     []string{"documents:search"},
		Collections: []string{"companies"},
		ExpiresAt:   pointer.Int64(time.Date(2222, 0, 1, 0, 0, 0, 0, time.UTC).Unix()),
	}
}

func createNewKey(id int64) *api.ApiKey {
	return &api.ApiKey{
		Id: pointer.Int64(id),
	}
}

func TestKeyCreate(t *testing.T) {
	newKey := createNewKeySchema()
	expectedResult := &api.ApiKey{
		Id:    pointer.Int64(1),
		Value: pointer.String("k8pX5hD0793d8YQC5aD1aEPd7VleSuGP"),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := &api.ApiKey{}
	assert.Nil(t, copier.Copy(mockedResult, expectedResult))

	mockAPIClient.EXPECT().
		CreateKeyWithResponse(gomock.Not(gomock.Nil()),
			api.CreateKeyJSONRequestBody(*newKey)).
		Return(&api.CreateKeyResponse{
			JSON201: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Keys().Create(context.Background(), newKey)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestKeyCreateOnApiClientErrorReturnsError(t *testing.T) {
	newKey := createNewKeySchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateKeyWithResponse(gomock.Not(gomock.Nil()),
			api.CreateKeyJSONRequestBody(*newKey)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Create(context.Background(), newKey)
	assert.NotNil(t, err)
}

func TestKeyCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newKey := createNewKeySchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateKeyWithResponse(gomock.Not(gomock.Nil()),
			api.CreateKeyJSONRequestBody(*newKey)).
		Return(&api.CreateKeyResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Create(context.Background(), newKey)
	assert.NotNil(t, err)
}

func TestKeysRetrieve(t *testing.T) {
	expectedResult := []*api.ApiKey{
		createNewKey(1),
		createNewKey(2),
		createNewKey(3),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := []*api.ApiKey{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetKeysWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetKeysResponse{
			JSON200: &api.ApiKeysResponse{
				Keys: mockedResult,
			},
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Keys().Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestKeysRetrieveOnApiClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetKeysWithResponse(gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestKeysRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetKeysWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetKeysResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestKeysGenerateScopedSearchKey(t *testing.T) {
	// setup example from the docs
	searchKey := "RN23GFr1s6jQ9kgSNg2O7fYcAUXU7127"
	scopedSearchKey := "SC9sT0hncHFwTHNFc3U3d3psRDZBUGNXQUViQUdDNmRHSmJFQnNnczJ4VT1STjIzeyJmaWx0ZXJfYnkiOiJjb21wYW55X2lkOjEyNCJ9"

	scopedKey, err := (&keys{}).GenerateScopedSearchKey(searchKey, map[string]interface{}{
		"filter_by": "company_id:124",
	})

	assert.NoError(t, err)
	assert.Equal(t, scopedKey, scopedSearchKey)
}
