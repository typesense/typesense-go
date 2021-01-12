// +build integration

package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthStatus(t *testing.T) {
	t.Parallel()
	healthy, err := typesenseClient.Health(2 * time.Second)
	assert.NoError(t, err)
	assert.True(t, healthy)
}
