package typesense

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"
)

type OperationsInterface interface {
	Snapshot(snapshotPath string) (bool, error)
	Vote() (bool, error)
}

type operations struct {
	apiClient APIClientInterface
}

func (o *operations) Snapshot(snapshotPath string) (bool, error) {
	response, err := o.apiClient.TakeSnapshotWithResponse(context.Background(),
		&api.TakeSnapshotParams{SnapshotPath: snapshotPath})
	if err != nil {
		return false, err
	}
	if response.JSON201 == nil {
		return false, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201.Success, nil
}

func (o *operations) Vote() (bool, error) {
	response, err := o.apiClient.VoteWithResponse(context.Background())
	if err != nil {
		return false, err
	}
	if response.JSON200 == nil {
		return false, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Success, nil
}
