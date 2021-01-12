package typesense

import (
	"context"
	"time"
)

func (c *Client) Health(timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := c.apiClient.HealthWithResponse(ctx)
	if err != nil {
		return false, err
	}
	if response.JSON200 == nil {
		return false, &httpError{status: response.StatusCode(), body: response.Body}
	}
	return response.JSON200.Ok, nil
}
