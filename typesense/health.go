package typesense

import (
	"context"
	"time"
)

// Checks if Typesense server is ready to accept requests.
//
// HTTP: GET /health
//
// See: https://typesense.org/docs/latest/api/cluster-operations.html#health
func (c *Client) Health(ctx context.Context, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	response, err := c.apiClient.HealthWithResponse(ctx)
	if err != nil {
		return false, err
	}
	if response.JSON200 == nil {
		return false, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Ok, nil
}
