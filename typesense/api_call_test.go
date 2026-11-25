package typesense

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestApiCallDoesNotRetryAnUnreplayableBody(t *testing.T) {
	const payload = "body data"

	tests := []struct {
		name string
		body func(t *testing.T) io.Reader
	}{
		{
			name: "reader net/http cannot replay",
			body: func(t *testing.T) io.Reader {
				return struct{ io.Reader }{strings.NewReader(payload)}
			},
		},
		{
			name: "file on disk",
			body: func(t *testing.T) io.Reader {
				t.Helper()
				path := filepath.Join(t.TempDir(), "documents.jsonl")
				require.NoError(t, os.WriteFile(path, []byte(payload), 0o600))
				file, err := os.Open(path)
				require.NoError(t, err)
				t.Cleanup(func() { file.Close() })
				return file
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestURLHistory := make([]string, 0, 2)
			received := make([]string, 0, 2)

			servers, serverURLs := instantiateServers([]serverHandler{
				func(w http.ResponseWriter, r *http.Request) {
					appendHistory(&requestURLHistory, r)
					bodyBytes, _ := io.ReadAll(r.Body)
					received = append(received, string(bodyBytes))
					w.WriteHeader(500)
				},
				func(w http.ResponseWriter, r *http.Request) {
					appendHistory(&requestURLHistory, r)
					bodyBytes, _ := io.ReadAll(r.Body)
					received = append(received, string(bodyBytes))
					w.WriteHeader(201)
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
			req, err := http.NewRequest(http.MethodPost, "http://example.com", test.body(t))
			require.NoError(t, err)
			assert.Nil(t, req.GetBody)

			res, err := apiCall.Do(req)
			require.NoError(t, err)
			require.NotNil(t, res)
			defer res.Body.Close()
			assert.Equal(t, 500, res.StatusCode)
			assert.Equal(t, serverURLs[:1], requestURLHistory)
			assert.Equal(t, []string{payload}, received)
		})
	}
}

func openTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	path := filepath.Join(t.TempDir(), "documents.jsonl")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	file, err := os.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { file.Close() })
	return file
}

func TestApiCallFailsOverAnUnreplayableBodyOnADialError(t *testing.T) {
	const payload = "body data"

	tests := []struct {
		name string
		body func(t *testing.T) io.Reader
	}{
		{
			name: "reader net/http cannot replay",
			body: func(*testing.T) io.Reader { return struct{ io.Reader }{strings.NewReader(payload)} },
		},
		{
			name: "file on disk",
			body: func(t *testing.T) io.Reader { return openTempFile(t, payload) },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestURLHistory := make([]string, 0, 2)
			received := make([]string, 0, 2)

			servers, serverURLs := instantiateServers([]serverHandler{
				func(w http.ResponseWriter, r *http.Request) {
					appendHistory(&requestURLHistory, r)
					w.WriteHeader(500)
				},
				func(w http.ResponseWriter, r *http.Request) {
					appendHistory(&requestURLHistory, r)
					bodyBytes, _ := io.ReadAll(r.Body)
					received = append(received, string(bodyBytes))
					w.WriteHeader(201)
				},
			})
			defer servers[1].Close()
			// Node 1 refuses the connection.
			servers[0].Close()

			apiCall := newAPICall(
				&ClientConfig{
					Nodes:             serverURLs,
					ConnectionTimeout: 5 * time.Second,
				},
			)
			body := test.body(t)
			req, err := http.NewRequest(http.MethodPost, "http://example.com", body)
			require.NoError(t, err)
			assert.Nil(t, req.GetBody)

			res, err := apiCall.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Equal(t, 201, res.StatusCode)
			assert.Equal(t, serverURLs[1:], requestURLHistory)
			assert.Equal(t, []string{payload}, received)
			if file, ok := body.(*os.File); ok {
				_, err := file.Read(make([]byte, 1))
				assert.ErrorIs(t, err, os.ErrClosed)
			}
		})
	}
}

func TestApiCallReturnsTheDialErrorWhenEveryNodeRefusesAnUnreplayableBody(t *testing.T) {
	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(201) },
		func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(201) },
	})
	for _, server := range servers {
		server.Close()
	}

	apiCall := newAPICall(
		&ClientConfig{
			Nodes:             serverURLs,
			ConnectionTimeout: 5 * time.Second,
		},
	)
	file := openTempFile(t, "body data")
	req, err := http.NewRequest(http.MethodPost, "http://example.com", file)
	require.NoError(t, err)

	res, err := apiCall.Do(req)
	assert.Nil(t, res)
	var opErr *net.OpError
	require.ErrorAs(t, err, &opErr)
	assert.Equal(t, "dial", opErr.Op)
	_, err = file.Read(make([]byte, 1))
	assert.ErrorIs(t, err, os.ErrClosed)
}

func TestApiCallDoesNotFailOverAPartlyReadBodyOnADialError(t *testing.T) {
	hosts := make([]string, 0, 2)
	apiCall := NewAPICall(
		doerFunc(func(r *http.Request) (*http.Response, error) {
			hosts = append(hosts, r.URL.Host)
			_, _ = r.Body.Read(make([]byte, 4))
			return nil, &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
		}),
		&ClientConfig{Nodes: []string{"http://node-a", "http://node-b"}},
	)
	req, err := http.NewRequest(http.MethodPost, "http://example.com", struct{ io.Reader }{strings.NewReader("body data")})
	require.NoError(t, err)

	res, err := apiCall.Do(req)
	assert.Nil(t, res)
	var opErr *net.OpError
	require.ErrorAs(t, err, &opErr)
	assert.Equal(t, []string{"node-a"}, hosts)
}

func TestApiCallKeepsAFailedOverBodyOpenUntilTheNodeHasReadIt(t *testing.T) {
	payload := strings.Repeat(`{"id":"1","title":"a document"}`+"\n", 1<<18)
	received := make(chan int, 1)

	servers, serverURLs := instantiateServers([]serverHandler{
		func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(201) },
		func(w http.ResponseWriter, r *http.Request) {
			// Answer before reading the body, so Do returns while it is still being sent.
			require.NoError(t, http.NewResponseController(w).EnableFullDuplex())
			w.WriteHeader(200)
			w.(http.Flusher).Flush()
			n, _ := io.Copy(io.Discard, r.Body)
			received <- int(n)
		},
	})
	defer servers[1].Close()
	// Node 1 refuses the connection.
	servers[0].Close()

	apiCall := newAPICall(
		&ClientConfig{
			Nodes:             serverURLs,
			ConnectionTimeout: 5 * time.Second,
		},
	)
	req, err := http.NewRequest(http.MethodPost, "http://example.com", openTempFile(t, payload))
	require.NoError(t, err)

	res, err := apiCall.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, len(payload), <-received)
}

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func TestApiCallTreatsStatusZeroAsAFailedAttempt(t *testing.T) {
	tests := []struct {
		name      string
		body      io.Reader
		wantHosts []string
	}{
		{name: "replayable body", body: strings.NewReader("body data"), wantHosts: []string{"node-a", "node-b"}},
		{name: "unreplayable body", body: struct{ io.Reader }{strings.NewReader("body data")}, wantHosts: []string{"node-a"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hosts := make([]string, 0, 2)
			apiCall := NewAPICall(
				doerFunc(func(r *http.Request) (*http.Response, error) {
					hosts = append(hosts, r.URL.Host)
					return &http.Response{StatusCode: 0, Body: http.NoBody, Request: r}, nil
				}),
				&ClientConfig{Nodes: []string{"http://node-a", "http://node-b"}},
			)
			req, err := http.NewRequest(http.MethodPost, "http://example.com", test.body)
			require.NoError(t, err)

			res, err := apiCall.Do(req)
			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Equal(t, 0, res.StatusCode)
			assert.Equal(t, test.wantHosts, hosts)
		})
	}
}
