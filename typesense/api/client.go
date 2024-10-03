package api

import (
	"context"
	"net/http"
)

const APIKeyHeader = "X-TYPESENSE-API-KEY" // #nosec G101

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = []RequestEditorFn{func(_ context.Context, req *http.Request) error {
			req.Header.Add(APIKeyHeader, apiKey)
			return nil
		}}
		return nil
	}
}

// Manually defining this unreferenced schema here instead of disabling oapi-codegen schema pruning

type DocumentIndexParameters struct {
	DirtyValues *DirtyValues `json:"dirty_values,omitempty"`
}
