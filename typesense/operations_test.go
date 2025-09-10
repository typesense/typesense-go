package typesense

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/mocks"
	"go.uber.org/mock/gomock"
)

const snapshotPath = "/tmp/typesense-data-snapshot"

func TestSnapshot(t *testing.T) {

	tests := []struct {
		ok bool
	}{
		{
			ok: true,
		},
		{
			ok: false,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

		expectedParams := &api.TakeSnapshotParams{SnapshotPath: snapshotPath}
		mockAPIClient.EXPECT().
			TakeSnapshotWithResponse(gomock.Not(gomock.Nil()), expectedParams).
			Return(&api.TakeSnapshotResponse{
				JSON201: &api.SuccessStatus{Success: tt.ok},
			}, nil).
			Times(1)

		client := NewClient(WithAPIClient(mockAPIClient))
		result, err := client.Operations().Snapshot(context.Background(), snapshotPath)
		assert.NoError(t, err)
		assert.Conditionf(t, func() bool {
			return result == tt.ok
		}, "snapshot status expected to be %v", tt.ok)
	}
}

func TestSnapshotOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	expectedParams := &api.TakeSnapshotParams{SnapshotPath: snapshotPath}
	mockAPIClient.EXPECT().
		TakeSnapshotWithResponse(gomock.Not(gomock.Nil()), expectedParams).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Operations().Snapshot(context.Background(), snapshotPath)
	assert.Error(t, err)
	assert.False(t, result)
}

func TestSnapshotOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	expectedParams := &api.TakeSnapshotParams{SnapshotPath: snapshotPath}
	mockAPIClient.EXPECT().
		TakeSnapshotWithResponse(gomock.Not(gomock.Nil()), expectedParams).
		Return(&api.TakeSnapshotResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Operations().Snapshot(context.Background(), snapshotPath)
	assert.Error(t, err)
	assert.False(t, result)
}

func TestVote(t *testing.T) {

	tests := []struct {
		ok bool
	}{
		{
			ok: true,
		},
		{
			ok: false,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

		mockAPIClient.EXPECT().
			VoteWithResponse(gomock.Not(gomock.Nil())).
			Return(&api.VoteResponse{
				JSON200: &api.SuccessStatus{Success: tt.ok},
			}, nil).
			Times(1)

		client := NewClient(WithAPIClient(mockAPIClient))
		result, err := client.Operations().Vote(context.Background())
		assert.NoError(t, err)
		assert.Conditionf(t, func() bool {
			return result == tt.ok
		}, "vote status expected to be %v", tt.ok)
	}
}

func TestVoteOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		VoteWithResponse(gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Operations().Vote(context.Background())
	assert.Error(t, err)
	assert.False(t, result)
}

func TestVoteOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockAPIClientInterface(ctrl)

	mockAPIClient.EXPECT().
		VoteWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.VoteResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Operations().Vote(context.Background())
	assert.Error(t, err)
	assert.False(t, result)
}
