package typesense

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/typesense/typesense-go/v4/typesense/api/circuit"
)

type APICall struct {
	client               circuit.HTTPRequestDoer
	nearestNode          *Node
	nodes                []Node
	currentNodeIndex     int
	healthcheckInterval  time.Duration
	numRetriesPerRequest int
	retryInterval        time.Duration
}

type Node struct {
	isHealthy           bool
	index               interface{}
	url                 string
	lastAccessTimestamp int64
}

var apiCallTimeNow = time.Now // for test stubbing

const (
	HEALTHY   = true
	UNHEALTHY = false
)

type APICallOption func(*APICall)

func NewAPICall(client circuit.HTTPRequestDoer, config *ClientConfig) *APICall {
	apiCall := &APICall{
		currentNodeIndex:     -1,
		healthcheckInterval:  config.HealthcheckInterval,
		client:               client,
		numRetriesPerRequest: config.NumRetries,
		retryInterval:        config.RetryInterval,
	}

	// default numRetries is the number of nodes (+1 if nearestNode is specified)
	if config.NumRetries == 0 {
		apiCall.numRetriesPerRequest = len(config.Nodes)
		if config.NearestNode != "" {
			apiCall.numRetriesPerRequest++
		}
	}

	apiCall.initializeNodesMetadata(config)

	return apiCall
}

func (a *APICall) Do(req *http.Request) (*http.Response, error) {
	// Default is to not load balance for backward compatibility
	if len(a.nodes) == 0 {
		res, err := a.client.Do(req)
		return res, err
	}

	var lastResponse *http.Response
	var lastError error

	var unreplayable *unreplayableBody
	if req.GetBody == nil && req.Body != nil && req.Body != http.NoBody {
		unreplayable = &unreplayableBody{body: req.Body}
		req.Body = unreplayable
		defer unreplayable.release()
	}

	for numTries := 0; numTries < a.numRetriesPerRequest; numTries++ {
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			req.Body = body
		} else if unreplayable != nil {
			unreplayable.startAttempt()
		}

		node := a.getNextNode()

		replaceRequestHostname(req, node.url)

		response, err := a.client.Do(req)

		// return early if request is aborted
		if errors.Is(err, context.Canceled) {
			return nil, err
		}

		// On a connection error, status 0 or a 5xx, retry with the next node
		if err != nil || response.StatusCode < 1 || response.StatusCode >= 500 {
			lastResponse = response
			lastError = err

			setNodeHealthCheck(node, UNHEALTHY)
			if unreplayable != nil && !unreplayable.resendAfter(err) {
				break
			}
			time.Sleep(a.retryInterval)
			continue
		} else if response.StatusCode >= 1 && response.StatusCode <= 499 {
			// Treat any status code > 0 and < 500 to be an indication that node is healthy
			// We exclude 0 since some clients return 0 when request fails
			setNodeHealthCheck(node, HEALTHY)
			return response, err
		}
	}

	return lastResponse, lastError
}

// unreplayableBody wraps a request body that net/http can't replay. The
// transport closes the body even when it fails to dial, so Close is held back
// until a byte is read, which keeps the body whole for the next node.
type unreplayableBody struct {
	body io.ReadCloser

	mu        sync.Mutex
	read      bool
	released  bool
	closeHeld bool
}

func (b *unreplayableBody) Read(p []byte) (int, error) {
	n, err := b.body.Read(p)
	if n > 0 {
		b.mu.Lock()
		b.read = true
		b.mu.Unlock()
	}
	return n, err
}

func (b *unreplayableBody) Close() error {
	b.mu.Lock()
	if !b.read && !b.released {
		b.closeHeld = true
		b.mu.Unlock()
		return nil
	}
	b.mu.Unlock()
	return b.body.Close()
}

// startAttempt drops a Close held from the previous attempt, because the
// transport sending this attempt closes the body itself.
func (b *unreplayableBody) startAttempt() {
	b.mu.Lock()
	b.closeHeld = false
	b.mu.Unlock()
}

// resendAfter reports whether the body can go to another node after err, which
// is only when the node couldn't be dialed and none of the body was read.
func (b *unreplayableBody) resendAfter(err error) bool {
	var opErr *net.OpError
	if !errors.As(err, &opErr) || opErr.Op != "dial" {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return !b.read
}

// release lets every later Close through and forwards one that was held.
func (b *unreplayableBody) release() {
	b.mu.Lock()
	b.released = true
	held := b.closeHeld
	b.mu.Unlock()
	if held {
		_ = b.body.Close()
	}
}

func (a *APICall) getNextNode() *Node {
	if a.nearestNode != nil && (a.nearestNode.isHealthy || a.nodeDueForHealthcheck(a.nearestNode)) {
		return a.nearestNode
	}

	candidateNode := &a.nodes[0]
	for i := 0; i <= len(a.nodes); i++ {
		a.currentNodeIndex = (a.currentNodeIndex + 1) % len(a.nodes)
		candidateNode = &a.nodes[a.currentNodeIndex]
		if candidateNode.isHealthy || a.nodeDueForHealthcheck(candidateNode) {
			return candidateNode
		}
	}
	// None of the nodes are marked healthy, but some of them could have become healthy since last health check.
	// So we will just return the next node.
	return candidateNode
}

func (a *APICall) initializeNodesMetadata(config *ClientConfig) {
	if config.NearestNode != "" {
		a.nearestNode = &Node{index: "nearestNode", url: config.NearestNode}
		setNodeHealthCheck(a.nearestNode, HEALTHY)
	}
	a.nodes = make([]Node, 0, len(config.Nodes))
	for i, v := range config.Nodes {
		a.nodes = append(a.nodes, Node{isHealthy: true, index: i, url: v, lastAccessTimestamp: apiCallTimeNow().UnixMilli()})
	}
}

func replaceRequestHostname(req *http.Request, urlToReplace string) {
	newURL, _ := url.Parse(urlToReplace)

	req.URL.Scheme = newURL.Scheme
	req.URL.Host = newURL.Host
	req.Host = newURL.Host
}

func setNodeHealthCheck(node *Node, isHealthy bool) {
	node.isHealthy = isHealthy
	node.lastAccessTimestamp = apiCallTimeNow().UnixMilli()
}

func (a *APICall) nodeDueForHealthcheck(node *Node) bool {
	return apiCallTimeNow().UnixMilli()-node.lastAccessTimestamp > a.healthcheckInterval.Milliseconds()
}
