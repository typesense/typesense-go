package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type OperationsInterface interface {
	Snapshot(ctx context.Context, snapshotPath string) (bool, error)
	Vote(ctx context.Context) (bool, error)
}

type operations struct {
	apiClient APIClientInterface
}

func (o *operations) Snapshot(ctx context.Context, snapshotPath string) (bool, error) {
	response, err := o.apiClient.TakeSnapshotWithResponse(ctx,
		&api.TakeSnapshotParams{SnapshotPath: snapshotPath})
	if err != nil {
		return false, err
	}
	if response.JSON201 == nil {
		return false, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201.Success, nil
}

func (o *operations) Vote(ctx context.Context) (bool, error) {
	response, err := o.apiClient.VoteWithResponse(ctx)
	if err != nil {
		return false, err
	}
	if response.JSON200 == nil {
		return false, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Success, nil
}
