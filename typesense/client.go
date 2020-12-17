package typesense

import (
	"fmt"

	"github.com/v-byte-cpu/typesense-go/typesense/api"
)

type Client struct {
	apiClient   api.ClientWithResponsesInterface
	Collections CollectionsInterface
}

// Collections('name') -> new object (or from map)
// create internal map for collections
// Documents <- struct field (interface)

// client.Collections('name').Retrieve()
// client.Collections('name').Delete()
// client.Collections('name').Documents.Create(document)
// client.Collections('name').Documents('124').Retrieve()

type httpError struct {
	status int
	body   []byte
}

func (e *httpError) Error() string {
	return fmt.Sprintf("status: %v response: %s", e.status, string(e.body))
}

type ClientOption func(*Client)

func WithApiClient(apiClient api.ClientWithResponsesInterface) ClientOption {
	return func(c *Client) {
		c.apiClient = apiClient
	}
}

// TODO WithServer option (server string)
// TODO WithConnectionTimeout option (seconds int)
// TODO WithApiKey option (apiKey string)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}
	//implement option pattern
	for _, opt := range opts {
		opt(c)
	}
	c.Collections = &collections{c.apiClient}
	return c
}
