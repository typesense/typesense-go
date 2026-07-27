package typesense

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

func updateExistingSchema() *api.CollectionUpdateSchema {
	return &api.CollectionUpdateSchema{
		Fields: []api.Field{
			{
				Name: "url",
				Drop: pointer.True(),
			},
			{
				Name:  "url",
				Type:  "string",
				Index: pointer.False(),
			},
		},
		Metadata: &map[string]interface{}{
			"revision": 2,
		},
	}
}

func updateExistingCollection() *api.CollectionUpdateSchema {
	return updateExistingSchema()
}

func TestCollectionRetrieve(t *testing.T) {
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollection("companies")

	mockAPIClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Retrieve(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		GetCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.GetCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Retrieve(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionDelete(t *testing.T) {
	expectedResult := createNewCollection("companies")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := createNewCollection("companies")

	mockAPIClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.DeleteCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Delete(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCollectionDeleteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Delete(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionDeleteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		DeleteCollectionWithResponse(gomock.Not(gomock.Nil()), "companies").
		Return(&api.DeleteCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Delete(context.Background())
	assert.NotNil(t, err)
}

func TestCollectionUpdate(t *testing.T) {
	updateSchema := updateExistingSchema()
	expectedResult := updateExistingCollection()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)
	mockedResult := updateExistingCollection()

	mockAPIClient.EXPECT().
		UpdateCollectionWithResponse(gomock.Not(gomock.Nil()), "companies",
			api.UpdateCollectionJSONRequestBody(*updateSchema)).
		Return(&api.UpdateCollectionResponse{
			JSON200: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Collection("companies").Update(context.Background(), updateSchema)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

// A collection update that changes only the metadata must not send a `fields` key at
// all: Typesense rejects both `"fields":null` and `"fields":[]` with a 400. This relies
// on the `omitempty` tag generated for CollectionUpdateSchema.Fields.
func TestCollectionUpdateSchemaOmitsEmptyFields(t *testing.T) {
	tests := []struct {
		name       string
		fields     []api.Field
		wantFields bool
	}{
		{
			name:   "nil fields",
			fields: nil,
		},
		{
			name:   "empty fields",
			fields: []api.Field{},
		},
		{
			name:       "populated fields",
			fields:     []api.Field{{Name: "country", Drop: pointer.True()}},
			wantFields: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(&api.CollectionUpdateSchema{
				Fields:   tt.fields,
				Metadata: &map[string]interface{}{"revision": "2"},
			})
			assert.NoError(t, err)

			payload := map[string]json.RawMessage{}
			assert.NoError(t, json.Unmarshal(body, &payload))

			_, hasFields := payload["fields"]
			assert.Equal(t, tt.wantFields, hasFields, "body: %s", body)
		})
	}
}

func TestCollectionUpdateOnApiClientErrorReturnsError(t *testing.T) {
	updateSchema := updateExistingSchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpdateCollectionWithResponse(gomock.Not(gomock.Nil()), "companies",
			api.UpdateCollectionJSONRequestBody(*updateSchema)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("companies").Update(context.Background(), updateSchema)
	assert.Error(t, err)
}

func TestCollectionUpdateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	updateSchema := updateExistingSchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		UpdateCollectionWithResponse(gomock.Not(gomock.Nil()), "non_existent",
			api.UpdateCollectionJSONRequestBody(*updateSchema)).
		Return(&api.UpdateCollectionResponse{
			HTTPResponse: &http.Response{
				StatusCode: 404,
			},
			Body: []byte("Collection not found"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Collection("non_existent").Update(context.Background(), updateSchema)
	assert.Error(t, err)
}
