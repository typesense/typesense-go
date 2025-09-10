//go:generate go run go.uber.org/mock/mockgen -source=http_client.go -destination=mocks/mock_circuit.go -package=mocks

package circuit

import (
	"net/http"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Breaker defines contract for circuit breakers to implement
type Breaker interface {
	Execute(req func() error) error
}

type HTTPClient struct {
	client  HTTPRequestDoer
	breaker Breaker
}

// assert that HttpClient implements api.HttpRequestDoer interface
var _ api.HttpRequestDoer = (*HTTPClient)(nil)

type ClientOption func(*HTTPClient)

func WithHTTPRequestDoer(client HTTPRequestDoer) ClientOption {
	return func(c *HTTPClient) {
		c.client = client
	}
}

func WithCircuitBreaker(cb Breaker) ClientOption {
	return func(c *HTTPClient) {
		c.breaker = cb
	}
}

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	client := &HTTPClient{}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (hc *HTTPClient) Do(req *http.Request) (response *http.Response, err error) {
	err = hc.breaker.Execute(func() (err error) {
		response, err = hc.client.Do(req)
		return
	})
	return
}
