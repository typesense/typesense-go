//go:generate go run go.uber.org/mock/mockgen -destination=mocks/mock_client.go -package=mocks -source client.go

package typesense

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/circuit"
)

type APIClientInterface interface {
	api.ClientWithResponsesInterface
	api.ClientInterface
}

type Client struct {
	apiConfig    *ClientConfig
	apiClient    APIClientInterface
	collections  CollectionsInterface
	aliases      AliasesInterface
	MultiSearch  MultiSearchInterface
	synonymSets  SynonymSetsInterface
	curationSets CurationSetsInterface
}

func (c *Client) Collections() CollectionsInterface {
	return c.collections
}

func GenericCollection[T any](c *Client, collectionName string) CollectionInterface[T] {
	return &collection[T]{apiClient: c.apiClient, name: collectionName}
}

func (c *Client) Collection(collectionName string) CollectionInterface[map[string]any] {
	return GenericCollection[map[string]any](c, collectionName)
}

func (c *Client) Aliases() AliasesInterface {
	return c.aliases
}

func (c *Client) Alias(aliasName string) AliasInterface {
	return &alias{apiClient: c.apiClient, name: aliasName}
}

func (c *Client) Analytics() AnalyticsInterface {
	return &analytics{apiClient: c.apiClient}
}

func (c *Client) Stemming() StemmingInterface {
	return &stemming{apiClient: c.apiClient}
}

func (c *Client) Conversations() ConversationsInterface {
	return &conversations{apiClient: c.apiClient}
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

func (c *Client) Presets() PresetsInterface {
	return &presets{apiClient: c.apiClient}
}

func (c *Client) Preset(presetName string) PresetInterface {
	return &preset{apiClient: c.apiClient, presetName: presetName}
}

func (c *Client) NLSearchModels() NLSearchModelsInterface {
	return &nlSearchModels{apiClient: c.apiClient}
}

func (c *Client) NLSearchModel(modelID string) NLSearchModelInterface {
	return &nlSearchModel{apiClient: c.apiClient, modelID: modelID}
}

func (c *Client) SynonymSets() SynonymSetsInterface {
	return c.synonymSets
}

func (c *Client) SynonymSet(synonymSetName string) SynonymSetInterface {
	return &synonymSet{apiClient: c.apiClient, synonymSetName: synonymSetName}
}

func (c *Client) CurationSets() CurationSetsInterface {
	return c.curationSets
}

func (c *Client) CurationSet(curationSetName string) CurationSetInterface {
	return &curationSet{apiClient: c.apiClient, curationSetName: curationSetName}
}

func (c *Client) Stopwords() StopwordsInterface {
	return &stopwords{apiClient: c.apiClient}
}

func (c *Client) Stopword(stopwordsSetId string) StopwordInterface {
	return &stopword{apiClient: c.apiClient, stopwordsSetId: stopwordsSetId}
}

func (c *Client) Stats() StatsInterface {
	return &stats{apiClient: c.apiClient}
}

func (c *Client) Metrics() MetricsInterface {
	return &metrics{apiClient: c.apiClient}
}

// Debug retrieves debug information from the Typesense server
func (c *Client) Debug(ctx context.Context) (*api.DebugResponse, error) {
	return c.apiClient.DebugWithResponse(ctx)
}

type HTTPError struct {
	Status int
	Body   []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("status: %v response: %s", e.Status, string(e.Body))
}

const (
	defaultRetryInterval       = 100 * time.Millisecond
	defaultHealthcheckInterval = 1 * time.Minute
	defaultConnectionTimeout   = 5 * time.Second
	defaultCircuitBreakerName  = "typesenseClient"
)

type ClientConfig struct {
	ServerURL                   string
	NearestNode                 string // optional
	Nodes                       []string
	NumRetries                  int
	RetryInterval               time.Duration
	HealthcheckInterval         time.Duration
	APIKey                      string
	ConnectionTimeout           time.Duration
	CircuitBreakerName          string
	CircuitBreakerMaxRequests   uint32
	CircuitBreakerInterval      time.Duration
	CircuitBreakerTimeout       time.Duration
	CircuitBreakerReadyToTrip   circuit.GoBreakerReadyToTripFunc
	CircuitBreakerOnStateChange circuit.GoBreakerOnStateChangeFunc
	CustomHTTPClient            *http.Client
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

// WithNearestNode sets the Load Balanced endpoint.
func WithNearestNode(nodeURL string) ClientOption {
	return func(c *Client) {
		c.apiConfig.NearestNode = nodeURL
	}
}

// WithNodes sets multiple hostnames to load balance reads & writes across all nodes.
func WithNodes(nodeURLs []string) ClientOption {
	return func(c *Client) {
		c.apiConfig.Nodes = nodeURLs
	}
}

// WithNumRetries sets the number of retries per request.
// Default value is the number of nodes (+1 if nearestNode is specified).
func WithNumRetries(num int) ClientOption {
	return func(c *Client) {
		c.apiConfig.NumRetries = num
	}
}

// WithRetryInterval sets the wait time between each retry.
// Default value is 100 milliseconds.
func WithRetryInterval(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.apiConfig.RetryInterval = duration
	}
}

// WithHealthcheckInterval sets the wait time for an unhealthy node to become healthy again.
// A node is marked as unhealthy if its response status code is 5xx or has an error (e.g. timeout).
// Default value is 1 minute.
func WithHealthcheckInterval(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.apiConfig.HealthcheckInterval = duration
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
		c.apiConfig.NearestNode = config.NearestNode
		c.apiConfig.Nodes = config.Nodes
		c.apiConfig.NumRetries = config.NumRetries
		c.apiConfig.RetryInterval = config.RetryInterval
		c.apiConfig.HealthcheckInterval = config.HealthcheckInterval
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

func WithCustomHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.apiConfig.CustomHTTPClient = client
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{apiConfig: &ClientConfig{
		RetryInterval:             defaultRetryInterval,
		HealthcheckInterval:       defaultHealthcheckInterval,
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
		client := c.apiConfig.CustomHTTPClient
		if client == nil {
			client = &http.Client{
				Timeout: c.apiConfig.ConnectionTimeout,
			}
		}
		httpClient := circuit.NewHTTPClient(
			circuit.WithHTTPRequestDoer(
				NewAPICall(
					client,
					c.apiConfig,
				)),
			circuit.WithCircuitBreaker(cb),
		)
		serverURL := ""

		switch {
		case c.apiConfig.ServerURL != "":
			serverURL = c.apiConfig.ServerURL
		case c.apiConfig.NearestNode != "":
			serverURL = c.apiConfig.NearestNode
		default:
			if len(c.apiConfig.Nodes) != 0 {
				serverURL = c.apiConfig.Nodes[0]
			}
		}

		apiClient, _ := api.NewClientWithResponses(serverURL,
			api.WithAPIKey(c.apiConfig.APIKey),
			api.WithHTTPClient(httpClient))
		c.apiClient = apiClient
	}
	c.collections = &collections{c.apiClient}
	c.aliases = &aliases{c.apiClient}
	c.MultiSearch = &multiSearch{c.apiClient}
	c.synonymSets = &synonymSets{c.apiClient}
	c.curationSets = &curationSets{c.apiClient}
	return c
}
