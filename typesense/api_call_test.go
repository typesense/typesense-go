package typesense

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

type serverHandler func(http.ResponseWriter, *http.Request)

func newHTTPRequest(t *testing.T, urls ...string) *http.Request {
	t.Helper()
	url := "http://example.com/collections/1?test=1"
	if len(urls) != 0 {
		url = urls[0]
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	return req
}

func newAPICall(apiConfig *ClientConfig) *APICall {
	return NewAPICall(
		&http.Client{
			Timeout: apiConfig.ConnectionTimeout,
		},
		apiConfig,
	)
}

func appendHistory(history *[]string, r *http.Request) {
	*history = append(*history, "http://"+r.Host)
}

func instantiateServers(handlers []serverHandler) ([]*httptest.Server, []string) {
	servers := make([]*httptest.Server, 0, len(handlers))
	serverURLs := make([]string, 0, len(handlers))
	for _, handler := range handlers {
		server := httptest.NewServer(http.HandlerFunc(handler))
		servers = append(servers, server)
		serverURLs = append(serverURLs, server.URL)
	}
	return servers, serverURLs
}

func freezeUnixMilli(msec int64) {
	apiCallTimeNow = func() time.Time {
		return time.UnixMilli(msec)
	}
}

// * When NearestNode is not specified
func TestApiCallDoesNotRetryWhenStatusCodeIs3xxOr4xx(t *testing.T) {

	requestURLHistory := make([]string, 0, 2)

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(301)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(409)
		},
	})
	for _, server := range servers {
		defer server.Close()
	}

	apiCall := newAPICall(
		&ClientConfig{
			Nodes:             serverURLs,
			ConnectionTimeout: 5 * time.Second,
		},
	)
	req := newHTTPRequest(t)

	res, err := apiCall.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 301, res.StatusCode)

	res2, err2 := apiCall.Do(req)
	assert.NoError(t, err2)
	assert.Equal(t, 409, res2.StatusCode)

	assert.Equal(t, serverURLs, requestURLHistory)
}

type timeoutError struct {
	err     string
	timeout bool
}

func (e *timeoutError) Error() string {
	return e.err
}

func (e *timeoutError) Timeout() bool {
	return e.timeout
}

type mockClient struct{}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	(&http.Client{}).Do(req)
	return nil, &timeoutError{
		err:     "context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
		timeout: true,
	}
}
func TestApiCallSelectNextNodeWhenTimeOut(t *testing.T) {
	requestURLHistory := make([]string, 0, 3)

	servers, serverURLs := instantiateServers([]serverHandler{
		func(_ http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
		},
		func(_ http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
		},
		func(_ http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
		},
	})
	for _, server := range servers {
		defer server.Close()
	}

	apiCall := NewAPICall(
		&mockClient{},
		&ClientConfig{
			Nodes:             serverURLs,
			ConnectionTimeout: 10 * time.Millisecond,
		},
	)
	req := newHTTPRequest(t)

	res, err := apiCall.Do(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, serverURLs, requestURLHistory)
}
func TestApiCallRemoveAndAddUnhealthyNodeIntoRotation(t *testing.T) {
	requestURLHistory := make([]string, 0, 8)
	var count int

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			if count > 1 {
				// will response successful code after 2 failed request
				w.WriteHeader(201)
			} else {
				count++
				w.WriteHeader(500)
			}
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(501)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(202)
		},
	})
	for _, server := range servers {
		defer server.Close()
	}
	defer func() {
		apiCallTimeNow = time.Now
	}()

	apiCall := newAPICall(
		&ClientConfig{
			Nodes:               serverURLs,
			HealthcheckInterval: 20 * time.Millisecond,
			ConnectionTimeout:   5 * time.Second,
		},
	)
	req := newHTTPRequest(t)

	freezeUnixMilli(0)

	res, err := apiCall.Do(req) // node 0 and node 1 will fail, node 2 will succeed

	assert.NoError(t, err)
	assert.Equal(t, 202, res.StatusCode)
	assert.Equal(t, serverURLs, requestURLHistory[:3])

	freezeUnixMilli(10)

	res2, err := apiCall.Do(req) // request should still be made to node 2

	assert.NoError(t, err)
	assert.Equal(t, 202, res2.StatusCode)
	assert.Equal(t, serverURLs[2], requestURLHistory[3])

	freezeUnixMilli(25)

	res3, err := apiCall.Do(req) // node 0 and 1 added back into rotation but still fail

	assert.NoError(t, err)
	assert.Equal(t, 202, res3.StatusCode)
	assert.Equal(t, serverURLs, requestURLHistory[4:7])

	freezeUnixMilli(50)

	res4, err := apiCall.Do(req) // node 0 will succeed

	assert.NoError(t, err)
	assert.Equal(t, 201, res4.StatusCode)
	assert.Equal(t, serverURLs[0], requestURLHistory[len(requestURLHistory)-1])
}

// * When NearestNode is specified
func TestApiCallWithNearestNode(t *testing.T) {
	requestURLHistory := make([]string, 0, 10)
	var count int

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			if count > 1 {
				// will response successful code after 2 failed request
				w.WriteHeader(201)
			} else {
				count++
				w.WriteHeader(501)
			}
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(502)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(503)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(202)
		},
	})
	for _, server := range servers {
		defer server.Close()
	}
	defer func() {
		apiCallTimeNow = time.Now
	}()

	apiCall := newAPICall(
		&ClientConfig{
			NearestNode:         serverURLs[0],
			Nodes:               serverURLs[1:],
			HealthcheckInterval: 20 * time.Millisecond,
			ConnectionTimeout:   5 * time.Second,
		},
	)

	req := newHTTPRequest(t)

	freezeUnixMilli(0)

	res, err := apiCall.Do(req) // nearest node, node 0 and node 1 will fail, node 2 will succeed

	assert.NoError(t, err)
	assert.Equal(t, 202, res.StatusCode)
	assert.Equal(t, serverURLs, requestURLHistory[:4])

	freezeUnixMilli(10)

	res2, err := apiCall.Do(req) // request should still be made to node 2

	assert.NoError(t, err)
	assert.Equal(t, 202, res2.StatusCode)
	assert.Equal(t, serverURLs[3], requestURLHistory[4])

	freezeUnixMilli(25)

	res3, err := apiCall.Do(req) // nearest node, node 0 and 1 added back into rotation but still fail

	assert.NoError(t, err)
	assert.Equal(t, 202, res3.StatusCode)
	assert.Equal(t, serverURLs, requestURLHistory[5:9])

	freezeUnixMilli(50)

	res4, err := apiCall.Do(req) // nearest node will succeed

	assert.NoError(t, err)
	assert.Equal(t, 201, res4.StatusCode)
	assert.Equal(t, serverURLs[0], requestURLHistory[9])

	res5, err := apiCall.Do(req) // request should still be made to nearest node
	assert.NoError(t, err)
	assert.Equal(t, 201, res5.StatusCode)
	assert.Equal(t, serverURLs[0], requestURLHistory[len(requestURLHistory)-1])
}

func TestApiCallCompatibleWithServerURL(t *testing.T) {
	var lastRequestURL string

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, r *http.Request) {
			lastRequestURL = "http://" + r.Host + r.RequestURI
			w.WriteHeader(201)
		},
	})
	defer servers[0].Close()

	apiCall := newAPICall(
		&ClientConfig{
			ServerURL:         serverURLs[0],
			ConnectionTimeout: 5 * time.Second,
		},
	)
	req := newHTTPRequest(t, serverURLs[0]+"/collections/1?test=1")

	res, err := apiCall.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, serverURLs[0]+"/collections/1?test=1", lastRequestURL)
}

func TestApiCallCanReplaceRequestHostName(t *testing.T) {
	var lastRequestURL string

	servers, serverURLs := instantiateServers([]serverHandler{
		func(_ http.ResponseWriter, r *http.Request) {
			lastRequestURL = "http://" + r.Host + r.RequestURI
		},
	})
	defer servers[0].Close()

	apiCall := newAPICall(
		&ClientConfig{
			Nodes:             serverURLs,
			ConnectionTimeout: 5 * time.Second,
		},
	)
	req := newHTTPRequest(t)

	res, err := apiCall.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	assert.Equal(t, serverURLs[0]+"/collections/1?test=1", lastRequestURL)
}

func TestApiCallCanAbortRequest(t *testing.T) {
	requestURLHistory := make([]string, 0)
	ctx, cancel := context.WithCancel(context.Background())

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(http.StatusBadGateway)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(http.StatusBadGateway)
		},
		func(_ http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			cancel()
			// block until the client closes the connection
			<-r.Context().Done()
		},
		func(_ http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
		},
	})
	for _, server := range servers {
		defer server.Close()
	}

	client := NewClient(
		WithNearestNode(serverURLs[0]),
		WithNodes(serverURLs[1:]),
		WithRetryInterval(0),
	)

	res, err := client.Collections().Retrieve(ctx, &api.GetCollectionsParams{})

	assert.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, res)
	assert.Equal(t, requestURLHistory, serverURLs[:3])
}

func TestApiCallRetryWithRequestBody(t *testing.T) {

	requestURLHistory := make([]string, 0, 2)

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			w.WriteHeader(501)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			data := r.Body
			bodyBytes, _ := io.ReadAll(data)
			assert.Equal(t, string(bodyBytes), "body data")
			w.WriteHeader(201)
		},
		func(w http.ResponseWriter, r *http.Request) {
			appendHistory(&requestURLHistory, r)
			data := r.Body
			bodyBytes, _ := io.ReadAll(data)
			assert.Equal(t, string(bodyBytes), "body data")
			w.WriteHeader(203)
		},
	})
	for _, server := range servers {
		defer server.Close()
	}

	apiCall := newAPICall(
		&ClientConfig{
			Nodes:             serverURLs,
			ConnectionTimeout: 5 * time.Second,
		},
	)
	req, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer([]byte("body data")))
	assert.NoError(t, err)

	res, err := apiCall.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)

	res2, err2 := apiCall.Do(req)
	assert.NoError(t, err2)
	assert.Equal(t, 203, res2.StatusCode)

	assert.Equal(t, serverURLs, requestURLHistory)
}
