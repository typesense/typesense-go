package circuit

import (
	"time"

	"github.com/sony/gobreaker"
)

type GoBreaker struct {
	cb       *gobreaker.CircuitBreaker
	settings *gobreaker.Settings
}

// assert that GoBreaker implements CircuitBreaker interface
var _ Breaker = (*GoBreaker)(nil)

type GoBreakerOption func(*GoBreaker)

// WithGoBreakerName sets the name of the CircuitBreaker.
// Default value is "GoBreaker".
func WithGoBreakerName(name string) GoBreakerOption {
	return func(breaker *GoBreaker) {
		breaker.settings.Name = name
	}
}

// WithGoBreakerMaxRequests sets the maximum number of requests allowed to pass
// through when the CircuitBreaker is half-open. If MaxRequests is 0,
// CircuitBreaker allows only 1 request.
// Default value is 50 requests.
func WithGoBreakerMaxRequests(maxRequests uint32) GoBreakerOption {
	return func(breaker *GoBreaker) {
		breaker.settings.MaxRequests = maxRequests
	}
}

// WithGoBreakerInterval sets the cyclic period of the closed state for CircuitBreaker
// to clear the internal Counts, described in gobreaker documentation. If Interval is 0,
// CircuitBreaker doesn't clear the internal Counts during the closed state.
// Default value is 2 minutes.
func WithGoBreakerInterval(interval time.Duration) GoBreakerOption {
	return func(breaker *GoBreaker) {
		breaker.settings.Interval = interval
	}
}

// WithGoBreakerTimeout sets the period of the open state, after which the state of
// CircuitBreaker becomes half-open. If Timeout is 0, the timeout value of CircuitBreaker is set
// to 60 seconds.
// Default value is 1 minute.
func WithGoBreakerTimeout(timeout time.Duration) GoBreakerOption {
	return func(breaker *GoBreaker) {
		breaker.settings.Timeout = timeout
	}
}

type GoBreakerReadyToTripFunc func(counts gobreaker.Counts) bool

// WithGoBreakerReadyToTrip sets the function that is called with a copy of Counts
// whenever a request fails in the closed state.
// If ReadyToTrip returns true, CircuitBreaker will be placed into the open state.
// If ReadyToTrip is nil, default ReadyToTrip is used. Default ReadyToTrip returns true when
// number of requests more than 100 and the percent of failures is more than 50 percents.
func WithGoBreakerReadyToTrip(readyToTrip GoBreakerReadyToTripFunc) GoBreakerOption {
	return func(breaker *GoBreaker) {
		breaker.settings.ReadyToTrip = readyToTrip
	}
}

type GoBreakerOnStateChangeFunc func(name string, from gobreaker.State, to gobreaker.State)

// WithGoBreakerOnStateChange sets the function that is called whenever
// the state of CircuitBreaker changes.
func WithGoBreakerOnStateChange(onStateChange GoBreakerOnStateChangeFunc) GoBreakerOption {
	return func(breaker *GoBreaker) {
		breaker.settings.OnStateChange = onStateChange
	}
}

func DefaultReadyToTrip(counts gobreaker.Counts) bool {
	return counts.Requests > 100 &&
		(float64(counts.TotalFailures)/float64(counts.Requests)) > 0.5
}

const (
	DefaultGoBreakerName        = "GoBreaker"
	DefaultGoBreakerMaxRequests = uint32(50)
	DefaultGoBreakerInterval    = 2 * time.Minute
	DefaultGoBreakerTimeout     = 1 * time.Minute
)

func NewGoBreaker(opts ...GoBreakerOption) *GoBreaker {
	gb := &GoBreaker{
		settings: &gobreaker.Settings{
			Name:        DefaultGoBreakerName,
			MaxRequests: DefaultGoBreakerMaxRequests,
			Interval:    DefaultGoBreakerInterval,
			Timeout:     DefaultGoBreakerTimeout,
			ReadyToTrip: DefaultReadyToTrip,
		},
	}
	for _, opt := range opts {
		opt(gb)
	}
	gb.cb = gobreaker.NewCircuitBreaker(*gb.settings)
	return gb
}

func (gb *GoBreaker) Execute(req func() error) error {
	_, err := gb.cb.Execute(func() (interface{}, error) {
		return nil, req()
	})
	return err
}
