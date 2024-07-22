//go:build integration
// +build integration

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricsRetrieve(t *testing.T) {
	result, err := typesenseClient.Metrics().Retrieve(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result["system_memory_total_bytes"])
}
