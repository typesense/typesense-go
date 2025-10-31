package typesense

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"github.com/typesense/typesense-go/v4/typesense/api/circuit"
)

func TestHttpError(t *testing.T) {
	err := &HTTPError{Status: 200, Body: []byte("error message body")}
	assert.Equal(t, "status: 200 response: error message body", err.Error())
}

func getAPIClient(t *testing.T, apiClient APIClientInterface) *api.Client {
	t.Helper()
	assert.NotNil(t, apiClient)
	assert.IsType(t, &api.ClientWithResponses{}, apiClient)
	clientWithResponses := apiClient.(*api.ClientWithResponses)
	assert.IsType(t, &api.Client{}, clientWithResponses.ClientInterface)
	client := clientWithResponses.ClientInterface.(*api.Client)
	return client
}

func TestClientConfigOptions(t *testing.T) {
	readyToTrip := func(counts gobreaker.Counts) bool {
		return counts.Requests > 10 &&
			(float64(counts.TotalFailures)/float64(counts.Requests)) > 0.4
	}
	onStateChange := func(_ string, _ gobreaker.State, _ gobreaker.State) {}
	tests := []struct {
		name    string
		options []ClientOption
		verify  func(t *testing.T, client *Client)
	}{
		{
			name:    "WithDefaultConfig",
			options: []ClientOption{},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, defaultRetryInterval, client.apiConfig.RetryInterval)
				assert.Equal(t, defaultHealthcheckInterval, client.apiConfig.HealthcheckInterval)
				assert.Equal(t, defaultConnectionTimeout, client.apiConfig.ConnectionTimeout)
				assert.Equal(t, defaultCircuitBreakerName, client.apiConfig.CircuitBreakerName)
				assert.Equal(t, circuit.DefaultGoBreakerMaxRequests, client.apiConfig.CircuitBreakerMaxRequests)
				assert.Equal(t, circuit.DefaultGoBreakerInterval, client.apiConfig.CircuitBreakerInterval)
				assert.Equal(t, circuit.DefaultGoBreakerTimeout, client.apiConfig.CircuitBreakerTimeout)
				assert.Equal(t,
					reflect.ValueOf(circuit.DefaultReadyToTrip).Pointer(),
					reflect.ValueOf(client.apiConfig.CircuitBreakerReadyToTrip).Pointer(),
					"readyToTrip is not valid")
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithServer",
			options: []ClientOption{
				WithServer("http://example.com"),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, "http://example.com", client.apiConfig.ServerURL)
				apiClient := getAPIClient(t, client.apiClient)
				assert.Equal(t, "http://example.com/", apiClient.Server)
			},
		},
		{
			name: "WithNearestNode",
			options: []ClientOption{
				WithNearestNode("http://localhost:8108"),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, "http://localhost:8108", client.apiConfig.NearestNode)
				apiClient := getAPIClient(t, client.apiClient)
				assert.Equal(t, "http://localhost:8108/", apiClient.Server)
			},
		},
		{
			name: "WithNodes",
			options: []ClientOption{
				WithNodes([]string{"http://localhost:3000", "http://localhost:3001"}),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, []string{"http://localhost:3000", "http://localhost:3001"}, client.apiConfig.Nodes)
				apiClient := getAPIClient(t, client.apiClient)
				assert.Equal(t, "http://localhost:3000/", apiClient.Server)
			},
		},
		{
			name: "WithNumRetries",
			options: []ClientOption{
				WithNumRetries(10),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, 10, client.apiConfig.NumRetries)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithRetryInterval",
			options: []ClientOption{
				WithRetryInterval(10 * time.Second),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, 10*time.Second, client.apiConfig.RetryInterval)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithHealthcheckInterval",
			options: []ClientOption{
				WithHealthcheckInterval(10 * time.Second),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, 10*time.Second, client.apiConfig.HealthcheckInterval)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithAPIKey",
			options: []ClientOption{
				WithAPIKey("API_KEY"),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, "API_KEY", client.apiConfig.APIKey)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithConnectionTimeout",
			options: []ClientOption{
				WithConnectionTimeout(10 * time.Second),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, 10*time.Second, client.apiConfig.ConnectionTimeout)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCircuitBreakerName",
			options: []ClientOption{
				WithCircuitBreakerName("typesense"),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, "typesense", client.apiConfig.CircuitBreakerName)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCircuitBreakerMaxRequests",
			options: []ClientOption{
				WithCircuitBreakerMaxRequests(100),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, uint32(100), client.apiConfig.CircuitBreakerMaxRequests)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCircuitBreakerInterval",
			options: []ClientOption{
				WithCircuitBreakerInterval(30 * time.Second),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, 30*time.Second, client.apiConfig.CircuitBreakerInterval)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCircuitBreakerTimeout",
			options: []ClientOption{
				WithCircuitBreakerTimeout(45 * time.Second),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, 45*time.Second, client.apiConfig.CircuitBreakerTimeout)
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCircuitBreakerReadyToTrip",
			options: []ClientOption{
				WithCircuitBreakerReadyToTrip(readyToTrip),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t,
					reflect.ValueOf(readyToTrip).Pointer(),
					reflect.ValueOf(client.apiConfig.CircuitBreakerReadyToTrip).Pointer(),
					"readyToTrip is not valid")
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCircuitBreakerOnStateChange",
			options: []ClientOption{
				WithCircuitBreakerOnStateChange(onStateChange),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t,
					reflect.ValueOf(onStateChange).Pointer(),
					reflect.ValueOf(client.apiConfig.CircuitBreakerOnStateChange).Pointer(),
					"onStateChange is not valid")
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithConfig",
			options: []ClientOption{
				WithClientConfig(&ClientConfig{
					ServerURL:                   "http://example.com",
					APIKey:                      "API_KEY_1",
					ConnectionTimeout:           5 * time.Second,
					CircuitBreakerName:          "typesense_2",
					CircuitBreakerMaxRequests:   100,
					CircuitBreakerInterval:      30 * time.Second,
					CircuitBreakerTimeout:       45 * time.Second,
					CircuitBreakerReadyToTrip:   readyToTrip,
					CircuitBreakerOnStateChange: onStateChange,
				}),
			},
			verify: func(t *testing.T, client *Client) {
				assert.Equal(t, "http://example.com", client.apiConfig.ServerURL)
				assert.Equal(t, "API_KEY_1", client.apiConfig.APIKey)
				assert.Equal(t, 5*time.Second, client.apiConfig.ConnectionTimeout)
				assert.Equal(t, "typesense_2", client.apiConfig.CircuitBreakerName)
				assert.Equal(t, uint32(100), client.apiConfig.CircuitBreakerMaxRequests)
				assert.Equal(t, 30*time.Second, client.apiConfig.CircuitBreakerInterval)
				assert.Equal(t, 45*time.Second, client.apiConfig.CircuitBreakerTimeout)
				assert.Equal(t,
					reflect.ValueOf(readyToTrip).Pointer(),
					reflect.ValueOf(client.apiConfig.CircuitBreakerReadyToTrip).Pointer(),
					"readyToTrip is not valid")
				assert.Equal(t,
					reflect.ValueOf(onStateChange).Pointer(),
					reflect.ValueOf(client.apiConfig.CircuitBreakerOnStateChange).Pointer(),
					"onStateChange is not valid")
				assert.NotNil(t, client.apiClient)
			},
		},
		{
			name: "WithCustomHTTPClient",
			options: []ClientOption{
				WithCustomHTTPClient(&http.Client{
					Timeout: 10 * time.Second,
				}),
			},
			verify: func(t *testing.T, client *Client) {
				assert.NotNil(t, client.apiConfig.CustomHTTPClient)
				assert.Equal(t, 10*time.Second, client.apiConfig.CustomHTTPClient.Timeout)
				assert.NotNil(t, client.apiClient)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.options...)
			tt.verify(t, client)
		})
	}
}
