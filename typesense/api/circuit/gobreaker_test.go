package circuit

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestGoBreakerExecute(t *testing.T) {
	gb := NewGoBreaker()
	i := 0
	err := gb.Execute(func() error {
		i++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}

func TestGoBreakerExecuteReturnsError(t *testing.T) {
	gb := NewGoBreaker()
	err := gb.Execute(func() error {
		return errors.New("execute error")
	})
	assert.Error(t, err)
}

func TestGoBreakerSettingsOptions(t *testing.T) {
	readyToTrip := func(counts gobreaker.Counts) bool {
		return counts.Requests > 10 &&
			(float64(counts.TotalFailures)/float64(counts.Requests)) > 0.4
	}
	onStateChange := func(_ string, _ gobreaker.State, _ gobreaker.State) {}
	tests := []struct {
		name    string
		options []GoBreakerOption
		verify  func(t *testing.T, breaker *GoBreaker)
	}{
		{
			name:    "WithDefaultSettings",
			options: []GoBreakerOption{},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t, DefaultGoBreakerName, breaker.settings.Name)
				assert.Equal(t, DefaultGoBreakerMaxRequests, breaker.settings.MaxRequests)
				assert.Equal(t, DefaultGoBreakerInterval, breaker.settings.Interval)
				assert.Equal(t, DefaultGoBreakerTimeout, breaker.settings.Timeout)
				assert.Equal(t,
					reflect.ValueOf(DefaultReadyToTrip).Pointer(),
					reflect.ValueOf(breaker.settings.ReadyToTrip).Pointer(),
					"readyToTrip function is not default",
				)
				assert.NotNil(t, breaker.cb)
			},
		},
		{
			name: "WithGoBreakerName",
			options: []GoBreakerOption{
				WithGoBreakerName("goBreakerName"),
			},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t, "goBreakerName", breaker.settings.Name)
				assert.NotNil(t, breaker.cb)
			},
		},
		{
			name: "WithGoBreakerMaxRequests",
			options: []GoBreakerOption{
				WithGoBreakerMaxRequests(100),
			},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t, uint32(100), breaker.settings.MaxRequests)
				assert.NotNil(t, breaker.cb)
			},
		},
		{
			name: "WithGoBreakerInterval",
			options: []GoBreakerOption{
				WithGoBreakerInterval(1 * time.Minute),
			},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t, 1*time.Minute, breaker.settings.Interval)
				assert.NotNil(t, breaker.cb)
			},
		},
		{
			name: "WithGoBreakerTimeout",
			options: []GoBreakerOption{
				WithGoBreakerTimeout(30 * time.Second),
			},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t, 30*time.Second, breaker.settings.Timeout)
				assert.NotNil(t, breaker.cb)
			},
		},
		{
			name: "WithGoBreakerReadyToTrip",
			options: []GoBreakerOption{
				WithGoBreakerReadyToTrip(readyToTrip),
			},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t,
					reflect.ValueOf(readyToTrip).Pointer(),
					reflect.ValueOf(breaker.settings.ReadyToTrip).Pointer(),
					"readyToTrip is not valid")
				assert.NotNil(t, breaker.cb)
			},
		},
		{
			name: "WithGoBreakerOnStateChange",
			options: []GoBreakerOption{
				WithGoBreakerOnStateChange(onStateChange),
			},
			verify: func(t *testing.T, breaker *GoBreaker) {
				assert.Equal(t,
					reflect.ValueOf(onStateChange).Pointer(),
					reflect.ValueOf(breaker.settings.OnStateChange).Pointer(),
					"onStateChange is not valid")
				assert.NotNil(t, breaker.cb)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			breaker := NewGoBreaker(tt.options...)
			tt.verify(t, breaker)
		})
	}
}
