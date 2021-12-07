//go:generate mockgen -destination=mocks/mock_client.go -package=mocks -source client.go

package typesense

import (
	"fmt"
	"net/http"
	"time"

	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/circuit"
)

type APIClientInterface interface {
	api.ClientWithResponsesInterface
	api.ClientInterface
}

type Client struct {
	apiConfig   *ClientConfig
	apiClient   APIClientInterface
	collections CollectionsInterface
	aliases     AliasesInterface
	MultiSearch MultiSearchInterface
}

func (c *Client) Collections() CollectionsInterface {
	return c.collections
}

func (c *Client) Collection(collectionName string) CollectionInterface {
	return &collection{apiClient: c.apiClient, name: collectionName}
}

func (c *Client) Aliases() AliasesInterface {
	return c.aliases
}

func (c *Client) Alias(aliasName string) AliasInterface {
	return &alias{apiClient: c.apiClient, name: aliasName}
}

func (c *Client) Keys() KeysInterface {
	return &keys{apiClient: c.apiClient}
}

func (c *Client) Key(keyID int64) KeyInterface {
	return &key{apiClient: c.apiClient, keyID: keyID}
}

func (c *Client) Operations() OperationsInterface {
	return &operations{apiClient: c.apiClient}
}

type HTTPError struct {
	Status int
	Body   []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("status: %v response: %s", e.Status, string(e.Body))
}

const (
	defaultConnectionTimeout  = 5 * time.Second
	defaultCircuitBreakerName = "typesenseClient"
)

type ClientConfig struct {
	ServerURL                   string
	APIKey                      string
	ConnectionTimeout           time.Duration
	CircuitBreakerName          string
	CircuitBreakerMaxRequests   uint32
	CircuitBreakerInterval      time.Duration
	CircuitBreakerTimeout       time.Duration
	CircuitBreakerReadyToTrip   circuit.GoBreakerReadyToTripFunc
	CircuitBreakerOnStateChange circuit.GoBreakerOnStateChangeFunc
}

type ClientOption func(*Client)

// WithAPIClient sets low-level API client
func WithAPIClient(apiClient APIClientInterface) ClientOption {
	return func(c *Client) {
		c.apiClient = apiClient
	}
}

// WithServer sets the API server URL
func WithServer(serverURL string) ClientOption {
	return func(c *Client) {
		c.apiConfig.ServerURL = serverURL
	}
}

// WithAPIKey sets the API token.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiConfig.APIKey = apiKey
	}
}

// WithConnectionTimeout sets the connection timeout of http client.
// Default value is 5 seconds.
func WithConnectionTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.apiConfig.ConnectionTimeout = timeout
	}
}

// WithCircuitBreakerName sets the name of the CircuitBreaker.
// Default value is "typesenseClient".
func WithCircuitBreakerName(name string) ClientOption {
	return func(c *Client) {
		c.apiConfig.CircuitBreakerName = name
	}
}

// WithCircuitBreakerMaxRequests sets the maximum number of requests allowed to pass
// through when the CircuitBreaker is half-open. If MaxRequests is 0,
// CircuitBreaker allows only 1 request.
// Default value is 50 requests.
func WithCircuitBreakerMaxRequests(maxRequests uint32) ClientOption {
	return func(c *Client) {
		c.apiConfig.CircuitBreakerMaxRequests = maxRequests
	}
}

// WithCircuitBreakerInterval sets the cyclic period of the closed state for CircuitBreaker
// to clear the internal Counts, described in gobreaker documentation. If Interval is 0,
// CircuitBreaker doesn't clear the internal Counts during the closed state.
// Default value is 2 minutes.
func WithCircuitBreakerInterval(interval time.Duration) ClientOption {
	return func(c *Client) {
		c.apiConfig.CircuitBreakerInterval = interval
	}
}

// WithCircuitBreakerTimeout sets the period of the open state, after which the state of
// CircuitBreaker becomes half-open. If Timeout is 0, the timeout value of CircuitBreaker is set
// to 60 seconds.
// Default value is 1 minute.
func WithCircuitBreakerTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.apiConfig.CircuitBreakerTimeout = timeout
	}
}

// WithCircuitBreakerReadyToTrip sets the function that is called with a copy of Counts
// whenever a request fails in the closed state.
// If ReadyToTrip returns true, CircuitBreaker will be placed into the open state.
// If ReadyToTrip is nil, default ReadyToTrip is used. Default ReadyToTrip returns true when
// number of requests more than 100 and the percent of failures is more than 50 percents.
func WithCircuitBreakerReadyToTrip(readyToTrip circuit.GoBreakerReadyToTripFunc) ClientOption {
	return func(c *Client) {
		c.apiConfig.CircuitBreakerReadyToTrip = readyToTrip
	}
}

// WithCircuitBreakerOnStateChange sets the function that is called whenever
// the state of CircuitBreaker changes.
func WithCircuitBreakerOnStateChange(onStateChange circuit.GoBreakerOnStateChangeFunc) ClientOption {
	return func(c *Client) {
		c.apiConfig.CircuitBreakerOnStateChange = onStateChange
	}
}

// WithClientConfig allows to pass all configs at once
func WithClientConfig(config *ClientConfig) ClientOption {
	return func(c *Client) {
		c.apiConfig.ServerURL = config.ServerURL
		c.apiConfig.APIKey = config.APIKey
		c.apiConfig.ConnectionTimeout = config.ConnectionTimeout
		c.apiConfig.CircuitBreakerName = config.CircuitBreakerName
		c.apiConfig.CircuitBreakerMaxRequests = config.CircuitBreakerMaxRequests
		c.apiConfig.CircuitBreakerInterval = config.CircuitBreakerInterval
		c.apiConfig.CircuitBreakerTimeout = config.CircuitBreakerTimeout
		c.apiConfig.CircuitBreakerReadyToTrip = config.CircuitBreakerReadyToTrip
		c.apiConfig.CircuitBreakerOnStateChange = config.CircuitBreakerOnStateChange
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{apiConfig: &ClientConfig{
		ConnectionTimeout:         defaultConnectionTimeout,
		CircuitBreakerName:        defaultCircuitBreakerName,
		CircuitBreakerMaxRequests: circuit.DefaultGoBreakerMaxRequests,
		CircuitBreakerInterval:    circuit.DefaultGoBreakerInterval,
		CircuitBreakerTimeout:     circuit.DefaultGoBreakerTimeout,
		CircuitBreakerReadyToTrip: circuit.DefaultReadyToTrip,
	}}
	// implement option pattern
	for _, opt := range opts {
		opt(c)
	}
	if c.apiClient == nil {
		cb := circuit.NewGoBreaker(
			circuit.WithGoBreakerName(c.apiConfig.CircuitBreakerName),
			circuit.WithGoBreakerMaxRequests(c.apiConfig.CircuitBreakerMaxRequests),
			circuit.WithGoBreakerInterval(c.apiConfig.CircuitBreakerInterval),
			circuit.WithGoBreakerTimeout(c.apiConfig.CircuitBreakerTimeout),
			circuit.WithGoBreakerReadyToTrip(c.apiConfig.CircuitBreakerReadyToTrip),
			circuit.WithGoBreakerOnStateChange(c.apiConfig.CircuitBreakerOnStateChange),
		)
		httpClient := circuit.NewHTTPClient(
			circuit.WithHTTPRequestDoer(&http.Client{
				Timeout: c.apiConfig.ConnectionTimeout,
			}),
			circuit.WithCircuitBreaker(cb),
		)
		apiClient, _ := api.NewClientWithResponses(c.apiConfig.ServerURL,
			api.WithAPIKey(c.apiConfig.APIKey),
			api.WithHTTPClient(httpClient))
		c.apiClient = apiClient
	}
	c.collections = &collections{c.apiClient}
	c.aliases = &aliases{c.apiClient}
	c.MultiSearch = &multiSearch{c.apiClient}
	return c
}
