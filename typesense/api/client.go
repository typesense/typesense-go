package api

import (
	"context"
	"net/http"
)

const APIKeyHeader = "X-TYPESENSE-API-KEY"

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = []RequestEditorFn{func(ctx context.Context, req *http.Request) error {
			req.Header.Add(APIKeyHeader, apiKey)
			return nil
		}}
		return nil
	}
}
