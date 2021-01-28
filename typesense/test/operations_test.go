// +build integration,docker

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshot(t *testing.T) {
	t.Skip("snapshot blocks other write operations, fix it")
	snapshotPath := newUUIDName("/tmp/typesense-data-snapshot")
	success, err := typesenseClient.Operations().Snapshot(snapshotPath)
	assert.NoError(t, err)
	assert.True(t, success)
}

func TestVote(t *testing.T) {
	success, err := typesenseClient.Operations().Vote()
	assert.NoError(t, err)
	assert.False(t, success)
}
