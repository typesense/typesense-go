package circuit

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v4/typesense/api/circuit/mocks"
	"go.uber.org/mock/gomock"
)

func newHTTPRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	assert.NoError(t, err)
	return req
}

func newHTTPResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
	}
}

func applyFunc(f func() error) error {
	return f()
}

func TestClientDoOnCircuitBreakerErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCircuitBreaker := mocks.NewMockBreaker(ctrl)
	mockRequestDoer := mocks.NewMockHTTPRequestDoer(ctrl)

	mockCircuitBreaker.EXPECT().Execute(gomock.Not(gomock.Nil())).
		Return(errors.New("circuit breaker error"))

	client := NewHTTPClient(
		WithCircuitBreaker(mockCircuitBreaker),
		WithHTTPRequestDoer(mockRequestDoer),
	)

	_, err := client.Do(newHTTPRequest(t))
	assert.Error(t, err)
}

func TestClientDoOnHttpClientErrorReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCircuitBreaker := mocks.NewMockBreaker(ctrl)
	mockRequestDoer := mocks.NewMockHTTPRequestDoer(ctrl)

	gomock.InOrder(
		mockCircuitBreaker.EXPECT().Execute(gomock.Not(gomock.Nil())).
			DoAndReturn(applyFunc),
		mockRequestDoer.EXPECT().Do(newHTTPRequest(t)).
			Return(nil, errors.New("http client error")),
	)

	client := NewHTTPClient(
		WithCircuitBreaker(mockCircuitBreaker),
		WithHTTPRequestDoer(mockRequestDoer),
	)

	_, err := client.Do(newHTTPRequest(t))
	assert.Error(t, err)
}

func TestClientDoReturnsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCircuitBreaker := mocks.NewMockBreaker(ctrl)
	mockRequestDoer := mocks.NewMockHTTPRequestDoer(ctrl)

	gomock.InOrder(
		mockCircuitBreaker.EXPECT().Execute(gomock.Not(gomock.Nil())).
			DoAndReturn(applyFunc),
		mockRequestDoer.EXPECT().Do(newHTTPRequest(t)).
			Return(newHTTPResponse(), nil),
	)

	client := NewHTTPClient(
		WithCircuitBreaker(mockCircuitBreaker),
		WithHTTPRequestDoer(mockRequestDoer),
	)

	response, err := client.Do(newHTTPRequest(t))
	assert.NoError(t, err)
	assert.Equal(t, newHTTPResponse(), response)
}

func TestClientDoOn5xxStatusErrorSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCircuitBreaker := mocks.NewMockBreaker(ctrl)
	mockRequestDoer := mocks.NewMockHTTPRequestDoer(ctrl)

	gomock.InOrder(
		mockCircuitBreaker.EXPECT().Execute(gomock.Not(gomock.Nil())).
			DoAndReturn(applyFunc),
		mockRequestDoer.EXPECT().Do(newHTTPRequest(t)).
			Return(&http.Response{
				StatusCode: http.StatusInternalServerError,
			}, nil),
	)

	client := NewHTTPClient(
		WithCircuitBreaker(mockCircuitBreaker),
		WithHTTPRequestDoer(mockRequestDoer),
	)

	_, err := client.Do(newHTTPRequest(t))
	assert.NoError(t, err)
}
