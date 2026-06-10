package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type OperationsInterface interface {
	// Creates a point-in-time snapshot of a Typesense node's state and data in the specified directory.
	//
	// Creates a point-in-time snapshot of a Typesense node's state and data in the specified directory. You can then backup the snapshot directory that gets created and later restore it as a data directory, as needed.
	//
	// HTTP: POST /operations/snapshot
	//
	// See: https://typesense.org/docs/latest/api/cluster-operations.html
	Snapshot(ctx context.Context, snapshotPath string) (bool, error)
	// Triggers a follower node to initiate the raft voting process, which triggers leader re-election.
	//
	// Triggers a follower node to initiate the raft voting process, which triggers leader re-election. The follower node that you run this operation against will become the new leader, once this command succeeds.
	//
	// HTTP: POST /operations/vote
	//
	// See: https://typesense.org/docs/latest/api/cluster-operations.html
	Vote(ctx context.Context) (bool, error)
}

type operations struct {
	apiClient APIClientInterface
}

// Creates a point-in-time snapshot of a Typesense node's state and data in the specified directory.
//
// Creates a point-in-time snapshot of a Typesense node's state and data in the specified directory. You can then backup the snapshot directory that gets created and later restore it as a data directory, as needed.
//
// HTTP: POST /operations/snapshot
//
// See: https://typesense.org/docs/latest/api/cluster-operations.html
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

// Triggers a follower node to initiate the raft voting process, which triggers leader re-election.
//
// Triggers a follower node to initiate the raft voting process, which triggers leader re-election. The follower node that you run this operation against will become the new leader, once this command succeeds.
//
// HTTP: POST /operations/vote
//
// See: https://typesense.org/docs/latest/api/cluster-operations.html
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
