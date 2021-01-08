// Package circuit implements the Circuit Breaker pattern for http client.
// It will wrap a http request and monitor for
// failures and/or time outs. When a threshold of failures or time outs has been
// reached, future requests will not run. During this state, the
// breaker will allow limited number of http requests to run and, if they are successful,
// will start performing all http requests again.
package circuit
