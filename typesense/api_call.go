package typesense

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
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
	var bodyBytes []byte

	if req.GetBody != nil {
		// Store body in case we need to retry
		reqBody, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		defer reqBody.Close()

		bodyBytes, err = io.ReadAll(reqBody)
		if err != nil {
			return nil, err
		}
	}

	for numTries := 0; numTries < a.numRetriesPerRequest; numTries++ {
		node := a.getNextNode()

		replaceRequestHostname(req, node.url)

		if bodyBytes != nil {
			// Create a new io.ReadCloser for each retry
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		response, err := a.client.Do(req)

		// return early if request is aborted
		if errors.Is(err, context.Canceled) {
			return nil, err
		}

		// If connection timeouts or status 5xx, retry with the next node
		if err != nil || response.StatusCode >= 500 {
			lastResponse = response
			lastError = err

			setNodeHealthCheck(node, UNHEALTHY)
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
