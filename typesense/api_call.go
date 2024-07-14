package typesense

import (
	"net/http"
	"net/url"
	"time"

	"github.com/typesense/typesense-go/typesense/api/circuit"
)

type ApiCall struct {
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

type ApiCallOption func(*ApiCall)

func NewApiCall(client circuit.HTTPRequestDoer, config *ClientConfig) *ApiCall {
	apiCall := &ApiCall{
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

func (a *ApiCall) Do(req *http.Request) (*http.Response, error) {
	// Default is to not load balance for backward compatibility
	if len(a.nodes) == 0 {
		res, err := a.client.Do(req)
		return res, err
	}

	var lastResponse *http.Response
	var lastError error

	for numTries := 0; numTries < a.numRetriesPerRequest; numTries++ {
		node := a.getNextNode()

		replaceRequestHostname(req, node.url)

		response, err := a.client.Do(req)

		// If conection timeouts or status 5xx, retry
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

func (a *ApiCall) getNextNode() *Node {
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

func (a *ApiCall) initializeNodesMetadata(config *ClientConfig) {
	if config.NearestNode != "" {
		a.nearestNode = &Node{index: "nearestNode", url: config.NearestNode}
		setNodeHealthCheck(a.nearestNode, HEALTHY)
	}
	a.nodes = make([]Node, 0, len(config.Nodes))
	for i, v := range config.Nodes {
		a.nodes = append(a.nodes, Node{isHealthy: true, index: i, url: v, lastAccessTimestamp: apiCallTimeNow().UnixMilli()})
	}
}

func replaceRequestHostname(req *http.Request, URL string) {
	newURL, _ := url.Parse(URL)

	req.URL.Scheme = newURL.Scheme
	req.URL.Host = newURL.Host
	req.Host = newURL.Host
}

func setNodeHealthCheck(node *Node, isHealthy bool) {
	node.isHealthy = isHealthy
	node.lastAccessTimestamp = apiCallTimeNow().UnixMilli()
}

func (a *ApiCall) nodeDueForHealthcheck(node *Node) bool {
	return apiCallTimeNow().UnixMilli()-node.lastAccessTimestamp > a.healthcheckInterval.Milliseconds()
}
